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

type UserConfig struct {
	Db  DbConfig  `yaml:"db"`
	Gui GuiConfig `yaml:"gui"`
	Os  OsConfig  `yaml:"os"`
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
}

type OsConfig struct {
	// Command for retrieving current session
	FetchCurrentSessionCmd string `yaml:"fetchCurrentSessionCmd"`
}
