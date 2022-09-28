package parser

import (
	"reflect"
	"testing"
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
