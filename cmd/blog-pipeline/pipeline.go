package main

import (
	"context"
	"fmt"
	"log"

	"dagger.io/dagger"
	"github.com/mmcdole/gofeed"
	"github.com/tuannvm/blogenetes/github-publisher"
	"github.com/tuannvm/blogenetes/markdown"
	"github.com/tuannvm/blogenetes/rss"
	"github.com/tuannvm/blogenetes/summarizer"
)

// Config holds the pipeline configuration
type Config struct {
	RSSURL    string
	Owner     string
	Repo      string
	Branch    string
	FilePath  string
	CommitMsg string
}

// Pipeline represents the blog pipeline
type Pipeline struct {
	client *dagger.Client
}

// NewPipeline creates a new blog pipeline
func NewPipeline(client *dagger.Client) *Pipeline {
	return &Pipeline{
		client: client,
	}
}

// Run executes the blog pipeline with the given configuration
func (p *Pipeline) Run(ctx context.Context, cfg Config) (string, error) {
	// Validate required fields
	if cfg.Owner == "" || cfg.Repo == "" {
		return "", fmt.Errorf("GitHub owner and repo are required")
	}

	log.Println("🚀 Starting blog pipeline...")

	// Step 1: Fetch RSS feed
	log.Println("📡 Fetching RSS feed from:", cfg.RSSURL)
	rssClient := rss.New(ctx, p.client)
	feed, err := rssClient.Fetch(ctx, cfg.RSSURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	if len(feed.Items) == 0 {
		return "", fmt.Errorf("no items found in the feed")
	}

	// Use the first item for demonstration
	item := feed.Items[0]
	log.Printf("📰 Found article: %s", item.Title)

	// Step 2: Summarize content
	log.Println("🧠 Summarizing content...")
	summarizerClient := summarizer.New(ctx, p.client)
	summary, err := summarizerClient.Summarize(ctx, item.Description)
	if err != nil {
		return "", fmt.Errorf("failed to summarize content: %w", err)
	}

	// Step 3: Format as markdown
	log.Println("📄 Formatting as markdown...")
	markdownClient := markdown.New(ctx, p.client)
	mdContent, err := markdownClient.Format(ctx, markdown.Content{
		Title:   item.Title,
		Content: summary.Summary,
		Link:    item.Link,
	})
	if err != nil {
		return "", fmt.Errorf("failed to format markdown: %w", err)
	}

	// Step 4: Publish to GitHub
	log.Println("🚀 Publishing to GitHub...")
	publisher := github_publisher.New(ctx, p.client)
	_, err = publisher.Publish(ctx, github_publisher.PublishOpts{
		Content:   mdContent,
		Owner:     cfg.Owner,
		Repo:      cfg.Repo,
		Branch:    cfg.Branch,
		Path:      cfg.FilePath,
		CommitMsg: cfg.CommitMsg,
	})
	if err != nil {
		return "", fmt.Errorf("failed to publish to GitHub: %w", err)
	}

	return fmt.Sprintf("Successfully published to %s/%s/%s", cfg.Owner, cfg.Repo, cfg.FilePath), nil
}
