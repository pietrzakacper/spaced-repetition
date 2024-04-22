package flashcard

import (
	"strings"
	"time"
)

type flashcard struct {
	id           string
	front        string
	back         string
	flagged      bool
	creationDate time.Time

	supermemo *supermemo
}

type DTO struct {
	Id      string
	Front   string
	Back    string
	Flagged bool
}

type Record struct {
	Id               string
	Front            string
	Back             string
	Flagged          bool
	CreationDate     time.Time
	LastReviewDate   time.Time
	NextReviewOffset int
	RepetitionCount  int
	EF               float64
	Deleted          bool
}

const card_front_back_size_limit = 300

func withSizeLimit(s string) string {
	if len(s) > card_front_back_size_limit {
		return s[:card_front_back_size_limit]
	}

	return s
}

func (dto *DTO) ToCard() *flashcard {
	return &flashcard{
		front:        withSizeLimit(strings.Trim(dto.Front, " ")),
		back:         withSizeLimit(strings.Trim(dto.Back, " ")),
		flagged:      dto.Flagged,
		creationDate: time.Now(),
		supermemo:    InitSupermemo(),
	}
}

func (card *flashcard) ToDTO() *DTO {
	return &DTO{Id: card.id, Front: card.front, Back: card.back, Flagged: card.flagged}
}

func (record *Record) ToCard() *flashcard {
	return &flashcard{
		id:           record.Id,
		front:        record.Front,
		back:         record.Back,
		flagged:      record.Flagged,
		creationDate: record.CreationDate,
		supermemo: &supermemo{
			RepetitionCount:  record.RepetitionCount,
			LastReviewDate:   record.LastReviewDate,
			NextReviewOffset: record.NextReviewOffset,
			EF:               record.EF,
		},
	}
}

func (card *flashcard) ToRecord() *Record {
	return &Record{
		Id:               card.id,
		Front:            card.front,
		Back:             card.back,
		Flagged:          card.flagged,
		CreationDate:     card.creationDate,
		LastReviewDate:   card.supermemo.LastReviewDate,
		NextReviewOffset: card.supermemo.NextReviewOffset,
		RepetitionCount:  card.supermemo.RepetitionCount,
		EF:               card.supermemo.EF,
	}
}

func (card *flashcard) IsNew() bool {
	return card.supermemo.IsNew()
}

func (card *flashcard) IsDueToReview() bool {
	return card.supermemo.IsDueToReview()
}

func (card *flashcard) PostponeReview(byDays int) {
	card.supermemo.LastReviewDate = time.Now()
	card.supermemo.NextReviewOffset = byDays
}

func (card *flashcard) Flag(flagged bool) {
	card.flagged = flagged
}
