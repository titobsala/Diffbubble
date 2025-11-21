# Side-by-Side Git Diff TUI

A Terminal User Interface (TUI) application written in Go to display a side-by-side git diff with synchronized scrolling, multi-file navigation, and beautiful color-coded statistics.

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

## How to Run

1.  **Prerequisites:** Ensure you have [Go](https://go.dev/doc/install) and [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) installed on your system.
2.  **Navigate to Project:** Open your terminal in the project's root directory.
3.  **Install Dependencies:** If you haven't already, run the following command to download the necessary Go modules:
    ```sh
    go mod tidy
    ```
4.  **Run the Application:** Execute the following command to start the TUI. The application shows the output of `git diff`, so make sure you have staged or unstaged changes in your local repository to see a result.
    ```sh
    go run main.go
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