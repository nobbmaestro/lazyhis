package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/help"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisfilter"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisquery"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
)

type QueryHistoryCallback func(query []string, mode config.FilterMode) []model.History

type Option func(*Model)

type Model struct {
	cfg            config.GuiConfig
	records        []model.History
	SelectedRecord model.History
	formatter      formatters.Formatter

	height int
	width  int

	table  histable.Model
	input  hisquery.Model
	help   help.Model
	filter hisfilter.Model

	version      string
	initialQuery []string
	queryHistory QueryHistoryCallback
	UserAction   Action
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewGui(cb QueryHistoryCallback, cfg config.GuiConfig, opts ...Option) Model {
	m := Model{
		SelectedRecord: model.History{},
		UserAction:     ActionNone,
		height:         10,
		width:          100,
		queryHistory:   cb,
		cfg:            cfg,
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.records = m.queryHistory(m.initialQuery, m.cfg.InitialFilterMode)

	m.input = hisquery.New(
		hisquery.WithFocus(),
		hisquery.WithValue(strings.Join(m.initialQuery, " ")),
		hisquery.WithStyles(hisquery.NewStyles(m.cfg.Theme)),
	)

	rows := m.formatter.HistoryToTableRows(m.records)
	cols := histable.NewColumns(m.cfg.ColumnLayout, m.cfg.ShowColumnLabels, m.width)
	m.table = histable.New(
		histable.WithRows(rows),
		histable.WithColumns(cols),
		histable.WithStyles(histable.NewStyles(m.cfg.Theme)),
		histable.WithGotoBottom(),
	)

	m.help = help.New(
		help.WithStyles(help.NewStyles(m.cfg.Theme)),
	)

	m.filter = hisfilter.New(
		hisfilter.WithValues(m.cfg.InitialFilterMode, m.cfg.CyclicFilterModes),
		hisfilter.WithStyles(hisfilter.NewStyles(m.cfg.Theme)),
	)

	return m
}

func WithVersion(ver string) Option {
	return func(m *Model) {
		m.version = ver
	}
}

func WithInitialQuery(query []string) Option {
	return func(m *Model) {
		m.initialQuery = query
	}
}

func WithFormatter(fmt formatters.Formatter) Option {
	return func(m *Model) {
		m.formatter = fmt
	}
}
