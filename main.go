package main

import (
	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/nobbmaestro/lazyhis/pkg/db"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
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

	ctx := context.NewContext()
	ctx = context.WithService(ctx, historyService)
	cmd.SetContext(ctx)
	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
