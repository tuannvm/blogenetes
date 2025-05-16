package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/github_publisher"
)

func main() {
	// Parse command line flags
	repo := flag.String("repo", "", "GitHub repository name")
	owner := flag.String("owner", "", "GitHub repository owner")
	branch := flag.String("branch", "main", "Git branch")
	path := flag.String("path", "", "Path to the file in the repository")
	message := flag.String("message", "Update via Dagger", "Commit message")
	content := flag.String("content", "", "Content to publish (if empty, will read from stdin)")
	flag.Parse()

	// Validate required flags
	if *repo == "" || *owner == "" || *path == "" {
		fmt.Println("Error: repo, owner, and path are required")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		fmt.Printf("Error connecting to Dagger: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Create a new file with the content
	tmpFile := client.Directory().
		WithNewFile("content", *content).
		File("content")

	// Create publisher and publish
	publisher := github_publisher.New().WithClient(client)
	_, err = publisher.Publish(
		ctx,
		tmpFile,
		*repo,
		*owner,
		*branch,
		*path,
		*message,
	)

	if err != nil {
		fmt.Printf("Error publishing to GitHub: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully published to GitHub")
	_, err = publisher.Publish(
		ctx,
		tmpFile,
		*repo,
		*owner,
		*branch,
		*path,
		*message,
	)

	if err != nil {
		fmt.Printf("Error publishing to GitHub: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully published to GitHub")
}
