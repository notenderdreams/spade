package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainInfoCommand() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Aliases:   []string{"i"},
		Usage:     "Show details of a chain",
		ArgsUsage: "<name>",
		Action:    cmdChainInfo,
	}
}

func cmdChainInfo(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd chain info <name>")
		return nil
	}

	chain, err := db.GetChain(args[0])
	if err != nil {
		return err
	}
	if chain == nil {
		utils.PrintErr(fmt.Sprintf("no chain found: %q", args[0]))
		return nil
	}

	fmt.Println(utils.InfoHeaderStyle.Render("Chain: ") + utils.InfoNameStyle.Render(chain.Name))
	fmt.Println(utils.InfoHeaderStyle.Render("Steps:"))
	for _, step := range chain.Steps {
		prefix := utils.InfoValueStyle.Render(fmt.Sprintf("  [%d]", step.Seq))
		name := utils.InfoNameStyle.Render(step.Script.Name)
		cmd := utils.InfoValueStyle.Render(step.Script.Command)
		args := utils.RenderArgs(step.Script.Args)
		fmt.Printf("%s %s  %s %s\n", prefix, name, cmd, args)
	}
	return nil
}