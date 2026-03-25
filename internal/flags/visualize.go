package flags

import "github.com/urfave/cli/v2"

// VisualizeFlags returns flags for the visualize command
func VisualizeFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"n"},
			Usage:   "Limit the number of commits to display",
			Value:   20,
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Show commits from a specific branch",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "author",
			Aliases: []string{"a"},
			Usage:   "Show commits by a specific author",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: ascii, compact, json",
			Value:   "ascii",
		},
	}
}
