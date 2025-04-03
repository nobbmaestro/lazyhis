package hisfilter

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

var filterModeNames = map[config.FilterMode]string{
	config.NoFilter:             "-",
	config.SuccessFilter:        "SUCCESS",
	config.WorkdirFilter:        "WORKDIR",
	config.SessionFilter:        "SESSION",
	config.UniqueFilter:         "UNIQUE",
	config.WorkdirSessionFilter: "WDIR + SESS",
}

type Option func(*Model)

type Styles struct {
	TextStyle lipgloss.Style
}

type Model struct {
	Mode   config.FilterMode
	Modes  []config.FilterMode
	styles Styles
}

func New(opts ...Option) Model {
	m := Model{}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func NewStyles(theme config.GuiTheme) Styles {
	return Styles{
		TextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilterFgColor)),
	}
}

func WithStyles(styles Styles) Option {
	return func(m *Model) {
		m.styles = styles
	}
}

func WithValues(mode config.FilterMode, modes []config.FilterMode) Option {
	return func(m *Model) {
		m.Mode = mode
		m.Modes = modes
	}
}

func (m Model) View() string {
	if len(m.Modes) == 0 {
		return ""
	}

	itemWidth := int(18)

	return m.styles.TextStyle.Render(
		lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(itemWidth).
			Render(utils.CenterString(filterModeNames[m.Mode], 11, "[ %-*s ]")),
	)
}

func (m *Model) NextMode() {
	m.Mode = utils.Cycle(m.Mode, m.Modes, true)
}

func (m *Model) PrevMode() {
	m.Mode = utils.Cycle(m.Mode, m.Modes, false)
}
