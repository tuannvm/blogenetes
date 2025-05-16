// Package main provides a Dagger module for text summarization
//
// This module provides functionality to summarize text content
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger/summarizer/internal/dagger"
)

// Summary represents the structured output of the summarization process
type Summary struct {
	// A concise summary of the input text
	Summary string `json:"summary"`
	// Key points extracted from the text
	KeyPoints []string `json:"key_points"`
}

// Summarizer handles text summarization
type Summarizer struct {
	dag *dagger.Client
}

// NewSummarizer creates a new instance of the Summarizer
//
// Parameters:
// - dag: Dagger client instance
//
// Returns:
// - *Summarizer: A new Summarizer instance
func NewSummarizer(dag *dagger.Client) *Summarizer {
	return &Summarizer{dag: dag}
}

// firstN returns the first n characters of a string, or the entire string if it's shorter than n
func firstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Summarize generates a summary and key points from the given text using Dagger's LLM
func (s *Summarizer) Summarize(ctx context.Context, text string) (*Summary, error) {
	// Input validation
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Truncate text if too long (keep first 8000 chars to fit in context)
	const maxInputLength = 8000
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "... [truncated]"
	}

	// Create a prompt for the LLM
	prompt := fmt.Sprintf(`You are an expert at summarizing content and extracting key points.

Please analyze the following text and provide:
1. A concise 2-3 sentence summary
2. 3-5 key points as bullet points

Text to summarize:
%s

Format your response as a valid JSON object with this exact structure (no extra text or markdown):
{
  "summary": "Your summary here",
  "key_points": ["Point 1", "Point 2", "Point 3"]
}`, text)

	// Execute the LLM with the prompt
	llm := s.dag.LLM().WithPrompt(prompt)

	// Get the LLM response
	response, err := llm.LastReply(ctx)
	if err != nil {
		// Fallback to basic summarization if LLM fails
		return s.basicSummarize(text), nil
	}

	// Parse the JSON response
	var llmResponse struct {
		Summary   string   `json:"summary"`
		KeyPoints []string `json:"key_points"`
	}

	// Clean up the response - sometimes the LLM includes markdown code blocks
	cleanResponse := response
	if len(cleanResponse) > 0 && cleanResponse[0] == '`' {
		// Try to extract JSON from markdown code block
		start := strings.Index(cleanResponse, "{")
		end := strings.LastIndex(cleanResponse, "}")
		if start > 0 && end > start {
			cleanResponse = cleanResponse[start : end+1]
		}
	}

	// Try to parse the JSON response
	if err := json.Unmarshal([]byte(cleanResponse), &llmResponse); err != nil {
		// If parsing fails, fall back to basic summarization
		return s.basicSummarize(text), nil
	}

	// Ensure we have valid data
	if llmResponse.Summary == "" {
		llmResponse.Summary = firstN(text, 200)
	}

	if len(llmResponse.KeyPoints) == 0 {
		llmResponse.KeyPoints = []string{
			"No key points were extracted",
			"Please check the original content for details",
		}
	}

	return &Summary{
		Summary:   llmResponse.Summary,
		KeyPoints: llmResponse.KeyPoints,
	}, nil
}

// basicSummarize provides a fallback implementation when LLM is not available
func (s *Summarizer) basicSummarize(text string) *Summary {
	// Simple implementation that splits the text into sentences
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	// Create a basic summary (first few sentences)
	summary := ""
	if len(sentences) > 0 {
		summary = strings.TrimSpace(sentences[0])
	}

	// Extract key points (first few sentences after the summary)
	keyPoints := make([]string, 0, 3)
	for i := 1; i < min(4, len(sentences)); i++ {
		if trimmed := strings.TrimSpace(sentences[i]); trimmed != "" {
			keyPoints = append(keyPoints, trimmed+".")
		}
	}

	// Ensure we have at least one key point
	if len(keyPoints) == 0 && summary != "" {
		keyPoints = append(keyPoints, firstN(summary, 100))
	}

	return &Summary{
		Summary:   summary,
		KeyPoints: keyPoints,
	}
}

// ProcessSummarize is the Dagger function that will be called from the pipeline
// It wraps the Summarize function to make it compatible with Dagger's pipeline steps
func (s *Summarizer) ProcessSummarize(ctx context.Context, text string) (*Summary, error) {
	return s.Summarize(ctx, text)
}
