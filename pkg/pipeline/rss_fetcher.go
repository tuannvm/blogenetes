package pipeline

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/rss"
)

type RSSFetcher struct {
	Client *dagger.Client
	rss    *rss.RSS
}

func NewRSSFetcher(client *dagger.Client) *RSSFetcher {
	return &RSSFetcher{
		Client: client,
		rss:    rss.New(),
	}
}

// Fetch fetches and parses an RSS feed
func (r *RSSFetcher) Fetch(ctx context.Context, url string) (*dagger.Container, error) {
	// Fetch the feed using our RSS module
	feed, err := r.rss.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch and parse RSS feed: %w", err)
	}

	if len(feed.Items) == 0 {
		return nil, fmt.Errorf("no items found in the feed")
	}

	// Get the latest item
	item := feed.Items[0]

	// Create a container with the content
	container := r.Client.Container().
		From("alpine:latest").
		WithNewFile("/tmp/content.txt", fmt.Sprintf("# %s\n\n%s", item.Title, item.Content))

	// Add debug information
	container = container.WithExec([]string{"sh", "-c", fmt.Sprintf(`
		echo "✅ Successfully fetched and parsed RSS feed"
		echo "📰 Title: %s"
		echo "🔗 Link: %s"
		echo "📅 Published: %s"
		echo "📝 Content length: %d characters"
		wc -c /tmp/content.txt
	`,
		strings.ReplaceAll(item.Title, "'", "'\\''"),
		item.Link,
		item.Published.Format("2006-01-02 15:04:05"),
		len(item.Content),
	)})

	return container, nil
}
