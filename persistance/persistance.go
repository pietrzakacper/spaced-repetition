package persistance

import (
	"fmt"
	"os"
	"parser"
	"path/filepath"
)

type Persistance struct {
	store string
}

func Create(store string) Persistance {
	return Persistance{store: store}
}

func (p *Persistance) Read() []string {
	cwd, _ := os.Getwd()
	storePath := filepath.Join(cwd, "../", p.store+".csv")

	f, _ := os.Open(storePath)

	defer f.Close()

	stream := parser.TextToLines(parser.FileToChannel(f))

	lines := make([]string, 0)

	for line := range stream {
		lines = append(lines, line)
	}

	return lines
}

func (p *Persistance) Add(line string) {
	cwd, _ := os.Getwd()
	storePath := filepath.Join(cwd, "../", p.store+".csv")
	f, _ := os.OpenFile(storePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	defer f.Close()

	_, err := f.WriteString(line + "\n")

	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}
