package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

// GithubPublisher is a module for publishing content to GitHub
type GithubPublisher struct {
	// Client is the Dagger client
	Client *dagger.Client
}

// New creates a new instance of the GithubPublisher
func New() *GithubPublisher {
	return &GithubPublisher{}
}

// WithClient sets the Dagger client for the publisher
func (m *GithubPublisher) WithClient(client *dagger.Client) *GithubPublisher {
	m.Client = client
	return m
}

// Publish commits and pushes content to a GitHub repository
//
// Parameters:
// - ctx: The Dagger context
// - content: The content to publish (as a Dagger file)
// - repo: The repository name
// - owner: The repository owner
// - branch: The branch to push to
// - path: The path where to create the file
// - message: The commit message
//
// Returns:
// - A container with the result of the operation
// - An error if something goes wrong
func (m *GithubPublisher) Publish(
	ctx context.Context,
	content *dagger.File,
	repo string,
	owner string,
	branch string,
	path string,
	message string,
) (*dagger.Container, error) {
	// Get the Dagger client
	client := m.Client
	if client == nil {
		var err error
		client, err = dagger.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Dagger: %w", err)
		}
		defer client.Close()
	}

	// Set default values
	authorName := "Blogenetes"
	authorEmail := "blogenetes@example.com"
	committerName := authorName
	committerEmail := authorEmail

	// Create a container with git
	container := client.Container().
		From("alpine/git:latest").
		WithEnvVariable("GIT_AUTHOR_NAME", authorName).
		WithEnvVariable("GIT_AUTHOR_EMAIL", authorEmail).
		WithEnvVariable("GIT_COMMITTER_NAME", committerName).
		WithEnvVariable("GIT_COMMITTER_EMAIL", committerEmail)

	// Get the GitHub token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	// Clone the repository
	repoURL := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", owner, githubToken, owner, repo)
	container = container.WithExec([]string{"git", "clone", "--depth", "1", "--branch", branch, repoURL, "/repo"})

	// Create the destination directory if it doesn't exist
	destPath := fmt.Sprintf("/repo/%s", path)
	container = container.WithExec([]string{"mkdir", "-p", destPath})

	// Create parent directories if needed
	if lastSlash := strings.LastIndex(path, "/"); lastSlash != -1 {
		dirPath := fmt.Sprintf("/repo/%s", path[:lastSlash])
		container = container.WithExec([]string{"mkdir", "-p", dirPath})
	}

	// Copy the content to the destination
	container = container.WithMountedFile(destPath, content)

	// Configure git
	container = container.WithWorkdir("/repo")
	container = container.WithExec([]string{"git", "config", "user.name", authorName})
	container = container.WithExec([]string{"git", "config", "user.email", authorEmail})

	// Add and commit the changes
	container = container.WithExec([]string{"git", "add", "."})
	container = container.WithExec([]string{"git", "commit", "-m", message})

	// Push the changes
	container = container.WithExec([]string{"git", "push", "origin", fmt.Sprintf("HEAD:%s", branch)})

	return client.Container().
		From("alpine:latest").
		WithExec([]string{"echo", "✅ Successfully published to GitHub"}), nil
}
