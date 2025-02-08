package service

import (
	"errors"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
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

func (s *HistoryService) AddHistoryIfUnique(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	tmuxSession *string,
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
	)
}

func (s *HistoryService) AddHistory(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	tmuxSession *string,
) (*model.History, error) {
	var (
		commandID     *uint
		pathID        *uint
		tmuxSessionID *uint
	)

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
