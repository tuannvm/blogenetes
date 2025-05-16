package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/tuannvm/blogenetes/modules/shared"
)

type RSS struct{}

// Fetch fetches and parses an RSS feed from the given URL
func (r *RSS) Fetch(ctx context.Context, url string) (*shared.RSSFeed, error) {
	// Fetch the RSS feed
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the RSS feed
	var rssFeed struct {
		Channel struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			Items       []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
				Content     string `xml:"encoded"`
				Published   string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Convert to our shared types
	feed := &shared.RSSFeed{
		Title:       rssFeed.Channel.Title,
		Link:        rssFeed.Channel.Link,
		Description: rssFeed.Channel.Description,
	}

	for _, item := range rssFeed.Channel.Items {
		feed.Items = append(feed.Items, shared.RSSItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Content:     item.Content,
			Published:   item.Published,
		})
	}

	return feed, nil
}

// ProcessRSS is the Dagger function that will be called from the pipeline
func (r *RSS) ProcessRSS(ctx context.Context, url string) (*shared.RSSFeed, error) {
	return r.Fetch(ctx, url)
}
