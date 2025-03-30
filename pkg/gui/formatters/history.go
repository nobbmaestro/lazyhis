package formatters

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type getHistoryFieldByColumn func(history model.History) string

var columnToGetter = map[config.Column]getHistoryFieldByColumn{
	config.ColumnCommand:    extractCommandFromHistory,
	config.ColumnExecutedAt: extractExecutedAtFromHistory,
	config.ColumnExecutedIn: extractExecutedInFromHistory,
	config.ColumnExitCode:   extractExitCodeFromHistory,
	config.ColumnID:         extractIDFromHistory,
	config.ColumnPath:       extractPathFromHistory,
	config.ColumnSession:    extractSessionFromHistory,
}

type Formatter struct {
	columns []config.Column
}

type Option func(*Formatter)

func NewFormatter(opts ...Option) Formatter {
	m := Formatter{}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func WithColumns(columns []config.Column) Option {
	return func(f *Formatter) {
		f.columns = columns
	}
}

func (f Formatter) HistoryToTableRows(records []model.History) []table.Row {
	rows := make([]table.Row, len(records))

	for i, history := range records {
		row := make([]string, len(f.columns))

		for j, column := range f.columns {
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

func extractIDFromHistory(history model.History) string {
	return strconv.FormatUint(uint64(history.ID), 10)
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
