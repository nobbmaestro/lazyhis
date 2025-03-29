package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

type Option func(*Model)

type Model struct {
	help.Model
}

func New(opts ...Option) Model {
	m := Model{help.New()}

	m.ShortSeparator = " |"

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func NewStyles() help.Styles {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#606064"))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#464959"))

	return help.Styles{
		ShortKey:       keyStyle,
		ShortDesc:      descStyle,
		ShortSeparator: sepStyle,
		FullKey:        keyStyle,
		FullDesc:       descStyle,
		FullSeparator:  sepStyle,
	}
}

func WithStyles(styles help.Styles) Option {
	return func(m *Model) {
		m.Styles = styles
	}
}
