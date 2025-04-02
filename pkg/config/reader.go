package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadUserConfig(path string) *UserConfig {
	cfg := GetDefaultUserConfig()

	userConfig, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer userConfig.Close()

	err = yaml.NewDecoder(userConfig).Decode(cfg)
	if err != nil {
		return cfg
	}

	return cfg
}
