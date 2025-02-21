package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type RepositoryProvider struct {
	CommandRepo *repository.CommandRepository
	HistoryRepo *repository.HistoryRepository
	PathRepo    *repository.PathRepository
	TmuxRepo    *repository.TmuxSessionRepository
}

type HistoryService struct {
	repos *RepositoryProvider
}

func NewHistoryService(repos *RepositoryProvider) *HistoryService {
	return &HistoryService{repos: repos}
}

func (s *HistoryService) SearchHistory(
	keywords []string,
	exitCode int,
	path string,
	tmuxSession string,
	limit int,
	offset int,
) ([]model.History, error) {
	results, err := s.repos.HistoryRepo.QueryHistory(
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
	if s.repos.CommandRepo.Exists(&model.Command{Command: strings.Join(command, " ")}) {
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

	return s.repos.HistoryRepo.Create(history)
}

func (s *HistoryService) EditHistory(
	historyID int,
	exitCode *int,
	executedIn *int,
	path *string,
	tmuxSession *string,
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

	if tmuxSession != nil {
		tmuxSessionRecord, err := s.AddTmuxSession(tmuxSession)
		if err != nil {
			return nil, err
		}
		if tmuxSessionRecord != nil {
			history.TmuxSession = tmuxSessionRecord
		}
	}

	return s.repos.HistoryRepo.Update(history)
}

func (s *HistoryService) PruneHistory(excludeCommands []string) error {
	records, err := s.GetAllCommands()
	if err != nil {
		return err
	}

	for _, record := range records {
		if utils.IsExcluded(record.Command, excludeCommands) {
			fmt.Println("Prune:", record.Command)

			_, err := s.repos.CommandRepo.Delete(&record)
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
	return s.repos.CommandRepo.GetOrCreate(
		&model.Command{Command: strings.Join(command, " ")},
	)
}

func (s *HistoryService) AddTmuxSession(session *string) (*model.TmuxSession, error) {
	if session == nil {
		return nil, nil
	}
	return s.repos.TmuxRepo.GetOrCreate(&model.TmuxSession{Session: *session})
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

func (s *HistoryService) GetAllTmuxSessions() ([]model.TmuxSession, error) {
	return s.repos.TmuxRepo.GetAll()
}

func (s *HistoryService) GetAllPaths() ([]model.Path, error) {
	return s.repos.PathRepo.GetAll()
}
