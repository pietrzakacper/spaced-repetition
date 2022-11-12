package main

import (
	controller "controller"
	"io"
	"model"
	"net/http"
	"strings"
	"view"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flashcards := controller.GetAllFlashCards()

		view.RenderAllFlashcards(w, flashcards)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		body, _ := io.ReadAll(r.Body)
		entries := strings.Split(string(body), "&")

		card := model.Flashcard{}

		for _, entry := range entries {
			kvPair := strings.Split(entry, "=")

			if kvPair[0] == "Front" {
				card.Front = kvPair[1]
			}

			if kvPair[0] == "Back" {
				card.Back = kvPair[1]
			}
		}

		controller.AddCard(&card)

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

	http.ListenAndServe(":3000", nil)
}
