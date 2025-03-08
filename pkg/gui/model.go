package gui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
)

type QueryHistoryCallback func(keywords []string) []model.History

type Model struct {
	records        []model.History
	columns        []config.Column
	table          table.Model
	input          textinput.Model
	height         int
	width          int
	version        string
	queryHistory   QueryHistoryCallback
	SelectedRecord model.History
	UserAction     Action
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(
	columns []config.Column,
	queryHistory QueryHistoryCallback,
	version string,
	searchKeywords string,
) Model {
	input := textinput.New()
	input.SetValue(searchKeywords)
	input.Focus()

	records := queryHistory([]string{})

	content := formatters.NewHistoryTableContent(records, columns, 100)
	historyTable := table.New(
		table.WithColumns(content.Columns),
		table.WithRows(content.Rows),
		table.WithFocused(true),
	)
	historyTable.SetStyles(table.DefaultStyles())

	return Model{
		records:        records,
		columns:        columns,
		table:          historyTable,
		input:          input,
		height:         10,
		width:          10,
		version:        version,
		queryHistory:   queryHistory,
		SelectedRecord: model.History{},
		UserAction:     ActionNone,
	}
}
