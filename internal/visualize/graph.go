package visualize

import (
	"fmt"
	"strings"

	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// RenderOptions holds options for rendering
type RenderOptions struct {
	Format    string // "ascii", "json", "dot"
	Width     int
	ShowGraph bool
}

// GraphRenderer renders commit graphs
type GraphRenderer struct {
	options RenderOptions
}

// NewGraphRenderer creates a new graph renderer
func NewGraphRenderer(opts RenderOptions) *GraphRenderer {
	if opts.Width == 0 {
		opts.Width = 80
	}
	return &GraphRenderer{options: opts}
}

// Render renders a commit list as ASCII graph
func (r *GraphRenderer) Render(commits []*models.Commit) string {
	if len(commits) == 0 {
		return "No commits to display"
	}

	var sb strings.Builder
	width := r.options.Width

	// Header
	sb.WriteString("╭" + strings.Repeat("─", width-2) + "╮\n")
	header := fmt.Sprintf("│ %-10s %-8s %-15s %-12s %s",
		"Commit", "Author", "", "Date", "Message")
	sb.WriteString(padRight(header, width) + "│\n")
	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

	// Commits
	for _, commit := range commits {
		// Graph line
		line := r.renderLine(commit)
		sb.WriteString("│ " + line + "\n")
		
		// Merge info
		if len(commit.Parents) > 1 {
			// Show parent hashes (shortened)
			parentShortHashes := make([]string, len(commit.Parents))
			for i, p := range commit.Parents {
				if len(p) >= 7 {
					parentShortHashes[i] = p[:7]
				} else {
					parentShortHashes[i] = p
				}
			}
			mergeInfo := fmt.Sprintf("│   └─ Merge: %s", strings.Join(parentShortHashes, ", "))
			sb.WriteString(padRight(mergeInfo, width) + "│\n")
		}
	}

	// Footer
	sb.WriteString("╰" + strings.Repeat("─", width-2) + "╯\n")

	return sb.String()
}

// renderLine renders a single commit line
func (r *GraphRenderer) renderLine(commit *models.Commit) string {
	// Graph characters: │ ╭ ╮ ╰ ╯ ├ ┤ ┬ ┴ ┼ ─ • ○ ● ★
	graph := "●" // Current commit
	
	// Format: graph short_hash author date message
	msg := models.TruncateMessage(commit.Message, 30)
	return fmt.Sprintf("%s %s %-15s %s %s",
		graph,
		commit.ShortHash,
		truncate(commit.Author, 15),
		commit.Date.Format("2006-01-02"),
		msg,
	)
}

// RenderCompact renders a compact view
func (r *GraphRenderer) RenderCompact(commits []*models.Commit) string {
	if len(commits) == 0 {
		return "No commits to display"
	}

	var sb strings.Builder
	
	for _, commit := range commits {
		sb.WriteString(fmt.Sprintf("%s %s %s %s\n",
			commit.ShortHash,
			commit.Date.Format("2006-01-02"),
			truncate(commit.Author, 12),
			models.TruncateMessage(commit.Message, 40),
		))
	}
	
	return sb.String()
}

// RenderJSON renders as JSON
func (r *GraphRenderer) RenderJSON(commits []*models.Commit) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	
	for i, commit := range commits {
		sb.WriteString(fmt.Sprintf(`  {
    "hash": "%s",
    "short_hash": "%s",
    "author": "%s",
    "date": "%s",
    "message": "%s"
  }`,
			commit.Hash,
			commit.ShortHash,
			commit.Author,
			commit.Date.Format("2006-01-02"),
			escapeJSON(commit.Message),
		))
		if i < len(commits)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	
	sb.WriteString("]\n")
	return sb.String()
}

// Helper functions
func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func padRight(s string, width int) string {
	if len(s) >= width-1 {
		return s[:width-1]
	}
	return s + strings.Repeat(" ", width-1-len(s))
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
