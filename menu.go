package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuChoice int

const (
	menuNewWorkspace menuChoice = iota
	menuNoobs
)

var menuItems = []struct {
	title       string
	description string
}{
	{"Start New Workspace", "→ Configure directory, terminal, and agent"},
	{"TMUX Noobs", "→ Setup tmux configuration and tools"},
}

// updateMenu handles the landing menu update logic
func updateMenu(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Navigate to selected option
			switch menuChoice(m.menuCursor) {
			case menuNewWorkspace:
				// Go to form state
				m.currentState = stateForm
				return m, nil
			case menuNoobs:
				// Go to noobs setup state
				m.currentState = stateNoobs
				return m, nil
			}

		case "down", "j":
			m.menuCursor++
			if m.menuCursor >= len(menuItems) {
				m.menuCursor = 0
			}

		case "up", "k":
			m.menuCursor--
			if m.menuCursor < 0 {
				m.menuCursor = len(menuItems) - 1
			}
		}
	}

	return m, nil
}

// viewMenu renders the landing menu
func viewMenu(m model) string {
	var s strings.Builder

	// Animated ASCII art title with gradient
	// Use current animation frame from Update() - View only renders
	animatedTitle := renderAnimatedTitle(m.animFrame)

	// Center the title
	titleLines := strings.Split(strings.TrimRight(animatedTitle, "\n"), "\n")
	for _, line := range titleLines {
		centered := lipgloss.NewStyle().
			Width(60).
			Align(lipgloss.Center).
			Render(line)
		s.WriteString(centered)
		s.WriteString("\n")
	}
	s.WriteString("\n")

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Width(60)

	s.WriteString(subtitleStyle.Render("Interactive TUI for launching tmux workspaces"))
	s.WriteString("\n\n\n")

	// Menu items with proper alignment
	menuStyle := lipgloss.NewStyle().Width(60)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(60).
		PaddingLeft(6)

	for i, item := range menuItems {
		var menuLine string
		if i == m.menuCursor {
			menuLine = "  (•) " + item.title
		} else {
			menuLine = "  ( ) " + item.title
		}
		s.WriteString(menuStyle.Render(menuLine))
		s.WriteString("\n")
		s.WriteString(descStyle.Render(item.description))
		s.WriteString("\n")
		if i < len(menuItems)-1 {
			s.WriteString("\n")
		}
	}

	// Help text - centered
	s.WriteString("\n\n")
	helpCentered := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(60).
		Align(lipgloss.Center)
	s.WriteString(helpCentered.Render("↑↓/jk: navigate • enter: select • q: quit"))

	// Full-height container (leave margin for border)
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center)

	return fullHeightContainer.Render(s.String())
}
