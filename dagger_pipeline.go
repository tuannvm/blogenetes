package main

import (
	"context"
	"fmt"
	"log"

		"github.com/dagger/dagger-go"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Default().Writer()))
	if err != nil {
		log.Fatalf("failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// 1) Clone your blog repo
	blog := client.
		Git("https://github.com/your-username/your-blog.git").
		Branch("main").
		Tree()

	// 2) Spin up a Go container to fetch & summarise
	fetcher := client.Container().
		From("golang:1.20").
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod")).
		        WithMountedFile("/scripts/fetch_and_summary.go", client.Host().File("scripts/fetch_and_summary.go")).
		        WithMountedFile("/scripts/go.mod", client.Host().File("scripts/go.mod")).
		WithEnvVariable("RSS_URL", "https://some-blog.com/rss").
		WithSecretVariable("OPENAI_API_KEY", client.Host().Secret("OPENAI_API_KEY")).
		        WithExec([]string{"bash", "-lc", "cd /scripts && go mod tidy"}).
		        WithExec([]string{"bash", "-lc", "cd /scripts && go run fetch_and_summary.go"}).
		WithFile("/output/new_post.md")

	// 3) Inject the generated markdown into your blog tree
	updated := blog.WithNewFile("/content/new_post.md", fetcher.File("/output/new_post.md"))

	// 4) Commit & push back to repo
	gitToken := client.Host().Secret("GIT_TOKEN")
	runner := updated.
		WithExec([]string{"git", "config", "user.name", "dagger-bot"}).
		WithExec([]string{"git", "config", "user.email", "bot@example.com"}).
		WithExec([]string{"git", "add", "."}).
		WithExec([]string{"git", "commit", "-m", "🔖 Add summary of latest external blog post"}).
		WithSecretVariable("GIT_TOKEN", gitToken).
		WithExec([]string{"git", "push", fmt.Sprintf("https://x-access-token:%s@github.com/your-username/your-blog.git", gitToken), "main"})

	if exit, err := runner.ExitCode(ctx); err != nil || exit != 0 {
		log.Fatalf("push failed (exit=%d): %v", exit, err)
	}

	fmt.Println("✅ Pipeline completed successfully")
}