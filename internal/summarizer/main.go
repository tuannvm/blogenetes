package main

import (
	"context"
	"fmt"
	"os"
	go_openai "github.com/sashabaranov/go-openai"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: summarizer <article-file>")
		os.Exit(1)
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	article := string(data)
	client := go_openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(context.Background(), go_openai.ChatCompletionRequest{
		Model:    go_openai.GPT3Dot5Turbo,
		Messages: []go_openai.ChatCompletionMessage{
			{Role: go_openai.ChatMessageRoleSystem, Content: "You are a helpful assistant that summarizes articles in a detail manner."},
			{Role: go_openai.ChatMessageRoleUser, Content: article},
		},
	})
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		fmt.Println("no completion choices returned")
		os.Exit(1)
	}
	fmt.Print(resp.Choices[0].Message.Content)
}
