package flashcard

import (
	"time"
)

type flashcard struct {
	id           string
	front        string
	back         string
	creationDate time.Time

	answerSubmitted bool
	savedAnswer     int

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

func (dto *DTO) ToCard() *flashcard {
	return &flashcard{
		front:        dto.Front,
		back:         dto.Back,
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

type FlashcardSerialized struct {
	Id               string
	Front            string
	Back             string
	CreationDate     int64
	LastReviewDate   int64
	RepetitionCount  int
	NextReviewOffset int
	EF               float64
	Deleted          bool
}

func (f *FlashcardSerialized) ToRecord() Record {
	creationDate := time.UnixMicro(f.CreationDate)
	lastReviewDate := time.UnixMicro(f.LastReviewDate)

	return Record{
		Id:               f.Id,
		Front:            f.Front,
		Back:             f.Back,
		CreationDate:     creationDate,
		LastReviewDate:   lastReviewDate,
		RepetitionCount:  f.RepetitionCount,
		NextReviewOffset: f.NextReviewOffset,
		EF:               f.EF,
		Deleted:          f.Deleted,
	}
}
