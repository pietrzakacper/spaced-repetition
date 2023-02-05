package controller

import (
	"csv"
	"flashcard"
	"fmt"
	"io"
	"strings"
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

func (c *FlashcardsController) getAllCards() []flashcard.Record {
	notDeletedCards := make([]flashcard.Record, 0)

	for _, card := range c.store.ReadAll() {
		if card.Deleted == false {
			notDeletedCards = append(notDeletedCards, card)
		}
	}

	return notDeletedCards
}

func (c *FlashcardsController) ShowHome() {
	records := c.getAllCards()

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
	defer c.view.GoToHome()
	entriesChan, err := csv.ParseCSVStream(csvStream, []string{"front", "back"})

	if err != nil {
		fmt.Println("Error importing cards:", err)
		return
	}

	for entry := range entriesChan {
		fmt.Println(entry)
		card := (&flashcard.DTO{Front: strings.Trim(entry[0], " "), Back: strings.Trim(entry[1], " ")}).ToCard()

		record := card.ToRecord()

		c.store.Add(record)
	}
}

func (c *FlashcardsController) CreateMemorizingSession() {
	records := c.getAllCards()

	c.session = flashcard.CreateSession(records, flashcard.LearnNew)

	c.view.GoToQuest()
}

func (c *FlashcardsController) CreateReviewSession() {
	records := c.getAllCards()

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
		[]int{0, 2, 3, 5},
	)
}

func (c *FlashcardsController) SubmitAnswer(answer int) {
	card := c.session.SubmitAnswer(answer)

	c.store.Update(card.ToRecord())

	c.view.GoToQuest()
}

func (c *FlashcardsController) ShowCards() {
	records := c.getAllCards()

	flashcardDTOs := make([]flashcard.DTO, len(records))

	for i, r := range records {
		card := r.ToCard()

		flashcardDTOs[i] = *card.ToDTO()
	}

	c.view.RenderCards(flashcardDTOs)
}

func (c *FlashcardsController) DeleteCard(cardId string) {
	card, err := c.store.Find(cardId)

	if err != nil {
		fmt.Println(err)
		c.view.GoToCards()
		return
	}

	card.Deleted = true

	c.store.Update(&card)

	c.view.GoToCards()
}

func (c *FlashcardsController) EditCard(cardId, front, back string) {
	card, err := c.store.Find(cardId)

	if err != nil {
		fmt.Println(err)
		c.view.GoToCards()
		return
	}

	if len(front) > 0 {
		card.Front = front
	}

	if len(back) > 0 {
		card.Back = back
	}

	c.store.Update(&card)

	c.view.GoToCards()
}
