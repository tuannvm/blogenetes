package main

import (
	"context"
	"flag"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	// Parse command line flags
	rssURL := flag.String("rss", "https://example.com/feed", "RSS feed URL")
	repoOwner := flag.String("owner", "", "GitHub repository owner")
	repoName := flag.String("repo", "", "GitHub repository name")
	branch := flag.String("branch", "main", "GitHub branch")
	filePath := flag.String("path", "content/blog/latest.md", "File path in the repository")
	commitMsg := flag.String("message", "Update blog content", "Commit message")
	flag.Parse()

	// Initialize Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Create and run the pipeline
	pipeline := NewPipeline(client)
	output, err := pipeline.Run(ctx, Config{
		RSSURL:    *rssURL,
		Owner:     *repoOwner,
		Repo:      *repoName,
		Branch:    *branch,
		FilePath:  *filePath,
		CommitMsg: *commitMsg,
	})

	if err != nil {
		log.Fatalf("❌ Pipeline failed: %v", err)
	}

	log.Printf("✅ %s", output)
}
