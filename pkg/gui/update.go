package gui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
)

type ExitCode int

const (
	ExitNone ExitCode = iota
	ExitAcceptSelected
	ExitPrefillSelected
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.setSelectedRecord()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.updateTableContent()
	}
	return m, nil
}

func (m Model) onKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.ActionAcceptSelected):
		return m.onUserActionAcceptSelected()
	case key.Matches(msg, m.keys.ActionPrefillSelected):
		return m.onUserActionPrefillSelected()
	case key.Matches(msg, m.keys.ActionDeleteSelected):
		return m.onUserActionDeleteSelected()
	case key.Matches(msg, m.keys.ActionMoveDown):
		return m.onUserActionMoveDown()
	case key.Matches(msg, m.keys.ActionMoveUp):
		return m.onUserActionMoveUp()
	case key.Matches(msg, m.keys.ActionJumpDown):
		return m.onUserActionJumpDown()
	case key.Matches(msg, m.keys.ActionJumpUp):
		return m.onUserActionJumpUp()
	case key.Matches(msg, m.keys.ActionNextFilter):
		return m.onUserActionNextFilter()
	case key.Matches(msg, m.keys.ActionPrevFilter):
		return m.onUserActionPrevFilter()
	case key.Matches(msg, m.keys.ActionQuit):
		return m.onUserActionQuit()
	case key.Matches(msg, m.keys.ActionShowHelp):
		return m.onUserShowHelp()
	default:
		m.input, _ = m.input.Update(msg)
		m.updateTableContent()
	}

	return m, nil
}

func (m *Model) onUserActionMoveDown() (tea.Model, tea.Cmd) {
	if m.table.Cursor() == 0 {
		return m.onUserActionQuit()
	}
	m.table.MoveDown(1)
	return m, nil
}

func (m *Model) onUserActionMoveUp() (tea.Model, tea.Cmd) {
	m.table.MoveUp(1)
	return m, nil
}

func (m *Model) onUserActionJumpDown() (tea.Model, tea.Cmd) {
	m.table.MoveDown(10)
	return m, nil
}

func (m *Model) onUserActionJumpUp() (tea.Model, tea.Cmd) {
	m.table.MoveUp(10)
	return m, nil
}

func (m *Model) onUserActionNextFilter() (tea.Model, tea.Cmd) {
	m.filter.NextMode()
	m.updateTableContent()
	return m, nil
}

func (m *Model) onUserActionPrevFilter() (tea.Model, tea.Cmd) {
	m.filter.PrevMode()
	m.updateTableContent()
	return m, nil
}

func (m *Model) onUserActionAcceptSelected() (tea.Model, tea.Cmd) {
	return m.quit(ExitAcceptSelected)
}

func (m *Model) onUserActionPrefillSelected() (tea.Model, tea.Cmd) {
	return m.quit(ExitPrefillSelected)
}

func (m *Model) onUserActionDeleteSelected() (tea.Model, tea.Cmd) {
	if err := m.app.DeleteHistory(&m.SelectedRecord); err == nil {
		m.updateTableContent()
	}
	return m, nil
}

func (m *Model) onUserActionQuit() (tea.Model, tea.Cmd) {
	return m.quit(ExitNone)
}

func (m *Model) onUserShowHelp() (tea.Model, tea.Cmd) {
	m.help.ShowAll = !m.help.ShowAll
	return m, nil
}

func (m *Model) updateTableContent() {
	m.records = m.doSearchHistory(
		strings.Fields(m.input.Value()),
		[]config.FilterMode{m.filter.Mode},
	)

	rows := m.formatter.HistoryToTableRows(m.records)
	cols := histable.NewColumns(m.cfg.ColumnLayout, m.cfg.ShowColumnLabels, m.width-2*BorderPadding)

	m.table = histable.New(
		histable.WithRows(rows),
		histable.WithColumns(cols),
		histable.WithGotoBottom(),
		histable.WithStyles(histable.NewStyles(m.cfg.Theme)),
	)
}

func (m *Model) setSelectedRecord() {
	if cursor := m.table.Cursor(); len(m.records) > 0 && cursor < len(m.records) {
		m.SelectedRecord = m.records[cursor]
	}
}

func (m *Model) quit(exitCode ExitCode) (tea.Model, tea.Cmd) {
	m.ExitCode = exitCode
	return m, tea.Quit
}
