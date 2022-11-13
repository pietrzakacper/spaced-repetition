package supermemo

import (
	"fmt"
	"math"
	"time"
)

type Memorizable struct {
	RepetitionCount  int
	LastReviewDate   time.Time
	NextReviewOffset int
	EF               float64
}

func Create() *Memorizable {
	return &Memorizable{EF: 2.5, NextReviewOffset: 0, RepetitionCount: 0, LastReviewDate: time.Now()}
}

func (m *Memorizable) IsNew() bool {
	return m.RepetitionCount == 0
}

func (m *Memorizable) SubmitRepetition(qualityOfResponse int) {
	m.RepetitionCount += 1
	m.LastReviewDate = time.Now()

	nextOffset := calculateNextReviewOffset(m.RepetitionCount, m.EF)
	m.NextReviewOffset = int(math.Round(nextOffset))

	m.EF = calculateNextEF(m.EF, qualityOfResponse)
}

func (m *Memorizable) IsDueToReview() bool {
	if m.IsNew() {
		return false
	}

	nowInSeconds := time.Now().Unix()

	lastRepInSeconds := m.LastReviewDate.Unix()

	// convert to seconds
	offsetInSeconds := int64(m.NextReviewOffset * 24 * 60 * 60)

	return (lastRepInSeconds + offsetInSeconds) <= nowInSeconds
}

func calculateNextReviewOffset(repetitionCount int, EF float64) float64 {
	if repetitionCount < 1 {
		fmt.Println("Error: repetitionCount cannot be less than 1")
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
	if qualityOfResponse > 5 {
		fmt.Println("Error: qualityOfResponse cannot be more than 5")
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
