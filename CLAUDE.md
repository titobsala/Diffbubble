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
go build -o diffbubble main.go

# Run the built binary
./diffbubble
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
  - `DiffMode`: Enum for diff modes (DiffAll, DiffStaged, DiffUnstaged)
  - `FileStatus`: Enum for file states (Modified, Added, Deleted, Renamed)
  - `FileStat`: Metadata about changed files (path, status, additions, deletions)
- Functions:
  - `Diff()`: Returns full unified diff output (legacy, kept for compatibility)
  - `GetModifiedFiles(mode)`: Returns list of all modified files with stats based on diff mode
  - `GetFileDiff(filepath, contextLines, mode)`: Returns unified diff for a specific file with context and mode support

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
  - **Feature toggles**:
    - `showLineNumbers bool` for line number display (default: true)
    - `fullContext bool` for context mode (false = focus mode, true = full context)
  - **CLI options**:
    - `diffMode git.DiffMode` for staged/unstaged/all changes
    - `initialFile string` for pre-selecting a file on startup
- Update logic:
  - **Async loading**: `filesLoadedMsg` and `fileDiffLoadedMsg` for non-blocking git operations
  - **File navigation**: j/k keys navigate file list when focused, load new diff on selection
  - **Initial file selection**: If `--file` flag provided, searches for and selects that file on startup
  - **Line number toggle**: 'n' key toggles line numbers and re-renders diff
  - **Context toggle**: 'c' key toggles between focus mode and full context
  - **Focus switching**: Tab key switches between file list and diff panes
  - **Synchronized scrolling**: Diff panes scroll together via `YOffset` syncing
  - Window resize handling for all three panes
- View rendering:
  - **3-column layout**: File list sidebar | Left diff | Right diff
  - Sidebar shows: status icon (M/A/D/R), filename, +/- stats
  - Dynamic width calculations based on terminal size
  - Header, body, and footer sections
  - Error state handling
- CLI Flags (v0.2.0+):
  - `--help, -h`: Show help message
  - `--version, -v`: Show version information
  - `--file=<filename>`: Open with specific file pre-selected
  - `--staged`: Show only staged changes (git diff --cached)
  - `--unstaged`: Show only unstaged changes

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

## Release Automation

### GoReleaser Configuration (`.goreleaser.yml`)
The project uses GoReleaser for automated multi-platform builds:
- **Platforms**: Linux, macOS, Windows
- **Architectures**: amd64, ARM64
- **Artifacts**:
  - Compressed archives (.tar.gz for Unix, .zip for Windows)
  - SHA256 checksums
  - Automatic changelog generation from git commits
- **GitHub Release**: Automatically creates releases with binaries attached

### GitHub Actions Workflows

#### CI Workflow (`.github/workflows/ci.yml`)
Runs on every push and pull request:
- **Multi-platform testing**: Ubuntu, macOS, Windows
- **Go version**: 1.25
- **Checks**:
  - `go test -v -race -coverprofile=coverage.txt`
  - `go build -v ./...`
  - `go vet ./...`
  - `go fmt ./...` (enforces code formatting)

#### Release Workflow (`.github/workflows/release.yml`)
Triggers on version tags (v*.*.*):
1. Runs tests to ensure code quality
2. Builds binaries for all platforms using GoReleaser
3. Creates GitHub release with:
   - Changelog from commits
   - All platform binaries
   - Checksums file
   - Installation instructions

### Creating a Release
```bash
# 1. Ensure code is formatted
go fmt ./...

# 2. Commit all changes
git add .
git commit -m "feat: your changes"

# 3. Create and push tag
git tag -a v0.x.x -m "Release v0.x.x"
git push origin main
git push origin v0.x.x

# 4. GitHub Actions automatically builds and publishes
```

### Code Formatting
**Important**: All code must be formatted with `go fmt` before committing:
```bash
# Format all files
go fmt ./...

# The CI will fail if code is not formatted
```

## Features Roadmap

### ✅ Implemented Features

#### v0.1.0
- Side-by-side diff display with synchronized scrolling
- Multi-file navigation with sidebar
- Line number toggle ('n' key)
- Context mode toggle ('c' key)
- Mouse support
- Focus indicators
- Color-coded statistics

#### v0.2.0
- ✅ CLI flag: `--file=<filename>` - Pre-select specific file on startup
- ✅ CLI flag: `--staged` - Show only staged changes
- ✅ CLI flag: `--unstaged` - Show only unstaged changes
- ✅ Release automation with GoReleaser
- ✅ CI/CD with GitHub Actions
- ✅ Multi-platform binary builds

### 🎯 High Priority (v0.3.0)

#### 1. Theme System 🎨
**Priority**: HIGH | **Effort**: Medium | **User Value**: Very High

- Implement `--theme=<name>` CLI flag
- Built-in themes:
  - `dark` (current/default)
  - `light` - Light background with dark text
  - `high-contrast` - Maximum contrast for accessibility
  - `solarized` - Popular solarized color scheme
  - `dracula` - Dark theme with purple accents
  - `github` - GitHub-style diff colors
- Theme configuration structure in code
- Syntax highlighting integration for different file types
- Custom theme support via config file (future)

**Implementation Notes**:
- Add `ui/themes.go` with theme definitions
- Update `ui/styles.go` to use theme colors
- Add theme parameter to render functions
- Store current theme in model state

#### 2. Configuration File Support ⚙️
**Priority**: HIGH | **Effort**: Medium | **User Value**: High

- Config file location: `~/.config/diffbubble/config.yaml`
- Per-repository config: `.diffbubble.yml` in repo root
- Configurable settings:
  - Default theme
  - Line numbers (on/off by default)
  - Context mode (focus/full by default)
  - Diff mode (all/staged/unstaged)
  - Custom key bindings
  - Custom colors (advanced)
- Config priority: CLI flags > repo config > user config > defaults

**Example config.yaml**:
```yaml
theme: dracula
line_numbers: true
context_mode: focus
diff_mode: all
key_bindings:
  search: "/"
  next_file: "j"
  prev_file: "k"
```

#### 3. Search Functionality 🔍
**Priority**: HIGH | **Effort**: Medium | **User Value**: Very High

- Press `/` to enter search mode
- Real-time search highlighting
- Navigate results with `n` (next) and `N` (previous)
- Case-sensitive and case-insensitive modes
- Regex support (toggle with flag)
- Search scope: current file or all files
- Search status in footer (e.g., "Match 3 of 15")

**Implementation Notes**:
- Add search state to model (query, matches, current index)
- Update render to highlight search matches
- Add search input mode (like vim)
- Preserve scroll position when navigating matches

### 📝 Medium Priority (v0.4.0)

#### 4. Copy to Clipboard 📋
**Priority**: MEDIUM | **Effort**: Low-Medium | **User Value**: Medium

- `y` key to copy current line
- `Y` key to copy entire hunk
- Visual mode to select multiple lines (like vim)
- Copy to system clipboard (cross-platform)
- Show confirmation message after copy
- Support for different clipboard formats (plain text, markdown)

**Library**: Use `github.com/atotto/clipboard` for cross-platform support

#### 5. Export Functionality 💾
**Priority**: MEDIUM | **Effort**: Medium | **User Value**: Medium

- Export current diff to file
- Export formats:
  - Plain text (.diff, .patch)
  - HTML with syntax highlighting and styles
  - Markdown with code blocks
  - PDF (via HTML intermediate)
- CLI flag: `--export=<format>` for non-interactive export
- Interactive: Press `e` to open export menu
- Customizable export templates

#### 6. Advanced Git Features 🔧
**Priority**: MEDIUM | **Effort**: Medium-High | **User Value**: High

- `--ignore-whitespace` flag (git diff -w)
- `--ignore-all-space` flag (git diff -b)
- `--word-diff` mode - Show word-level changes
- Compare specific commits: `--compare=<commit1>..<commit2>`
- Compare branches: `--compare=main..feature`
- Interactive staging (press `s` on hunk to stage)
- Unstage hunks (press `u` on staged hunk)
- File history navigation (see previous versions)

**Implementation Notes**:
- Extend git.DiffMode to support more options
- Add interactive staging commands
- Integrate with git add -p style interaction

### 🌟 Nice to Have (v0.5.0+)

#### 7. Syntax Highlighting 🌈
**Priority**: LOW-MEDIUM | **Effort**: High | **User Value**: High

- Language-specific syntax highlighting
- Library options:
  - `github.com/alecthomas/chroma` (Pygments port)
  - `github.com/tree-sitter/go-tree-sitter` (more accurate)
- Detect language from file extension
- Highlight within diff context (color both diff and syntax)
- Theme-aware highlighting
- Fallback to current coloring if language unsupported

#### 8. Git Blame Integration 👤
**Priority**: LOW | **Effort**: Medium | **User Value**: Medium

- Toggle blame view with `b` key
- Show commit hash, author, and date for each line
- Inline blame info (compact mode)
- Full blame panel (detailed mode)
- Click on commit to see full commit message
- Navigate to commit in browser (if remote)

#### 9. Bookmarks & Favorites ⭐
**Priority**: LOW | **Effort**: Low | **User Value**: Low-Medium

- Press `m` to bookmark current file
- Press `'` to show bookmark list
- Jump to bookmarked files quickly
- Persistent bookmarks (saved per repository)
- Bookmark groups/categories
- Show bookmark indicator in file list

#### 10. Performance Optimizations ⚡
**Priority**: MEDIUM (for large diffs) | **Effort**: High | **User Value**: Medium

- Virtual scrolling for huge files (10,000+ lines)
- Lazy loading of diff content
- Incremental rendering
- Diff caching (remember parsed diffs)
- Async file list loading with progress indicator
- Background pre-loading of adjacent files
- Memory usage optimization for large repos

#### 11. Diff Algorithm Selection 🧮
**Priority**: LOW | **Effort**: Low-Medium | **User Value**: Low

- Support different diff algorithms:
  - Myers (default)
  - Patience (better for code)
  - Histogram (faster patience)
  - Minimal (smallest diff)
- CLI flag: `--diff-algorithm=<algorithm>`
- Config file setting
- Per-file-type defaults

#### 12. Remote & Branch Diff 🌐
**Priority**: MEDIUM | **Effort**: Medium | **User Value**: High

- Diff between local and remote branches
- CLI: `--remote-diff=origin/main`
- Pull request diff view
- Show commits in branch
- Fetch remote before diffing
- Compare any two refs (branches, tags, commits)

#### 13. Split View Modes 📐
**Priority**: LOW | **Effort**: Medium | **User Value**: Low-Medium

- Vertical split (current default)
- Horizontal split mode (top/bottom)
- Unified diff mode (single pane, like git diff)
- Toggle between modes with hotkey
- Remember preferred mode in config

#### 14. Diff Statistics & Analytics 📊
**Priority**: LOW | **Effort**: Low | **User Value**: Low

- Total lines changed across all files
- Complexity metrics
- Chart/graph of changes by file
- Language breakdown
- Churn analysis (files changed most often)
- Export statistics to JSON/CSV

### 🔮 Experimental Ideas (Research Needed)

#### 15. AI-Powered Features 🤖
- AI-generated commit message suggestions
- Automated code review comments
- Change summarization
- Detect potential bugs in diff
- Integration with GitHub Copilot/Claude

#### 16. Collaborative Features 👥
- Share diff view with others (via web link)
- Real-time collaborative review
- Comment threads on lines
- Review approval workflow
- Integration with code review tools

#### 17. IDE Integration 🔌
- VS Code extension
- Neovim plugin
- JetBrains plugin
- Integrate as external diff tool in gitconfig

### 📊 Priority Matrix

**Recommended Implementation Order** (considering value vs effort):

**v0.3.0** (Q1 2025):
1. Theme System
2. Config File Support
3. Search Functionality

**v0.4.0** (Q2 2025):
4. Copy to Clipboard
5. Advanced Git Features (--ignore-whitespace, etc.)
6. Export Functionality

**v0.5.0** (Q3 2025):
7. Syntax Highlighting
8. Performance Optimizations
9. Remote & Branch Diff

**Future** (as requested):
- Git Blame Integration
- Bookmarks
- Diff Algorithm Selection
- Split View Modes
- Statistics
- Experimental features

### 💡 Contributing Ideas

If you want to contribute a feature:
1. Check if it's already planned above
2. Open an issue to discuss the feature
3. Wait for feedback before implementing
4. Follow the architecture patterns in this document
5. Add tests for new functionality
6. Update this document with implementation details
