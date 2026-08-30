package formatters_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/formatters"
)

func ptr[T any](v T) *T { return &v }

// identityOptions returns format options that produce predictable, raw string output.
func identityOptions() formatters.FormatOptions {
	return formatters.FormatOptions{
		Command:    func(s string) string { return s },
		ExecutedAt: func(t time.Time) string { return "T" },
		ExecutedIn: func(i int) string { return "D" },
		ExitCode:   func(i int) string { return "E" },
		Id:         func(u uint) string { return "I" },
		Path:       func(s string) string { return s },
		Session:    func(s string) string { return s },
	}
}

func makeHistory(cmd string, exitCode int, path string, session string) model.History {
	return model.History{
		ExitCode:   ptr(exitCode),
		ExecutedIn: ptr(100),
		Command:    &model.Command{Command: cmd},
		Path:       &model.Path{Path: path},
		Session:    &model.Session{Session: session},
	}
}

func TestHistoryToFormatString(t *testing.T) {
	h := makeHistory("git status", 0, "/home/user/project", "session-a")

	cases := []struct {
		description string
		format      string
		want        string
	}{
		{
			"single command placeholder",
			"{COMMAND}",
			"git status",
		},
		{
			"single path placeholder",
			"{PATH}",
			"/home/user/project",
		},
		{
			"single session placeholder",
			"{SESSION}",
			"session-a",
		},
		{
			"multiple placeholders",
			"{COMMAND} {PATH}",
			"git status /home/user/project",
		},
		{
			"static text without placeholders",
			"no placeholders here",
			"no placeholders here",
		},
		{
			"unknown placeholder is left as-is",
			"{UNKNOWN}",
			"{UNKNOWN}",
		},
		{
			"mixed known and unknown placeholders",
			"{COMMAND} {UNKNOWN}",
			"git status {UNKNOWN}",
		},
		{
			"opaque placeholders use formatter options",
			"{EXIT_CODE} {EXECUTED_IN} {EXECUTED_AT} {ID}",
			"E D T I",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			f := formatters.NewFormatter(
				formatters.WithFormat(tc.format),
				formatters.WithOptions(identityOptions()),
			)

			got := f.HistoryToFormatString([]model.History{h})
			if len(got) != 1 {
				t.Fatalf("expected 1 result, got %d", len(got))
			}

			if got[0] != tc.want {
				t.Errorf("got %q, want %q", got[0], tc.want)
			}
		})
	}
}

func TestHistoryToFormatString_NilFields(t *testing.T) {
	// History with no associations set — nil pointer fields produce empty strings.
	h := model.History{}

	f := formatters.NewFormatter(
		formatters.WithFormat("{COMMAND} {PATH} {SESSION} {EXIT_CODE} {EXECUTED_IN}"),
		formatters.WithOptions(identityOptions()),
	)

	got := f.HistoryToFormatString([]model.History{h})

	if strings.TrimSpace(got[0]) != "" || strings.Count(got[0], " ") != 4 {
		t.Errorf("nil fields: got %q, want five space-separated empty strings", got[0])
	}
}

func TestHistoryToFormatString_MultipleRecords(t *testing.T) {
	records := []model.History{
		makeHistory("git status", 0, "/a", "s1"),
		makeHistory("ls -la", 1, "/b", "s2"),
	}

	f := formatters.NewFormatter(
		formatters.WithFormat("{COMMAND}"),
		formatters.WithOptions(identityOptions()),
	)

	got := f.HistoryToFormatString(records)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}

	if got[0] != "git status" || got[1] != "ls -la" {
		t.Errorf("got %v, want [\"git status\" \"ls -la\"]", got)
	}
}

func TestHistoryToTableRows(t *testing.T) {
	h := makeHistory("git status", 0, "/home/user/project", "session-a")

	columns := []config.Column{
		config.ColumnCommand,
		config.ColumnPath,
		config.ColumnSession,
	}
	f := formatters.NewFormatter(
		formatters.WithColumns(columns),
		formatters.WithOptions(identityOptions()),
	)

	rows := f.HistoryToTableRows([]model.History{h})
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	if len(row) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(row))
	}

	if row[0] != "git status" {
		t.Errorf("column 0 (COMMAND): got %q, want %q", row[0], "git status")
	}
	if row[1] != "/home/user/project" {
		t.Errorf("column 1 (PATH): got %q, want %q", row[1], "/home/user/project")
	}
	if row[2] != "session-a" {
		t.Errorf("column 2 (SESSION): got %q, want %q", row[2], "session-a")
	}
}

func TestHistoryToTableRows_ColumnOrder(t *testing.T) {
	// Verify that reversing column order in the config reverses the row output.
	h := makeHistory("git status", 0, "/home/user/project", "session-a")

	forward := formatters.NewFormatter(
		formatters.WithColumns(
			[]config.Column{config.ColumnCommand, config.ColumnPath},
		),
		formatters.WithOptions(identityOptions()),
	)
	reversed := formatters.NewFormatter(
		formatters.WithColumns(
			[]config.Column{config.ColumnPath, config.ColumnCommand},
		),
		formatters.WithOptions(identityOptions()),
	)

	fwd := forward.HistoryToTableRows([]model.History{h})[0]
	rev := reversed.HistoryToTableRows([]model.History{h})[0]

	if fwd[0] != rev[1] || fwd[1] != rev[0] {
		t.Errorf("column order not respected: forward=%v, reversed=%v", fwd, rev)
	}
}
