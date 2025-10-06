package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed tmux-config/tmux.conf
var tmuxConfigTemplate string

// updateNoobs handles the TMUX Noobs setup view
func updateNoobs(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Go back to menu
			m.currentState = stateMenu
			return m, nil

		case "up", "k":
			if m.noobsCursor > 0 {
				m.noobsCursor--
			}

		case "down", "j":
			// We have 1 option for now (install .tmux.conf)
			if m.noobsCursor < 0 {
				m.noobsCursor++
			}

		case "enter":
			// Handle selection
			switch m.noobsCursor {
			case 0:
				// Install .tmux.conf
				return m, installTmuxConfig()
			}
		}
	}

	return m, nil
}

// installTmuxConfig installs the embedded tmux.conf to ~/.tmux.conf
func installTmuxConfig() tea.Cmd {
	return func() tea.Msg {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return statusMsg{text: "Error: Could not find home directory", isError: true}
		}

		tmuxConfPath := filepath.Join(homeDir, ".tmux.conf")

		// Check if file already exists
		if _, err := os.Stat(tmuxConfPath); err == nil {
			// Create timestamped backup
			timestamp := time.Now().Format("20060102-150405")
			backupPath := fmt.Sprintf("%s.backup.%s", tmuxConfPath, timestamp)

			// Copy existing config to backup
			originalContent, err := os.ReadFile(tmuxConfPath)
			if err != nil {
				return statusMsg{text: "Error: Could not read existing .tmux.conf", isError: true}
			}
			if err := os.WriteFile(backupPath, originalContent, 0644); err != nil {
				return statusMsg{text: "Error: Could not backup existing .tmux.conf", isError: true}
			}
		}

		// Write the new config
		if err := os.WriteFile(tmuxConfPath, []byte(tmuxConfigTemplate), 0644); err != nil {
			return statusMsg{text: "Error: Could not write .tmux.conf", isError: true}
		}

		return statusMsg{text: "✓ .tmux.conf installed successfully!", isError: false}
	}
}

// viewNoobs renders the TMUX Noobs setup view
func viewNoobs(m model) string {
	var s strings.Builder

	s.WriteString(renderStaticGradientTitle(tmuxNoobsAscii))
	s.WriteString("\n")

	// Description
	descStyle := lipgloss.NewStyle().
		Foreground(lightGray).
		Align(lipgloss.Center)

	s.WriteString(descStyle.Render("Get started with tmux - install config and learn the basics"))
	s.WriteString("\n\n")

	// Setup options
	s.WriteString(sectionTitleStyle.Render("Setup Options:") + "\n\n")

	// Option 1: Install .tmux.conf
	optionStyle := lipgloss.NewStyle().Foreground(lightGray)
	selectedStyle := lipgloss.NewStyle().
		Foreground(pinkColor).
		Bold(true)

	cursor := " "
	if m.noobsCursor == 0 {
		cursor = "›"
		s.WriteString(selectedStyle.Render(cursor + " Install optimized .tmux.conf") + "\n")
	} else {
		s.WriteString(optionStyle.Render(cursor + " Install optimized .tmux.conf") + "\n")
	}

	s.WriteString("\n")

	// Coming soon features
	s.WriteString(sectionTitleStyle.Render("Coming Soon:") + "\n\n")

	featureStyle := lipgloss.NewStyle().Foreground(grayColor)
	s.WriteString(featureStyle.Render("  • Setup essential tmux plugins (TPM)") + "\n")
	s.WriteString(featureStyle.Render("  • Configure custom key bindings") + "\n")
	s.WriteString(featureStyle.Render("  • Install recommended CLI tools") + "\n")
	s.WriteString(featureStyle.Render("  • Interactive tmux tutorial") + "\n\n")

	// Show status message if available
	if m.statusMessage.text != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(greenColor).
			Bold(true)
		if m.statusMessage.isError {
			statusStyle = statusStyle.Foreground(redColor)
		}
		s.WriteString(statusStyle.Render(m.statusMessage.text) + "\n\n")
	}

	s.WriteString(helpStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"))

	// Full-height container
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return fullHeightContainer.Render(s.String())
}
