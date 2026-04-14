package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-repo-visualize/internal/git"
	"github.com/h1s97x/gh-repo-visualize/internal/visualize"
)

// Branch handles the branch command
func Branch(c *cli.Context) error {
	// Get flags
	limit := c.Int("limit")
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
		colorEnabled = isTTY()
	}

	// Create git client
	client := git.NewClient("")

	// Check if we're in a git repo
	if !client.IsGitRepo() {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
	}

	// Get branch graph
	graph, err := client.GetBranchGraph(limit)
	if err != nil {
		return fmt.Errorf("failed to get branch graph: %w", err)
	}

	if len(graph.Branches) == 0 {
		fmt.Println("No branches found")
		return nil
	}

	// Render output
	renderer := visualize.NewBranchGraphRenderer(visualize.RenderOptions{
		Format: format,
		Width:  80,
		Color:  colorEnabled,
	})

	switch format {
	case "json":
		fmt.Println(renderer.RenderJSON(graph))
	case "tree":
		fmt.Println(renderer.RenderASCII(graph))
	default:
		fmt.Println(renderer.Render(graph))
	}

	return nil
}
