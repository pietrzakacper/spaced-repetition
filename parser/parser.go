package parser

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

func ParseLine(line string) (front string, back string) {
	return "", ""
}

func ParseCSV() bool {
	return false
}
