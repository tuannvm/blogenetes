package main

import (
	"fmt"
	"os"

	"github.com/tuannvm/blogenetes/cmd/blogenetes"
)

func main() {
	// Delegate to the actual command implementation
	if err := blogenetes.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
