package gui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
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
	ActionShowHelp
)

type keyMap struct {
	ActionAcceptSelected  key.Binding
	ActionPrefillSelected key.Binding
	ActionNextFilter      key.Binding
	ActionPrevFilter      key.Binding
	ActionJumpDown        key.Binding
	ActionJumpUp          key.Binding
	ActionMoveDown        key.Binding
	ActionMoveUp          key.Binding
	ActionQuit            key.Binding
	ActionShowHelp        key.Binding
}

var Keys = keyMap{
	ActionMoveUp: key.NewBinding(
		key.WithKeys("ctrl+p", "up"),
		key.WithHelp(" ⌃p", "move up"),
	),
	ActionMoveDown: key.NewBinding(
		key.WithKeys("ctrl+n", "down"),
		key.WithHelp(" ⌃n", "move down"),
	),
	ActionJumpUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp(" ⌃u", "jump up"),
	),
	ActionJumpDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp(" ⌃d", "jump down"),
	),
	ActionAcceptSelected: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("  ↵", "accept"),
	),
	ActionPrefillSelected: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp(" ⌃o", "prefill"),
	),
	ActionQuit: key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+q", "esc"),
		key.WithHelp(" ⌃q", "quit"),
	),
	ActionShowHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("  ?", "toggle help"),
	),
	ActionNextFilter: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("  ⇥", "next filter"),
	),
	ActionPrevFilter: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp(" ⇧⇥", "prev filter"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ActionShowHelp,
		k.ActionAcceptSelected,
		k.ActionNextFilter,
		k.ActionQuit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ActionAcceptSelected, k.ActionPrefillSelected},
		{k.ActionNextFilter, k.ActionPrevFilter},
		{k.ActionMoveUp, k.ActionMoveDown},
		{k.ActionJumpDown, k.ActionJumpUp},
		{k.ActionShowHelp, k.ActionQuit},
	}
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
	switch {
	case key.Matches(msg, Keys.ActionAcceptSelected):
		return m.onUserActionAcceptSelected()
	case key.Matches(msg, Keys.ActionPrefillSelected):
		return m.onUserActionPrefillSelected()
	case key.Matches(msg, Keys.ActionMoveDown):
		return m.onUserActionMoveDown()
	case key.Matches(msg, Keys.ActionMoveUp):
		return m.onUserActionMoveUp()
	case key.Matches(msg, Keys.ActionJumpDown):
		return m.onUserActionJumpDown()
	case key.Matches(msg, Keys.ActionJumpUp):
		return m.onUserActionJumpUp()
	case key.Matches(msg, Keys.ActionNextFilter):
		return m.onUserActionNextFilter()
	case key.Matches(msg, Keys.ActionPrevFilter):
		return m.onUserActionPrevFilter()
	case key.Matches(msg, Keys.ActionQuit):
		return m.onUserActionQuit()
	case key.Matches(msg, Keys.ActionShowHelp):
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

func (m *Model) onUserShowHelp() (tea.Model, tea.Cmd) {
	m.help.ShowAll = !m.help.ShowAll
	return m, nil
}

func (m *Model) updateTableWidth() {
	columns := histable.NewColumns(m.cfg.ColumnLayout, m.cfg.ShowColumnLabels, m.width)
	m.table.SetColumns(columns)
}

func (m *Model) updateTableContent() {
	m.records = m.queryHistory(strings.Fields(m.input.Value()), m.filter.Mode)

	rows := formatters.HistoryToTableRows(m.records, m.cfg.ColumnLayout)
	cols := histable.NewColumns(m.cfg.ColumnLayout, m.cfg.ShowColumnLabels, m.width)
	m.table = histable.New(
		histable.WithRows(rows),
		histable.WithColumns(cols),
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
