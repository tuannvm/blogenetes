package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/openai/openai-go"
)

func main() {
	rssURL := os.Getenv("RSS_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if rssURL == "" || apiKey == "" {
		fmt.Fprintln(os.Stderr, "RSS_URL or OPENAI_API_KEY not set")
		os.Exit(1)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse RSS feed: %v\n", err)
		os.Exit(1)
	}
	if len(feed.Items) == 0 {
		fmt.Fprintln(os.Stderr, "no items found in RSS feed")
		os.Exit(1)
	}
	entry := feed.Items[0]

	slug := strings.ToLower(entry.Title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")

	resp, err := http.Get(entry.Link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch article: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse HTML: %v\n", err)
		os.Exit(1)
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
	resp, err := client.CreateChatCompletion(context.Background(), &openai.ChatCompletionRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
		MaxTokens:   200,
		Temperature: 0.7,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "OpenAI API error: %v\n", err)
		os.Exit(1)
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)

	date := time.Now().UTC().Format("2006-01-02")
	outDir := "/output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		os.Exit(1)
	}
	filename := fmt.Sprintf("%s/%s-%s.md", outDir, date, slug)
	md := fmt.Sprintf("---\ntitle: \"Summary: %s\"\ndate: %s\n---\n\n%s\n", entry.Title, date, summary)

	if err := ioutil.WriteFile(filename, []byte(md), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Wrote summary to", filename)
}
