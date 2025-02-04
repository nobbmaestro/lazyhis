package repository

import (
	"gorm.io/gorm"
)

type PathRepository struct {
	db *gorm.DB
}

func NewPathRepository(db *gorm.DB) *PathRepository {
	return &PathRepository{db: db}
}
