package main

import (
	"fmt"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
)

func isInTmux() bool {
	return os.Getenv("TMUX") != ""
}

func sessionExists(session string) bool {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return false
	}
	return tmux.HasSession(session)
}

func getTmuxSessions() []string {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return []string{}
	}
	sessions, err := tmux.ListSessions()
	if err != nil {
		return []string{}
	}
	names := make([]string, len(sessions))
	for i, s := range sessions {
		names[i] = s.Name
	}
	return names
}

func launchTmuxSession(dir, session, terminal, agent, sessionID, customCmd string) error {
	absDir := os.ExpandEnv(dir)

	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return fmt.Errorf("failed to initialize tmux: %w", err)
	}

	if tmux.HasSession(session) {
		return fmt.Errorf("session '%s' already exists - choose a different name", session)
	}

	// Create detached session with starting directory
	sess, err := tmux.NewSession(&gotmux.SessionOptions{
		Name:           session,
		StartDirectory: absDir,
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Get the auto-created window and pane
	windows, err := sess.ListWindows()
	if err != nil {
		return fmt.Errorf("failed to list windows: %w", err)
	}
	if len(windows) == 0 {
		return fmt.Errorf("no windows found in session")
	}
	window := windows[0]

	panes, err := window.ListPanes()
	if err != nil {
		return fmt.Errorf("failed to list panes: %w", err)
	}
	if len(panes) == 0 {
		return fmt.Errorf("no panes found in window")
	}
	pane0 := panes[0]

	// Pane 0 (left): vinw - start the command
	_, err = tmux.Command("send-keys", "-R", "-t", pane0.Id, fmt.Sprintf("cd %s && vinw", absDir), "Enter")
	if err != nil {
		return fmt.Errorf("failed to start vinw: %w", err)
	}

	// Split horizontally - right pane will be larger
	_, err = tmux.Command("split-window", "-h", "-c", absDir, "-t", pane0.Id)
	if err != nil {
		return fmt.Errorf("failed to create right split: %w", err)
	}

	// Resize left pane (vinw) to 43 columns
	_, err = tmux.Command("resize-pane", "-t", pane0.Id, "-x", "43")
	if err != nil {
		return fmt.Errorf("failed to resize left pane: %w", err)
	}

	// Refresh panes list and get pane 1
	panes, err = window.ListPanes()
	if err != nil {
		return fmt.Errorf("failed to list panes after first split: %w", err)
	}
	if len(panes) < 2 {
		return fmt.Errorf("expected 2 panes after first split")
	}
	pane1 := panes[1]

	// Pane 1 (top-right): vinw-viewer - 48% of right side height
	_, err = tmux.Command("send-keys", "-R", "-t", pane1.Id, fmt.Sprintf("cd %s && vinw-viewer %s", absDir, sessionID), "Enter")
	if err != nil {
		return fmt.Errorf("failed to start vinw-viewer: %w", err)
	}

	// Split vertically from pane 1 - bottom gets 52% height
	_, err = tmux.Command("split-window", "-v", "-p", "52", "-c", absDir, "-t", pane1.Id)
	if err != nil {
		return fmt.Errorf("failed to create bottom-right split: %w", err)
	}

	// Refresh panes list and get pane 2
	panes, err = window.ListPanes()
	if err != nil {
		return fmt.Errorf("failed to list panes after second split: %w", err)
	}
	if len(panes) < 3 {
		return fmt.Errorf("expected 3 panes after second split")
	}
	pane2 := panes[2]

	// Pane 2 (bottom-right top): custom command, or terminal/nextui if requested
	if customCmd != "" {
		// Run custom command
		_, err = tmux.Command("send-keys", "-R", "-t", pane2.Id, fmt.Sprintf("cd %s && %s", absDir, customCmd), "Enter")
		if err != nil {
			return fmt.Errorf("failed to start custom command: %w", err)
		}
	} else if terminal == "nextui" {
		_, err = tmux.Command("send-keys", "-R", "-t", pane2.Id, fmt.Sprintf("cd %s && nextui", absDir), "Enter")
		if err != nil {
			return fmt.Errorf("failed to start nextui: %w", err)
		}
	}

	// Split horizontally from pane 2 - 50/50 split
	_, err = tmux.Command("split-window", "-h", "-c", absDir, "-t", pane2.Id)
	if err != nil {
		return fmt.Errorf("failed to create agent pane: %w", err)
	}

	// Refresh panes list and get pane 3
	panes, err = window.ListPanes()
	if err != nil {
		return fmt.Errorf("failed to list panes after third split: %w", err)
	}
	if len(panes) < 4 {
		return fmt.Errorf("expected 4 panes after third split")
	}
	pane3 := panes[3]

	// Pane 3 (bottom-right bottom): agent if requested
	if agent != "none" && agent != "" {
		_, err = tmux.Command("send-keys", "-R", "-t", pane3.Id, fmt.Sprintf("cd %s && %s", absDir, agent), "Enter")
		if err != nil {
			return fmt.Errorf("failed to start agent: %w", err)
		}
	}

	// Focus on vinw-viewer pane
	_, err = tmux.Command("select-pane", "-t", pane1.Id)
	if err != nil {
		return fmt.Errorf("failed to select pane: %w", err)
	}

	// Attach or switch to session
	if isInTmux() {
		err = tmux.SwitchClient(&gotmux.SwitchClientOptions{
			TargetSession: session,
		})
		if err != nil {
			return fmt.Errorf("failed to switch client: %w", err)
		}
	} else {
		err = sess.Attach()
		if err != nil {
			return fmt.Errorf("failed to attach session: %w", err)
		}
	}

	return nil
}
