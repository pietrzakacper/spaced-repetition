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
	ShowCards()
	DeleteCard(cardId string) error
	EditCard(cardId, front, back string) error
}

type View interface {
	GoToHome()
	GoToAnswer()
	GoToQuest()
	GoToCards()
	RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int)
	RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int)
	RenderCardAnswer(card *flashcard.DTO, cardNumber int, totalCardsInSession int, answerOptions []int)
	RenderCards(cards []flashcard.DTO)
}

type Persistance interface {
	Create(name string, userId string) Store
}

type Store interface {
	ReadAll() []flashcard.Record
	Add(record *flashcard.Record)
	Update(record *flashcard.Record)
	Find(cardId string) (flashcard.Record, error)
}
