package persistance

import (
	"os"
	"parser"
	"path/filepath"
)

type Persistance struct {
	store string
	lines []string
}

func Create(store string) Persistance {
	return Persistance{store: store, lines: make([]string, 0)}
}

func (p *Persistance) Read() []string {
	cwd, _ := os.Getwd()
	storePath := filepath.Join(cwd, "../", p.store)

	stream := parser.TextToLines(parser.FileToStream(storePath))

	for line := range stream {
		p.lines = append(p.lines, line)
	}

	return p.lines
}
