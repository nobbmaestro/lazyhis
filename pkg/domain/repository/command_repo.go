package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type CommandRepository struct {
	*BaseRepository[model.Command]
}

func NewCommandRepository(db *gorm.DB) *CommandRepository {
	return &CommandRepository{
		BaseRepository: &BaseRepository[model.Command]{db: db},
	}
}
