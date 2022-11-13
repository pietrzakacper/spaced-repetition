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
		card := (&flashcard{}).fromRecord(&r)

		if card.memorizable.IsNew() {
			newCardsCount++
		}

		if card.memorizable.IsDueToReview() {
			dueToReviewCount++
		}

		flashcardDTOs[i] = *card.toDTO()
	}

	c.view.RenderHome(flashcardDTOs, newCardsCount, dueToReviewCount)
}

func (c *FlashcardsController) AddCard(front string, back string) {
	card := flashcard{
		front:        front,
		back:         back,
		creationDate: time.Now(),
		memorizable:  supermemo.Create(),
	}

	record := card.toRecord()

	c.store.Add(record)

	c.view.GoToHome()
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	entriesChan := csv.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		card := flashcard{
			front:        entry[0],
			back:         entry[1],
			creationDate: time.Now(),
			memorizable:  supermemo.Create(),
		}

		record := card.toRecord()

		c.store.Add(record)
	}

	c.view.GoToHome()
}

func (c *FlashcardsController) CreateMemorizingSession() {
	records := c.store.ReadAll()

	c.session = CreateSession(records, LearnNew)

	c.view.GoToQuest()

}

func (c *FlashcardsController) CreateReviewSession() {
	records := c.store.ReadAll()

	c.session = CreateSession(records, Review)

	c.view.GoToQuest()
}

func (c *FlashcardsController) ShowQuest() {
	if c.session.hasEnded() {
		c.view.GoToHome()
		return
	}

	card := c.session.currentCard()

	c.view.RenderCardQuestion(
		card.toDTO(),
		c.session.memorizedCount,
		len(c.session.cardsToMemorize),
	)
}

func (c *FlashcardsController) ShowAnswer() {
	if c.session.hasEnded() {
		c.view.GoToHome()
	}

	card := c.session.currentCard()

	c.view.RenderCardAnswer(
		card.toDTO(),
		c.session.memorizedCount+1,
		len(c.session.cardsToMemorize),
		[]int{0, 1, 2, 3, 4, 5},
	)
}

func (c *FlashcardsController) SubmitAnswer(answer int) {
	card := c.session.currentCard()
	card.memorizable.SubmitRepetition(answer)

	c.store.Update(card.toRecord())

	c.session.goToNext()

	c.view.GoToQuest()
}
