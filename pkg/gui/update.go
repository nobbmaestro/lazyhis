package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
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
	ActionNextFilter
	ActionPrevFilter
)

var keyToAction = map[tea.KeyType]Action{
	tea.KeyCtrlP:    ActionMoveUp,
	tea.KeyUp:       ActionMoveUp,
	tea.KeyCtrlN:    ActionMoveDown,
	tea.KeyDown:     ActionMoveDown,
	tea.KeyCtrlD:    ActionJumpDown,
	tea.KeyCtrlU:    ActionJumpUp,
	tea.KeyEnter:    ActionAcceptSelected,
	tea.KeyCtrlO:    ActionPrefillSelected,
	tea.KeyCtrlC:    ActionQuit,
	tea.KeyCtrlQ:    ActionQuit,
	tea.KeyEsc:      ActionQuit,
	tea.KeyTab:      ActionNextFilter,
	tea.KeyShiftTab: ActionPrevFilter,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.onKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.updateTableWidth()
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
	case ActionNextFilter:
		return m.onUserActionNextFilter()
	case ActionPrevFilter:
		return m.onUserActionPrevFilter()
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
	m.currentFilterMode = utils.Cycle(m.currentFilterMode, m.filterModes, true)
	m.updateTableContent()
	return m, nil
}

func (m *Model) onUserActionPrevFilter() (tea.Model, tea.Cmd) {
	m.currentFilterMode = utils.Cycle(m.currentFilterMode, m.filterModes, false)
	m.updateTableContent()
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

func (m *Model) updateTableWidth() {
	columns := formatters.GenerateTableColumnsFromColumns(m.columns, m.width)
	m.table.SetColumns(columns)
}

func (m *Model) updateTableContent() {
	m.records = m.queryHistory(strings.Fields(m.input.Value()), m.currentFilterMode)

	content := formatters.NewHistoryTableContent(m.records, m.columns, m.width)
	m.table = histable.New(
		histable.WithColumns(content.Columns),
		histable.WithRows(content.Rows),
	)
	m.table.GotoBottom()
}

func (m *Model) setSelectedRecord() {
	if cursor := m.table.Cursor(); cursor < len(m.records) {
		m.SelectedRecord = m.records[cursor]
	}
}

func (m *Model) setUserAction(action Action) {
	m.UserAction = action
}
