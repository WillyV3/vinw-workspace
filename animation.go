package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	animationFPS     = 60
	animationFrames  = 90 // 1.5 seconds at 60fps for smooth entrance
)

//go:embed ascii/ascii-art.txt
var asciiArtRaw string

//go:embed ascii/choos-dir-ascii.txt
var chooseDirAscii string

//go:embed ascii/config-wokrspace-ascii.txt
var configWorkspaceAscii string

//go:embed ascii/preview-ascii.txt
var previewAscii string

//go:embed ascii/tmux-noobs.txt
var tmuxNoobsAscii string

// ASCII art split into lines for easy manipulation
var asciiArtLines []string

func init() {
	// Load ASCII art lines on package initialization
	asciiArtLines = strings.Split(strings.TrimSpace(asciiArtRaw), "\n")
}

// animationMsg is sent on each frame tick for smooth animation
type animationMsg struct{}

// tickAnimation creates a command that sends animation frames at 60 FPS
func tickAnimation() tea.Cmd {
	return tea.Tick(time.Second/animationFPS, func(time.Time) tea.Msg {
		return animationMsg{}
	})
}

// makeGradientRamp creates a smooth color gradient between two colors
// Returns an array of lipgloss styles, one per step
func makeGradientRamp(colorStart, colorEnd string, steps int) []lipgloss.Style {
	start, _ := colorful.Hex(colorStart)
	end, _ := colorful.Hex(colorEnd)

	styles := make([]lipgloss.Style, steps)
	for i := 0; i < steps; i++ {
		// Blend from start to end color
		t := float64(i) / float64(steps-1)
		color := start.BlendLuv(end, t)
		hexColor := fmt.Sprintf("#%02x%02x%02x",
			uint8(color.R*255),
			uint8(color.G*255),
			uint8(color.B*255))
		styles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor))
	}
	return styles
}

// easeOutCubic provides smooth deceleration
// t goes from 0.0 to 1.0, returns eased value
func easeOutCubic(t float64) float64 {
	t--
	return t*t*t + 1
}

// renderAnimatedTitle renders the ASCII art with gradient and reveal animation
// frame: current animation frame (0 to animationFrames)
func renderAnimatedTitle(frame int) string {
	if len(asciiArtLines) == 0 {
		return ""
	}

	// Calculate progress with easing (0.0 to 1.0)
	progress := float64(frame) / float64(animationFrames)
	if progress > 1.0 {
		progress = 1.0
	}
	easedProgress := easeOutCubic(progress)

	// Get max line length for gradient
	maxLineLen := 0
	for _, line := range asciiArtLines {
		if len(line) > maxLineLen {
			maxLineLen = len(line)
		}
	}

	// Create gradient from purple to cyan
	gradient := makeGradientRamp("#B14FFF", "#00FFA3", maxLineLen)

	var result strings.Builder

	// Fade in opacity (faster than character reveal)
	fadeProgress := progress * 1.5
	if fadeProgress > 1.0 {
		fadeProgress = 1.0
	}

	// Render each line with progressive character reveal
	for _, line := range asciiArtLines {
		// Calculate how many characters to show (left-to-right reveal)
		charsToShow := int(float64(len(line)) * easedProgress)

		// Render only the revealed portion of this line
		for i, char := range line {
			if i >= charsToShow {
				break // Don't render unrevealed characters
			}

			if i < len(gradient) {
				if fadeProgress < 1.0 {
					// Dim by blending with black
					fgColor := fmt.Sprintf("%v", gradient[i].GetForeground())
					baseColor, _ := colorful.Hex(fgColor)
					dimmed := baseColor.BlendLuv(colorful.Color{R: 0, G: 0, B: 0}, 1.0-fadeProgress)
					dimHex := fmt.Sprintf("#%02x%02x%02x",
						uint8(dimmed.R*255),
						uint8(dimmed.G*255),
						uint8(dimmed.B*255))
					dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(dimHex))
					result.WriteString(dimStyle.Render(string(char)))
				} else {
					// Full brightness
					result.WriteString(gradient[i].Render(string(char)))
				}
			} else {
				result.WriteRune(char)
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderStaticGradientTitle renders ASCII art with gradient colors (no animation)
func renderStaticGradientTitle(asciiText string) string {
	if asciiText == "" {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(asciiText), "\n")
	if len(lines) == 0 {
		return ""
	}

	// Get max line length in runes (characters) for gradient
	maxLineLen := 0
	for _, line := range lines {
		runeCount := len([]rune(line))
		if runeCount > maxLineLen {
			maxLineLen = runeCount
		}
	}

	// Create gradient from purple to cyan
	gradient := makeGradientRamp("#B14FFF", "#00FFA3", maxLineLen)

	var result strings.Builder

	// Render each line with gradient (use rune index, not byte index)
	for _, line := range lines {
		charIndex := 0
		for _, char := range line {
			if charIndex < len(gradient) {
				result.WriteString(gradient[charIndex].Render(string(char)))
			} else {
				result.WriteRune(char)
			}
			charIndex++
		}
		result.WriteString("\n")
	}

	return result.String()
}
