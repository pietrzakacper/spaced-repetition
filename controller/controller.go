package controller

import (
	"io"
	"model"
)

type Controller interface {
	GetAllFlashCards() ([]model.Flashcard, uint, uint)
	AddCard(front string, back string)
	ImportCards(csvStream io.Reader)
	CreateMemorizingSession(count uint) *MemorizingSession
	CreateReviewSession(count uint) *MemorizingSession
}
