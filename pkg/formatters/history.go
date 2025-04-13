package formatters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type FormatterFunc[T any] func(T) string

type FormatOptions struct {
	Command    FormatterFunc[string]
	ExecutedAt FormatterFunc[time.Time]
	ExecutedIn FormatterFunc[int]
	ExitCode   FormatterFunc[int]
	Id         FormatterFunc[uint]
	Path       FormatterFunc[string]
	Session    FormatterFunc[string]
}

func DefaultTuiFormatOptions() FormatOptions {
	return FormatOptions{
		Command:    func(s string) string { return s },
		ExecutedAt: func(t time.Time) string { return strconv.FormatInt(t.Unix(), 10) },
		ExecutedIn: func(i int) string { return strconv.FormatInt(int64(i), 10) },
		ExitCode:   func(i int) string { return strconv.Itoa(i) },
		Id:         func(u uint) string { return strconv.FormatUint(uint64(u), 10) },
		Path:       func(s string) string { return s },
		Session:    func(s string) string { return s },
	}
}

func DefaultGuiFormatOptions() FormatOptions {
	return FormatOptions{
		Command:    func(s string) string { return s },
		ExecutedAt: func(t time.Time) string { return utils.HumanizeTimeAgo(t) },
		ExecutedIn: func(i int) string { return utils.HumanizeDuration(int64(i)) },
		ExitCode:   func(v int) string { return fmt.Sprintf("%3d", v) },
		Id:         func(v uint) string { return strconv.FormatUint(uint64(v), 10) },
		Path:       func(t string) string { return t },
		Session:    func(t string) string { return t },
	}
}

type Formatter struct {
	columns []config.Column
	opts    FormatOptions
	format  string
}

type Option func(*Formatter)

func NewFormatter(opts ...Option) Formatter {
	f := Formatter{
		columns: []config.Column{
			config.ColumnID,
			config.ColumnExecutedAt,
			config.ColumnExecutedIn,
			config.ColumnCommand,
			config.ColumnExitCode,
			config.ColumnSession,
			config.ColumnPath,
		},
		opts: FormatOptions{
			Command:    func(s string) string { return s },
			ExecutedAt: func(t time.Time) string { return t.String() },
			ExecutedIn: func(i int) string { return strconv.Itoa(i) },
			ExitCode:   func(v int) string { return strconv.Itoa(v) },
			Id:         func(v uint) string { return strconv.FormatUint(uint64(v), 10) },
			Path:       func(t string) string { return t },
			Session:    func(t string) string { return t },
		},
		format: "",
	}

	for _, opt := range opts {
		opt(&f)
	}

	return f
}

func WithColumns(columns []config.Column) Option {
	return func(f *Formatter) {
		f.columns = columns
	}
}

func WithOptions(opts FormatOptions) Option {
	return func(f *Formatter) {
		f.opts = opts
	}
}

func WithFormat(format string) Option {
	return func(f *Formatter) {
		f.format = format
	}
}

func (f Formatter) formatHistoryByColumn(
	record model.History,
	column config.Column,
) string {
	switch column {
	case config.ColumnID:
		return f.opts.Id(record.ID)
	case config.ColumnExecutedAt:
		return f.opts.ExecutedAt(record.CreatedAt)
	case config.ColumnExecutedIn:
		if v := record.ExecutedIn; v != nil {
			return f.opts.ExecutedIn(*v)
		}
	case config.ColumnCommand:
		if v := record.Command; v != nil {
			return f.opts.Command(v.Command)
		}
	case config.ColumnExitCode:
		if v := record.ExitCode; v != nil {
			return f.opts.ExitCode(*v)
		}
	case config.ColumnSession:
		if v := record.Session; v != nil {
			return f.opts.Session(v.Session)
		}
	case config.ColumnPath:
		if v := record.Path; v != nil {
			return f.opts.Path(v.Path)
		}
	}
	return ""
}

func (f Formatter) HistoryToTableRows(records []model.History) []table.Row {
	rows := make([]table.Row, len(records))
	for i, record := range records {
		row := make([]string, len(f.columns))
		for j, column := range f.columns {
			row[j] = f.formatHistoryByColumn(record, column)
		}
		rows[i] = row
	}
	return rows
}

func findColumnsInFormat(format string) []config.Column {
	columns := []config.Column{}
	for _, column := range []config.Column{
		config.ColumnCommand,
		config.ColumnExecutedAt,
		config.ColumnExecutedIn,
		config.ColumnExitCode,
		config.ColumnID,
		config.ColumnPath,
		config.ColumnSession,
	} {
		if strings.Contains(format, "{"+string(column)+"}") {
			columns = append(columns, column)
		}
	}
	return columns
}

func (f Formatter) HistoryToFormatString(records []model.History) []string {
	rows := make([]string, len(records))
	for i, record := range records {
		rowStr := f.format
		for _, column := range findColumnsInFormat(f.format) {
			placeholder := "{" + string(column) + "}"
			replacement := f.formatHistoryByColumn(record, column)
			rowStr = strings.ReplaceAll(rowStr, placeholder, replacement)
		}
		rows[i] = rowStr
	}
	return rows
}
