package main

import (
	"controller"
	"io"
	"net/http"
	"strings"
	"view"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flashcards, newCardsCount, dueToReviewCount := controller.GetAllFlashCards()

		view.RenderAllFlashcards(w, flashcards, newCardsCount, dueToReviewCount)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		body, _ := io.ReadAll(r.Body)
		entries := strings.Split(string(body), "&")

		var (
			front string
			back  string
		)

		// @TODO use form parsing
		for _, entry := range entries {
			kvPair := strings.Split(entry, "=")

			if kvPair[0] == "Front" {
				front = kvPair[1]
			}

			if kvPair[0] == "Back" {
				back = kvPair[1]
			}
		}

		controller.AddCard(front, back)

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

		controller.ImportCards(file)
	})

	var session *controller.MemorizingSession

	http.HandleFunc("/learnNew", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		session = controller.CreateMemorizingSession(10)
		w.Header().Add("Location", "/quest")
		w.WriteHeader(303)
	})

	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		session = controller.CreateReviewSession(10)
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

		// @TODO move views to controllers
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
