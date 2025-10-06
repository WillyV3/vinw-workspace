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

// WorkspaceCommand represents a custom command to run in terminal pane
type WorkspaceCommand struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// WorkspaceConfig stores custom commands for workspaces
type WorkspaceConfig struct {
	Commands []WorkspaceCommand `json:"commands"`
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

// getVinwDir returns the ~/.vinw directory (shared with vinw)
func getVinwDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	vinwDir := filepath.Join(homeDir, ".vinw")
	if err := os.MkdirAll(vinwDir, 0755); err != nil {
		return "", err
	}
	return vinwDir, nil
}

// loadWorkspaceCommands loads custom commands from ~/.vinw/workspace.conf
func loadWorkspaceCommands() ([]WorkspaceCommand, error) {
	vinwDir, err := getVinwDir()
	if err != nil {
		return []WorkspaceCommand{}, err
	}

	confFile := filepath.Join(vinwDir, "workspace.conf")

	// Return empty list if file doesn't exist yet
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return []WorkspaceCommand{}, nil
	}

	data, err := os.ReadFile(confFile)
	if err != nil {
		return []WorkspaceCommand{}, err
	}

	var config WorkspaceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return []WorkspaceCommand{}, err
	}

	return config.Commands, nil
}

// saveWorkspaceCommands saves custom commands to ~/.vinw/workspace.conf
func saveWorkspaceCommands(commands []WorkspaceCommand) error {
	vinwDir, err := getVinwDir()
	if err != nil {
		return err
	}

	confFile := filepath.Join(vinwDir, "workspace.conf")

	config := WorkspaceConfig{Commands: commands}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(confFile, data, 0644)
}
