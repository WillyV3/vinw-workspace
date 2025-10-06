package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) loadDirectory(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	m.files = []fileEntry{}
	m.directory = path

	// Add parent directory option if not root
	if path != "/" && path != filepath.Dir(path) {
		parent := filepath.Dir(path)
		m.files = append(m.files, fileEntry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		})
	}

	// Sort directories first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		m.files = append(m.files, fileEntry{
			Name:  entry.Name(),
			Path:  filepath.Join(path, entry.Name()),
			IsDir: entry.IsDir(),
		})
	}

	m.filteredFiles = m.files
	m.cursor = 0
	m.searching = false
	m.searchInput.Blur()
	m.searchInput.SetValue("")
	m.updateViewport()
}

func fuzzyMatch(pattern, text string) bool {
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	if pattern == "" {
		return true
	}

	patternIdx := 0
	for _, char := range text {
		if patternIdx < len(pattern) && rune(pattern[patternIdx]) == char {
			patternIdx++
		}
	}

	return patternIdx == len(pattern)
}

func (m *model) filterFiles() {
	query := m.searchInput.Value()
	if query == "" {
		m.filteredFiles = m.files
	} else {
		m.filteredFiles = []fileEntry{}
		for _, file := range m.files {
			if fuzzyMatch(query, file.Name) {
				m.filteredFiles = append(m.filteredFiles, file)
			}
		}
	}

	m.cursor = 0
	m.updateViewport()
}

func (m *model) updateViewport() {
	if len(m.filteredFiles) == 0 {
		return
	}

	// Account for full-height container layout:
	// - Container height: m.height - 2 (margin for border)
	// - Border: 2 (1 top + 1 bottom)
	// - Padding: 2 (1 top + 1 bottom)
	// - Title + spacing: 3 lines
	// - Current dir + spacing: 3 lines
	// - Separators: 4 lines (2 separators + surrounding spacing)
	// - Help text: 1 line
	// Total overhead: ~15 lines
	availableHeight := m.height - 17
	if m.creatingNewDir || m.searching {
		availableHeight -= 3 // Extra lines for input
	}
	if availableHeight < 5 {
		availableHeight = 5 // Minimum visible items
	}

	if m.cursor >= len(m.filteredFiles) {
		m.cursor = len(m.filteredFiles) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor < m.viewportStart {
		m.viewportStart = m.cursor
	}
	if m.cursor >= m.viewportStart+availableHeight {
		m.viewportStart = m.cursor - availableHeight + 1
	}

	if m.viewportStart < 0 {
		m.viewportStart = 0
	}

	m.viewportEnd = m.viewportStart + availableHeight
	if m.viewportEnd > len(m.filteredFiles) {
		m.viewportEnd = len(m.filteredFiles)
	}
}

func (m model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Directory browsing mode (when focusIndex == 0)
	if m.focusIndex == 0 {
		// Handle search mode
		if m.searching {
			switch msg.String() {
			case "enter":
				// Navigate INTO the directory (stay in browser to continue exploring)
				if m.cursor < len(m.filteredFiles) {
					entry := m.filteredFiles[m.cursor]
					if entry.IsDir {
						m.loadDirectory(entry.Path)
					}
				}
				// Exit search mode but stay in directory browser
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.updateViewport() // Restore viewport to full size
				return m, nil
			case "tab":
				// SELECT this as final directory choice and move to next field
				if m.cursor < len(m.filteredFiles) {
					entry := m.filteredFiles[m.cursor]
					if entry.IsDir {
						m.loadDirectory(entry.Path)
					}
				}
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.focusIndex = 1
				m.inputs[0].Focus()
				m.inputs[0].PromptStyle = focusedStyle
				m.inputs[0].TextStyle = focusedStyle
				m.updateViewport() // Restore viewport to full size
				return m, nil
			case "esc":
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.filterFiles()
				m.updateViewport() // Restore viewport to full size
				return m, nil
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.updateViewport()
				}
				return m, nil
			case "down", "j":
				if m.cursor < len(m.filteredFiles)-1 {
					m.cursor++
					m.updateViewport()
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterFiles()
				return m, cmd
			}
		} else if m.creatingNewDir {
			switch msg.String() {
			case "enter":
				dirName := strings.TrimSpace(m.newDirInput.Value())
				if dirName != "" {
					fullPath := filepath.Join(m.directory, dirName)
					err := os.MkdirAll(fullPath, 0755)
					if err == nil {
						// Create new directory and load it
						m.loadDirectory(fullPath)
						m.creatingNewDir = false
						m.newDirInput.Blur()
						m.newDirInput.SetValue("")
						m.focusIndex = 1
						m.inputs[0].Focus()
						m.inputs[0].PromptStyle = focusedStyle
						m.inputs[0].TextStyle = focusedStyle
						return m, nil
					}
				}
				return m, nil
			case "esc":
				m.creatingNewDir = false
				m.newDirInput.Blur()
				m.newDirInput.SetValue("")
				m.updateViewport() // Restore viewport to full size
				return m, nil
			default:
				var cmd tea.Cmd
				m.newDirInput, cmd = m.newDirInput.Update(msg)
				return m, cmd
			}
		} else {
			// Normal directory navigation
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				return m, tea.Quit
			case "enter", "tab":
				// Select highlighted directory and load its contents
				if m.cursor < len(m.filteredFiles) {
					entry := m.filteredFiles[m.cursor]
					if entry.IsDir {
						m.loadDirectory(entry.Path)
					}
				}
				// Move to next field
				m.focusIndex = 1
				m.inputs[0].Focus()
				m.inputs[0].PromptStyle = focusedStyle
				m.inputs[0].TextStyle = focusedStyle
				return m, nil
			case "right":
				// Navigate INTO selected directory
				if m.cursor < len(m.filteredFiles) {
					entry := m.filteredFiles[m.cursor]
					if entry.IsDir {
						m.loadDirectory(entry.Path)
					}
				}
				return m, nil
			case "left":
				// Navigate UP to parent directory
				parent := filepath.Dir(m.directory)
				if parent != m.directory {
					m.loadDirectory(parent)
				}
				return m, nil
			case "n":
				m.creatingNewDir = true
				m.newDirInput.Focus()
				m.updateViewport() // Recalculate viewport with reduced space
				return m, nil
			case "s", " ":
				m.searching = true
				m.searchInput.Focus()
				m.updateViewport() // Recalculate viewport with reduced space
				return m, nil
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.updateViewport()
				}
				return m, nil
			case "down", "j":
				if m.cursor < len(m.filteredFiles)-1 {
					m.cursor++
					m.updateViewport()
				}
				return m, nil
			}
		}
		return m, nil
	}

	// Rest of the form (session name, terminal, agent, button)

	// When session input is focused, block all navigation except esc
	if m.focusIndex == 1 && m.inputs[0].Focused() {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Blur the input, go back to directory browser
			m.focusIndex = 0
			m.inputs[0].Blur()
			m.inputs[0].PromptStyle = noStyle
			m.inputs[0].TextStyle = noStyle
			m.loadDirectory(m.directory)
			return m, nil
		case "enter":
			// Move to next field (terminal selection)
			m.inputs[0].Blur()
			m.inputs[0].PromptStyle = noStyle
			m.inputs[0].TextStyle = noStyle
			m.focusIndex = 2
			return m, nil
		default:
			// Allow only text input, block all navigation
			var cmd tea.Cmd
			m.inputs[0], cmd = m.inputs[0].Update(msg)
			return m, cmd
		}
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.focusIndex > 0 {
			// Go back to directory browser and reload current directory
			m.focusIndex = 0
			m.inputs[0].Blur()
			m.inputs[0].PromptStyle = noStyle
			m.inputs[0].TextStyle = noStyle
			// Reload the directory to reset browser state
			m.loadDirectory(m.directory)
			return m, nil
		}
		return m, tea.Quit

	case "c":
		// Open custom commands list (only when not in directory browser)
		if m.focusIndex > 0 {
			m.currentState = stateCommands
			// Update list height to fit terminal
			m.commandsList.SetSize(m.width-10, m.height-15)
			return m, nil
		}

	case "tab", "shift+tab":
		s := msg.String()
		maxIndex := len(m.inputs) + 2

		if s == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}

		if m.focusIndex > maxIndex {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = maxIndex
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i <= len(m.inputs)-1; i++ {
			if i == m.focusIndex {
				cmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
			} else {
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
		}

		return m, tea.Batch(cmds...)

	case "left", "right":
		if m.focusIndex == 2 {
			if msg.String() == "right" {
				m.terminalCursor++
				if m.terminalCursor >= len(m.terminalOptions) {
					m.terminalCursor = 0
				}
			} else {
				m.terminalCursor--
				if m.terminalCursor < 0 {
					m.terminalCursor = len(m.terminalOptions) - 1
				}
			}
		} else if m.focusIndex == 3 {
			if msg.String() == "right" {
				m.agentCursor++
				if m.agentCursor >= len(m.agentOptions) {
					m.agentCursor = 0
				}
			} else {
				m.agentCursor--
				if m.agentCursor < 0 {
					m.agentCursor = len(m.agentOptions) - 1
				}
			}
		}
		return m, nil

	case "enter":
		if m.focusIndex == 4 {
			m.currentState = statePreview
			return m, nil
		}

		m.focusIndex++
		if m.focusIndex > 4 {
			m.focusIndex = 0
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i <= len(m.inputs)-1; i++ {
			if i == m.focusIndex {
				cmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
			} else {
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
		}
		return m, tea.Batch(cmds...)
	}

	// Note: Session input update is handled in the focused block above
	return m, nil
}
