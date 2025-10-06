package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

		// TODO: Add navigation for setup options
		// - Install .tmux.conf template
		// - Install recommended tools
		// - Configure tmux settings
		}
	}

	return m, nil
}

// viewNoobs renders the TMUX Noobs setup view
func viewNoobs(m model) string {
	var s strings.Builder

	s.WriteString(renderStaticGradientTitle(tmuxNoobsAscii))
	s.WriteString("\n")

	// Coming soon badge
	comingSoonStyle := lipgloss.NewStyle().
		Foreground(purpleColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(60)

	s.WriteString(comingSoonStyle.Render("Coming Soon"))
	s.WriteString("\n\n")

	// Feature list
	s.WriteString(sectionTitleStyle.Render("Planned Features:") + "\n\n")

	featureStyle := lipgloss.NewStyle().Foreground(lightGray)
	s.WriteString(featureStyle.Render("  ✓ Install optimized .tmux.conf") + "\n")
	s.WriteString(featureStyle.Render("  ✓ Setup essential tmux plugins (TPM)") + "\n")
	s.WriteString(featureStyle.Render("  ✓ Configure custom key bindings") + "\n")
	s.WriteString(featureStyle.Render("  ✓ Install recommended CLI tools") + "\n")
	s.WriteString(featureStyle.Render("  ✓ Interactive tmux tutorial") + "\n\n")

	s.WriteString(helpStyle.Render("esc: back to menu • q: quit"))

	// Full-height container
	fullHeightContainer := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purpleColor)

	return fullHeightContainer.Render(s.String())
}
