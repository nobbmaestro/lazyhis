package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
	"gorm.io/gorm"
)

type HistoryService struct {
	historyRepo *repository.HistoryRepository
	commandRepo *repository.CommandRepository
	pathRepo    *repository.PathRepository
	tmuxRepo    *repository.TmuxSessionRepository
}

func NewHistoryService(
	historyRepo *repository.HistoryRepository,
	commandRepo *repository.CommandRepository,
	pathRepo *repository.PathRepository,
	tmuxRepo *repository.TmuxSessionRepository,
) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		commandRepo: commandRepo,
		pathRepo:    pathRepo,
		tmuxRepo:    tmuxRepo,
	}
}

func (s *HistoryService) SearchHistory(
	keywords []string,
	exitCode int,
	path string,
	tmuxSession string,
	limit int,
	offset int,
) ([]model.History, error) {
	results, err := s.historyRepo.QueryHistory(
		keywords,
		exitCode,
		path,
		tmuxSession,
		limit,
		offset,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return results, nil
}

func (s *HistoryService) AddHistoryIfUnique(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	tmuxSession *string,
	excludeCommands *[]string,
) (*model.History, error) {
	commandRecord, err := s.commandRepo.Get(
		&model.Command{Command: strings.Join(command, " ")},
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if commandRecord != nil {
		return nil, nil
	}

	return s.AddHistory(
		command,
		exitCode,
		executedIn,
		path,
		tmuxSession,
		excludeCommands,
	)
}

func (s *HistoryService) AddHistory(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	tmuxSession *string,
	excludeCommands *[]string,
) (*model.History, error) {
	var (
		commandID     *uint
		pathID        *uint
		tmuxSessionID *uint
	)

	if excludeCommands != nil &&
		utils.IsExcluded(strings.Join(command, " "), *excludeCommands) {
		return nil, nil
	}

	commandRecord, err := s.AddCommand(command)
	if err != nil {
		return nil, err
	}
	if commandRecord != nil {
		commandID = &commandRecord.ID
	}

	tmuxSessionRecord, err := s.AddTmuxSession(tmuxSession)
	if err != nil {
		return nil, err
	}
	if tmuxSessionRecord != nil {
		tmuxSessionID = &tmuxSessionRecord.ID
	}

	pathRecord, err := s.AddPath(path)
	if err != nil {
		return nil, err
	}
	if pathRecord != nil {
		pathID = &pathRecord.ID
	}

	history := &model.History{
		ExitCode:      exitCode,
		ExecutedIn:    executedIn,
		CommandID:     commandID,
		PathID:        pathID,
		TmuxSessionID: tmuxSessionID,
	}

	return s.historyRepo.Create(history)
}

func (s *HistoryService) PruneHistory(excludeCommands []string) error {
	records, err := s.GetAllCommands()
	if err != nil {
		return err
	}

	for _, record := range records {
		if utils.IsExcluded(record.Command, excludeCommands) {
			fmt.Println("Prune:", record.Command)

			_, err := s.commandRepo.Delete(&record)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *HistoryService) AddCommand(command []string) (*model.Command, error) {
	if len(command) == 0 {
		return nil, errors.New("Command cannot be empty")
	}
	return s.commandRepo.GetOrCreate(
		&model.Command{Command: strings.Join(command, " ")},
	)
}

func (s *HistoryService) AddTmuxSession(session *string) (*model.TmuxSession, error) {
	if session == nil {
		return nil, nil
	}
	return s.tmuxRepo.GetOrCreate(&model.TmuxSession{Session: *session})
}

func (s *HistoryService) AddPath(path *string) (*model.Path, error) {
	if path == nil {
		return nil, nil
	}
	return s.pathRepo.GetOrCreate(&model.Path{Path: *path})
}

func (s *HistoryService) GetAllHistory() ([]model.History, error) {
	return s.historyRepo.GetAll()
}

func (s *HistoryService) GetAllCommands() ([]model.Command, error) {
	return s.commandRepo.GetAll()
}

func (s *HistoryService) GetAllTmuxSessions() ([]model.TmuxSession, error) {
	return s.tmuxRepo.GetAll()
}

func (s *HistoryService) GetAllPaths() ([]model.Path, error) {
	return s.pathRepo.GetAll()
}
