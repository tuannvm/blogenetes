package main

import (
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Parse command line flags
	rssURL := flag.String("rss", "https://www.npr.org/rss/rss.php?id=1001", "RSS feed URL")
	owner := flag.String("owner", "your-org", "GitHub owner")
	repo := flag.String("repo", "your-repo", "GitHub repository name")
	branch := flag.String("branch", "main", "GitHub branch to commit to")
	path := flag.String("path", "posts/post.md", "Path to save the post")
	message := flag.String("message", "New post via agent", "Commit message")
	flag.Parse()

	if *rssURL == "" || *owner == "" || *repo == "" {
		log.Fatal("flags --rss, --owner, and --repo are required")
	}

	// Create a new instance of our Blogenetes module
	blog := New()

	// Execute the pipeline
	ctx := context.Background()
	container, err := blog.ProcessRSSFeed(
		ctx,
		*rssURL,
		*owner,
		*repo,
		*branch,
		*path,
		*message,
	)

	if err != nil {
		log.Fatalf("pipeline failed: %v", err)
	}

	// Get the output
	output, err := container.Stdout(ctx)
	if err != nil {
		log.Fatalf("failed to get container output: %v", err)
	}

	fmt.Println("Pipeline execution completed successfully!")
	fmt.Println("Output:")
	fmt.Println(output)
}
