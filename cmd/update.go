package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Aliases:   []string{"u"},
		Usage:     "Update an existing script",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "New command to run",
			},
			&cli.StringFlag{
				Name:    "args",
				Aliases: []string{"a"},
				Usage:   "New args (space separated)",
			},
		},
		Action: cmdUpdate,
	}
}

func cmdUpdate(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd update <name> [--command <cmd>] [--args <args>]")
		return nil
	}

	newCmd := c.String("command")
	newArgs := c.String("args")

	if newCmd == "" && newArgs == "" {
		utils.PrintErr("at least one of --command or --args is required")
		return nil
	}

	name := args[0]
	script, err := db.GetScript(name)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no script found: %q", name))
		return nil
	}

	if newCmd != "" {
		script.Command = newCmd
	}
	if newArgs != "" {
		script.Args = splitArgs(newArgs)
	}

	if err := db.UpdateScript(*script); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("updated script %q", name))
	return nil
}

func splitArgs(s string) []string {
	if s == "" {
		return nil
	}
	parts := []string{}
	current := ""
	inQuote := false
	quoteChar := rune(0)

	for _, ch := range s {
		switch {
		case inQuote && ch == quoteChar:
			inQuote = false
		case !inQuote && (ch == '"' || ch == '\''):
			inQuote = true
			quoteChar = ch
		case !inQuote && ch == ' ':
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
