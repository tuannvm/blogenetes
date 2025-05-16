package main

import (
	"context"
	"fmt"
	"log"

	"dagger/rss/internal/dagger"

	"github.com/mmcdole/gofeed"
)

// Alias the internal Dagger client type for clarity
type DaggerClient = dagger.Client

// RSSFeed is a minimal wrapper for gofeed.Feed that's compatible with Dagger
type RSSFeed struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Link        string    `json:"link,omitempty"`
	Items       []RSSItem `json:"items,omitempty"`
}

// RSSItem is a minimal wrapper for gofeed.Item that's compatible with Dagger
type RSSItem struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
	Content     string `json:"content,omitempty"`
	Published   string `json:"published,omitempty"`
}

type RSS struct {
	dag *DaggerClient
}

func NewRSS(dag *dagger.Client) *RSS {
	return &RSS{
		dag: dag,
	}
}

// Fetch fetches and parses an RSS/Atom feed from the given URL
// This is the main entry point that will be called from the Dagger pipeline
func (r *RSS) Fetch(ctx context.Context, url string) (*RSSFeed, error) {
	// Log the start of the fetch operation
	log.Printf("Fetching RSS feed from: %s", url)

	// Create a new parser for each request
	parser := gofeed.NewParser()

	// Parse the feed directly from URL
	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	log.Printf("Processing %d items from feed: %s", len(feed.Items), feed.Title)

	// Convert to our minimal wrapper
	result := &RSSFeed{
		Title:       feed.Title,
		Description: feed.Description,
		Link:        feed.Link,
		Items:       make([]RSSItem, 0, len(feed.Items)),
	}

	// Convert items to our minimal wrapper
	for _, item := range feed.Items {
		result.Items = append(result.Items, RSSItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			Content:     item.Content,
			Published:   item.Published,
		})
	}

	return result, nil
}
