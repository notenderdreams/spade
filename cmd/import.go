package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"spd/db"
	"spd/utils"

	"github.com/urfave/cli/v3"
)

func newImportCommand() *cli.Command {
	return &cli.Command{
		Name:      "import",
		Aliases:   []string{"im"},
		Usage:     "Import scripts from a JSON file",
		ArgsUsage: "[file]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "overwrite",
				Aliases: []string{"o"},
				Usage:   "Overwrite existing scripts",
			},
		},
		Action: cmdImport,
	}
}

func cmdImport(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()
	if len(args) < 1 {
		utils.PrintErr("usage: spd import <file>")
		return nil
	}

	path := args[0]

	overwrite := c.Bool("overwrite")

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var file db.ExportFile
	if err := json.Unmarshal(data, &file); err != nil {
		return err
	}

	// For now
	if file.Version != exportVersion {
		utils.PrintErr(fmt.Sprintf("unsupported export version %q (expected %q)", file.Version, exportVersion))
		return nil
	}

	imported, skipped := 0, 0
	for _, s := range file.Scripts {
		existing, err := db.GetScript(s.Name)
		if err != nil {
			return err
		}
		if existing != nil {
			if !overwrite {
				utils.PrintInfo(fmt.Sprintf("skipping %q — already exists", s.Name))
				skipped++
				continue
			}
			if err := db.UpdateScript(s); err != nil {
				return err
			}
		} else {
			if err := db.AddScript(s); err != nil {
				return err
			}
		}
		imported++
	}

	utils.PrintOK(fmt.Sprintf("imported %d scripts, skipped %d", imported, skipped))
	return nil
}
