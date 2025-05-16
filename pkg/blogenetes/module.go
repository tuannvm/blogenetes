package blogenetes

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type Blogenetes struct{}

func New() *Blogenetes {
	return &Blogenetes{}
}

// ProcessRSSFeed processes an RSS feed through the pipeline
func (b *Blogenetes) ProcessRSSFeed(
	ctx context.Context,
	rssURL string,
	githubOwner string,
	githubRepo string,
	githubBranch string,
	path string,
	message string,
) (*dagger.Container, error) {
	// Create a new Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Dagger: %w", err)
	}
	defer client.Close()

	// Run the pipeline in a container
	return client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			`echo "Processing RSS feed from %s" && `+
			`echo "Will commit to: %s/%s@%s" && `+
			`echo "Path: %s" && `+
			`echo "Commit message: %s"`,
			rssURL, githubOwner, githubRepo, githubBranch, path, message,
		)}),
		nil
}
