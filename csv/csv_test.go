package csv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextToLines(t *testing.T) {
	textChan := make(chan string, 8)
	linesChan := TextToLines(textChan)

	textChan <- "a 1\nb 2\n"
	textChan <- "c"
	textChan <- " 3\n"
	textChan <- "d "
	textChan <- "4"
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
