package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	MachineID string   `json:"machine_id"`
	Paths     []string `json:"paths"`
	GistToken string   `json:"gist_token"`
	GistID    string   `json:"gist_id"`
}

func LoadConfig(filepath string) (*AppConfig, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
