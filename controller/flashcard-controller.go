package controller

import (
	"csv"
	"flashcard"
	"io"
)

type FlashcardsController struct {
	view    View
	store   Store
	session *flashcard.MemorizingSession
}

func CreateFlashcardsController(view View, persistance Persistance) *FlashcardsController {
	return &FlashcardsController{
		view:  view,
		store: persistance.Create("db"),
	}
}

func (c *FlashcardsController) ShowHome() {
	records := c.store.ReadAll()

	flashcardDTOs := make([]flashcard.DTO, len(records))

	newCardsCount := 0
	dueToReviewCount := 0

	for i, r := range records {
		card := r.ToCard()

		if card.IsNew() {
			newCardsCount++
		}

		if card.IsDueToReview() {
			dueToReviewCount++
		}

		flashcardDTOs[i] = *card.ToDTO()
	}

	c.view.RenderHome(flashcardDTOs, newCardsCount, dueToReviewCount)
}

func (c *FlashcardsController) AddCard(front string, back string) {
	card := (&flashcard.DTO{Front: front, Back: back}).ToCard()

	record := card.ToRecord()

	c.store.Add(record)

	c.view.GoToHome()
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	entriesChan := csv.ParseCSVStream(csvStream)

	for entry := range entriesChan {
		card := (&flashcard.DTO{Front: entry[0], Back: entry[1]}).ToCard()

		record := card.ToRecord()

		c.store.Add(record)
	}

	c.view.GoToHome()
}

func (c *FlashcardsController) CreateMemorizingSession() {
	records := c.store.ReadAll()

	c.session = flashcard.CreateSession(records, flashcard.LearnNew)

	c.view.GoToQuest()

}

func (c *FlashcardsController) CreateReviewSession() {
	records := c.store.ReadAll()

	c.session = flashcard.CreateSession(records, flashcard.Review)

	c.view.GoToQuest()
}

func (c *FlashcardsController) ShowQuest() {
	if c.session.HasEnded() {
		if c.session.HasAnyFailedCards() {
			c.session.ReviewFailedCardsAgain()
		} else {
			c.view.GoToHome()
			return
		}
	}

	card := c.session.CurrentCard()

	c.view.RenderCardQuestion(
		card.ToDTO(),
		c.session.CurrentCardNumber(),
		c.session.TotalCardsNumber(),
	)
}

func (c *FlashcardsController) ShowAnswer() {
	if c.session.HasEnded() {
		c.view.GoToQuest()
	}

	card := c.session.CurrentCard()

	c.view.RenderCardAnswer(
		card.ToDTO(),
		c.session.CurrentCardNumber(),
		c.session.TotalCardsNumber(),
		[]int{0, 1, 2, 3, 4, 5},
	)
}

func (c *FlashcardsController) SubmitAnswer(answer int) {
	card := c.session.SubmitAnswer(answer)

	c.store.Update(card.ToRecord())

	c.view.GoToQuest()
}
