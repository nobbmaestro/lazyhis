package repository

import (
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"

	"gorm.io/gorm"
)

type CommandRepository struct {
	*GenericRepository[model.Command]
}

func NewCommandRepository(db *gorm.DB) *CommandRepository {
	return &CommandRepository{
		GenericRepository: &GenericRepository[model.Command]{db: db},
	}
}

func (r *CommandRepository) QueryCommands(keywords []string) ([]model.Command, error) {
	var commands []model.Command
	err := r.db.Model(&model.Command{}).
		Where("command LIKE ?", "%"+strings.Join(keywords, " ")+"%").
		Find(&commands).Error
	return commands, err
}
