package service_test

import (
	"testing"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ptr[T any](v T) *T { return &v }

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Command{},
		&model.Path{},
		&model.Session{},
		&model.History{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func newService(db *gorm.DB) *service.HistoryService {
	return service.NewHistoryService(&service.RepositoryProvider{
		CommandRepo: repository.NewCommandRepository(db),
		HistoryRepo: repository.NewHistoryRepository(db),
		PathRepo:    repository.NewPathRepository(db),
		SessionRepo: repository.NewSessionRepository(db),
	})
}

func TestAddHistory(t *testing.T) {
	t.Run("empty command returns error", func(t *testing.T) {
		svc := newService(setupDB(t))

		_, err := svc.AddHistory([]string{}, nil, nil, nil, nil)

		if err == nil {
			t.Error("expected error for empty command, got nil")
		}
	})

	t.Run("creates history with all fields", func(t *testing.T) {
		svc := newService(setupDB(t))

		h, err := svc.AddHistory(
			[]string{"git", "status"},
			ptr(0),
			ptr(42),
			ptr("/home/user/project"),
			ptr("session-a"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if h == nil {
			t.Fatal("expected non-nil history")
		}
		if h.CommandID == nil {
			t.Error("CommandID should be set")
		}
		if h.PathID == nil {
			t.Error("PathID should be set")
		}
		if h.SessionID == nil {
			t.Error("SessionID should be set")
		}
		if h.ExitCode == nil || *h.ExitCode != 0 {
			t.Errorf("ExitCode: got %v, want 0", h.ExitCode)
		}
		if h.ExecutedIn == nil || *h.ExecutedIn != 42 {
			t.Errorf("ExecutedIn: got %v, want 42", h.ExecutedIn)
		}
	})

	t.Run("nil path leaves PathID unset", func(t *testing.T) {
		svc := newService(setupDB(t))

		h, err := svc.AddHistory([]string{"ls"}, ptr(0), nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if h.PathID != nil {
			t.Errorf("expected nil PathID, got %v", h.PathID)
		}
	})

	t.Run("nil session leaves SessionID unset", func(t *testing.T) {
		svc := newService(setupDB(t))

		h, err := svc.AddHistory([]string{"ls"}, ptr(0), nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if h.SessionID != nil {
			t.Errorf("expected nil SessionID, got %v", h.SessionID)
		}
	})

	t.Run("duplicate command reuses existing command row", func(t *testing.T) {
		db := setupDB(t)
		svc := newService(db)

		for i := 0; i < 2; i++ {
			if _, err := svc.AddHistory([]string{"git", "status"}, ptr(0), nil, nil, nil); err != nil {
				t.Fatalf("call %d: %v", i+1, err)
			}
		}

		var commands []model.Command
		db.Find(&commands)
		if len(commands) != 1 {
			t.Errorf("expected 1 command row, got %d", len(commands))
		}

		var histories []model.History
		db.Find(&histories)
		if len(histories) != 2 {
			t.Errorf("expected 2 history rows, got %d", len(histories))
		}
	})

	t.Run("duplicate path reuses existing path row", func(t *testing.T) {
		db := setupDB(t)
		svc := newService(db)

		for _, cmd := range []string{"ls", "pwd"} {
			if _, err := svc.AddHistory([]string{cmd}, nil, nil, ptr("/home/user"), nil); err != nil {
				t.Fatalf("AddHistory(%q): %v", cmd, err)
			}
		}

		var paths []model.Path
		db.Find(&paths)
		if len(paths) != 1 {
			t.Errorf("expected 1 path row, got %d", len(paths))
		}
	})

	t.Run("duplicate session reuses existing session row", func(t *testing.T) {
		db := setupDB(t)
		svc := newService(db)

		for _, cmd := range []string{"ls", "pwd"} {
			if _, err := svc.AddHistory([]string{cmd}, nil, nil, nil, ptr("session-a")); err != nil {
				t.Fatalf("AddHistory(%q): %v", cmd, err)
			}
		}

		var sessions []model.Session
		db.Find(&sessions)
		if len(sessions) != 1 {
			t.Errorf("expected 1 session row, got %d", len(sessions))
		}
	})

	t.Run("distinct commands each get their own history row", func(t *testing.T) {
		db := setupDB(t)
		svc := newService(db)

		commands := [][]string{
			{"git", "status"},
			{"ls", "-la"},
			{"docker", "ps"},
		}
		for _, cmd := range commands {
			if _, err := svc.AddHistory(cmd, ptr(0), nil, nil, nil); err != nil {
				t.Fatalf("AddHistory(%v): %v", cmd, err)
			}
		}

		var histories []model.History
		db.Find(&histories)
		if len(histories) != 3 {
			t.Errorf("expected 3 history rows, got %d", len(histories))
		}
	})
}
