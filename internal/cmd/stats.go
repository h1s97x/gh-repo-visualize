package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-repo-visualize/internal/git"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
	"github.com/h1s97x/gh-repo-visualize/internal/visualize"
)

// Stats handles the stats command
func Stats(c *cli.Context) error {
	// Get flags
	byAuthor := c.Bool("by-author")
	byDay := c.Bool("by-day")
	branch := c.String("branch")
	author := c.String("author")
	since := c.String("since")
	until := c.String("until")
	format := c.String("format")
	enableColor := c.Bool("color")
	disableColor := c.Bool("no-color")

	// Determine color setting
	colorEnabled := false
	if disableColor {
		colorEnabled = false
	} else if enableColor {
		colorEnabled = true
	} else {
		// Auto-detect: enable colors if stdout is a TTY
		colorEnabled = isTTY()
	}

	// Create git client
	client := git.NewClient("")
	
	// Check if we're in a git repo
	if !client.IsGitRepo() {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
	}

	// Get commits
	opts := &git.LogOptions{
		Limit:  0, // No limit for stats
		Branch: branch,
		Author: author,
		Since:  since,
		Until:  until,
	}

	commits, err := client.GetCommits(opts)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("No commits found matching the criteria")
		return nil
	}

	// Calculate stats
	stats := models.NewStats()
	stats.Calculate(commits)

	// Render output
	renderer := visualize.NewStatsRendererWithColor(colorEnabled)

	if format == "json" {
		fmt.Println(renderer.RenderJSON(stats))
		return nil
	}

	// Default summary
	fmt.Println(renderer.Render(stats))

	// By author breakdown
	if byAuthor || (!byAuthor && !byDay) {
		fmt.Println(renderer.RenderByAuthor(stats))
	}

	// By day breakdown
	if byDay || (!byAuthor && !byDay) {
		fmt.Println(renderer.RenderByDay(stats))
	}

	return nil
}
