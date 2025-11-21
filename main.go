package main

import (
	"bytes"
	"fmt"

	"diffbuble/git"
	"diffbuble/parser"
	"diffbuble/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const appTitle = "Git Diff Side-by-Side"

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
	showLineNumbers bool
	fullContext     bool // false = focus mode (default), true = full context mode
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
	return loadFilesCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "n":
			// Toggle line numbers
			m.showLineNumbers = !m.showLineNumbers
			if len(m.currentRows) > 0 {
				m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers))
				m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers))
			}
			return m, nil

		case "c":
			// Toggle context mode (focus vs full context)
			m.fullContext = !m.fullContext
			// Reload current file's diff with new context
			if len(m.files) > 0 && m.selectedFile >= 0 && m.selectedFile < len(m.files) {
				return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext)
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
					return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext)
				}
				return m, nil
			}
			// Otherwise scroll diff

		case "k", "up":
			if m.focus == focusFileList && len(m.files) > 0 {
				// Navigate file list
				if m.selectedFile > 0 {
					m.selectedFile--
					return m, loadFileDiffCmd(m.files[m.selectedFile].Path, m.fullContext)
				}
				return m, nil
			}
			// Otherwise scroll diff
		}

	case filesLoadedMsg:
		m.files = msg.files
		m.err = msg.err

		if m.err == nil && len(m.files) > 0 {
			// Update file list viewport content
			if m.ready {
				m.fileListView.SetContent(ui.RenderFileList(m.files, m.selectedFile))
			}
			// Auto-load first file's diff
			m.selectedFile = 0
			return m, loadFileDiffCmd(m.files[0].Path, m.fullContext)
		}
		return m, nil

	case fileDiffLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.currentRows = msg.rows
			m.err = nil

			// Update diff viewports
			m.leftView.SetContent(ui.RenderSide(m.currentRows, ui.SideLeft, m.showLineNumbers))
			m.rightView.SetContent(ui.RenderSide(m.currentRows, ui.SideRight, m.showLineNumbers))

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
		headerHeight := 3  // Title + margin + buffer
		footerHeight := 3  // Footer can wrap to 2-3 lines in narrow terminals
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

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := ui.TitleStyle.Render(appTitle)
	focusOnFileList := m.focus == focusFileList
	footer := ui.RenderFooter(m.showLineNumbers, m.fullContext, focusOnFileList, m.winWidth)

	if m.err != nil {
		errorBox := ui.ErrorBox(m.err, m.winWidth)
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

	return lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
}

func loadFilesCmd() tea.Cmd {
	return func() tea.Msg {
		files, err := git.GetModifiedFiles()
		return filesLoadedMsg{files: files, err: err}
	}
}

func loadFileDiffCmd(filepath string, fullContext bool) tea.Cmd {
	return func() tea.Msg {
		contextLines := 0 // default
		if fullContext {
			contextLines = -1 // full context
		}

		diffOutput, err := git.GetFileDiff(filepath, contextLines)
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

func main() {
	p := tea.NewProgram(
		model{
			showLineNumbers: true, // Default on
			focus:           focusFileList,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
