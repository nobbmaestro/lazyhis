package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadUserConfig(path string) *UserConfig {
	cfg := GetDefaultUserConfig()

	userConfig, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer func() {
		if err := userConfig.Close(); err != nil {
			log.Printf("failed to close file: %v\n", err)
		}
	}()

	err = yaml.NewDecoder(userConfig).Decode(cfg)
	if err != nil {
		return cfg
	}

	return cfg
}
