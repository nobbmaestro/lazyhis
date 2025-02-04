package main

import (
	"context"

	"github.com/nobbmaestro/lazyhis/cmd"
	"github.com/nobbmaestro/lazyhis/pkg/db"
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
	serviceCtx := service.NewServiceContext(database)

	cmd.SetContext(context.WithValue(
		context.Background(),
		service.ServiceCtxKey,
		serviceCtx),
	)
	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
