# Blogenetes

A modular blog generation and publishing pipeline using Dagger.

## Overview

Blogenetes is a Go application that automates the process of fetching content from RSS feeds, processing it, and publishing it to a GitHub repository. The application is built using a modular architecture with Dagger for containerized pipeline execution.

## Architecture

The application is organized into the following modules:

1. **RSS Module**: Fetches and parses RSS feeds
2. **Summarizer Module**: Processes content to generate summaries
3. **Markdown Module**: Converts content to markdown format
4. **GitHub Module**: Handles committing and pushing content to GitHub

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker or another container runtime
- Dagger CLI installed
- GitHub personal access token with repo access

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/blogenetes.git
   cd blogenetes
   ```

2. Build the application:
   ```bash
   go build -o bin/blogenetes
   ```

### Configuration

Set up the required environment variables:

```bash
export GITHUB_TOKEN=your_github_token
```

### Usage

Run the application with the required flags:

```bash
./bin/blogenetes \
  --rss "https://example.com/feed.xml" \
  --title "My Blog Post" \
  --owner your-github-username \
  --repo your-repo-name \
  --branch main \
  --path "content/posts/new-post.md" \
  --message "Add new blog post"
```

### Dagger Integration

Each module is designed to work with Dagger's pipeline system. The main pipeline is defined in `main.go` and executes the following steps:

1. Fetches RSS feed content
2. Processes the latest feed item
3. Generates a summary
4. Converts to markdown
5. Commits to GitHub

## Development

### Module Structure

Each module follows this structure:

```
modules/
  <module-name>/
    internal/     # Internal package code
    main.go       # Module entry point
    dagger.json   # Dagger module configuration
```

### Adding a New Module

1. Create a new directory under `modules/`
2. Add a `dagger.json` file with the module configuration
3. Implement the module logic in `main.go`
4. Update the main pipeline to use the new module

## License

MIT
