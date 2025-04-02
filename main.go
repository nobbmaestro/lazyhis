package main

import (
	"os"
	"path/filepath"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/db"
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
	confPath = filepath.Join(os.Getenv("HOME"), ".config", "lazyhis", "lazyhis.yml")
)

func main() {
	cfg, err := config.ReadUserConfig(confPath)
	if err != nil {
		return
	}

	database, err := db.New()
	if err != nil {
		return
	}

	logger, err := log.New(cfg.Log)
	if err != nil {
		return
	}
	defer logger.Close()

	historyService := service.NewHistoryService(
		&service.RepositoryProvider{
			CommandRepo: repository.NewCommandRepository(database),
			HistoryRepo: repository.NewHistoryRepository(database),
			PathRepo:    repository.NewPathRepository(database),
			SessionRepo: repository.NewSessionRepository(database),
		},
		&cfg.Db,
		logger.Logger,
	)

	reg := registry.NewRegistry(
		registry.WithConfig(cfg),
		registry.WithConfigPath(confPath),
		registry.WithLogger(logger.Logger),
		registry.WithService(historyService),
	)

	cmd.SetContext(reg.Context)
	cmd.SetVersionInfo(version, commit, date)

	err = cmd.Execute()
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
}
