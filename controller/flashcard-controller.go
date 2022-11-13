package controller

import (
	"csv"
	"io"
	"supermemo"
	"time"
)

type FlashcardsController struct {
	view    View
	store   Store
	session *MemorizingSession
}

func CreateFlashcardsController(view View, persistance Persistance) *FlashcardsController {
	return &FlashcardsController{
		view:  view,
		store: persistance.Create("db"),
	}
}

func (c *FlashcardsController) ShowHome() {
	records := c.store.ReadAll()

	flashcardDTOs := make([]FlashcardDTO, len(records))

	newCardsCount := 0
	dueToReviewCount := 0

	for i, r := range records {
		flashcard := &flashcard{
			id:           r.Id,
			front:        r.Front,
			back:         r.Back,
			creationDate: r.CreationDate,
			memorizable: &supermemo.Memorizable{
				RepetitionCount:  r.RepetitionCount,
				LastReviewDate:   r.LastReviewDate,
				NextReviewOffset: r.NextReviewOffset,
				EF:               r.EF,
			},
		}

		if flashcard.memorizable.IsNew() {
			newCardsCount++
		}

		if flashcard.memorizable.IsDueToReview() {
			dueToReviewCount++
		}

		flashcardDTOs[i] = FlashcardDTO{Front: flashcard.front, Back: flashcard.back}
	}

	c.view.RenderHome(flashcardDTOs, newCardsCount, dueToReviewCount)
}

func (c *FlashcardsController) AddCard(front string, back string) {
	card := flashcard{
		id:           "",
		front:        front,
		back:         back,
		creationDate: time.Now(),
		memorizable:  supermemo.Create(),
	}

	record := &FlashcardRecord{
		Id:               card.id,
		Front:            card.front,
		Back:             card.back,
		CreationDate:     card.creationDate,
		LastReviewDate:   card.memorizable.LastReviewDate,
		NextReviewOffset: card.memorizable.NextReviewOffset,
		RepetitionCount:  card.memorizable.RepetitionCount,
		EF:               card.memorizable.EF,
	}

	c.store.Add(record)

	c.view.GoToHome()
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	entriesChan := csv.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		card := flashcard{
			id:           "",
			front:        entry[0],
			back:         entry[1],
			creationDate: time.Now(),
			memorizable:  supermemo.Create(),
		}

		record := &FlashcardRecord{
			Id:               card.id,
			Front:            card.front,
			Back:             card.back,
			CreationDate:     card.creationDate,
			LastReviewDate:   card.memorizable.LastReviewDate,
			NextReviewOffset: card.memorizable.NextReviewOffset,
			RepetitionCount:  card.memorizable.RepetitionCount,
			EF:               card.memorizable.EF,
		}

		c.store.Add(record)
	}

	c.view.GoToHome()
}

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize []FlashcardRecord
}

func (c *FlashcardsController) CreateMemorizingSession(count int) {
	records := c.store.ReadAll()

	flashcardsInSession := make([]FlashcardRecord, 0)

	for _, r := range records {
		if len(flashcardsInSession) == count {
			break
		}

		memorizable := &supermemo.Memorizable{
			RepetitionCount:  r.RepetitionCount,
			LastReviewDate:   r.LastReviewDate,
			NextReviewOffset: r.NextReviewOffset,
			EF:               r.EF,
		}

		if memorizable.IsNew() {
			flashcardsInSession = append(flashcardsInSession, r)
		}
	}

	c.view.GoToQuest()

	c.session = &MemorizingSession{memorizedCount: 0, cardsToMemorize: flashcardsInSession}
}

func (c *FlashcardsController) CreateReviewSession(count int) {
	records := c.store.ReadAll()

	flashcardsInSession := make([]FlashcardRecord, 0)

	for _, r := range records {
		if len(flashcardsInSession) == count {
			break
		}

		memorizable := &supermemo.Memorizable{
			RepetitionCount:  r.RepetitionCount,
			LastReviewDate:   r.LastReviewDate,
			NextReviewOffset: r.NextReviewOffset,
			EF:               r.EF,
		}

		if memorizable.IsDueToReview() {
			flashcardsInSession = append(flashcardsInSession, r)
		}
	}

	c.view.GoToQuest()

	c.session = &MemorizingSession{memorizedCount: 0, cardsToMemorize: flashcardsInSession}
}

func (c *FlashcardsController) ShowQuest() {
	m := c.session

	if m.memorizedCount == len(m.cardsToMemorize) {
		c.view.GoToHome()
		return
	}

	card := m.cardsToMemorize[m.memorizedCount]

	cardDTO := FlashcardDTO{Front: card.Front, Back: card.Back}

	c.view.RenderCardQuestion(&cardDTO, m.memorizedCount+1, len(m.cardsToMemorize))
}

func (c *FlashcardsController) ShowAnswer() {
	m := c.session

	if m.memorizedCount == len(m.cardsToMemorize) {
		c.view.GoToHome()
	}

	card := m.cardsToMemorize[m.memorizedCount]

	cardDTO := FlashcardDTO{Front: card.Front, Back: card.Back}

	c.view.RenderCardAnswer(
		&cardDTO,
		m.memorizedCount+1,
		len(m.cardsToMemorize),
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
	record := m.cardsToMemorize[m.memorizedCount]

	qualityOfResponse := supermemo.QualityOfResponse(answerFeedback[answer])

	card := &flashcard{
		id:           record.Id,
		front:        record.Front,
		back:         record.Back,
		creationDate: record.CreationDate,
		memorizable: &supermemo.Memorizable{
			RepetitionCount:  record.RepetitionCount,
			LastReviewDate:   record.LastReviewDate,
			NextReviewOffset: record.NextReviewOffset,
			EF:               record.EF,
		},
	}

	card.memorizable.SubmitRepetition(qualityOfResponse)

	newRecord := &FlashcardRecord{
		Id:               record.Id,
		Front:            card.front,
		Back:             card.back,
		CreationDate:     card.creationDate,
		LastReviewDate:   card.memorizable.LastReviewDate,
		NextReviewOffset: card.memorizable.NextReviewOffset,
		RepetitionCount:  card.memorizable.RepetitionCount,
		EF:               card.memorizable.EF,
	}

	c.store.Update(newRecord)

	m.memorizedCount++

	c.view.GoToQuest()
}
