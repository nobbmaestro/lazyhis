package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type SessionRepository struct {
	*GenericRepository[model.Session]
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		GenericRepository: &GenericRepository[model.Session]{db: db},
	}
}

func (r *SessionRepository) QuerySessions(
	session string,
) ([]model.Session, error) {
	var sessions []model.Session
	err := r.db.Model(&model.Session{}).
		Where("session LIKE ?", session).
		Find(&sessions).Error
	return sessions, err
}
