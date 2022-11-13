package controller

import (
	"io"
)

type Controller interface {
	ShowHome()
	ShowQuest()
	ShowAnswer()
	AddCard(front string, back string)
	ImportCards(csvStream io.Reader)
	CreateMemorizingSession(count uint)
	CreateReviewSession(count uint)
	SubmitAnswer(answer string)
}

type FlashcardDTO struct {
	Front string
	Back  string
}

type View interface {
	GoToHome()
	GoToAnswer()
	GoToQuest()
	RenderHome(cards []FlashcardDTO, newCardsCount uint, dueToReviewCount uint)
	RenderCardQuestion(card *FlashcardDTO, cardNumber int, totalCardsInSession int)
	RenderCardAnswer(card *FlashcardDTO, cardNumber int, totalCardsInSession int, answerOptions []string)
}
