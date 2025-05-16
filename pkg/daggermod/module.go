package daggermod

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/summarizer"
)

// Blogenetes represents the Dagger module
type Blogenetes struct {
	Client *dagger.Client
}

// New creates a new instance of the Dagger module
func New() *Blogenetes {
	// Create a new Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(fmt.Errorf("failed to connect to Dagger: %w", err))
	}

	return &Blogenetes{Client: client}
}

// Close cleans up the Dagger client
func (b *Blogenetes) Close() error {
	if b.Client != nil {
		return b.Client.Close()
	}
	return nil
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
	// Step 1: Fetch RSS feed
	fetchStep := b.Client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			`echo "Fetching RSS feed from %s" && `+
			`curl -s %s > /tmp/feed.xml && `+
			`echo "RSS feed fetched successfully"`,
			rssURL, rssURL,
		)})

	// Step 2: Extract content
	extractStep := b.Client.Container().
		From("alpine:latest").
		WithFile("/tmp/feed.xml", fetchStep.File("/tmp/feed.xml")).
		WithExec([]string{"sh", "-c", `
			echo "Extracting content from RSS feed..."
			# Extract first item's title and description
			CONTENT=$(grep -oP '(?<=<title>).*?(?=</title>)' /tmp/feed.xml | head -n 2 | tail -n 1 || echo "No title found")
			DESC=$(grep -oP '(?<=<description>).*?(?=</description>)' /tmp/feed.xml | head -n 1 || echo "No description found")
			echo -e "$CONTENT\n\n$DESC" > /tmp/content.txt
			echo "Content extracted successfully"
		`})

	// Step 3: Summarize content
	summarizer := summarizer.New(b.Client)
	content, err := extractStep.File("/tmp/content.txt").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	summary, err := summarizer.Summarize(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("summarization failed: %w", err)
	}

	// Step 4: Generate markdown
	markdownContent := fmt.Sprintf(
		"# %s\n\n## Summary\n\n%s\n\n## Key Points\n\n- %s",
		rssURL,
		summary.Summary,
		strings.Join(summary.KeyPoints, "\n- "),
	)

	// Step 5: Commit to GitHub
	return b.Client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf(`
			set -e
			echo "Cloning repository..."
			git clone https://github.com/%s/%s.git /repo
			cd /repo
			
			echo "Creating markdown file..."
			mkdir -p "$(dirname %s)"
			echo '%s' > %s
			
			echo "Configuring git..."
			git config --global user.name "Blogenetes Bot"
			git config --global user.email "bot@blogenetes.com"
			
			echo "Committing changes..."
			git checkout -b %s 2>/dev/null || git checkout %s
			git add %s
			git diff-index --quiet HEAD || git commit -m "%s"
			
			echo "Pushing changes..."
			echo "To push changes to GitHub, uncomment and configure the git push command below:"
			echo "# git push https://$GITHUB_TOKEN@github.com/%s/%s.git %s"
			
			echo "\nGenerated markdown content:"
			echo "------------------------"
			cat %s
		`,
			githubOwner, githubRepo,
			path, // For mkdir -p
			markdownContent, path,
			githubBranch, githubBranch, // For git checkout
			path,
			message,
			githubOwner, githubRepo, githubBranch,
			path,
		)}),
		nil
}
