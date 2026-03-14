package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func newRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove a saved script",
		ArgsUsage: "<name>",
		Action:    cmdRemove,
	}
}

func cmdRemove(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	payload := map[string]any{
		"raw_args": args,
	}

	if len(args) > 0 {
		payload["name"] = args[0]
	}

	return printInvocation(c.Name, payload)
}
