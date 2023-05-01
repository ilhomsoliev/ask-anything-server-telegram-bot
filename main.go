package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"io"
	"net/http"
	"os"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w,"Hello")
	})
	
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		
			body, err := ioutil.ReadAll(r.Body)
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
			fmt.Println(message)
			
	})

	fmt.Println("Server has started")
	http.ListenAndServe(":8080", nil)
}

type Question struct {
	Question string `json:"question"`
}

type Answer struct {
	Answer string `json:"answer"`
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
