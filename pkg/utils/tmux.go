package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetCurrentTmuxSession() (string, error) {
	_, err := exec.LookPath("tmux")
	if err != nil {
		return "", os.ErrNotExist
	}

	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get tmux session: %w", err)
	}

	return strings.TrimSpace(string(stdout)), nil
}
