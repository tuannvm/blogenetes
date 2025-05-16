package main

import (
	"context"
	"github.com/tuannvm/blogenetes/modules/shared"
	"strings"
)

type Summarizer struct{}

// Summarize generates a summary of the given text
func (s *Summarizer) Summarize(ctx context.Context, text string) (*shared.Summary, error) {
	// This is a simple summarization implementation
	// In a real-world scenario, you would use an LLM or other NLP technique
	
	sentences := strings.Split(text, ". ")
	if len(sentences) == 0 {
		return &shared.Summary{
			OriginalText: text,
			Summary:     text,
			KeyPoints:  []string{text},
		}, nil
	}

	// Simple summarization: take first 3 sentences as summary
	summary := ""
	if len(sentences) > 3 {
		summary = strings.Join(sentences[:3], ". ") + "."
	} else {
		summary = text
	}

	// Simple key points: first 3 sentences as bullet points
	keyPoints := make([]string, 0, 3)
	for i := 0; i < 3 && i < len(sentences); i++ {
		keyPoints = append(keyPoints, strings.TrimSpace(sentences[i])+".")
	}

	return &shared.Summary{
		OriginalText: text,
		Summary:     summary,
		KeyPoints:  keyPoints,
	}, nil
}

// ProcessSummarize is the Dagger function that will be called from the pipeline
func (s *Summarizer) ProcessSummarize(ctx context.Context, text string) (*shared.Summary, error) {
	return s.Summarize(ctx, text)
}
