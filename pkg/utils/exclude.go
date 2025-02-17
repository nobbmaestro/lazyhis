package utils

import "regexp"

func IsExcluded(command string, excludeCommands []string) bool {
	for _, pattern := range excludeCommands {
		if matched, _ := regexp.MatchString(pattern, command); matched {
			return true
		}
	}
	return false
}
