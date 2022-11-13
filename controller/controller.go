package controller

import (
	"flashcard"
	"io"
)

type Controller interface {
	ShowHome()
	ShowQuest()
	ShowAnswer()
	AddCard(front string, back string)
	ImportCards(csvStream io.Reader)
	CreateMemorizingSession()
	CreateReviewSession()
	SubmitAnswer(answer int)
}

type View interface {
	GoToHome()
	GoToAnswer()
	GoToQuest()
	RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int)
	RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int)
	RenderCardAnswer(card *flashcard.DTO, cardNumber int, totalCardsInSession int, answerOptions []int)
}

type Persistance interface {
	Create(name string) Store
}

type Store interface {
	ReadAll() []flashcard.Record
	Add(record *flashcard.Record)
	Update(record *flashcard.Record)
}
