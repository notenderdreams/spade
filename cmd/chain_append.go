package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainAppendCommand() *cli.Command {
	return &cli.Command{
		Name:      "append",
		Aliases:   []string{"ap"},
		Usage:     "Append a script to an existing chain",
		ArgsUsage: "<chain> <script>",
		Action:    cmdChainAppend,
	}
}

func cmdChainAppend(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: spd chain append <chain> <script>")
		return nil
	}

	chainName, scriptName := args[0], args[1]

	chain, err := db.GetChain(chainName)
	if err != nil {
		return err
	}
	if chain == nil {
		utils.PrintErr(fmt.Sprintf("no chain found: %q", chainName))
		return nil
	}

	script, err := db.GetScript(scriptName)
	if err != nil {
		return err
	}
	if script == nil {
		utils.PrintErr(fmt.Sprintf("no script found: %q", scriptName))
		return nil
	}

	if err := db.AppendChainStep(chainName, scriptName); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("appended %q to chain %q", scriptName, chainName))
	return nil
}