package gui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

var filterModeNames = map[config.FilterMode]string{
	config.NoFilter:          "-",
	config.ExitFilter:        "EXIT",
	config.PathFilter:        "PATH",
	config.SessionFilter:     "SESS",
	config.UniqueFilter:      "UNIQUE",
	config.PathSessionFilter: "PATH + SESS",
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500"))
	// cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
)

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderHistoryTable(),
		m.renderInput(),
		m.renderFooter(),
	)
}

func (m *Model) renderHistoryTable() string {
	m.table.SetWidth(m.width)
	m.table.SetHeight(m.height - 4)

	return lipgloss.NewStyle().PaddingBottom(1).Render(m.table.View())
}

func (m Model) renderFilterModeName() string {
	if len(m.filterModes) == 0 {
		return ""
	}

	itemWidth := int(18)

	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(itemWidth).
		Render(utils.CenterString(filterModeNames[m.currentFilterMode], 11, "[ %-*s ]"))
}

func (m Model) renderInput() string {
	filterModeName := m.renderFilterModeName()
	row := lipgloss.JoinHorizontal(lipgloss.Top, filterModeName, m.input.View())
	return lipgloss.Place(m.width, 1, lipgloss.Left, lipgloss.Bottom, row)
}

func (m Model) renderFooter() string {
	versionWidth := 35

	help := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.width - versionWidth).
		Render(m.help.View(Keys))

	version := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(versionWidth).
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")).
		Render(strings.Join([]string{"lazyhis", m.version}, " "))

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, help, version)

	return lipgloss.Place(m.width, 2, lipgloss.Left, lipgloss.Bottom, row)
}
