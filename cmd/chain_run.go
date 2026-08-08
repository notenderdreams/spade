package cmd

import (
	"context"
	"fmt"
	"os"
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
			&cli.StringFlag{
				Name:    "workspace",
				Aliases: []string{"w"},
				Usage:   "Run all steps from this directory",
			},
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
			&cli.BoolFlag{
				Name:    "detached",
				Aliases: []string{"d"},
				Usage:   "Run chain in background (detached)",
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

	return executeChainByName(args[0], args[1:], c.Bool("dry-run"), c.Bool("confirm"), c.Bool("stop-on-error"), c.Bool("detached"), c.String("workspace"))
}

func executeChainByName(name string, commandArgs []string, dryRun, confirm, stopOnError, detached bool, workspace string) error {
	chain, err := db.GetChain(name)
	if err != nil {
		return err
	}
	if chain == nil {
		utils.PrintErr(fmt.Sprintf("no chain found: %q", name))
		return nil
	}

	if detached {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		chainArgs := []string{"chain", "run", name}
		chainArgs = append(chainArgs, commandArgs...)
		if workspace != "" {
			chainArgs = append(chainArgs, "-w", workspace)
		}
		if stopOnError {
			chainArgs = append(chainArgs, "-se")
		}

		id, logFile, err := utils.ExecuteDetachedInWorkspace(workspace, name, exePath, chainArgs...)
		if err != nil {
			return err
		}
		utils.PrintOK(fmt.Sprintf("Started background chain [%d] (Log: %s)", id, logFile))
		return nil
	}

	for i := 0; i < len(chain.Steps); {
		step := chain.Steps[i]
		utils.PrintInfo(fmt.Sprintf("running [%d] %s", step.Seq, step.Script.Name))

		if step.PipeToNext {
			pipeline := []utils.PipelineCommand{}
			pipelineStopOnError := false
			j := i
			for {
				current := chain.Steps[j]
				pipelineStopOnError = pipelineStopOnError || current.StopOnError
				cmd, args, dir, err := prepareScript(&current.Script, commandArgs, workspace)
				if err != nil {
					utils.PrintErr(fmt.Sprintf("[%d] %s failed: %s", current.Seq, current.Script.Name, err.Error()))
					if stopOnError || pipelineStopOnError {
						return nil
					}
					break
				}
				if dryRun || confirm {
					utils.PrintDryRun(cmd, args)
				}
				if confirm && !utils.Confirm(fmt.Sprintf("run %q?", current.Script.Name)) {
					continue
				}
				pipeline = append(pipeline, utils.PipelineCommand{Command: cmd, Args: args, Workspace: dir})
				if !current.PipeToNext || j+1 >= len(chain.Steps) {
					break
				}
				j++
			}
			if !dryRun && len(pipeline) > 0 {
				if err := utils.ExecutePipeline(pipeline); err != nil {
					utils.PrintErr(fmt.Sprintf("pipeline failed: %s", err.Error()))
					if stopOnError || pipelineStopOnError {
						return nil
					}
				}
			}
			i = j + 1
			continue
		}

		if err := invokeScript(&step.Script, commandArgs, dryRun, confirm, false, step.Script.Name, workspace); err != nil {
			utils.PrintErr(fmt.Sprintf("[%d] %s failed: %s", step.Seq, step.Script.Name, err.Error()))
			if stopOnError || step.StopOnError {
				return nil
			}
		}
		i++
	}
	return nil
}
