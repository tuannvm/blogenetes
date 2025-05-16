package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/tuannvm/blogenetes/modules/github_publisher/github_publisher"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	publisher := github_publisher.New().WithClient(client)
	fmt.Println("GitHub Publisher module initialized")
	_ = publisher
}
