package cmd

import (
	"encoding/json"
	"fmt"
)

func printInvocation(command string, payload map[string]any) error {
	body := map[string]any{
		"command": command,
		"args":    payload,
	}

	encoded, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(encoded))
	return nil
}
