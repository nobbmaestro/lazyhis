package hisquery

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nobbmaestro/lazyhis/pkg/config"
)

type Option func(*Model)

type Style struct {
	Prompt    string
	TextStyle lipgloss.Style
}

type Model struct {
	textinput.Model
}

func NewStyles(theme config.GuiTheme) Style {
	return Style{
		Prompt:    "  ",
		TextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.InputFgColor)),
	}
}

func New(opts ...Option) Model {
	m := Model{textinput.New()}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func WithValue(value string) Option {
	return func(m *Model) {
		m.SetValue(value)
	}
}

func WithFocus() Option {
	return func(m *Model) {
		m.Focus()
	}
}

func WithStyles(style Style) Option {
	return func(m *Model) {
		m.Prompt = style.Prompt
		m.TextStyle = style.TextStyle
	}
}

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	updatedModel, cmd := m.Model.Update(msg)
	return Model{updatedModel}, cmd
}
