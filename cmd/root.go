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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"dr"},
				Usage:   "Print the resolved command without executing it",
			},
			&cli.BoolFlag{
				Name:    "confirm",
				Aliases: []string{"c"},
				Usage:   "Show resolved command and prompt before executing",
			},
		},
		Action: cmdRoot,
		CommandNotFound: func(ctx context.Context, c *cli.Command, name string) {
			_ = ctx
			if err := handleScriptInvocation(name, c.Args().Slice(), false, false); err != nil {
				utils.PrintErr(err.Error())
				os.Exit(1)
			}
		},
		Commands: []*cli.Command{
			newAddCommand(),
			newInfoCommand(),
			newListCommand(),
			newRenameCommand(),
			newUpdateCommand(),
			newRemoveCommand(),
			newExportCommand(),
			newImportCommand(),
			newChainCommand(),
		},
	}
}

func cmdRoot(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) > 0 {
		if err := handleScriptInvocation(args[0], args[1:], c.Bool("dry-run"), c.Bool("confirm")); err != nil {
			utils.PrintErr(err.Error())
			os.Exit(1)
		}
		return nil
	}

	return cli.ShowRootCommandHelp(c)
}

func handleScriptInvocation(name string, commandArgs []string, dryRun, confirm bool) error {
	script, err := db.GetScript(name)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no command or script found: %q", name))
		return nil
	}

	if err := invokeScript(script, commandArgs, dryRun, confirm, name); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
	return nil
}

func invokeScript(s *db.Script, commandArgs []string, dryRun, confirm bool, stepLabel string) error {
	cmd, expandedArgs, err := utils.SubstitutePlaceholders(s.Command, s.Args, commandArgs)
	if err != nil {
		return err
	}

	if s.Runner != "" {
		expandedArgs = append([]string{cmd}, expandedArgs...)
		cmd = s.Runner
	}

	if dryRun || confirm {
		utils.PrintDryRun(cmd, expandedArgs)
	}
	if dryRun {
		return nil
	}
	if confirm && !utils.Confirm(fmt.Sprintf("run %q?", stepLabel)) {
		return nil
	}

	return utils.ExecuteCommand(cmd, expandedArgs...)
}
