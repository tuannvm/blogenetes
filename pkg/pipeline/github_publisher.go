package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

type GitHubPublisher struct {
	Client *dagger.Client
}

func NewGitHubPublisher(client *dagger.Client) *GitHubPublisher {
	return &GitHubPublisher{Client: client}
}

// Publish commits and pushes the content to a GitHub repository
func (g *GitHubPublisher) Publish(
	ctx context.Context,
	content string,
	repo string,
	owner string,
	branch string,
	path string,
	message string,
) (*dagger.Container, error) {
	// Get GitHub token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create a container with git installed
	container := g.Client.Container().
		From("alpine/git:latest").
		WithEnvVariable("GIT_AUTHOR_NAME", "Blogenetes").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "blogenetes@example.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "Blogenetes").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "blogenetes@example.com")

	// Clone the repository
	repoURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", githubToken, owner, repo)
	container = container.
		WithExec([]string{"git", "clone", "--depth", "1", "--branch", branch, repoURL, "/repo"}).
		WithWorkdir("/repo")

	// Create the directory structure if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "." {
		container = container.WithExec([]string{"mkdir", "-p", dir})
	}

	// Create the file using echo since WithNewFile has issues with the current Dagger SDK
	container = container.
		WithExec([]string{"sh", "-c", fmt.Sprintf("echo '%s' > %s", content, path)}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/repo"}).
		WithExec([]string{"git", "config", "--global", "user.name", "Blogenetes"}).
		WithExec([]string{"git", "config", "--global", "user.email", "blogenetes@example.com"}).
		WithExec([]string{"git", "add", path}).
		WithExec([]string{"git", "commit", "-m", message}).
		WithExec([]string{"git", "remote", "set-url", "origin", fmt.Sprintf("https://%s@github.com/%s/%s.git", githubToken, owner, repo)}).
		WithExec([]string{"git", "push", "origin", fmt.Sprintf("HEAD:%s", branch)})

	return container, nil
}
