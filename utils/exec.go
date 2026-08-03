package utils

import (
	"os"
	"os/exec"
)

func ExecuteCommand(command string, args ...string) error {
	return ExecuteCommandInWorkspace("", command, args...)
}

// ExecuteCommandInWorkspace runs a command with workspace as its working
// directory. An empty workspace preserves the caller's current directory.
func ExecuteCommandInWorkspace(workspace, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if workspace != "" {
		cmd.Dir = workspace
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
