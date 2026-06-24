package repository

import (
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"gorm.io/gorm"
)

type HistoryRepository struct {
	*GenericRepository[model.History]
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{
		GenericRepository: &GenericRepository[model.History]{db: db},
	}
}

func (r *HistoryRepository) Get(record *model.History) (*model.History, error) {
	var result model.History

	err := r.db.
		Preload("Command").
		Preload("Path").
		Preload("Session").
		Where(record).
		First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *HistoryRepository) GetByID(id uint) (*model.History, error) {
	var result model.History

	err := r.db.
		Preload("Command").
		Preload("Path").
		Preload("Session").
		Where("id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *HistoryRepository) GetAll() ([]model.History, error) {
	var records []model.History

	err := r.db.
		Preload("Command").
		Preload("Path").
		Preload("Session").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *HistoryRepository) GetLast() (model.History, error) {
	var record model.History

	err := r.db.
		Where("exit_code != ?", -1). // skip pending (not yet terminated) commands
		Preload("Command").
		Preload("Path").
		Preload("Session").
		Last(&record).Error
	if err != nil {
		return record, err
	}

	return record, nil
}

func (r *HistoryRepository) QueryHistory(
	keywords []string,
	exitCode int,
	path string,
	session string,
	limit int,
	offset int,
	unique bool,
) ([]model.History, error) {
	var histories []model.History

	query := applyHistoryFilters(
		r.db.Model(&model.History{}),
		r.db,
		keywords,
		path,
		session,
		exitCode,
	)

	if unique {
		subQuery := applyHistoryFilters(
			r.db.Model(&model.History{}).Select("MAX(id) as id").Group("command_id"),
			r.db, keywords, path, session, exitCode,
		)
		query = query.Where("id IN (?)", subQuery)
	}

	if limit != -1 {
		query = query.Limit(limit)
	}

	if offset != -1 {
		query = query.Offset(offset)
	}

	err := query.
		Order("id DESC").
		Preload("Command").
		Preload("Path").
		Preload("Session").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	return histories, nil
}

func applyHistoryFilters(
	q *gorm.DB,
	db *gorm.DB,
	keywords []string,
	path, session string,
	exitCode int,
) *gorm.DB {
	if len(keywords) > 0 {
		kq := db.Model(&model.Command{}).Select("id")
		for _, k := range keywords {
			kq = kq.Where("command LIKE ?", "%"+k+"%")
		}
		q = q.Where("command_id IN (?)", kq)
	}

	if path != "" {
		q = q.Where(
			"path_id IN (?)",
			db.Model(&model.Path{}).Select("id").Where("path LIKE ?", path),
		)
	}

	if session != "" {
		q = q.Where(
			"session_id IN (?)",
			db.Model(&model.Session{}).Select("id").Where("session LIKE ?", session),
		)
	}

	if exitCode != -1 {
		q = q.Where("exit_code = ?", exitCode)
	}

	return q
}
