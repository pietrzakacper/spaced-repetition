package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func TextToLines(textChan <-chan string) <-chan string {
	linesChan := make(chan string)

	go func() {
		line := ""
		for textChunk := range textChan {
			for _, char := range textChunk {
				if string(char) == "\n" {
					linesChan <- line
					line = ""
					continue
				}

				line += string(char)
			}
		}

		if line != "" {
			linesChan <- line
		}

		close(linesChan)
	}()

	return linesChan
}

func ParseLine(line string) (front string, back string, err error) {
	splitArr := strings.Split(line, ",")

	if len(splitArr) != 2 {
		return "", "", errors.New("Incorrect number of elements in line")
	}

	for i, el := range splitArr {
		splitArr[i] = strings.Trim(el, " ")
	}

	return splitArr[0], splitArr[1], nil
}

type TwoColumnEntry struct {
	First  string
	Second string
}

func ParseCSVStream(textChannel <-chan string) []TwoColumnEntry {
	linesChannel := TextToLines(textChannel)

	entries := []TwoColumnEntry{}

	for line := range linesChannel {
		first, back, err := ParseLine(line)

		if err != nil {
			fmt.Printf("Couldn't parse line: %s \nError: %v", line, err)
		}

		entries = append(entries, TwoColumnEntry{first, back})
	}

	return entries
}

func ParseCSVLines(lines []string) []TwoColumnEntry {
	entries := make([]TwoColumnEntry, len(lines))

	for lineIndex, line := range lines {
		first, back, err := ParseLine(line)

		if err != nil {
			fmt.Printf("Couldn't parse line: %s \nError: %v", line, err)
		}

		entries[lineIndex] = TwoColumnEntry{first, back}
	}

	return entries
}

func FileToStream(filename string) chan string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open the file due to this %s error \n", err)
	}

	buff := make([]byte, 100)

	textChannel := make(chan string)

	go func() {
		for {
			// read content to buffer
			readTotal, err := file.Read(buff)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			fileContent := string(buff[:readTotal])

			textChannel <- fileContent
		}

		close(textChannel)

		file.Close()
	}()

	return textChannel
}
