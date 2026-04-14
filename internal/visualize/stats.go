package visualize

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// StatsRenderer renders commit statistics
type StatsRenderer struct {
	width  int
	colors *ColorScheme
}

// NewStatsRenderer creates a new stats renderer
func NewStatsRenderer() *StatsRenderer {
	return &StatsRenderer{width: 80, colors: NewColorScheme(false)}
}

// NewStatsRendererWithColor creates a new stats renderer with color support
func NewStatsRendererWithColor(color bool) *StatsRenderer {
	return &StatsRenderer{width: 80, colors: NewColorScheme(color)}
}

// Render renders statistics as ASCII table
func (r *StatsRenderer) Render(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	borderTop := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(borderTop + "\n")
	
	title := r.colors.Header.Render("Commit Statistics")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", title))
	
	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

	// Summary
	totalCommits := r.colors.Feat.Render(fmt.Sprintf("%d", stats.TotalCommits))
	sb.WriteString(fmt.Sprintf("│ Total Commits: %-63s │\n", totalCommits))
	
	if stats.DateRange != nil {
		dateRange := r.colors.Meta.Render(stats.DateRange.Start + " to " + stats.DateRange.End)
		sb.WriteString(fmt.Sprintf("│ Date Range: %-66s │\n", dateRange))
		avg := r.colors.Meta.Render(fmt.Sprintf("%.1f", stats.AvgPerDay))
		sb.WriteString(fmt.Sprintf("│ Avg Commits/Day: %-59s │\n", avg))
	}
	
	if stats.MostActiveDay != "" {
		day := r.colors.Meta.Render(stats.MostActiveDay)
		sb.WriteString(fmt.Sprintf("│ Most Active Day: %-60s │\n", day))
	}
	
	hour := r.colors.Meta.Render(fmt.Sprintf("%02d:00", stats.MostActiveHour))
	sb.WriteString(fmt.Sprintf("│ Most Active Hour: %-59s │\n", hour))

	sb.WriteString(borderMid + "\n")

	return sb.String()
}

// RenderByAuthor renders statistics grouped by author
func (r *StatsRenderer) RenderByAuthor(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	borderTop := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(borderTop + "\n")
	
	title := r.colors.Header.Render("Commits by Author")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", title))
	
	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

	// Sort authors by commits
	authors := make([]*models.AuthorStats, len(stats.Authors))
	copy(authors, stats.Authors)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	for _, author := range authors {
		name := r.colors.Feat.Render(truncate(author.Name, 25))
		count := r.colors.Fix.Render(fmt.Sprintf("%5d", author.Commits))
		percent := r.colors.Meta.Render(author.Percent)
		line := fmt.Sprintf("  %s %s commits (%s)", name, count, percent)
		sb.WriteString(fmt.Sprintf("│%-78s│\n", line))
		
		// Bar chart
		bar := r.colors.RenderBar(author.Commits, stats.TotalCommits, 50)
		sb.WriteString(fmt.Sprintf("│   %-74s│\n", bar))
	}

	borderBottom := r.colors.Border.Render("╰") + strings.Repeat("─", width-2) + r.colors.Border.Render("╯")
	sb.WriteString(borderBottom + "\n")

	return sb.String()
}

// RenderByDay renders statistics grouped by day
func (r *StatsRenderer) RenderByDay(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	borderTop := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(borderTop + "\n")
	
	title := r.colors.Header.Render("Commits by Day")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", title))
	
	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

	// Sort dates
	dates := make([]string, 0, len(stats.ByDate))
	for date := range stats.ByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Show last 14 days max
	start := 0
	if len(dates) > 14 {
		start = len(dates) - 14
	}

	maxCommits := 0
	for _, date := range dates[start:] {
		if stats.ByDate[date] > maxCommits {
			maxCommits = stats.ByDate[date]
		}
	}

	for i := start; i < len(dates); i++ {
		date := dates[i]
		count := stats.ByDate[date]
		dateStr := r.colors.Meta.Render(date)
		countStr := r.colors.Fix.Render(fmt.Sprintf("%5d", count))
		bar := r.colors.RenderBar(count, maxCommits, 30)
		line := fmt.Sprintf("  %s %s %s", dateStr, countStr, bar)
		sb.WriteString(fmt.Sprintf("│%-78s│\n", line))
	}

	borderBottom := r.colors.Border.Render("╰") + strings.Repeat("─", width-2) + r.colors.Border.Render("╯")
	sb.WriteString(borderBottom + "\n")

	return sb.String()
}

// RenderJSON renders statistics as JSON (no colors)
func (r *StatsRenderer) RenderJSON(stats *models.Stats) string {
	var sb strings.Builder
	
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"total_commits\": %d,\n", stats.TotalCommits))
	
	if stats.DateRange != nil {
		sb.WriteString(fmt.Sprintf("  \"date_range\": {\"start\": \"%s\", \"end\": \"%s\"},\n",
			stats.DateRange.Start, stats.DateRange.End))
	}
	
	sb.WriteString(fmt.Sprintf("  \"avg_per_day\": %.1f,\n", stats.AvgPerDay))
	sb.WriteString(fmt.Sprintf("  \"most_active_day\": \"%s\",\n", stats.MostActiveDay))
	sb.WriteString(fmt.Sprintf("  \"most_active_hour\": %d,\n", stats.MostActiveHour))
	
	sb.WriteString("  \"authors\": [\n")
	for i, author := range stats.Authors {
		sb.WriteString(fmt.Sprintf("    {\"name\": \"%s\", \"commits\": %d, \"percent\": \"%s\"}",
			author.Name, author.Commits, author.Percent))
		if i < len(stats.Authors)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  ]\n")
	
	sb.WriteString("}\n")
	return sb.String()
}
