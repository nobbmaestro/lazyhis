package app

import (
	"os"
	"slices"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

type HistoryOptions struct {
	Command             []string
	ExitCode            *int
	Path                *string
	Session             *string
	ExecutedIn          *int
	MaxNumSearchResults *int
	OffsetSearchResults *int
	UniqueSearchResults *bool
	Filters             []config.FilterMode
}

func ptr[T any](v T) *T {
	return &v
}

func defaultHistoryOptions() HistoryOptions {
	return HistoryOptions{
		Command:             []string{},
		ExitCode:            ptr(-1),
		Path:                ptr(""),
		Session:             ptr(""),
		ExecutedIn:          ptr(-1),
		MaxNumSearchResults: ptr(-1),
		OffsetSearchResults: ptr(-1),
		UniqueSearchResults: ptr(false),
		Filters:             []config.FilterMode{},
	}
}

type HistoryOption func(*HistoryOptions)

func WithQuery(query []string) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.Command = query
	}
}

func WithExitCode(exitCode int) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.ExitCode = ptr(exitCode)
	}
}

func WithPath(path string) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.Path = ptr(path)
	}
}

func WithSession(session string) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.Session = ptr(session)
	}
}

func WithExecutedIn(executedIn int) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.ExecutedIn = ptr(executedIn)
	}
}

func WithMaxNumSearchResults(value int) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.MaxNumSearchResults = ptr(value)
	}
}

func WithOffsetSearchResults(value int) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.OffsetSearchResults = ptr(value)
	}
}

func WithUniqueSearchResults(value bool) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.UniqueSearchResults = ptr(value)
	}
}

func WithFilters(filters []config.FilterMode) HistoryOption {
	return func(opts *HistoryOptions) {
		opts.Filters = filters
	}
}

func applyPathFilter(filters []config.FilterMode) string {
	if slices.Contains(filters, config.WorkdirFilter) ||
		slices.Contains(filters, config.WorkdirSessionFilter) {
		if p, err := os.Getwd(); err == nil {
			return p
		}
	}
	return ""
}

func applySessionFilter(
	filters []config.FilterMode,
	sessionCmd string,
) string {
	if slices.Contains(filters, config.WorkdirSessionFilter) ||
		slices.Contains(filters, config.SessionFilter) {
		if s, err := utils.RunCommand(strings.Split(sessionCmd, " ")); err == nil {
			return s
		}
	}
	return ""
}

func applySuccessFilter(filters []config.FilterMode) int {
	if slices.Contains(filters, config.SuccessFilter) {
		return 0
	}
	return -1
}

func applyUniqueCommandFilter(filters []config.FilterMode) bool {
	return slices.Contains(filters, config.UniqueFilter)
}

func applyFilters(
	opts *HistoryOptions,
	fetchSessionString string,
) {
	for _, filter := range opts.Filters {
		switch filter {
		case config.WorkdirFilter:
			opts.Path = ptr(applyPathFilter(opts.Filters))
		case config.WorkdirSessionFilter:
			opts.Path = ptr(applyPathFilter(opts.Filters))
			opts.Session = ptr(applySessionFilter(opts.Filters, fetchSessionString))
		case config.SessionFilter:
			opts.Session = ptr(applySessionFilter(opts.Filters, fetchSessionString))
		case config.SuccessFilter:
			opts.ExitCode = ptr(applySuccessFilter(opts.Filters))
		case config.UniqueFilter:
			opts.UniqueSearchResults = ptr(applyUniqueCommandFilter(opts.Filters))
		}
	}
}
