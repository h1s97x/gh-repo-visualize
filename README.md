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
- 🌳 **Branch Graph Visualization** - Visualize branch structure and relationships
- 👥 **Contributor Leaderboard** - Rank contributors by commit count
- 📅 **Date Range Filtering** - Filter commits by date with --since/--until flags
- 📤 **Export Formats** - Export data as CSV, Markdown, or HTML
- 🎨 **Terminal Colors** - Rich color output with multiple color schemes
- 🔍 **Flexible Filtering** - Filter by branch, author, or limit count (case-insensitive)
- ⚡ **Fast and Lightweight** - Pure Go implementation

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

# Color output (auto-detected in terminal)
gh repo visualize --color     # Force enable colors
gh repo visualize --no-color  # Disable colors
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

# Color output
gh repo visualize stats --color     # Force enable colors
gh repo visualize stats --no-color  # Disable colors
```

### Show Branch Graph

```bash
# Display branch visualization
gh repo visualize branch

# Show local and remote branches
gh repo visualize branch --all

# Show branch with merged status
gh repo visualize branch --merged

# Show unmerged branches only
gh repo visualize branch --unmerged

# Output as JSON
gh repo visualize branch -f json
```

### Show Contributor Leaderboard

```bash
# Display top 10 contributors
gh repo visualize leaderboard

# Show top 5 contributors
gh repo visualize leaderboard -n 5

# Use alias
gh repo visualize top -n 20

# Output as JSON
gh repo visualize leaderboard -f json
```

### Date Range Filtering

```bash
# Filter commits since a specific date
gh repo visualize --since "2024-01-01"
gh repo visualize stats --since "2024-01-01"

# Filter commits until a specific date
gh repo visualize --until "2024-12-31"

# Combined date range
gh repo visualize --since "2024-06-01" --until "2024-12-31"
```

### Export Data

```bash
# Export commit graph as CSV
gh repo visualize -f csv -o commits.csv

# Export commit graph as Markdown
gh repo visualize -f markdown -o commits.md

# Export commit graph as HTML
gh repo visualize -f html -o commits.html

# Export statistics
gh repo visualize stats -f csv -o stats.csv
```

## Flags

### `gh repo visualize`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-n` | Number of commits to display | `20` |
| `--branch` | `-b` | Filter by branch | (all branches) |
| `--author` | `-a` | Filter by author (case-insensitive) | (all authors) |
| `--since` | | Show commits since this date (YYYY-MM-DD) | (all time) |
| `--until` | | Show commits until this date (YYYY-MM-DD) | (all time) |
| `--format` | `-f` | Output format: `ascii`, `compact`, `json`, `csv`, `markdown`, `html` | `ascii` |
| `--output` | `-o` | Output file path for exports | (stdout) |
| `--color` | | Enable colored output | auto-detect |
| `--no-color` | | Disable colored output | auto-detect |

### `gh repo visualize stats`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--by-author` | | Show breakdown by author | `false` |
| `--by-day` | | Show breakdown by day | `false` |
| `--branch` | `-b` | Filter by branch | (all branches) |
| `--author` | `-a` | Filter by author (case-insensitive) | (all authors) |
| `--since` | | Show commits since this date (YYYY-MM-DD) | (all time) |
| `--until` | | Show commits until this date (YYYY-MM-DD) | (all time) |
| `--format` | `-f` | Output format: `ascii`, `json`, `csv` | `ascii` |
| `--output` | `-o` | Output file path for exports | (stdout) |
| `--color` | | Enable colored output | auto-detect |
| `--no-color` | | Disable colored output | auto-detect |

### `gh repo visualize branch`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--all` | `-a` | Show all branches (local and remote) | `false` |
| `--merged` | | Show only merged branches | `false` |
| `--unmerged` | | Show only unmerged branches | `false` |
| `--format` | `-f` | Output format: `ascii`, `json` | `ascii` |
| `--color` | | Enable colored output | auto-detect |
| `--no-color` | | Disable colored output | auto-detect |

### `gh repo visualize leaderboard` / `gh repo visualize top`

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-n` | Number of contributors to display | `10` |
| `--since` | | Show commits since this date (YYYY-MM-DD) | (all time) |
| `--until` | | Show commits until this date (YYYY-MM-DD) | (all time) |
| `--format` | `-f` | Output format: `ascii`, `compact`, `json` | `ascii` |
| `--color` | | Enable colored output | auto-detect |
| `--no-color` | | Disable colored output | auto-detect |

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

### Branch Graph Output

```
╭────────────────────────────────────────────────────────╮
│                    Branch Visualization                 │
├────────────────────────────────────────────────────────┤
│ ● main        (HEAD) - 7 commits ahead, 2 behind      │
│   ↳ feature/auth    - 3 commits ahead                 │
│   ↳ feature/api     - 5 commits ahead                 │
│   ↳ bugfix/login    - merged into main                │
╰────────────────────────────────────────────────────────╯
```

### Leaderboard Output

```
╭────────────────────────────────────────────────────────╮
│              Top Contributors                           │
├────────────────────────────────────────────────────────┤
│   🥇 John Doe           156 commits   32.5%             │
│   🥈 Jane Smith          98 commits   20.4%             │
│   🥉 Bob Wilson          76 commits   15.8%             │
│   4. Alice Brown         52 commits   10.8%             │
│   5. Charlie Davis       45 commits   9.4%              │
╰────────────────────────────────────────────────────────╯
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

## Roadmap / Planned Features

We're always looking to improve! Check out our planned features:

- [#5](https://github.com/h1s97x/gh-repo-visualize/issues/5) - Branch graph visualization
- [#6](https://github.com/h1s97x/gh-repo-visualize/issues/6) - Date range filtering (--since/--until)
- [#7](https://github.com/h1s97x/gh-repo-visualize/issues/7) - Export to CSV, Markdown, HTML
- [#8](https://github.com/h1s97x/gh-repo-visualize/issues/8) - Contributor leaderboard

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
