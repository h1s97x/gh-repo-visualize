package visualize

import (
	"fmt"
	"sort"
	"strings"

	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// StatsRenderer renders commit statistics
type StatsRenderer struct {
	width int
}

// NewStatsRenderer creates a new stats renderer
func NewStatsRenderer() *StatsRenderer {
	return &StatsRenderer{width: 80}
}

// Render renders statistics as ASCII table
func (r *StatsRenderer) Render(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	sb.WriteString("╭" + strings.Repeat("─", width-2) + "╮\n")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", "Commit Statistics"))
	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

	// Summary
	sb.WriteString(fmt.Sprintf("│ Total Commits: %-63d │\n", stats.TotalCommits))
	
	if stats.DateRange != nil {
		sb.WriteString(fmt.Sprintf("│ Date Range: %-66s │\n", 
			stats.DateRange.Start+" to "+stats.DateRange.End))
		sb.WriteString(fmt.Sprintf("│ Avg Commits/Day: %-59.1f │\n", stats.AvgPerDay))
	}
	
	if stats.MostActiveDay != "" {
		sb.WriteString(fmt.Sprintf("│ Most Active Day: %-60s │\n", stats.MostActiveDay))
	}
	
	sb.WriteString(fmt.Sprintf("│ Most Active Hour: %02d:00%-56s │\n", stats.MostActiveHour, ""))

	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

	return sb.String()
}

// RenderByAuthor renders statistics grouped by author
func (r *StatsRenderer) RenderByAuthor(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	sb.WriteString("╭" + strings.Repeat("─", width-2) + "╮\n")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", "Commits by Author"))
	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

	// Sort authors by commits
	authors := make([]*models.AuthorStats, len(stats.Authors))
	copy(authors, stats.Authors)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	for _, author := range authors {
		line := fmt.Sprintf("  %-25s %5d commits (%s)", 
			truncate(author.Name, 25), 
			author.Commits, 
			author.Percent)
		sb.WriteString(fmt.Sprintf("│%-78s│\n", line))
		
		// Bar chart
		bar := r.renderBar(author.Commits, stats.TotalCommits, 50)
		sb.WriteString(fmt.Sprintf("│   %-74s│\n", bar))
	}

	sb.WriteString("╰" + strings.Repeat("─", width-2) + "╯\n")

	return sb.String()
}

// RenderByDay renders statistics grouped by day
func (r *StatsRenderer) RenderByDay(stats *models.Stats) string {
	var sb strings.Builder
	width := r.width

	// Header
	sb.WriteString("╭" + strings.Repeat("─", width-2) + "╮\n")
	sb.WriteString(fmt.Sprintf("│ %-76s │\n", "Commits by Day"))
	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")

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
	for date, count := range stats.ByDate {
		if date >= dates[start] && count > maxCommits {
			maxCommits = count
		}
	}

	for i := start; i < len(dates); i++ {
		date := dates[i]
		count := stats.ByDate[date]
		line := fmt.Sprintf("  %s %5d %s", date, count, r.renderSimpleBar(count, maxCommits, 30))
		sb.WriteString(fmt.Sprintf("│%-78s│\n", line))
	}

	sb.WriteString("╰" + strings.Repeat("─", width-2) + "╯\n")

	return sb.String()
}

// RenderJSON renders statistics as JSON
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

// renderBar renders a bar chart
func (r *StatsRenderer) renderBar(value, total, width int) string {
	if total == 0 {
		return ""
	}
	
	filled := value * width / total
	bar := strings.Repeat("█", filled)
	bar += strings.Repeat("░", width-filled)
	return bar
}

// renderSimpleBar renders a simple bar
func (r *StatsRenderer) renderSimpleBar(value, max, width int) string {
	if max == 0 {
		return ""
	}
	
	filled := value * width / max
	return strings.Repeat("▓", filled) + strings.Repeat("░", width-filled)
}
