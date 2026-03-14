package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8E8E8")).Bold(true).PaddingRight(2)
	nameStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#C084FC")).PaddingRight(2)
	cmdStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC")).PaddingRight(2)
	argsStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true).PaddingRight(2)
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

	if len(scripts) == 0 {
		fmt.Println(argsStyle.Render("No scripts saved yet"))
		return nil
	}

	nameW, cmdW := len("NAME"), len("COMMAND")
	for _, s := range scripts {
		if len(s.Name) > nameW {
			nameW = len(s.Name)
		}
		if len(s.Command) > cmdW {
			cmdW = len(s.Command)
		}
	}

	pad := func(s string, w int) string {
		return s + strings.Repeat(" ", max(0, w-len(s)))
	}

	fmt.Println(headerStyle.Render(pad("NAME", nameW+2)) + headerStyle.Render(pad("COMMAND", cmdW+2)) + headerStyle.Render("ARGS"))

	for _, s := range scripts {
		fmt.Println(nameStyle.Render(pad(s.Name, nameW+2)) + cmdStyle.Render(pad(s.Command, cmdW+2)) + argsStyle.Render(strings.Join(s.Args, " ")))
	}

	return nil
}
