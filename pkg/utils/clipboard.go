package utils

import "golang.design/x/clipboard"

func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func CopyToClipboard(command string) {
	clipboard.Write(clipboard.FmtText, []byte(command))
}
