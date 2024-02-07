package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Prompt the user to input their API key
	apiKey := getUserInput("Enter your OpenAI API key: ")

	ctx := context.Background()

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
	}

	for {
		userMessage := getUserInput("You: ")
		if userMessage == "quit" {
			break
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: userMessage,
		})

		reqBody := &CompletionRequest{
			Model:    "gpt-3.5-turbo", // Specify the model here
			Messages: messages,
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatalln(err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}
		json.Unmarshal(respBytes, &result)

		assistantMessage := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println("Assistant: " + assistantMessage)

		messages = append(messages, Message{
			Role:    "assistant",
			Content: assistantMessage,
		})
	}
}

// getUserInput prompts the user for input and returns the entered value
func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
