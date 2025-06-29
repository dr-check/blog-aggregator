package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	configPath := filepath.Join(homeDir, configFileName)
	return configPath, nil
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read from filepath: %w", err)
	}
	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
