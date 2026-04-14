package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-repo-visualize/internal/git"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
	"github.com/h1s97x/gh-repo-visualize/internal/visualize"
)

// Leaderboard handles the leaderboard command
func Leaderboard(c *cli.Context) error {
	// Get flags
	limit := c.Int("limit")
	format := c.String("format")
	branch := c.String("branch")
	since := c.String("since")
	until := c.String("until")
	enableColor := c.Bool("color")
	disableColor := c.Bool("no-color")

	// Determine color setting
	colorEnabled := false
	if disableColor {
		colorEnabled = false
	} else if enableColor {
		colorEnabled = true
	} else {
		colorEnabled = isTTY()
	}

	// Create git client
	client := git.NewClient("")

	// Check if we're in a git repo
	if !client.IsGitRepo() {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
	}

	// Get commits (no limit for stats)
	opts := &git.LogOptions{
		Limit:  0,
		Branch: branch,
		Since:  since,
		Until:  until,
	}

	commits, err := client.GetCommits(opts)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("No commits found")
		return nil
	}

	// Calculate stats
	stats := models.NewStats()
	stats.Calculate(commits)

	// Render output
	renderer := visualize.NewLeaderboardRenderer(visualize.RenderOptions{
		Format: format,
		Width:  60,
		Color:  colorEnabled,
		Limit:  limit,
	})

	switch format {
	case "json":
		fmt.Println(renderer.RenderJSON(stats))
	case "compact":
		fmt.Println(renderer.RenderCompact(stats))
	default:
		fmt.Println(renderer.Render(stats))
	}

	return nil
}
