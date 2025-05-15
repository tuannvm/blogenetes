package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/tuannvm/blogenetes/internal/dagger"
)

func main() {
	ctx := context.Background()

	openAIClient := dagger.SetupOpenAIClient()
	daggerClient, err := dagger.SetupDaggerClient()
	if err != nil {
		log.Fatalf("failed to connect to Dagger: %v", err)
	}
	defer daggerClient.Close()

	rssURL := flag.String("rss", "https://www.npr.org/rss/rss.php?id=1001", "Comma-separated RSS feed URLs (default top 10 sources)")
	owner := flag.String("owner", "", "GitHub owner")
	repo := flag.String("repo", "", "GitHub repository name")
	branch := flag.String("branch", "main", "GitHub branch to commit to")
	path := flag.String("path", "post.md", "Path in the repo for the markdown file")
	message := flag.String("message", "Add new post via agent", "Commit message")
	flag.Parse()

	if *rssURL == "" || *owner == "" || *repo == "" {
		log.Fatal("flags --rss, --owner, and --repo are required")
	}

	// Initialize agents
	rssAgent := dagger.NewAgentRSS()
	summarizer := dagger.NewAgentSummarizer(openAIClient)
	markdownAgent := dagger.NewAgentMarkdown()
	githubAgent := dagger.NewAgentGitHub()

	// Fetch articles
	rssURLs := strings.Split(*rssURL, ",")
	items, err := rssAgent.FetchArticles(ctx, rssURLs)
	if err != nil {
		log.Fatalf("failed to fetch articles: %v", err)
	}
	if len(items) == 0 {
		log.Fatal("no articles found at RSS feed")
	}
	latest := items[0]

	// Summarize (Dagger container step)
	summary, err := summarizer.Summarize(ctx, daggerClient, latest.Content)
	if err != nil {
		log.Fatalf("failed to summarize article: %v", err)
	}

	// Convert to Markdown
	mdBytes, err := markdownAgent.ToMarkdown(ctx, latest.Title, summary)
	if err != nil {
		log.Fatalf("failed to generate markdown: %v", err)
	}

	// Commit to GitHub
	err = githubAgent.CommitToRepo(ctx, *owner, *repo, *branch, *path, string(mdBytes), *message)
	if err != nil {
		log.Fatalf("failed to commit to repo: %v", err)
	}

	fmt.Println("Successfully published markdown to GitHub")
}
