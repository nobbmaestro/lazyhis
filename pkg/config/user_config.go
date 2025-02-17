package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UserConfig struct {
	// Config relating to history storage
	Db DbConfig `yaml:"db"`
}

type DbConfig struct {
	// List of excluded commands
	ExcludeCommands []string `yaml:"excludeCommands"`
}

var confPath string

func init() {
	confPath = filepath.Join(os.Getenv("HOME"), ".config", "lazyhis", "lazyhis.yml")
}

func LoadUserConfig() (*UserConfig, error) {
	cfg := GetDefaultUserConfig()

	userConfig, err := os.Open(confPath)
	if err != nil {
		fmt.Println("User config file does not exist:", confPath)
		return cfg, nil
	}
	defer userConfig.Close()

	err = yaml.NewDecoder(userConfig).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetUserConfigPath() string {
	return confPath
}

func GetDefaultUserConfig() *UserConfig {
	return &UserConfig{
		Db: DbConfig{
			ExcludeCommands: []string{},
		},
	}
}
