package main

import (
	controller "controller"
	"net/http"
	"view"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flashcards := controller.GetAllFlashCards()

		view.RenderAllFlashcards(w, flashcards)
	})

	http.ListenAndServe(":3000", nil)
}
