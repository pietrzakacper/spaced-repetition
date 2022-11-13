package controller

import (
	"io"
	"model"
	"parser"
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
		flashcard := &model.Flashcard{
			Id:    r.Id,
			Front: r.Front,
			Back:  r.Back,
			Metadata: &model.Metadata{
				CreationDate:   r.CreationDate,
				LastReviewDate: r.LastReviewDate,
				Memorizable: &supermemo.Memorizable{
					RepetitionCount:  r.RepetitionCount,
					NextReviewOffset: r.NextReviewOffset,
					EF:               r.EF,
				},
			},
		}

		if flashcard.Metadata.Memorizable.IsNew() {
			newCardsCount++
		}

		if dueToReview(flashcard) {
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

	lastRepInSeconds := card.Metadata.LastReviewDate.Unix()

	// convert to seconds
	offsetInSeconds := int64(card.Metadata.Memorizable.NextReviewOffset * 24 * 60 * 60)

	return (lastRepInSeconds + offsetInSeconds) <= nowInSeconds
}

func (c *FlashcardsController) AddCard(front string, back string) {
	card := model.Flashcard{Id: "", Front: front, Back: back, Metadata: &model.Metadata{
		CreationDate:   time.Now(),
		LastReviewDate: time.Now(),
		Memorizable:    supermemo.Create(),
	}}

	record := &FlashcardRecord{
		Id:               card.Id,
		Front:            card.Front,
		Back:             card.Back,
		CreationDate:     card.Metadata.CreationDate,
		LastReviewDate:   card.Metadata.LastReviewDate,
		NextReviewOffset: card.Metadata.Memorizable.NextReviewOffset,
		RepetitionCount:  card.Metadata.Memorizable.RepetitionCount,
		EF:               card.Metadata.Memorizable.EF,
	}

	c.store.Add(record)

	c.view.GoToHome()
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	entriesChan := parser.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		card := model.Flashcard{Id: "", Front: entry[0], Back: entry[1], Metadata: &model.Metadata{
			CreationDate:   time.Now(),
			LastReviewDate: time.Now(),
			Memorizable:    supermemo.Create(),
		}}

		record := &FlashcardRecord{
			Id:               card.Id,
			Front:            card.Front,
			Back:             card.Back,
			CreationDate:     card.Metadata.CreationDate,
			LastReviewDate:   card.Metadata.LastReviewDate,
			NextReviewOffset: card.Metadata.Memorizable.NextReviewOffset,
			RepetitionCount:  card.Metadata.Memorizable.RepetitionCount,
			EF:               card.Metadata.Memorizable.EF,
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

		card := &model.Flashcard{
			Id:    r.Id,
			Front: r.Front,
			Back:  r.Back,
			Metadata: &model.Metadata{
				CreationDate:   r.CreationDate,
				LastReviewDate: r.LastReviewDate,
				Memorizable: &supermemo.Memorizable{
					RepetitionCount:  r.RepetitionCount,
					NextReviewOffset: r.NextReviewOffset,
					EF:               r.EF,
				},
			},
		}

		if card.Metadata.Memorizable.IsNew() {
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

		// @TODO create only memorizable here
		card := &model.Flashcard{
			Id:    r.Id,
			Front: r.Front,
			Back:  r.Back,
			Metadata: &model.Metadata{
				CreationDate:   r.CreationDate,
				LastReviewDate: r.LastReviewDate,
				Memorizable: &supermemo.Memorizable{
					RepetitionCount:  r.RepetitionCount,
					NextReviewOffset: r.NextReviewOffset,
					EF:               r.EF,
				},
			},
		}

		if dueToReview(card) {
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

	card := &model.Flashcard{
		Id:    record.Id,
		Front: record.Front,
		Back:  record.Back,
		Metadata: &model.Metadata{
			CreationDate:   record.CreationDate,
			LastReviewDate: record.LastReviewDate,
			Memorizable: &supermemo.Memorizable{
				RepetitionCount:  record.RepetitionCount,
				NextReviewOffset: record.NextReviewOffset,
				EF:               record.EF,
			},
		},
	}

	card.Metadata.Memorizable.SubmitRepetition(qualityOfResponse)
	card.Metadata.LastReviewDate = time.Now()

	newRecord := &FlashcardRecord{
		Id:               record.Id,
		Front:            card.Front,
		Back:             card.Back,
		CreationDate:     card.Metadata.CreationDate,
		LastReviewDate:   card.Metadata.LastReviewDate,
		NextReviewOffset: card.Metadata.Memorizable.NextReviewOffset,
		RepetitionCount:  card.Metadata.Memorizable.RepetitionCount,
		EF:               card.Metadata.Memorizable.EF,
	}

	c.store.Update(newRecord)

	m.memorizedCount++

	c.view.GoToQuest()
}
