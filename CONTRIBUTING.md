# Contributing to gh-repo-visualize

Thank you for your interest in contributing to gh-repo-visualize! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

1. A clear title and description
2. Steps to reproduce the issue
3. Expected behavior
4. Actual behavior
5. Your environment (OS, Go version, etc.)

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:

1. A clear title describing the feature
2. Use case and benefits
3. Possible implementation approach (optional)

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes following conventional commits
6. Push to your branch
7. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Git
- Make (optional)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gh-repo-visualize.git
cd gh-repo-visualize

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

### Project Structure

```
gh-repo-visualize/
├── cmd/gh-repo-visualize/     # CLI entry point
├── internal/
│   ├── models/               # Data models
│   ├── git/                  # Git operations
│   ├── visualize/            # Rendering logic
│   ├── cmd/                  # Command handlers
│   ├── flags/                # CLI flags
│   └── errors/               # Error handling
└── .github/workflows/        # CI/CD
```

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `go vet` before committing
- Write tests for new functionality

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance

Example:
```
feat: add support for custom date formats

Add --date-format flag to allow users to customize
the date format in the output.
```

### Code Organization

- Keep the `cmd/` directory minimal - only CLI definitions
- Put business logic in `internal/` packages
- Use clear, descriptive names
- Add comments for exported functions

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./internal/git
```

## Release Process

Releases are automated via GitHub Actions:

1. Create a new tag: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. GitHub Actions will build and create a release

## Questions?

Feel free to open an issue for any questions or discussions.

Thank you for contributing! 🎉
