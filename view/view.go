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

	html += `
		<form action="/add" method="POST">
			<input type="text" name="Front"/>
			<input type="text" name="Back"/>
			<br/>
			<input type="submit"/>
		</form>
	
	`

	io.WriteString(w, html)
}
