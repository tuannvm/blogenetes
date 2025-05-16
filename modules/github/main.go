package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type GitHub struct{}

// CommitToRepo commits content to a GitHub repository
func (g *GitHub) CommitToRepo(ctx context.Context, client *dagger.Client, owner, repo, branch, path, content, message string) error {
	// Get GitHub token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is not set")
	}

	repoURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", githubToken, owner, repo)

	// Create a directory with the content
	dir := client.Directory().WithNewFile(path, content)

	// Configure git and commit
	container := client.Container().
		From("alpine/git").
		WithMountedDirectory("/src", dir).
		WithWorkdir("/src").
		WithEnvVariable("GIT_AUTHOR_NAME", "Dagger Bot").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "dagger@example.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "Dagger Bot").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "dagger@example.com").
		WithExec([]string{"git", "init"}).
		WithExec([]string{"git", "remote", "add", "origin", repoURL}).
		WithExec([]string{"git", "fetch", "origin", branch}).
		WithExec([]string{"git", "checkout", "-B", branch, "origin/" + branch}).
		WithExec([]string{"git", "add", path}).
		WithExec([]string{"git", "commit", "-m", message}).
		WithExec([]string{"git", "push", "origin", branch})

	_, err := container.Stdout(ctx)
	return err
}

// ProcessGitHub is the Dagger function that will be called from the pipeline
func (g *GitHub) ProcessGitHub(ctx context.Context, client *dagger.Client, owner, repo, branch, path, content, message string) error {
	return g.CommitToRepo(ctx, client, owner, repo, branch, path, content, message)
}
