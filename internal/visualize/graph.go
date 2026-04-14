package visualize

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// RenderOptions holds options for rendering
type RenderOptions struct {
	Format    string // "ascii", "json", "dot"
	Width     int
	ShowGraph bool
	Color     bool   // Enable color output
	Limit     int    // Limit for leaderboard, etc.
}

// GraphRenderer renders commit graphs
type GraphRenderer struct {
	options RenderOptions
	colors  *ColorScheme
}

// NewGraphRenderer creates a new graph renderer
func NewGraphRenderer(opts RenderOptions) *GraphRenderer {
	if opts.Width == 0 {
		opts.Width = 80
	}
	colors := NewColorScheme(opts.Color)
	return &GraphRenderer{options: opts, colors: colors}
}

// Render renders a commit list as ASCII graph
func (r *GraphRenderer) Render(commits []*models.Commit) string {
	if len(commits) == 0 {
		return "No commits to display"
	}

	var sb strings.Builder
	width := r.options.Width

	// Header
	border := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(border + "\n")
	
	header := fmt.Sprintf("│ %-10s %-8s %-15s %-12s %s",
		"Commit", "Author", "", "Date", "Message")
	if r.colors.Enabled {
		header = r.colors.Header.Render(header)
	}
	sb.WriteString(padRight(header, width) + "│\n")
	
	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

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
			mergeStyle := r.colors.Merge
			if !r.colors.Enabled {
				mergeStyle = lipgloss.NewStyle()
			}
			mergeInfo := "│   " + mergeStyle.Render("└─ Merge: ") + strings.Join(parentShortHashes, ", ")
			sb.WriteString(padRight(mergeInfo, width) + "│\n")
		}
	}

	// Footer
	borderBottom := r.colors.Border.Render("╰") + strings.Repeat("─", width-2) + r.colors.Border.Render("╯")
	sb.WriteString(borderBottom + "\n")

	return sb.String()
}

// renderLine renders a single commit line
func (r *GraphRenderer) renderLine(commit *models.Commit) string {
	// Graph characters with color
	graph := "●"
	if r.colors.Enabled {
		graph = r.colors.Feat.Render("●")
	}
	
	// Format: graph short_hash author date message
	msg := models.TruncateMessage(commit.Message, 30)
	msg = r.colors.GetCommitStyle(commit.Message).Render(msg)
	
	hashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
	if r.colors.Enabled {
		hashStyle = hashStyle.Bold(true)
	}
	
	return fmt.Sprintf("%s %s %-15s %s %s",
		graph,
		hashStyle.Render(commit.ShortHash),
		r.colors.Meta.Render(truncate(commit.Author, 15)),
		r.colors.Meta.Render(commit.Date.Format("2006-01-02")),
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
		msg := r.colors.GetCommitStyle(commit.Message).Render(models.TruncateMessage(commit.Message, 40))
		sb.WriteString(fmt.Sprintf("%s %s %s %s\n",
			commit.ShortHash,
			commit.Date.Format("2006-01-02"),
			truncate(commit.Author, 12),
			msg,
		))
	}
	
	return sb.String()
}

// RenderJSON renders as JSON (no colors)
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

// RenderCSV renders commits as CSV format
func (r *GraphRenderer) RenderCSV(commits []*models.Commit) string {
	var sb strings.Builder
	
	// Header
	sb.WriteString("Hash,ShortHash,Author,Email,Date,Message,Parents\n")
	
	for _, commit := range commits {
		parents := strings.Join(commit.Parents, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,\"%s\",%s\n",
			commit.Hash,
			commit.ShortHash,
			escapeCSV(commit.Author),
			escapeCSV(commit.Email),
			commit.Date.Format("2006-01-02 15:04:05"),
			escapeCSV(commit.Message),
			parents,
		))
	}
	
	return sb.String()
}

// RenderMarkdown renders commits as Markdown table
func (r *GraphRenderer) RenderMarkdown(commits []*models.Commit) string {
	var sb strings.Builder
	
	// Title
	sb.WriteString("# Commit History\n\n")
	
	// Summary
	sb.WriteString(fmt.Sprintf("**Total Commits:** %d\n\n", len(commits)))
	
	// Table header
	sb.WriteString("| Commit | Author | Date | Message |\n")
	sb.WriteString("|--------|--------|------|--------|\n")
	
	// Table rows
	for _, commit := range commits {
		shortMsg := models.TruncateMessage(commit.Message, 50)
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
			commit.ShortHash,
			escapeMarkdown(commit.Author),
			commit.Date.Format("2006-01-02"),
			escapeMarkdown(shortMsg),
		))
	}
	
	return sb.String()
}

// RenderHTML renders commits as HTML table
func (r *GraphRenderer) RenderHTML(commits []*models.Commit) string {
	var sb strings.Builder
	
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Commit History</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 40px; }
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
th { background-color: #f4f4f4; }
tr:nth-child(even) { background-color: #fafafa; }
.commit { font-family: monospace; color: #0366d6; }
.message { max-width: 400px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
</head>
<body>
<h1>Commit History</h1>
<p><strong>Total Commits:</strong> `)
	sb.WriteString(fmt.Sprintf("%d", len(commits)))
	sb.WriteString(`</p>
<table>
<thead>
<tr>
<th>Commit</th>
<th>Author</th>
<th>Date</th>
<th>Message</th>
</tr>
</thead>
<tbody>
`)
	
	for _, commit := range commits {
		shortMsg := models.TruncateMessage(commit.Message, 60)
		sb.WriteString(fmt.Sprintf(`<tr>
<td class="commit">%s</td>
<td>%s</td>
<td>%s</td>
<td class="message" title="%s">%s</td>
</tr>
`,
			commit.ShortHash,
			escapeHTML(commit.Author),
			commit.Date.Format("2006-01-02"),
			escapeHTML(commit.Message),
			escapeHTML(shortMsg),
		))
	}
	
	sb.WriteString(`</tbody>
</table>
</body>
</html>
`)
	
	return sb.String()
}

// escapeCSV escapes a string for CSV output
func escapeCSV(s string) string {
	s = strings.ReplaceAll(s, "\"", "\"\"")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// escapeMarkdown escapes a string for Markdown output
func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// escapeHTML escapes a string for HTML output
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
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
	// Strip ANSI codes for length calculation
	clean := stripANSI(s)
	if len(clean) >= width-1 {
		return s
	}
	padding := width - 1 - len(clean)
	return s + strings.Repeat(" ", padding)
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	var result strings.Builder
	inANSI := false
	for _, r := range s {
		if r == '\x1b' {
			inANSI = true
			continue
		}
		if inANSI && r == 'm' {
			inANSI = false
			continue
		}
		if !inANSI {
			result.WriteRune(r)
		}
	}
	return result.String()
}
