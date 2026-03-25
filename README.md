# gh-repo-visualize

<p align="center">
  <strong>A GitHub CLI extension for visualizing Git commit history</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#examples">Examples</a> •
  <a href="#contributing">Contributing</a>
</p>

---

## Features

- 📊 **Commit Graph Visualization** - Display Git commit history as ASCII art
- 📈 **Commit Statistics** - Analyze commit patterns by author, day, and time
- 🎨 **Multiple Output Formats** - ASCII, compact, and JSON formats
- 🔍 **Flexible Filtering** - Filter by branch, author, or limit count
- ⚡ **Fast and Lightweight** - Pure Go implementation with no dependencies

## Installation

### Using gh CLI (Recommended)

```bash
gh extension install h1s97x/gh-repo-visualize
```

### Manual Installation

Download the latest release for your platform:

```bash
# Linux/macOS
curl -sL https://github.com/h1s97x/gh-repo-visualize/releases/latest/download/gh-repo-visualize-$(uname -s)-$(uname -m).tar.gz | tar xz
chmod +x gh-repo-visualize
sudo mv gh-repo-visualize /usr/local/bin/

# Or using Go
go install github.com/h1s97x/gh-repo-visualize/cmd/gh-repo-visualize@latest
```

## Usage

### Show Commit Graph

```bash
# Display commit graph for current repository
gh repo visualize

# Limit to last 50 commits
gh repo visualize -n 50

# Show commits from a specific branch
gh repo visualize -b main

# Show commits by a specific author
gh repo visualize -a "John Doe"

# Output in different formats
gh repo visualize -f json     # JSON format
gh repo visualize -f compact  # Compact one-line format
gh repo visualize -f ascii    # ASCII art (default)
```

### Show Commit Statistics

```bash
# Display commit statistics
gh repo visualize stats

# Show breakdown by author
gh repo visualize stats --by-author

# Show breakdown by day
gh repo visualize stats --by-day

# Output as JSON
gh repo visualize stats -f json
```

## Flags

### `gh repo visualize`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-n` | Number of commits to display | `20` |
| `--branch` | `-b` | Filter by branch | (all branches) |
| `--author` | `-a` | Filter by author | (all authors) |
| `--format` | `-f` | Output format: `ascii`, `compact`, `json` | `ascii` |

### `gh repo visualize stats`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--by-author` | | Show breakdown by author | `false` |
| `--by-day` | | Show breakdown by day | `false` |
| `--branch` | `-b` | Filter by branch | (all branches) |
| `--author` | `-a` | Filter by author | (all authors) |
| `--format` | `-f` | Output format: `ascii`, `json` | `ascii` |

## Examples

### Commit Graph Output

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Commit     Author                   Date         Message                   │
├──────────────────────────────────────────────────────────────────────────────┤
│ ● 53a8d8c John Doe         2024-01-15 feat: Add user authentication
│ ● c3ee2f7 Jane Smith       2024-01-14 fix: Resolve login issue
│ ● a1b2c3d John Doe         2024-01-13 docs: Update README
│   └─ Merge: feature/auth
╰──────────────────────────────────────────────────────────────────────────────╯
```

### Statistics Output

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Commit Statistics                                                            │
├──────────────────────────────────────────────────────────────────────────────┤
│ Total Commits: 156                                                           │
│ Date Range: 2023-06-01 to 2024-01-15                                         │
│ Avg Commits/Day: 0.7                                                         │
│ Most Active Day: Tuesday                                                     │
│ Most Active Hour: 14:00                                                      │
├──────────────────────────────────────────────────────────────────────────────┤

╭──────────────────────────────────────────────────────────────────────────────╮
│ Commits by Author                                                            │
├──────────────────────────────────────────────────────────────────────────────┤
│   John Doe                   89 commits (57.1%)
│   ███████████████████████████████████████████████████████░░░░░░░░░░░░░░░░░░░░
│   Jane Smith                 45 commits (28.8%)
│   ██████████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
╰──────────────────────────────────────────────────────────────────────────────╯
```

### JSON Output

```json
[
  {
    "hash": "53a8d8c99b1ee274dcddc15ce0730b8bebf1ef3e",
    "short_hash": "53a8d8c",
    "author": "John Doe",
    "date": "2024-01-15",
    "message": "feat: Add user authentication"
  }
]
```

## Development

### Prerequisites

- Go 1.22 or higher
- Make (optional)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/h1s97x/gh-repo-visualize.git
cd gh-repo-visualize

# Build
make build

# Run tests
make test

# Install locally
make install
```

### Project Structure

```
gh-repo-visualize/
├── cmd/gh-repo-visualize/     # CLI entry point
│   ├── main.go               
│   └── version.go            
├── internal/
│   ├── models/               # Data models
│   ├── git/                  # Git operations
│   ├── visualize/            # Rendering logic
│   ├── cmd/                  # Command handlers
│   ├── flags/                # CLI flags
│   └── errors/               # Error handling
├── .github/workflows/        # CI/CD
├── go.mod
├── Makefile
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [urfave/cli](https://github.com/urfave/cli)
- Inspired by `git log --graph`
