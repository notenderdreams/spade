package main

import (
	"context"
	"os"
	"spd/cmd"
	"spd/utils"
)

func main() {
	root := cmd.NewRootCommand()
	root.EnableShellCompletion = true
	if err := root.Run(context.Background(), os.Args); err != nil {
		utils.PrintErr(err.Error())
		os.Exit(1)
	}
}
