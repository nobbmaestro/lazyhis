package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadUserConfig(path string) (*UserConfig, error) {
	cfg := GetDefaultUserConfig()

	userConfig, err := os.Open(path)
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
