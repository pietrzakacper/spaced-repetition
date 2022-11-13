package supermemo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Memorizable struct {
	repetitionCount      uint
	nextRepetitionOffset Days
	ef                   float64
}

// integer between 0-5
type QualityOfResponse byte
type Days uint

func Create() *Memorizable {
	return &Memorizable{ef: 2.5, nextRepetitionOffset: 0, repetitionCount: 0}
}

func (m *Memorizable) IsNew() bool {
	return m.repetitionCount == 0
}

func (m *Memorizable) GetNextRepetitionDaysOffset() Days {
	return m.nextRepetitionOffset
}

func (m *Memorizable) SubmitRepetition(qualityOfResponse QualityOfResponse) {
	m.repetitionCount += 1

	nextOffset := calculateNextRepetitionOffset(m.repetitionCount, m.ef)
	// make sure the days offset is in days
	m.nextRepetitionOffset = Days(math.Round(nextOffset))

	m.ef = calculateNextEF(m.ef, qualityOfResponse)
}

func (m *Memorizable) Serialize() string {
	return fmt.Sprintf("%d", m.repetitionCount) +
		"|" + fmt.Sprintf("%d", m.nextRepetitionOffset) +
		"|" + fmt.Sprintf("%.2f", m.ef)
}

func Deserialize(serializedCard string) *Memorizable {
	components := strings.Split(serializedCard, "|")

	if len(components) != 3 {
		fmt.Println("Error: Incorrect number of components in a serializedCard", serializedCard)
		return nil
	}

	repetitionCountParsed, err := strconv.ParseInt(components[0], 10, 64)

	if err != nil {
		fmt.Println("Error: Failed to parse repetitionCount", components[0])
		return nil
	}

	nextRepetitionOffsetParsed, err := strconv.ParseInt(components[1], 10, 64)

	if err != nil {
		fmt.Println("Error: Failed to parse nextRepetitionOffset", components[1])
		return nil
	}

	efParsed, err := strconv.ParseFloat(components[2], 64)

	if err != nil {
		fmt.Println("Error: Failed to parse ef", components[2])
		return nil
	}

	return &Memorizable{
		repetitionCount:      uint(repetitionCountParsed),
		nextRepetitionOffset: Days(nextRepetitionOffsetParsed),
		ef:                   efParsed,
	}
}

func calculateNextRepetitionOffset(repetitionCount uint, EF float64) float64 {
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

	return calculateNextRepetitionOffset(repetitionCount-1, EF) * EF
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
		// @TODO verify if it works
		return 1.3
	}

	return newEF
}
