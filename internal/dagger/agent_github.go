package dagger

import (
	"context"
	"os"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

// AgentGitHub publishes markdown files to a GitHub repository.
type AgentGitHub struct {
	Client *github.Client
}

// NewAgentGitHub creates a GitHub agent using GITHUB_TOKEN from env.
func NewAgentGitHub() *AgentGitHub {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("GITHUB_TOKEN is not set")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	return &AgentGitHub{Client: client}
}

// CommitToRepo commits a markdown file to the repo under given path.
func (a *AgentGitHub) CommitToRepo(ctx context.Context, owner, repo, branch, path, content, message string) error {
	optsGet := &github.RepositoryContentGetOptions{Ref: branch}
	existingFile, _, resp, err := a.Client.Repositories.GetContents(ctx, owner, repo, path, optsGet)
	if err != nil && resp.StatusCode != 404 {
		return err
	}
	if existingFile != nil {
		// Update existing file
		updateOpts := &github.RepositoryContentFileOptions{
			Message: github.String(message),
			Content: []byte(content),
			SHA:     existingFile.SHA,
			Branch:  github.String(branch),
		}
		_, _, err := a.Client.Repositories.UpdateFile(ctx, owner, repo, path, updateOpts)
		return err
	}
	// Create new file
	createOpts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		Branch:  github.String(branch),
	}
	_, _, err = a.Client.Repositories.CreateFile(ctx, owner, repo, path, createOpts)
	return err
}
