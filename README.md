# vinw-workspace

Opinionated tmux workspace launcher featuring [vinw](https://github.com/willyv3/vinw) file browser.

## What is this?

A simple TUI that launches tmux sessions with a fixed, powerful layout designed for modern terminal-based development workflows. Perfect for developers transitioning from VSCode who want the power of tmux without the configuration headaches.

## Installation

```bash
brew install willyv3/tap/vinw-workspace
```

This will automatically install:
- `vinw` (file browser)
- `tmux` (terminal multiplexer)

No other dependencies required! No Python, no pip, no configuration.

### Optional Tools

```bash
# Next.js project scaffolder
brew install willyv3/tap/nextui

# Coding agents (install as needed)
# claude, opencode, crush, codex, etc.
```

## Usage

```bash
vinw-workspace
```

The TUI will guide you through:
1. Setting your project directory
2. Naming your session
3. Choosing terminal mode (shell or nextui)
4. Selecting a coding agent (or none)
5. Preview and launch

## The Layout

```
┌─────────┬──────────────────────────────────┐
│         │  vinw-viewer (file preview)      │
│  vinw   │                                  │
│  tree   │                                  │
│         ├──────────────┬───────────────────┤
│         │  terminal/   │  coding agent     │
│         │  nextui      │  (your choice)    │
└─────────┴──────────────┴───────────────────┘
```

**Panes:**
- **[a]** vinw file browser (left sidebar)
- **[v]** vinw-viewer for file preview (top right)
- **[t]** Terminal or nextui (bottom left)
- **[c]** Coding agent of your choice (bottom right)

**Features:**
- vinw and vinw-viewer are synced via session ID
- Select files in vinw → instantly previewed in viewer
- Persistent layout across sessions
- Ready-to-code environment in seconds

## Configuration

vinw-workspace stores config in `~/.vinw-workspace/`

### Adding Custom Applications

Edit `~/.vinw-workspace/config.json`:

```json
{
  "terminal_options": [
    "shell",
    "nextui",
    "my-custom-tool"
  ],
  "agent_options": [
    "claude",
    "opencode",
    "crush",
    "codex",
    "none",
    "my-agent"
  ]
}
```

The config file is automatically created with defaults on first run. Add your preferred tools to the arrays and they'll appear as options in the TUI.

### Configuration Files

vinw-workspace stores only your preferences in:
```
~/.vinw-workspace/config.json
```

Sessions are created directly with tmux - no intermediate files.

## Keyboard Shortcuts

### Input Screen
- `Tab` / `↑` / `↓` - Navigate between fields
- `j` / `k` - Select options in radio lists
- `Enter` - Next field / Preview
- `Esc` - Quit

### Preview Screen
- `Enter` or `l` - Launch tmux session
- `Esc` - Back to input
- `q` - Quit

## Use Cases

**For VSCode refugees:**
- Get a familiar file tree (vinw) + preview pane setup
- Bottom terminal split like VSCode
- Coding agent integration for AI assistance

**For tmux beginners:**
- No manual tmux configuration needed
- One fixed layout that just works
- Sessions are easily reproducible

**For power users:**
- Fast workspace setup
- Customizable via config.json
- Pure Go implementation - no Python dependencies

## How It Works

1. You configure your workspace in the TUI
2. vinw-workspace executes tmux commands directly
3. Your session is created with the fixed layout
4. vinw and vinw-viewer communicate via session ID (using Skate)
5. Each directory gets a unique session ID (deterministic hash)

## Session Management

Multiple workspaces in different directories are isolated. Each gets a unique session ID based on the directory path.

```bash
# Project A
cd ~/projects/webapp
vinw-workspace

# Project B (separate session)
cd ~/projects/api
vinw-workspace
```

Sessions don't interfere - each has its own vinw ↔ viewer connection.

## Examples

**Web development:**
```
Directory: ~/code/my-app
Session: webapp
Terminal: shell
Agent: claude
```

**Next.js project creation:**
```
Directory: ~/code/new-project
Session: nextjs
Terminal: nextui  (scaffolds the project)
Agent: opencode   (helps with code)
```

**System monitoring:**
```
Directory: ~
Session: monitor
Terminal: shell
Agent: none
```

## Troubleshooting

**"vinw not found"**
```bash
brew install willyv3/tap/vinw
```

**"Can't connect vinw to viewer"**
- Make sure both are running in the same tmux session
- Check that session IDs match (shown at launch)

**"Missing agent/tool"**
- Install the tool first
- Or select "none" / "shell" if you don't need it

## Philosophy

This tool is intentionally opinionated:
- **One layout only** - It's tested, it works, it's productive
- **Minimal config** - Just directory, session name, and tool choices
- **vinw-first** - Built around the vinw workflow
- **Terminal-native** - For developers who want to stay in the terminal

If you want full customization, use tmuxomatic directly. If you want a great default that just works, use vinw-workspace.

## Contributing

This is a personal tool made public. It's designed for a specific workflow. If it works for you, great! If not, fork it or use tmuxomatic.

## License

MIT

## Links

- [vinw](https://github.com/willyv3/vinw) - The file browser
- [nextui](https://github.com/willyv3/nextui) - Next.js scaffolder
- [tmux](https://github.com/tmux/tmux) - Terminal multiplexer
