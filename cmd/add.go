package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
			&cli.BoolFlag{
				Name:    "relative-path",
				Aliases: []string{"rp", "reletive-path"},
				Usage:   "Attach current directory to the command path",
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
	rawArgs := append([]string(nil), args...)

	runner := c.String("runner")
	useRelativePath := c.Bool("relative-path")
	command := args[1]
	commandArgs := append([]string(nil), args[2:]...)
	if runner != "" {
		command = runner
		commandArgs = append([]string(nil), args[1:]...)
	}

	if useRelativePath {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		command = filepath.Join(cwd, command)
	}

	payload := map[string]any{
		"raw_args":      rawArgs,
		"runner":        runner,
		"relative_path": useRelativePath,
		"name":          args[0],
		"command":       command,
		"command_args":  commandArgs,
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
