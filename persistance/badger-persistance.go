package persistance

import (
	"controller"
	"encoding/json"
	"flashcard"
	"log"
	"os"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

type BadgerPersistance struct {
	inMemory bool
}

type BadgerStore struct {
	db *badger.DB
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

func (f *FlashcardSerialized) toRecord() flashcard.Record {
	creationDate := time.UnixMicro(f.CreationDate)
	lastReviewDate := time.UnixMicro(f.LastReviewDate)

	return flashcard.Record{
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

func (p *BadgerPersistance) Create(name string, userId string) controller.Store {
	cwd, _ := os.Getwd()

	var opts badger.Options

	if p.inMemory {
		opts = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opts = badger.DefaultOptions(cwd + "/badger-" + name)
	}

	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal(err)
	}

	return &BadgerStore{db}
}

func (b *BadgerStore) ReadAll() []flashcard.Record {
	records := make([]flashcard.Record, 0)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			err := it.Item().Value(func(v []byte) error {
				f := &FlashcardSerialized{}
				err := json.Unmarshal(v, f)

				if err != nil {
					return err
				}

				records = append(records, f.toRecord())

				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return []flashcard.Record{}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreationDate.UnixMicro() < records[j].CreationDate.UnixMicro()
	})

	return records
}

func (b *BadgerStore) Add(record *flashcard.Record) {
	newId := uuid.New().String()

	b.setRecord(newId, record)
}

func (b *BadgerStore) Update(record *flashcard.Record) {
	b.setRecord(record.Id, record)
}
func (s *BadgerStore) Find(cardId string) (flashcard.Record, error) {
	f := &FlashcardSerialized{}

	err := s.db.View(func(txn *badger.Txn) error {
		record, err := txn.Get([]byte(cardId))

		if err != nil {
			return err
		}

		err = record.Value(
			func(v []byte) error {
				err := json.Unmarshal(v, f)

				if err != nil {
					return err
				}

				return nil
			})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return f.toRecord(), err
	}

	return f.toRecord(), nil
}

func (b *BadgerStore) setRecord(id string, record *flashcard.Record) {
	flashcardSerialized := &FlashcardSerialized{
		Id:               id,
		Front:            record.Front,
		Back:             record.Back,
		CreationDate:     record.CreationDate.UnixMicro(),
		LastReviewDate:   record.LastReviewDate.UnixMicro(),
		RepetitionCount:  record.RepetitionCount,
		NextReviewOffset: record.NextReviewOffset,
		EF:               record.EF,
		Deleted:          record.Deleted,
	}

	flashcardSerializedStr, err := json.Marshal(flashcardSerialized)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(id), []byte(flashcardSerializedStr))

		if err != nil {
			log.Fatal(err)
			return err
		}

		err = txn.Commit()

		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
}
