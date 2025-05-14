package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/sashabaranov/go-openai"
)

func main() {
	fmt.Println("Starting blog fetcher and summarizer...")
	
	rssURL := os.Getenv("RSS_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	
	fmt.Printf("RSS_URL is set: %v\n", rssURL != "")
	fmt.Printf("OPENAI_API_KEY is set: %v\n", apiKey != "")
	
	if rssURL == "" || apiKey == "" {
		log.Fatal("environment variables RSS_URL and OPENAI_API_KEY must be set")
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		log.Fatalf("failed to parse RSS feed: %v", err)
	}
	if len(feed.Items) == 0 {
		log.Fatal("no items found in RSS feed")
	}
	entry := feed.Items[0]

	slug := strings.ToLower(entry.Title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")

	httpResp, err := http.Get(entry.Link)
	if err != nil {
		log.Fatalf("failed to fetch article: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("failed to fetch article: status %s", httpResp.Status)
	}
	defer httpResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResp.Body)
	if err != nil {
		log.Fatalf("failed to parse HTML: %v", err)
	}

	var builder strings.Builder
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			builder.WriteString(text)
			builder.WriteString("\n\n")
		}
	})

	client := openai.NewClient(apiKey)
	prompt := fmt.Sprintf("Summarise the following blog post in 3–5 bullet points:\n\n%s", builder.String())
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   200,
			Temperature: 0.7,
		},
	)
	if err != nil {
		log.Fatalf("OpenAI API error: %v", err)
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)

	date := time.Now().UTC().Format("2006-01-02")
	outDir := "/output"
	fmt.Printf("Creating output directory: %s\n", outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}
	filename := fmt.Sprintf("%s/new_post.md", outDir)
	fmt.Printf("Will write to file: %s\n", filename)
	md := fmt.Sprintf("---\ntitle: \"Summary: %s\"\ndate: %s\n---\n\n%s\n", entry.Title, date, summary)

	if err := os.WriteFile(filename, []byte(md), 0644); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	// Verify file was written
	if _, err := os.Stat(filename); err != nil {
		log.Fatalf("failed to verify file exists after writing: %v", err)
	}

	fmt.Println("Successfully wrote summary to", filename)
	fmt.Println("Process completed successfully!")
}
