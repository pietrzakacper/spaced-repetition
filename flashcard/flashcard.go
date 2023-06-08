package flashcard

import (
	"strings"
	"time"
)

type flashcard struct {
	id           string
	front        string
	back         string
	creationDate time.Time

	supermemo *supermemo
}

type DTO struct {
	Id    string
	Front string
	Back  string
}

type Record struct {
	Id               string
	Front            string
	Back             string
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
		creationDate: time.Now(),
		supermemo:    InitSupermemo(),
	}
}

func (card *flashcard) ToDTO() *DTO {
	return &DTO{Id: card.id, Front: card.front, Back: card.back}
}

func (record *Record) ToCard() *flashcard {
	return &flashcard{
		id:           record.Id,
		front:        record.Front,
		back:         record.Back,
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
