package csv

import (
	"fmt"
	"io"
	"strings"
)

func TextToLines(textChan <-chan string) <-chan string {
	linesChan := make(chan string)

	go func() {
		line := make([]rune, 0)
		for textChunk := range textChan {
			for _, char := range textChunk {
				if string(char) == "\n" {
					linesChan <- string(line)
					line = make([]rune, 0)
					continue
				}

				line = append(line, char)
			}
		}

		if len(line) != 0 {
			linesChan <- string(line)
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

func MakeLine(entries []string) string {
	return strings.Join(entries, ",") + "\n"
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
