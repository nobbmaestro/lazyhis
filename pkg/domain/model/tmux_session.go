package model

import (
	"gorm.io/gorm"
)

type TmuxSession struct {
	gorm.Model
	Session string `gorm:"unique"`
}
