# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release

## [1.0.0] - 2024-01-15

### Added
- `gh repo visualize` command to display commit graph
  - `--limit/-n` flag to limit number of commits
  - `--branch/-b` flag to filter by branch
  - `--author/-a` flag to filter by author
  - `--format/-f` flag for output format (ascii, compact, json)
- `gh repo visualize stats` command to display commit statistics
  - `--by-author` flag for author breakdown
  - `--by-day` flag for daily breakdown
- ASCII art visualization with box drawing characters
- JSON output format for scripting and integration
- Compact one-line output format
- Cross-platform builds (Linux, macOS, Windows)
- GitHub Actions CI/CD pipeline
- Automated releases with goreleaser

### Technical
- Go 1.22+ support
- urfave/cli/v2 framework
- Clean architecture with cmd/ and internal/ separation
- Custom error types for better error handling

[Unreleased]: https://github.com/h1s97x/gh-repo-visualize/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/h1s97x/gh-repo-visualize/releases/tag/v1.0.0
