package app

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type App struct {
	Service *service.HistoryService
	config  *config.UserConfig
	logger  *slog.Logger
}

type Option func(*App)

type Params struct {
	Command    []string
	ExitCode   int
	Path       string
	Session    string
	ExecutedIn int

	MaxNumSearchResults int
	OffsetSearchResults int
	UniqueSearchResults bool
	AddUniqueOnly       bool

	DryRun  bool
	Verbose bool
}

func NewApp(opts ...Option) App {
	app := App{}

	for _, opt := range opts {
		opt(&app)
	}

	return app
}

func WithService(service *service.HistoryService) Option {
	return func(app *App) {
		app.Service = service
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(app *App) {
		app.logger = logger
	}
}

func WithConfig(config *config.UserConfig) Option {
	return func(app *App) {
		app.config = config
	}
}

func (app App) SearchHistory(
	keywords []string,
	exitCode int,
	path string,
	session string,
	executedIn int,
	maxNumSearchResults int,
	offsetSearchResults int,
	uniqueSearchResults bool,
) []model.History {
	records, err := app.Service.SearchHistory(
		keywords,
		exitCode,
		path,
		session,
		maxNumSearchResults,
		offsetSearchResults,
		uniqueSearchResults,
	)
	if err != nil {
		app.logger.Error(err.Error())
		return []model.History{}
	}
	return records
}

func (app App) SearchHistoryWithFilters(
	keywords []string,
	filters []config.FilterMode,
) []model.History {
	records, err := app.Service.SearchHistory(
		keywords,
		applySuccessFilter(filters),
		applyPathFilter(filters),
		applySessionFilter(filters, app.config.Os.FetchCurrentSessionCmd),
		-1, //maxNumSearchResults
		-1, //offsetSearchResults
		applyUniqueCommandFilter(filters),
	)
	if err != nil {
		app.logger.Error(err.Error())
		return []model.History{}
	}
	return records
}

func (app App) AddHistory(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	session *string,
	dryRun bool,
	verbose bool,
	addUniqueOnly bool,
) (*model.History, error) {
	if utils.IsExcludedCommand(
		command,
		app.config.Db.ExcludePrefix,
		app.config.Db.ExcludeCommands,
	) {
		return nil, nil
	}

	app.logger.Debug("Add", "dry", dryRun, "command", strings.Join(command, " "))

	if verbose || dryRun {
		fmt.Println(strings.Join(command, " "))
	}

	if dryRun {
		return nil, nil
	}

	if addUniqueOnly && app.Service.CommandExists(command) {
		return nil, nil
	}

	if *session == "" {
		cmd := strings.Fields(app.config.Os.FetchCurrentSessionCmd)
		if currentSession, err := utils.RunCommand(cmd); err == nil {
			*session = currentSession
		}
	}

	return app.Service.AddHistory(
		command,
		exitCode,
		executedIn,
		path,
		session,
	)
}

func (app App) EditHistory(
	historyID int,
	exitCode *int,
	executedIn *int,
	path *string,
	session *string,
) (*model.History, error) {
	return app.Service.EditHistory(
		historyID,
		exitCode,
		executedIn,
		path,
		session,
	)
}

func (app App) DeleteHistory() {
}

func (app App) PruneHistory(dryRun bool, verboseMode bool) error {
	records, err := app.Service.GetAllCommands()
	if err != nil {
		return err
	}

	for _, record := range records {
		if !utils.MatchesExclusionPatterns(record.Command, app.config.Db.ExcludeCommands) {
			continue
		}

		app.logger.Debug("Prune", "dry", dryRun, "command", record.Command)

		if dryRun || verboseMode {
			fmt.Println(record.Command)
		}

		if dryRun {
			continue
		}

		err := app.Service.DeleteCommand(&record)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyPathFilter(filters []config.FilterMode) string {
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

func applySuccessFilter(filters []config.FilterMode) int {
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
