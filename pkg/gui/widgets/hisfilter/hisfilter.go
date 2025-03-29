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

type Model struct {
	Mode  config.FilterMode
	Modes []config.FilterMode
}

func New(mode config.FilterMode, modes []config.FilterMode) Model {
	return Model{
		Mode:  mode,
		Modes: modes,
	}
}

func (m Model) View() string {
	if len(m.Modes) == 0 {
		return ""
	}

	itemWidth := int(18)

	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(itemWidth).
		Render(utils.CenterString(filterModeNames[m.Mode], 11, "[ %-*s ]"))
}

func (m *Model) NextMode() {
	m.Mode = utils.Cycle(m.Mode, m.Modes, true)
}

func (m *Model) PrevMode() {
	m.Mode = utils.Cycle(m.Mode, m.Modes, false)
}
