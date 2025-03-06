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
		Preload("TmuxSession").
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
		Preload("TmuxSession").
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
		Preload("TmuxSession").
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
		Preload("TmuxSession").
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
	tmuxSession string,
	limit int,
	offset int,
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

	if tmuxSession != "" {
		query = query.Where("tmux_session_id IN (?)", r.db.Model(&model.TmuxSession{}).
			Select("id").
			Where("session LIKE ?", tmuxSession))
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

	err := query.
		Order("id DESC").
		Preload("Command").
		Preload("Path").
		Preload("TmuxSession").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	return histories, nil
}
