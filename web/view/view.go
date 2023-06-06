package view

import (
	"flashcard"
	t "html/template"
	"net/http"
	"os"
)

type HttpView struct {
	w         http.ResponseWriter
	userEmail string
}

func CreateHttpView(userEmail string) *HttpView {
	return &HttpView{w: nil, userEmail: userEmail}
}

func (v *HttpView) SetRequestContext(w http.ResponseWriter) {
	v.w = w
}

func (v *HttpView) GoToHome() {
	v.w.Header().Add("Location", "/")
	v.w.WriteHeader(303)
}

func (v *HttpView) GoToCards() {
	v.w.Header().Add("Location", "/cards")
	v.w.WriteHeader(303)
}

func (v *HttpView) UpdateClientSession(session *flashcard.MemorizingSessionDTO) {
	cookies := EncodeMemorizingSessionToCookies(session)
	for cookieName, cookieValue := range cookies {
		v.w.Header().Add("Set-Cookie", cookieName+"="+cookieValue+" Max-Age=1800 Path=/")
	}
}

func (v *HttpView) OnSessionFinished() {
	v.w.Header().Add("Set-Cookie", "sessionFinished=true Path=/")
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
	// @TODO set it in all places
	Id    string
	Front string
	Back  string
	Kind  byte
}

func (v *HttpView) RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int) {
	template := t.Must(t.ParseFiles("templates/home.html"))

	recentCards := make([]cardView, 0)

	for i := len(cards) - 1; i >= 0 && len(recentCards) < 5; i-- {
		dto := cards[i]

		cardView := cardView{Front: dto.Front, Back: dto.Back, Kind: getCardKind(dto.Id)}

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
	ExtraRound          bool
	Kind                byte
}

func (v *HttpView) RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int, extraRound bool) {
	template := t.Must(t.ParseFiles("templates/quest.html"))

	data := QuestData{
		CardNumber:          cardNumber,
		TotalCardsInSession: totalCardsInSession,
		Front:               card.Front,
		ExtraRound:          extraRound,
		Kind:                getCardKind(card.Id),
	}

	template.Execute(v.w, data)
}

var answerLabels = map[int]string{
	0: "blackout",
	1: "slipped my mind",
	2: "almost got it",
	3: "barely correct",
	4: "correct",
	5: "too easy!",
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
	Kind                byte
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
		Kind:                getCardKind(card.Id),
	}

	template.Execute(v.w, data)
}

type CardsData struct {
	Cards []cardView
}

func (v *HttpView) RenderCards(cards []flashcard.DTO) {
	template := t.Must(t.ParseFiles("templates/cards.html"))

	cardViews := make([]cardView, len(cards))

	for index, dto := range cards {
		cardViews[index] = cardView{
			Id: dto.Id, Front: dto.Front, Back: dto.Back, Kind: getCardKind(dto.Id),
		}
	}

	template.Execute(v.w, CardsData{Cards: cardViews})
}

func getCardKind(cardId string) byte {
	// integer between 0-3
	kind := byte(0)

	// compute short hash from id
	for _, char := range cardId {
		kind = (kind + byte(char)) % 4
	}

	return kind
}

type LoginData struct {
	AppUrl string
}

func (v *HttpView) RenderLogin() {
	template := t.Must(t.ParseFiles("templates/login.html"))

	appUrl := os.Getenv("APP_URL")

	if appUrl == "" {
		appUrl = "http://localhost:3000"
	}

	template.Execute(v.w, LoginData{AppUrl: appUrl})
}
