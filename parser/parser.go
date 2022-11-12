package parser

import (
	"fmt"
	"io"
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

func ParseLine(line string) []string {
	splitArr := strings.Split(line, ",")

	for i, el := range splitArr {
		splitArr[i] = strings.Trim(el, " ")
	}

	return splitArr
}

func ParseCSVStream(textStream io.Reader) <-chan []string {
	linesChannel := TextToLines(FileToChannel(textStream))

	entryChan := make(chan []string)

	columnLine := ParseLine(<-linesChannel)

	fmt.Printf("Column names are: %s %s\n", columnLine[0], columnLine[1])

	go func() {
		for line := range linesChannel {
			entryChan <- ParseLine(line)
		}

		close(entryChan)
	}()

	return entryChan
}

func ParseCSVLines(lines []string) [][]string {
	entries := make([][]string, len(lines))

	for lineIndex, line := range lines {
		entries[lineIndex] = ParseLine(line)
	}

	return entries
}

func FileToChannel(file io.Reader) chan string {
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
	}()

	return textChannel
}
