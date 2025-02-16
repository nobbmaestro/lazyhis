package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbPath string

func init() {
	dbPath = filepath.Join(os.Getenv("HOME"), ".lazyhis.db")
}

func CreateDatabase() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: newLogger})
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
