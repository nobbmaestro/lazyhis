package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type TmuxSessionRepository struct {
	*BaseRepository[model.TmuxSession]
}

func NewTmuxSessionRepository(db *gorm.DB) *TmuxSessionRepository {
	return &TmuxSessionRepository{
		BaseRepository: &BaseRepository[model.TmuxSession]{db: db},
	}
}
