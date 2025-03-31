package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/nobbmaestro/lazyhis/pkg/config"
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

func NewStyles(theme config.GuiTheme) help.Styles {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.HelpAccentColor))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.HelpFgColor))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.BorderColor))

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
