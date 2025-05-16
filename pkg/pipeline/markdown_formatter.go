package pipeline

import (
	"context"
	"fmt"
	"time"

	"dagger.io/dagger"
)

type MarkdownFormatter struct {
	Client *dagger.Client
}

func NewMarkdownFormatter(client *dagger.Client) *MarkdownFormatter {
	return &MarkdownFormatter{Client: client}
}

// Format formats the content as markdown
func (m *MarkdownFormatter) Format(ctx context.Context, content, title string) (string, error) {
	// Get current date in YYYY-MM-DD format
	now := time.Now().Format("2006-01-02")
	
	// Format the markdown with proper frontmatter
	formatted := fmt.Sprintf(`---
title: "%s"
date: %s
---

%s`, title, now, content)
	
	return formatted, nil
}
