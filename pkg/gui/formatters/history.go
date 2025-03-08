package formatters

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type getHistoryFieldByColumn func(history model.History) string

var tableColumnNames = map[config.Column]string{
	config.ColumnCommand:    "Command",
	config.ColumnExecutedAt: "Executed",
	config.ColumnExecutedIn: "Duration",
	config.ColumnExitCode:   "Exit",
	config.ColumnPath:       "Path",
	config.ColumnSession:    "Session",
}

var tableColumnWidth = map[config.Column]int{
	config.ColumnCommand:    100,
	config.ColumnExecutedAt: 10,
	config.ColumnExecutedIn: 10,
	config.ColumnExitCode:   5,
	config.ColumnPath:       50,
	config.ColumnSession:    15,
}

var columnToGetter = map[config.Column]getHistoryFieldByColumn{
	config.ColumnCommand:    extractCommandFromHistory,
	config.ColumnExecutedAt: extractExecutedAtFromHistory,
	config.ColumnExecutedIn: extractExecutedInFromHistory,
	config.ColumnExitCode:   extractExitCodeFromHistory,
	config.ColumnPath:       extractPathFromHistory,
	config.ColumnSession:    extractSessionFromHistory,
}

type HistoryTableContent struct {
	Columns []table.Column
	Rows    []table.Row
}

func NewHistoryTableContent(
	records []model.History,
	columns []config.Column,
) HistoryTableContent {
	return HistoryTableContent{
		Columns: generateTableColumnsFromColumns(columns),
		Rows:    generateTableRowsFromHistory(records, columns),
	}
}

func generateTableColumnsFromColumns(columns []config.Column) []table.Column {
	tableColumns := make([]table.Column, len(columns))
	for i, column := range columns {
		tableColumns[i].Title = tableColumnNames[column]
		tableColumns[i].Width = tableColumnWidth[column]
	}
	return tableColumns
}

func generateTableRowsFromHistory(
	records []model.History,
	columns []config.Column,
) []table.Row {
	rows := make([]table.Row, len(records))

	for i, history := range records {
		row := make([]string, len(columns))

		for j, column := range columns {
			if getter, ok := columnToGetter[column]; ok {
				row[j] = getter(history)
			} else {
				row[j] = ""
			}
		}
		rows[i] = row
	}
	return rows
}

func extractCommandFromHistory(history model.History) string {
	if history.Command != nil {
		return history.Command.Command
	}
	return ""
}

func extractExecutedAtFromHistory(history model.History) string {
	return utils.HumanizeTimeAgo(history.CreatedAt)
}

func extractExecutedInFromHistory(history model.History) string {
	if history.ExecutedIn != nil {
		return utils.HumanizeDuration(int64(*history.ExecutedIn))
	}
	return ""
}

func extractExitCodeFromHistory(history model.History) string {
	if history.ExitCode != nil {
		return fmt.Sprintf("%3d", *history.ExitCode)
	}
	return ""
}

func extractPathFromHistory(history model.History) string {
	if history.Path != nil {
		return utils.HumanizePath(history.Path.Path)
	}
	return ""
}

func extractSessionFromHistory(history model.History) string {
	if history.Session != nil {
		return history.Session.Session
	}
	return ""
}
