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

func (r *PathRepository) QueryPaths(path string) ([]model.Path, error) {
	var paths []model.Path
	err := r.db.Model(&model.Path{}).
		Where("path LIKE ?", path).
		Find(&paths).Error
	return paths, err
}
