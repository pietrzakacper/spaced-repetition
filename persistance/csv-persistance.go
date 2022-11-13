package persistance

import (
	"controller"
	"fmt"
	"os"
	"parser"
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

func (p *CSVPersistance) Create(name string) controller.Store {
	cwd, _ := os.Getwd()
	filepath := filepath.Join(cwd, "../", name+".csv")

	return &CSVStore{filepath}
}

var daysPrecision = "2006-02-01"

func (s *CSVStore) ReadAll() []controller.FlashcardRecord {
	f, _ := os.Open(s.filepath)
	defer f.Close()

	stream := parser.TextToLines(parser.FileToChannel(f))

	records := make([]controller.FlashcardRecord, 0)

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

		record := controller.FlashcardRecord{
			Id:               id,
			Front:            front,
			Back:             back,
			CreationDate:     creationDate,
			LastReviewDate:   lastReviewDate,
			RepetitionCount:  int(repetitionCount),
			NextReviewOffset: int(nextReviewOffset),
			EF:               ef,
		}

		records = append(records, record)
	}

	return records
}

func (s *CSVStore) Add(record *controller.FlashcardRecord) {
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
	}

	line := parser.MakeLine(entries)

	_, err := f.WriteString(line)

	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}

func (s *CSVStore) Update(record *controller.FlashcardRecord) {
	f, _ := os.OpenFile(s.filepath, os.O_RDWR, 0644)
	defer f.Close()

	stream := parser.TextToLines(parser.FileToChannel(f))

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
			}
		}

		newContent += parser.MakeLine(append(idSlice, dataSlice...))
	}

	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(newContent)
}
