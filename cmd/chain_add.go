package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newChainAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Create a new chain",
		ArgsUsage: "<name> <script1> [-> <script2> ...]",
		Action:    cmdChainAdd,
	}
}

func cmdChainAdd(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: spd chain add <name> <script1> <script2> ...")
		return nil
	}

	name := args[0]
	scriptNames := args[1:]
	pipeAfter := make([]bool, len(scriptNames))
	cleanNames := make([]string, 0, len(scriptNames))
	lastWasPipe := false
	for _, token := range scriptNames {
		if token == "->" {
			if len(cleanNames) == 0 || lastWasPipe {
				utils.PrintErr("pipe operator must connect two scripts")
				return nil
			}
			pipeAfter[len(cleanNames)-1] = true
			lastWasPipe = true
			continue
		}
		cleanNames = append(cleanNames, token)
		lastWasPipe = false
	}
	if len(cleanNames) < 1 || lastWasPipe {
		utils.PrintErr("pipe operator must connect two scripts")
		return nil
	}
	scriptNames = cleanNames

	existing, err := db.GetChain(name)
	if err != nil {
		return err
	}
	if existing != nil {
		utils.PrintErr(fmt.Sprintf("chain %q already exists", name))
		return nil
	}

	for _, s := range scriptNames {
		script, err := db.GetScript(s)
		if err != nil {
			return err
		}
		if script == nil {
			utils.PrintErr(fmt.Sprintf("no script found: %q", s))
			return nil
		}
	}

	if err := db.AddChain(name, scriptNames, pipeAfter); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("created chain %q with %d steps", name, len(scriptNames)))
	return nil
}
