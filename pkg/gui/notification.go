package gui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const duration = 5 * time.Second

type Notification struct {
	message   string
	createdAt time.Time
}

type ClearNotification struct {
	createdAt time.Time
}

func Notify(s string) tea.Cmd {
	n := Notification{message: s, createdAt: time.Now()}

	return tea.Batch(
		func() tea.Msg { return n },
		tea.Tick(duration, func(time.Time) tea.Msg {
			return ClearNotification{createdAt: n.createdAt}
		}),
	)
}
