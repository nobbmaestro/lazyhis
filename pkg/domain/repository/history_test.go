package repository_test

import (
	"testing"

	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Fixture layout (IDs assigned in insertion order):
//
//	ID | Command                  | ExitCode | Path                 | Session
//	 1 | git status               |    0     | /home/user/project   | session-a
//	 2 | ls -la                   |    0     | /home/user/other     | session-a
//	 3 | docker ps                |    1     | /home/user/project   | session-b
//	 4 | git commit -m 'wip'      |    0     | /home/user/project   | session-a
//	 5 | git status               |    0     | /home/user/other     | session-b
//	 6 | git log                  |    0     | /home/user/project   | session-a
//	 7 | git status               |    1     | /home/user/project   | session-a
//	 8 | ls -la                   |    0     | /home/user/project   | session-b

func ptr[T any](v T) *T { return &v }

func setupTestDB(t *testing.T) *gorm.DB {
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

type fixtureEntry struct {
	command  string
	exitCode int
	path     string
	session  string
}

func seedFixtures(t *testing.T, db *gorm.DB) {
	t.Helper()

	entries := []fixtureEntry{
		{"git status", 0, "/home/user/project", "session-a"},
		{"ls -la", 0, "/home/user/other", "session-a"},
		{"docker ps", 1, "/home/user/project", "session-b"},
		{"git commit -m 'wip'", 0, "/home/user/project", "session-a"},
		{"git status", 0, "/home/user/other", "session-b"},
		{"git log", 0, "/home/user/project", "session-a"},
		{"git status", 1, "/home/user/project", "session-a"},
		{"ls -la", 0, "/home/user/project", "session-b"},
	}
	for _, e := range entries {
		cmd := &model.Command{Command: e.command}
		db.Where(cmd).FirstOrCreate(cmd)

		path := &model.Path{Path: e.path}
		db.Where(path).FirstOrCreate(path)

		sess := &model.Session{Session: e.session}
		db.Where(sess).FirstOrCreate(sess)

		h := &model.History{
			ExitCode:  ptr(e.exitCode),
			CommandID: &cmd.ID,
			PathID:    &path.ID,
			SessionID: &sess.ID,
		}
		if err := db.Create(h).Error; err != nil {
			t.Fatalf("failed to seed fixture %+v: %v", e, err)
		}
	}
}

func ids(records []model.History) []uint {
	out := make([]uint, len(records))
	for i, r := range records {
		out[i] = r.ID
	}
	return out
}

func equalIDs(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestQueryHistory(t *testing.T) {
	cases := []struct {
		name     string
		keywords []string
		exitCode int
		path     string
		session  string
		limit    int
		offset   int
		unique   bool
		wantIDs  []uint // expected IDs in order (DESC)
	}{
		{
			name:     "no filters returns all ordered by id desc",
			keywords: []string{},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:     "single keyword match",
			keywords: []string{"docker"},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{3},
		},
		{
			name:     "multiple keywords are ANDed",
			keywords: []string{"git", "commit"},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{4},
		},
		{
			name:     "keyword with no match returns empty",
			keywords: []string{"nonexistent"},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{},
		},
		{
			name:     "exit code 0 filter",
			keywords: []string{},
			exitCode: 0, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{8, 6, 5, 4, 2, 1},
		},
		{
			name:     "exit code 1 filter",
			keywords: []string{},
			exitCode: 1, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{7, 3},
		},
		{
			name:     "path filter",
			keywords: []string{},
			exitCode: -1, path: "/home/user/project", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{8, 7, 6, 4, 3, 1},
		},
		{
			name:     "session filter",
			keywords: []string{},
			exitCode: -1, path: "", session: "session-a",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{7, 6, 4, 2, 1},
		},
		{
			name:     "path and session combined",
			keywords: []string{},
			exitCode: -1, path: "/home/user/project", session: "session-a",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{7, 6, 4, 1},
		},
		{
			name:     "keyword and exit code combined",
			keywords: []string{"git"},
			exitCode: 0, path: "", session: "",
			limit: -1, offset: -1, unique: false,
			wantIDs: []uint{6, 5, 4, 1},
		},
		{
			name:     "unique returns max id per command",
			keywords: []string{},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: true,
			wantIDs: []uint{8, 7, 6, 4, 3},
		},
		{
			name:     "unique with keyword filter",
			keywords: []string{"git"},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: -1, unique: true,
			wantIDs: []uint{7, 6, 4},
		},
		{
			name:     "unique with exit code filter",
			keywords: []string{},
			exitCode: 0, path: "", session: "",
			limit: -1, offset: -1, unique: true,
			// Only entries with exitCode=0 are candidates for MAX(id):
			// git status exit=0: IDs 1,5 → max=5; ls -la exit=0: IDs 2,8 → max=8
			// git commit exit=0: ID 4 → max=4; git log exit=0: ID 6 → max=6
			// docker ps has no exit=0 entry → excluded
			// Outer query also filters exit=0, so ID 1 (git status, exit=0) is
			// excluded because 5 is the max for that command — not 1.
			wantIDs: []uint{8, 6, 5, 4},
		},
		{
			name:     "limit restricts results",
			keywords: []string{},
			exitCode: -1, path: "", session: "",
			limit: 3, offset: -1, unique: false,
			wantIDs: []uint{8, 7, 6},
		},
		{
			name:     "offset skips leading results",
			keywords: []string{},
			exitCode: -1, path: "", session: "",
			limit: -1, offset: 2, unique: false,
			wantIDs: []uint{6, 5, 4, 3, 2, 1},
		},
		{
			name:     "limit and offset paginate correctly",
			keywords: []string{},
			exitCode: -1, path: "", session: "",
			limit: 2, offset: 2, unique: false,
			wantIDs: []uint{6, 5},
		},
		{
			name:     "all filters combined",
			keywords: []string{"git"},
			exitCode: 0, path: "/home/user/project", session: "session-a",
			limit: -1, offset: -1, unique: false,
			// git commands with exitCode=0 on /home/user/project in session-a:
			// ID 1 (git status), ID 4 (git commit), ID 6 (git log)
			wantIDs: []uint{6, 4, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			seedFixtures(t, db)
			repo := repository.NewHistoryRepository(db)

			got, err := repo.QueryHistory(
				tc.keywords, tc.exitCode, tc.path, tc.session,
				tc.limit, tc.offset, tc.unique,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotIDs := ids(got)
			if !equalIDs(gotIDs, tc.wantIDs) {
				t.Errorf("got IDs %v, want %v", gotIDs, tc.wantIDs)
			}
		})
	}
}

func TestGetLast(t *testing.T) {
	t.Run("returns most recent non-pending entry", func(t *testing.T) {
		db := setupTestDB(t)
		seedFixtures(t, db)

		// Insert a pending entry (exit_code = -1) that should be skipped.
		cmd := &model.Command{Command: "pending command"}
		db.FirstOrCreate(cmd)
		db.Create(&model.History{
			ExitCode:  ptr(-1),
			CommandID: &cmd.ID,
		})

		repo := repository.NewHistoryRepository(db)
		last, err := repo.GetLast()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// ID 8 is the last non-pending entry from our fixture.
		if last.ID != 8 {
			t.Errorf("got ID %d, want 8", last.ID)
		}
	})

	t.Run("returns error on empty database", func(t *testing.T) {
		db := setupTestDB(t)
		repo := repository.NewHistoryRepository(db)

		_, err := repo.GetLast()

		if err == nil {
			t.Error("expected error on empty database, got nil")
		}
	})

	t.Run("preloads associated command", func(t *testing.T) {
		db := setupTestDB(t)
		seedFixtures(t, db)

		repo := repository.NewHistoryRepository(db)
		last, err := repo.GetLast()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if last.Command == nil || last.Command.Command == "" {
			t.Error("expected Command to be preloaded")
		}
	})
}
