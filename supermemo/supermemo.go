package supermemo

import (
	"fmt"
	"math"
)

type Memorizable struct {
	RepetitionCount  int64
	NextReviewOffset Days
	EF               float64
}

// integer between 0-5
type QualityOfResponse byte
type Days int64

func Create() *Memorizable {
	return &Memorizable{EF: 2.5, NextReviewOffset: 0, RepetitionCount: 0}
}

func (m *Memorizable) IsNew() bool {
	return m.RepetitionCount == 0
}

func (m *Memorizable) GetNextRepetitionDaysOffset() Days {
	return m.NextReviewOffset
}

func (m *Memorizable) SubmitRepetition(qualityOfResponse QualityOfResponse) {
	m.RepetitionCount += 1

	nextOffset := calculateNextReviewOffset(m.RepetitionCount, m.EF)
	// make sure the days offset is in days
	m.NextReviewOffset = Days(math.Round(nextOffset))

	m.EF = calculateNextEF(m.EF, qualityOfResponse)
}

func calculateNextReviewOffset(repetitionCount int64, EF float64) float64 {
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

func calculateNextEF(oldEF float64, qualityOfResponse QualityOfResponse) float64 {
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
