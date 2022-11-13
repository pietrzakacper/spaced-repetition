package model

import (
	"supermemo"
	"time"
)

type Metadata struct {
	CreationDate       time.Time
	LastRepetitionDate time.Time
	Memorizable        *supermemo.Memorizable
}

type Flashcard struct {
	Front    string
	Back     string
	Metadata *Metadata
}

var timeLayout = "2006-02-01"

func (f *Flashcard) Serialize() *[]string {
	return &[]string{
		f.Front,
		f.Back,
		f.Metadata.CreationDate.Format(timeLayout),
		f.Metadata.LastRepetitionDate.Format(timeLayout),
		f.Metadata.Memorizable.Serialize(),
	}
}

func Deserialize(serialized []string) *Flashcard {
	creationDate, _ := time.Parse(timeLayout, serialized[2])
	lastRepetitionDate, _ := time.Parse(timeLayout, serialized[3])

	return &Flashcard{
		Front: serialized[0],
		Back:  serialized[1],
		Metadata: &Metadata{
			CreationDate:       creationDate,
			LastRepetitionDate: lastRepetitionDate,
			Memorizable:        supermemo.Deserialize(serialized[4]),
		},
	}
}
