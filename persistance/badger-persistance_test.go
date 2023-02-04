package persistance

import (
	"flashcard"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestCreate(t *testing.T) {
	p := BadgerPersistance{inMemory: true}

	store := p.Create("db")

	store.Add(&flashcard.Record{
		Id:               "",
		Front:            "Hello",
		Back:             "Hola",
		CreationDate:     time.UnixMilli(int64(1675453002000)),
		LastReviewDate:   time.UnixMilli(int64(1675453002000)),
		NextReviewOffset: 1,
		RepetitionCount:  2,
		EF:               1.5,
	})

	store.Add(&flashcard.Record{
		Id:               "",
		Front:            "Bye",
		Back:             "Adios",
		CreationDate:     time.UnixMilli(int64(1675276600000)),
		LastReviewDate:   time.UnixMilli(int64(1675276600000)),
		NextReviewOffset: 123,
		RepetitionCount:  200,
		EF:               1.5353531,
	})

	cards := store.ReadAll()

	assert.Equal(t, len(cards), 2, "Cards length must be 2")

	helloCardIndex := slices.IndexFunc(cards, func(c flashcard.Record) bool { return c.Front == "Hello" })

	assert.Equal(t, "Hello", cards[helloCardIndex].Front)
	assert.Equal(t, "Hola", cards[helloCardIndex].Back)
	assert.Equal(t, "2023-03-02", cards[helloCardIndex].CreationDate.Format(daysPrecision))
	assert.Equal(t, "2023-03-02", cards[helloCardIndex].LastReviewDate.Format(daysPrecision))
	assert.Equal(t, 1, cards[helloCardIndex].NextReviewOffset)
	assert.Equal(t, 2, cards[helloCardIndex].RepetitionCount)
	assert.Equal(t, 1.5, cards[helloCardIndex].EF)

	byeCardIndex := slices.IndexFunc(cards, func(c flashcard.Record) bool { return c.Front == "Bye" })

	assert.Equal(t, "Bye", cards[byeCardIndex].Front)
	assert.Equal(t, "Adios", cards[byeCardIndex].Back)
	assert.Equal(t, "2023-01-02", cards[byeCardIndex].CreationDate.Format(daysPrecision))
	assert.Equal(t, "2023-01-02", cards[byeCardIndex].LastReviewDate.Format(daysPrecision))
	assert.Equal(t, 123, cards[byeCardIndex].NextReviewOffset)
	assert.Equal(t, 200, cards[byeCardIndex].RepetitionCount)
	assert.Equal(t, 1.5353531, cards[byeCardIndex].EF)

	cards[helloCardIndex].NextReviewOffset = 2
	cards[helloCardIndex].RepetitionCount = 3
	cards[helloCardIndex].EF = 2.113

	store.Update(&cards[helloCardIndex])

	cards = store.ReadAll()

	assert.Equal(t, len(cards), 2, "Cards length must be 2")

	helloCardIndex = slices.IndexFunc(cards, func(c flashcard.Record) bool { return c.Front == "Hello" })

	assert.Equal(t, "Hello", cards[helloCardIndex].Front)
	assert.Equal(t, "Hola", cards[helloCardIndex].Back)
	assert.Equal(t, "2023-03-02", cards[helloCardIndex].CreationDate.Format(daysPrecision))
	assert.Equal(t, "2023-03-02", cards[helloCardIndex].LastReviewDate.Format(daysPrecision))
	assert.Equal(t, 2, cards[helloCardIndex].NextReviewOffset)
	assert.Equal(t, 3, cards[helloCardIndex].RepetitionCount)
	assert.Equal(t, 2.113, cards[helloCardIndex].EF)
}
