package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Println("Too few arguments!")
		fmt.Println("Usage <path to zoom dir> <filterKeyword> <openAiToken> <spaced-rep-url> <space-rep-cookie>")
		return
	}

	pathToZoom, filterKeyword, openAiToken, spacedRepUrl, spacedRepCookie :=
		os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]

	onNewMatchingFile := func(path string) {
		content, _ := os.ReadFile(path)
		fmt.Printf("File content: %s\n", content)

		csvContent := transformChatToFlashcardsCSV(string(content), openAiToken)
		fmt.Printf("CSV content from OpenAI: \n---\n%s\n---\n", csvContent)

		csv := "back,front\n" + csvContent

		uploadCSVToSpacedRepetition(csv, spacedRepUrl, spacedRepCookie)
	}

	watchDir(pathToZoom, filterKeyword, onNewMatchingFile)
}

func watchDir(path, filterKeyword string, onNewMatchingFile func(path string)) {
	w := watcher.New()

	w.FilterOps(watcher.Create)

	rx := regexp.MustCompile(filterKeyword)

	go func() {
		fmt.Println("watcher started")
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				fmt.Println(event.Path)
				if f, _ := os.Stat(event.Path); f.IsDir() {
					fmt.Println("Directory detected")
					w.AddRecursive(event.Path)
					continue
				}

				if rx.MatchString(event.Path) {
					onNewMatchingFile(event.Path)

				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				fmt.Println("watcher closed")
				return
			}
		}
	}()

	if err := w.AddRecursive(path); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Second); err != nil {
		log.Fatalln(err)
	}
}

var client *openai.Client

func transformChatToFlashcardsCSV(chat string, openAiToken string) string {
	if client == nil {
		client = openai.NewClient(openAiToken)
	}

	resp, err := client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt(chat),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

var prompt = func(chat string) string {
	return fmt.Sprintf(`I have a language learning transcript file and need to convert it into a CSV format suitable for Anki flashcards. Each line in the transcript consists of a phrase in a foreign language and its English translation, separated by a timestamp and speaker identification.
 For the CSV:
 
 Include only two columns for each entry: the foreign language phrase and its English translation.
 Exclude any Polish phrases and their translations.
 Omit timestamps, speaker names, and any non-phrase text.
 Do not include headers or additional text in the CSV.
 Separate each phrase pair with a semicolon and enclose each pair in quotes.
 Example Output:
 "tan grande como pensaba","as big as I thought"
 "muy raras veces","very rarely"
 "el papel","the role"
 "actuar","to act"
 
 Transcript:
 
 %s
 
 Note that I only want you to respond with a csv formated date. Example of the perfect response is:
 "tan grande como pensaba","as big as I thought"
 "muy raras veces","very rarely"
 "el papel","the role"
 "actuar","to act"
 "las botas","the boots"
 
 Start your response with data.`, chat)
}

func uploadCSVToSpacedRepetition(csv, spacedRepUrl, spacedRepCookie string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/import", spacedRepUrl), strings.NewReader(csv))

	if err != nil {
		log.Fatalf("Error constructing request: %v\n", err)
	}

	req.Header.Add("Content-Type", "text/csv")
	req.Header.Add("Cookie", spacedRepCookie)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Error uploading CSV to Spaced Repetition: %v\n", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Error uploading CSV to Spaced Repetition: %v\n", res.Status)
	}

	res.Body.Close()
}
