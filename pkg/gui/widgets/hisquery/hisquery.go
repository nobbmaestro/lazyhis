package hisquery

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Option func(*Model)

type Model struct {
	textinput.Model
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

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	updatedModel, cmd := m.Model.Update(msg)
	return Model{updatedModel}, cmd
}
