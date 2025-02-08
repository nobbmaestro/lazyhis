package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type HistoryRepository struct {
	*BaseRepository[model.History]
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{
		BaseRepository: &BaseRepository[model.History]{db: db},
	}
}
