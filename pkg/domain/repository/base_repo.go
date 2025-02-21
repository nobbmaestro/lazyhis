package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type Model interface {
	model.Command | model.Path | model.TmuxSession | model.History
}

type Repository[T Model] interface {
	Create(record *T) (*T, error)
	Update(record *T) (*T, error)
	Delete(record *T) (*T, error)
	Get(record *T) (*T, error)
	GetOrCreate(record *T) (*T, error)
	GetAll() ([]T, error)
}

type BaseRepository[T Model] struct {
	db *gorm.DB
}

func (r *BaseRepository[T]) Create(record *T) (*T, error) {
	return record, r.db.Create(record).Error
}

func (r *BaseRepository[T]) Update(record *T) (*T, error) {
	return record, r.db.Model(record).Updates(record).Error
}

func (r *BaseRepository[T]) Delete(record *T) (*T, error) {
	return record, r.db.Unscoped().Delete(record).Error // This perform HARD delete
}

func (r *BaseRepository[T]) Get(record *T) (*T, error) {
	var result T

	err := r.db.Where(record).First(&result).Error
	if nil != err {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) GetOrCreate(record *T) (*T, error) {
	err := r.db.Where(record).FirstOrCreate(record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *BaseRepository[T]) GetAll() ([]T, error) {
	var records []T

	err := r.db.Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *BaseRepository[T]) Exists(record *T) bool {
	result, err := r.Get(record)
	if err != nil {
		return false
	}
	return result != nil
}
