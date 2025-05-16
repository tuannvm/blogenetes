package shared

import "time"

// RSSFeed represents the structure of an RSS feed
type RSSFeed struct {
	Title       string
	Link        string
	Description string
	Items       []RSSItem
}

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Published   time.Time `json:"published"`
}

// Summary represents the summarized content
type Summary struct {
	OriginalText string
	Summary     string
	KeyPoints  []string
}

// MarkdownContent represents the final markdown content
type MarkdownContent struct {
	Title   string
	Date    string
	Content string
}
