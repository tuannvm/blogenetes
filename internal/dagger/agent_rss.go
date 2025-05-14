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

// FetchArticles fetches and parses articles from the given RSS URL.
func (a *AgentRSS) FetchArticles(ctx context.Context, rssURL string) ([]*gofeed.Item, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(rssURL)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}
