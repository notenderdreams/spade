package cmd

import (
	"context"
	"os"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainCommand() *cli.Command {
	return &cli.Command{
		Name:    "chain",
		Aliases: []string{"ch"},
		Usage:   "Manage and run script chains",
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
		Action: func(ctx context.Context, c *cli.Command) error {
			_ = ctx
			args := c.Args().Slice()
			if len(args) > 0 {
				return executeChainByName(args[0], args[1:], c.Bool("dry-run"), c.Bool("confirm"), c.Bool("stop-on-error"))
			}
			return cli.ShowSubcommandHelp(c)
		},
		CommandNotFound: func(ctx context.Context, c *cli.Command, name string) {
			_ = ctx
			if err := executeChainByName(name, c.Args().Slice(), c.Bool("dry-run"), c.Bool("confirm"), c.Bool("stop-on-error")); err != nil {
				utils.PrintErr(err.Error())
				os.Exit(1)
			}
		},
		Commands: []*cli.Command{
			newChainAddCommand(),
			newChainRemoveCommand(),
			newChainListCommand(),
			newChainInfoCommand(),
			newChainRunCommand(),
			newChainRenameCommand(),
			newChainAppendCommand(),
		},
	}
}
