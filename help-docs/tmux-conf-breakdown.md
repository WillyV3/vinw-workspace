# tmux.conf Configuration Breakdown

This optimized tmux configuration includes:

---

## Mouse Support

```bash
set -g mouse on
```

Enables mouse support for:
- Clicking to select panes
- Scrolling through history
- Resizing panes by dragging borders
- Selecting text (hold Shift to copy to system clipboard)

## Color Support

```bash
set -g default-terminal "tmux-256color"
```

Enables 256 color support for better visuals:
- Required for modern color schemes
- Improves syntax highlighting in editors
- Better gradient support in TUIs

## Window/Pane Indexing

```bash
set -g base-index 1
setw -g pane-base-index 1
```

Start counting windows and panes at 1 instead of 0:
- More intuitive keyboard navigation (1 is left of 2)
- Matches physical keyboard layout

## Pane Titles

```bash
set -g set-titles on
set -g set-titles-string '#{pane_title}'
```

Display pane titles in terminal window title bar:
- Helps identify which tmux session is active
- Updates dynamically as you switch panes

## History Buffer

```bash
set -g history-limit 50000
```

Increase scrollback buffer to 50,000 lines:
- Default is usually 2,000 lines
- Useful for reviewing long outputs
- Can scroll back through extensive logs

## Config Reload

```bash
bind-key r source-file ~/.tmux.conf \
  \; display-message "~/.tmux.conf reloaded."
```

**Press `<prefix>r`** to reload configuration:
- Apply changes without restarting tmux
- Shows confirmation message on success

## Status Bar

```bash
set -g status-interval 1
set -g status-justify centre
set -g status-right '%a%l:%M:%S %p %Y-%m-%d'
```

Configures the bottom status bar:
- Updates every second (shows live clock)
- Window list centered
- Right side shows: Day, Time (with seconds), Date

## Plugin Manager (TPM)

```bash
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'catppuccin/tmux'
set -g @plugin 'tmux-plugins/tmux-logging'
set -g @plugin 'b0o/tmux-autoreload'
set -g @plugin 'laktak/extrakto'
```

### Plugins included:

- **tpm**: Plugin manager itself
- **tmux-sensible**: Sane defaults everyone can agree on
- **catppuccin**: Beautiful pastel color scheme
- **tmux-logging**: Save pane output to file
- **tmux-autoreload**: Auto-reload config on changes
- **extrakto**: Fuzzy find/copy text from panes

### Plugin Commands:

- **Install plugins**: `<prefix>I` (shift+i)
- **Update plugins**: `<prefix>U` (shift+u)
- **Remove plugins**: `<prefix>alt+u`

---

**Press `<prefix>?`** to see all keybindings in tmux

*Default prefix: Ctrl+b*
