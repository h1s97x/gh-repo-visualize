package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-repo-visualize/internal/git"
	"github.com/h1s97x/gh-repo-visualize/internal/visualize"
)

// Visualize handles the visualize command
func Visualize(c *cli.Context) error {
	// Get flags
	limit := c.Int("limit")
	branch := c.String("branch")
	author := c.String("author")
	format := c.String("format")
	enableColor := c.Bool("color")
	disableColor := c.Bool("no-color")

	// Determine color setting
	// Priority: --no-color > --color > auto-detect TTY
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
		Limit:  limit,
		Branch: branch,
		Author: author,
	}

	commits, err := client.GetCommits(opts)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("No commits found matching the criteria")
		return nil
	}

	// Render output
	renderer := visualize.NewGraphRenderer(visualize.RenderOptions{
		Format: format,
		Width:  80,
		Color:  colorEnabled,
	})

	switch format {
	case "json":
		fmt.Println(renderer.RenderJSON(commits))
	case "compact":
		fmt.Println(renderer.RenderCompact(commits))
	default:
		fmt.Println(renderer.Render(commits))
	}

	return nil
}

// isTTY checks if stdout is a terminal
func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	// Check if it's a character device (terminal)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
