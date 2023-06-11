package flashcard

import (
	"log"
	"math"
	"time"
)

type supermemo struct {
	RepetitionCount  int
	LastReviewDate   time.Time
	NextReviewOffset int
	EF               float64
}

func InitSupermemo() *supermemo {
	return &supermemo{EF: 2.5, NextReviewOffset: 0, RepetitionCount: 0, LastReviewDate: time.Now()}
}

func (m *supermemo) IsNew() bool {
	return m.RepetitionCount == 0
}

func (m *supermemo) SubmitRepetition(qualityOfResponse int) {
	if qualityOfResponse < 3 {
		m.RepetitionCount = 1
	} else {
		m.RepetitionCount += 1
	}

	m.LastReviewDate = time.Now()

	nextOffset := calculateNextReviewOffset(m.RepetitionCount, m.EF)
	m.NextReviewOffset = int(math.Round(nextOffset))

	m.EF = calculateNextEF(m.EF, qualityOfResponse)
}

func (m *supermemo) IsDueToReview() bool {
	if m.IsNew() {
		return false
	}

	nowInSeconds := time.Now().Unix()

	lastRepInSeconds := startOfDay(m.LastReviewDate).Unix()

	// convert to seconds
	offsetInSeconds := int64(m.NextReviewOffset * 24 * 60 * 60)

	return (lastRepInSeconds + offsetInSeconds) <= nowInSeconds
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func calculateNextReviewOffset(repetitionCount int, EF float64) float64 {
	if repetitionCount < 1 {
		log.Println("Error: repetitionCount cannot be less than 1")
		return 0
	}

	if repetitionCount == 1 {
		return 1
	}

	if repetitionCount == 2 {
		return 6
	}

	return calculateNextReviewOffset(repetitionCount-1, EF) * EF
}

func calculateNextEF(oldEF float64, qualityOfResponse int) float64 {
	if qualityOfResponse < 0 || qualityOfResponse > 5 {
		log.Println("Error: qualityOfResponse cannot must be within 0 and 5")
		return 0
	}

	newEF := oldEF + (float64(0.1) -
		(float64(5)-float64(qualityOfResponse))*
			(float64(0.08)+(float64(5)-float64(qualityOfResponse))*float64(0.02)))

	if newEF < 1.3 {
		return 1.3
	}

	return newEF
}
