package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"spd/db"
	"spd/utils"

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
			if err := handleScriptInvocation(name, c.Args().Slice()); err != nil {
				utils.PrintErr(err.Error())
				os.Exit(1)
			}
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
		if err := handleScriptInvocation(args[0], args[1:]); err != nil {
			utils.PrintErr(err.Error())
			os.Exit(1)
		}
		return nil
	}

	return cli.ShowRootCommandHelp(c)
}

func handleScriptInvocation(name string, commandArgs []string) error {
	script, err := db.GetScript(name)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no command or script found: %q", name))
		return nil
	}

	cmd, expandedArgs, err := utils.SubstitutePlaceholders(script.Command, script.Args, commandArgs)

	if err != nil {
		return err
	}

	if err := utils.ExecuteCommand(cmd, expandedArgs...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
	return nil
}
