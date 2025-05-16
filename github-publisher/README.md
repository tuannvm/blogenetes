# GitHub Publisher Module

A Dagger module for publishing files to GitHub repositories with proper Git history and authentication.

## Features

- Securely publish files to GitHub repositories
- Configurable Git author and committer information
- Handles Git operations in an isolated container
- Supports GitHub token authentication
- Preserves file permissions and directory structure

## Prerequisites

- [Dagger](https://docs.dagger.io/install) installed and configured
- A GitHub personal access token with `repo` scope

## Installation

```bash
dagger mod install github.com/tuannvm/blogenetes/modules/github_publisher
```

## Usage

```go
import "github.com/your-org/github-publisher"

// Create a new publisher instance
publisher := &GithubPublisher{}

// Publish a file
result, err := publisher.Publish(
    ctx,
    content,        // *dagger.File - The file to publish
    "my-repo",      // string - Repository name
    "my-org",       // string - Repository owner (username or org name)
    "main",         // string - Branch name
    "path/to/file", // string - Path in the repository
    "Update file",   // string - Commit message
    "",              // string - Author name (optional, defaults to "Blogenetes")
    "",              // string - Author email (optional, defaults to "blogenetes@example.com")
    "",              // string - Committer name (optional, defaults to author name)
    "",              // string - Committer email (optional, defaults to author email)
)
```

## Publishing to Daggerverse

To publish this module to Daggerverse and make it available to others, follow these steps:

1. **Login to Dagger** (if not already logged in):
   ```bash
   dagger login
   ```

2. **Navigate to the module directory**:
   ```bash
   cd /path/to/blogenetes/modules/github_publisher
   ```

3. **Publish the module**:
   ```bash
   dagger mod publish
   ```

4. **After publishing**, you can use the module from anywhere with:
   ```bash
   dagger mod install github.com/tuannvm/blogenetes/modules/github_publisher@<version>
   ```
   Replace `<version>` with the published version (e.g., `v0.1.0`).

## Environment Variables

- `GITHUB_TOKEN`: Required. A GitHub personal access token with `repo` scope

## Example

```go
// Create a file with some content
content := dag.Directory().WithNewFile("example.txt", "Hello, World!")

// Publish the file
result, err := publisher.Publish(
    ctx,
    content,
    "my-repo",
    "my-org",
    "main",
    "examples/hello.txt",
    "Add hello world example",
    "CI Bot",
    "ci@example.com",
    "",  // Committer name will be same as author
    "",  // Committer email will be same as author
)
```

## Development

### Building the Module

```bash
# Build the module
dagger mod sync
```

### Testing

```bash
# Run tests
dagger call --content=file:/path/to/test.txt --repo=test-repo --owner=your-username --branch=main --path=test.txt --message="Test commit"
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
