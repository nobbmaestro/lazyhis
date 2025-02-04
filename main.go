package main

import (
	"context"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/db"
	"github.com/nobbmaestro/lazyhis/domain/service"
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
	serviceCtx := service.NewServiceContext(database)

	cmd.SetContext(context.WithValue(
		context.Background(),
		service.ServiceCtxKey,
		serviceCtx),
	)
	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
