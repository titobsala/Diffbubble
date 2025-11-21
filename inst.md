Technical Specification: Side-by-Side Git Diff TUI

1. Project Overview

The goal is to build a Terminal User Interface (TUI) application in Go that displays a git diff in a side-by-side split view.

Left Pane: Shows the "Old" (Base) version of the file.

Right Pane: Shows the "New" (Head) version of the file.

Behavior: Scrolling must be synchronized (locking the Y-offset of both panes).

2. Technology Stack

Language: Go (Golang)

Framework: Bubble Tea (The Elm Architecture for TUI)

Styling: Lip Gloss (Layouts and Colors)

Components: Bubbles (Specifically the viewport component)

3. Core Logic: The "Flush" Algorithm

The most complex part of this application is converting a standard Unified Diff (which comes in a single stream of lines starting with + or -) into two distinct, aligned text blocks.

3.1 The Parsing Rules

We cannot simply stream lines to two viewports. We must parse them into "hunks" to ensure visual alignment.

Iterate through the raw diff line by line.

Maintain two buffers (temporary lists): pendingMinus (lines removed) and pendingPlus (lines added).

Trigger a "Flush": When we encounter a "Context" line (a line starting with a space) or a new Hunk Header (@@ ... @@), we must write the contents of our temporary buffers to the final strings.

Alignment Logic (Crucial): During a flush, if pendingMinus has 2 lines and pendingPlus has 5 lines, we must add 3 empty newlines to the "Left" side so the subsequent context lines start at the exact same vertical position in both viewports.

3.2 Pseudo-Code for Alignment

maxLen = max(len(pendingMinus), len(pendingPlus))

for i from 0 to maxLen:
    // Left Side
    if i < len(pendingMinus):
        finalLeftString += pendingMinus[i]
    else:
        finalLeftString += "\n" // Padding for alignment

    // Right Side
    if i < len(pendingPlus):
        finalRightString += pendingPlus[i]
    else:
        finalRightString += "\n" // Padding for alignment


4. Architecture & State

The Bubble Tea model struct requires the following state:

type model struct {
    ready     bool            // Has the terminal size been initialized?
    leftView  viewport.Model  // The bubbles component for the left pane
    rightView viewport.Model  // The bubbles component for the right pane
    diffData  DiffContent     // Holds the parsed Left and Right strings
    winWidth  int             // Current terminal width
    winHeight int             // Current terminal height
}


5. Reference Implementation

Below is a working prototype using mock data. Use this as the ground truth for how to implement the synchronization and styling.

package main

import (
	"fmt"
	"strings"

	"[github.com/charmbracelet/bubbles/viewport](https://github.com/charmbracelet/bubbles/viewport)"
	tea "[github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)"
	"[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)"
)

// ---------------------------------------------------------
// 1. CONSTANTS & MOCK DATA
// ---------------------------------------------------------

// This mocks the output of 'git diff'. In a real app, you would use:
// out, _ := exec.Command("git", "diff").Output()
const mockGitDiff = `diff --git a/main.go b/main.go
index 88a2s..d1234 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,5 @@
 package main
 
-import "fmt"
+import "log"
 
 func main() {
@@ -10,4 +10,8 @@
-    fmt.Println("Hello World")
-    fmt.Println("This line is removed")
+    log.Println("Hello Charm")
+    log.Println("This is a new feature")
+    log.Println("Added extra logging")
+    log.Println("More lines added than removed")
 
     // Context line that should align
     return
}`

// ---------------------------------------------------------
// 2. DIFF PARSING LOGIC (The "Hard" Part)
// ---------------------------------------------------------

// DiffContent holds the aligned strings for left (old) and right (new) views
type DiffContent struct {
	Left  string
	Right string
}

// parseDiff converts a Unified Diff string into two aligned side-by-side strings.
// This is a simplified algorithm. Robust diffing algorithms (like Myers) are complex,
// but this "flush on context" approach works for 90% of visual use cases.
func parseDiff(diff string) DiffContent {
	var leftBuf, rightBuf strings.Builder
	lines := strings.Split(diff, "\n")

	// Temporary buffers for the current "hunk" of changes
	var pendingMinus []string
	var pendingPlus []string

	// Helper to flush pending changes to the main buffers with alignment
	flush := func() {
		maxLen := len(pendingMinus)
		if len(pendingPlus) > maxLen {
			maxLen = len(pendingPlus)
		}

		for i := 0; i < maxLen; i++ {
			// Left Side (Old)
			if i < len(pendingMinus) {
				leftBuf.WriteString(pendingMinus[i] + "\n")
			} else {
				leftBuf.WriteString("\n") // Padding
			}

			// Right Side (New)
			if i < len(pendingPlus) {
				rightBuf.WriteString(pendingPlus[i] + "\n")
			} else {
				rightBuf.WriteString("\n") // Padding
			}
		}
		// Reset buckets
	
pendingMinus = nil
pendingPlus = nil
	}

	for _, line := range lines {
		// Ignore git metadata headers for the view content
		if strings.HasPrefix(line, "diff") || 
		   strings.HasPrefix(line, "index") || 
		   strings.HasPrefix(line, "---") || 
		   strings.HasPrefix(line, "+++") {
			continue
		}

		// Hunk Header (@@ -1,5 +1,5 @@)
		if strings.HasPrefix(line, "@@") {
			flush()
			// You could render the hunk header here if you wanted
			separator := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", 30))
			leftBuf.WriteString(fmt.Sprintf("%s\n%s\n", separator, line))
			rightBuf.WriteString(fmt.Sprintf("%s\n%s\n", separator, line))
			continue
		}

		// Logic based on first character
		if len(line) > 0 {
			switch line[0] {
			case '-':
				// Removal
			
pendingMinus = append(pendingMinus, line)
			case '+':
				// Addition
			
pendingPlus = append(pendingPlus, line)
			case ' ':
				// Context: Context acts as a barrier that forces a flush of previous changes
				flush()
				leftBuf.WriteString(line + "\n")
				rightBuf.WriteString(line + "\n")
			}
		} else {
			// Empty lines in diff are usually context
			flush()
			leftBuf.WriteString("\n")
			rightBuf.WriteString("\n")
		}
	}
	flush() // Flush any remaining changes at EOF

	return DiffContent{
		Left:  leftBuf.String(),
		Right: rightBuf.String(),
	}
}

// ---------------------------------------------------------
// 3. BUBBLE TEA MODEL & STYLING
// ---------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")),
			Padding(0, 1).
			MarginBottom(1)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	// Style for added lines (greenish)
	addStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D"))
	// Style for removed lines (redish)
	delStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E05252"))
)

type model struct {
	ready     bool
	leftView  viewport.Model
	rightView viewport.Model
	diffData  DiffContent
	winWidth  int
	winHeight int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height
		
		// Calculate dimensions
		headerHeight := 2
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.ready = true
			// Parse the diff once on startup
			m.diffData = parseDiff(mockGitDiff)

			// Initialize viewports
			m.leftView = viewport.New(msg.Width/2-2, msg.Height-verticalMarginHeight)
			m.rightView = viewport.New(msg.Width/2-2, msg.Height-verticalMarginHeight)
			
			// Set content with simple syntax highlighting
			m.leftView.SetContent(highlight(m.diffData.Left, "left"))
			m.rightView.SetContent(highlight(m.diffData.Right, "right"))
		} else {
			// Handle resize
			m.leftView.Width = msg.Width/2 - 2
			m.leftView.Height = msg.Height - verticalMarginHeight
			m.rightView.Width = msg.Width/2 - 2
			m.rightView.Height = msg.Height - verticalMarginHeight
		}
	}

	// Sync Scrolling: Update both viewports with the same message
	// This ensures if you press "down" on one, both move.
	m.leftView, cmd = m.leftView.Update(msg)
	cmds = append(cmds, cmd)

	m.rightView, _ = m.rightView.Update(msg)
	// We don't append the right view's command to avoid duplicate key handling 
	// artifacts if they both processed the same key, though for viewports it's usually fine.
	// To keep them perfectly synced, we force the Y offset to match.
	m.rightView.YOffset = m.leftView.YOffset

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Header
	header := titleStyle.Render("Git Diff Side-by-Side")

	// Apply borders
	leftBox := borderStyle.Width(m.winWidth/2 - 3).Render(m.leftView.View())
	rightBox := borderStyle.Width(m.winWidth/2 - 3).Render(m.rightView.View())

	// Join them horizontally
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	
	// Footer
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Scroll with j/k/arrows • q to quit")

	return lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
}

// highlight applies simple color based on the first character of the line
func highlight(content, side string) string {
	lines := strings.Split(content, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if len(line) == 0 {
			sb.WriteString("\n")
			continue
		}
		
		// Check the first char to determine color
		// Note: The parser preserved the +/- signs so we can check them here
		firstChar := line[0]
		
		if firstChar == '-' && side == "left" {
			sb.WriteString(delStyle.Render(line) + "\n")
		} else if firstChar == '+' && side == "right" {
			sb.WriteString(addStyle.Render(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}

func main() {
	p := tea.NewProgram(
		model{},
		tea.WithAltScreen(),       // Use full screen
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

This is a new line.