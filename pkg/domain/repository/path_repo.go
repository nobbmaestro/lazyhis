package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type PathRepository struct {
	*BaseRepository[model.Path]
}

func NewPathRepository(db *gorm.DB) *PathRepository {
	return &PathRepository{BaseRepository: &BaseRepository[model.Path]{db: db}}
}
