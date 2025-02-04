package repository

import (
	"gorm.io/gorm"
)

type CommandRepository struct {
	db *gorm.DB
}

func NewCommandRepository(db *gorm.DB) *CommandRepository {
	return &CommandRepository{db: db}
}
