package cmd

import (
	"context"
	"fmt"
	"os"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

// Version is overridden by the release build's linker flags.
var Version = "0.1.1"

func NewRootCommand() *cli.Command {
	return &cli.Command{
		Name:    "spade",
		Usage:   "A simple CLI tool to manage and run scripts",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "workspace",
				Aliases: []string{"w"},
				Usage:   "Run the script from this directory",
			},
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
			&cli.BoolFlag{
				Name:    "detached",
				Aliases: []string{"d"},
				Usage:   "Run process in background (detached)",
			},
		},
		Action: cmdRoot,
		CommandNotFound: func(ctx context.Context, c *cli.Command, name string) {
			_ = ctx
			if err := handleScriptInvocation(name, c.Args().Slice(), false, false, c.Bool("detached"), c.String("workspace")); err != nil {
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
			newBgCommand(),
		},
	}
}

func cmdRoot(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) > 0 {
		if err := handleScriptInvocation(args[0], args[1:], c.Bool("dry-run"), c.Bool("confirm"), c.Bool("detached"), c.String("workspace")); err != nil {
			utils.PrintErr(err.Error())
			os.Exit(1)
		}
		return nil
	}

	return cli.ShowRootCommandHelp(c)
}

func handleScriptInvocation(name string, commandArgs []string, dryRun, confirm, detached bool, workspace string) error {
	script, err := db.GetScript(name)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no command or script found: %q", name))
		return nil
	}

	if err := invokeScript(script, commandArgs, dryRun, confirm, detached, name, workspace); err != nil {
		return err
	}
	return nil
}

func invokeScript(s *db.Script, commandArgs []string, dryRun, confirm, detached bool, stepLabel, workspace string) error {
	cmd, expandedArgs, workspace, err := prepareScript(s, commandArgs, workspace)
	if err != nil {
		return err
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

	if detached {
		id, logFile, err := utils.ExecuteDetachedInWorkspace(workspace, stepLabel, cmd, expandedArgs...)
		if err != nil {
			return err
		}
		utils.PrintOK(fmt.Sprintf("Started background process [%d] (Log: %s)", id, logFile))
		return nil
	}

	return utils.ExecuteCommandInWorkspace(workspace, cmd, expandedArgs...)
}

func prepareScript(s *db.Script, commandArgs []string, workspace string) (string, []string, string, error) {
	if workspace == "" {
		workspace = s.Workspace
	}
	cmd, expandedArgs, err := utils.SubstitutePlaceholders(s.Command, s.Args, commandArgs)
	if err != nil {
		return "", nil, "", err
	}

	cmd, err = utils.ScriptExistsOrBackup(s.Name, cmd)
	if err != nil {
		return "", nil, "", err
	}

	if s.Runner != "" {
		expandedArgs = append([]string{cmd}, expandedArgs...)
		cmd = s.Runner
	}

	return cmd, expandedArgs, workspace, nil
}
