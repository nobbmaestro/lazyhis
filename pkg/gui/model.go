package gui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/gui/widgets/histable"
)

type QueryHistoryCallback func(keywords []string, mode config.FilterMode) []model.History

type Model struct {
	currentFilterMode config.FilterMode
	filterModes       []config.FilterMode
	columns           []config.Column
	records           []model.History
	table             histable.Model
	input             textinput.Model
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
	input := textinput.New()
	input.SetValue(searchKeywords)
	input.Focus()

	records := queryHistory([]string{}, cfg.InitialFilterMode)

	content := formatters.NewHistoryTableContent(records, cfg.ColumnLayout, 100)
	historyTable := histable.New(
		histable.WithColumns(content.Columns),
		histable.WithRows(content.Rows),
	)
	historyTable.SetStyles(table.DefaultStyles())
	historyTable.GotoBottom()

	return Model{
		columns:           cfg.ColumnLayout,
		currentFilterMode: cfg.InitialFilterMode,
		filterModes:       cfg.CyclicFilterModes,
		records:           records,
		table:             historyTable,
		input:             input,
		height:            10,
		width:             10,
		version:           version,
		queryHistory:      queryHistory,
		SelectedRecord:    model.History{},
		UserAction:        ActionNone,
	}
}
