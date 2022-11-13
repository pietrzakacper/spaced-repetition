package flashcard

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize []*flashcard
	failedCards     []*flashcard
}

type SessionType int

const (
	Review   SessionType = 0
	LearnNew SessionType = 1
)

var cardsInSession = 10

func CreateSession(records []Record, sessionType SessionType) *MemorizingSession {
	cardsToMemorize := make([]*flashcard, 0)

	for _, r := range records {
		if len(cardsToMemorize) == cardsInSession {
			break
		}

		card := r.ToCard()

		card.answerSubmitted = false

		if sessionType == Review && card.supermemo.IsDueToReview() {
			cardsToMemorize = append(cardsToMemorize, card)
		}

		if sessionType == LearnNew && card.supermemo.IsNew() {
			cardsToMemorize = append(cardsToMemorize, card)
		}
	}

	return &MemorizingSession{
		memorizedCount:  0,
		cardsToMemorize: cardsToMemorize,
		failedCards:     make([]*flashcard, 0),
	}
}

func (m *MemorizingSession) SubmitAnswer(answer int) *flashcard {
	card := m.CurrentCard()

	if answer < 4 {
		m.failedCards = append(m.failedCards, card)
	}

	if card.answerSubmitted == false {
		m.CurrentCard().supermemo.SubmitRepetition(answer)
		card.answerSubmitted = true
	}

	m.GoToNext()

	return card
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
