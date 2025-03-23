package config

type Column string

const (
	ColumnCommand    Column = "COMMAND"
	ColumnExecutedAt Column = "EXECUTED_AT"
	ColumnExecutedIn Column = "EXECUTED_IN"
	ColumnExitCode   Column = "EXIT_CODE"
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

type UserConfig struct {
	Db  DbConfig  `yaml:"db"`
	Gui GuiConfig `yaml:"gui"`
	Os  OsConfig  `yaml:"os"`
	Log LogConfig `yaml:"log"`
}

type GuiConfig struct {
	// List of columns to render
	ColumnLayout []Column `yaml:"columnLayout"`
	// Option for display only unique commands
	ShowUniqueCommands bool `yaml:"showUniqueCommands"`
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
	LogFile string `yaml:"logFile"`
}
