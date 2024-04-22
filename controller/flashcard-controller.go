package controller

import (
	"encoding/csv"
	"flashcard"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/pietrzakacper/tracethat.dev/reporters/golang/tt"
)

type FlashcardsController struct {
	view  View
	store Store
}

func CreateFlashcardsController(view View, store Store) *FlashcardsController {
	return &FlashcardsController{view, store}
}

func (c *FlashcardsController) getAllCards() []flashcard.Record {
	defer tt.LogWithTime("get all cards")()
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
	defer tt.LogWithTime("show-home", flashcardDTOs)()

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

func normalizeColumn(col string) string {
	return strings.ToLower(strings.Trim(col, " "))
}

func (c *FlashcardsController) ImportCards(csvStream io.Reader) {
	defer c.view.GoToHome()
	r := csv.NewReader(csvStream)
	columns, err := r.Read()

	if err != nil {
		log.Println("Error importing cards:", err)
		return
	}

	frontIndex := -1
	backIndex := -1
	for colIndex := range columns {
		if normalizeColumn(columns[colIndex]) == "front" {
			frontIndex = colIndex
		} else if normalizeColumn(columns[colIndex]) == "back" {
			backIndex = colIndex
		}
	}

	if frontIndex < 0 || backIndex < 0 {
		log.Println("Incorrect columns in CSV:", columns)
		return
	}

	var wg sync.WaitGroup

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading CSV line", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			card := (&flashcard.DTO{Front: line[frontIndex], Back: line[backIndex]}).ToCard()

			record := card.ToRecord()

			c.store.Add(record)
		}()
	}

	wg.Wait()
}

func (c *FlashcardsController) CreateMemorizingSession() {
	records := c.getAllCards()

	session := flashcard.CreateSession(records, flashcard.LearnNew)

	c.view.UpdateClientSession(session.ToDTO())
	c.view.GoToQuest()
}

func (c *FlashcardsController) CreateReviewSession() {
	records := c.getAllCards()

	session := flashcard.CreateSession(records, flashcard.Review)

	c.view.UpdateClientSession(session.ToDTO())
	c.view.GoToQuest()
}

func (c *FlashcardsController) ShowQuest(sessionDTO *flashcard.MemorizingSessionDTO) {
	session := sessionDTO.ToMemorizingSession()

	if !session.IsValid() {
		c.view.GoToHome()
		return
	}

	if session.HasEnded() {
		if session.HasAnyFailedCards() {
			session.ReviewFailedCardsAgain()
			c.view.UpdateClientSession(session.ToDTO())
		} else {
			c.view.OnSessionFinished()
			c.view.GoToHome()
			return
		}
	}

	card := session.CurrentCard()

	c.view.RenderCardQuestion(
		card.ToDTO(),
		session.CurrentCardNumber(),
		session.TotalCardsNumber(),
		session.IsExtraRound(),
	)
}

func (c *FlashcardsController) ShowAnswer(sessionDTO *flashcard.MemorizingSessionDTO) {
	session := sessionDTO.ToMemorizingSession()

	if !session.IsValid() {
		c.view.GoToHome()
		return
	}

	if session.HasEnded() {
		c.view.GoToQuest()
		return
	}

	card := session.CurrentCard()

	c.view.RenderCardAnswer(
		card.ToDTO(),
		session.CurrentCardNumber(),
		session.TotalCardsNumber(),
		[]int{0, 2, 4, 5},
	)
}

func (c *FlashcardsController) SubmitAnswer(sessionDTO *flashcard.MemorizingSessionDTO, answer int) {
	session := sessionDTO.ToMemorizingSession()

	if !session.IsValid() {
		c.view.GoToHome()
		return
	}

	if session.HasEnded() {
		c.view.GoToQuest()
		return
	}

	card := session.SubmitAnswer(answer)

	if session.HasEnded() && !session.HasAnyFailedCards() {
		// if this is the last card in session, we want the DB update to happen before redirect
		// to show consistent view to the user
		c.store.Update(card.ToRecord())
		c.view.OnSessionFinished()
		c.view.GoToHome()
	} else {
		c.view.UpdateClientSession(session.ToDTO())
		c.view.GoToQuest()
		go c.store.Update(card.ToRecord())
	}

}

func (c *FlashcardsController) ShowCards() {
	records := c.getAllCards()

	flashcardDTOs := make([]flashcard.DTO, len(records))

	for i, r := range records {
		card := r.ToCard()

		// add cards in reversed order (from Newest to Oldest)
		flashcardDTOs[len(records)-1-i] = *card.ToDTO()
	}

	c.view.RenderCards(flashcardDTOs)
}

func (c *FlashcardsController) DeleteCard(cardId string) error {
	go func() {
		card, err := c.store.Find(cardId)

		if err != nil {
			log.Println(err)
		}

		card.Deleted = true

		c.store.Update(&card)
	}()
	return nil
}

func (c *FlashcardsController) EditCard(cardId, front, back string) error {
	go func() {
		card, err := c.store.Find(cardId)

		if err != nil {
			log.Println(err)
		}

		if len(front) > 0 {
			card.Front = front
		}

		if len(back) > 0 {
			card.Back = back
		}

		c.store.Update(&card)
	}()

	return nil
}

func (c *FlashcardsController) SpreadCardsInTime() error {
	records := c.getAllCards()

	counter := 0
	const cardsPerDay = 10

	var wg sync.WaitGroup
	for _, r := range records {
		card := r.ToCard()

		if card.IsDueToReview() {
			postponeByDays := counter / cardsPerDay
			card.PostponeReview(postponeByDays)
			counter++

			wg.Add(1)
			go func() {
				c.store.Update(card.ToRecord())
				wg.Done()
			}()
		}
	}

	wg.Wait()
	return nil
}

func (c *FlashcardsController) FlagCard(
	cardId string,
	flagged bool,
	memorizingSession *flashcard.MemorizingSessionDTO) error {
	// if is in the middle of memorizing session
	if len(memorizingSession.CardsToMemorize) > 0 {
		session := memorizingSession.ToMemorizingSession()
		session.CurrentCard().Flag(flagged)
		c.view.UpdateClientSession(session.ToDTO())
	}

	go func() {
		card, err := c.store.Find(cardId)

		if err != nil {
			log.Println(err)
		}

		card.Flagged = flagged

		c.store.Update(&card)
	}()

	return nil
}
