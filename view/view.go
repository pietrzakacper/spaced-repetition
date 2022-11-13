package view

import (
	"flashcard"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type HttpView struct {
	w http.ResponseWriter
}

func (v *HttpView) SetRequestContext(w http.ResponseWriter) {
	v.w = w
}

func (v *HttpView) GoToHome() {
	v.w.Header().Add("Location", "/")
	v.w.WriteHeader(303)
}

func (v *HttpView) GoToQuest() {
	v.w.Header().Add("Location", "/quest")
	v.w.WriteHeader(303)
}

func (v *HttpView) GoToAnswer() {
	v.w.Header().Add("Location", "/answer")
	v.w.WriteHeader(303)
}

func (v *HttpView) RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int) {
	v.w.Header().Add("Content-Type", "text/html")

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

	io.WriteString(v.w, html)
}

func (v *HttpView) RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int) {
	v.w.Header().Add("Content-Type", "text/html")

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

	io.WriteString(v.w, html)
}

var answerFeedback = map[int]string{
	0: "Complete Blackout",
	1: "Slipped my mind",
	2: "Ah shit, I knew it!",
	3: "Barely correct bro",
	4: "I remembered correctly:)",
	5: "Too easy!",
}

func (v *HttpView) RenderCardAnswer(card *flashcard.DTO, cardNumber int, totalCardsInSession int, answerOptions []int) {
	v.w.Header().Add("Content-Type", "text/html")

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
			`<label for="%d">%s</label>
			<input type="submit" id="%d" value="%d" name="answerFeedback"/>`,
			option,
			answerFeedback[option],
			option,
			option,
		)
	}

	html += "</form>"

	html += "</body></html>"

	io.WriteString(v.w, html)
}
