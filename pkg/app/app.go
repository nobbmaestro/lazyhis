package app

import (
	"os"
	"slices"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type App struct {
	service    *service.HistoryService
	sessionCmd string
}

type Option func(*App)

func NewApp(opts ...Option) App {
	app := App{}

	for _, opt := range opts {
		opt(&app)
	}

	return app
}

func WithService(service *service.HistoryService) Option {
	return func(app *App) {
		app.service = service
	}
}

func (app App) SearchHistory(keywords []string, filters []config.FilterMode) []model.History {
	records, err := app.service.SearchHistory(
		keywords,
		applyExitCodeFilter(filters),
		applyPathFilter(filters),
		applySessionFilter(filters, app.sessionCmd),
		-1, //maxNumSearchResults
		-1, //offsetSearchResults
		applyUniqueCommandFilter(filters),
	)
	if err != nil {
		return []model.History{}
	}
	return records
}

func (app App) AddHistory() {
}

func (app App) EditHistory() {
}

func (app App) DeleteHistory() {
}

func applyPathFilter(
	filters []config.FilterMode,
) string {
	if slices.Contains(filters, config.WorkdirFilter) ||
		slices.Contains(filters, config.WorkdirSessionFilter) {
		if p, err := os.Getwd(); err == nil {
			return p
		}
	}
	return ""
}

func applySessionFilter(
	filters []config.FilterMode,
	sessionCmd string,
) string {
	if slices.Contains(filters, config.WorkdirSessionFilter) ||
		slices.Contains(filters, config.SessionFilter) {
		if s, err := utils.RunCommand(strings.Split(sessionCmd, " ")); err == nil {
			return s
		}
	}
	return ""
}

func applyExitCodeFilter(filters []config.FilterMode) int {
	if slices.Contains(filters, config.SuccessFilter) {
		return 0
	}
	return -1
}

func applyUniqueCommandFilter(filters []config.FilterMode) bool {
	if slices.Contains(filters, config.UniqueFilter) {
		return true
	}
	return false
}
