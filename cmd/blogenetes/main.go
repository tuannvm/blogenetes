package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/pipeline"
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
		return fmt.Errorf("flags --rss, --owner, and --repo are required")
	}

	// Check for GitHub token
	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return fmt.Errorf("failed to connect to Dagger: %w", err)
	}
	defer client.Close()

	// Create and run the pipeline
	p, err := pipeline.NewPipeline(client)
	if err != nil {
		log.Fatalf("Failed to create pipeline: %v", err)
	}

	// Run the full pipeline
	log.Println("🚀 Starting Blogenetes pipeline...")
	container, err := p.Run(
		ctx,
		*rssURL,
		*owner,
		*repo,
		*branch,
		*path,
		*message,
	)

	if err != nil {
		return fmt.Errorf("❌ Pipeline failed: %w", err)
	}

	// Get the output
	output, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("❌ Failed to get container output: %w", err)
	}

	fmt.Println("\n✅ Pipeline execution completed successfully!")
	fmt.Println("📝 Output:")
	fmt.Println(output)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	log.Printf("✅ Successfully published to %s/%s/%s\n", *owner, *repo, *path)
}
