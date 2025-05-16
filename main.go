// Package main is the root module for the Blogenetes Dagger pipeline.
// It composes all the individual modules (RSS, summarizer, markdown, github-publisher)
// into a complete workflow.
package main

import (
	"context"

	"dagger.io/dagger"
)

// Blogenetes is the root module that composes all the individual modules
type Blogenetes struct {
	dag *dagger.Client
}

// New creates a new instance of the Blogenetes module
func NewBlogenetes(dag *dagger.Client) *Blogenetes {
	return &Blogenetes{
		dag: dag,
	}
}

// RunPipeline runs the complete blog generation pipeline
func (b *Blogenetes) RunPipeline(
	ctx context.Context,
	rssURL string,
	githubRepo string,
	githubToken *dagger.Secret,
	openaiKey *dagger.Secret,
) (string, error) {
	// Initialize the pipeline
	pipeline := b.dag.Pipeline("blogenetes-pipeline", dagger.PipelineOpts{
		Description: "Generate blog posts from RSS feeds with AI summarization",
	})

	// Step 1: Fetch RSS feed
	rssJSON, err := pipeline.Pipeline("fetch-rss").RSS().Fetch(ctx, rssURL)
	if err != nil {
		return "", err
	}

	// Step 2: Summarize content
	summarized, err := pipeline.Pipeline("summarize-content").Summarizer().Summarize(ctx, rssJSON, openaiKey)
	if err != nil {
		return "", err
	}

	// Step 3: Generate markdown
	markdownFile, err := pipeline.Pipeline("generate-markdown").Markdown().Generate(ctx, summarized)
	if err != nil {
		return "", err
	}

	// Step 4: Publish to GitHub
	result, err := pipeline.Pipeline("publish-github").GithubPublisher().Publish(
		ctx,
		markdownFile,
		githubRepo,
		"main",
		"feed.md",
		githubToken,
	)
	if err != nil {
		return "", err
	}

	return result, nil
}
