package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/titobsala/Diffbubble/config"
	"github.com/titobsala/Diffbubble/git"
	"github.com/titobsala/Diffbubble/parser"
	"github.com/titobsala/Diffbubble/search"
	"github.com/titobsala/Diffbubble/ui"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	appTitle = "Git Diff Side-by-Side"
	version  = "0.3.1"
)

type focusPane int

const (
	focusFileList focusPane = iota
	focusDiff
)

type model struct {
	ready     bool
	winWidth  int
	winHeight int
	err       error

	// File list (sidebar)
	files        []git.FileStat
	selectedFile int
	fileListView viewport.Model
	focus        focusPane

	// Diff views (current file)
	currentRows []parser.DiffRow
	leftView    viewport.Model
	rightView   viewport.Model

	// Feature toggles
	showLineNumbers  bool
	fullContext      bool         // false = focus mode (default), true = full context mode
	diffMode         git.DiffMode // Which changes to show (all, staged, unstaged)
	initialFile      string       // File to pre-select on startup (if specified)
	currentThemeIdx  int          // Current theme index for 't' key cycling
	themeChangeMsg   string       // Brief message shown when theme changes
	themeChangeTicks int          // Counter to clear theme change message

	// Search state
	searchMode       bool            // Whether search mode is active
	searchInput      textinput.Model // Text input for search query
	searchMatches    []search.Match  // All matches found
	currentMatchIdx  int             // Index of current match being viewed (-1 if none)
	searchInAllFiles bool            // Whether to search across all files
}

// Message types for async operations
type filesLoadedMsg struct {
	files []git.FileStat
	err   error
}

type fileDiffLoadedMsg struct {
	rows []parser.DiffRow
	err  error
}

func (m model) Init() tea.Cmd {
	return loadFilesCmd(m.diffMode, m.initialFile)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()

		// Handle search mode input
		if m.searchMode {
			switch k {
			case "esc":
				// Exit search mode
				m.searchMode = false
				m.searchInput.Reset()
				return m, nil

			case "enter":
				m.searchMode = false
				return m, nil

			default:
				// Pass input to text input
				var newCmd tea.Cmd
				m.searchInput, newCmd = m.searchInput.Update(msg)

				// Perform dynamic search as user types
				m.performSearch()

				return m, newCmd
			}
		}

		// Normal mode key handling
		switch k {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Exit search mode (already handled above but just in case)
			if m.searchMode {
				m.searchMode = false
				m.searchInput.Reset()
				return m, nil
			}

			// Clear search matches if any exist
			if len(m.searchMatches) > 0 || m.searchInput.Value() != "" {
				m.searchMatches = nil
				m.currentMatchIdx = -1
				m.searchInput.Reset()

				// Refresh viewports to remove highlights
				if len(m.currentRows) > 0 {
					m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers))
					m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers))
				}
				return m, nil
			}

			// Otherwise quit
			return m, tea.Quit

		case "/":
			// Enter search mode
			m.searchMode = true
			m.searchInput.Focus()
			m.searchInput.Reset()
			return m, nil

		case "n":
			// Navigate to next match (if matches exist), otherwise toggle line numbers
			if len(m.searchMatches) > 0 && m.currentMatchIdx >= 0 {
				// Next match
				m.currentMatchIdx = (m.currentMatchIdx + 1) % len(m.searchMatches)
				match := m.searchMatches[m.currentMatchIdx]
				pos := search.GetMatchPosition(match)
				m.leftView.YOffset = pos
				m.rightView.YOffset = pos
				return m, nil
			}
			// Toggle line numbers
			m.showLineNumbers = !m.showLineNumbers
			if len(m.currentRows) > 0 {
				searchHighlights := convertSearchMatches(m.searchMatches, m.currentMatchIdx)
				m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers, searchHighlights...))
				m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers, searchHighlights...))
			}
			return m, nil

		case "N":
			// Navigate to previous match
			if len(m.searchMatches) > 0 && m.currentMatchIdx >= 0 {
				m.currentMatchIdx--
				if m.currentMatchIdx < 0 {
					m.currentMatchIdx = len(m.searchMatches) - 1
				}
				match := m.searchMatches[m.currentMatchIdx]
				pos := search.GetMatchPosition(match)
				m.leftView.YOffset = pos
				m.rightView.YOffset = pos
				return m, nil
			}
			return m, nil

		case "c":
			// Toggle context mode (focus vs full context)
			m.fullContext = !m.fullContext
			// Reload current file's diff with new context
			if len(m.files) > 0 && m.selectedFile >= 0 && m.selectedFile < len(m.files) {
				return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext, m.diffMode)
			}
			return m, nil

		case "t":
			// Cycle through themes
			themes := ui.ListThemes()
			m.currentThemeIdx = (m.currentThemeIdx + 1) % len(themes)
			newTheme := themes[m.currentThemeIdx]
			ui.SetTheme(newTheme)
			updateSearchStyles(&m.searchInput)

			// Show theme change message
			m.themeChangeMsg = fmt.Sprintf("Theme: %s", newTheme)
			m.themeChangeTicks = 3 // Show for 3 ticks

			// Re-render current diff with new theme
			if len(m.currentRows) > 0 {
				searchHighlights := convertSearchMatches(m.searchMatches, m.currentMatchIdx)
				m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers, searchHighlights...))
				m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers, searchHighlights...))
			}
			if len(m.files) > 0 && m.ready {
				m.fileListView.SetContent(ui.RenderFileList(m.files, m.selectedFile))
			}
			return m, nil

		case "tab":
			// Switch focus between file list and diff
			if m.focus == focusFileList {
				m.focus = focusDiff
			} else {
				m.focus = focusFileList
			}
			return m, nil

		case "j", "down":
			if m.focus == focusFileList && len(m.files) > 0 {
				// Navigate file list
				if m.selectedFile < len(m.files)-1 {
					m.selectedFile++
					return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext, m.diffMode)
				}
				return m, nil
			}
			// Otherwise scroll diff

		case "k", "up":
			if m.focus == focusFileList && len(m.files) > 0 {
				// Navigate file list
				if m.selectedFile > 0 {
					m.selectedFile--
					return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext, m.diffMode)
				}
				return m, nil
			}
			// Otherwise scroll diff
		}

	case filesLoadedMsg:
		m.files = msg.files
		m.err = msg.err

		if m.err == nil && len(m.files) == 0 {
			// No files found - provide helpful context-specific message
			switch m.diffMode {
			case git.DiffStaged:
				m.err = fmt.Errorf("No staged changes found.\n\nTry one of the following:\n  • Run 'git add <file>' to stage some changes\n  • Use --unstaged to see unstaged changes\n  • Remove --staged flag to see all changes")
			case git.DiffUnstaged:
				m.err = fmt.Errorf("No unstaged changes found.\n\nTry one of the following:\n  • Use --staged to see staged changes\n  • Remove --unstaged flag to see all changes\n  • Make some changes to your working directory")
			default:
				m.err = fmt.Errorf("No changes found in the repository.\n\nMake sure you have:\n  • Modified some files in your working directory\n  • Staged some changes with 'git add'\n  • Checked that you're in a git repository")
			}
		}

		if m.err == nil && len(m.files) > 0 {
			// Update file list viewport content
			if m.ready {
				m.fileListView.SetContent(ui.RenderFileList(m.files, m.selectedFile))
			}

			// Select initial file (either specified via --file flag or default to first)
			m.selectedFile = 0
			if m.initialFile != "" {
				// Find the specified file in the list
				for i, file := range m.files {
					if file.Path == m.initialFile {
						m.selectedFile = i
						break
					}
				}
			}

			return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext, m.diffMode)
		}
		return m, nil

	case fileDiffLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.currentRows = msg.rows
			m.err = nil

			// Update diff viewports
			searchHighlights := convertSearchMatches(m.searchMatches, m.currentMatchIdx)
			m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers, searchHighlights...))
			m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers, searchHighlights...))

			// Update file list to show new selection
			if len(m.files) > 0 {
				m.fileListView.SetContent(ui.RenderFileList(m.files, m.selectedFile))
			}

			// Reset scroll position
			m.leftView.YOffset = 0
			m.rightView.YOffset = 0
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height

		// Calculate dimensions: 20-40-40 split
		// Increased margin to account for header, footer, borders, and potential text wrapping
		headerHeight := 3 // Title + margin + buffer
		footerHeight := 3 // Footer can wrap to 2-3 lines in narrow terminals
		verticalMarginHeight := headerHeight + footerHeight

		// 20% for sidebar, 40% for each diff pane
		sidebarWidth := msg.Width * 20 / 100
		diffPaneWidth := msg.Width * 40 / 100

		// Account for borders (subtract a bit for padding)
		if sidebarWidth > 4 {
			sidebarWidth -= 4
		}
		if diffPaneWidth > 2 {
			diffPaneWidth -= 2
		}

		if !m.ready {
			m.ready = true

			// Initialize three viewports
			m.fileListView = viewport.New(sidebarWidth, msg.Height-verticalMarginHeight)
			m.leftView = viewport.New(diffPaneWidth, msg.Height-verticalMarginHeight)
			m.rightView = viewport.New(diffPaneWidth, msg.Height-verticalMarginHeight)
		} else {
			// Handle resize
			m.fileListView.Width = sidebarWidth
			m.fileListView.Height = msg.Height - verticalMarginHeight
			m.leftView.Width = diffPaneWidth
			m.leftView.Height = msg.Height - verticalMarginHeight
			m.rightView.Width = diffPaneWidth
			m.rightView.Height = msg.Height - verticalMarginHeight
		}

		// Update file list content
		if len(m.files) > 0 {
			m.fileListView.SetContent(ui.RenderFileList(m.files, m.selectedFile))
		}
	}

	// Update viewports based on focus
	if m.focus == focusFileList {
		m.fileListView, cmd = m.fileListView.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Sync scrolling for diff panes
		m.leftView, cmd = m.leftView.Update(msg)
		cmds = append(cmds, cmd)

		m.rightView, _ = m.rightView.Update(msg)
		m.rightView.YOffset = m.leftView.YOffset
	}

	// Decrement theme change message counter
	if m.themeChangeTicks > 0 {
		m.themeChangeTicks--
		if m.themeChangeTicks == 0 {
			m.themeChangeMsg = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := ui.TitleStyle.Render(appTitle)

	// Add theme change notification if active
	if m.themeChangeMsg != "" {
		themeMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9E2AF")).
			Bold(true).
			Render(" " + m.themeChangeMsg)
		header = lipgloss.JoinHorizontal(lipgloss.Center, header, themeMsg)
	}

	focusOnFileList := m.focus == focusFileList

	// Prepare search info for footer
	searchInfo := ""
	if len(m.searchMatches) > 0 && m.currentMatchIdx >= 0 {
		searchInfo = fmt.Sprintf("Match %d of %d", m.currentMatchIdx+1, len(m.searchMatches))
	} else if len(m.searchMatches) == 0 && m.searchInput.Value() != "" && !m.searchMode {
		searchInfo = "No matches found"
	}

	footer := ui.RenderFooter(m.showLineNumbers, m.fullContext, focusOnFileList, m.searchMode, searchInfo, m.winWidth)

	// Show search input if in search mode
	var searchBar string
	if m.searchMode {
		searchBar = ui.SearchInputStyle.Render(m.searchInput.View())
	}

	if m.err != nil {
		errorBox := ui.ErrorBox(m.err, m.winWidth)
		if searchBar != "" {
			return lipgloss.JoinVertical(lipgloss.Top, header, errorBox, searchBar, footer)
		}
		return lipgloss.JoinVertical(lipgloss.Top, header, errorBox, footer)
	}

	// Render file list sidebar with focus-aware styling
	fileListContent := m.fileListView.View()
	var sidebarBox string
	if focusOnFileList {
		sidebarBox = ui.FileListStyleFocused.Width(m.fileListView.Width).Height(m.fileListView.Height).Render(fileListContent)
	} else {
		sidebarBox = ui.FileListStyle.Width(m.fileListView.Width).Height(m.fileListView.Height).Render(fileListContent)
	}

	// Render diff panes with focus-aware styling
	var leftBox, rightBox string
	if focusOnFileList {
		// Diff panes are unfocused
		leftBox = ui.BorderStyleUnfocused.Width(m.leftView.Width).Render(m.leftView.View())
		rightBox = ui.BorderStyleUnfocused.Width(m.rightView.Width).Render(m.rightView.View())
	} else {
		// Diff panes are focused
		leftBox = ui.BorderStyleFocused.Width(m.leftView.Width).Render(m.leftView.View())
		rightBox = ui.BorderStyleFocused.Width(m.rightView.Width).Render(m.rightView.View())
	}

	// Join horizontally: sidebar | left diff | right diff
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, leftBox, rightBox)

	if searchBar != "" {
		return lipgloss.JoinVertical(lipgloss.Top, header, body, searchBar, footer)
	}
	return lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
}

// convertSearchMatches converts search.Match to ui.SearchMatch format
func convertSearchMatches(matches []search.Match, currentMatchIdx int) []ui.SearchMatch {
	var result []ui.SearchMatch
	for i, match := range matches {
		result = append(result, ui.SearchMatch{
			RowIndex:  match.RowIndex,
			Side:      match.Side,
			Column:    match.Column,
			Length:    match.Length,
			IsCurrent: i == currentMatchIdx,
		})
	}
	return result
}

func loadFilesCmd(mode git.DiffMode, initialFile string) tea.Cmd {
	return func() tea.Msg {
		files, err := git.GetModifiedFiles(mode)
		return filesLoadedMsg{files: files, err: err}
	}
}

func loadFileDiffCmd(filepath string, fullContext bool, mode git.DiffMode) tea.Cmd {
	return func() tea.Msg {
		contextLines := 0 // default
		if fullContext {
			contextLines = -1 // full context
		}

		diffOutput, err := git.GetFileDiff(filepath, contextLines, mode)
		if err != nil {
			return fileDiffLoadedMsg{err: err}
		}

		rows, parseErr := parser.Parse(bytes.NewReader(diffOutput))
		if parseErr != nil {
			return fileDiffLoadedMsg{err: parseErr}
		}

		return fileDiffLoadedMsg{rows: rows}
	}
}

func printVersion() {
	fmt.Printf("diffbubble version %s\n", version)
}

func printHelp() {
	fmt.Println("diffbubble - A Terminal UI for side-by-side git diffs")
	fmt.Printf("\nVersion: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  diffbubble [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -h, --help                    Show this help message")
	fmt.Println("  -v, --version                 Show version information")
	fmt.Println("  --file=<filename>             Open with specific file selected")
	fmt.Println("  --staged                      Show only staged changes (git diff --cached)")
	fmt.Println("  --unstaged                    Show only unstaged changes")
	fmt.Println("  --theme=<name>                Color theme (default: dark)")
	fmt.Println("  --list-themes                 List all available themes")
	fmt.Println("  --show-theme-colors <name>    Preview colors for a specific theme")
	fmt.Println("\nAvailable Themes:")
	fmt.Println("  dark, light, high-contrast, solarized, dracula, github,")
	fmt.Println("  catppuccin, tokyo-night, one-dark")
	fmt.Println("\nDescription:")
	fmt.Println("  diffbubble displays git diffs in a beautiful side-by-side format with")
	fmt.Println("  multi-file navigation, synchronized scrolling, and customizable themes.")
	fmt.Println("\nExamples:")
	fmt.Println("  diffbubble                               # Show all changes")
	fmt.Println("  diffbubble --staged                      # Show only staged changes")
	fmt.Println("  diffbubble --unstaged                    # Show only unstaged changes")
	fmt.Println("  diffbubble --file=README.md              # Open with README.md selected")
	fmt.Println("  diffbubble --theme=catppuccin            # Use Catppuccin theme")
	fmt.Println("  diffbubble --theme=tokyo-night --staged  # Tokyo Night theme, staged only")
	fmt.Println("  diffbubble --list-themes                 # List all available themes")
	fmt.Println("  diffbubble --show-theme-colors dracula   # Preview Dracula theme colors")
	fmt.Println("\nKeyboard Controls:")
	fmt.Println("  tab          Switch focus between file list and diff panes")
	fmt.Println("  j/k, ↓/↑     Navigate files (when file list focused) or scroll diff")
	fmt.Println("  n            Toggle line numbers on/off")
	fmt.Println("  c            Toggle between focus mode and full context")
	fmt.Println("  t            Cycle through themes interactively")
	fmt.Println("  q, esc       Quit the application")
	fmt.Println("\nRequires:")
	fmt.Println("  - A git repository with changes to display")
	fmt.Println("  - Git must be installed and available in PATH")
}

func printThemeList() {
	fmt.Println("Available themes:")
	fmt.Println()
	themes := ui.ListThemes()
	for _, theme := range themes {
		fmt.Printf("  • %s\n", theme)
	}
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  diffbubble --theme=<name>")
	fmt.Println("  diffbubble --show-theme-colors <name>  # Preview theme colors")
}

func printThemeColors(themeName string) {
	if !ui.ValidateTheme(themeName) {
		fmt.Printf("Error: Unknown theme '%s'\n", themeName)
		fmt.Printf("Available themes: %v\n", ui.ListThemes())
		os.Exit(1)
	}

	ui.SetTheme(themeName)
	theme := ui.GetTheme()

	fmt.Printf("Theme: %s\n", theme.Name)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()

	// Helper function to display color with ANSI codes
	colorBox := func(color, name string) {
		// Convert hex to ANSI 24-bit color
		r, g, b := hexToRGB(color)
		fmt.Printf("  %-20s ", name+":")
		fmt.Printf("\033[48;2;%d;%d;%dm   \033[0m", r, g, b) // Background box
		fmt.Printf("  %s\n", color)
	}

	fmt.Println("Diff Colors:")
	colorBox(theme.AdditionFg, "  Addition (text)")
	colorBox(theme.AdditionBg, "  Addition (bg)")
	colorBox(theme.DeletionFg, "  Deletion (text)")
	colorBox(theme.DeletionBg, "  Deletion (bg)")
	colorBox(theme.ContextFg, "  Context")
	colorBox(theme.HeaderFg, "  Headers")
	fmt.Println()

	fmt.Println("UI Colors:")
	colorBox(theme.FocusedBorderColor, "  Focused border")
	colorBox(theme.BorderColor, "  Border")
	colorBox(theme.TitleFg, "  Title")
	fmt.Println()

	fmt.Println("File List Colors:")
	colorBox(theme.ModifiedFg, "  Modified")
	colorBox(theme.AddedFg, "  Added")
	colorBox(theme.DeletedFg, "  Deleted")
	fmt.Println()

	fmt.Println("General:")
	colorBox(theme.Background, "  Background")
	colorBox(theme.Foreground, "  Foreground")
	fmt.Println()
}

func hexToRGB(hex string) (int, int, int) {
	// Remove # if present
	hex = strings.TrimPrefix(hex, "#")

	// Parse hex color
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func (m *model) performSearch() {
	query := m.searchInput.Value()
	if query == "" {
		m.searchMatches = nil
		m.currentMatchIdx = -1
	} else {
		// Search current file
		if len(m.currentRows) > 0 {
			fileName := ""
			if len(m.files) > 0 && m.selectedFile >= 0 && m.selectedFile < len(m.files) {
				fileName = m.files[m.selectedFile].Path
			}
			// Use case-insensitive search for now (false)
			m.searchMatches = search.SearchInRows(m.currentRows, query, fileName, false)
			if len(m.searchMatches) > 0 {
				m.currentMatchIdx = 0
				// Scroll to first match
				match := m.searchMatches[0]
				pos := search.GetMatchPosition(match)
				if match.Side == "left" {
					m.leftView.YOffset = pos
					m.rightView.YOffset = pos
				} else {
					m.rightView.YOffset = pos
					m.leftView.YOffset = pos
				}
			} else {
				m.currentMatchIdx = -1
			}
		}
	}

	// Refresh viewports to show/hide highlights
	if len(m.currentRows) > 0 {
		searchHighlights := convertSearchMatches(m.searchMatches, m.currentMatchIdx)
		m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers, searchHighlights...))
		m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers, searchHighlights...))
	}
}

func updateSearchStyles(ti *textinput.Model) {
	theme := ui.GetTheme()
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Foreground))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ContextFg))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FocusedBorderColor))
}

func main() {
	// Load configuration file (user + repo)
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		fmt.Println("Using default configuration...")
		cfg = &config.Config{} // Use empty config
		*cfg = config.DefaultConfig()
	}
	cfg.Validate()

	var (
		showVersion     bool
		showHelp        bool
		selectedFile    string
		showStaged      bool
		showUnstaged    bool
		themeName       string
		listThemes      bool
		showThemeColors string
	)

	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.StringVar(&selectedFile, "file", "", "Open with specific file selected")
	flag.BoolVar(&showStaged, "staged", false, "Show only staged changes")
	flag.BoolVar(&showUnstaged, "unstaged", false, "Show only unstaged changes")
	flag.StringVar(&themeName, "theme", cfg.Theme, "Color theme")
	flag.BoolVar(&listThemes, "list-themes", false, "List all available themes")
	flag.StringVar(&showThemeColors, "show-theme-colors", "", "Show color preview for a theme")
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if listThemes {
		printThemeList()
		os.Exit(0)
	}

	if showThemeColors != "" {
		printThemeColors(showThemeColors)
		os.Exit(0)
	}

	// Validate and set theme (CLI flag or config file)
	if !ui.ValidateTheme(themeName) {
		fmt.Printf("Error: Invalid theme '%s'. Available themes: %v\n", themeName, ui.ListThemes())
		os.Exit(1)
	}
	ui.SetTheme(themeName)

	// Determine diff mode based on flags (CLI flags override config)
	diffMode := git.DiffAll
	if showStaged && showUnstaged {
		fmt.Println("Error: Cannot use both --staged and --unstaged flags together")
		os.Exit(1)
	} else if showStaged {
		diffMode = git.DiffStaged
	} else if showUnstaged {
		diffMode = git.DiffUnstaged
	} else {
		// Use config diff mode if no CLI flag provided
		switch cfg.DiffMode {
		case "staged":
			diffMode = git.DiffStaged
		case "unstaged":
			diffMode = git.DiffUnstaged
		default:
			diffMode = git.DiffAll
		}
	}

	// Determine initial context mode and line numbers from config
	fullContext := cfg.ContextMode == "full"

	// Find initial theme index for 't' key cycling
	themeIdx := 0
	themes := ui.ListThemes()
	for i, t := range themes {
		if t == themeName {
			themeIdx = i
			break
		}
	}

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100
	ti.Width = 50
	updateSearchStyles(&ti)

	p := tea.NewProgram(
		model{
			showLineNumbers:  cfg.LineNumbers, // From config
			fullContext:      fullContext,     // From config
			focus:            focusFileList,
			diffMode:         diffMode,
			initialFile:      selectedFile,
			currentThemeIdx:  themeIdx,
			searchInput:      ti,
			currentMatchIdx:  -1,   // No match selected initially
			searchInAllFiles: true, // Default to searching all files
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
