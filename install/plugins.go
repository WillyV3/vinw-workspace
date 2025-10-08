package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// IsTPMInstalled checks if TPM (Tmux Plugin Manager) is already installed
func IsTPMInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	tpmPath := filepath.Join(homeDir, ".tmux", "plugins", "tpm")
	info, err := os.Stat(tpmPath)
	return err == nil && info.IsDir()
}

// IsCatppuccinInstalled checks if Catppuccin theme is already installed
func IsCatppuccinInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	catppuccinPath := filepath.Join(homeDir, ".config", "tmux", "plugins", "catppuccin")
	info, err := os.Stat(catppuccinPath)
	return err == nil && info.IsDir()
}

// InstallTPM clones the TPM repository to ~/.tmux/plugins/tpm
func InstallTPM() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	tpmPath := filepath.Join(homeDir, ".tmux", "plugins", "tpm")

	// Check if already installed
	if IsTPMInstalled() {
		return nil // Already installed, no error
	}

	// Create plugins directory if it doesn't exist
	pluginsDir := filepath.Join(homeDir, ".tmux", "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Clone TPM repository
	cmd := exec.Command("git", "clone", "https://github.com/tmux-plugins/tpm", tpmPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone TPM: %w\n%s", err, string(output))
	}

	return nil
}

// InstallCatppuccin clones the Catppuccin theme to ~/.config/tmux/plugins/catppuccin
func InstallCatppuccin() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	catppuccinPath := filepath.Join(homeDir, ".config", "tmux", "plugins", "catppuccin")

	// Check if already installed
	if IsCatppuccinInstalled() {
		return nil // Already installed, no error
	}

	// Create plugins directory if it doesn't exist
	pluginsDir := filepath.Join(homeDir, ".config", "tmux", "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Clone Catppuccin repository
	cmd := exec.Command("git", "clone", "https://github.com/catppuccin/tmux.git", catppuccinPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone Catppuccin: %w\n%s", err, string(output))
	}

	return nil
}

// ArePluginsInstalled checks if both TPM and Catppuccin are installed
func ArePluginsInstalled() bool {
	return IsTPMInstalled() && IsCatppuccinInstalled()
}
