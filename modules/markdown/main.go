package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tuannvm/blogenetes/modules/shared"
)

type Markdown struct{}

// GenerateMarkdown creates a markdown document from the summary
func (m *Markdown) GenerateMarkdown(ctx context.Context, title string, summary *shared.Summary) (*shared.MarkdownContent, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("*Generated on: %s*\n\n", time.Now().Format(time.RFC1123)))

	sb.WriteString("## Summary\n\n")
	sb.WriteString(summary.Summary)
	sb.WriteString("\n\n")

	sb.WriteString("## Key Points\n\n")
	for _, point := range summary.KeyPoints {
		sb.WriteString(fmt.Sprintf("- %s\n", point))
	}

	return &shared.MarkdownContent{
		Title:   title,
		Date:    time.Now().Format("2006-01-02"),
		Content: sb.String(),
	}, nil
}

// ProcessMarkdown is the Dagger function that will be called from the pipeline
func (m *Markdown) ProcessMarkdown(ctx context.Context, title string, summary *shared.Summary) (*shared.MarkdownContent, error) {
	return m.GenerateMarkdown(ctx, title, summary)
}
