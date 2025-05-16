package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/shared"
)

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

// Summarizer handles text summarization using Dagger's LLM functionality
type Summarizer struct {
	Client *dagger.Client
}

// New creates a new Summarizer with a Dagger client
func New(client *dagger.Client) *Summarizer {
	return &Summarizer{Client: client}
}

// Summarize generates a summary and key points from the given text using Dagger's LLM
func (s *Summarizer) Summarize(ctx context.Context, text string) (*shared.Summary, error) {
	fmt.Println("🔍 Starting LLM summarization...")

	// Validate input text
	if text == "" {
		return nil, fmt.Errorf("empty text provided for summarization")
	}

	// Truncate text if too long (keep first 8000 chars to fit in context)
	if len(text) > 8000 {
		text = text[:8000] + "\n[Content truncated for brevity]"
	}

	// Debug: Print first 200 chars of input
	fmt.Printf("📝 Input text (first 200 chars): %s...\n", text[:min(200, len(text))])


	// Create the LLM prompt with the text directly embedded
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
	fmt.Println("🤖 Sending request to LLM...")
	llm := s.Client.LLM().
		WithPrompt(prompt)

	// Get the LLM response
	response, err := llm.LastReply(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
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

	// Parse the JSON response
	var llmResponse struct {
		Summary   string   `json:"summary"`
		KeyPoints []string `json:"key_points"`
	}

	if err := json.Unmarshal([]byte(cleanResponse), &llmResponse); err != nil {
		// If parsing fails, try to extract a basic summary from the response
		if strings.Contains(response, "I'm sorry") || strings.Contains(response, "unable") {
			// If the LLM indicates it can't complete the request, return a helpful error
			return nil, fmt.Errorf("LLM was unable to process the request. Please try again or check your prompt.")
		}
		
		// Fallback to using the first 200 characters as a summary
		return &shared.Summary{
			OriginalText: text,
			Summary:      fmt.Sprintf("Summary: %s", firstN(text, 200)),
			KeyPoints:    []string{"Could not extract key points from the response"},
		}, nil
	}

	// Validate the response
	if llmResponse.Summary == "" {
		llmResponse.Summary = firstN(text, 200)
	}

	if len(llmResponse.KeyPoints) == 0 {
		llmResponse.KeyPoints = []string{
			"No key points were extracted",
			"Please check the original content for details",
		}
	}

	// Ensure we have at least 3 key points
	for len(llmResponse.KeyPoints) < 3 {
		llmResponse.KeyPoints = append(llmResponse.KeyPoints, 
			fmt.Sprintf("Additional point %d", len(llmResponse.KeyPoints)+1))
	}

	return &shared.Summary{
		OriginalText: text,
		Summary:      llmResponse.Summary,
		KeyPoints:    llmResponse.KeyPoints,
	}, nil
}

// ProcessSummarize is the Dagger function that will be called from the pipeline
// It wraps the Summarize function to make it compatible with Dagger's pipeline steps
func (s *Summarizer) ProcessSummarize(ctx context.Context, text string) (*shared.Summary, error) {
	return s.Summarize(ctx, text)
}
