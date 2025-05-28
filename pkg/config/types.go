package config

import path "github.com/nobbmaestro/lazyhis/pkg/config/parsers"

type Column string

const (
	ColumnCommand    Column = "COMMAND"
	ColumnExecutedAt Column = "EXECUTED_AT"
	ColumnExecutedIn Column = "EXECUTED_IN"
	ColumnExitCode   Column = "EXIT_CODE"
	ColumnID         Column = "ID"
	ColumnPath       Column = "PATH"
	ColumnSession    Column = "SESSION"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARNING"
	LevelError LogLevel = "ERROR"
)

type FilterMode string

const (
	NoFilter             FilterMode = "NO_FILTER"
	SuccessFilter        FilterMode = "SUCCESS_FILTER"
	WorkdirFilter        FilterMode = "WORKDIR_FILTER"
	SessionFilter        FilterMode = "SESSION_FILTER"
	UniqueFilter         FilterMode = "UNIQUE_FILTER"
	WorkdirSessionFilter FilterMode = "WORKDIR_SESSION_FILTER"
)

type UserConfig struct {
	Db  DbConfig  `yaml:"db"`
	Gui GuiConfig `yaml:"gui"`
	Os  OsConfig  `yaml:"os"`
	Log LogConfig `yaml:"log"`
}

type GuiConfig struct {
	// Option for hiding column labels
	ShowColumnLabels bool `yaml:"showColumnLabels"`
	// List of columns to render
	ColumnLayout []Column `yaml:"columnLayout"`
	// List of filter modes to cycle through
	CyclicFilterModes []FilterMode `yaml:"cyclicFilterModes"`
	// List of persistent filter modes
	PersistentFilterModes []FilterMode `yaml:"persistentFilterModes"`
	// Gui Theme
	Theme GuiTheme `yaml:"theme"`
	// Gui Keys
	Keys GuiKeys `yaml:"keys"`
}

type DbConfig struct {
	// List of excluded commands
	ExcludeCommands []string `yaml:"excludeCommands"`
	// Ignore commands starting with this prefix
	ExcludePrefix string `yaml:"excludePrefix"`
}

type OsConfig struct {
	// Command for retrieving current session
	FetchCurrentSessionCmd string `yaml:"fetchCurrentSessionCmd"`
}

type LogConfig struct {
	// Option for enabling logging
	LogEnabled bool `yaml:"logEnabled"`
	// Option for configuring log level
	LogLevel LogLevel `yaml:"logLevel"`
	// Path to the log file
	LogFile path.Path `yaml:"logFile"`
}

type GuiTheme struct {
	// Shared colors
	BorderColor string `yaml:"borderColor"`
	// Table colors
	TableCursorBgColor string `yaml:"tableCursorBgColor"`
	TableCursorFgColor string `yaml:"tableCursorFgColor"`
	TableLabelsFgColor string `yaml:"tableLabelsFgColor"`
	// Query input colors
	FilterFgColor string `yaml:"filterFgColor"`
	InputFgColor  string `yaml:"inputFgColor"`
	// Footer colors
	HelpAccentColor string `yaml:"helpAccentColor"`
	HelpFgColor     string `yaml:"helpFgColor"`
	VersionFgColor  string `yaml:"versionFgColor"`
}

type GuiKeys struct {
	AcceptSelected  []string `yaml:"acceptSelected"`
	PrefillSelected []string `yaml:"prefillSelected"`
	DeleteSelected  []string `yaml:"deleteSelected"`
	CopySelected    []string `yaml:"copySelected"`
	NextFilter      []string `yaml:"nextFilter"`
	PrevFilter      []string `yaml:"prevFilter"`
	JumpDown        []string `yaml:"jumpDown"`
	JumpUp          []string `yaml:"jumpUp"`
	MoveDown        []string `yaml:"moveDown"`
	MoveUp          []string `yaml:"moveUp"`
	Quit            []string `yaml:"quit"`
	ShowHelp        []string `yaml:"showHelp"`
}
