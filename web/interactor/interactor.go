package interactor

import (
	"controller"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"web/view"
)

type Interactor interface {
	Start(c controller.Controller)
}

type HttpInteractor struct {
	view *view.HttpView
}

func CreateHttpInteractor(view *view.HttpView) HttpInteractor {
	return HttpInteractor{view: view}
}

func (i HttpInteractor) Start(c controller.Controller) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		c.ShowHome()
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		r.ParseForm()

		c.AddCard(r.Form.Get("Front"), r.Form.Get("Back"))
	})

	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("fileToUpload")

		if err != nil {
			i.view.GoToHome()
			return
		}

		c.ImportCards(file)
	})

	http.HandleFunc("/learnNew", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		c.CreateMemorizingSession()
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		if r.Method != "POST" {
			return
		}

		c.CreateReviewSession()
	})

	http.HandleFunc("/quest", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		c.ShowQuest()
	})

	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		c.ShowAnswer()
	})

	http.HandleFunc("/submitAnswer", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		r.ParseForm()
		answerStr := r.Form.Get("answerFeedback")

		answer, _ := strconv.ParseInt(answerStr, 10, 32)

		c.SubmitAnswer(int(answer))
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		i.view.SetRequestContext(w)

		c.ShowCards()
	})

	http.HandleFunc("/delete-card", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		cardId := r.URL.Query().Get("id")

		i.view.SetRequestContext(w)

		c.DeleteCard(cardId)
	})

	http.HandleFunc("/edit-card", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		i.view.SetRequestContext(w)

		cardId := r.URL.Query().Get("id")

		r.ParseForm()

		c.EditCard(cardId, r.Form.Get("Front"), r.Form.Get("Back"))
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}
