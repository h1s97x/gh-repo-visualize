package visualize

import (
	"fmt"
	"sort"
	"strings"

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

// RenderCSV renders statistics as CSV format
func (r *StatsRenderer) RenderCSV(stats *models.Stats) string {
	var sb strings.Builder
	
	// Summary section
	sb.WriteString("# Summary\n")
	sb.WriteString("Metric,Value\n")
	sb.WriteString(fmt.Sprintf("Total Commits,%d\n", stats.TotalCommits))
	if stats.DateRange != nil {
		sb.WriteString(fmt.Sprintf("Date Range Start,%s\n", stats.DateRange.Start))
		sb.WriteString(fmt.Sprintf("Date Range End,%s\n", stats.DateRange.End))
	}
	sb.WriteString(fmt.Sprintf("Average Per Day,%.1f\n", stats.AvgPerDay))
	sb.WriteString(fmt.Sprintf("Most Active Day,%s\n", stats.MostActiveDay))
	sb.WriteString(fmt.Sprintf("Most Active Hour,%d\n", stats.MostActiveHour))
	
	// Authors section
	sb.WriteString("\n# Authors\n")
	sb.WriteString("Author,Commits,Percentage\n")
	for _, author := range stats.Authors {
		sb.WriteString(fmt.Sprintf("%s,%d,%s\n", author.Name, author.Commits, author.Percent))
	}
	
	// By Date section
	if len(stats.ByDate) > 0 {
		sb.WriteString("\n# Commits By Date\n")
		sb.WriteString("Date,Commits\n")
		dates := make([]string, 0, len(stats.ByDate))
		for date := range stats.ByDate {
			dates = append(dates, date)
		}
		sort.Strings(dates)
		for _, date := range dates {
			sb.WriteString(fmt.Sprintf("%s,%d\n", date, stats.ByDate[date]))
		}
	}
	
	return sb.String()
}

// RenderMarkdown renders statistics as Markdown format
func (r *StatsRenderer) RenderMarkdown(stats *models.Stats) string {
	var sb strings.Builder
	
	// Title
	sb.WriteString("# Commit Statistics\n\n")
	
	// Summary table
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Total Commits | %d |\n", stats.TotalCommits))
	if stats.DateRange != nil {
		sb.WriteString(fmt.Sprintf("| Date Range | %s to %s |\n", stats.DateRange.Start, stats.DateRange.End))
	}
	sb.WriteString(fmt.Sprintf("| Average Per Day | %.1f |\n", stats.AvgPerDay))
	sb.WriteString(fmt.Sprintf("| Most Active Day | %s |\n", stats.MostActiveDay))
	sb.WriteString(fmt.Sprintf("| Most Active Hour | %02d:00 |\n", stats.MostActiveHour))
	
	// Authors table
	sb.WriteString("\n## Commits by Author\n\n")
	sb.WriteString("| Author | Commits | Percentage |\n")
	sb.WriteString("|--------|---------|------------|\n")
	for _, author := range stats.Authors {
		sb.WriteString(fmt.Sprintf("| %s | %d | %s |\n", author.Name, author.Commits, author.Percent))
	}
	
	// By Date table
	if len(stats.ByDate) > 0 {
		sb.WriteString("\n## Commits by Day\n\n")
		sb.WriteString("| Date | Commits |\n")
		sb.WriteString("|------|---------|\n")
		dates := make([]string, 0, len(stats.ByDate))
		for date := range stats.ByDate {
			dates = append(dates, date)
		}
		sort.Strings(dates)
		for _, date := range dates {
			sb.WriteString(fmt.Sprintf("| %s | %d |\n", date, stats.ByDate[date]))
		}
	}
	
	return sb.String()
}

// RenderHTML renders statistics as HTML format
func (r *StatsRenderer) RenderHTML(stats *models.Stats) string {
	var sb strings.Builder
	
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Commit Statistics</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 40px; max-width: 800px; }
h1, h2 { color: #333; }
table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }
th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
th { background-color: #f4f4f4; }
tr:nth-child(even) { background-color: #fafafa; }
.metric { font-weight: bold; color: #0366d6; }
</style>
</head>
<body>
<h1>Commit Statistics</h1>

<h2>Summary</h2>
<table>
<thead>
<tr><th>Metric</th><th>Value</th></tr>
</thead>
<tbody>
`)
	
	sb.WriteString(fmt.Sprintf("<tr><td class=\"metric\">Total Commits</td><td>%d</td></tr>\n", stats.TotalCommits))
	if stats.DateRange != nil {
		sb.WriteString(fmt.Sprintf("<tr><td class=\"metric\">Date Range</td><td>%s to %s</td></tr>\n", stats.DateRange.Start, stats.DateRange.End))
	}
	sb.WriteString(fmt.Sprintf("<tr><td class=\"metric\">Average Per Day</td><td>%.1f</td></tr>\n", stats.AvgPerDay))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"metric\">Most Active Day</td><td>%s</td></tr>\n", stats.MostActiveDay))
	sb.WriteString(fmt.Sprintf("<tr><td class=\"metric\">Most Active Hour</td><td>%02d:00</td></tr>\n", stats.MostActiveHour))
	
	sb.WriteString(`</tbody>
</table>

<h2>Commits by Author</h2>
<table>
<thead>
<tr><th>Author</th><th>Commits</th><th>Percentage</th></tr>
</thead>
<tbody>
`)
	
	for _, author := range stats.Authors {
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>\n", author.Name, author.Commits, author.Percent))
	}
	
	sb.WriteString(`</tbody>
</table>
</body>
</html>
`)
	
	return sb.String()
}
