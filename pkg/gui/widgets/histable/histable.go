package histable

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/nobbmaestro/lazyhis/pkg/config"
)

const rowPadding = 100

type Model struct {
	table.Model
}

func New(opts ...table.Option) Model {
	return Model{Model: table.New(opts...)}
}

func NewColumns(
	columns []config.Column,
	showLabels bool,
	titles map[config.Column]string,
	widths map[config.Column]int,
	width int,
) []table.Column {
	tableColumns := make([]table.Column, len(columns))

	newTableColumnWidth := calculateTableColumnWidth(columns, width, widths)
	for i, column := range columns {
		columnTitle := titles[column]
		if !showLabels {
			columnTitle = ""
		}
		tableColumns[i].Title = columnTitle
		tableColumns[i].Width = newTableColumnWidth[column]
	}
	return tableColumns
}

func NewStyles(theme config.GuiTheme) table.Styles {
	return table.Styles{
		Header: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color(theme.TableLabelsFgColor)).
			Bold(true),

		Cell: lipgloss.NewStyle().
			Padding(0, 1),

		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TableCursorFgColor)).
			Background(lipgloss.Color(theme.TableCursorBgColor)).
			Bold(true),
	}
}

func WithRows(rows []table.Row) table.Option {
	return func(m *table.Model) {
		paddedRows := append(
			make([]table.Row, rowPadding),
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

func WithCursor(cursor int) table.Option {
	return func(m *table.Model) {
		m.GotoBottom()
		m.MoveUp(cursor)
	}
}

func WithGotoBottom() table.Option {
	return func(m *table.Model) {
		m.GotoBottom()
	}
}

func (m Model) Cursor() int {
	return len(m.Rows()) - 1 - rowPadding - m.realCursor()
}

// Non-reversed (Real) Curosr for internal usage
func (m Model) realCursor() int {
	return max(m.Model.Cursor()-rowPadding, 0)
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

func calculateTableColumnWidth(
	columns []config.Column,
	totalWidth int,
	widths map[config.Column]int,
) map[config.Column]int {
	newTableColumnWidth := make(map[config.Column]int)

	totalStaticWidth := 0
	for _, column := range columns {
		if column != config.ColumnCommand {
			totalStaticWidth += widths[column]
		}
	}

	remainingWidth := totalWidth - totalStaticWidth - 2*(len(columns)-1)
	remainingWidth = max(remainingWidth, 0)

	for _, column := range columns {
		if column == config.ColumnCommand {
			newTableColumnWidth[column] = remainingWidth
		} else {
			newTableColumnWidth[column] = widths[column]
		}
	}

	return newTableColumnWidth
}
