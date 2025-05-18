package gui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/app"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/help"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisfilter"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/hisquery"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type QueryHistoryCallback func(query []string, filters []config.FilterMode) []model.History

type Option func(*Model)

type Model struct {
	app *app.App
	cfg *config.GuiConfig

	records        []model.History
	SelectedRecord model.History
	formatter      formatters.Formatter

	height int
	width  int

	table  histable.Model
	input  hisquery.Model
	help   help.Model
	filter hisfilter.Model
	keys   keyMap

	version      string
	initialQuery []string
	ExitCode     ExitCode
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewGui(app *app.App, cfg *config.GuiConfig, opts ...Option) Model {
	m := Model{
		SelectedRecord: model.History{},
		ExitCode:       ExitNone,
		height:         1000,
		width:          1000,
		app:            app,
		cfg:            cfg,
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.records = m.doSearchHistory(
		m.initialQuery,
		[]config.FilterMode{utils.SafeIndex(m.cfg.CyclicFilterModes, 0)},
	)

	m.input = hisquery.New(
		hisquery.WithFocus(),
		hisquery.WithValue(strings.Join(m.initialQuery, " ")),
		hisquery.WithStyles(hisquery.NewStyles(m.cfg.Theme)),
	)

	m.keys = createKeyMap(cfg.Keys)

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
		hisfilter.WithValues(
			utils.SafeIndex(m.cfg.CyclicFilterModes, 0),
			m.cfg.CyclicFilterModes,
		),
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

func (m Model) doSearchHistory(
	query []string,
	filters []config.FilterMode,
) []model.History {
	return m.app.SearchHistory(
		app.WithQuery(query),
		app.WithFilters(append(filters, m.cfg.PersistentFilterModes...)),
	)
}
