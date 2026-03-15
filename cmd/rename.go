package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newRenameCommand() *cli.Command {
	return &cli.Command{
		Name: "rename",
		Aliases: []string {"rnm"} ,
		Usage: "Rename a saved scirpt",
		ArgsUsage: "<old_name> <new_name>",
		Action: cmdRename,	
	}
}

func cmdRename(ctx context.Context, c *cli.Command) error {
	_ = ctx 
	
	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: add rename <old_name> <new_nmae>")
		return nil
	}

	oldName, newName := args[0], args[1] 

	script, err := db.GetScript(oldName)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no script found: %q", oldName))
		return nil
	}

	existing, err := db.GetScript(newName)
	if err != nil {
		return err
	}
	if existing != nil {
		utils.PrintErr(fmt.Sprintf("script %q already exists", newName))
		return nil
	}

	if err := db.RenameScript(oldName, newName); err != nil {
		return err
	}
	utils.PrintOK(fmt.Sprintf("renamed %q → %q", oldName, newName))
	return nil
}