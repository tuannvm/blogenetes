package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Default().Writer()))
	if err != nil {
		log.Fatalf("failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Clone blog repo
	blog := client.Git("https://github.com/your-username/your-blog.git").Branch("main").Tree()

	// Prepare secrets
	openAISecret := client.SetSecret("OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY"))
	gitSecret := client.SetSecret("GIT_TOKEN", os.Getenv("GIT_TOKEN"))

	// Fetch & summarise
	scriptsDir := client.Host().Directory("scripts")
	
	// First test Go environment
	tester := client.Container().
		From("golang:1.20").
		WithMountedDirectory("/scripts", scriptsDir).
		WithExec([]string{"/bin/sh", "-c", "/scripts/test_go_env.sh"})
		
	// Get stdout for debugging
	testOutput, err := tester.Stdout(ctx)
	if err != nil {
		log.Printf("Go environment test failed: %v", err)
	} else {
		log.Printf("Go environment test result: %s", testOutput)
	}
	
	fetcher := client.Container().
		From("golang:1.20").
		WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod")).
		WithMountedDirectory("/scripts", scriptsDir).
		WithEnvVariable("RSS_URL", "https://some-blog.com/rss").
		WithEnvVariable("PATH", "/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin").
		WithEnvVariable("GOGC", "off").
		WithEnvVariable("GOPATH", "/go").
		WithSecretVariable("OPENAI_API_KEY", openAISecret).
		WithExec([]string{"/bin/sh", "-c", "mkdir -p /output && cd /scripts && which go && go mod tidy && go run fetch_and_summary.go"})

	fetchFile := fetcher.File("/output/new_post.md")

	// Read contents and inject generated markdown
	contents, err := fetchFile.Contents(ctx)
	if err != nil {
		log.Fatalf("failed to read generated summary file: %v", err)
	}
	updated := blog.WithNewFile("content/new_post.md", contents)

	// Commit & push via git container
	gitContainer := client.Container().
		From("alpine/git:latest").
		WithDirectory("/repo", updated).
		WithWorkdir("/repo").
		WithSecretVariable("GIT_TOKEN", gitSecret).
		WithExec([]string{"git", "config", "user.name", "dagger-bot"}).
		WithExec([]string{"git", "config", "user.email", "bot@example.com"}).
		WithExec([]string{"git", "add", "."}).
		WithExec([]string{"git", "commit", "-m", "🔖 Add summary of latest external blog post"}).
		WithExec([]string{"git", "push", fmt.Sprintf("https://x-access-token:%s@github.com/your-username/your-blog.git", gitSecret), "main"})

	// Actually run the last git push command and fail on error
	_, err = gitContainer.Sync(ctx)
	if err != nil {
		log.Fatalf("git push failed: %v", err)
	}

	fmt.Println("✅ Pipeline completed successfully")
}
