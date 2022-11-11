package controller

import (
	"model"
	"parser"
	p "persistance"
)

var store = "flashcards.csv"

func GetAllFlashCards() []model.Flashcard {
	persistance := p.Create(store)

	lines := persistance.Read()

	csvEntries := parser.ParseCSVLines(lines)

	flashcards := make([]model.Flashcard, len(csvEntries))

	for i, entry := range csvEntries {
		flashcards[i] = model.Flashcard{Front: entry.First, Back: entry.Second}
	}

	return flashcards
}
