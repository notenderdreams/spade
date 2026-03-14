package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

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
	if len(args) < 1 {
		return fmt.Errorf("usage: spade remove <name>")
	}

	deleted, err := db.DeleteScript(args[0])
	if err != nil {
		return err
	}

	payload := map[string]any{
		"raw_args": args,
		"name":     args[0],
		"deleted":  deleted,
	}
	if !deleted {
		payload["message"] = fmt.Sprintf("script %q not found", args[0])
	}

	return utils.PrintInvocation(c.Name, payload)
}
