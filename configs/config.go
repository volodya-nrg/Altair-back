package configs

import (
	"altair/pkg/manager"
	"encoding/json"
	"fmt"
	"os"
)

// Config - структура гл. конфига
type Config struct {
	Mode          string          `json:"mode"`
	TokenPassword string          `json:"tokenPassword"` // секретный ключ на выдаваемый токен
	AdminEmail    string          `json:"adminEmail"`
	Domain        string          `json:"domain"`
	DB            DbConfig        `json:"db"`
	Mediafire     MediafireConfig `json:"mediafire"`
	Email         EmailConfig     `json:"email"`
	SMS           SMSConfig       `json:"sms"`
	Socials       SocialsConfig   `json:"socials"`
}

// Cfg - гл. переменная конфига
var Cfg Config

// Load - ф-ия загрузки данных для конфига
func Load(configPath string) error {
	if configPath == "" {
		return manager.ErrConfigPath
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
