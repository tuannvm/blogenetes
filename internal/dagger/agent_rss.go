package dagger

import (
	"context"

	"github.com/mmcdole/gofeed"
)

// AgentRSS fetches and parses articles from an RSS feed.
type AgentRSS struct{}

// NewAgentRSS creates a new RSS agent.
func NewAgentRSS() *AgentRSS {
	return &AgentRSS{}
}

// FetchArticles fetches and parses articles from the given RSS URLs.
func (a *AgentRSS) FetchArticles(ctx context.Context, rssURLs []string) ([]*gofeed.Item, error) {
	parser := gofeed.NewParser()
	var allItems []*gofeed.Item
	for _, url := range rssURLs {
		feed, err := parser.ParseURL(url)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, feed.Items...)
	}
	return allItems, nil
}
