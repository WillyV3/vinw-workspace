# vinw-workspace Build Plan

## What We're Building
Simple TUI launcher for tmux sessions with vinw. One fixed layout, minimal config.

## The Layout (Fixed)
```
┌─────────┬──────────────────────────────────┐
│         │  vinw-viewer (file preview)      │
│  vinw   │                                  │
│  tree   │                                  │
│         ├──────────────┬───────────────────┤
│         │  terminal/   │  coding agent     │
│         │  nextui      │  (user picks)     │
└─────────┴──────────────┴───────────────────┘

Windowgram:
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc
```

## TUI Flow
1. User inputs project directory
2. User inputs session name
3. User picks terminal option (shell OR nextui)
4. User picks agent option (claude, opencode, crush, codex, OR none)
5. Show preview → Launch button
6. Generate tmuxomatic file → Launch tmux

## File Structure
```
vinw-workspace/
├── main.go              # TUI entry point (ONE FILE - allowed >100 lines)
├── deps.go              # Check dependencies (vinw, tmux, agents)
├── template.go          # Generate .tmo files
├── go.mod
├── go.sum
├── README.md
└── PLAN.md             # This file
```

## Component Strategy (From Examples)

### Base: result/main.go
- Simple radio button selection
- Cursor navigation with j/k/up/down
- Clean exit on choice

### Add: textinputs/main.go patterns
- Directory input (textinput.Model)
- Session name input (textinput.Model)
- Focus management between fields

### Add: window-size/main.go patterns
- Handle tea.WindowSizeMsg
- Responsive layout (lipgloss)

### NO custom components - Use ONLY:
- `textinput.Model` - For directory/session inputs
- Radio rendering (like result.go) - For terminal/agent selection
- `lipgloss` - For layout and styling
- Standard Update/View pattern

## Model Structure
```go
type model struct {
    width, height    int              // Window size
    focusIndex       int              // Which field is focused
    inputs           []textinput.Model // [0]=directory, [1]=session
    terminalCursor   int              // Radio selection 0=shell, 1=nextui
    agentCursor      int              // Radio selection (claude/opencode/crush/codex/none)
    state            string           // "input" or "preview"
}
```

## Radio Options
```go
var terminalOptions = []string{"shell", "nextui"}
var agentOptions = []string{"claude", "opencode", "crush", "codex", "none"}
```

## Key Bindings
- `Tab` / `Shift+Tab` - Navigate between fields
- `Up/Down` or `j/k` - Navigate radio options
- `Enter` - Advance to next field / Launch
- `Esc` - Back / Cancel
- `Ctrl+C` - Quit

## Dependencies (deps.go)
Check for:
1. `vinw` binary (required)
2. `tmux` binary (required)
3. `tmuxomatic` (python package - warn if missing)
4. `nextui` binary (optional - only if selected)
5. Agent binaries (optional - only if selected)

## Template Generation (template.go)
```go
func generateTmuxomaticFile(dir, session, terminal, agent string) string {
    // Returns tmuxomatic file content as string
    // Saves to ~/.vinw-workspace/{session}.tmo
}
```

Generated file example:
```
window myproject_vinw

aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaavvvvvvvvvvvvvvvvvvvvvvvvv
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc
aaaattttttttttttttccccccccccc

  dir ~/code/myproject
a run vinw
v run vinw-viewer <session-id>
t run nextui               # OR empty if shell
c run opencode             # OR empty if none
v foc
```

## Launch Flow (main.go)
```go
1. Generate session ID (like vinw does)
2. Write .tmo file to ~/.vinw-workspace/
3. Execute: tmuxomatic ~/.vinw-workspace/{session}.tmo
4. Exit TUI
```

## Layout Pattern (Responsive)
Using lipgloss.JoinVertical and JoinHorizontal:
```go
func (m model) View() string {
    if m.state == "preview" {
        return renderPreview(m)
    }

    header := renderHeader()
    inputs := renderInputs(m)
    radios := renderRadios(m)
    footer := renderFooter()

    // Adjust based on m.width, m.height
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header, inputs, radios, footer,
    )
}
```

## Preview Screen
Show:
1. Generated windowgram (ASCII art)
2. All configurations:
   - Directory: ~/code/myproject
   - Session: myproject_vinw
   - Terminal: nextui
   - Agent: opencode
3. [Launch Session] button
4. Dependencies status (✓ or ✗)

## Implementation Order
1. ✅ Create directory structure
2. ✅ Copy starter files (result.go, textinputs.go)
3. Create deps.go (dependency checking)
4. Create template.go (tmuxomatic file generation)
5. Create main.go:
   - Basic structure from result.go
   - Add textinput fields from textinputs.go
   - Add radio selections (terminal, agent)
   - Add window size handling
   - Add preview state
   - Add launch logic
6. Create go.mod
7. Test with vinw
8. Create README.md
9. Create release.sh (copy from vinw)

## Code Standards
- Files <100 lines EXCEPT main.go (allowed up to ~200)
- No custom bubble components
- Use only patterns from examples
- Self-documenting code
- Minimal nesting
- Dynamic list delegates based on window size

## Homebrew Distribution
```ruby
class VinwWorkspace < Formula
  desc "Opinionated tmux workspace launcher with vinw"
  homepage "https://github.com/willyv3/vinw-workspace"

  depends_on "willyv3/tap/vinw"
  depends_on "tmux"
  depends_on "willyv3/tap/nextui" => :optional

  def install
    system "go", "build", "-o", "vinw-workspace"
    bin.install "vinw-workspace"
  end

  def caveats
    <<~EOS
      Requires tmuxomatic:
        pip3 install tmuxomatic
    EOS
  end
end
```

## Ready to Build
All patterns identified. No unknowns. Clean, simple, focused on the one task: launch vinw workspace in tmux.
