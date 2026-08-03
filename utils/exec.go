package utils

import (
	"os"
	"os/exec"
)

type PipelineCommand struct {
	Command   string
	Args      []string
	Workspace string
}

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

func ExecutePipeline(commands []PipelineCommand) error {
	if len(commands) == 0 {
		return nil
	}
	processes := make([]*exec.Cmd, len(commands))
	for i, spec := range commands {
		cmd := exec.Command(spec.Command, spec.Args...)
		if spec.Workspace != "" {
			cmd.Dir = spec.Workspace
		}
		cmd.Stderr = os.Stderr
		processes[i] = cmd
	}
	for i := 0; i < len(processes)-1; i++ {
		pipe, err := processes[i].StdoutPipe()
		if err != nil {
			return err
		}
		processes[i+1].Stdin = pipe
	}
	processes[0].Stdin = os.Stdin
	processes[len(processes)-1].Stdout = os.Stdout
	for _, cmd := range processes {
		if err := cmd.Start(); err != nil {
			return err
		}
	}
	var firstErr error
	for _, cmd := range processes {
		if err := cmd.Wait(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
