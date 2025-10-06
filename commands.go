package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// commandItem implements list.Item for workspace commands
type commandItem struct {
	name        string
	command     string
	description string
}

func (i commandItem) FilterValue() string { return i.name }
func (i commandItem) Title() string       { return i.name }
func (i commandItem) Description() string { return i.description }

// commandDelegate renders list items
type commandDelegate struct{}

func (d commandDelegate) Height() int                             { return 2 }
func (d commandDelegate) Spacing() int                            { return 1 }
func (d commandDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d commandDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(commandItem)
	if !ok {
		return
	}

	// Styles
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Render based on selection
	var title, desc string
	if index == m.Index() {
		title = selectedTitleStyle.Render("▸ " + i.Title())
		desc = descStyle.Render("  " + i.Description())
	} else {
		title = titleStyle.Render("  " + i.Title())
		desc = descStyle.Render("  " + i.Description())
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

// updateCommands handles the custom commands list view
func updateCommands(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't let list handle keys if we're adding a command
		if m.addingCommand {
			switch msg.String() {
			case "esc":
				m.addingCommand = false
				m.commandNameInput.Blur()
				m.commandCmdInput.Blur()
				m.commandDescInput.Blur()
				return m, nil
			case "enter":
				// Save the new command
				name := strings.TrimSpace(m.commandNameInput.Value())
				cmd := strings.TrimSpace(m.commandCmdInput.Value())
				desc := strings.TrimSpace(m.commandDescInput.Value())

				if name != "" && cmd != "" {
					newCmd := WorkspaceCommand{
						Name:        name,
						Command:     cmd,
						Description: desc,
					}

					// Add to list
					m.workspaceCommands = append(m.workspaceCommands, newCmd)

					// Save to disk
					saveWorkspaceCommands(m.workspaceCommands)

					// Update list items
					items := make([]list.Item, len(m.workspaceCommands))
					for i, c := range m.workspaceCommands {
						items[i] = commandItem{
							name:        c.Name,
							command:     c.Command,
							description: c.Description,
						}
					}
					m.commandsList.SetItems(items)

					// Reset inputs
					m.commandNameInput.SetValue("")
					m.commandCmdInput.SetValue("")
					m.commandDescInput.SetValue("")
					m.addingCommand = false
					m.commandNameInput.Blur()
					m.commandCmdInput.Blur()
					m.commandDescInput.Blur()
				}
				return m, nil
			case "tab":
				// Cycle through inputs
				if m.commandNameInput.Focused() {
					m.commandNameInput.Blur()
					m.commandCmdInput.Focus()
				} else if m.commandCmdInput.Focused() {
					m.commandCmdInput.Blur()
					m.commandDescInput.Focus()
				} else {
					m.commandDescInput.Blur()
					m.commandNameInput.Focus()
				}
				return m, nil
			default:
				// Handle input updates
				var cmd tea.Cmd
				if m.commandNameInput.Focused() {
					m.commandNameInput, cmd = m.commandNameInput.Update(msg)
				} else if m.commandCmdInput.Focused() {
					m.commandCmdInput, cmd = m.commandCmdInput.Update(msg)
				} else if m.commandDescInput.Focused() {
					m.commandDescInput, cmd = m.commandDescInput.Update(msg)
				}
				return m, cmd
			}
		}

		// Normal list navigation
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			// Go back to form
			m.currentState = stateForm
			m.focusIndex = 1 // Back to session name
			return m, nil
		case "a":
			// Start adding a command
			m.addingCommand = true
			m.commandNameInput.Focus()
			return m, nil
		case "d":
			// Delete selected command
			if len(m.workspaceCommands) > 0 {
				idx := m.commandsList.Index()
				if idx >= 0 && idx < len(m.workspaceCommands) {
					// Remove from slice
					m.workspaceCommands = append(m.workspaceCommands[:idx], m.workspaceCommands[idx+1:]...)

					// Save to disk
					saveWorkspaceCommands(m.workspaceCommands)

					// Update list items
					items := make([]list.Item, len(m.workspaceCommands))
					for i, c := range m.workspaceCommands {
						items[i] = commandItem{
							name:        c.Name,
							command:     c.Command,
							description: c.Description,
						}
					}
					m.commandsList.SetItems(items)
				}
			}
			return m, nil
		case "enter":
			// Select command to run (mark it in model)
			if len(m.workspaceCommands) > 0 {
				idx := m.commandsList.Index()
				if idx >= 0 && idx < len(m.workspaceCommands) {
					m.selectedCommandIdx = idx
				}
			}
			// Go back to form
			m.currentState = stateForm
			m.focusIndex = 1
			return m, nil
		}
	}

	// Update list
	var cmd tea.Cmd
	m.commandsList, cmd = m.commandsList.Update(msg)
	return m, cmd
}

// viewCommands renders the custom commands list view
func viewCommands(m model) string {
	if m.addingCommand {
		// Show add command form
		var s strings.Builder

		addTitleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(pinkColor).
			Padding(0, 0, 1, 0)

		s.WriteString(addTitleStyle.Render("➕ Add Custom Command"))
		s.WriteString("\n\n")

		// Form fields with labels
		labelStyle := lipgloss.NewStyle().Foreground(purpleColor).Bold(true)

		s.WriteString(labelStyle.Render("Name:") + "\n")
		s.WriteString(m.commandNameInput.View() + "\n\n")

		s.WriteString(labelStyle.Render("Command:") + "\n")
		s.WriteString(m.commandCmdInput.View() + "\n\n")

		s.WriteString(labelStyle.Render("Description:") + "\n")
		s.WriteString(m.commandDescInput.View() + "\n\n")

		s.WriteString(helpStyle.Render("tab: next field • enter: save • esc: cancel"))

		// Full-height container
		fullHeightContainer := lipgloss.NewStyle().
			Width(m.width - 2).
			Height(m.height - 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleColor)

		return fullHeightContainer.Render(s.String())
	}

	// Show list view
	var s strings.Builder

	listTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(pinkColor).
		Padding(0, 0, 1, 0)

	s.WriteString(listTitleStyle.Render("⚡ Custom Commands"))
	s.WriteString("\n\n")

	if len(m.workspaceCommands) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(grayColor).
			Align(lipgloss.Center).
			Width(60)
		s.WriteString(emptyStyle.Render("No custom commands yet"))
		s.WriteString("\n\n")
		s.WriteString(emptyStyle.Render("Press 'a' to add your first command"))
		s.WriteString("\n\n")
	} else {
		// Render list
		s.WriteString(m.commandsList.View())
		s.WriteString("\n\n")
	}

	// Help text
	helpText := "a: add • d: delete • enter: select • esc: back"
	if len(m.workspaceCommands) == 0 {
		helpText = "a: add command • esc: back"
	}
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
