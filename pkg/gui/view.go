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

	itemWidth := int(0.1 * float64(m.width))

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
	itemWidth := int(0.5 * float64(m.width))

	helpText := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(itemWidth).
		Render("Press [ctrl+q] to quit, [Tab] to cycle filter, [Enter] to select.")

	title := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(itemWidth).
		Render(titleStyle.Render(strings.Join([]string{"LazyHis", m.version}, " ")))

	row := lipgloss.JoinHorizontal(lipgloss.Top, helpText, title)

	return lipgloss.Place(m.width, 2, lipgloss.Left, lipgloss.Bottom, row)
}
