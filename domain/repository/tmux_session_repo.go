package repository

import (
	"gorm.io/gorm"
)

type TmuxSessionRepository struct {
	db *gorm.DB
}

func NewTmuxSessionRepository(db *gorm.DB) *TmuxSessionRepository {
	return &TmuxSessionRepository{db: db}
}
