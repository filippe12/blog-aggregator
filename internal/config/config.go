package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUserName = username
	write(cfg)
}

func Read() Config {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Open(filepath.Join(homedir, configFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	config := Config{}
	if err = decoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	return config
}

func write(cfg *Config) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.OpenFile(
		filepath.Join(homedir, configFileName),
		os.O_WRONLY|os.O_TRUNC,
		0600,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	if err = encoder.Encode(cfg); err != nil {
		log.Fatal(err)
	}
}
