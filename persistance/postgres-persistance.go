package persistance

import (
	"controller"
	"database/sql"
	"flashcard"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type PostgresPersistance struct {
}

type PostgresStore struct {
	db     *sql.DB
	userId string
}

func (p *PostgresPersistance) Create(name string, userId string) controller.Store {
	// @TODO move outside
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresStore{db, userId}
}

func (p *PostgresStore) ReadAll() []flashcard.Record {
	records := make([]flashcard.Record, 0)

	rows, err := p.db.Query("SELECT id, front, back, repetition_count, next_review_offset, ef, deleted, creation_date, last_review_date FROM flashcards WHERE user_id=$1", p.userId)

	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
		return records
	}

	var id, front, back string
	var creationDate, lastReviewDate time.Time
	var nextReviewOffset, repetitionCount int
	var ef float64
	var deleted bool

	for rows.Next() {
		rows.Scan(&id, &front, &back, &repetitionCount, &nextReviewOffset, &ef, &deleted, &creationDate, &lastReviewDate)
		records = append(records, flashcard.Record{
			Id:               id,
			Front:            front,
			Back:             back,
			CreationDate:     creationDate,
			LastReviewDate:   lastReviewDate,
			NextReviewOffset: nextReviewOffset,
			EF:               ef,
			Deleted:          deleted,
		})
	}

	return records
}

func (p *PostgresStore) Add(record *flashcard.Record) {

}

func (p *PostgresStore) Update(record *flashcard.Record) {
}

func (s *PostgresStore) Find(cardId string) (flashcard.Record, error) {
	f := &flashcard.FlashcardSerialized{}

	return f.ToRecord(), nil
}
