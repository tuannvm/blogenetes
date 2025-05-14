.PHONY: lint fmt build all

LINT_TARGETS := $(shell find . -name '*.go')
FMT_TARGETS := $(LINT_TARGETS)
GOFMT := gofmt
GOLINT := golangci-lint

all: fmt lint build

lint:
	@command -v $(GOLINT) >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1; }
	$(GOLINT) run ./...

fmt:
	$(GOFMT) -s -w $(FMT_TARGETS)

build:
	go build ./...

