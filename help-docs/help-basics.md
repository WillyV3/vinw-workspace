# Tmux Basics

Everything you need to know to get started with tmux.

## What is Tmux?

**Tmux** (Terminal Multiplexer) lets you:
- Run multiple terminal sessions in one window
- Split terminals horizontally and vertically
- Detach and reattach to sessions
- Keep processes running when you disconnect
- Organize your workflow into sessions

**Why use tmux?**
- Work on multiple projects simultaneously
- Never lose your terminal state
- SSH into servers and keep sessions alive
- Pair program with others
- Create complex terminal layouts

## Core Concepts

### Sessions
A **session** is a collection of windows and panes
- Think of it as a workspace for a project
- Survives disconnects and terminal closures
- Can have multiple attached clients

**Commands:**
```bash
tmux new -s myproject          # Create session named "myproject"
tmux ls                         # List all sessions
tmux attach -t myproject        # Attach to session
tmux kill-session -t myproject  # Kill session
```

### Windows
A **window** is like a tab in a browser
- Each session has one or more windows
- Switch between them with `prefix + number`
- Rename with `prefix + ,`

**Commands:**
```bash
prefix + c      # Create new window
prefix + ,      # Rename window
prefix + w      # List windows
prefix + n      # Next window
prefix + p      # Previous window
prefix + 1-9    # Switch to window number
```

### Panes
A **pane** is a split section of a window
- Divide screen horizontally or vertically
- Each pane runs independently
- Resize and rearrange as needed

**Commands:**
```bash
prefix + "      # Split horizontal
prefix + %      # Split vertical
prefix + arrow  # Navigate panes
prefix + x      # Close pane
prefix + z      # Zoom/unzoom pane
prefix + space  # Cycle layouts
```

## The Prefix Key

**Default prefix:** `Ctrl + b`

**Why a prefix?**
- Prevents conflicts with shell shortcuts
- Indicates you're talking to tmux, not the shell
- Press prefix, then the command key

**Example:**
To create a new window:
1. Press `Ctrl + b` (prefix)
2. Release both keys
3. Press `c`

**Change prefix (optional):**
```bash
# In ~/.tmux.conf
unbind C-b
set -g prefix C-a
bind C-a send-prefix
```

## Configuration Files

### Where Configs Live

**System-wide:** `/etc/tmux.conf`
- Affects all users
- Rarely modified

**User config:** `~/.tmux.conf`
- Your personal settings
- This is what we installed
- Takes precedence over system config

### How Configs Are Loaded

1. Tmux starts
2. Reads `/etc/tmux.conf` (if exists)
3. Reads `~/.tmux.conf` (if exists)
4. User config overrides system config

### Reload Config Without Restart

**Method 1:** Key binding
```
prefix + r
```

**Method 2:** Command line
```bash
tmux source-file ~/.tmux.conf
```

**Method 3:** Inside tmux
```bash
:source-file ~/.tmux.conf
```

## Essential Commands

### Session Management
```bash
# Create new session
tmux new -s work

# Detach from session
prefix + d

# List sessions
tmux ls

# Attach to last session
tmux attach

# Attach to specific session
tmux attach -t work

# Kill session
tmux kill-session -t work

# Rename session
prefix + $
```

### Window Management
```bash
# New window
prefix + c

# Close window
prefix + &

# Next/previous
prefix + n
prefix + p

# Switch by number
prefix + 0-9

# Rename window
prefix + ,

# Find window
prefix + f
```

### Pane Management
```bash
# Split horizontal
prefix + "

# Split vertical
prefix + %

# Navigate
prefix + arrow keys

# Resize
prefix + Ctrl + arrow keys

# Close pane
prefix + x

# Zoom pane (fullscreen)
prefix + z

# Cycle layouts
prefix + space
```

### Copy Mode
```bash
# Enter copy mode
prefix + [

# Navigate
arrow keys, Page Up/Down

# Start selection
space

# Copy selection
Enter

# Paste
prefix + ]

# Exit copy mode
q or Esc
```

## Common Workflows

### Development Workflow
```bash
# Create session for project
tmux new -s myapp

# Split into 3 panes
prefix + "      # Split horizontal
prefix + %      # Split vertical on bottom

# Layout:
# ┌─────────────┐
# │   Editor    │
# ├──────┬──────┤
# │ Logs │ Shell│
# └──────┴──────┘

# Pane 1: vim/code editor
# Pane 2: tail -f logs
# Pane 3: git commands

# Detach when done
prefix + d

# Reattach later
tmux attach -t myapp
```

### Server Monitoring
```bash
# Create monitoring session
tmux new -s monitor

# 4-pane layout
prefix + "      # Split horizontal
prefix + %      # Split current pane vertical
prefix + ↑      # Move to top pane
prefix + %      # Split vertical

# Layout:
# ┌──────┬──────┐
# │ htop │ logs │
# ├──────┼──────┤
# │ tail │ shell│
# └──────┴──────┘
```

### Pair Programming
```bash
# Person 1 creates session
tmux new -s pair

# Person 2 attaches
tmux attach -t pair

# Both see same screen
# Both can type
# Perfect for remote pairing
```

## Tips & Tricks

**Quickly swap panes:**
```
prefix + {      # Swap with previous
prefix + }      # Swap with next
```

**Synchronize panes:**
```
prefix + :
setw synchronize-panes on
```
Now typing in one pane types in all!

**Break pane to new window:**
```
prefix + !
```

**Join pane from another window:**
```
prefix + :
join-pane -s :2
```

**Rename session:**
```
prefix + $
```

**Clock mode:**
```
prefix + t
```

**Show all key bindings:**
```
prefix + ?
```

## Configuration Best Practices

**1. Use comments**
```bash
# This reloads the config
bind-key r source-file ~/.tmux.conf
```

**2. Group related settings**
```bash
# Color settings
set -g default-terminal "tmux-256color"
set -g status-style bg=black,fg=white
```

**3. Test before committing**
- Make small changes
- Reload and test
- Keep backup of working config

**4. Use plugins wisely**
- Only install what you need
- Too many plugins slow startup
- Keep it simple

## Troubleshooting

**Config not loading?**
```bash
# Check for syntax errors
tmux source-file ~/.tmux.conf
```

**Keybinding not working?**
```bash
# List all bindings
tmux list-keys
```

**Colors look wrong?**
```bash
# Verify terminal support
echo $TERM
```

**Plugin not working?**
```bash
# Reinstall TPM
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
```

## Next Steps

1. **Practice:** Use tmux daily for a week
2. **Customize:** Add your own key bindings
3. **Explore:** Try new plugins
4. **Share:** Show teammates your workflow

Press **Tab** to search the full tmux man page!
