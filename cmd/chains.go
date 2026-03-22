package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"
	"strings"

	"github.com/urfave/cli/v3"
)

func newChainCommand() *cli.Command {
	return &cli.Command{
		Name:    "chain",
		Aliases: []string{"ch"},
		Usage:   "Manage and run script chains",
		Commands: []*cli.Command{
			newChainAddCommand(),
			newChainRemoveCommand(),
			newChainListCommand(),
			newChainInfoCommand(),
			newChainRunCommand(),
		},
	}
}

func newChainAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Create a new chain",
		ArgsUsage: "<name> <script1> <script2> ...",
		Action:    cmdChainAdd,
	}
}

func newChainRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "Remove a chain",
		ArgsUsage: "<name>",
		Action:    cmdChainRemove,
	}
}

func newChainListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all chains",
		Action:  cmdChainList,
	}
}

func newChainInfoCommand() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Aliases:   []string{"i"},
		Usage:     "Show details of a chain",
		ArgsUsage: "<name>",
		Action:    cmdChainInfo,
	}
}

func newChainRunCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "Run a chain",
		ArgsUsage: "<name>",
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
		Action: cmdChainRun,
	}
}

// --- actions ---

func cmdChainAdd(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 2 {
		utils.PrintErr("usage: spd chain add <name> <script1> <script2> ...")
		return nil
	}

	name := args[0]
	scriptNames := args[1:]

	existing, err := db.GetChain(name)
	if err != nil {
		return err
	}
	if existing != nil {
		utils.PrintErr(fmt.Sprintf("chain %q already exists", name))
		return nil
	}

	// verify all scripts exist before inserting
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

	if err := db.AddChain(name, scriptNames); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("created chain %q with %d steps", name, len(scriptNames)))
	return nil
}

func cmdChainRemove(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd chain remove <name>")
		return nil
	}

	deleted, err := db.DeleteChain(args[0])
	if err != nil {
		return err
	}
	if !deleted {
		utils.PrintInfo(fmt.Sprintf("chain %q not found", args[0]))
		return nil
	}

	utils.PrintOK(fmt.Sprintf("removed chain %q", args[0]))
	return nil
}

func cmdChainList(ctx context.Context, c *cli.Command) error {
	_ = ctx

	chains, err := db.GetAllChains()
	if err != nil {
		return err
	}
	if len(chains) == 0 {
		fmt.Println(utils.InfoValueStyle.Render("no chains saved yet"))
		return nil
	}

	nameW := len("NAME")
	for _, ch := range chains {
		if len(ch.Name) > nameW {
			nameW = len(ch.Name)
		}
	}

	pad := func(s string, w int) string {
		return s + strings.Repeat(" ", max(0, w-len(s)))
	}

	fmt.Println(
		utils.InfoHeaderStyle.Render(pad("NAME", nameW+2)) +
			utils.InfoHeaderStyle.Render("STEPS"),
	)

	for _, ch := range chains {
		names := make([]string, len(ch.Steps))
		for i, st := range ch.Steps {
			names[i] = st.Script.Name
		}
		fmt.Println(
			utils.InfoNameStyle.Render(pad(ch.Name, nameW+2)) +
				utils.InfoValueStyle.Render(strings.Join(names, " → ")),
		)
	}
	return nil
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

func cmdChainRun(ctx context.Context, c *cli.Command) error {
	_ = ctx
	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd chain run <name>")
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

	dryRun := c.Bool("dry-run")
	confirm := c.Bool("confirm")
	stopOnError := c.Bool("stop-on-error")

	for _, step := range chain.Steps {
		utils.PrintInfo(fmt.Sprintf("running [%d] %s", step.Seq, step.Script.Name))

		err := invokeScript(&step.Script, []string{}, dryRun, confirm, step.Script.Name)
		if err != nil {
			utils.PrintErr(fmt.Sprintf("[%d] %s failed: %s", step.Seq, step.Script.Name, err.Error()))
			if stopOnError {
				return nil
			}
		}
	}
	return nil
}