package main

import (
	"context"
	"fmt"
	"os"
	"spd/cmd"
)

func main() {
	root := cmd.NewRootCommand()
	if err := root.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
