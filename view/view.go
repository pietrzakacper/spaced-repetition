package view

import (
	"io"
	"model"
	"net/http"
)

func RenderAllFlashcards(w http.ResponseWriter, cards []model.Flashcard) {
	html := ""

	for _, card := range cards {
		html += "<p>Front: " + card.Front + ", Back: " + card.Back + "</p>\n"
	}

	io.WriteString(w, html)
}
