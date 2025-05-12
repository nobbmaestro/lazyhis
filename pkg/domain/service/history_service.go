package service

import (
	"errors"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
)

type RepositoryProvider struct {
	CommandRepo *repository.CommandRepository
	HistoryRepo *repository.HistoryRepository
	PathRepo    *repository.PathRepository
	SessionRepo *repository.SessionRepository
}

type HistoryService struct {
	repos *RepositoryProvider
}

func NewHistoryService(
	repos *RepositoryProvider,
) *HistoryService {
	return &HistoryService{repos: repos}
}

func (s *HistoryService) SearchHistory(
	keywords []string,
	exitCode int,
	path string,
	session string,
	limit int,
	offset int,
	unique bool,
) ([]model.History, error) {
	results, err := s.repos.HistoryRepo.QueryHistory(
		keywords,
		exitCode,
		path,
		session,
		limit,
		offset,
		unique,
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *HistoryService) AddHistory(
	command []string,
	exitCode *int,
	executedIn *int,
	path *string,
	session *string,
) (*model.History, error) {
	var (
		commandID *uint
		pathID    *uint
		sessionID *uint
	)

	commandRecord, err := s.AddCommand(command)
	if err != nil {
		return nil, err
	}
	if commandRecord != nil {
		commandID = &commandRecord.ID
	}

	sessionRecord, err := s.AddSession(session)
	if err != nil {
		return nil, err
	}
	if sessionRecord != nil {
		sessionID = &sessionRecord.ID
	}

	pathRecord, err := s.AddPath(path)
	if err != nil {
		return nil, err
	}
	if pathRecord != nil {
		pathID = &pathRecord.ID
	}

	history := &model.History{
		ExitCode:   exitCode,
		ExecutedIn: executedIn,
		CommandID:  commandID,
		PathID:     pathID,
		SessionID:  sessionID,
	}

	return s.repos.HistoryRepo.Create(history)
}

func (s *HistoryService) EditHistory(
	historyID int,
	exitCode *int,
	executedIn *int,
	path *string,
	session *string,
) (*model.History, error) {
	history, err := s.repos.HistoryRepo.GetByID(uint(historyID))
	if err != nil {
		return nil, err
	}

	if exitCode != nil {
		history.ExitCode = exitCode
	}

	if executedIn != nil {
		history.ExecutedIn = executedIn
	}

	if path != nil {
		pathRecord, err := s.AddPath(path)
		if err != nil {
			return nil, err
		}
		if pathRecord != nil {
			history.Path = pathRecord
		}
	}

	if session != nil {
		sessionRecord, err := s.AddSession(session)
		if err != nil {
			return nil, err
		}
		if sessionRecord != nil {
			history.Session = sessionRecord
		}
	}

	return s.repos.HistoryRepo.Update(history)
}

func (s *HistoryService) DeleteCommand(record *model.Command) error {
	_, err := s.repos.CommandRepo.Delete(record)
	return err
}

func (s *HistoryService) AddCommand(command []string) (*model.Command, error) {
	if len(command) == 0 {
		return nil, errors.New("Command cannot be empty")
	}
	return s.repos.CommandRepo.GetOrCreate(
		&model.Command{Command: strings.Join(command, " ")},
	)
}

func (s *HistoryService) AddSession(session *string) (*model.Session, error) {
	if session == nil {
		return nil, nil
	}
	return s.repos.SessionRepo.GetOrCreate(&model.Session{Session: *session})
}

func (s *HistoryService) AddPath(path *string) (*model.Path, error) {
	if path == nil {
		return nil, nil
	}
	return s.repos.PathRepo.GetOrCreate(&model.Path{Path: *path})
}

func (s *HistoryService) GetAllHistory() ([]model.History, error) {
	return s.repos.HistoryRepo.GetAll()
}

func (s *HistoryService) GetAllCommands() ([]model.Command, error) {
	return s.repos.CommandRepo.GetAll()
}

func (s *HistoryService) GetAllSessions() ([]model.Session, error) {
	return s.repos.SessionRepo.GetAll()
}

func (s *HistoryService) GetAllPaths() ([]model.Path, error) {
	return s.repos.PathRepo.GetAll()
}

func (s *HistoryService) GetLastHistory() (model.History, error) {
	return s.repos.HistoryRepo.GetLast()
}

func (s HistoryService) CommandExists(command []string) bool {
	return s.repos.CommandRepo.Exists(&model.Command{Command: strings.Join(command, " ")})
}
