package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "Remove a chain",
		ArgsUsage: "<name>",
		Action:    cmdChainRemove,
	}
}

func cmdChainRemove(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd chain remove <name>")
		return nil
	}

	deleted, err := db.DeleteChain(args[0])
	if err != nil {
		return err
	}
	if !deleted {
		utils.PrintInfo(fmt.Sprintf("chain %q not found", args[0]))
		return nil
	}

	utils.PrintOK(fmt.Sprintf("removed chain %q", args[0]))
	return nil
}