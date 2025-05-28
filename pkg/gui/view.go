package gui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	BorderPadding = 2
	BottomPadding = 5
)

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderTable(),
		m.renderInput(),
		m.renderFooter(),
	)
}

func (m *Model) renderTable() string {

	m.table.SetWidth(m.width - BorderPadding)
	m.table.SetHeight(m.height - BottomPadding - BorderPadding)

	return lipgloss.NewStyle().
		PaddingBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.cfg.Theme.BorderColor)).
		Render(m.table.View())
}

func (m Model) renderInput() string {
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.cfg.Theme.BorderColor))

	return inputStyle.Render(
		lipgloss.Place(
			m.width-BorderPadding,
			1,
			lipgloss.Left,
			lipgloss.Bottom,
			lipgloss.JoinHorizontal(lipgloss.Top, m.filter.View(), m.input.View()),
		),
	)
}

func (m Model) renderFooter() string {
	versionWidth := 35

	help := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.width - versionWidth - BorderPadding).
		Render(m.help.View(m.keys))

	version := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(versionWidth).
		Bold(true).
		Foreground(lipgloss.Color(m.cfg.Theme.VersionFgColor)).
		Render(strings.Join([]string{"lazyhis", *m.app.GetVersion()}, " "))

	return lipgloss.Place(
		m.width,
		1,
		lipgloss.Left,
		lipgloss.Bottom,
		lipgloss.JoinHorizontal(lipgloss.Bottom, help, version),
	)
}
