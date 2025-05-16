package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
	"dagger/rss/internal/dagger"
)

// Alias the internal Dagger client type for clarity
type DaggerClient = dagger.Client

// RSSFeed aliases gofeed.Feed for better semantics
type RSSFeed = gofeed.Feed

// RSSItem aliases gofeed.Item for better semantics
type RSSItem = gofeed.Item

type RSS struct {
	dag    *DaggerClient
	parser *gofeed.Parser
}

// New creates a new RSS fetcher with the given Dagger client
func New(ctx context.Context, client *DaggerClient) (*RSS, error) {
	if client == nil {
		return nil, fmt.Errorf("dagger client cannot be nil")
	}

	return &RSS{
		dag:    client,
		parser: gofeed.NewParser(),
	}, nil
}

// Fetch fetches and parses an RSS/Atom feed from the given URL
// This is the main entry point that will be called from the Dagger pipeline
func (r *RSS) Fetch(ctx context.Context, url string) (*RSSFeed, error) {
	// Log the start of the fetch operation
	log.Printf("Fetching RSS feed from: %s", url)

	// Parse the feed directly from URL
	feed, err := r.parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	log.Printf("Processing %d items from feed: %s", len(feed.Items), feed.Title)

	// Return the parsed feed directly
	result := feed

	return result, nil
}
