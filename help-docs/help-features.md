# Feature Deep Dive

Learn how to use each feature of your tmux configuration.

## Catppuccin Theme

**What it is:** A soothing pastel theme for tmux with excellent readability

**Features:**
- Consistent color palette across panes
- Easy-to-read status bar
- Distinct active/inactive pane indicators
- Matches other Catppuccin-themed tools

**Customization:**
Edit the theme settings in your config:
```bash
# In ~/.tmux.conf
set -g @catppuccin_flavor 'mocha'  # or latte, frappe, macchiato
```

## Tmux Logging

**Activate logging:** `prefix + shift + p`
**Stop logging:** `prefix + shift + p` (again)
**Save complete history:** `prefix + alt + shift + p`

**Use cases:**
- Debugging build processes
- Capturing terminal sessions for documentation
- Saving important command output
- Creating tutorials

**Where logs are saved:**
- Default: `~/tmux-${session}-${window}-${pane}.log`
- Timestamped and unique per pane

## Tmux Autoreload

**What it does:** Watches ~/.tmux.conf for changes and reloads automatically

**How it works:**
1. Edit your tmux.conf
2. Save the file
3. Config reloads instantly
4. No manual `prefix + r` needed

**Perfect for:**
- Testing new key bindings
- Tweaking status bar settings
- Experimenting with plugins
- Learning tmux configuration

## Extrakto

**Activate:** `prefix + tab`

**What it does:**
- Scans current pane for extractable text
- Presents fuzzy-searchable list
- Copy selected item to clipboard

**Extracts:**
- URLs (http://, https://)
- File paths (/home/user/file.txt)
- Git commit hashes
- IP addresses
- Email addresses
- Custom patterns

**Workflow:**
1. `prefix + tab` to activate
2. Type to filter results
3. Arrow keys to select
4. Enter to copy to clipboard

## Tmux Wizard Integration

**Launch:** `prefix + q`

**What it provides:**
- Interactive session browser
- Visual pane layouts
- Quick workspace creation
- Integrated with this vinw tool

**Benefits:**
- No need to remember session names
- Visual preview of layouts
- Faster workflow switching
- Mouse-friendly interface

## Mouse Mode Features

With mouse mode enabled, you can:

**Select Panes:**
- Click any pane to focus it
- No keyboard shortcuts needed

**Resize Panes:**
- Click and drag pane borders
- Visual feedback as you resize
- Precise control

**Scroll History:**
- Mouse wheel to scroll back
- Works in any pane
- Natural scrolling

**Select Text:**
- Click and drag to select
- Double-click to select word
- Triple-click to select line

**Copy Text:**
- Select text with mouse
- Automatically copies to tmux buffer
- Paste with `prefix + ]`

## Custom Key Bindings

### Reload Config
```
bind-key r source-file ~/.tmux.conf \; display-message "~/.tmux.conf reloaded."
```
- Quick feedback message
- Instant config changes
- No session restart

### Wizard Popup
```
bind-key q display-popup -E -w 70% -h 60% ...
```
- 70% of terminal width
- 60% of terminal height
- Centered popup
- Closes on exit

## Status Bar Customization

**Current Format:**
```
[Time] [Date]
```

**Customize:**
```bash
# Add more info:
set -g status-right '#[fg=white]%a %l:%M:%S %p #[fg=blue]%Y-%m-%d #[fg=green]#H'

# Show hostname, CPU, memory:
set -g status-right '#[fg=yellow]#(hostname) #[fg=green]CPU:#(top -l 1 | grep "CPU usage" | awk "{print \$3}")'
```

**Refresh Rate:**
```bash
set -g status-interval 1  # Update every second
```

## Plugin Management

### Install New Plugins

1. Add to ~/.tmux.conf:
```bash
set -g @plugin 'user/plugin-name'
```

2. Reload config: `prefix + r`

3. Install plugin: `prefix + I` (capital i)

### Update Plugins

Press `prefix + U` to update all plugins

### Remove Plugins

1. Remove from ~/.tmux.conf
2. Reload config
3. Press `prefix + alt + u` to uninstall

## Advanced Tips

**Create Custom Bindings:**
```bash
# Split panes with current directory
bind '"' split-window -c "#{pane_current_path}"
bind % split-window -h -c "#{pane_current_path}"

# Quick pane switching
bind -n M-Left select-pane -L
bind -n M-Right select-pane -R
bind -n M-Up select-pane -U
bind -n M-Down select-pane -D
```

**Session Management:**
```bash
# Save and restore sessions
set -g @plugin 'tmux-plugins/tmux-resurrect'
set -g @plugin 'tmux-plugins/tmux-continuum'
```

Press **Tab** to learn tmux basics!
