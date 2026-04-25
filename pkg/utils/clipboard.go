package utils

import (
	"os"

	"golang.design/x/clipboard"
)

var clipboardAvailable bool

func init() {
	if os.Getenv("CI") != "" {
		clipboardAvailable = false
		return
	}

	err := clipboard.Init()
	if err != nil {
		clipboardAvailable = false
		return
	}

	clipboardAvailable = true
}

func CopyToClipboard(command string) {
	if !clipboardAvailable {
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(command))
}
