# Side-by-Side Git Diff TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev)
[![CI](https://github.com/titobsala/Diffbubble/workflows/CI/badge.svg)](https://github.com/titobsala/Diffbubble/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/titobsala/Diffbubble)](https://github.com/titobsala/Diffbubble/releases)

A Terminal User Interface (TUI) application written in Go to display a side-by-side git diff with synchronized scrolling, multi-file navigation, customizable themes, search functionality, and beautiful color-coded statistics.

> **Note:** Currently at v0.3.2. Fixes Go module installation issue and includes search functionality with real-time highlighting and match navigation!

## Features

- **Multi-file navigation**: Sidebar showing all modified files with colored stats
- **Side-by-side diff display**: View old and new versions simultaneously
- **Synchronized scrolling**: Both panes scroll together for easy comparison
- **Search functionality**: Press '/' to search, real-time highlighting, navigate with 'n'/'N'
- **Customizable themes**: 9 built-in themes with interactive cycling (press 't')
- **Configuration file support**: User and per-repository config files
- **Line numbers toggle**: Show/hide line numbers with 'n' key
- **Context mode toggle**: Switch between focus mode (changes only) and full context (entire file)
- **Beautiful statistics**: Color-coded additions (green), deletions (red), and delta (yellow)
- **Focus indicators**: Visual cues show which pane is active (file list or diff)
- **Syntax highlighting**: Added lines in green, removed lines in red
- **Mouse support**: Click and scroll with your mouse
- **Keyboard navigation**: Use arrow keys, j/k, or mouse wheel

## Installation

### Option 1: Download Binary (Easiest)

Download the latest pre-built binary for your platform from the [Releases page](https://github.com/titobsala/Diffbubble/releases/latest).

**Supported platforms:**
- Linux (x86_64, ARM64)
- macOS (x86_64, ARM64/Apple Silicon)
- Windows (x86_64)

After downloading:
```sh
# Linux/macOS: Make it executable and move to PATH
chmod +x diffbubble
sudo mv diffbubble /usr/local/bin/

# Windows: Add the directory to your PATH
```

### Option 2: Install with Go

If you have Go installed, you can install diffbubble directly:

```sh
go install github.com/titobsala/Diffbubble@latest
```

**Note:** Make sure `$GOPATH/bin` (typically `~/go/bin`) is in your PATH. Add this to your shell config:
```sh
export PATH="$PATH:$HOME/go/bin"
```

Then run it from anywhere:

```sh
diffbubble
```

### Option 3: Build from Source

1.  **Prerequisites:** Ensure you have [Go](https://go.dev/doc/install) (1.25+) and [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) installed.
2.  **Clone the repository:**
    ```sh
    git clone https://github.com/titobsala/Diffbubble.git
    cd Diffbubble
    ```
3.  **Install dependencies:**
    ```sh
    go mod tidy
    ```
4.  **Build the binary:**
    ```sh
    go build -o diffbubble main.go
    ```
5.  **Run it:**
    ```sh
    ./diffbubble
    ```

### Option 4: Run without installing

For development or quick testing:

```sh
git clone https://github.com/titobsala/Diffbubble.git
cd Diffbubble
go run main.go
```

## Themes

Diffbubble comes with 9 beautiful built-in themes to customize your diff viewing experience:

- **dark** (default) - Classic dark theme with green/red diffs
- **light** - Clean light background for bright environments
- **high-contrast** - Maximum contrast for accessibility
- **solarized** - Popular Solarized Dark color scheme
- **dracula** - Dracula theme with purple accents
- **github** - GitHub-style light diffs
- **catppuccin** - Catppuccin Mocha theme
- **tokyo-night** - Tokyo Night color scheme
- **one-dark** - Atom's One Dark theme

### Using Themes

**List all available themes:**
```sh
diffbubble --list-themes
```

**Preview a theme's colors:**
```sh
diffbubble --show-theme-colors catppuccin
```

**Start with a specific theme:**
```sh
diffbubble --theme=tokyo-night
```

**Cycle through themes interactively:**
Press `t` while the app is running to cycle through all available themes in real-time!

### Theme Configuration

You can set a default theme in your configuration file:

**User config:** `~/.config/diffbubble/config.yaml`
```yaml
theme: catppuccin
```

**Repository config:** `.diffbubble.yml` (in your git repo root)
```yaml
theme: github
```

See `.config.example.yaml` for a complete configuration example.

## Usage

### Basic Usage

Navigate to any git repository with changes and run:

```sh
diffbubble
```

### CLI Options

```sh
diffbubble [flags]
```

**Available flags:**
- `--help, -h` - Show help message
- `--version, -v` - Show version information
- `--file=<filename>` - Open with specific file selected
- `--staged` - Show only staged changes (git diff --cached)
- `--unstaged` - Show only unstaged changes
- `--theme=<name>` - Set color theme (default: dark)
- `--list-themes` - List all available themes
- `--show-theme-colors <name>` - Preview colors for a specific theme

**Examples:**
```sh
# Show all changes (staged + unstaged)
diffbubble

# Show only staged changes
diffbubble --staged

# Show only unstaged changes
diffbubble --unstaged

# Open with README.md selected
diffbubble --file=README.md

# Use a specific theme
diffbubble --theme=catppuccin

# Combine flags
diffbubble --staged --file=main.go --theme=tokyo-night

# List all available themes
diffbubble --list-themes

# Preview theme colors
diffbubble --show-theme-colors dracula
```

## Controls

### Navigation
-   **File list navigation:** When file list is focused, use `j`/`k` or `↑`/`↓` to select different files
-   **Diff scrolling:** When diff is focused, use `j`/`k` or `↑`/`↓` to scroll through the diff. Both panes scroll simultaneously.
-   **Switch pane:** Press `tab` to switch focus between file list and diff panes (purple border indicates focused pane)

### Search
-   **Enter search:** Press `/` to activate search mode and type your query
-   **Execute search:** Press `Enter` to search and highlight matches
-   **Navigate matches:** Press `n` for next match, `N` for previous match
-   **Exit search:** Press `Esc` to cancel search mode
-   **Match status:** Current match shown in gold/underline, others in orange; footer shows "Match X of Y"

### Toggles
-   **Line numbers:** Press `n` to toggle line numbers on/off (or next match when search is active)
-   **Context mode:** Press `c` to toggle between focus mode (changes only) and full context (entire file)
-   **Theme cycling:** Press `t` to cycle through all available themes interactively

### General
-   **Quit:** Press `q`, `esc`, or `ctrl+c` to exit the application

### File List
The sidebar shows:
- Status icon: **M** (modified in yellow), **A** (added in green), **D** (deleted in red)
- Filename
- **+n** additions in green
- **-n** deletions in red
- **(±delta)** net change in yellow

## Acknowledgments

This project is built with the excellent TUI libraries from [Charm](https://github.com/charmbracelet):

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The Elm Architecture framework for Go
- [Bubbles](https://github.com/charmbracelet/bubbles) - Reusable TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling and layout

Special thanks to the Charm team for creating such wonderful tools for building terminal interfaces!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.