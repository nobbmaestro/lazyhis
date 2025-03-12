package gui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500"))
	// cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
)

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderHistoryTable(),
		m.input.View(),
		m.renderFooter(),
	)
}

func (m *Model) renderHistoryTable() string {
	m.table.SetWidth(m.width)
	m.table.SetHeight(m.height - 4)

	return lipgloss.NewStyle().PaddingBottom(1).Render(m.table.View())
}

func (m Model) renderFooter() string {
	itemWidth := int(0.5 * float64(m.width))

	helpText := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(itemWidth).
		Render("Press [ctrl+q] to quit, [Enter] to select.")

	title := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(itemWidth).
		Render(titleStyle.Render(strings.Join([]string{"LazyHis", m.version}, " ")))

	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, helpText, title)

	return lipgloss.Place(m.width, 2, lipgloss.Left, lipgloss.Bottom, bottomRow)
}
