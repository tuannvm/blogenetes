package dagger

import (
	"context"
	"log"
	"os"

	dagger "dagger.io/dagger"

	"github.com/joho/godotenv"
	go_openai "github.com/sashabaranov/go-openai"
)

// SetupOpenAIClient loads environment variables and initializes the OpenAI client.
func SetupOpenAIClient() *go_openai.Client {
	_ = godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}
	return go_openai.NewClient(apiKey)
}

// SetupDaggerClient connects to the Dagger engine and returns the Dagger client.
func SetupDaggerClient() (*dagger.Client, error) {
	_ = godotenv.Load()
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return nil, err
	}
	return client, nil
}
