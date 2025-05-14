package main

import (
	"os"
	"path/filepath"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/pkg/app"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/db"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/log"
	"github.com/nobbmaestro/lazyhis/pkg/registry"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var (
	dbPath   = filepath.Join(os.Getenv("HOME"), ".lazyhis.db")
	confPath = filepath.Join(os.Getenv("HOME"), ".config", "lazyhis", "lazyhis.yml")
)

func main() {
	cfg := config.ReadUserConfig(confPath)

	logger, err := log.New(cfg.Log)
	if err != nil {
		return
	}
	defer logger.Close()

	database, err := db.New(
		dbPath,
		db.WithLogger(db.DefaultLogger()),
		db.WithForeignKeysOn(),
		db.WithAutoMigrate(
			model.History{},
			model.Command{},
			model.Session{},
			model.Path{},
		),
	)
	if err != nil {
		return
	}

	historyService := service.NewHistoryService(
		&service.RepositoryProvider{
			CommandRepo: repository.NewCommandRepository(database),
			HistoryRepo: repository.NewHistoryRepository(database),
			PathRepo:    repository.NewPathRepository(database),
			SessionRepo: repository.NewSessionRepository(database),
		},
	)

	app := app.NewApp(
		app.WithService(historyService),
		app.WithLogger(logger.Logger),
		app.WithConfig(cfg),
	)

	reg := registry.NewRegistry(
		registry.WithApp(&app),
		registry.WithConfig(cfg),
		registry.WithConfigPath(confPath),
	)

	cmd.SetContext(reg.Context)
	cmd.SetVersionInfo(version, commit, date)

	err = cmd.Execute()
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
}
