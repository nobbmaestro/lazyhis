package repository

import (
	"github.com/nobbmaestro/lazyhis/domain/model"
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

func (r *HistoryRepository) QueryTmuxSessions(
	session string,
) ([]model.TmuxSession, error) {
	var sessions []model.TmuxSession
	err := r.db.Model(&model.TmuxSession{}).
		Where("session LIKE ?", session).
		Find(&sessions).Error
	return sessions, err
}
