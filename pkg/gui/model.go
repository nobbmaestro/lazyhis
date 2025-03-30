package gui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/help"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisfilter"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisquery"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
)

type QueryHistoryCallback func(keywords []string, mode config.FilterMode) []model.History

type Model struct {
	formatter      formatters.Formatter
	cfg            config.GuiConfig
	records        []model.History
	table          histable.Model
	input          hisquery.Model
	help           help.Model
	filter         hisfilter.Model
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

	formatter := formatters.NewFormatter(formatters.WithColumns(cfg.ColumnLayout))

	rows := formatter.HistoryToTableRows(records)
	cols := histable.NewColumns(cfg.ColumnLayout, cfg.ShowColumnLabels, 100)
	historyTable := histable.New(
		histable.WithRows(rows),
		histable.WithColumns(cols),
		histable.WithStyles(table.DefaultStyles()),
		histable.WithGotoBottom(),
	)

	help := help.New(
		help.WithStyles(help.NewStyles()),
	)

	historyFilter := hisfilter.New(
		cfg.InitialFilterMode,
		cfg.CyclicFilterModes,
	)

	return Model{
		formatter:      formatter,
		cfg:            cfg,
		records:        records,
		table:          historyTable,
		input:          historyQuery,
		filter:         historyFilter,
		help:           help,
		height:         10,
		width:          10,
		version:        version,
		queryHistory:   queryHistory,
		SelectedRecord: model.History{},
		UserAction:     ActionNone,
	}
}
