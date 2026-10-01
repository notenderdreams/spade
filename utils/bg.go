package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ProcessInfo struct {
	ID        int       `json:"id"`
	PID       int       `json:"pid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	Workspace string    `json:"workspace"`
	LogFile   string    `json:"log_file"`
	StartTime time.Time `json:"start_time"`
}

func getBgDir() string {
	dir := filepath.Join(getSpadeConfigDir(), "bg")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func getBgFile() string {
	return filepath.Join(getBgDir(), "processes.json")
}

func loadBgProcesses() ([]ProcessInfo, error) {
	filePath := getBgFile()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []ProcessInfo{}, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var list []ProcessInfo
	if err := json.Unmarshal(data, &list); err != nil {
		return []ProcessInfo{}, nil
	}
	return list, nil
}

func saveBgProcesses(list []ProcessInfo) error {
	filePath := getBgFile()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func IsProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func ExecuteDetachedInWorkspace(workspace, name, command string, args ...string) (int, string, error) {
	list, err := loadBgProcesses()
	if err != nil {
		list = []ProcessInfo{}
	}

	maxID := 0
	for _, p := range list {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	nextID := maxID + 1

	logFile := filepath.Join(getBgDir(), fmt.Sprintf("bg_%d.log", nextID))
	f, err := os.Create(logFile)
	if err != nil {
		return 0, "", err
	}

	cmd := exec.Command(command, args...)
	if workspace != "" {
		cmd.Dir = workspace
	}
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		f.Close()
		return 0, "", err
	}

	info := ProcessInfo{
		ID:        nextID,
		PID:       cmd.Process.Pid,
		Name:      name,
		Command:   command,
		Args:      args,
		Workspace: workspace,
		LogFile:   logFile,
		StartTime: time.Now(),
	}

	list = append(list, info)
	if err := saveBgProcesses(list); err != nil {
		return 0, "", err
	}

	return nextID, logFile, nil
}

func ListBgProcesses() ([]ProcessInfo, error) {
	list, err := loadBgProcesses()
	if err != nil {
		return nil, err
	}

	active := make([]ProcessInfo, 0, len(list))
	for _, p := range list {
		if IsProcessAlive(p.PID) {
			active = append(active, p)
		}
	}

	if len(active) != len(list) {
		_ = saveBgProcesses(active)
	}

	return active, nil
}

func KillBgProcess(target string) (*ProcessInfo, error) {
	list, err := loadBgProcesses()
	if err != nil {
		return nil, err
	}

	targetID, idErr := strconv.Atoi(target)

	var targetProc *ProcessInfo
	remaining := make([]ProcessInfo, 0, len(list))

	for _, p := range list {
		match := false
		if idErr == nil && p.ID == targetID {
			match = true
		} else if strings.EqualFold(p.Name, target) {
			match = true
		}

		if match && targetProc == nil {
			pCopy := p
			targetProc = &pCopy
		} else {
			remaining = append(remaining, p)
		}
	}

	if targetProc == nil {
		return nil, fmt.Errorf("no running background process found matching %q", target)
	}

	if IsProcessAlive(targetProc.PID) {
		_ = syscall.Kill(-targetProc.PID, syscall.SIGKILL)
		proc, err := os.FindProcess(targetProc.PID)
		if err == nil {
			_ = proc.Kill()
		}
	}

	_ = saveBgProcesses(remaining)
	return targetProc, nil
}

type ClearBgLogsResult struct {
	Cleared int
	Skipped int
}

func ClearBgLogs(target string, force bool) (*ClearBgLogsResult, error) {
	bgDir := getBgDir()
	activeProcs, err := ListBgProcesses()
	if err != nil {
		return nil, err
	}

	activeLogs := make(map[string]bool)
	for _, p := range activeProcs {
		activeLogs[filepath.Clean(p.LogFile)] = true
	}

	res := &ClearBgLogsResult{}

	// If a specific target is provided (id, process name, or log file name)
	if target != "" {
		allProcs, _ := loadBgProcesses()
		targetID, idErr := strconv.Atoi(target)

		// Try finding by process ID or Name in registered processes
		var matchedProc *ProcessInfo
		for _, p := range allProcs {
			if (idErr == nil && p.ID == targetID) || strings.EqualFold(p.Name, target) {
				pCopy := p
				matchedProc = &pCopy
				break
			}
		}

		var targetLogPath string
		if matchedProc != nil {
			targetLogPath = matchedProc.LogFile
			if IsProcessAlive(matchedProc.PID) && !force {
				return nil, fmt.Errorf("process [%d] (%s) is still running; use --force to delete its log", matchedProc.ID, matchedProc.Name)
			}
		} else {
			// Check if target is a file in bgDir or target is an ID like "1" -> "bg_1.log"
			candidates := []string{
				filepath.Join(bgDir, target),
			}
			if idErr == nil {
				candidates = append(candidates, filepath.Join(bgDir, fmt.Sprintf("bg_%d.log", targetID)))
			}
			if !strings.HasSuffix(target, ".log") {
				candidates = append(candidates, filepath.Join(bgDir, target+".log"))
			}

			for _, c := range candidates {
				if info, err := os.Stat(c); err == nil && !info.IsDir() {
					targetLogPath = c
					break
				}
			}
		}

		if targetLogPath == "" {
			return nil, fmt.Errorf("no log file found matching %q", target)
		}

		if activeLogs[filepath.Clean(targetLogPath)] && !force {
			return nil, fmt.Errorf("log file %q belongs to an active process; use --force to delete it", filepath.Base(targetLogPath))
		}

		if err := os.Remove(targetLogPath); err != nil {
			return nil, err
		}
		res.Cleared = 1
		return res, nil
	}

	// No target provided: clear all log files
	entries, err := os.ReadDir(bgDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}

		filePath := filepath.Join(bgDir, entry.Name())
		if activeLogs[filepath.Clean(filePath)] && !force {
			res.Skipped++
			continue
		}

		if err := os.Remove(filePath); err == nil {
			res.Cleared++
		}
	}

	return res, nil
}
