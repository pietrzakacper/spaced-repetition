package persistance

import (
	"controller"
	"csv"
	"errors"
	"flashcard"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CSVPersistance struct {
}

type CSVStore struct {
	filepath string
}

func (p *CSVPersistance) Create(name string, userId string) controller.Store {
	cwd, _ := os.Getwd()
	filepath := filepath.Join(cwd, name+"_"+userId+".csv")

	if _, err := os.Stat(filepath); err != nil {
		os.Create(filepath)
	}

	return &CSVStore{filepath}
}

var daysPrecision = "2006-02-01"

func (s *CSVStore) ReadAll() []flashcard.Record {
	f, _ := os.Open(s.filepath)
	defer f.Close()

	stream := csv.TextToLines(csv.FileToChannel(f))

	records := make([]flashcard.Record, 0)

	for line := range stream {
		parts := strings.Split(line, ",")
		id := parts[0]
		front := parts[1]
		back := parts[2]
		creationDate, _ := time.Parse(daysPrecision, parts[3])
		lastReviewDate, _ := time.Parse(daysPrecision, parts[4])
		repetitionCount, _ := strconv.ParseInt(parts[5], 10, 64)
		nextReviewOffset, _ := strconv.ParseInt(parts[6], 10, 64)
		ef, _ := strconv.ParseFloat(parts[7], 64)

		deleted := false

		if len(parts) >= 9 {
			deleted, _ = strconv.ParseBool(parts[8])
		}

		record := flashcard.Record{
			Id:               id,
			Front:            front,
			Back:             back,
			CreationDate:     creationDate,
			LastReviewDate:   lastReviewDate,
			RepetitionCount:  int(repetitionCount),
			NextReviewOffset: int(nextReviewOffset),
			EF:               ef,
			Deleted:          deleted,
		}

		records = append(records, record)
	}

	return records
}

func (s *CSVStore) Find(cardId string) (flashcard.Record, error) {
	f, _ := os.Open(s.filepath)
	defer f.Close()

	stream := csv.TextToLines(csv.FileToChannel(f))

	for line := range stream {
		parts := strings.Split(line, ",")
		id := parts[0]

		if id == cardId {
			front := parts[1]
			back := parts[2]
			creationDate, _ := time.Parse(daysPrecision, parts[3])
			lastReviewDate, _ := time.Parse(daysPrecision, parts[4])
			repetitionCount, _ := strconv.ParseInt(parts[5], 10, 64)
			nextReviewOffset, _ := strconv.ParseInt(parts[6], 10, 64)
			ef, _ := strconv.ParseFloat(parts[7], 64)

			record := flashcard.Record{
				Id:               id,
				Front:            front,
				Back:             back,
				CreationDate:     creationDate,
				LastReviewDate:   lastReviewDate,
				RepetitionCount:  int(repetitionCount),
				NextReviewOffset: int(nextReviewOffset),
				EF:               ef,
			}

			return record, nil
		}
	}

	return flashcard.Record{}, errors.New("No card found for id: " + cardId)
}

func (s *CSVStore) Add(record *flashcard.Record) {
	f, _ := os.OpenFile(s.filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	defer f.Close()

	newId := uuid.New().String()

	entries := []string{
		newId,
		record.Front,
		record.Back,
		record.CreationDate.Format(daysPrecision),
		record.LastReviewDate.Format(daysPrecision),
		strconv.FormatInt(int64(record.RepetitionCount), 10),
		strconv.FormatInt(int64(record.NextReviewOffset), 10),
		strconv.FormatFloat(record.EF, 'e', 2, 64),
		"false",
	}

	line := csv.MakeLine(entries)

	_, err := f.WriteString(line)

	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}

func (s *CSVStore) Update(record *flashcard.Record) {
	f, _ := os.OpenFile(s.filepath, os.O_RDWR, 0644)
	defer f.Close()

	stream := csv.TextToLines(csv.FileToChannel(f))

	newContent := ""

	for line := range stream {
		parts := strings.Split(line, ",")
		idSlice := parts[0:1]

		var dataSlice = parts[1:]

		if idSlice[0] == record.Id {
			dataSlice = []string{
				record.Front,
				record.Back,
				record.CreationDate.Format(daysPrecision),
				record.LastReviewDate.Format(daysPrecision),
				strconv.FormatInt(int64(record.RepetitionCount), 10),
				strconv.FormatInt(int64(record.NextReviewOffset), 10),
				strconv.FormatFloat(record.EF, 'e', -1, 64),
				strconv.FormatBool(record.Deleted),
			}
		}

		newContent += csv.MakeLine(append(idSlice, dataSlice...))
	}

	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(newContent)
}
