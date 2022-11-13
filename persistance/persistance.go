package persistance

import (
	"fmt"
	"os"
	"parser"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Persistance struct {
	storePath string
}

var idSeparator = "___"

func Create(store string) Persistance {
	cwd, _ := os.Getwd()
	storePath := filepath.Join(cwd, "../", store+".csv")

	return Persistance{storePath: storePath}
}

type Record struct {
	persistance *Persistance
	id          string
	Data        string
}

func (p *Persistance) Read() []*Record {
	f, _ := os.Open(p.storePath)

	defer f.Close()

	stream := parser.TextToLines(parser.FileToChannel(f))

	records := make([]*Record, 0)

	for line := range stream {
		parts := strings.Split(line, idSeparator)
		id := parts[0]
		data := parts[1]
		records = append(records, &Record{id: id, Data: data, persistance: p})
	}

	return records
}

func (p *Persistance) Add(data string) string {
	f, _ := os.OpenFile(p.storePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	defer f.Close()

	id := uuid.New().String()

	_, err := f.WriteString(id + idSeparator + data + "\n")

	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}

	return id
}

func (r *Record) Save() {
	f, _ := os.OpenFile(r.persistance.storePath, os.O_RDWR, 0644)

	defer f.Close()

	stream := parser.TextToLines(parser.FileToChannel(f))

	newContent := ""

	for line := range stream {
		parts := strings.Split(line, idSeparator)
		id := parts[0]

		var data string

		if id == r.id {
			data = r.Data
		} else {
			data = parts[1]
		}

		newContent += id + idSeparator + data + "\n"
	}

	// @TODO that's turbo inefficient, but maybe it's enough
	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(newContent)
}
