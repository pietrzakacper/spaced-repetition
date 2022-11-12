package controller

import (
	"io"
	"model"
	"parser"
	"persistance"
	"supermemo"
	"time"
)

var store = persistance.Create("db")

func GetAllFlashCards() ([]model.Flashcard, uint) {
	lines := store.Read()

	csvEntries := parser.ParseCSVLines(lines)

	flashcards := make([]model.Flashcard, len(csvEntries))

	newCardsCount := uint(0)

	for i, entry := range csvEntries {
		flashcards[i] = *model.Deserialize(entry)

		if flashcards[i].Metadata.Memorizable.IsNew() {
			newCardsCount++
		}
	}

	return flashcards, newCardsCount
}

func AddCard(front string, back string) {
	card := model.Flashcard{Front: front, Back: back, Metadata: &model.Metadata{
		CreationDate:       time.Now(),
		LastRepetitionDate: time.Now(),
		Memorizable:        supermemo.Create(),
	}}

	entries := card.Serialize()

	store.Add(parser.MakeLine(*entries))
}

func ImportCards(csvStream io.Reader) {
	entriesChan := parser.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		AddCard(entry[0], entry[1])
	}
}

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize *[]model.Flashcard
}

func CreateMemorizingSession(count uint) *MemorizingSession {
	lines := store.Read()

	csvEntries := parser.ParseCSVLines(lines)

	flashcards := make([]model.Flashcard, 0)

	for _, entry := range csvEntries {
		if uint(len(flashcards)) == count {
			break
		}

		card := *model.Deserialize(entry)

		if card.Metadata.Memorizable.IsNew() {
			flashcards = append(flashcards, card)
		}
	}

	return &MemorizingSession{memorizedCount: 0, cardsToMemorize: &flashcards}
}

func (m *MemorizingSession) GetCurrentQuest() (int, int, *model.Flashcard) {
	return m.memorizedCount + 1, len(*m.cardsToMemorize), &(*m.cardsToMemorize)[m.memorizedCount]
}

var answerFeedback = map[string]int{
	"Complete Blackout":        0,
	"Slipped my mind":          1,
	"Ah shit, I knew it!":      2,
	"Barely correct bro":       3,
	"I remembered correctly:)": 4,
	"Too easy!":                5,
}

func (m *MemorizingSession) GetAnswerFeedbackOptions() []string {
	keys := make([]string, 0, len(answerFeedback))
	for k := range answerFeedback {
		keys = append(keys, k)
	}

	return keys
}

func (m *MemorizingSession) SubmitAnswer(answer string) {
	_, _, card := m.GetCurrentQuest()

	qualityOfResponse := supermemo.QualityOfResponse(answerFeedback[answer])

	card.Metadata.Memorizable.SubmitRepetition(qualityOfResponse)

	// @TODO persist card

	m.memorizedCount++
}
