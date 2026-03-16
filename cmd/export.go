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

const exportVersion = "1"

func newExportCommand() *cli.Command {
	return &cli.Command{
		Name:      "export",
		Aliases:   []string{"ex"},
		Usage:     "Export all scripts to a JSON file",
		ArgsUsage: "[file]",
		Action:    cmdExport,
	}
}

func cmdExport(ctx context.Context, c *cli.Command) error {
	_ = ctx

	args := c.Args().Slice()

	path := "scripts.spade"
	var scripts []db.Script
	var err error

	// support: export, export <file>, export <name> <file>
	switch len(args) {
	case 0:
		scripts, err = db.GetAllScripts()
	case 1:
		// try as script name first
		s, serr := db.GetScript(args[0])
		if serr != nil {
			return serr
		}
		if s != nil {
			scripts = []db.Script{*s}
			path = fmt.Sprint(s.Name, ".spade")
		} else {
			// treat as path
			path = args[0]
			scripts, err = db.GetAllScripts()
		}
	default:
		// arg[0] is name, arg[1] is path
		path = args[1]
		s, serr := db.GetScript(args[0])
		if serr != nil {
			return serr
		}
		if s == nil {
			utils.PrintErr(fmt.Sprintf("no script found: %q", args[0]))
			return nil
		}
		scripts = []db.Script{*s}
	}

	if err != nil {
		return err
	}

	out := db.ExportFile{
		Version: exportVersion,
		Scripts: scripts,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	utils.PrintOK(fmt.Sprintf("exported %d script(s) → %s", len(scripts), path))
	return nil
}
