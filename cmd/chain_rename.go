package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainRenameCommand() *cli.Command {
	return &cli.Command{
		Name:      "rename",
		Aliases:   []string{"rnm"},
		Usage:     "Rename a chain",
		ArgsUsage: "<old_name> <new_name>",
		Action:    cmdChainRename,
	}
}

func cmdChainRename(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: spd chain rename <old_name> <new_name>")
		return nil
	}

	oldName, newName := args[0], args[1]

	chain, err := db.GetChain(oldName)
	if err != nil {
		return err
	}
	if chain == nil {
		utils.PrintErr(fmt.Sprintf("no chain found: %q", oldName))
		return nil
	}

	existing, err := db.GetChain(newName)
	if err != nil {
		return err
	}
	if existing != nil {
		utils.PrintErr(fmt.Sprintf("chain %q already exists", newName))
		return nil
	}

	if err := db.RenameChain(oldName, newName); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("renamed %q → %q", oldName, newName))
	return nil
}