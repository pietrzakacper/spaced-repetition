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
	CreateMemorizingSession(count int)
	CreateReviewSession(count int)
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
	NextReviewOffset int
	RepetitionCount  int
	EF               float64
}

type Store interface {
	ReadAll() []FlashcardRecord
	Add(record *FlashcardRecord)
	Update(record *FlashcardRecord)
}

type flashcard struct {
	id           string
	front        string
	back         string
	creationDate time.Time
	memorizable  *supermemo.Memorizable
}
