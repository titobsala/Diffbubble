package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// getGitDiff executes the 'git diff' command and returns its output as a string.
// If an error occurs, it returns a formatted error string.
func getGitDiff() string {
	cmd := exec.Command("git", "diff")
	out, err := cmd.Output()
	if err != nil {
		// Return a user-friendly error message to be displayed in the TUI
		return fmt.Sprintf("Error running 'git diff':\n\n%s\n\nPlease ensure 'git' is installed and you are in a git repository.", err)
	}
	return string(out)
}

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

	// Line number counters for each side
	leftLineNum := 1
	rightLineNum := 1

	// Helper to flush pending changes to the main buffers with alignment
	flush := func() {
		maxLen := len(pendingMinus)
		if len(pendingPlus) > maxLen {
			maxLen = len(pendingPlus)
		}

		for i := 0; i < maxLen; i++ {
			// Left Side (Old)
			if i < len(pendingMinus) {
				leftBuf.WriteString(fmt.Sprintf("%4d %s\n", leftLineNum, pendingMinus[i]))
				leftLineNum++
			} else {
				leftBuf.WriteString("\n") // Padding
			}

			// Right Side (New)
			if i < len(pendingPlus) {
				rightBuf.WriteString(fmt.Sprintf("%4d %s\n", rightLineNum, pendingPlus[i]))
				rightLineNum++
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
			leftBuf.WriteString(fmt.Sprintf("%s\n    %s\n", separator, line))
			rightBuf.WriteString(fmt.Sprintf("%s\n    %s\n", separator, line))
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
				leftBuf.WriteString(fmt.Sprintf("%4d %s\n", leftLineNum, line))
				rightBuf.WriteString(fmt.Sprintf("%4d %s\n", rightLineNum, line))
				leftLineNum++
				rightLineNum++
			}
		} else {
			// Empty lines in diff are usually context
			flush()
			leftBuf.WriteString(fmt.Sprintf("%4d \n", leftLineNum))
			rightBuf.WriteString(fmt.Sprintf("%4d \n", rightLineNum))
			leftLineNum++
			rightLineNum++
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
			Background(lipgloss.Color("#7D56F4")).
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
			// Get the git diff and parse it.
			diffContent := getGitDiff()
			m.diffData = parseDiff(diffContent)

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

// highlight applies simple color based on the diff marker character
func highlight(content, side string) string {
	lines := strings.Split(content, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if len(line) == 0 {
			sb.WriteString("\n")
			continue
		}

		// Line format: "NNNN [+- ]content" where NNNN is 4-digit line number
		// Check the character after line number (position 5) to determine color
		var markerChar byte = ' '
		if len(line) >= 6 {
			markerChar = line[5] // Position after "NNNN "
		}

		if markerChar == '-' && side == "left" {
			sb.WriteString(delStyle.Render(line) + "\n")
		} else if markerChar == '+' && side == "right" {
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
