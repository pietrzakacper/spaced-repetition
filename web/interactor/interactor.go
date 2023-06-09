package interactor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"persistance"
	"strconv"
	"strings"
	"sync"
	"user"
	"web/view"

	"google.golang.org/api/oauth2/v1"
)

type Interactor interface {
	Start()
}

type UserSessionWithLock struct {
	session *user.UserSession
	lock    *sync.Mutex
}

type HttpInteractor struct {
	sessionFactory *user.UserSessionFactory

	sessions           map[string]*UserSessionWithLock
	sessionsCommonLock sync.Mutex
}

func CreateHttpInteractor(sessionFactory *user.UserSessionFactory) HttpInteractor {
	return HttpInteractor{
		sessionFactory,
		map[string]*UserSessionWithLock{},
		sync.Mutex{},
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
func (i HttpInteractor) authenticateUser(w http.ResponseWriter, r *http.Request) (*user.UserSession, error) {
	if os.Getenv("LOCAL_DEV") == "true" {
		if c := i.sessions["LOCAL"]; c != nil {
			c.session.SetRequestContext(w)
			return c.session, nil
		}
		c := &UserSessionWithLock{
			session: i.sessionFactory.Create(user.UserContext{
				Email: "kacpietrzak@gmail.com",
			}),
			lock: &sync.Mutex{},
		}

		i.sessions["LOCAL"] = c

		c.session.SetRequestContext(w)

		return c.session, nil
	}

	cookies := r.Cookies()

	cookieMap := make(map[string]string, 1)

	for _, str := range cookies {
		// ignore cookie attributes
		cookieMap[str.Name] = strings.Split(str.Value, "; ")[0]
	}

	authToken := cookieMap["sessionToken"]
	w.Header().Add("Referrer-Policy", "no-referrer-when-downgrade")

	if authToken == "" {
		// user logged out
		anonymousView := view.CreateHttpView("")
		anonymousView.SetRequestContext(w)
		anonymousView.RenderLogin()

		return nil, errors.New("Unauthorized")
	}

	i.sessionsCommonLock.Lock()
	sl := i.sessions[authToken]

	if sl == nil {
		sl = &UserSessionWithLock{
			session: nil,
			lock:    &sync.Mutex{},
		}
		i.sessions[authToken] = sl
	}
	i.sessionsCommonLock.Unlock()

	sl.lock.Lock()
	defer sl.lock.Unlock()

	// try finding session in the memory
	if sl.session != nil {
		sl.session.SetRequestContext(w)
		return sl.session, nil
	}

	// try finding the session in DB
	db := (&persistance.PostgresPersistance{}).Create("anonymous@spaced.sh")
	userIdFromDb, err := db.FindUserIdByToken(authToken)

	if err == nil {
		sl.session = i.sessionFactory.Create(user.UserContext{
			Email: userIdFromDb,
		})
		sl.session.SetRequestContext(w)

		return sl.session, nil
	}

	// if there was no saved session for the user, we need to verify the token and save it
	tokenInfo, err := verifyIdToken(authToken)

	if err != nil {
		fmt.Printf("/login error: %v\n", err)

		// user logged out
		anonymousView := view.CreateHttpView("")
		anonymousView.SetRequestContext(w)
		anonymousView.RenderLogin()

		return nil, errors.New("Unauthorized")
	}

	sl.session = i.sessionFactory.Create(user.UserContext{
		Email: tokenInfo.Email,
	})

	db.UpsertSession(authToken, tokenInfo.Email)

	sl.session.SetRequestContext(w)

	return sl.session, nil
}

func (i HttpInteractor) Start() {
	// @TODO invalidate cache on deployment
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			// user logged in
			c.ShowHome()
			return
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		const max_request_size = 4096
		r.Body = http.MaxBytesReader(w, r.Body, max_request_size)
		r.ParseForm()
		token := r.Form.Get("credential")
		_, err := verifyIdToken(token)

		if err == nil {
			// keep logged in for 90 days
			w.Header().Add("Set-Cookie", "sessionToken="+token+"; Max-Age=7776000; Path=/; SameSite=Strict; HttpOnly; Secure;")
		}

		w.Header().Add("Location", "/")
		w.WriteHeader(303)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			const max_request_size = 1024
			r.Body = http.MaxBytesReader(w, r.Body, max_request_size)
			r.ParseForm()

			front, back := r.Form.Get("front"), r.Form.Get("back")

			if front == "" || back == "" {
				c.GoToHome()
				return
			}

			c.AddCard(front, back)
		}
	})

	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			const max_request_size = 2 * 1024 * 1024
			r.Body = http.MaxBytesReader(w, r.Body, max_request_size)
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
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateMemorizingSession()
		}
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateReviewSession()
		}
	})

	http.HandleFunc("/quest", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {

			c.ShowQuest(view.DecodeCookiesToMemorizingSession(r))
		}
	})

	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {

			c.ShowAnswer(view.DecodeCookiesToMemorizingSession(r))
		}
	})

	http.HandleFunc("/submitAnswer", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			const max_request_size = 1024
			r.Body = http.MaxBytesReader(w, r.Body, max_request_size)
			r.ParseForm()
			answerStr := r.Form.Get("answerFeedback")

			answer, _ := strconv.ParseInt(answerStr, 10, 32)

			c.SubmitAnswer(view.DecodeCookiesToMemorizingSession(r), int(answer))
		}
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowCards()
		}
	})

	http.HandleFunc("/delete-card", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
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
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			const max_request_size = 1024
			r.Body = http.MaxBytesReader(w, r.Body, max_request_size)
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
