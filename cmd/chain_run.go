package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainRunCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "Run a chain",
		ArgsUsage: "<name> [args...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"dr"},
				Usage:   "Print resolved commands without executing",
			},
			&cli.BoolFlag{
				Name:    "confirm",
				Aliases: []string{"c"},
				Usage:   "Prompt before each step",
			},
			&cli.BoolFlag{
				Name:    "stop-on-error",
				Aliases: []string{"se"},
				Usage:   "Stop chain if a step fails",
			},
		},
		Action: cmdChainRun,
	}
}

func cmdChainRun(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd chain run <name> [args...]")
		return nil
	}

	return executeChainByName(args[0], args[1:], c.Bool("dry-run"), c.Bool("confirm"), c.Bool("stop-on-error"))
}

func executeChainByName(name string, commandArgs []string, dryRun, confirm, stopOnError bool) error {
	chain, err := db.GetChain(name)
	if err != nil {
		return err
	}
	if chain == nil {
		utils.PrintErr(fmt.Sprintf("no chain found: %q", name))
		return nil
	}

	for _, step := range chain.Steps {
		utils.PrintInfo(fmt.Sprintf("running [%d] %s", step.Seq, step.Script.Name))

		if err := invokeScript(&step.Script, commandArgs, dryRun, confirm, step.Script.Name); err != nil {
			utils.PrintErr(fmt.Sprintf("[%d] %s failed: %s", step.Seq, step.Script.Name, err.Error()))
			if stopOnError {
				return nil
			}
		}
	}
	return nil
}