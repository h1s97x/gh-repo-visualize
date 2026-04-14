package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	internal_cmd "github.com/h1s97x/gh-repo-visualize/internal/cmd"
	"github.com/h1s97x/gh-repo-visualize/internal/flags"
)

func main() {
	app := &cli.App{
		Name:    "gh-repo-visualize",
		Usage:   "Visualize Git commit history as ASCII graphs",
		Version: Version,
		Description: `A GitHub CLI extension that generates text-based visualizations
   of Git commit history, including commit graphs and statistics.

Examples:
   # Show commit graph for current repo
   gh repo visualize

   # Show last 50 commits
   gh repo visualize -n 50

   # Show commits from a specific branch
   gh repo visualize -b main

   # Show commits by a specific author
   gh repo visualize -a "John Doe"

   # Output as JSON
   gh repo visualize -f json

   # Show commit statistics
   gh repo visualize stats

   # Show stats by author
   gh repo visualize stats --by-author`,
		// Global flags for default action
		Flags:  flags.VisualizeFlags(),
		Action: internal_cmd.Visualize,
		Commands: []*cli.Command{
			{
				Name:        "stats",
				Usage:       "Display commit statistics",
				Description: "Shows statistics about commits, including breakdown by author and day.",
				Flags:       flags.StatsFlags(),
				Action:      internal_cmd.Stats,
			},
			{
				Name:        "branch",
				Usage:       "Display branch graph",
				Description: "Shows branch visualization including local and remote branches.",
				Flags:       flags.BranchFlags(),
				Action:      internal_cmd.Branch,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
