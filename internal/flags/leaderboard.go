package flags

import "github.com/urfave/cli/v2"

// LeaderboardFlags returns flags for the leaderboard command
func LeaderboardFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"n"},
			Usage:   "Number of top contributors to display",
			Value:   10,
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: ascii, compact, json",
			Value:   "ascii",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Show leaderboard for a specific branch",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "since",
			Usage:   "Show commits since this date",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "until",
			Usage:   "Show commits until this date",
			Value:   "",
		},
		&cli.BoolFlag{
			Name:  "color",
			Usage: "Force enable color output",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "Force disable color output",
			Value: false,
		},
	}
}
