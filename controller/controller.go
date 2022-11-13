package controller

import (
	"io"
	"supermemo"
	"time"
)

type Controller interface {
	ShowHome()
	ShowQuest()
	ShowAnswer()
	AddCard(front string, back string)
	ImportCards(csvStream io.Reader)
	CreateMemorizingSession(count int64)
	CreateReviewSession(count int64)
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
	RenderHome(cards []FlashcardDTO, newCardsCount int, dueToReviewCount int)
	RenderCardQuestion(card *FlashcardDTO, cardNumber int, totalCardsInSession int)
	RenderCardAnswer(card *FlashcardDTO, cardNumber int, totalCardsInSession int, answerOptions []string)
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
	NextReviewOffset supermemo.Days
	RepetitionCount  int64
	EF               float64
}

type Store interface {
	ReadAll() []FlashcardRecord
	Add(record *FlashcardRecord)
	Update(record *FlashcardRecord)
}
