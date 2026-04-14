package flags

import "github.com/urfave/cli/v2"

// StatsFlags returns flags for the stats command
func StatsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "by-author",
			Usage: "Show breakdown by author",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "by-day",
			Usage: "Show breakdown by day",
			Value: false,
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Show stats for a specific branch",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "author",
			Aliases: []string{"a"},
			Usage:   "Show stats for a specific author",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "since",
			Usage:   "Show stats since this date (ISO 8601 or relative like '2 weeks ago')",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "until",
			Usage:   "Show stats until this date (ISO 8601 or relative like 'yesterday')",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: ascii, json",
			Value:   "ascii",
		},
	}
}
