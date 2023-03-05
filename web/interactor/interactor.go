package interactor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"user"
	"web/view"

	"google.golang.org/api/oauth2/v1"
)

type Interactor interface {
	Start()
}

type HttpInteractor struct {
	sessionFactory *user.UserSessionFactory

	sessions map[string]*user.UserSession
}

func CreateHttpInteractor(sessionFactory *user.UserSessionFactory) HttpInteractor {
	return HttpInteractor{
		// @TODO create a view per session
		sessionFactory,
		map[string]*user.UserSession{},
	}
}

type EditCardPayload struct {
	Front string
	Back  string
}

func verifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(http.DefaultClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

// @TODO think about making a server middleware
func (i HttpInteractor) authenticateUser(w http.ResponseWriter, r *http.Request, keepLocationOnAuthFail ...bool) (*user.UserSession, error) {
	cookies := r.Cookies()

	cookieMap := make(map[string]string, 1)

	for _, str := range cookies {
		cookieMap[str.Name] = str.Value
	}

	authToken := cookieMap["sessionToken"]
	w.Header().Add("Referrer-Policy", "no-referrer-when-downgrade")

	c := i.sessions[authToken]

	if c != nil {
		fmt.Println("Found session")
		c.SetRequestContext(w)
		return c, nil
	}

	// @TODO investigate double call
	if authToken == "" {
		fmt.Printf("/login error: No sessionToken")

		if len(keepLocationOnAuthFail) == 0 || !keepLocationOnAuthFail[0] {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
		}

		return nil, errors.New("Unauthorized")
	}

	tokenInfo, err := verifyIdToken(authToken)

	if err != nil {
		fmt.Printf("/login error: %v\n", err)

		if len(keepLocationOnAuthFail) == 0 || !keepLocationOnAuthFail[0] {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
		}

		return nil, errors.New("Unauthorized")
	}

	c = i.sessionFactory.Create(user.UserContext{
		Id: tokenInfo.UserId, Email: tokenInfo.Email,
	})
	i.sessions[authToken] = c

	c.SetRequestContext(w)

	return c, nil
}

func (i HttpInteractor) Start() {
	// @TODO invalidate cache on deployment
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r, true); err == nil {
			// user logged in
			c.ShowHome()
			return
		}

		// user logged out
		anonymousView := view.CreateHttpView("")
		anonymousView.SetRequestContext(w)
		anonymousView.RenderLogin()
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		token := r.Form.Get("credential")
		_, err := verifyIdToken(token)

		if err == nil {
			w.Header().Add("Set-Cookie", "sessionToken="+token)
		}

		w.Header().Add("Location", "/")
		w.WriteHeader(303)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseForm()

			c.AddCard(r.Form.Get("Front"), r.Form.Get("Back"))
		}
	})

	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseMultipartForm(10 << 20)

			file, _, err := r.FormFile("fileToUpload")

			if err != nil {
				c.GoToHome()
				return
			}

			c.ImportCards(file)
		}
	})

	http.HandleFunc("/learnNew", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateMemorizingSession()
		}
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateReviewSession()
		}
	})

	http.HandleFunc("/quest", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowQuest()
		}
	})

	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowAnswer()
		}
	})

	http.HandleFunc("/submitAnswer", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseForm()
			answerStr := r.Form.Get("answerFeedback")

			answer, _ := strconv.ParseInt(answerStr, 10, 32)

			c.SubmitAnswer(int(answer))
		}
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowCards()
		}
	})

	http.HandleFunc("/delete-card", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			cardId := r.URL.Query().Get("id")

			err := c.DeleteCard(cardId)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	})

	http.HandleFunc("/edit-card", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			cardId := r.URL.Query().Get("id")

			payload := EditCardPayload{}
			err := json.NewDecoder(r.Body).Decode(&payload)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = c.EditCard(cardId, payload.Front, payload.Back)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}
