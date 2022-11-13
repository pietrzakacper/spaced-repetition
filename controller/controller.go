package controller

import (
	"io"
	"time"
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

type FlashcardDTO struct {
	Front string
	Back  string
}

type View interface {
	GoToHome()
	GoToAnswer()
	GoToQuest()
	RenderHome(cards []FlashcardDTO, newCardsCount int, dueToReviewCount int)
	RenderCardQuestion(card *FlashcardDTO, cardNumber int, totalCardsInSession int)
	RenderCardAnswer(card *FlashcardDTO, cardNumber int, totalCardsInSession int, answerOptions []int)
}

type Persistance interface {
	Create(name string) Store
}

type FlashcardRecord struct {
	Id               string
	Front            string
	Back             string
	CreationDate     time.Time
	LastReviewDate   time.Time
	NextReviewOffset int
	RepetitionCount  int
	EF               float64
}

type Store interface {
	ReadAll() []FlashcardRecord
	Add(record *FlashcardRecord)
	Update(record *FlashcardRecord)
}
