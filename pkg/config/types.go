package config

type UserConfig struct {
	// Config relating to history storage
	Db DbConfig `yaml:"db"`
}

type DbConfig struct {
	// List of excluded commands
	ExcludeCommands []string `yaml:"excludeCommands"`
}
