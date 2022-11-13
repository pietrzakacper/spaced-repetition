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

func GetAllFlashCards() ([]model.Flashcard, uint, uint) {
	records := store.Read()

	flashcards := make([]model.Flashcard, len(records))

	newCardsCount := uint(0)
	dueToReviewCount := uint(0)

	for i, r := range records {
		entry := parser.ParseLine(r.Data)
		flashcards[i] = *model.Deserialize(entry)

		if flashcards[i].Metadata.Memorizable.IsNew() {
			newCardsCount++
		}

		if dueToReview(&flashcards[i]) {
			dueToReviewCount++
		}
	}

	return flashcards, newCardsCount, dueToReviewCount
}

func dueToReview(card *model.Flashcard) bool {
	if card.Metadata.Memorizable.IsNew() {
		return false
	}

	nowInSeconds := time.Now().Unix()

	lastRepInSeconds := card.Metadata.LastRepetitionDate.Unix()

	// convert to seconds
	offsetInSeconds := int64(card.Metadata.Memorizable.GetNextRepetitionDaysOffset() * 24 * 60 * 60)

	return (lastRepInSeconds + offsetInSeconds) <= nowInSeconds
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
	cardsToMemorize *[]FlashcardRecord
}

type FlashcardRecord struct {
	flashcard *model.Flashcard
	record    *persistance.Record
}

func CreateMemorizingSession(count uint) *MemorizingSession {
	records := store.Read()

	flashcards := make([]FlashcardRecord, 0)

	for _, record := range records {
		if uint(len(flashcards)) == count {
			break
		}

		entry := parser.ParseLine(record.Data)

		card := *model.Deserialize(entry)

		if card.Metadata.Memorizable.IsNew() {
			flashcards = append(flashcards, FlashcardRecord{
				flashcard: &card,
				record:    record,
			})
		}
	}

	return &MemorizingSession{memorizedCount: 0, cardsToMemorize: &flashcards}
}

func CreateReviewSession(count uint) *MemorizingSession {
	// @TODO handle failed answers
	records := store.Read()

	flashcards := make([]FlashcardRecord, 0)

	for _, record := range records {
		if uint(len(flashcards)) == count {
			break
		}

		entry := parser.ParseLine(record.Data)

		card := *model.Deserialize(entry)

		if dueToReview(&card) {
			flashcards = append(flashcards, FlashcardRecord{
				flashcard: &card,
				record:    record,
			})
		}
	}

	return &MemorizingSession{memorizedCount: 0, cardsToMemorize: &flashcards}
}

func (m *MemorizingSession) GetCurrentQuest() (int, int, *model.Flashcard) {
	if m.memorizedCount == len(*m.cardsToMemorize) {
		return 0, 0, nil
	}

	return m.memorizedCount + 1, len(*m.cardsToMemorize), (*m.cardsToMemorize)[m.memorizedCount].flashcard
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
	card := (*m.cardsToMemorize)[m.memorizedCount]

	qualityOfResponse := supermemo.QualityOfResponse(answerFeedback[answer])

	card.flashcard.Metadata.Memorizable.SubmitRepetition(qualityOfResponse)
	card.flashcard.Metadata.LastRepetitionDate = time.Now()

	card.record.Data = parser.MakeLine(*card.flashcard.Serialize())
	card.record.Save()

	m.memorizedCount++
}
