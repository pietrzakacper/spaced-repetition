package persistance

import (
	"controller"
	"database/sql"
	"errors"
	"flashcard"
	"log"
	"os"
	"sort"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresPersistance struct {
}

type PostgresStore struct {
	db     *sql.DB
	userId string
}

var db *sql.DB

func (p *PostgresPersistance) Create(userId string) controller.Store {
	if db == nil {
		var err error
		db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(0)

		if err != nil {
			log.Fatal(err)
		}
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

const cards_limit_per_user = 20000

func (p *PostgresStore) Add(r *flashcard.Record) {
	_, err := p.db.Exec(`
		INSERT into flashcards
		(id, creation_date, user_id, front, back, repetition_count, next_review_offset, ef, deleted, last_review_date)
		SELECT gen_random_uuid(), CURRENT_TIMESTAMP, $1, $2, $3, $4, $5, $6, $7, $8
		WHERE (SELECT COUNT(*) FROM flashcards WHERE user_id=$1) < $9`,
		p.userId,
		r.Front,
		r.Back,
		r.RepetitionCount,
		r.NextReviewOffset,
		r.EF,
		r.Deleted,
		r.LastReviewDate,
		cards_limit_per_user,
	)

	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
	defer rows.Close()

	if err != nil || !rows.Next() {
		log.Println(err)
	}

	r := flashcard.Record{}

	rows.Scan(&r.Id, &r.Front, &r.Back, &r.RepetitionCount, &r.NextReviewOffset, &r.EF, &r.Deleted, &r.CreationDate, &r.LastReviewDate)

	return r, nil
}

func (p *PostgresStore) FindUserIdByToken(token string) (string, error) {
	rows, err := p.db.Query(`SELECT user_id FROM sessions WHERE token=$1`, token)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return "", errors.New("Unknown error")
	}

	if !rows.Next() {
		return "", errors.New("No token found")
	}

	var userId string
	rows.Scan(&userId)

	if len(userId) > 0 {
		return userId, nil
	}

	return "", errors.New("user_id empty")
}

func (p *PostgresStore) UpsertSession(token string, userId string) {
	_, err := p.db.Exec(`
		INSERT INTO sessions (token, user_id)
		VALUES($1, $2) 
		ON CONFLICT (user_id) 
		DO 
		UPDATE SET token = $1;
	`, token, userId)

	if err != nil {
		log.Println(err)
		return
	}
}
