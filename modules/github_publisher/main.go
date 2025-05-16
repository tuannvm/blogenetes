// Package main provides a Dagger module for publishing content to GitHub
//
// This module allows you to commit and push content to a GitHub repository
// using Dagger's pipeline capabilities.
package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GithubPublisher is the main module type that provides GitHub publishing functionality
// GithubPublisher represents a client for publishing content to GitHub
type GithubPublisher struct {
	dag *dagger.Client
}

// Publish commits and pushes content to a GitHub repository
//
// Parameters:
// - ctx: The context for the operation
// - content: The file or directory to publish
// - repo: The name of the GitHub repository
// - owner: The owner of the GitHub repository
// - branch: The branch to push to (defaults to "main")
// - path: The path in the repository to publish to
// - message: The commit message
// - githubToken: The GitHub token with repository write access
//
// Returns:
// - *dagger.Container: The container used for Git operations
// - error: Any error that occurred
func (m *GithubPublisher) Publish(
	ctx context.Context,
	content *dagger.Directory,
	repo string,
	owner string,
	branch string,
	path string,
	message string,
	githubToken *dagger.Secret,
) (*dagger.Container, error) {
	if branch == "" {
		branch = "main"
	}

	// Get the GitHub token
	tokenValue, err := githubToken.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read GitHub token: %w", err)
	}

	// Initialize a new container with git
	ctr := m.dag.Container().From("alpine:latest")

	// Install git and required dependencies
	ctr = ctr.WithExec([]string{"apk", "add", "--no-cache", "git", "openssh"})

	// Configure git
	ctr = ctr.WithExec([]string{"git", "config", "--global", "user.name", "Dagger Bot"})
	ctr = ctr.WithExec([]string{"git", "config", "--global", "user.email", "bot@dagger.io"})

	// Clone the repository
	repoURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", tokenValue, owner, repo)
	ctr = ctr.WithWorkdir("/src")
	ctr = ctr.WithExec([]string{"git", "clone", "-b", branch, "--single-branch", repoURL, "."})

	// Add the content
	ctr = ctr.WithDirectory(path, content, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{".git"},
	})

	// Configure git to use the token for authentication
	ctr = ctr.WithExec([]string{"git", "config", "--global", "url.https://github.com/.insteadOf", "git@github.com:"})
	ctr = ctr.WithExec([]string{"git", "config", "--global", "url.https://oauth2:${GITHUB_TOKEN}@github.com/.insteadOf", "https://github.com/"})

	// Stage and commit changes
	ctr = ctr.WithEnvVariable("GITHUB_TOKEN", tokenValue)
	ctr = ctr.WithExec([]string{"git", "add", "."})
	ctr = ctr.WithExec([]string{"git", "status"})
	ctr = ctr.WithExec([]string{"git", "commit", "-m", message, "--allow-empty"})

	// Push changes
	ctr = ctr.WithExec([]string{"git", "push", "origin", fmt.Sprintf("HEAD:%s", branch)})

	return ctr, nil
}

// NewGithubPublisher creates a new instance of the GithubPublisher
//
// Parameters:
// - dag: Dagger client instance
//
// Returns:
// - *GithubPublisher: A new GithubPublisher instance
func NewGithubPublisher(dag *dagger.Client) *GithubPublisher {
	return &GithubPublisher{
		dag: dag,
	}
}

func main() {
	// This is a Dagger module, not a standalone program
	// The main function is required to make it a valid Go program
}
