package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	TerminalOptions []string `json:"terminal_options"`
	AgentOptions    []string `json:"agent_options"`
}

var defaultConfig = Config{
	TerminalOptions: []string{"shell", "nextui"},
	AgentOptions:    []string{"claude", "opencode", "crush", "codex", "none"},
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".vinw-workspace")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

func loadConfig() (Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return defaultConfig, nil
	}

	configFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := saveConfig(defaultConfig); err != nil {
			return defaultConfig, nil
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return defaultConfig, nil
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultConfig, nil
	}

	return config, nil
}

func saveConfig(config Config) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}
