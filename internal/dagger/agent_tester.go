package dagger

import (
	"io/ioutil"
	"path/filepath"
)

// AgentTester handles writing markdown files.
type AgentTester struct{}

// NewAgentTester creates a new AgentTester.
func NewAgentTester() *AgentTester {
	return &AgentTester{}
}

// WriteFile writes content to a specified directory and filename.
func (a *AgentTester) WriteFile(dir, filename string, content []byte) error {
	path := filepath.Join(dir, filename)
	return ioutil.WriteFile(path, content, 0644)
}
