package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func RunCommand(cmdArgs []string) (string, error) {
	if len(cmdArgs) == 0 {
		return "", fmt.Errorf("command is empty")
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command '%s': %w", cmdArgs, err)
	}

	ret := strings.TrimSpace(string(stdout))
	ret = regexp.MustCompile("'").ReplaceAllString(ret, "")

	return ret, nil
}
