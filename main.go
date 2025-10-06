package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	purpleColor = lipgloss.Color("62")  // Accent purple (matches menu border)
	greenColor  = lipgloss.Color("42")  // Success green
	pinkColor   = lipgloss.Color("205") // Focused pink
	grayColor   = lipgloss.Color("240") // Muted gray
	lightGray   = lipgloss.Color("252") // Bright gray
	redColor    = lipgloss.Color("196") // Error red
	cyanColor   = lipgloss.Color("51")  // Cyan accent

	// Core styles
	focusedStyle         = lipgloss.NewStyle().Foreground(pinkColor)
	blurredStyle         = lipgloss.NewStyle().Foreground(grayColor)
	helpStyle            = lipgloss.NewStyle().Foreground(lightGray)
	errorStyle           = lipgloss.NewStyle().Foreground(redColor)
	successStyle         = lipgloss.NewStyle().Foreground(greenColor)
	radioSelectedStyle   = lipgloss.NewStyle().Foreground(greenColor)
	radioUnselectedStyle = lipgloss.NewStyle().Foreground(grayColor)
	focusedLabelStyle    = lipgloss.NewStyle().Foreground(pinkColor).Bold(true)
	blurredLabelStyle    = lipgloss.NewStyle().Foreground(grayColor)
	noStyle              = lipgloss.NewStyle()
	cursorStyle          = focusedStyle
	selectedStyle        = lipgloss.NewStyle().Background(lipgloss.Color("237"))
	folderStyle          = lipgloss.NewStyle().Foreground(cyanColor).Bold(true)
	fileStyle            = lipgloss.NewStyle().Foreground(grayColor)

	// Container styles (matches menu)
	containerStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleColor)

	titleStyle = lipgloss.NewStyle().
			Foreground(pinkColor).
			Bold(true).
			Padding(0, 0, 1, 0)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(purpleColor).
				Bold(true)
)

type fileEntry struct {
	Name  string
	Path  string
	IsDir bool
}

type state int

const (
	stateMenu state = iota
	stateForm
	statePreview
	stateLaunching
	stateNoobs
	stateCommands
)

type model struct {
	width, height      int
	currentState       state
	menuCursor         int
	animFrame          int
	focusIndex         int
	inputs             []textinput.Model
	terminalCursor     int
	agentCursor        int
	terminalOptions    []string
	agentOptions       []string
	directory          string
	files              []fileEntry
	filteredFiles      []fileEntry
	cursor             int
	viewportStart      int
	viewportEnd        int
	newDirInput        textinput.Model
	searchInput        textinput.Model
	creatingNewDir     bool
	searching          bool
	err                error
	shouldLaunch       bool
	launchDir          string
	launchSession      string
	launchTerminal     string
	launchAgent        string
	launchSessionID    string
	launchCustomCmd    string
	commandsList       list.Model
	workspaceCommands  []WorkspaceCommand
	selectedCommandIdx int
	addingCommand      bool
	commandNameInput   textinput.Model
	commandCmdInput    textinput.Model
	commandDescInput   textinput.Model
	noobsCursor        int
	statusMessage      statusMsg
}

type statusMsg struct {
	text    string
	isError bool
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

	// Load workspace commands
	workspaceCommands, _ := loadWorkspaceCommands()

	// Create list items from commands
	items := make([]list.Item, len(workspaceCommands))
	for i, cmd := range workspaceCommands {
		items[i] = commandItem{
			name:        cmd.Name,
			command:     cmd.Command,
			description: cmd.Description,
		}
	}

	// Create commands list
	commandsList := list.New(items, commandDelegate{}, 0, 0)
	commandsList.Title = "Custom Commands"
	commandsList.SetShowStatusBar(false)
	commandsList.SetFilteringEnabled(false)

	// Create command input fields
	commandNameInput := textinput.New()
	commandNameInput.Placeholder = "e.g., 'Dev Server'"
	commandNameInput.CharLimit = 50
	commandNameInput.Width = 40

	commandCmdInput := textinput.New()
	commandCmdInput.Placeholder = "e.g., 'npm run dev'"
	commandCmdInput.CharLimit = 100
	commandCmdInput.Width = 40

	commandDescInput := textinput.New()
	commandDescInput.Placeholder = "e.g., 'Start development server'"
	commandDescInput.CharLimit = 100
	commandDescInput.Width = 40

	m := model{
		inputs:             []textinput.Model{sessionInput},
		terminalCursor:     0,
		agentCursor:        0,
		menuCursor:         0,
		animFrame:          0,
		currentState:       stateMenu,
		terminalOptions:    config.TerminalOptions,
		agentOptions:       config.AgentOptions,
		directory:          homeDir,
		newDirInput:        newDirInput,
		searchInput:        searchInput,
		commandsList:       commandsList,
		workspaceCommands:  workspaceCommands,
		selectedCommandIdx: -1,
		addingCommand:      false,
		commandNameInput:   commandNameInput,
		commandCmdInput:    commandCmdInput,
		commandDescInput:   commandDescInput,
		width:              80,
		height:             24,
	}

	m.loadDirectory(homeDir)
	return m
}

func (m model) Init() tea.Cmd {
	// Start with animation tick for smooth entrance
	return tea.Batch(
		textinput.Blink,
		tickAnimation(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()
		// Update commands list size
		m.commandsList.SetSize(m.width-10, m.height-15)
		return m, nil

	case animationMsg:
		// Advance animation frame
		if m.animFrame < animationFrames {
			m.animFrame++
			return m, tickAnimation()
		}
		return m, nil

	case statusMsg:
		m.statusMessage = msg
		return m, nil

	case tea.KeyMsg:
		// Delegate to appropriate update function based on current state
		switch m.currentState {
		case stateMenu:
			return updateMenu(msg, m)
		case stateForm:
			return m.updateInput(msg)
		case statePreview:
			return m.updatePreview(msg)
		case stateNoobs:
			return updateNoobs(msg, m)
		case stateCommands:
			return updateCommands(msg, m)
		}
	}

	return m, nil
}

func (m model) updatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.currentState = stateForm
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

		// Include custom command if selected
		if m.selectedCommandIdx >= 0 && m.selectedCommandIdx < len(m.workspaceCommands) {
			m.launchCustomCmd = m.workspaceCommands[m.selectedCommandIdx].Command
		}

		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	// Delegate to appropriate view function based on current state
	switch m.currentState {
	case stateMenu:
		return viewMenu(m)
	case stateForm:
		return m.viewInput()
	case statePreview:
		return m.viewPreview()
	case stateNoobs:
		return viewNoobs(m)
	case stateCommands:
		return viewCommands(m)
	default:
		return "Unknown state"
	}
}

func (m model) viewInput() string {
	var s strings.Builder

	// Directory browser when focusIndex == 0
	if m.focusIndex == 0 {
		// Content width: window - margin - border - padding
		// = (m.width - 2) - 2 - 4 = m.width - 8
		contentWidth := m.width - 8
		if contentWidth < 40 {
			contentWidth = 40
		}

		// Render ASCII title with gradient
		s.WriteString(renderStaticGradientTitle(chooseDirAscii))
		s.WriteString("\n")

		homeDir, _ := os.UserHomeDir()
		displayDir := m.directory
		if strings.HasPrefix(displayDir, homeDir) {
			displayDir = "~" + strings.TrimPrefix(displayDir, homeDir)
		}

		// Truncate path if too long
		if len(displayDir) > contentWidth-10 {
			displayDir = "..." + displayDir[len(displayDir)-(contentWidth-13):]
		}
		s.WriteString(blurredStyle.Render("Current: ") + successStyle.Render(displayDir))
		s.WriteString("\n\n")

		if m.creatingNewDir {
			s.WriteString(sectionTitleStyle.Render("Create Directory") + "\n")
			s.WriteString(blurredStyle.Render("In: "+displayDir) + "\n")
			s.WriteString("Name: " + m.newDirInput.View() + "\n")
		} else if m.searching {
			s.WriteString(sectionTitleStyle.Render("Search") + "\n")
			s.WriteString(m.searchInput.View() + "\n")
		}

		// Separator with consistent width
		sepStyle := lipgloss.NewStyle().Foreground(purpleColor)
		s.WriteString("\n" + sepStyle.Render(strings.Repeat("─", contentWidth)) + "\n\n")

		// Show files and directories (viewport)
		for i := m.viewportStart; i < m.viewportEnd; i++ {
			if i >= len(m.filteredFiles) {
				break
			}
			entry := m.filteredFiles[i]
			line := entry.Name

			// Truncate long filenames
			maxNameLen := contentWidth - 6
			if len(line) > maxNameLen {
				line = line[:maxNameLen-3] + "..."
			}

			if entry.IsDir {
				line = folderStyle.Render("📂 " + line)
			} else {
				line = fileStyle.Render("   " + line)
			}

			if i == m.cursor {
				line = selectedStyle.Render("▸ " + line)
			} else {
				line = "  " + line
			}
			s.WriteString(line + "\n")
		}

		s.WriteString("\n" + sepStyle.Render(strings.Repeat("─", contentWidth)) + "\n\n")

		// Help based on mode
		if m.creatingNewDir {
			s.WriteString(helpStyle.Render("enter: create • esc: cancel"))
		} else if m.searching {
			totalFiles := len(m.filteredFiles)
			if totalFiles == 0 {
				s.WriteString(helpStyle.Render("No matches • esc: exit search"))
			} else {
				s.WriteString(helpStyle.Render(fmt.Sprintf("(%d/%d) ↑↓: navigate • enter: explore • tab: select • esc: exit", m.cursor+1, totalFiles)))
			}
		} else {
			s.WriteString(helpStyle.Render(fmt.Sprintf("(%d/%d) ↑↓/jk: nav • →: explore • ←: up • enter/tab: select • s: search • n: new", m.cursor+1, len(m.filteredFiles))))
		}

		// Full-height container (leave 1 char margin for border on each side)
		fullHeightContainer := lipgloss.NewStyle().
			Width(m.width - 2).
			Height(m.height - 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleColor)

		return fullHeightContainer.Render(s.String())
	}

	// Rest of the form - ASCII title with gradient
	s.WriteString(renderStaticGradientTitle(configWorkspaceAscii))
	s.WriteString("\n")

	// Display directory
	homeDir, _ := os.UserHomeDir()
	displayDir := m.directory
	if strings.HasPrefix(displayDir, homeDir) {
		displayDir = "~" + strings.TrimPrefix(displayDir, homeDir)
	}
	s.WriteString(blurredStyle.Render("Directory: ") + successStyle.Render(displayDir))
	s.WriteString("\n\n")

	// Session name input
	s.WriteString(m.inputs[0].View())
	s.WriteString("\n\n")

	// Terminal selection
	terminalLabel := "Terminal:"
	terminalHint := ""
	if m.focusIndex == 2 {
		terminalLabel = focusedLabelStyle.Render("› Terminal:")
		terminalHint = "    " + blurredStyle.Render("(Press 'c' to add custom commands)")
	} else {
		terminalLabel = blurredLabelStyle.Render("  Terminal:")
	}
	s.WriteString(terminalLabel + terminalHint + "\n")

	// Column layout for narrow terminals
	if m.width < 100 {
		// Two columns
		var col1, col2 strings.Builder
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
				label = "nextui"
			} else if opt == "shell" {
				label = "shell"
			}

			line := fmt.Sprintf("  %s %s", style.Render(cursor), style.Render(label))
			if i%2 == 0 {
				col1.WriteString(line + "\n")
			} else {
				col2.WriteString(line + "\n")
			}
		}

		terminalCols := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width((m.width-8)/2).Render(col1.String()),
			lipgloss.NewStyle().Width((m.width-8)/2).Render(col2.String()),
		)
		s.WriteString(terminalCols)
	} else {
		// Single column
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
	}

	s.WriteString("\n")

	// Agent selection
	agentLabel := "Coding Agent:"
	if m.focusIndex == 3 {
		agentLabel = focusedLabelStyle.Render("› Coding Agent:")
	} else {
		agentLabel = blurredLabelStyle.Render("  Coding Agent:")
	}
	s.WriteString(agentLabel + "\n")

	// Column layout for narrow terminals
	if m.width < 100 {
		// Two columns
		var col1, col2 strings.Builder
		for i, opt := range m.agentOptions {
			cursor := "○"
			style := radioUnselectedStyle

			if m.agentCursor == i {
				cursor = "●"
				if m.focusIndex == 3 {
					style = radioSelectedStyle
				}
			}

			line := fmt.Sprintf("  %s %s", style.Render(cursor), style.Render(opt))
			if i%2 == 0 {
				col1.WriteString(line + "\n")
			} else {
				col2.WriteString(line + "\n")
			}
		}

		agentCols := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width((m.width-8)/2).Render(col1.String()),
			lipgloss.NewStyle().Width((m.width-8)/2).Render(col2.String()),
		)
		s.WriteString(agentCols)
	} else {
		// Single column
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
	}

	s.WriteString("\n")

	// Show selected custom command if any
	if m.selectedCommandIdx >= 0 && m.selectedCommandIdx < len(m.workspaceCommands) {
		selectedCmd := m.workspaceCommands[m.selectedCommandIdx]
		cmdLabel := blurredLabelStyle.Render("  Custom Command:")
		s.WriteString(cmdLabel + "\n")
		s.WriteString(fmt.Sprintf("    %s\n", successStyle.Render("✓ "+selectedCmd.Name)))
		s.WriteString(fmt.Sprintf("      %s\n", blurredStyle.Render(selectedCmd.Command)))
		s.WriteString("\n")
	}

	// Launch button
	buttonStyle := blurredStyle
	if m.focusIndex == 4 {
		buttonStyle = focusedStyle
	}
	s.WriteString(buttonStyle.Render("[ Preview & Launch ]"))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	helpText := "tab/↑↓: nav • ←→: select • c: custom commands • enter: next • esc: back"
	s.WriteString(helpStyle.Render(helpText))

	// Full-height container
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return fullHeightContainer.Render(s.String())
}

func (m model) viewPreview() string {
	var s strings.Builder

	// Title at top
	s.WriteString(renderStaticGradientTitle(previewAscii))
	s.WriteString("\n\n")

	sessionName := m.inputs[0].Value()
	sessionAlreadyExists := sessionExists(sessionName)

	deps := checkDependencies(
		m.terminalOptions[m.terminalCursor],
		m.agentOptions[m.agentCursor],
	)

	// Get custom command if selected
	customCmd := ""
	if m.selectedCommandIdx >= 0 && m.selectedCommandIdx < len(m.workspaceCommands) {
		customCmd = m.workspaceCommands[m.selectedCommandIdx].Command
	}

	// Layout diagram
	s.WriteString(sectionTitleStyle.Render("Pane Layout") + "\n")
	layoutBox := `┌─────┬──────────────┐
│vinw │ vinw-viewer  │
│     ├──────┬───────┤
│     │ term │ agent │
└─────┴──────┴───────┘`
	s.WriteString(blurredStyle.Render(layoutBox))
	s.WriteString("\n\n")

	// Configuration
	homeDir, _ := os.UserHomeDir()
	displayDir := m.directory
	if strings.HasPrefix(displayDir, homeDir) {
		displayDir = "~" + strings.TrimPrefix(displayDir, homeDir)
	}

	s.WriteString(sectionTitleStyle.Render("Configuration") + "\n")
	s.WriteString(fmt.Sprintf("  %s %s\n", blurredStyle.Render("Directory:"), successStyle.Render(displayDir)))
	s.WriteString(fmt.Sprintf("  %s %s\n", blurredStyle.Render("Session:"), successStyle.Render(sessionName)))

	terminalDisplay := m.terminalOptions[m.terminalCursor]
	if customCmd != "" {
		terminalDisplay = "custom"
	}
	s.WriteString(fmt.Sprintf("  %s %s\n", blurredStyle.Render("Terminal:"), successStyle.Render(terminalDisplay)))
	s.WriteString(fmt.Sprintf("  %s %s\n", blurredStyle.Render("Agent:"), successStyle.Render(m.agentOptions[m.agentCursor])))

	// Session conflict warning
	if sessionAlreadyExists {
		s.WriteString("\n")
		s.WriteString(errorStyle.Render("  ⚠ Session exists! Choose different name"))
	}

	s.WriteString("\n\n")

	// Dependencies
	s.WriteString(sectionTitleStyle.Render("Dependencies") + "\n")
	for _, dep := range deps {
		if dep.Available {
			s.WriteString(fmt.Sprintf("  %s %s\n", successStyle.Render("✓"), dep.Name))
		} else {
			s.WriteString(fmt.Sprintf("  %s %s\n", errorStyle.Render("✗"), dep.Name))
		}
	}

	s.WriteString("\n")

	canLaunch := allDependenciesAvailable(deps) && !sessionAlreadyExists

	if canLaunch {
		s.WriteString(successStyle.Render("✓ Ready to launch!"))
		s.WriteString("\n\n")
		s.WriteString(focusedStyle.Render("[ Launch Session ]"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("enter/l: launch • esc: back • q: quit"))
	} else if sessionAlreadyExists {
		s.WriteString(helpStyle.Render("esc: back to change name • q: quit"))
	} else {
		tmuxMissing := false
		for _, dep := range deps {
			if dep.Name == "tmux" && !dep.Available {
				tmuxMissing = true
				break
			}
		}

		if tmuxMissing {
			s.WriteString(errorStyle.Render("⚠ tmux not installed"))
			s.WriteString("\n\n")
			s.WriteString(blurredStyle.Render("Install: brew install tmux (macOS)"))
			s.WriteString("\n")
			s.WriteString(blurredStyle.Render("        apt install tmux (Debian/Ubuntu)"))
		} else {
			s.WriteString(errorStyle.Render("⚠ Missing dependencies"))
		}
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("esc: back • q: quit"))
	}

	// Full-height container (account for border)
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width - 4).   // -4 for border (2) + margin (2)
		Height(m.height - 4). // -4 for border (2) + margin (2)
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return fullHeightContainer.Render(s.String())
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
			m.launchCustomCmd,
		); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
