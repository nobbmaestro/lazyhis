package model

import "gorm.io/gorm"

type History struct {
	gorm.Model
	ExitCode   *int
	ExecutedIn *int
	CommandID  *uint
	Command    *Command `gorm:"foreignKey:CommandID;constraint:OnDelete:CASCADE;"`
	PathID     *uint
	Path       *Path `gorm:"foreignKey:PathID"`
	SessionID  *uint
	Session    *Session `gorm:"foreignKey:SessionID"`
}
