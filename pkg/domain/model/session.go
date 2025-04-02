package model

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Session string `gorm:"unique"`
}
