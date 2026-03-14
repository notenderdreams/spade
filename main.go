package main

import (
	"context"
	"os"

	"spd/cmd"
)

func main() {
	root := cmd.NewRootCommand()
	if err := root.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
