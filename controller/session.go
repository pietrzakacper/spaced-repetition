package controller

type MemorizingSession struct {
	memorizedCount  int
	cardsToMemorize []flashcard
}

type SessionType int

const (
	Review   SessionType = 0
	LearnNew SessionType = 1
)

var cardsInSession = 10

func CreateSession(records []FlashcardRecord, sessionType SessionType) *MemorizingSession {
	cardsToMemorize := make([]flashcard, 0)

	for _, r := range records {
		if len(cardsToMemorize) == cardsInSession {
			break
		}

		card := (&flashcard{}).fromRecord(&r)

		if sessionType == Review && card.memorizable.IsDueToReview() {
			cardsToMemorize = append(cardsToMemorize, *card)
		}

		if sessionType == LearnNew && card.memorizable.IsNew() {
			cardsToMemorize = append(cardsToMemorize, *card)
		}
	}

	return &MemorizingSession{memorizedCount: 0, cardsToMemorize: cardsToMemorize}
}

func (m *MemorizingSession) currentCard() *flashcard {
	return &m.cardsToMemorize[m.memorizedCount]
}

func (m *MemorizingSession) goToNext() {
	m.memorizedCount++
}

func (m *MemorizingSession) hasEnded() bool {
	return m.memorizedCount >= len(m.cardsToMemorize)
}
