package controller

import (
	"io"
	"model"
	"parser"
	"persistance"
	"supermemo"
	"time"
)

type FlashcardsController struct {
	view    View
	store   persistance.Persistance
	session *MemorizingSession
}

func CreateFlashcardsController(view View) *FlashcardsController {
	return &FlashcardsController{
		view:  view,
		store: persistance.Create("db"),
	}
}

func (c *FlashcardsController) ShowHome() {
	records := c.store.Read()

	flashcardDTOs := make([]FlashcardDTO, len(records))

	newCardsCount := uint(0)
	dueToReviewCount := uint(0)

	for i, r := range records {
		entry := parser.ParseLine(r.Data)
		flashcard := *model.Deserialize(entry)

		if flashcard.Metadata.Memorizable.IsNew() {
			newCardsCount++
		}

		if dueToReview(&flashcard) {
			dueToReviewCount++
		}

		flashcardDTOs[i] = FlashcardDTO{Front: flashcard.Front, Back: flashcard.Back}
	}

	c.view.RenderHome(flashcardDTOs, newCardsCount, dueToReviewCount)
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

func (c *FlashcardsController) AddCard(front string, back string) {
	card := model.Flashcard{Front: front, Back: back, Metadata: &model.Metadata{
		CreationDate:       time.Now(),
		LastRepetitionDate: time.Now(),
		Memorizable:        supermemo.Create(),
	}}

	entries := card.Serialize()

	c.store.Add(parser.MakeLine(*entries))

	c.view.GoToHome()
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	entriesChan := parser.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		card := model.Flashcard{Front: entry[0], Back: entry[1], Metadata: &model.Metadata{
			CreationDate:       time.Now(),
			LastRepetitionDate: time.Now(),
			Memorizable:        supermemo.Create(),
		}}

		entries := card.Serialize()

		c.store.Add(parser.MakeLine(*entries))
	}

	c.view.GoToHome()
}

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize *[]FlashcardRecord
}

type FlashcardRecord struct {
	flashcard *model.Flashcard
	record    *persistance.Record
}

func (c *FlashcardsController) CreateMemorizingSession(count uint) {
	records := c.store.Read()

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

	c.view.GoToQuest()

	c.session = &MemorizingSession{memorizedCount: 0, cardsToMemorize: &flashcards}
}

func (c *FlashcardsController) CreateReviewSession(count uint) {
	records := c.store.Read()

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

	c.view.GoToQuest()

	c.session = &MemorizingSession{memorizedCount: 0, cardsToMemorize: &flashcards}
}

func (c *FlashcardsController) ShowQuest() {
	m := c.session

	if m.memorizedCount == len(*m.cardsToMemorize) {
		c.view.GoToHome()
		return
	}

	card := (*m.cardsToMemorize)[m.memorizedCount].flashcard

	cardDTO := FlashcardDTO{Front: card.Front, Back: card.Back}

	c.view.RenderCardQuestion(&cardDTO, m.memorizedCount+1, len(*m.cardsToMemorize))
}

func (c *FlashcardsController) ShowAnswer() {
	m := c.session

	if m.memorizedCount == len(*m.cardsToMemorize) {
		c.view.GoToHome()
	}

	card := (*m.cardsToMemorize)[m.memorizedCount].flashcard

	cardDTO := FlashcardDTO{Front: card.Front, Back: card.Back}

	c.view.RenderCardAnswer(
		&cardDTO,
		m.memorizedCount+1,
		len(*m.cardsToMemorize),
		m.getAnswerFeedbackOptions(),
	)
}

var answerFeedback = map[string]int{
	"Complete Blackout":        0,
	"Slipped my mind":          1,
	"Ah shit, I knew it!":      2,
	"Barely correct bro":       3,
	"I remembered correctly:)": 4,
	"Too easy!":                5,
}

func (m *MemorizingSession) getAnswerFeedbackOptions() []string {
	keys := make([]string, 0, len(answerFeedback))
	for k := range answerFeedback {
		keys = append(keys, k)
	}

	return keys
}

func (c *FlashcardsController) SubmitAnswer(answer string) {
	m := c.session
	card := (*m.cardsToMemorize)[m.memorizedCount]

	qualityOfResponse := supermemo.QualityOfResponse(answerFeedback[answer])

	card.flashcard.Metadata.Memorizable.SubmitRepetition(qualityOfResponse)
	card.flashcard.Metadata.LastRepetitionDate = time.Now()

	card.record.Data = parser.MakeLine(*card.flashcard.Serialize())
	card.record.Save()

	m.memorizedCount++

	c.view.GoToQuest()
}
