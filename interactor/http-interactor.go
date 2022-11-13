package interactor

import (
	"controller"
	"net/http"
	"view"
)

type HttpInteractor struct {
}

func (HttpInteractor) Start(c controller.Controller) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flashcards, newCardsCount, dueToReviewCount := c.GetAllFlashCards()

		view.RenderAllFlashcards(w, flashcards, newCardsCount, dueToReviewCount)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		r.ParseForm()

		c.AddCard(r.Form.Get("Front"), r.Form.Get("Back"))

		w.Header().Add("Location", "/")
		w.WriteHeader(303)
	})

	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, _, _ := r.FormFile("fileToUpload")

		w.Header().Add("Location", "/")
		w.WriteHeader(303)

		c.ImportCards(file)
	})

	var session *controller.MemorizingSession

	http.HandleFunc("/learnNew", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		session = c.CreateMemorizingSession(10)
		w.Header().Add("Location", "/quest")
		w.WriteHeader(303)
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		session = c.CreateReviewSession(10)
		w.Header().Add("Location", "/quest")
		w.WriteHeader(303)
	})

	http.HandleFunc("/quest", func(w http.ResponseWriter, r *http.Request) {
		cardNumber, totalCardsCount, card := session.GetCurrentQuest()

		if card == nil {
			w.Header().Add("Location", "/")
			w.WriteHeader(303)
			return
		}

		view.RenderCardQuestion(w, card, cardNumber, totalCardsCount)
	})

	http.HandleFunc("/answer", func(w http.ResponseWriter, r *http.Request) {
		cardNumber, totalCardsCount, card := session.GetCurrentQuest()

		view.RenderCardAnswer(w, card, cardNumber, totalCardsCount, session.GetAnswerFeedbackOptions())
	})

	http.HandleFunc("/submitAnswer", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		answer := r.Form.Get("answerFeedback")

		session.SubmitAnswer(answer)

		w.Header().Add("Location", "/quest")
		w.WriteHeader(303)
	})

	http.ListenAndServe(":3000", nil)
}
