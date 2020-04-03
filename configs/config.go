package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var Cfg Config

type Config struct {
	AppMode   string
	DB        DbConfig
	Mediafire MediafireConfig
}

func Load(configPath string) error {
	if configPath == "" {
		return errors.New("empty config path")
	}

	file, _ := os.Open(configPath)
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&Cfg); err != nil {
		cwd, _ := os.Getwd()
		sErr := fmt.Errorf(
			"cannot load config from file: %s, cwd: %s, config path: %s",
			err.Error(),
			cwd,
			configPath,
		)
		return sErr
	}

	return nil
}
