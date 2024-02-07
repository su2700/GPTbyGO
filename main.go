// Import necessary packages
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

// Define the structure of a message
type Message struct {
	Role    string `json:"role"`    // Role can be 'system', 'user', or 'assistant'
	Content string `json:"content"` // Content is the text of the message
}

// Define the structure of a completion request
type CompletionRequest struct {
	Model    string    `json:"model"`    // Model is the ID of the language model to use
	Messages []Message `json:"messages"` // Messages is an array of message objects
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Prompt the user to input their API key
	apiKey := getUserInput("Enter your OpenAI API key: ")

	ctx := context.Background()

	// Initialize the conversation with a system message
	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
	}

	// Main loop for the interactive chat
	for {
		// Get the user's message
		userMessage := getUserInput("You: ")
		// If the user types 'quit', break the loop and end the program
		if userMessage == "quit" {
			break
		}
		// Add the user's message to the conversation
		messages = append(messages, Message{
			Role:    "user",
			Content: userMessage,
		})

		// Prepare the request body
		reqBody := &CompletionRequest{
			Model:    "gpt-3.5-turbo", // Specify the model here
			Messages: messages,
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatalln(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			log.Fatalln(err)
		}

		// Set the necessary headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		// Read the response body
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// Parse the response body
		var result map[string]interface{}
		json.Unmarshal(respBytes, &result)

		// Extract the assistant's message from the response
		assistantMessage := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println("Assistant: " + assistantMessage)

		// Add the assistant's message to the conversation
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
