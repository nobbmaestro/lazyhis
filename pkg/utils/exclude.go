package utils

import (
	"regexp"
	"strings"
)

func MatchesExclusionPatterns(command string, excludeCommands []string) bool {
	for _, pattern := range excludeCommands {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return true
		}
	}
	return false
}

func ContainsExclusionPrefix(command []string, excludePrefix string) bool {
	if len(command) == 0 {
		return true
	}
	return strings.HasPrefix(command[0], excludePrefix)
}

func IsExcludedCommand(
	command []string,
	excludePrefix string,
	excludeCommands []string,
) bool {
	if ContainsExclusionPrefix(command, excludePrefix) {
		return true
	}

	if MatchesExclusionPatterns(strings.Join(command, " "), excludeCommands) {
		return true
	}

	return false
}
