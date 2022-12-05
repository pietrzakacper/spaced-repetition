package view

import (
	"flashcard"
	t "html/template"
	"net/http"
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

type HomeData struct {
	DueToReviewCount int
	NewCardsCount    int
	AllCardsCount    int
	Cards            []cardView
}

type cardView struct {
	Front string
	Back  string
	Kind  byte
}

func (v *HttpView) RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int) {
	template := t.Must(t.ParseFiles("templates/home.html"))

	recentCards := make([]cardView, 0)

	kind := byte(0)

	for i := len(cards) - 1; i >= 0 && len(recentCards) < 5; i-- {
		dto := cards[i]
		for _, char := range dto.Id {
			kind += byte(char)
		}
		kind := kind % 4

		cardView := cardView{Front: dto.Front, Back: dto.Back, Kind: kind}

		recentCards = append(recentCards, cardView)
	}

	data := HomeData{
		DueToReviewCount: dueToReviewCount,
		NewCardsCount:    newCardsCount,
		AllCardsCount:    len(cards),
		Cards:            recentCards,
	}

	template.Execute(v.w, data)
}

type QuestData struct {
	CardNumber          int
	TotalCardsInSession int
	Front               string
}

func (v *HttpView) RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int) {
	template := t.Must(t.ParseFiles("templates/quest.html"))

	data := QuestData{
		CardNumber:          cardNumber,
		TotalCardsInSession: totalCardsInSession,
		Front:               card.Front,
	}

	template.Execute(v.w, data)
}

var answerLabels = map[int]string{
	0: "Blackout",
	1: "Slipped my mind",
	2: "Almost got it",
	3: "Barely correct",
	4: "I remembered correctly:)",
	5: "Too easy!",
}

type Answer struct {
	Value int
	Label string
}
type AnswerData struct {
	CardNumber          int
	TotalCardsInSession int
	Back                string
	Answers             []Answer
}

func (v *HttpView) RenderCardAnswer(card *flashcard.DTO, cardNumber int, totalCardsInSession int, answerOptions []int) {
	template := t.Must(t.ParseFiles("templates/answer.html"))

	Answers := make([]Answer, len(answerOptions))

	for i, value := range answerOptions {
		Answers[i] = Answer{Value: value, Label: answerLabels[value]}
	}

	data := AnswerData{
		CardNumber:          cardNumber,
		TotalCardsInSession: totalCardsInSession,
		Back:                card.Back,
		Answers:             Answers,
	}

	template.Execute(v.w, data)
}
