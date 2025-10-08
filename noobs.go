package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/willyv3/vinw-workspace/install"
)

//go:embed tmux-config/tmux.conf
var tmuxConfigTemplate string

//go:embed ascii/tmux-conf-configuration-breakdown.txt
var tmuxConfBreakdownAscii string

//go:embed ascii/tmux-basics.txt
var tmuxBasicsAscii string

//go:embed ascii/config-features.txt
var configFeaturesAscii string

//go:embed ascii/install-tmux.txt
var installTmuxAscii string

//go:embed help-docs/*
var helpDocs embed.FS

// Dracula theme colors for help content
var (
	draculaForeground = lipgloss.Color("#f8f8f2") // Main text
	draculaCyan       = lipgloss.Color("#8be9fd") // Section headers
	draculaGreen      = lipgloss.Color("#50fa7b") // Bullet points
	draculaOrange     = lipgloss.Color("#ffb86c") // Code/commands
	draculaPurple     = lipgloss.Color("#bd93f9") // Emphasis
	draculaPink       = lipgloss.Color("#ff79c6") // Scroll indicator (unique)
	draculaComment    = lipgloss.Color("#6272a4") // Muted text
)

// dialogContent holds the content for a confirmation dialog
type dialogContent struct {
	title       string  // Dialog title (e.g., "⚠️  Install tmux.conf?")
	description string  // Multi-line explanation text
	confirmKey  string  // Key to confirm (e.g., "y")
	cancelKey   string  // Key to cancel (e.g., "n")
	confirmText string  // Confirm button text (e.g., "Yes, install")
	cancelText  string  // Cancel button text (e.g., "No, cancel")
	onConfirm   tea.Cmd // Command to run on confirmation
}

// helpDoc represents a help documentation entry
type helpDoc struct {
	title    string // Menu item title
	asciiArt string // ASCII art for header
	mdFile   string // Path to markdown file relative to help-docs/
}

// noobsHelpContent provides documentation for each noobs option
var noobsHelpContent = map[int]helpDoc{
	0: { // Install tmux
		title:    "Install tmux",
		asciiArt: installTmuxAscii,
		mdFile:   "install-tmux.md",
	},
	1: { // Install tmux.conf
		title:    "Install optimized .tmux.conf",
		asciiArt: tmuxConfBreakdownAscii,
		mdFile:   "tmux-conf-breakdown.md",
	},
	2: { // Tmux Basics
		title:    "Tmux Basics",
		asciiArt: tmuxBasicsAscii,
		mdFile:   "help-basics.md",
	},
	3: { // Cool features in provided config
		title:    "Cool features in provided config",
		asciiArt: configFeaturesAscii,
		mdFile:   "help-features.md",
	},
}

// renderMarkdownHelp reads a markdown file and renders it with glamour
func renderMarkdownHelp(mdFile string, viewportWidth int) (string, error) {
	// Read markdown file from embedded help-docs directory
	helpDocsPath := filepath.Join("help-docs", mdFile)
	mdContent, err := helpDocs.ReadFile(helpDocsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", mdFile, err)
	}

	// Create glamour renderer with dark theme
	// Account for viewport width minus gutter
	const glamourGutter = 2
	glamourWidth := viewportWidth - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(glamourWidth),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create glamour renderer: %w", err)
	}

	// Render markdown to styled string
	rendered, err := renderer.Render(string(mdContent))
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return rendered, nil
}

// openHelp opens the help documentation for the currently selected noobs menu item
func openHelp(m model) (model, tea.Cmd) {
	if helpData, exists := noobsHelpContent[m.noobsCursor]; exists {
		m.currentState = stateNoobsHelp

		// Initialize viewport immediately with current dimensions
		// Container uses m.height-2, with Padding(1,2) and Border
		// Content area = m.height - 2 - 2 (padding) - 2 (border) = m.height - 6
		// Header = 4 lines (ASCII art 3 + % 1)
		// Footer = 3 lines (blank 1 + scroll 1 + quit 1)
		// Viewport = content area - header - footer = m.height - 6 - 4 - 3 = m.height - 13
		viewportWidth := m.width - 8    // m.width - 2 (container) - 4 (padding) - 2 (border)
		viewportHeight := m.height - 13 // m.height - 6 (container content) - 4 (header) - 3 (footer)

		m.helpViewport = viewport.New(viewportWidth, viewportHeight)

		// Render markdown with glamour
		rendered, err := renderMarkdownHelp(helpData.mdFile, viewportWidth)
		if err != nil {
			// Fallback to error message
			m.helpViewport.SetContent(fmt.Sprintf("Error loading help: %v", err))
		} else {
			m.helpViewport.SetContent(rendered)
		}
		m.helpReady = true
	}
	return m, nil
}

// updateNoobs handles the TMUX Noobs setup view
func updateNoobs(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	// If dialog is showing, delegate to dialog handler
	if m.showingNoobsDialog {
		return updateNoobsDialog(msg, m)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Go back to menu
			m.currentState = stateMenu
			return m, nil

		case "?":
			// Open help for current selection
			return openHelp(m)

		case "up", "k":
			if m.noobsCursor > 0 {
				m.noobsCursor--
			}

		case "down", "j":
			// Navigate through help menu items (0, 1, 2)
			if m.noobsCursor < len(noobsHelpContent)-1 {
				m.noobsCursor++
			}

		case "enter":
			// Handle selection - show dialog for confirmation
			switch m.noobsCursor {
			case 0:
				// Install tmux - show installation method selection
				methods := install.GetTmuxInstallMethods()
				// Convert install.InstallMethod to model.installMethod
				m.installMethods = make([]installMethod, 0)
				for _, method := range methods {
					if method.Available {
						m.installMethods = append(m.installMethods, installMethod{
							name:        method.Name,
							description: method.Description,
							command:     method.Command,
							args:        method.Args,
							available:   method.Available,
						})
					}
				}
				m.installMethodCursor = 0
				m.currentState = stateInstallSelection
				return m, nil

			case 1:
				// Install .tmux.conf - check for plugins first
				if !install.ArePluginsInstalled() {
					// Plugins missing - show plugin installation dialog
					m.showingNoobsDialog = true
					m.noobsDialogContent = dialogContent{
						title: "📦 Install tmux plugins?",
						description: `The tmux.conf requires plugins to work properly:
  • TPM (Tmux Plugin Manager) - Plugin management
  • Catppuccin - Theme for status bar

Without these plugins, tmux will show 127 errors.

Install plugins before proceeding with tmux.conf?`,
						confirmKey:  "y",
						cancelKey:   "n",
						confirmText: "Yes, install plugins",
						cancelText:  "Skip plugins (minimal config)",
						onConfirm:   installPluginsThenConfig(),
					}
					return m, nil
				}

				// Plugins already installed - show normal tmux.conf dialog
				m.showingNoobsDialog = true
				m.noobsDialogContent = dialogContent{
					title: "⚠️  Install tmux.conf?",
					description: `This will:
  • Install optimized tmux configuration to ~/.tmux.conf
  • Back up existing config (if found) with timestamp
  • Enable mouse support, 256 colors, and TPM plugins

Your current config will be safely backed up before any changes.`,
					confirmKey:  "y",
					cancelKey:   "n",
					confirmText: "Yes, install",
					cancelText:  "No, cancel",
					onConfirm:   installTmuxConfig(),
				}
				return m, nil

		default:
			// For documentation-only items (no action), open help on Enter
			return openHelp(m)
		}
		}
	}

	return m, nil
}

// updateNoobsDialog handles dialog interactions (y/n confirmation)
func updateNoobsDialog(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case m.noobsDialogContent.confirmKey, "enter":
			// User confirmed - execute action and close dialog
			m.showingNoobsDialog = false
			return m, m.noobsDialogContent.onConfirm

		case m.noobsDialogContent.cancelKey, "esc":
			// User cancelled - close dialog and return to list
			m.showingNoobsDialog = false
			return m, nil
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

// installPluginsThenConfig installs TPM and Catppuccin plugins, then installs tmux.conf
func installPluginsThenConfig() tea.Cmd {
	return func() tea.Msg {
		// Install TPM
		if err := install.InstallTPM(); err != nil {
			return statusMsg{text: fmt.Sprintf("Error installing TPM: %v", err), isError: true}
		}

		// Install Catppuccin
		if err := install.InstallCatppuccin(); err != nil {
			return statusMsg{text: fmt.Sprintf("Error installing Catppuccin: %v", err), isError: true}
		}

		// Now install tmux.conf
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

		return statusMsg{text: "✓ Plugins and .tmux.conf installed successfully!", isError: false}
	}
}

// viewNoobs renders the TMUX Noobs setup view
func viewNoobs(m model) string {
	// If dialog is active, show it instead
	if m.showingNoobsDialog {
		return viewDialog(m, m.noobsDialogContent)
	}

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

	// Render menu options from help content
	optionStyle := lipgloss.NewStyle().Foreground(lightGray)
	selectedStyle := lipgloss.NewStyle().
		Foreground(pinkColor).
		Bold(true)

	// Iterate through available help docs to build menu
	hintStyle := lipgloss.NewStyle().
		Foreground(grayColor).
		Italic(true)

	for i := 0; i < len(noobsHelpContent); i++ {
		if helpData, exists := noobsHelpContent[i]; exists {
			cursor := " "
			if m.noobsCursor == i {
				cursor = "›"
				s.WriteString(selectedStyle.Render(cursor+" "+helpData.title) + "\n")

				// Context-aware hint based on item type
				var hint string
				if i == 0 || i == 1 {
					// Install actions - focus on the action
					hint = "  Press '?' for installation guide"
				} else {
					// Documentation items - learn about topic
					hint = "  Press '?' to learn about " + helpData.title
				}
				s.WriteString(hintStyle.Render(hint) + "\n")
			} else {
				s.WriteString(optionStyle.Render(cursor+" "+helpData.title) + "\n")
			}
		}
	}

	s.WriteString("\n")

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

	s.WriteString(helpStyle.Render("↑/↓: navigate • ?: help • enter: select • esc: back • q: quit"))

	// Full-height container
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width-2).
		Height(m.height-2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return fullHeightContainer.Render(s.String())
}

// viewDialog renders a confirmation dialog with the given content
func viewDialog(m model, content dialogContent) string {
	var s strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(pinkColor).
		Bold(true).
		Align(lipgloss.Center)

	s.WriteString(titleStyle.Render(content.title))
	s.WriteString("\n\n")

	// Description (wrap text for readability)
	descStyle := lipgloss.NewStyle().
		Foreground(lightGray).
		Width(60)

	s.WriteString(descStyle.Render(content.description))
	s.WriteString("\n\n")

	// Buttons
	confirmBtn := lipgloss.NewStyle().
		Foreground(greenColor).
		Bold(true).
		Render(fmt.Sprintf("[%s] %s", content.confirmKey, content.confirmText))

	cancelBtn := lipgloss.NewStyle().
		Foreground(grayColor).
		Render(fmt.Sprintf("[%s] %s", content.cancelKey, content.cancelText))

	s.WriteString("  " + confirmBtn + "  " + cancelBtn)
	s.WriteString("\n")

	// Dialog box with border
	dialogBox := lipgloss.NewStyle().
		Width(70).
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor).
		Align(lipgloss.Center)

	// Center in screen
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogBox.Render(s.String()),
	)

	return centered
}

// updateNoobsHelp handles the help view interactions
func updateNoobsHelp(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "?":
			// Go back to noobs list
			m.currentState = stateNoobs
			m.helpReady = false
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Resize viewport if window dimensions change
		if m.helpReady {
			m.helpViewport.Width = msg.Width - 8
			m.helpViewport.Height = msg.Height - 13
		}
	}

	// Update viewport (handles scrolling)
	m.helpViewport, cmd = m.helpViewport.Update(msg)
	return m, cmd
}

// viewNoobsHelp renders the scrollable help documentation
func viewNoobsHelp(m model) string {
	if !m.helpReady {
		return "Initializing help..."
	}

	var s strings.Builder

	// Get current help data
	helpData := noobsHelpContent[m.noobsCursor]

	// Render ASCII art title with gradient
	s.WriteString(renderStaticGradientTitle(helpData.asciiArt))

	// Scroll percentage indicator at top (bright and visible) - replaces blank line
	scrollPercent := fmt.Sprintf("%3.f%%", m.helpViewport.ScrollPercent()*100)
	scrollStyle := lipgloss.NewStyle().
		Foreground(draculaPink).
		Bold(true).
		Align(lipgloss.Left).
		Width(m.width - 8)
	s.WriteString(scrollStyle.Render(scrollPercent))
	s.WriteString("\n")

	// Scrollable viewport content
	s.WriteString(m.helpViewport.View())
	s.WriteString("\n")

	// Footer with navigation help
	s.WriteString(helpStyle.Render("↑/↓ j/k: scroll") + "\n")
	s.WriteString(helpStyle.Render("esc/?: back • q: quit"))

	// Container
	container := lipgloss.NewStyle().
		Width(m.width-2).
		Height(m.height-2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return container.Render(s.String())
}

// updateInstallSelection handles the installation method selection view
func updateInstallSelection(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		// Installation completed - go back to noobs menu and show status
		m.statusMessage = msg
		m.currentState = stateNoobs
		m.installMethodCursor = 0
		m.installMethods = nil
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Go back to noobs menu
			m.currentState = stateNoobs
			m.installMethodCursor = 0
			m.installMethods = nil
			return m, nil

		case "up", "k":
			if m.installMethodCursor > 0 {
				m.installMethodCursor--
			}

		case "down", "j":
			if m.installMethodCursor < len(m.installMethods)-1 {
				m.installMethodCursor++
			}

		case "enter":
			// Execute installation
			if m.installMethodCursor < len(m.installMethods) {
				selectedMethod := m.installMethods[m.installMethodCursor]
				return m, executeTmuxInstall(selectedMethod)
			}
		}
	}

	return m, nil
}

// executeTmuxInstall runs the installation command for tmux
func executeTmuxInstall(method installMethod) tea.Cmd {
	return func() tea.Msg {
		// Check if tmux is already installed
		if _, err := exec.LookPath("tmux"); err == nil {
			version, _ := exec.Command("tmux", "-V").Output()
			return statusMsg{
				text:    fmt.Sprintf("✓ tmux is already installed: %s", strings.TrimSpace(string(version))),
				isError: false,
			}
		}

		// Execute installation command
		cmd := exec.Command(method.command, method.args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return statusMsg{
				text:    fmt.Sprintf("✗ Installation failed: %s\n%s", err.Error(), string(output)),
				isError: true,
			}
		}

		// Verify installation
		if _, err := exec.LookPath("tmux"); err != nil {
			return statusMsg{
				text:    "✗ Installation completed but tmux not found in PATH",
				isError: true,
			}
		}

		version, _ := exec.Command("tmux", "-V").Output()
		return statusMsg{
			text:    fmt.Sprintf("✓ tmux installed successfully! %s", strings.TrimSpace(string(version))),
			isError: false,
		}
	}
}

// viewInstallSelection renders the installation method selection dialog
func viewInstallSelection(m model) string {
	var s strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(pinkColor).
		Bold(true).
		Align(lipgloss.Left)

	s.WriteString(titleStyle.Render("📦 Install tmux"))
	s.WriteString("\n\n")

	// Description
	descStyle := lipgloss.NewStyle().
		Foreground(lightGray)

	s.WriteString(descStyle.Render("Select your package manager to install tmux:"))
	s.WriteString("\n\n")

	// Installation methods list
	if len(m.installMethods) == 0 {
		noMethodsStyle := lipgloss.NewStyle().
			Foreground(draculaOrange)
		s.WriteString(noMethodsStyle.Render("⚠️  No supported package managers found on this system."))
		s.WriteString("\n\n")
		s.WriteString(descStyle.Render("For manual installation instructions, visit:"))
		s.WriteString("\n")
		linkStyle := lipgloss.NewStyle().Foreground(draculaCyan).Underline(true)
		s.WriteString(linkStyle.Render("https://github.com/tmux/tmux/wiki/Installing"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("[esc] back • [q] quit"))
	} else {
		optionStyle := lipgloss.NewStyle().Foreground(lightGray)
		selectedStyle := lipgloss.NewStyle().
			Foreground(pinkColor).
			Bold(true)
		descOptionStyle := lipgloss.NewStyle().
			Foreground(grayColor).
			MarginLeft(4)

		for i, method := range m.installMethods {
			cursor := "  "
			if m.installMethodCursor == i {
				cursor = "› "
				s.WriteString(selectedStyle.Render(cursor+method.name) + "\n")
				s.WriteString(selectedStyle.Render(descOptionStyle.Render(method.description)) + "\n")
			} else {
				s.WriteString(optionStyle.Render(cursor+method.name) + "\n")
				s.WriteString(optionStyle.Render(descOptionStyle.Render(method.description)) + "\n")
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")

		// Add wiki link for manual installation
		s.WriteString(descStyle.Render("For more installation options:"))
		s.WriteString("\n")
		linkStyle := lipgloss.NewStyle().Foreground(draculaCyan).Underline(true)
		s.WriteString(linkStyle.Render("https://github.com/tmux/tmux/wiki/Installing"))
		s.WriteString("\n\n")

		s.WriteString(helpStyle.Render("↑/↓: navigate • enter: install • esc: back • q: quit"))
	}

	// Dialog box with border - left aligned
	dialogBox := lipgloss.NewStyle().
		Width(70).
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor).
		Align(lipgloss.Left)

	// Center in screen
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogBox.Render(s.String()),
	)

	return centered
}
