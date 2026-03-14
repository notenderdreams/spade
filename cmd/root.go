package cmd

import (
	"context"
	"fmt"
	"os"
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
				fmt.Fprintln(os.Stderr, "Error:", err)
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
		return handleScriptInvocation(args[0], args[1:])
	}

	return cli.ShowRootCommandHelp(c)
}

func handleScriptInvocation(name string, commandArgs []string) error {
	script, err := db.GetScript(name)
	if err != nil {
		return err
	}
	if script == nil {
		return utils.PrintInvocation(name, map[string]any{
			"name":         name,
			"command_args": commandArgs,
			"message":      "no command or script found",
		})
	}

	effectiveArgs := append([]string(nil), script.Args...)
	effectiveArgs = append(effectiveArgs, commandArgs...)

	return utils.ExecuteCommand(script.Command, effectiveArgs...)
}
