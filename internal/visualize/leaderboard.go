package visualize

import (
	"fmt"
	"sort"
	"strings"

	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// LeaderboardRenderer renders contributor leaderboard
type LeaderboardRenderer struct {
	options RenderOptions
	colors  *ColorScheme
}

// NewLeaderboardRenderer creates a new leaderboard renderer
func NewLeaderboardRenderer(opts RenderOptions) *LeaderboardRenderer {
	if opts.Width == 0 {
		opts.Width = 60
	}
	colors := NewColorScheme(opts.Color)
	return &LeaderboardRenderer{options: opts, colors: colors}
}

// Render renders the leaderboard
func (r *LeaderboardRenderer) Render(stats *models.Stats) string {
	if stats == nil || len(stats.Authors) == 0 {
		return "No contributors found"
	}

	var sb strings.Builder
	width := r.options.Width

	// Header
	border := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(border + "\n")

	title := " Top Contributors "
	if r.colors.Enabled {
		title = r.colors.Header.Render(title)
	}
	sb.WriteString(fmt.Sprintf("│%s│\n", centered(title, width-2)))

	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

	// Sort authors by commits
	authors := make([]*models.AuthorStats, len(stats.Authors))
	copy(authors, stats.Authors)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	// Limit to configured value or default 10
	limit := r.options.Limit
	if limit <= 0 {
		limit = 10
	}
	if len(authors) < limit {
		limit = len(authors)
	}

	for i := 0; i < limit; i++ {
		author := authors[i]
		line := r.renderAuthorLine(i+1, author, stats.TotalCommits)
		sb.WriteString(fmt.Sprintf("│%s│\n", padRight(line, width-2)))
	}

	// Footer
	borderBottom := r.colors.Border.Render("╰") + strings.Repeat("─", width-2) + r.colors.Border.Render("╯")
	sb.WriteString(borderBottom + "\n")

	return sb.String()
}

// renderAuthorLine renders a single author line
func (r *LeaderboardRenderer) renderAuthorLine(rank int, author *models.AuthorStats, total int) string {
	// Medal emoji for top 3
	var medal string
	switch rank {
	case 1:
		medal = r.colors.Feat.Render("🥇")
	case 2:
		medal = r.colors.Feat.Render("🥈")
	case 3:
		medal = r.colors.Feat.Render("🥉")
	default:
		medal = fmt.Sprintf("%d.", rank)
	}

	nameStyle := r.colors.Meta.Render(truncate(author.Name, 15))
	countStyle := r.colors.Feat.Render(fmt.Sprintf("%d", author.Commits))
	percentStyle := r.colors.Meta.Render(author.Percent)

	return fmt.Sprintf("  %s %-15s %s commits  %s",
		medal,
		nameStyle,
		countStyle,
		percentStyle,
	)
}

// RenderCompact renders a compact leaderboard
func (r *LeaderboardRenderer) RenderCompact(stats *models.Stats) string {
	if stats == nil || len(stats.Authors) == 0 {
		return "No contributors found"
	}

	var sb strings.Builder

	// Sort authors by commits
	authors := make([]*models.AuthorStats, len(stats.Authors))
	copy(authors, stats.Authors)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	// Header
	sb.WriteString(" #  Author                  Commits  %\n")
	sb.WriteString(strings.Repeat("-", 45) + "\n")

	// Limit to top 10
	limit := 10
	if len(authors) < limit {
		limit = len(authors)
	}

	for i := 0; i < limit; i++ {
		author := authors[i]
		medal := ""
		switch i + 1 {
		case 1:
			medal = "🥇 "
		case 2:
			medal = "🥈 "
		case 3:
			medal = "🥉 "
		default:
			medal = fmt.Sprintf("%2d. ", i+1)
		}
		sb.WriteString(fmt.Sprintf("%s%-20s %5d     %s\n",
			medal,
			truncate(author.Name, 20),
			author.Commits,
			author.Percent,
		))
	}

	return sb.String()
}

// RenderJSON renders as JSON
func (r *LeaderboardRenderer) RenderJSON(stats *models.Stats) string {
	if stats == nil {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf(`  "total_contributors": %d,
  "total_commits": %d,
  "leaderboard": [
`, len(stats.Authors), stats.TotalCommits))

	// Sort authors by commits
	authors := make([]*models.AuthorStats, len(stats.Authors))
	copy(authors, stats.Authors)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	for i, author := range authors {
		if i >= 10 {
			break
		}
		sb.WriteString(fmt.Sprintf(`    {
      "rank": %d,
      "name": "%s",
      "commits": %d,
      "percent": "%s"
    }`,
			i+1,
			escapeJSON(author.Name),
			author.Commits,
			author.Percent,
		))
		if i < len(authors)-1 && i < 9 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("  ]\n")
	sb.WriteString("}\n")
	return sb.String()
}
