package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {

	godotenv.Load()

	apiKey := os.Getenv("API_KEY")

	/*	url := os.Getenv("URL_KEY")

		if url == "" {
			log.Fatalln("URL Key is missing")
		}*/
	if apiKey == "" {
		log.Fatalln("Api Key is missing")
	}

	//token := apiKey

	//client := &http.Client{}

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with most use cases
	model := client.GenerativeModel("gemini-1.5-flash")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	})

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {

		body1, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		questionActual := string(body1)

		fmt.Println("Question asked: " + questionActual)

		/*payload := struct {
			Model    string    `json:"model"`
			Stream   bool      `json:"stream"`
			Messages []Message `json:"messages"`
		}{
			Model:  "gpt-3.5-turbo",
			Stream: false,
			Messages: []Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: questionActual},
			},
		}

		// Convert Go struct to JSON byte array
		jsonData, err := json.Marshal(payload)
		if err != nil {
			// log.Fatalf("Failed to marshal JSON: %s", err)
			fmt.Println("Failed to marshal JSON: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		// Create a new POST request with the JSON data as the body
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			// log.Fatalf("Failed to create request: %s", err)
			fmt.Println("Failed to create request: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		// Set headers
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			// log.Fatalf("Failed to make a request: %s", err)
			fmt.Println("Failed to make a request: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}
		defer resp.Body.Close()
		*/
		resp, err := model.GenerateContent(ctx, genai.Text(questionActual))
		if err != nil {
			io.WriteString(w, "Some error on server side, please try again later")
			log.Fatal(err)
		}
		if resp != nil {
			candidates := resp.Candidates
			if candidates != nil {
				for _, candidate := range candidates {
					content := candidate.Content
					if content != nil {
						text := content.Parts[0]

						stringAsString := fmt.Sprintf("%v", text) //print(text)
						log.Println("Candidate text:", stringAsString)
						io.WriteString(w, stringAsString)
					} else {
						io.WriteString(w, "Some error on server side, please try again later")
					}
				}
			} else {
				io.WriteString(w, "Some error on server side, please try again later")
			}
		} else {
			io.WriteString(w, "Some error on server side, please try again later")
		}
		/*body, err := ioutil.ReadAll()
		if err != nil {
			// log.Fatalf("Failed to read response body: %s", err)
			fmt.Println("Failed to read response body: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}
		data := string(body)
		// Parse Response
		var chatCompletion ChatCompletion
		if err := json.Unmarshal([]byte(data), &chatCompletion); err != nil {
			log.Fatal("Error unmarshaling JSON:", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		// Get the content from the first choice (assuming there's always at least one choice)
		content := chatCompletion.Choices[0].MessageResponse.Content
		fmt.Println(content)

		io.WriteString(w, content)*/
	})

	fmt.Println("Server has started")
	http.ListenAndServe(":8080", nil)
}

// Define the struct based on the JSON structure
type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index           int             `json:"index"`
	MessageResponse MessageResponse `json:"message"`
	FinishReason    string          `json:"finish_reason"`
}

type MessageResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
