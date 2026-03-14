package utils

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#86EFAC")).Bold(true)
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FDA4AF")).Bold(true)
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC")).Bold(true)
	msgStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8E8E8"))
)

func PrintOK(msg string) {
	fmt.Println(okStyle.Render("OK") + "  " + msgStyle.Render(msg))
}

func PrintErr(msg string) {
	fmt.Println(errStyle.Render("ERR") + "  " + msgStyle.Render(msg))
}

func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render("INFO") + "  " + msgStyle.Render(msg))
}

func PrintInvocation(command string, payload map[string]any) error {
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
