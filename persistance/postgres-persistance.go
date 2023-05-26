package persistance

import (
	"controller"
	"database/sql"
	"flashcard"
	"log"
	"os"
	"sort"
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

	rows, err := p.db.Query(`
		SELECT
		id, front, back, repetition_count, next_review_offset, ef, deleted, creation_date, last_review_date
		FROM flashcards
		WHERE user_id=$1`,
		p.userId,
	)

	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		r := flashcard.Record{}

		rows.Scan(&r.Id, &r.Front, &r.Back, &r.RepetitionCount, &r.NextReviewOffset, &r.EF, &r.Deleted, &r.CreationDate, &r.LastReviewDate)
		records = append(records, r)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreationDate.UnixMicro() < records[j].CreationDate.UnixMicro()
	})

	return records
}

func (p *PostgresStore) Add(r *flashcard.Record) {
	_, err := p.db.Exec(`
		INSERT into flashcards
		(id, creation_date, user_id, front, back, repetition_count, next_review_offset, ef, deleted, last_review_date)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		time.Now(),
		p.userId,
		r.Front,
		r.Back,
		r.RepetitionCount,
		r.NextReviewOffset,
		r.EF,
		r.Deleted,
		r.LastReviewDate,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func (p *PostgresStore) Update(r *flashcard.Record) {
	_, err := p.db.Exec(`
		UPDATE flashcards
		SET front=$1, back=$2, repetition_count=$3, next_review_offset=$4, ef=$5, deleted=$6, last_review_date=$7
		WHERE id=$8`,
		r.Front,
		r.Back,
		r.RepetitionCount,
		r.NextReviewOffset,
		r.EF,
		r.Deleted,
		r.LastReviewDate,
		r.Id,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func (p *PostgresStore) Find(cardId string) (flashcard.Record, error) {
	rows, err := p.db.Query(`
		SELECT
		id, front, back, repetition_count, next_review_offset, ef, deleted, creation_date, last_review_date
		FROM flashcards
		WHERE id=$1`,
		cardId,
	)

	if err != nil || !rows.Next() {
		log.Fatalln(err)
	}

	r := flashcard.Record{}

	rows.Scan(&r.Id, &r.Front, &r.Back, &r.RepetitionCount, &r.NextReviewOffset, &r.EF, &r.Deleted, &r.CreationDate, &r.LastReviewDate)

	return r, nil
}
