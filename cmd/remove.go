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
		Aliases:   []string{"rm"},
		Usage:     "Remove a saved script",
		ArgsUsage: "<name>",
		Action:    cmdRemove,
	}
}

func cmdRemove(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd remove <name>")
		return nil
	}

	deleted, err := db.DeleteScript(args[0])
	if err != nil {
		return err
	}

	if !deleted {
		utils.PrintInfo(fmt.Sprintf("script %q not found", args[0]))
		return nil
	}

	utils.PrintOK(fmt.Sprintf("removed script %q", args[0]))
	return nil

}
