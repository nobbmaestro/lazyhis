package main

import (
	"os"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/nobbmaestro/lazyhis/pkg/db"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/log"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cfg, err := config.ReadUserConfig()
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

	ctx := context.NewContext()
	ctx = context.WithService(ctx, historyService)
	ctx = context.WithConfig(ctx, cfg)
	ctx = context.WithLogger(ctx, logger.Logger)

	cmd.SetContext(ctx)
	cmd.SetVersionInfo(version, commit, date)

	err = cmd.Execute()
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
}
