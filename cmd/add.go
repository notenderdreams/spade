package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func newAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new script",
		ArgsUsage: "<name> <command> [args...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "runner",
				Usage: "Optional runner to prepend when executing the script",
			},
		},
		Action: cmdAdd,
	}
}

func cmdAdd(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	payload := map[string]any{
		"raw_args": args,
		"runner":   c.String("runner"),
	}

	if len(args) > 0 {
		payload["name"] = args[0]
	}
	if len(args) > 1 {
		payload["command"] = args[1]
	}
	if len(args) > 2 {
		payload["command_args"] = args[2:]
	}

	return printInvocation(c.Name, payload)
}
