package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Option func(db *gorm.DB) error

func New(path string, opts ...Option) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err := opt(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func DefaultLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{LogLevel: logger.Silent},
	)
}

func WithLogger(l logger.Interface) Option {
	return func(db *gorm.DB) error {
		db.Config.Logger = l
		return nil
	}
}

func WithForeignKeysOn() Option {
	return func(db *gorm.DB) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			return err
		}
		return nil
	}
}

func WithAutoMigrate(models ...interface{}) Option {
	return func(db *gorm.DB) error {
		return db.AutoMigrate(models...)
	}
}
