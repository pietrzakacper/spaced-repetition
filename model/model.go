package model

import (
	"supermemo"
	"time"
)

type Metadata struct {
	CreationDate   time.Time
	LastReviewDate time.Time
	Memorizable    *supermemo.Memorizable
}

type Flashcard struct {
	Id       string
	Front    string
	Back     string
	Metadata *Metadata
}
