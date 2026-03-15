package cmd

import (
	"context"
	"fmt"
	"strings"

	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newInfoCommand() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Show details of a saved script",
		ArgsUsage: "<name>",
		Action:    cmdInfo,
	}
}

func cmdInfo(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd info <name>")
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

	renderedArgs := make([]string, len(script.Args))
	for i, a := range script.Args {
		renderedArgs[i] = utils.RenderArg(a)
	}

	fmt.Println(utils.InfoHeaderStyle.Render("Script: ") + utils.InfoNameStyle.Render(script.Name))
	fmt.Println(utils.InfoHeaderStyle.Render("Command:") + " " + utils.InfoValueStyle.Render(script.Command))
	if len(script.Args) > 0 {
		fmt.Println(utils.InfoHeaderStyle.Render("Args:   ") + " " + strings.Join(renderedArgs, " "))
	}
	return nil
}
