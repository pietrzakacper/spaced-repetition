package controller

import (
	"model"
	"parser"
	p "persistance"
)

var store = p.Create("flashcards.csv")

func GetAllFlashCards() []model.Flashcard {
	lines := store.Read()

	csvEntries := parser.ParseCSVLines(lines)

	flashcards := make([]model.Flashcard, len(csvEntries))

	for i, entry := range csvEntries {
		flashcards[i] = model.Flashcard{Front: entry.First, Back: entry.Second}
	}

	return flashcards
}

func AddCard(card *model.Flashcard) {
	line := card.Front + "," + card.Back
	store.Add(line)
}
