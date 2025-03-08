package gui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
)

type Action int

const (
	ActionNone Action = iota
	ActionAcceptSelected
	ActionPrefillSelected
	ActionQuit
	ActionMoveDown
	ActionMoveUp
	ActionJumpDown
	ActionJumpUp
)

var keyToAction = map[tea.KeyType]Action{
	tea.KeyCtrlP: ActionMoveUp,
	tea.KeyUp:    ActionMoveUp,
	tea.KeyCtrlN: ActionMoveDown,
	tea.KeyDown:  ActionMoveDown,
	tea.KeyCtrlD: ActionJumpDown,
	tea.KeyCtrlU: ActionJumpUp,
	tea.KeyEnter: ActionAcceptSelected,
	tea.KeyCtrlO: ActionPrefillSelected,
	tea.KeyCtrlC: ActionQuit,
	tea.KeyCtrlQ: ActionQuit,
	tea.KeyEsc:   ActionQuit,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	}
	return m, nil
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	userAction := keyToAction[msg.Type]

	switch userAction {
	case ActionAcceptSelected:
		return m.onUserActionAcceptSelected()
	case ActionPrefillSelected:
		return m.onUserActionPrefillSelected()
	case ActionQuit:
		return m.onUserActionQuit()
	case ActionMoveDown:
		return m.onUserActionMoveDown()
	case ActionMoveUp:
		return m.onUserActionMoveUp()
	case ActionJumpDown:
		return m.onUserActionJumpDown()
	case ActionJumpUp:
		return m.onUserActionJumpUp()
	default:
		m.input, _ = m.input.Update(msg)
		m.updateTableContent()
	}

	return m, nil
}

func (m *Model) onUserActionMoveDown() (tea.Model, tea.Cmd) {
	if m.table.Cursor() < len(m.records)-1 {
		m.table.MoveDown(1)
	}
	return m, nil
}

func (m *Model) onUserActionMoveUp() (tea.Model, tea.Cmd) {
	if m.table.Cursor() > 0 {
		m.table.MoveUp(1)
	}
	return m, nil
}

func (m *Model) onUserActionJumpDown() (tea.Model, tea.Cmd) {
	remaining := len(m.records) - 1 - m.table.Cursor()
	m.table.MoveDown(min(10, remaining))
	return m, nil
}

func (m *Model) onUserActionJumpUp() (tea.Model, tea.Cmd) {
	m.table.MoveUp(min(10, m.table.Cursor()))
	return m, nil
}

func (m *Model) onUserActionAcceptSelected() (tea.Model, tea.Cmd) {
	m.setUserAction(ActionAcceptSelected)
	m.setSelectedRecord()
	return m, tea.Quit
}

func (m *Model) onUserActionPrefillSelected() (tea.Model, tea.Cmd) {
	m.setUserAction(ActionPrefillSelected)
	m.setSelectedRecord()
	return m, tea.Quit
}

func (m *Model) onUserActionQuit() (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m *Model) updateTableContent() {
	m.records = m.queryHistory(strings.Fields(m.input.Value()))

	content := formatters.NewHistoryTableContent(m.records, m.columns, m.width)
	m.table = table.New(
		table.WithColumns(content.Columns),
		table.WithRows(content.Rows),
		table.WithFocused(true),
	)
}

func (m *Model) setSelectedRecord() {
	if cursor := m.table.Cursor(); cursor < len(m.records) {
		m.SelectedRecord = m.records[cursor]
	}
}

func (m *Model) setUserAction(action Action) {
	m.UserAction = action
}
