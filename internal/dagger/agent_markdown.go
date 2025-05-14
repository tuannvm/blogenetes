package dagger

import (
	"context"
	"fmt"
	"time"
)

// AgentMarkdown converts summaries into markdown files.
type AgentMarkdown struct{}

// NewAgentMarkdown creates a new markdown agent.
func NewAgentMarkdown() *AgentMarkdown {
	return &AgentMarkdown{}
}

// ToMarkdown takes a title and summary and returns markdown bytes.
func (a *AgentMarkdown) ToMarkdown(ctx context.Context, title, summary string) ([]byte, error) {
	frontMatter := fmt.Sprintf("---\ntitle: \"%s\"\ndate: \"%s\"\n---\n\n", title, time.Now().Format(time.RFC3339))
	return []byte(frontMatter + summary + "\n"), nil
}
