package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger/markdown/internal/dagger"
)

// Summary represents the structured summary of content
type Summary struct {
	Summary   string   `json:"summary"`
	KeyPoints []string `json:"key_points"`
}

// MarkdownContent represents the final markdown document
type MarkdownContent struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

type Markdown struct {
	dag *dagger.Client
}

// New creates a new Markdown instance with the given Dagger client
func New(client *dagger.Client) *Markdown {
	return &Markdown{
		dag: client,
	}
}

// GenerateMarkdown creates a markdown document from the summary
// Returns an error if the summary is nil or if the title is empty
func (m *Markdown) GenerateMarkdown(ctx context.Context, title string, summary *Summary) (*MarkdownContent, error) {
	// Input validation
	if summary == nil {
		return nil, fmt.Errorf("summary cannot be nil")
	}
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	// Initialize string builder with estimated capacity
	estimatedSize := len(title) + len(summary.Summary) + (50 * len(summary.KeyPoints)) + 100
	var sb strings.Builder
	sb.Grow(estimatedSize)

	// Write header
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	
	// Add generation timestamp
	generatedAt := time.Now()
	sb.WriteString(fmt.Sprintf("*Generated on: %s*\n\n", generatedAt.Format(time.RFC1123)))

	// Write summary section
	sb.WriteString("## Summary\n\n")
	if summary.Summary != "" {
		sb.WriteString(summary.Summary)
	} else {
		sb.WriteString("No summary available.")
	}
	sb.WriteString("\n\n")

	// Write key points if available
	if len(summary.KeyPoints) > 0 {
		sb.WriteString("## Key Points\n\n")
		for _, point := range summary.KeyPoints {
			if point != "" {
				sb.WriteString(fmt.Sprintf("- %s\n", point))
			}
		}
	}

	return &MarkdownContent{
		Title:   title,
		Date:    generatedAt.Format("2006-01-02"),
		Content: sb.String(),
	}, nil
}

// ProcessMarkdown is the Dagger function that will be called from the pipeline
func (m *Markdown) ProcessMarkdown(ctx context.Context, title string, summary *Summary) (*MarkdownContent, error) {
	return m.GenerateMarkdown(ctx, title, summary)
}
