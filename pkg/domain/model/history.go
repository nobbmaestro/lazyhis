package model

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	ExitCode      *int
	ExecutedAt    time.Time
	ExecutedIn    *int
	CommandID     *uint
	Command       *Command `gorm:"foreignKey:CommandID"`
	PathID        *uint
	Path          *Path `gorm:"foreignKey:PathID"`
	TmuxSessionID *uint
	TmuxSession   *TmuxSession `gorm:"foreignKey:TmuxSessionID"`
}
