package view

import (
	"fmt"
	"io"
	"model"
	"net/http"
	"strconv"
)

func RenderAllFlashcards(w http.ResponseWriter, cards []model.Flashcard, newCardsCount uint, dueToReviewCount uint) {
	w.Header().Add("Content-Type", "text/html")

	html := "<html><body>"

	html += fmt.Sprintf(`
		<form action="/review" method="POST">
			<b>Cards to review</b>: %d
			<input type="submit" value="Review"/>
		</form>
		<form action="/learnNew" method="POST">
			<b>New cards</b>: %d
			<input type="submit" value="Memorize"/>
		</form>
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
	`, dueToReviewCount, newCardsCount)

	html += "<b>All cards:</b> " + strconv.FormatInt(int64(len(cards)), 10)

	for _, card := range cards {
		html += "<p>Front: " + card.Front + ", Back: " + card.Back + "</p>\n"
	}

	html += "</body></html>"

	io.WriteString(w, html)
}

func RenderCardQuestion(w http.ResponseWriter, card *model.Flashcard, cardNumber int, totalCardsInSession int) {
	w.Header().Add("Content-Type", "text/html")

	html := "<html><body>"

	html += "<b>Card: </b>" +
		strconv.FormatInt(int64(cardNumber), 10) +
		"/" + strconv.FormatInt(int64(totalCardsInSession), 10)

	html += "<br/>"

	html += "<h2>" + card.Front + "</h2>"

	html += `
		<form action="/answer" method="POST">
			<label>Show Answer</label>
			<input type="submit"/>
		</form>
	`

	html += "</body></html>"

	io.WriteString(w, html)
}

// @TODO use HTML templates
func RenderCardAnswer(w http.ResponseWriter, card *model.Flashcard, cardNumber int, totalCardsInSession int, answerOptions []string) {
	w.Header().Add("Content-Type", "text/html")

	html := "<html><body>"

	html += "<b>Card: </b>" +
		strconv.FormatInt(int64(cardNumber), 10) +
		"/" + strconv.FormatInt(int64(totalCardsInSession), 10)

	html += "<br/>"

	html += "<h2>" + card.Back + "</h2>"

	html += `
		<form action="/submitAnswer" method="POST">
			<label>How's your memory?</label><br/>`

	for _, option := range answerOptions {
		html += fmt.Sprintf(
			`<input type="submit" value="%s" name="answerFeedback"/>`,
			option,
		)
	}

	html += "</form>"

	html += "</body></html>"

	io.WriteString(w, html)
}
