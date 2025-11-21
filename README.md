# Side-by-Side Git Diff TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev)

A Terminal User Interface (TUI) application written in Go to display a side-by-side git diff with synchronized scrolling, multi-file navigation, and beautiful color-coded statistics.

> **Note:** This is an early release (v0.1.0). Features and APIs may change as the project evolves.

## Features

- **Multi-file navigation**: Sidebar showing all modified files with colored stats
- **Side-by-side diff display**: View old and new versions simultaneously
- **Synchronized scrolling**: Both panes scroll together for easy comparison
- **Line numbers toggle**: Show/hide line numbers with 'n' key
- **Context mode toggle**: Switch between focus mode (changes only) and full context (entire file)
- **Beautiful statistics**: Color-coded additions (green), deletions (red), and delta (yellow)
- **Focus indicators**: Visual cues show which pane is active (file list or diff)
- **Syntax highlighting**: Added lines in green, removed lines in red
- **Mouse support**: Click and scroll with your mouse
- **Keyboard navigation**: Use arrow keys, j/k, or mouse wheel

## Installation

### Option 1: Install with Go (Recommended)

If you have Go installed, you can install diffbubble directly:

```sh
go install github.com/titobsala/Diffbubble@latest
```

Then run it from anywhere:

```sh
diffbubble
```

### Option 2: Build from Source

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

### Option 3: Run without installing

For development or quick testing:

```sh
git clone https://github.com/titobsala/Diffbubble.git
cd Diffbubble
go run main.go
```

## Usage

Navigate to any git repository with changes and run:

```sh
diffbubble
```

To see available options:

```sh
diffbubble --help
```

To check the version:

```sh
diffbubble --version
```

## Controls

### Navigation
-   **File list navigation:** When file list is focused, use `j`/`k` or `↑`/`↓` to select different files
-   **Diff scrolling:** When diff is focused, use `j`/`k` or `↑`/`↓` to scroll through the diff. Both panes scroll simultaneously.
-   **Switch pane:** Press `tab` to switch focus between file list and diff panes (purple border indicates focused pane)

### Toggles
-   **Line numbers:** Press `n` to toggle line numbers on/off (default: on)
-   **Context mode:** Press `c` to toggle between focus mode (changes only) and full context (entire file)

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