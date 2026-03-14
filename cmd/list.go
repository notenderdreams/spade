package cmd

import (
	"context"
	"spd/db"
	"spd/utils"
	"github.com/urfave/cli/v3"
)

func newListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List saved scripts",
		Action: cmdList,
	}
}

func cmdList(ctx context.Context, c *cli.Command) error {
	_ = ctx

	scripts, err := db.GetAllScripts()
	if err != nil {
		return err
	}

	return utils.PrintInvocation(c.Name, map[string]any{
		"raw_args": c.Args().Slice(),
		"scripts":  scripts,
	})
}
