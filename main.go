package main

import (
	"context"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/db"
	"github.com/nobbmaestro/lazyhis/domain/repository"
	"github.com/nobbmaestro/lazyhis/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/ctxreg"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
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
	cmd.SetContext(ctx)
	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
