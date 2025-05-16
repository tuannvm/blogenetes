package rss

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/tuannvm/blogenetes/modules/shared"
)

type RSS struct {
	parser *gofeed.Parser
}

// New creates a new RSS fetcher
func New() *RSS {
	return &RSS{
		parser: gofeed.NewParser(),
	}
}

// Fetch fetches and parses an RSS/Atom feed from the given URL
func (r *RSS) Fetch(ctx context.Context, url string) (*shared.RSSFeed, error) {
	// Parse the feed directly from URL
	feed, err := r.parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	// Convert to our shared types
	result := &shared.RSSFeed{
		Title:       feed.Title,
		Link:        feed.Link,
		Description: feed.Description,
		Items:       make([]shared.RSSItem, 0, len(feed.Items)),
	}

	// Convert feed items
	for _, item := range feed.Items {
		// Handle published date
		var pubDate time.Time
		switch {
		case item.PublishedParsed != nil:
			pubDate = *item.PublishedParsed
		case item.UpdatedParsed != nil:
			pubDate = *item.UpdatedParsed
		default:
			pubDate = time.Now()
		}

		// Use content if available, fallback to description
		content := item.Content
		if content == "" {
			content = item.Description
		}

		result.Items = append(result.Items, shared.RSSItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Content:     content,
			Published:   pubDate,
		})
	}

	return result, nil
}
