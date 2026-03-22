package cmd

import (
	"context"
	"fmt"
	"spd/db"
	"spd/utils"
	"strings"

	"github.com/urfave/cli/v3"
)

func newChainListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all chains",
		Action:  cmdChainList,
	}
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