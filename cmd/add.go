package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Add a new script",
		ArgsUsage: "<name> <command> [args...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "runner",
				Usage: "Optional runner to prepend when executing the script",
			},
			&cli.BoolFlag{
				Name:    "relative-path",
				Aliases: []string{"rp", "rel"},
				Usage:   "Attach current directory to the command path",
			},
			&cli.BoolFlag{
				Name:    "store",
				Aliases: []string{"s"},
				Usage:   "Store a backup of the script in the script vault",
			},
		},
		Action: cmdAdd,
	}
}

func cmdAdd(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: spade add <name> <command> [args...]")
		return nil
	}

	runner := c.String("runner")
	useRelativePath := c.Bool("relative-path")
	store := c.Bool("store")

	command := args[1]
	commandArgs := append([]string(nil), args[2:]...)

	if useRelativePath {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		command = filepath.Join(cwd, command)
	}

	existing, err := db.GetScript(args[0])
	if err != nil {
		return err
	}
	if existing != nil {
		utils.PrintErr(fmt.Sprintf("script %q already exists", args[0]))
		return nil
	}
	if store {
		if err := utils.StoreAtScriptDir(args[0], command); err != nil {
			utils.PrintErr(fmt.Sprintf("failed to store script %q: %v", args[0], err))
		}
	}

	script := db.Script{
		Name:    args[0],
		Command: command,
		Args:    commandArgs,
		Runner:  runner,
	}

	if err := db.AddScript(script); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("added script %q → %s", args[0], script.Command))
	if store {
		utils.PrintOK(fmt.Sprintf("backup stored → %s", utils.GetScriptPath(args[0], command)))
	}

	return nil
}
