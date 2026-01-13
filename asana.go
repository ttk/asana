package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ttk/asana/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "asana"
	app.Version = "0.2.1"
	app.Usage = "asana cui client ( https://github.com/thash/asana )"

	app.Commands = defs()
	app.Run(os.Args)
}

func defs() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Asana configuration. Your settings will be saved in ~/.asana.yml",
			Action: func(c *cli.Context) error {
				commands.Config(c)
				return nil
			},
		},
		{
			Name:    "workspaces",
			Aliases: []string{"w"},
			Usage:   "get workspaces",
			Action: func(c *cli.Context) error {
				commands.Workspaces(c)
				return nil
			},
		},
		{
			Name:    "tasks",
			Aliases: []string{"ts"},
			Usage:   "get tasks",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "no-cache, n", Usage: "without cache"},
				&cli.BoolFlag{Name: "refresh, r", Usage: "update cache"},
			},
			Action: func(c *cli.Context) error {
				commands.Tasks(c)
				return nil
			},
		},
		{
			Name:      "task",
			Aliases:   []string{"t"},
			Usage:     "get a task",
			ArgsUsage: "<task-index-or-id>",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "verbose output"},
				&cli.BoolFlag{Name: "with-comments", Aliases: []string{"c"}, Usage: "include comments"},
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Usage: "output as JSON"},
			},
			Action: func(c *cli.Context) error {
				commands.Task(c)
				return nil
			},
		},
		{
			Name:      "comment",
			Aliases:   []string{"cm"},
			Usage:     "Post comment",
			ArgsUsage: "<task-index-or-id>",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "comment", Aliases: []string{"c"}, Usage: "comment text (if not provided, opens editor)"},
			},
			Action: func(c *cli.Context) error {
				commands.Comment(c)
				return nil
			},
		},
		{
			Name:      "done",
			Usage:     "Complete task",
			ArgsUsage: "<task-index-or-id>",
			Action: func(c *cli.Context) error {
				commands.Done(c)
				return nil
			},
		},
		{
			Name:      "due",
			Usage:     "set due date",
			ArgsUsage: "<task-index-or-id> <date>",
			Action: func(c *cli.Context) error {
				commands.DueOn(c)
				return nil
			},
		},
		{
			Name:      "set-sec",
			Aliases:   []string{"ss"},
			Usage:     "set task section",
			ArgsUsage: "<task-index-or-id> <section>",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "project", Aliases: []string{"p"}, Usage: "project keyword to match"},
			},
			Action: func(c *cli.Context) error {
				commands.SetSection(c)
				return nil
			},
		},
		{
			Name:      "browse",
			Aliases:   []string{"b"},
			Usage:     "open a task in the web browser",
			ArgsUsage: "<task-index-or-id>",
			Action: func(c *cli.Context) error {
				commands.Browse(c)
				return nil
			},
		},
		{
			Name:      "download",
			Aliases:   []string{"dl"},
			Usage:     "download attachment from a task",
			ArgsUsage: "<task-index-or-id> <attachment-index> or <attachment-index>",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "output file path"},
			},
			Action: func(c *cli.Context) error {
				commands.Download(c)
				return nil
			},
		},
	}
}
