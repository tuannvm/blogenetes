// Package main is the root module for the Blogenetes Dagger pipeline.
// It composes all the individual modules (RSS, summarizer, markdown, github-publisher)
// into a complete workflow.
package main

import (
	"dagger.io/dagger"
)

// Blogenetes is the root module that composes all the individual modules
type Blogenetes struct {
	dag *dagger.Client
}

// New creates a new instance of the Blogenetes module
func NewBlogenetes(dag *dagger.Client) *Blogenetes {
	return &Blogenetes{
		dag: dag,
	}
}
