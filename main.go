package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle            = blurredStyle
	errorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	radioSelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	radioUnselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	blurredLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle              = lipgloss.NewStyle()
	cursorStyle          = focusedStyle
	focusedButtonStyle   = focusedStyle.Copy().Render("[ Preview & Launch ]")
	blurredButtonStyle   = blurredStyle.Copy().Render("[ Preview & Launch ]")
	selectedStyle        = lipgloss.NewStyle().Background(lipgloss.Color("237"))
	folderStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	fileStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#008080"))
)

type fileEntry struct {
	Name  string
	Path  string
	IsDir bool
}

type state int

const (
	stateInput state = iota
	statePreview
	stateLaunching
)

type model struct {
	width, height   int
	currentState    state
	focusIndex      int
	inputs          []textinput.Model
	terminalCursor  int
	agentCursor     int
	terminalOptions []string
	agentOptions    []string
	directory       string
	files           []fileEntry
	filteredFiles   []fileEntry
	cursor          int
	viewportStart   int
	viewportEnd     int
	newDirInput     textinput.Model
	searchInput     textinput.Model
	creatingNewDir  bool
	searching       bool
	err             error
	shouldLaunch    bool
	launchDir       string
	launchSession   string
	launchTerminal  string
	launchAgent     string
	launchSessionID string
}

func initialModel() model {
	config, _ := loadConfig()

	homeDir, _ := os.UserHomeDir()

	// New directory input
	newDirInput := textinput.New()
	newDirInput.Placeholder = "Enter new directory name..."
	newDirInput.CharLimit = 100
	newDirInput.Width = 50

	// Search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search directories..."
	searchInput.CharLimit = 100
	searchInput.Width = 50

	// Session name input
	sessionInput := textinput.New()
	sessionInput.Placeholder = "my_session"
	sessionInput.Prompt = "Session:   "
	sessionInput.SetValue("dev")
	sessionInput.Cursor.Style = cursorStyle
	sessionInput.CharLimit = 256

	m := model{
		inputs:          []textinput.Model{sessionInput},
		terminalCursor:  0,
		agentCursor:     0,
		currentState:    stateInput,
		terminalOptions: config.TerminalOptions,
		agentOptions:    config.AgentOptions,
		directory:       homeDir,
		newDirInput:     newDirInput,
		searchInput:     searchInput,
		width:           80,
		height:          24,
	}

	m.loadDirectory(homeDir)
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()
		return m, nil

	case tea.KeyMsg:
		switch m.currentState {
		case stateInput:
			return m.updateInput(msg)
		case statePreview:
			return m.updatePreview(msg)
		}
	}

	return m, nil
}

func (m model) updatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.currentState = stateInput
		// Focus back on session name input when returning
		m.focusIndex = 1
		m.inputs[0].Focus()
		m.inputs[0].PromptStyle = focusedStyle
		m.inputs[0].TextStyle = focusedStyle
		return m, nil

	case "enter", "l":
		// Check dependencies and session existence before launching
		deps := checkDependencies(
			m.terminalOptions[m.terminalCursor],
			m.agentOptions[m.agentCursor],
		)

		if !allDependenciesAvailable(deps) {
			// Don't launch if dependencies are missing
			return m, nil
		}

		sessionName := m.inputs[0].Value()
		if sessionExists(sessionName) {
			// Don't launch if session already exists
			return m, nil
		}

		// Store launch parameters and quit
		// The actual launch will happen after the TUI exits
		m.shouldLaunch = true
		m.launchDir = m.directory
		m.launchSession = sessionName
		m.launchTerminal = m.terminalOptions[m.terminalCursor]
		m.launchAgent = m.agentOptions[m.agentCursor]
		m.launchSessionID = generateSessionID(m.directory)

		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.currentState == statePreview {
		return m.viewPreview()
	}
	return m.viewInput()
}

func (m model) viewInput() string {
	var s strings.Builder

	// Full-width header with centered title
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2F4F4F")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	header := headerStyle.Render("Ⓥ Ⓘ Ⓝ Ⓦ   Ⓦ Ⓞ Ⓡ Ⓚ Ⓢ Ⓟ Ⓐ Ⓒ Ⓔ - Launch tmux workspaces")
	s.WriteString(header + "\n\n")

	// Directory browser when focusIndex == 0
	if m.focusIndex == 0 {
		s.WriteString("Choose Directory to launch all terminals in\n")
		homeDir, _ := os.UserHomeDir()
		displayDir := m.directory
		if strings.HasPrefix(displayDir, homeDir) {
			displayDir = "~" + strings.TrimPrefix(displayDir, homeDir)
		}
		s.WriteString(fmt.Sprintf("Current: %s\n", displayDir))

		if m.creatingNewDir {
			s.WriteString(fmt.Sprintf("Creating new directory in: %s\n", displayDir))
			s.WriteString("Name: " + m.newDirInput.View() + "\n")
		} else if m.searching {
			s.WriteString("Search: " + m.searchInput.View() + "\n")
		}

		separator := strings.Repeat("━", m.width-2)
		if len(separator) < 20 {
			separator = "━━━━━━━━━━━━━━━━━━━━"
		}
		s.WriteString(separator + "\n")

		// Show files and directories (viewport)
		for i := m.viewportStart; i < m.viewportEnd; i++ {
			if i >= len(m.filteredFiles) {
				break
			}
			entry := m.filteredFiles[i]
			line := entry.Name
			if entry.IsDir {
				line = folderStyle.Render(line + "/")
			} else {
				line = fileStyle.Render(line)
			}

			if i == m.cursor && !m.creatingNewDir && !m.searching {
				line = selectedStyle.Render(line)
			} else if i == m.cursor && m.searching {
				line = selectedStyle.Render(line)
			}
			s.WriteString(line + "\n")
		}

		s.WriteString(separator + "\n")

		// Help based on mode
		if m.creatingNewDir {
			s.WriteString(helpStyle.Render("Type folder name • enter: create • esc: cancel"))
		} else if m.searching {
			totalFiles := len(m.filteredFiles)
			if totalFiles == 0 {
				s.WriteString(helpStyle.Render("No matches • esc: exit search"))
			} else {
				s.WriteString(helpStyle.Render(fmt.Sprintf("(%d/%d) ↑↓: navigate • enter: select highlighted • esc: exit", m.cursor+1, totalFiles)))
			}
		} else {
			s.WriteString(helpStyle.Render(fmt.Sprintf("(%d/%d) ↑↓/jk: navigate • →: explore dir • ←: up • enter/tab: select highlighted • s/space: search • n: new", m.cursor+1, len(m.filteredFiles))))
		}

		return s.String()
	}

	// Rest of the form
	s.WriteString("Configure your tmux workspace\n\n")

	// Display directory (not an input, just the value)
	homeDir, _ := os.UserHomeDir()
	displayDir := m.directory
	if strings.HasPrefix(displayDir, homeDir) {
		displayDir = "~" + strings.TrimPrefix(displayDir, homeDir)
	}
	s.WriteString(fmt.Sprintf("Directory: %s\n", displayDir))

	// Display session name input
	s.WriteString(m.inputs[0].View())
	s.WriteString("\n")

	s.WriteString("\n")

	terminalLabel := "Terminal:"
	if m.focusIndex == 2 {
		terminalLabel = focusedLabelStyle.Render("› Terminal:")
	} else {
		terminalLabel = blurredLabelStyle.Render("  Terminal:")
	}
	s.WriteString(terminalLabel + "\n")

	for i, opt := range m.terminalOptions {
		cursor := "○"
		style := radioUnselectedStyle

		if m.terminalCursor == i {
			cursor = "●"
			if m.focusIndex == 2 {
				style = radioSelectedStyle
			}
		}

		label := opt
		if opt == "nextui" {
			label = "nextui (Next.js scaffolder)"
		} else if opt == "shell" {
			label = "shell (empty terminal)"
		}

		s.WriteString(fmt.Sprintf("    %s %s\n", style.Render(cursor), style.Render(label)))
	}

	s.WriteString("\n")

	agentLabel := "Coding Agent:"
	if m.focusIndex == 3 {
		agentLabel = focusedLabelStyle.Render("› Coding Agent:")
	} else {
		agentLabel = blurredLabelStyle.Render("  Coding Agent:")
	}
	s.WriteString(agentLabel + "\n")

	for i, opt := range m.agentOptions {
		cursor := "○"
		style := radioUnselectedStyle

		if m.agentCursor == i {
			cursor = "●"
			if m.focusIndex == 3 {
				style = radioSelectedStyle
			}
		}

		s.WriteString(fmt.Sprintf("    %s %s\n", style.Render(cursor), style.Render(opt)))
	}

	s.WriteString("\n")

	button := blurredButtonStyle
	if m.focusIndex == 4 {
		button = focusedButtonStyle
	}
	s.WriteString(button)
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	helpText := "Tab/↑/↓: navigate • ←/→: select option • Enter: next/launch • Esc: back"
	s.WriteString(helpStyle.Render(helpText))

	return s.String()
}

func (m model) viewPreview() string {
	var s strings.Builder

	// Full-width header with centered title
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2F4F4F")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	header := headerStyle.Render("P R E V I E W")
	s.WriteString(header + "\n\n")

	sessionName := m.inputs[0].Value()
	sessionAlreadyExists := sessionExists(sessionName)

	deps := checkDependencies(
		m.terminalOptions[m.terminalCursor],
		m.agentOptions[m.agentCursor],
	)

	s.WriteString(getPreviewContent(
		m.directory,
		sessionName,
		m.terminalOptions[m.terminalCursor],
		m.agentOptions[m.agentCursor],
	))

	// Check if session already exists
	if sessionAlreadyExists {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("⚠ Session '%s' already exists!", sessionName)))
		s.WriteString("\n")

		// Show existing sessions
		existingSessions := getTmuxSessions()
		if len(existingSessions) > 0 {
			s.WriteString("\nExisting sessions:\n")
			for _, sess := range existingSessions {
				if sess == sessionName {
					s.WriteString(fmt.Sprintf("  %s (current conflict)\n", errorStyle.Render("• "+sess)))
				} else {
					s.WriteString(fmt.Sprintf("  • %s\n", sess))
				}
			}
		}
	}

	s.WriteString("\n\nDependencies:\n")
	for _, dep := range deps {
		if dep.Available {
			s.WriteString(fmt.Sprintf("  %s %s\n", successStyle.Render("✓"), dep.Name))
		} else {
			s.WriteString(fmt.Sprintf("  %s %s (not found)\n", errorStyle.Render("✗"), dep.Name))
		}
	}

	s.WriteString("\n")

	canLaunch := allDependenciesAvailable(deps) && !sessionAlreadyExists

	if canLaunch {
		s.WriteString(focusedStyle.Render("[ Launch Session ]"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Enter/l: launch • Esc: back • q: quit"))
	} else if sessionAlreadyExists {
		s.WriteString(errorStyle.Render("Cannot launch: session name already in use"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Esc: back to change session name • q: quit"))
	} else {
		// Check if tmux is missing
		tmuxMissing := false
		for _, dep := range deps {
			if dep.Name == "tmux" && !dep.Available {
				tmuxMissing = true
				break
			}
		}

		if tmuxMissing {
			s.WriteString(errorStyle.Render("tmux is not installed"))
			s.WriteString("\n\n")
			s.WriteString("Install tmux for your platform:\n\n")

			// Create table using lipgloss
			tableStyle := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

			rows := [][]string{
				{"Arch Linux", "pacman -S tmux"},
				{"Debian or Ubuntu", "apt install tmux"},
				{"Fedora", "dnf install tmux"},
				{"RHEL or CentOS", "yum install tmux"},
				{"macOS (Homebrew)", "brew install tmux"},
				{"macOS (MacPorts)", "port install tmux"},
				{"openSUSE", "zypper install tmux"},
			}

			var tableContent strings.Builder
			headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
			cellStyle := lipgloss.NewStyle()

			// Header
			tableContent.WriteString(headerStyle.Render("Platform") + "              " + headerStyle.Render("Install Command") + "\n")

			// Rows
			for _, row := range rows {
				tableContent.WriteString(cellStyle.Render(fmt.Sprintf("%-20s  %s\n", row[0], row[1])))
			}

			s.WriteString(tableStyle.Render(tableContent.String()))
			s.WriteString("\n\n")
			s.WriteString(blurredStyle.Render("More info: https://github.com/tmux/tmux/wiki/Installing"))
		} else {
			s.WriteString(errorStyle.Render("Cannot launch: missing dependencies"))
		}
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Esc: back • q: quit"))
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if we should launch a session after the TUI exits
	if m, ok := finalModel.(model); ok && m.shouldLaunch {
		if err := launchTmuxSession(
			m.launchDir,
			m.launchSession,
			m.launchTerminal,
			m.launchAgent,
			m.launchSessionID,
		); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
