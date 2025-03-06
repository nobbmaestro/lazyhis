package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type TmuxSessionRepository struct {
	*GenericRepository[model.TmuxSession]
}

func NewTmuxSessionRepository(db *gorm.DB) *TmuxSessionRepository {
	return &TmuxSessionRepository{
		GenericRepository: &GenericRepository[model.TmuxSession]{db: db},
	}
}

func (r *TmuxSessionRepository) QueryTmuxSessions(
	session string,
) ([]model.TmuxSession, error) {
	var sessions []model.TmuxSession
	err := r.db.Model(&model.TmuxSession{}).
		Where("session LIKE ?", session).
		Find(&sessions).Error
	return sessions, err
}
