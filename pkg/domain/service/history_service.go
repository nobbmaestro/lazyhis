package service

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
)

type HistoryService struct {
	historyRepo *repository.HistoryRepository
	commandRepo *repository.CommandRepository
	pathRepo    *repository.PathRepository
	tmuxRepo    *repository.TmuxSessionRepository
}

func NewHistoryService(
	historyRepo *repository.HistoryRepository,
	commandRepo *repository.CommandRepository,
	pathRepo *repository.PathRepository,
	tmuxRepo *repository.TmuxSessionRepository,
) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		commandRepo: commandRepo,
		pathRepo:    pathRepo,
		tmuxRepo:    tmuxRepo,
	}
}
