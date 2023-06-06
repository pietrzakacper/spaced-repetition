package controller

import (
	"flashcard"
	"io"
)

type Controller interface {
	ShowHome()
	CreateMemorizingSession()
	CreateReviewSession()
	ShowQuest(sessionDTO *flashcard.MemorizingSessionDTO)
	ShowAnswer(sessionDTO *flashcard.MemorizingSessionDTO)
	SubmitAnswer(sessionDTO *flashcard.MemorizingSessionDTO, answer int)
	AddCard(front string, back string)
	ImportCards(csvStream io.Reader)
	ShowCards()
	DeleteCard(cardId string) error
	EditCard(cardId, front, back string) error
}

type View interface {
	UpdateClientSession(session *flashcard.MemorizingSessionDTO)
	GoToHome()
	GoToAnswer()
	GoToQuest()
	OnSessionFinished()
	GoToCards()
	RenderHome(cards []flashcard.DTO, newCardsCount int, dueToReviewCount int)
	RenderCardQuestion(card *flashcard.DTO, cardNumber int, totalCardsInSession int, extraRound bool)
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
