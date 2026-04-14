package flags

import "github.com/urfave/cli/v2"

// BranchFlags returns flags for the branch command
func BranchFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"n"},
			Usage:   "Limit the number of branches to display",
			Value:   20,
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: ascii, tree, json",
			Value:   "ascii",
		},
		&cli.BoolFlag{
			Name:  "color",
			Usage: "Force enable color output",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "Force disable color output (for pipe redirection)",
			Value: false,
		},
	}
}
