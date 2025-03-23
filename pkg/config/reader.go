package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var confPath string

func init() {
	confPath = filepath.Join(os.Getenv("HOME"), ".config", "lazyhis", "lazyhis.yml")
}

func ReadUserConfig() (*UserConfig, error) {
	cfg := GetDefaultUserConfig()

	userConfig, err := os.Open(confPath)
	if err != nil {
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
