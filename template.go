package main

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"
)

func generateSessionID(path string) string {
	absPath, _ := filepath.Abs(path)
	hash := sha256.Sum256([]byte(absPath))
	return fmt.Sprintf("%x", hash[:4])
}

func getPreviewContent(dir, session, terminal, agent, customCmd string) string {
	absDir, _ := filepath.Abs(dir)
	sessionID := generateSessionID(absDir)

	var terminalDisplay, agentDisplay string
	if customCmd != "" {
		terminalDisplay = fmt.Sprintf("custom: %s", customCmd)
	} else if terminal == "nextui" {
		terminalDisplay = "nextui (Next.js scaffolder)"
	} else {
		terminalDisplay = "shell (empty terminal)"
	}

	if agent == "none" {
		agentDisplay = "none (empty terminal)"
	} else {
		agentDisplay = agent
	}

	windowgram := `aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc`

	var sb strings.Builder
	sb.WriteString("Layout Preview\n\n")
	sb.WriteString(windowgram)
	sb.WriteString("\n\n")
	sb.WriteString("Configuration\n")
	sb.WriteString(fmt.Sprintf("  Directory:  %s\n", absDir))
	sb.WriteString(fmt.Sprintf("  Session:    %s\n", session))
	sb.WriteString(fmt.Sprintf("  Session ID: %s\n", sessionID))
	sb.WriteString(fmt.Sprintf("  Terminal:   %s\n", terminalDisplay))
	sb.WriteString(fmt.Sprintf("  Agent:      %s\n", agentDisplay))
	sb.WriteString("\n")
	sb.WriteString("Pane Layout\n")
	sb.WriteString("  [a] vinw file browser\n")
	sb.WriteString(fmt.Sprintf("  [v] vinw-viewer %s\n", sessionID))
	sb.WriteString(fmt.Sprintf("  [t] %s\n", terminalDisplay))
	sb.WriteString(fmt.Sprintf("  [c] %s\n", agentDisplay))

	return sb.String()
}
