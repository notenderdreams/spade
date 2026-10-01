package cmd

import (
	"context"
	"fmt"
	"spd/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

var (
	bgHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8E8E8")).Bold(true).PaddingRight(2)
	bgIDStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FCD34D")).PaddingRight(2)
	bgPIDStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).PaddingRight(2)
	bgNameStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#C084FC")).PaddingRight(2)
	bgCmdStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC")).PaddingRight(2)
	bgTimeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6EE7B7")).PaddingRight(2)
	bgLogStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true)
)

func newBgCommand() *cli.Command {
	return &cli.Command{
		Name:    "bg",
		Usage:   "Manage background processes",
		Commands: []*cli.Command{
			newBgListCommand(),
			newBgKillCommand(),
			newBgCleanCommand(),
		},
	}
}

func newBgListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List running background processes",
		Action:  cmdBgList,
	}
}

func newBgKillCommand() *cli.Command {
	return &cli.Command{
		Name:      "kill",
		Aliases:   []string{"k"},
		Usage:     "Kill a running background process by ID, PID, or Name",
		ArgsUsage: "<id|pid|name>",
		Action:    cmdBgKill,
	}
}

func cmdBgList(ctx context.Context, c *cli.Command) error {
	_ = ctx
	procs, err := utils.ListBgProcesses()
	if err != nil {
		return err
	}

	if len(procs) == 0 {
		utils.PrintInfo("No background processes running")
		return nil
	}

	idW, pidW, nameW, cmdW := len("ID"), len("PID"), len("NAME"), len("COMMAND")
	for _, p := range procs {
		idStr := fmt.Sprintf("[%d]", p.ID)
		if len(idStr) > idW {
			idW = len(idStr)
		}
		pidStr := fmt.Sprintf("%d", p.PID)
		if len(pidStr) > pidW {
			pidW = len(pidStr)
		}
		if len(p.Name) > nameW {
			nameW = len(p.Name)
		}
		fullCmd := p.Command
		if len(p.Args) > 0 {
			fullCmd += " " + strings.Join(p.Args, " ")
		}
		if len(fullCmd) > cmdW {
			cmdW = len(fullCmd)
		}
	}

	pad := func(s string, w int) string {
		return s + strings.Repeat(" ", max(0, w-len(s)))
	}

	fmt.Println(
		bgHeaderStyle.Render(pad("ID", idW+2)) +
			bgHeaderStyle.Render(pad("PID", pidW+2)) +
			bgHeaderStyle.Render(pad("NAME", nameW+2)) +
			bgHeaderStyle.Render(pad("COMMAND", cmdW+2)) +
			bgHeaderStyle.Render(pad("STARTED", 12)) +
			bgHeaderStyle.Render("LOG FILE"),
	)

	for _, p := range procs {
		idStr := pad(fmt.Sprintf("[%d]", p.ID), idW+2)
		pidStr := pad(fmt.Sprintf("%d", p.PID), pidW+2)
		nameStr := pad(p.Name, nameW+2)
		fullCmd := p.Command
		if len(p.Args) > 0 {
			fullCmd += " " + strings.Join(p.Args, " ")
		}
		cmdStr := pad(fullCmd, cmdW+2)
		timeStr := pad(p.StartTime.Format("15:04:05"), 12)

		fmt.Println(
			bgIDStyle.Render(idStr) +
				bgPIDStyle.Render(pidStr) +
				bgNameStyle.Render(nameStr) +
				bgCmdStyle.Render(cmdStr) +
				bgTimeStyle.Render(timeStr) +
				bgLogStyle.Render(p.LogFile),
		)
	}

	return nil
}

func cmdBgKill(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd bg kill <id|pid|name>")
		return nil
	}

	proc, err := utils.KillBgProcess(args[0])
	if err != nil {
		utils.PrintErr(err.Error())
		return nil
	}

	utils.PrintOK(fmt.Sprintf("Killed background process [%d] (PID %d: %s)", proc.ID, proc.PID, proc.Name))
	return nil
}

func newBgCleanCommand() *cli.Command {
	return &cli.Command{
		Name:      "clean",
		Aliases:   []string{"c"},
		Usage:     "Clear background process log files",
		ArgsUsage: "[id|name]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force removal of logs even for active background processes",
			},
		},
		Action: cmdBgClean,
	}
}

func cmdBgClean(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	target := ""
	if len(args) > 0 {
		target = args[0]
	}

	force := c.Bool("force")
	res, err := utils.ClearBgLogs(target, force)
	if err != nil {
		utils.PrintErr(err.Error())
		return nil
	}

	if res.Cleared == 0 && res.Skipped == 0 {
		utils.PrintInfo("No log files found to clear")
		return nil
	}

	if res.Cleared > 0 {
		utils.PrintOK(fmt.Sprintf("Cleared %d log file(s)", res.Cleared))
	}

	if res.Skipped > 0 {
		utils.PrintInfo(fmt.Sprintf("Skipped %d active process log(s) (use --force or -f to delete all)", res.Skipped))
	}

	return nil
}
