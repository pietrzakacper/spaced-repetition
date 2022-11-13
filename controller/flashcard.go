package controller

import (
	"supermemo"
	"time"
)

type flashcard struct {
	id           string
	front        string
	back         string
	creationDate time.Time
	memorizable  *supermemo.Memorizable
}

func (card *flashcard) toDTO() *FlashcardDTO {
	return &FlashcardDTO{Front: card.front, Back: card.back}
}

func (*flashcard) fromRecord(record *FlashcardRecord) *flashcard {
	return &flashcard{
		id:           record.Id,
		front:        record.Front,
		back:         record.Back,
		creationDate: record.CreationDate,
		memorizable: &supermemo.Memorizable{
			RepetitionCount:  record.RepetitionCount,
			LastReviewDate:   record.LastReviewDate,
			NextReviewOffset: record.NextReviewOffset,
			EF:               record.EF,
		},
	}
}

func (card *flashcard) toRecord() *FlashcardRecord {
	return &FlashcardRecord{
		Id:               card.id,
		Front:            card.front,
		Back:             card.back,
		CreationDate:     card.creationDate,
		LastReviewDate:   card.memorizable.LastReviewDate,
		NextReviewOffset: card.memorizable.NextReviewOffset,
		RepetitionCount:  card.memorizable.RepetitionCount,
		EF:               card.memorizable.EF,
	}
}
