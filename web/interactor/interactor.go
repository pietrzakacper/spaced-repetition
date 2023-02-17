package interactor

import (
	"controller"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"user"
	"web/view"
)

type Interactor interface {
	Start()
}

type HttpInteractor struct {
	view           *view.HttpView
	sessionFactory *user.UserSessionFactory

	sessions map[string]*controller.FlashcardsController
}

func CreateHttpInteractor(view *view.HttpView, sessionFactory *user.UserSessionFactory) HttpInteractor {
	return HttpInteractor{
		view,
		sessionFactory,
		map[string]*controller.FlashcardsController{},
	}
}

type EditCardPayload struct {
	Front string
	Back  string
}

// @TODO think about making a server middleware
func (i HttpInteractor) authenticateUser(w http.ResponseWriter, r *http.Request) (*controller.FlashcardsController, error) {
	authToken := r.Header.Get("Cookie")

	// @TODO investigate double call
	fmt.Printf("TOKEN: %v\n", authToken)
	// @TODO authenticate the user

	if authToken == "" {
		w.WriteHeader(401)
		// @TODO redirect to login
		return nil, errors.New("Unauthorized")
	}

	c := i.sessions[authToken]

	if c == nil {
		c = i.sessionFactory.Create(authToken)
		i.sessions[authToken] = c
	}

	return c, nil
}

func (i HttpInteractor) Start() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowHome()
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseForm()

			c.AddCard(r.Form.Get("Front"), r.Form.Get("Back"))
		}
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {

			r.ParseForm()

			c.AddCard(r.Form.Get("Front"), r.Form.Get("Back"))
		}
	})

	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseMultipartForm(10 << 20)

			file, _, err := r.FormFile("fileToUpload")

			if err != nil {
				i.view.GoToHome()
				return
			}

			c.ImportCards(file)
		}
	})

	http.HandleFunc("/learnNew", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateMemorizingSession()
		}
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		if c, err := i.authenticateUser(w, r); err == nil {
			c.CreateReviewSession()
		}
	})

	http.HandleFunc("/quest", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowQuest()
		}
	})

	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if c, err := i.authenticateUser(w, r); err == nil {
			c.ShowAnswer()
		}
	})

	http.HandleFunc("/submitAnswer", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if c, err := i.authenticateUser(w, r); err == nil {
			r.ParseForm()
			answerStr := r.Form.Get("answerFeedback")

			answer, _ := strconv.ParseInt(answerStr, 10, 32)

			c.SubmitAnswer(int(answer))
		}
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

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

			i.view.SetRequestContext(w)

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
		i.view.SetRequestContext(w)

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
