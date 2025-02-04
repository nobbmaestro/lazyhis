package model

import "gorm.io/gorm"

type Path struct {
	gorm.Model
	Path string `gorm:"unique"`
}
