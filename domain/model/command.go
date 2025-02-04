package model

import "gorm.io/gorm"

type Command struct {
	gorm.Model
	Command string `gorm:"unique"`
}
