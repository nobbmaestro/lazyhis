package db

import (
	"os"
	"path/filepath"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbPath = filepath.Join(os.Getenv("HOME"), ".lazyhis.db")

func CreateDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDatabaseConnection() (*gorm.DB, error) {
	db, err := CreateDatabase()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		model.History{},
		model.Command{},
		model.TmuxSession{},
		model.Path{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
