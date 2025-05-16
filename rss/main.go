package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger/rss/internal/dagger"

	"github.com/mmcdole/gofeed"
)

type RSS struct {
	dag *dagger.Client
}

func NewRSS(dag *dagger.Client) *RSS {
	return &RSS{
		dag: dag,
	}
}

// Fetch fetches and parses an RSS/Atom feed from the given URL and returns it as JSON
func (r *RSS) Fetch(ctx context.Context, url string) (string, error) {
	// Create a new parser for each request
	parser := gofeed.NewParser()

	// Parse the feed directly from URL
	feed, err := parser.ParseURL(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse feed: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal feed to JSON: %w", err)
	}

	return string(jsonData), nil
}
