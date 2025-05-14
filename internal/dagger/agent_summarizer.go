package dagger

import (
	"context"
	"fmt"

	go_openai "github.com/sashabaranov/go-openai"
)

// AgentSummarizer summarizes articles using an LLM.
type AgentSummarizer struct {
	Client *go_openai.Client
}

// NewAgentSummarizer creates a new summarization agent.
func NewAgentSummarizer(client *go_openai.Client) *AgentSummarizer {
	return &AgentSummarizer{Client: client}
}

// Summarize generates a summary for the given article text.
func (a *AgentSummarizer) Summarize(ctx context.Context, article string) (string, error) {
	resp, err := a.Client.CreateChatCompletion(ctx, go_openai.ChatCompletionRequest{
		Model:    go_openai.GPT3Dot5Turbo,
		Messages: []go_openai.ChatCompletionMessage{
			{Role: go_openai.ChatMessageRoleSystem, Content: "You are a helpful assistant that summarizes articles."},
			{Role: go_openai.ChatMessageRoleUser, Content: article},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}
	return resp.Choices[0].Message.Content, nil
}
