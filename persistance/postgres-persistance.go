package persistance

import (
	"context"
	"controller"
	"errors"
	"flashcard"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresPersistance struct {
}

type PostgresStore struct {
	db     *pgxpool.Pool
	userId string
}

var db *pgxpool.Pool

func (p *PostgresPersistance) Create(userId string) controller.Store {
	if db == nil {
		config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal(err)
		}
		config.MaxConns = 30
		config.MinConns = 3
		config.MaxConnLifetime = 10 * time.Minute
		config.MaxConnIdleTime = 55 * time.Second

		pool, err := pgxpool.ConnectConfig(context.Background(), config)

		if err != nil {
			log.Fatal(err)
		}

		db = pool
		go func() {
			for {
				// in short fly postgres uses HAProxy with 60 second TCP timeout
				time.Sleep(40 * time.Second)
				// let's make sure we keep the connections warm
				err := pool.Ping(context.Background())
				log.Printf(
					"ping err: %v, acq conns: %v, idle conns: %v, total conns: %v",
					err,
					pool.Stat().AcquiredConns(),
					pool.Stat().IdleConns(),
					pool.Stat().TotalConns(),
				)
			}
		}()
	}

	return &PostgresStore{db, userId}
}

func (p *PostgresStore) readAllImpl() ([]flashcard.Record, error) {
	records := make([]flashcard.Record, 0)

	rows, err := p.db.Query(context.Background(), `
		SELECT
		id, front, back, repetition_count, next_review_offset, ef, deleted, creation_date, last_review_date
		FROM flashcards
		WHERE user_id=$1`,
		p.userId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := flashcard.Record{}

		err = rows.Scan(&r.Id, &r.Front, &r.Back, &r.RepetitionCount, &r.NextReviewOffset, &r.EF, &r.Deleted, &r.CreationDate, &r.LastReviewDate)
		if err != nil {
			break
		}
		records = append(records, r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreationDate.UnixMicro() < records[j].CreationDate.UnixMicro()
	})

	return records, nil
}

func (p *PostgresStore) ReadAll() []flashcard.Record {
	defer timeTrack(time.Now(), "postgres.ReadAll: "+p.userId)

	tries := 0
	for {
		if tries >= 3 {
			break
		}

		tries++
		records, err := p.readAllImpl()

		if err != nil && err.Error() == "unexpected EOF" {
			log.Printf("=== EOF ===\n")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if err != nil {
			log.Printf("Error in ReadAll: %v", err)
			break
		}

		return records
	}

	log.Println("ReadAll failed...")

	return make([]flashcard.Record, 0)
}

const cards_limit_per_user = 20000

func (p *PostgresStore) Add(r *flashcard.Record) {
	defer timeTrack(time.Now(), "postgres.Add: "+r.Front)

	_, err := p.db.Exec(context.Background(), `
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (p *PostgresStore) Update(r *flashcard.Record) {
	defer timeTrack(time.Now(), "postgres.Update: "+r.Front)
	_, err := p.db.Exec(context.Background(), `
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
	rows, err := p.db.Query(context.Background(), `
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
	rows, err := p.db.Query(context.Background(), `SELECT user_id FROM sessions WHERE token=$1`, token)
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
	_, err := p.db.Exec(context.Background(), `
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
