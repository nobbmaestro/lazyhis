package histable

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func DefaultStyles() table.Styles {
	return table.Styles{
		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00e8c6")).
			Background(lipgloss.Color("00FFAA")),
		Header: lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:   lipgloss.NewStyle().Padding(0, 1),
	}
}

const paddedRows = 100

type Model struct {
	table.Model
}

func New(opts ...table.Option) Model {
	return Model{Model: table.New(opts...)}
}

func WithRows(rows []table.Row) table.Option {
	return func(m *table.Model) {
		paddedRows := append(
			make([]table.Row, paddedRows),
			reverse(rows)...,
		)
		m.SetRows(paddedRows)
	}
}

func WithColumns(columns []table.Column) table.Option {
	return func(m *table.Model) {
		m.SetColumns(columns)
	}
}

func WithStyles(styles table.Styles) table.Option {
	return func(m *table.Model) {
		m.SetStyles(styles)
	}
}

func WithGotoBottom() table.Option {
	return func(m *table.Model) {
		m.GotoBottom()
	}
}

func (m Model) Cursor() int {
	return len(m.Rows()) - 1 - paddedRows - m.realCursor()
}

// Non-reversed (Real) Curosr for internal usage
func (m Model) realCursor() int {
	return max(m.Model.Cursor()-paddedRows, 0)
}

func (m *Model) MoveUp(n int) {
	if m.realCursor()-n >= 0 {
		m.Model.MoveUp(n)
	}
}

func (m *Model) MoveDown(n int) {
	if m.realCursor()+n <= len(m.Rows()) {
		m.Model.MoveDown(n)
	}
}

func reverse[T any](rows []T) []T {
	newRows := make([]T, len(rows))
	for i, row := range rows {
		newRows[len(rows)-1-i] = row
	}
	return newRows
}
