package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// config

type RegistryConfig struct {
	User     string
	Password string
	Api      string
	Base     string
	Folder   string
}

type Ssh struct {
	KeyName string
	KeyPath string
}

type Config struct {
	Registry      RegistryConfig
	Ssh           Ssh
	IsInitialised bool
}

func GetConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open config %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config %s: %w", path, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse config %s: %w", path, err)
	}

	return config, nil
}
