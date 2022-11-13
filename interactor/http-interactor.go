package interactor

import (
	"controller"
	"net/http"
	"strconv"
	"view"
)

type HttpInteractor struct {
	view *view.HttpView
}

func CreateHttpInteractor(view *view.HttpView) HttpInteractor {
	return HttpInteractor{view: view}
}

func (i HttpInteractor) Start(c controller.Controller) {
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

		file, _, _ := r.FormFile("fileToUpload")

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

	http.ListenAndServe(":3000", nil)
}
