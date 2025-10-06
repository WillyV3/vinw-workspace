# Your Tmux Configuration

This installation includes a carefully crafted tmux configuration optimized for modern terminal workflows.

## Core Settings

### Mouse Support
```
set -g mouse on
```
**What it does:** Enables full mouse interaction with tmux
- Click to select panes
- Drag borders to resize panes
- Scroll with mouse wheel
- Click status bar elements

### 256 Color Terminal
```
set -g default-terminal "tmux-256color"
```
**What it does:** Enables full color support for themes and syntax highlighting
- Better looking status bars
- Accurate color rendering for editors
- Proper theme display

### Smart Indexing
```
set -g base-index 1
setw -g pane-base-index 1
```
**What it does:** Windows and panes start at 1 instead of 0
- More intuitive numbering (matches keyboard layout)
- Easier to remember window positions

### Increased History
```
set -g history-limit 50000
```
**What it does:** Stores 50,000 lines of scrollback per pane
- Review long build outputs
- Search through extensive logs
- Never lose important information

## Key Bindings

### Reload Configuration
```
prefix + r
```
Reloads ~/.tmux.conf without restarting tmux
- Test configuration changes instantly
- No need to kill and restart sessions

### Launch Tmux Wizard
```
prefix + q
```
Opens an interactive TUI popup for tmux management
- Manage sessions visually
- Quick navigation
- Launch integrated with your workflow

## Plugin System (TPM)

This config uses **Tmux Plugin Manager** for easy plugin management.

### Installed Plugins

**tmux-plugins/tmux-sensible**
- Sensible default settings
- Better key bindings
- Improved copy mode

**catppuccin/tmux**
- Beautiful, modern theme
- Consistent color palette
- Great readability

**tmux-plugins/tmux-logging**
- Save pane output to file
- Log complete sessions
- Capture terminal history

**b0o/tmux-autoreload**
- Automatically reload config on changes
- No manual reload needed
- Seamless development

**laktak/extrakto**
- Extract text, URLs, paths from pane
- Copy to clipboard with ease
- Keyboard-driven selection

## Status Bar

Custom status bar showing:
- Current time (updates every second)
- Current date
- Clean, minimal design
- Centered window list

## How to Customize

1. Edit `~/.tmux.conf`
2. Add your own settings or key bindings
3. Press `prefix + r` to reload
4. Changes take effect immediately

## Learn More

Press **Tab** to explore more features and tmux basics!
