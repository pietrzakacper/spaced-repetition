package csv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextToLines(t *testing.T) {
	textChan := make(chan []byte, 8)
	linesChan := TextToLines(textChan)

	textChan <- []byte("a 1\nb 2\n")
	textChan <- []byte("c")
	textChan <- []byte(" 3\n")
	textChan <- []byte("d ")
	textChan <- []byte("4")
	close(textChan)

	lines := make([]string, 0, 4)
	for line := range linesChan {
		lines = append(lines, line)
	}

	equal := reflect.DeepEqual(lines, []string{"a 1", "b 2", "c 3", "d 4"})

	if !equal {
		t.Fatalf("lines %v are not equal to expected result", lines)
	}
}

func TestParseLine(t *testing.T) {
	line := ParseLine("hello, hola")

	assert.Equal(t, "hello", line[0])
	assert.Equal(t, "hola", line[1])
}
