package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	apiKey := os.Getenv("API_KEY")
	url := "https://api.deepseek.com/v1/chat/completions"

	messages := []map[string]string{}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		messages = append(messages, map[string]string{"role": "user", "content": input})

		payload := map[string]any{
			"model":    "deepseek-chat",
			"messages": messages,
			"stream":   false,
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var result DeepSeekResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			return
		}

		if len(result.Choices) > 0 {
			fmt.Printf("AI: %s\n\n", result.Choices[0].Message.Content)
			messages = append(messages, map[string]string{"role": "assistant", "content": result.Choices[0].Message.Content})
		} else {
			fmt.Println("No response content received")
		}
	}
}
