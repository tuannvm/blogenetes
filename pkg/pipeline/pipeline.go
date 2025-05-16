package pipeline

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/summarizer"
)

type Pipeline struct {
	Client      *dagger.Client
	RSSFetcher  *RSSFetcher
	Formatter   *MarkdownFormatter
	Summarizer  *summarizer.Summarizer
	Publisher   *GitHubPublisher
}

func NewPipeline(client *dagger.Client) *Pipeline {
	return &Pipeline{
		Client:      client,
		RSSFetcher:  NewRSSFetcher(client),
		Formatter:   NewMarkdownFormatter(client),
		Summarizer:  summarizer.New(client),
		Publisher:   NewGitHubPublisher(client),
	}
}

// Run executes the full pipeline with proper Dagger pipeline steps
func (p *Pipeline) Run(
	ctx context.Context,
	rssURL string,
	githubOwner string,
	githubRepo string,
	githubBranch string,
	path string,
	message string,
) (*dagger.Container, error) {
	// Step 1: Fetch RSS content
	fmt.Println("⏳ Fetching RSS feed...")
	fetchContainer, err := p.RSSFetcher.Fetch(ctx, rssURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS: %w", err)
	}

	// Step 2: Extract content
	fmt.Println("📝 Extracting content...")
	content, err := fetchContainer.File("/tmp/content.txt").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// Step 3: Summarize content
	fmt.Println("🧠 Summarizing content...")
	summary, err := p.Summarizer.Summarize(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize: %w", err)
	}

	// Step 4: Format as markdown
	fmt.Println("📄 Formatting as markdown...")
	markdown, err := p.Formatter.Format(ctx, summary.Summary, "Generated Post")
	if err != nil {
		return nil, fmt.Errorf("failed to format markdown: %w", err)
	}

	// Step 5: Publish to GitHub
	fmt.Println("🚀 Publishing to GitHub...")
	return p.Publisher.Publish(
		ctx,
		markdown,
		githubRepo,
		githubOwner,
		githubBranch,
		path,
		message,
	)
}

// RunSummarizerOnly runs just the summarizer component
func (p *Pipeline) RunSummarizerOnly(ctx context.Context, content string) (string, error) {
	summary, err := p.Summarizer.Summarize(ctx, content)
	if err != nil {
		return "", fmt.Errorf("summarization failed: %w", err)
	}
	return summary.Summary, nil
}

// RunFormatterOnly runs just the formatter component
func (p *Pipeline) RunFormatterOnly(ctx context.Context, content string, title string) (string, error) {
	return p.Formatter.Format(ctx, content, title)
}

// RunPublisherOnly runs just the publisher component
func (p *Pipeline) RunPublisherOnly(
	ctx context.Context,
	content string,
	repo string,
	owner string,
	branch string,
	path string,
	message string,
) (*dagger.Container, error) {
	return p.Publisher.Publish(ctx, content, repo, owner, branch, path, message)
}
