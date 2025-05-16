package shared

// RSSFeed represents the structure of an RSS feed
type RSSFeed struct {
	Title       string
	Link        string
	Description string
	Items       []RSSItem
}

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string
	Link        string
	Description string
	Content     string
	Published   string
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
