package dagger

import (
	"context"
	"os"

	dagger "dagger.io/dagger"
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
func (a *AgentSummarizer) Summarize(ctx context.Context, client *dagger.Client, article string) (string, error) {
	// Write the article to a file in a Dagger directory
	dir := client.Directory().WithNewFile("article.txt", article)

	// Mount the source code and run a Go summarizer binary inside a container
	openaiSecret := client.SetSecret("OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY"))
container := client.Container().
		From("golang:1.24").
		WithMountedDirectory("/src", client.Host().Directory("."), dagger.ContainerWithMountedDirectoryOpts{Owner: "root"}).
		WithMountedDirectory("/work", dir).
		WithWorkdir("/src/internal/summarizer").
		WithExec([]string{"go", "build", "-o", "/tmp/summarizer", "main.go"}).
		WithWorkdir("/work").
		WithSecretVariable("OPENAI_API_KEY", openaiSecret).
		WithExec([]string{"/tmp/summarizer", "article.txt"})

	// Capture the output (summary)
	summary, err := container.Stdout(ctx)
	if err != nil {
		return "", err
	}
	return summary, nil
}
