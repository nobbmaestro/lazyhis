package gui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/help"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisquery"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
)

type QueryHistoryCallback func(keywords []string, mode config.FilterMode) []model.History

type Model struct {
	currentFilterMode config.FilterMode
	filterModes       []config.FilterMode
	columns           []config.Column
	records           []model.History
	table             histable.Model
	input             hisquery.Model
	help              help.Model
	height            int
	width             int
	version           string
	queryHistory      QueryHistoryCallback
	SelectedRecord    model.History
	UserAction        Action
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(
	cfg config.GuiConfig,
	queryHistory QueryHistoryCallback,
	version string,
	searchKeywords string,
) Model {
	records := queryHistory([]string{}, cfg.InitialFilterMode)

	historyQuery := hisquery.New(
		hisquery.WithValue(searchKeywords),
		hisquery.WithFocus(),
	)

	content := formatters.NewHistoryTableContent(records, cfg.ColumnLayout, 100)
	historyTable := histable.New(
		histable.WithColumns(content.Columns),
		histable.WithRows(content.Rows),
		histable.WithStyles(table.DefaultStyles()),
		histable.WithGotoBottom(),
	)

	help := help.New(
		help.WithStyles(help.NewStyles()),
	)

	return Model{
		columns:           cfg.ColumnLayout,
		currentFilterMode: cfg.InitialFilterMode,
		filterModes:       cfg.CyclicFilterModes,
		records:           records,
		table:             historyTable,
		input:             historyQuery,
		help:              help,
		height:            10,
		width:             10,
		version:           version,
		queryHistory:      queryHistory,
		SelectedRecord:    model.History{},
		UserAction:        ActionNone,
	}
}
