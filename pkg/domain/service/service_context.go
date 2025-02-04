package service

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"gorm.io/gorm"
)

type contextKey string

type ServiceContext struct {
	HistoryService *HistoryService
}

const ServiceCtxKey contextKey = "serviceCtx"

func NewServiceContext(db *gorm.DB) *ServiceContext {
	historyRepo := repository.NewHistoryRepository(db)
	commandRepo := repository.NewCommandRepository(db)
	pathRepo := repository.NewPathRepository(db)
	tmuxRepo := repository.NewTmuxSessionRepository(db)

	historyService := NewHistoryService(historyRepo, commandRepo, pathRepo, tmuxRepo)

	return &ServiceContext{
		HistoryService: historyService,
	}
}
