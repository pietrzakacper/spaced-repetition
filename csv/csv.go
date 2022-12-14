package csv

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func TextToLines(textChan <-chan []byte) <-chan string {
	linesChan := make(chan string)

	go func() {
		line := make([]byte, 0)
		for textChunk := range textChan {
			for _, char := range textChunk {
				if string(char) == "\n" {
					linesChan <- string(line)
					line = make([]byte, 0)
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

func ParseCSVStream(textStream io.Reader, extractColumns []string) (<-chan []string, error) {
	linesChannel := TextToLines(FileToChannel(textStream))

	entryChan := make(chan []string)

	columnPositions, err := findColumnPositions(extractColumns, <-linesChannel)

	if err != nil {
		return nil, err
	}

	go func() {
		for line := range linesChannel {
			entry := ParseLine(line)

			if len(entry) != len(extractColumns) {
				fmt.Println("Skipping entry from line: ", line)
				continue
			}

			onlyRequiredColsEntry := make([]string, len(extractColumns))

			for entryIndex, properColIndex := range columnPositions {
				onlyRequiredColsEntry[properColIndex] = entry[entryIndex]
			}

			entryChan <- onlyRequiredColsEntry
		}

		close(entryChan)
	}()

	return entryChan, nil
}

func findColumnPositions(extractColumns []string, columnLine string) (map[int]int, error) {
	columnEntries := ParseLine(columnLine)

	columnPositions := make(map[int]int)

	for columnIndex, columnName := range columnEntries {
		extractColumnsIndex := indexOf(columnName, extractColumns)

		if extractColumnsIndex > -1 {
			columnPositions[columnIndex] = extractColumnsIndex
		}
	}

	if len(columnPositions) < len(extractColumns) {
		return nil, errors.New("Cannot find required columns")
	}

	return columnPositions, nil
}

func indexOf(elToFind string, arr []string) int {
	for index, el := range arr {
		if strings.ToLower(el) == strings.ToLower(elToFind) {
			return index
		}
	}

	return -1
}

func FileToChannel(file io.Reader) chan []byte {

	textChannel := make(chan []byte)

	go func() {
		for {
			buff := make([]byte, 100)
			// read content to buffer
			readTotal, err := file.Read(buff)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			fileContent := buff[:readTotal]

			textChannel <- fileContent
		}

		close(textChannel)
	}()

	return textChannel
}
