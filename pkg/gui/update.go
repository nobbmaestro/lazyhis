package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlP, tea.KeyUp:
			m.table.MoveUp(1)

		case tea.KeyCtrlN, tea.KeyDown:
			m.table.MoveDown(1)

		case tea.KeyCtrlU:
			m.table.MoveUp(10)

		case tea.KeyCtrlD:
			m.table.MoveDown(10)

		case tea.KeyEnter:
			m.setSelectedRecord()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyCtrlQ, tea.KeyEsc:
			return m, tea.Quit

		default:
			m.input, _ = m.input.Update(msg)
			m.updateRecords()
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	}
	return m, nil
}

func (m *Model) updateRecords() {
	m.records = m.queryHistory(strings.Split(m.input.Value(), " "))
}

func (m *Model) setSelectedRecord() {
	if cursor := m.table.Cursor(); cursor < len(m.records) {
		m.selectedRecord = m.records[cursor]
	}
}
