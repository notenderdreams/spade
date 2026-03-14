package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func NewRootCommand() *cli.Command {
	return &cli.Command{
		Name:    "spade",
		Usage:   "A simple CLI tool to manage and run scripts",
		Version: "0.1.0",
		Action:  cmdRoot,
		CommandNotFound: func(ctx context.Context, c *cli.Command, name string) {
			_ = ctx

			args := c.Args().Slice()
			payload := map[string]any{
				"raw_args": args,
			}

			if len(args) > 1 {
				payload["command_args"] = args[1:]
			}

			_ = printInvocation(name, payload)
		},
		Commands: []*cli.Command{
			newAddCommand(),
			newListCommand(),
			newRemoveCommand(),
		},
	}
}

func cmdRoot(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) > 0 {
		payload := map[string]any{
			"raw_args": args,
		}

		if len(args) > 1 {
			payload["command_args"] = args[1:]
		}

		return printInvocation(args[0], payload)
	}

	return cli.ShowRootCommandHelp(c)
}
