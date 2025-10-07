# Install tmux

Get tmux installed on your system using your package manager.

---

## Quick Installation

This tool will help you install tmux using your system's package manager. When you press **Enter**, you'll be able to select your installation method from the available options.

**Supported installation methods:**
- **Homebrew** (macOS)
- **apt** (Ubuntu/Debian)
- **dnf** (Fedora/RHEL 8+)
- **yum** (CentOS/RHEL 7)
- **pacman** (Arch Linux)

Only methods available on your current system will be shown.

---

## What Gets Installed

When you install tmux, you get:

**Core functionality:**
- Terminal multiplexer binary (`tmux`)
- Man pages and documentation
- Default configuration

**Version:**
- Most package managers install the latest stable release
- Some distributions may have slightly older versions
- Check installed version: `tmux -V`

---

## After Installation

Once tmux is installed:

1. **Verify installation:**
   ```bash
   tmux -V
   ```

2. **Install configuration:**
   - Go back to menu
   - Select "Install optimized .tmux.conf"

3. **Start using tmux:**
   ```bash
   tmux new -s mysession
   ```

4. **Learn the basics:**
   - Select "Tmux Basics" from the menu
   - Press `prefix + ?` to see all keybindings

---

## Manual Installation

For more installation options or to build from source, visit:

**Official tmux wiki:**
https://github.com/tmux/tmux/wiki/Installing

**Common manual installations:**

### macOS (from source)
```bash
brew install libevent ncurses
git clone https://github.com/tmux/tmux.git
cd tmux
sh autogen.sh
./configure && make
sudo make install
```

### Linux (from source)
```bash
# Install dependencies (Ubuntu/Debian)
sudo apt-get install libevent-dev ncurses-dev build-essential

# Clone and build
git clone https://github.com/tmux/tmux.git
cd tmux
sh autogen.sh
./configure && make
sudo make install
```

### FreeBSD
```bash
pkg install tmux
```

### NetBSD
```bash
pkgin install tmux
```

---

## Troubleshooting

**Command not found after install:**
- Close and reopen your terminal
- Check PATH: `echo $PATH`
- Verify installation: `which tmux`

**Permission denied:**
- Installation requires sudo/root access
- Ensure your user has sudo privileges

**Old version installed:**
- Some distributions ship older versions
- Consider installing from source for latest
- Or use alternative package managers (Homebrew on Linux)

**Already installed:**
- This tool will detect existing installations
- You can reinstall to update version
- Or skip to configuring tmux

---

Press **Tab** to return to menu and begin installation!
