package app

import (
	"fmt"
	"log/slog"
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

func (app App) SearchHistory(opts ...HistoryOption) []model.History {
	historyOpts := defaultHistoryOptions()

	for _, opt := range opts {
		opt(&historyOpts)
	}

	applyFilters(
		&historyOpts,
		app.config.Os.FetchCurrentSessionCmd,
	)

	records, err := app.Service.SearchHistory(
		historyOpts.Command,
		*historyOpts.ExitCode,
		*historyOpts.Path,
		*historyOpts.Session,
		*historyOpts.MaxNumSearchResults,
		*historyOpts.OffsetSearchResults,
		*historyOpts.UniqueSearchResults,
	)
	if err != nil {
		app.logger.Error(err.Error())
		return []model.History{}
	}
	return records
}

func (app App) AddHistory(
	dryRun bool,
	verbose bool,
	addUniqueOnly bool,
	opts ...HistoryOption,
) (*model.History, error) {
	historyOpts := HistoryOptions{}

	for _, opt := range opts {
		opt(&historyOpts)
	}

	if utils.IsExcludedCommand(
		historyOpts.Command,
		app.config.Db.ExcludePrefix,
		app.config.Db.ExcludeCommands,
	) {
		return nil, nil
	}

	if addUniqueOnly && app.Service.CommandExists(historyOpts.Command) {
		return nil, nil
	}

	app.logger.Debug(
		"Add",
		"dry",
		dryRun,
		"command",
		strings.Join(historyOpts.Command, " "),
	)

	if verbose || dryRun {
		fmt.Println(strings.Join(historyOpts.Command, " "))
	}

	if dryRun {
		return nil, nil
	}

	if historyOpts.Session != nil && *historyOpts.Session == "" {
		*historyOpts.Session = app.GetCurrentSession()
	}

	return app.Service.AddHistory(
		historyOpts.Command,
		historyOpts.ExitCode,
		historyOpts.ExecutedIn,
		historyOpts.Path,
		historyOpts.Session,
	)
}

func (app App) EditHistory(
	historyID int,
	opts ...HistoryOption,
) (*model.History, error) {
	historyOpts := defaultHistoryOptions()

	for _, opt := range opts {
		opt(&historyOpts)
	}

	return app.Service.EditHistory(
		historyID,
		historyOpts.ExitCode,
		historyOpts.ExecutedIn,
		historyOpts.Path,
		historyOpts.Session,
	)
}

func (app App) PruneHistory(dryRun bool, verboseMode bool) error {
	records, err := app.Service.GetAllCommands()
	if err != nil {
		return err
	}

	for _, record := range records {
		if !utils.MatchesExclusionPatterns(
			record.Command,
			app.config.Db.ExcludeCommands,
		) {
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

func (app App) DeleteHistory(record *model.History) error {
	app.logger.Debug("Delete", "command", record.Command)
	return app.Service.DeleteCommand(record.Command)
}

func (app App) GetCurrentSession() string {
	cmd := strings.Fields(app.config.Os.FetchCurrentSessionCmd)
	if currentSession, err := utils.RunCommand(cmd); err == nil {
		return currentSession
	}
	return ""
}
