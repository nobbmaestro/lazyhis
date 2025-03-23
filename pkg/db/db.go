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

func New() (*gorm.DB, error) {
	db, err := createDatabase()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		model.History{},
		model.Command{},
		model.Session{},
		model.Path{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabase() (*gorm.DB, error) {
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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}
