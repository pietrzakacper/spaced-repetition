package flashcard

import (
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize []*flashcard
	failedCards     []*flashcard
	answers         map[string]int // key=cardId value=savedAnswer
}

type FlashcardInSessionDTO struct {
	Id               string
	Front            string
	Back             string
	CreationDate     time.Time
	RepetitionCount  int
	LastReviewDate   time.Time
	NextReviewOffset int
	EF               float64
}

type MemorizingSessionDTO struct {
	MemorizedCount  int
	CardsToMemorize []*FlashcardInSessionDTO
	FailedCards     []*FlashcardInSessionDTO
	Answers         map[string]int
}

func (f *flashcard) ToMemorizingSessionDTO() *FlashcardInSessionDTO {
	return &FlashcardInSessionDTO{
		Id:               f.id,
		Front:            f.front,
		Back:             f.back,
		CreationDate:     f.creationDate,
		RepetitionCount:  f.supermemo.RepetitionCount,
		LastReviewDate:   f.supermemo.LastReviewDate,
		NextReviewOffset: f.supermemo.NextReviewOffset,
		EF:               f.supermemo.EF,
	}
}

func (f *FlashcardInSessionDTO) ToFlashcard() *flashcard {
	return &flashcard{
		id:           f.Id,
		front:        f.Front,
		back:         f.Back,
		creationDate: f.CreationDate,
		supermemo: &supermemo{
			RepetitionCount:  f.RepetitionCount,
			LastReviewDate:   f.LastReviewDate,
			NextReviewOffset: f.NextReviewOffset,
			EF:               f.EF,
		},
	}
}

type SessionType int

const (
	Review   SessionType = 0
	LearnNew SessionType = 1
)

var cardsInSession = 10

func CreateSession(records []Record, sessionType SessionType) *MemorizingSession {
	// cards that are valid (DueToReview or New depending on sessionType)
	allValidCardsForSession := make([]*flashcard, 0)

	for _, r := range records {
		card := r.ToCard()

		if sessionType == Review && card.supermemo.IsDueToReview() {
			allValidCardsForSession = append(allValidCardsForSession, card)
		}

		if sessionType == LearnNew && card.supermemo.IsNew() {
			allValidCardsForSession = append(allValidCardsForSession, card)
		}
	}

	// sort allValidCardsForSession randomly
	slices.SortFunc(allValidCardsForSession, func(a, b *flashcard) bool {
		return rand.Intn(10) > 5
	})

	cardsToMemorize := make([]*flashcard, 0)

	for _, c := range allValidCardsForSession {
		if len(cardsToMemorize) == cardsInSession {
			break
		}

		cardsToMemorize = append(cardsToMemorize, c)
	}

	return &MemorizingSession{
		memorizedCount:  0,
		cardsToMemorize: cardsToMemorize,
		failedCards:     make([]*flashcard, 0),
		answers:         make(map[string]int),
	}
}

func (m *MemorizingSession) ToDTO() *MemorizingSessionDTO {
	cardsToMemorizeDTOs := make([]*FlashcardInSessionDTO, len(m.cardsToMemorize))

	for i, f := range m.cardsToMemorize {
		cardsToMemorizeDTOs[i] = f.ToMemorizingSessionDTO()
	}

	failedCardsDTOs := make([]*FlashcardInSessionDTO, len(m.failedCards))

	for i, f := range m.failedCards {
		failedCardsDTOs[i] = f.ToMemorizingSessionDTO()
	}

	return &MemorizingSessionDTO{
		CardsToMemorize: cardsToMemorizeDTOs,
		FailedCards:     failedCardsDTOs,
		MemorizedCount:  m.memorizedCount,
		Answers:         m.answers,
	}
}

func (m *MemorizingSessionDTO) ToMemorizingSession() *MemorizingSession {
	cardsToMemorize := make([]*flashcard, len(m.CardsToMemorize))

	for i, f := range m.CardsToMemorize {
		cardsToMemorize[i] = f.ToFlashcard()
	}

	failedCards := make([]*flashcard, len(m.FailedCards))

	for i, f := range m.FailedCards {
		failedCards[i] = f.ToFlashcard()
	}

	return &MemorizingSession{
		cardsToMemorize: cardsToMemorize,
		failedCards:     failedCards,
		memorizedCount:  m.MemorizedCount,
		answers:         m.Answers,
	}
}

func (m *MemorizingSession) SubmitAnswer(answer int) *flashcard {
	card := m.CurrentCard()

	if _, hasAnswer := m.answers[card.id]; !hasAnswer {
		m.answers[card.id] = answer
	}

	if answer < 4 {
		m.failedCards = append(m.failedCards, card)
	} else {
		m.CurrentCard().supermemo.SubmitRepetition(m.answers[card.id])
	}

	m.GoToNext()

	return card
}

func (m *MemorizingSession) IsValid() bool {
	return m.cardsToMemorize != nil &&
		m.failedCards != nil &&
		m.answers != nil &&
		len(m.cardsToMemorize) > 0 &&
		m.memorizedCount >= 0 &&
		m.memorizedCount <= len(m.cardsToMemorize) &&
		len(m.failedCards) <= len(m.cardsToMemorize)
}

func (m *MemorizingSession) CurrentCard() *flashcard {
	return m.cardsToMemorize[m.memorizedCount]
}

func (m *MemorizingSession) GoToNext() {
	m.memorizedCount++
}

func (m *MemorizingSession) HasEnded() bool {
	return m.memorizedCount >= len(m.cardsToMemorize)
}

func (m *MemorizingSession) CurrentCardNumber() int {
	return m.memorizedCount + 1
}

func (m *MemorizingSession) TotalCardsNumber() int {
	return len(m.cardsToMemorize)
}

func (m *MemorizingSession) HasAnyFailedCards() bool {
	return len(m.failedCards) > 0
}

func (m *MemorizingSession) ReviewFailedCardsAgain() {
	m.memorizedCount = 0
	m.cardsToMemorize = m.failedCards
	m.failedCards = make([]*flashcard, 0)
}

/* are we showing an extra round to memorize the failed cards? */
func (m *MemorizingSession) IsExtraRound() bool {
	_, hasAnswer := m.answers[m.CurrentCard().id]
	return hasAnswer
}
