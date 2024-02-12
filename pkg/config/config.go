package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
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

func GetConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Error(err, "Cannot open config.")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Error(err, "Cannot read config.")
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Error(err, "Cannot parse config.")
	}

	return config
}
