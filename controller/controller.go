package controller

import (
	"io"
	"model"
	"parser"
	"persistance"
	"supermemo"
)

var store = persistance.Create("db")

func GetAllFlashCards() []model.Flashcard {
	lines := store.Read()

	csvEntries := parser.ParseCSVLines(lines)

	flashcards := make([]model.Flashcard, len(csvEntries))

	for i, entry := range csvEntries {
		flashcards[i] = model.Flashcard{Front: entry[0], Back: entry[1]}
	}

	return flashcards
}

func AddCard(card *model.Flashcard) {
	line := card.Front + "," + card.Back + "," + supermemo.Create().Serialize()
	store.Add(line)
}

func ImportCards(csvStream io.Reader) {
	entriesChan := parser.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		AddCard(&model.Flashcard{Front: entry[0], Back: entry[1]})
	}
}
