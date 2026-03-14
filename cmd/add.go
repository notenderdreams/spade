package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"
	"strings"

	"github.com/urfave/cli/v3"
)

func newAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new script",
		ArgsUsage: "<name> <command> [args...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "runner",
				Usage: "Optional runner to prepend when executing the script",
			},
		},
		Action: cmdAdd,
	}
}

func cmdAdd(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("usage: spade add <name> <command> [args...]")
	}

	runner := c.String("runner")
	command := args[1]
	commandArgs := args[2:]
	if runner != "" {
		command = runner
		commandArgs = args[1:]
	}

	payload := map[string]any{
		"raw_args":     args,
		"runner":       runner,
		"name":         args[0],
		"command":      command,
		"command_args": commandArgs,
	}

	script := db.Script{
		Name:    args[0],
		Command: command,
		Args:    commandArgs,
	}

	if err := db.AddScript(script); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: scripts.name") {
			return fmt.Errorf("script %q already exists", args[0])
		}
		return err
	}

	payload["saved"] = true
	payload["script"] = script

	return utils.PrintInvocation(c.Name, payload)
}
