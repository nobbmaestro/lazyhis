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

	query := r.db.Model(&model.History{})

	if len(keywords) > 0 {
		keywordQuery := r.db.Model(&model.Command{}).Select("id")
		for _, keyword := range keywords {
			keywordQuery = keywordQuery.Where("command LIKE ?", "%"+keyword+"%")
		}
		query = query.Where("command_id IN (?)", keywordQuery)
	}

	if path != "" {
		query = query.Where("path_id IN (?)", r.db.Model(&model.Path{}).
			Select("id").
			Where("path LIKE ?", path))
	}

	if session != "" {
		query = query.Where("session_id IN (?)", r.db.Model(&model.Session{}).
			Select("id").
			Where("session LIKE ?", session))
	}

	if exitCode != -1 {
		query = query.Where("exit_code = ?", exitCode)
	}

	if limit != -1 {
		query = query.Limit(limit)
	}

	if offset != -1 {
		query = query.Offset(offset)
	}

	if unique {
		subQuery := r.db.Model(&model.History{}).
			Select("MAX(id) as id").
			Group("command_id")
		query = query.Where("id IN (?)", subQuery)
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
