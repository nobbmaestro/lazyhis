package config

type Column string

const (
	ColumnCommand     Column = "COMMAND"
	ColumnExecutedAt  Column = "EXECUTED_AT"
	ColumnExecutedIn  Column = "EXECUTED_IN"
	ColumnExitCode    Column = "EXIT_CODE"
	ColumnPath        Column = "PATH"
	ColumnTmuxSession Column = "TMUX_SESSION"
)

type UserConfig struct {
	Db  DbConfig  `yaml:"db"`
	Gui GuiConfig `yaml:"gui"`
}

type GuiConfig struct {
	ColumnLayout []Column `yaml:"columnLayout"`
}

type DbConfig struct {
	// List of excluded commands
	ExcludeCommands []string `yaml:"excludeCommands"`
}
