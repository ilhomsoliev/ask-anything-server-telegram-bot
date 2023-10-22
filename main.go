package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"github.com/joho/godotenv"
	"log"
	"os"
	"bytes"
)
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
func main() {

	godotenv.Load()
	
	apiKey := os.Getenv("API_KEY")

	url := "https://api.aiguoguo199.com/v1/chat/completions" // Change URL as needed
	
	if apiKey == "" {
		log.Fatalln("Api Key is missing")
	}

	token := apiKey

	client := &http.Client{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w,"Hello")
	})
	
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		body1, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
		}

		questionActual := string(body1)

		payload := struct {
			Model    string    `json:"model"`
			Stream   bool      `json:"stream"`
			Messages []Message `json:"messages"`
		}{
			Model:  "gpt-3.5-turbo",
			Stream: false,
			Messages: []Message{
				{Role: "system", Content: questionActual},
				{Role: "user", Content: "Hello"},
			},
		}
		// Convert Go struct to JSON byte array
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		// Create a new POST request with the JSON data as the body
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to create request: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}

		// Set headers
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to make a request: %s", err)
			io.WriteString(w, "Some error on server side, please try again later")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %s", err)
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
		fmt.Println(content)	
		io.WriteString(w, content)
			
			/*body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("could not read body: %s\n", err)
			}

			questionActual := string(body)
			// ->
			apiKey := os.Getenv("API_KEY")
			if apiKey == "" {
				log.Fatalln("Api Key is missing")
			}
			ctx := context.Background()
			client := gpt3.NewClient(apiKey)
			question := &Question{}
			question.Question = questionActual
			fmt.Println(questionActual)
			var message = ""
			client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
				Prompt: []string{
					question.Question,
				},
				MaxTokens:   gpt3.IntPtr(3000),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				message1 := string(resp.Choices[0].Text)
				message = message + message1
			})
			io.WriteString(w, message)
			fmt.Println(message)*/
			
	})

	fmt.Println("Server has started")
	http.ListenAndServe(":8080", nil)
}
// Define the struct based on the JSON structure
type ChatCompletion struct {
	ID      string  `json:"id"`
	Object  string  `json:"object"`
	Created int64   `json:"created"`
	Model   string  `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage   `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	MessageResponse      MessageResponse `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type MessageResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens    int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
