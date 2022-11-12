package view

import (
	"io"
	"model"
	"net/http"
)

func RenderAllFlashcards(w http.ResponseWriter, cards []model.Flashcard) {
	w.Header().Add("Content-Type", "text/html")

	html := "<html><body>"

	html += `
		<label>Add single card</label>
		<form action="/add" method="POST">
			<input type="text" name="Front"/>
			<input type="text" name="Back"/>
			<input type="submit" value="Add"/>
		</form>
	
		<label>Import cards from CSV</label>
		<form action="/import" method="POST" enctype="multipart/form-data">
			<input type="file" name="fileToUpload" id="fileToUpload"/>
			<input type="submit"/>
		</form>
	`

	for _, card := range cards {
		html += "<p>Front: " + card.Front + ", Back: " + card.Back + "</p>\n"
	}

	html += "</body></html>"

	io.WriteString(w, html)
}
