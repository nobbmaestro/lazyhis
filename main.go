package main

import (
	"context"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/db"
	"github.com/nobbmaestro/lazyhis/domain/repository"
	"github.com/nobbmaestro/lazyhis/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/ctxreg"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cfg, err := config.LoadUserConfig()
	if err != nil {
		return
	}

	database, err := db.NewDatabaseConnection()
	if err != nil {
		return
	}

	historyRepo := repository.NewHistoryRepository(database)
	commandRepo := repository.NewCommandRepository(database)
	pathRepo := repository.NewPathRepository(database)
	tmuxRepo := repository.NewTmuxSessionRepository(database)

	historyService := service.NewHistoryService(
		historyRepo,
		commandRepo,
		pathRepo,
		tmuxRepo,
	)

	ctx := context.Background()
	ctx = ctxreg.WithService(ctx, historyService)
	ctx = ctxreg.WithConfig(ctx, cfg)

	cmd.SetContext(ctx)
	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
