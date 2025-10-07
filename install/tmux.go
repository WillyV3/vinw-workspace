package install

import (
	"fmt"
	"os/exec"
	"runtime"
)

// InstallMethod represents a tmux installation method
type InstallMethod struct {
	Name        string   // Display name (e.g., "Homebrew (macOS)")
	Description string   // Short description
	Command     string   // Command to run (e.g., "brew")
	Args        []string // Command arguments
	Available   bool     // Whether this method is available on current system
}

// GetTmuxInstallMethods returns available installation methods for the current platform
func GetTmuxInstallMethods() []InstallMethod {
	methods := []InstallMethod{
		{
			Name:        "Homebrew (macOS)",
			Description: "Install via Homebrew package manager",
			Command:     "brew",
			Args:        []string{"install", "tmux"},
			Available:   runtime.GOOS == "darwin" && commandExists("brew"),
		},
		{
			Name:        "apt (Ubuntu/Debian)",
			Description: "Install via apt package manager",
			Command:     "sudo",
			Args:        []string{"apt-get", "install", "-y", "tmux"},
			Available:   runtime.GOOS == "linux" && commandExists("apt-get"),
		},
		{
			Name:        "dnf (Fedora/RHEL 8+)",
			Description: "Install via dnf package manager",
			Command:     "sudo",
			Args:        []string{"dnf", "install", "-y", "tmux"},
			Available:   runtime.GOOS == "linux" && commandExists("dnf"),
		},
		{
			Name:        "yum (CentOS/RHEL 7)",
			Description: "Install via yum package manager",
			Command:     "sudo",
			Args:        []string{"yum", "install", "-y", "tmux"},
			Available:   runtime.GOOS == "linux" && commandExists("yum"),
		},
		{
			Name:        "pacman (Arch Linux)",
			Description: "Install via pacman package manager",
			Command:     "sudo",
			Args:        []string{"pacman", "-S", "--noconfirm", "tmux"},
			Available:   runtime.GOOS == "linux" && commandExists("pacman"),
		},
		{
			Name:        "Build from source (Go)",
			Description: "Build latest tmux from source (requires dev tools)",
			Command:     "sh",
			Args: []string{"-c", `
				set -e
				cd /tmp && \
				git clone https://github.com/tmux/tmux.git tmux-build && \
				cd tmux-build && \
				sh autogen.sh && \
				./configure && \
				make && \
				sudo make install && \
				cd .. && \
				rm -rf tmux-build
			`},
			Available: commandExists("go") && commandExists("git") && commandExists("make"),
		},
	}

	return methods
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// InstallTmux executes the installation command for the selected method
func InstallTmux(method InstallMethod) (string, error) {
	if !method.Available {
		return "", fmt.Errorf("installation method %s is not available on this system", method.Name)
	}

	cmd := exec.Command(method.Command, method.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("installation failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// IsTmuxInstalled checks if tmux is already installed
func IsTmuxInstalled() bool {
	return commandExists("tmux")
}

// GetTmuxVersion returns the installed tmux version
func GetTmuxVersion() (string, error) {
	if !IsTmuxInstalled() {
		return "", fmt.Errorf("tmux is not installed")
	}

	cmd := exec.Command("tmux", "-V")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
