# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Terminal User Interface (TUI) application written in Go that displays side-by-side git diffs with synchronized scrolling and multi-file navigation, built using the Bubble Tea framework (Elm Architecture for Go).

## Development Commands

### Build and Run
```bash
# Install/update dependencies
go mod tidy

# Run the application (displays git diff of current working directory)
go run main.go

# Build binary
go build -o diffbuble main.go

# Run the built binary
./diffbuble
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./parser
go test ./ui
go test ./git
```

## Architecture

The application follows a clean separation of concerns with three main layers and a 3-column UI layout:

### 1. Git Layer (`git/`)
- `git/diff.go`: Git command execution and file metadata retrieval
- Key types:
  - `FileStatus`: Enum for file states (Modified, Added, Deleted, Renamed)
  - `FileStat`: Metadata about changed files (path, status, additions, deletions)
- Functions:
  - `Diff()`: Returns full unified diff output (legacy, kept for compatibility)
  - `GetModifiedFiles()`: Returns list of all modified files with stats (uses `git diff --numstat` and `--name-status`)
  - `GetFileDiff(filepath)`: Returns unified diff for a specific file

### 2. Parser Layer (`parser/`)
- `parser/parser.go`: Parses unified diff format into structured data
- Key types:
  - `LineKind`: Enum for line types (context, addition, deletion, header)
  - `DiffLine`: Single line with number, content, and kind
  - `DiffRow`: Paired left/right lines for side-by-side alignment
- `Parse()` function handles:
  - Buffering pending additions/deletions for proper alignment
  - Line numbering for left (old file) and right (new file) sides
  - Flushing algorithm pairs deletions with additions

### 3. UI Layer (`ui/`)
- `ui/render.go`: Converts parsed diff rows into styled strings for viewports
  - `RenderSide(rows, side, showLineNumbers)`: Renders one side (left or right) with optional line numbers
  - `RenderFileList(files, selectedIdx)`: Renders sidebar file list with stats and status indicators
  - `RenderFooter(showLineNumbers)`: Renders footer with keyboard hints and line number state
  - Dynamic line number width calculation based on max line number
  - Header rendering with separators
- `ui/styles.go`: Centralized lipgloss styles for colors and formatting
  - Addition lines: green (#43BF6D)
  - Deletion lines: red (#E05252)
  - Modified files: yellow (#F5C842)
  - Headers: gray with separators
  - File list: sidebar styles with selection highlighting

### 4. Main Application (`main.go`)
- Implements Bubble Tea's Model-View-Update pattern with 3-column layout
- Model state:
  - **File list**: `[]git.FileStat` stores all modified files, `selectedFile` tracks current selection
  - **Three viewports**: `fileListView` (sidebar), `leftView` and `rightView` (diff panes)
  - **Current diff**: `currentRows []parser.DiffRow` stores only the selected file's parsed diff
  - **Focus management**: `focus focusPane` tracks whether file list or diff has focus
  - **Feature toggles**: `showLineNumbers bool` for line number display (default: true)
- Update logic:
  - **Async loading**: `filesLoadedMsg` and `fileDiffLoadedMsg` for non-blocking git operations
  - **File navigation**: j/k keys navigate file list when focused, load new diff on selection
  - **Line number toggle**: 'n' key toggles line numbers and re-renders diff
  - **Focus switching**: Tab key switches between file list and diff panes
  - **Synchronized scrolling**: Diff panes scroll together via `YOffset` syncing
  - Window resize handling for all three panes
- View rendering:
  - **3-column layout**: File list sidebar | Left diff | Right diff
  - Sidebar shows: status icon (M/A/D/R), filename, +/- stats
  - Dynamic width calculations based on terminal size
  - Header, body, and footer sections
  - Error state handling

## Key Implementation Details

### Multi-File Architecture
- **Memory efficient**: Only the file list metadata is loaded upfront; diff content is loaded on-demand per file
- **Per-file loading**: When user selects a file, `GetFileDiff()` fetches that file's diff, parser processes it, viewports are updated
- **No scroll memory**: Scroll position resets to top when switching files (intentional design choice)
- **Async operations**: Both file list loading and per-file diff loading use Bubble Tea commands for non-blocking UI

### Focus Management
The application has two focus modes controlled by Tab key:
- `focusFileList`: j/k keys navigate the file list, selecting different files
- `focusDiff`: j/k keys scroll the diff panes
Only the focused viewport receives update messages

### Synchronized Scrolling
Both diff viewports receive the same input events when diff pane is focused, but only the left viewport's scroll offset is used. The right viewport's `YOffset` is explicitly set to match the left viewport after each update (main.go:189).

### Line Number Toggle
- 'n' key toggles `showLineNumbers` boolean
- When toggled, diff content is re-rendered immediately without re-fetching from git
- `RenderSide()` function conditionally includes line numbers based on the flag
- Line number width is calculated dynamically based on the maximum line number in each side

### Diff Alignment Algorithm
The parser maintains two queues (`pendingMinus` and `pendingPlus`) to align consecutive deletion and addition blocks. The `flush()` function pairs lines by index, creating rows where:
- Deletions appear only on the left
- Additions appear only on the right
- Context lines appear on both sides

### Line Numbering
Line numbers increment separately for left and right sides:
- Left side: tracks line numbers in the old file (deletions and context)
- Right side: tracks line numbers in the new file (additions and context)
- Headers don't have line numbers
- Can be toggled on/off with 'n' key

## Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework (Elm Architecture)
- `github.com/charmbracelet/bubbles`: Reusable TUI components (viewport)
- `github.com/charmbracelet/lipgloss`: Terminal styling and layout

## Keyboard Controls

### Global
- `q`, `esc`, `ctrl+c`: Quit application
- `tab`: Switch focus between file list and diff panes
- `n`: Toggle line numbers on/off

### When File List is Focused
- `j`, `↓`: Select next file (loads its diff)
- `k`, `↑`: Select previous file (loads its diff)

### When Diff Panes are Focused
- `j`, `k`, `↓`, `↑`: Scroll diff content
- Mouse wheel: Scroll diff content

## Notes

- The application requires a git repository with changes to display
- If no files are modified, the sidebar will show "No modified files"
- Mouse support is enabled via `tea.WithMouseCellMotion()`
- Full screen mode is enabled via `tea.WithAltScreen()`
- Diff content is fetched on-demand per file for memory efficiency

## Future Enhancements

The following features are planned for future implementation:

### CLI Flags
- `--file=<filename>`: Open application with a specific file selected
- File selection would be maintained when the flag is provided
- If file doesn't exist in modified files, default to first file

### Theme System
- Customizable color schemes for syntax highlighting
- User-configurable themes via config file
- Additional built-in themes (dark, light, high contrast)
- Support for custom color palettes for different file types

### Enhanced Styling Options
- Configurable diff colors (additions, deletions, context)
- Font styling options (bold, italic, underline)
- Border style customization
- Status indicator styling options

### Additional Features Under Consideration
- Search functionality within diff content
- Jump to specific line number
- Expand/collapse hunks
- Copy diff content to clipboard
- Export diff to file
- Support for staged vs unstaged diffs (toggle between them)
