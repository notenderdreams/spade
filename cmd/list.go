package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func newListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List saved scripts",
		Action: cmdList,
	}
}

func cmdList(ctx context.Context, c *cli.Command) error {
	_ = ctx

	return printInvocation(c.Name, map[string]any{
		"raw_args": c.Args().Slice(),
	})
}
