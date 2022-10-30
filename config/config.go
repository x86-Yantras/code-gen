package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/x86-Yantras/code-gen/internal/constants"
)

type Config struct {
	PackageManager string
	ServiceDir     string
	ReadmeFile     string
}

func New(configType string) (*Config, error) {
	configFile := fmt.Sprintf("%s/%s%s", constants.ConfigDir, configType, "_config.json")

	configByte, err := os.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = json.Unmarshal(configByte, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
