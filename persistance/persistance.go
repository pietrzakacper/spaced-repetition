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
	storePath := filepath.Join(cwd, "../", p.store)

	stream := parser.TextToLines(parser.FileToStream(storePath))

	lines := make([]string, 0)

	for line := range stream {
		lines = append(lines, line)
	}

	return lines
}

func (p *Persistance) Add(line string) {
	cwd, _ := os.Getwd()
	storePath := filepath.Join(cwd, "../", p.store)

	f, _ := os.OpenFile(storePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	defer f.Close()

	w, err := f.WriteString(line + "\n")

	fmt.Println(w, err)
}
