package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/titobsala/diffbubble/git"
	"github.com/titobsala/diffbubble/stats"
)

// RenderStatsView renders the full stats screen for branch comparison.
func RenderStatsView(branchStats *stats.BranchStats, width, height int) string {
	if branchStats == nil {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, "No statistics available")
	}

	var sections []string

	// Title with branch comparison
	titleText := fmt.Sprintf("Repository Statistics: %s %s", branchStats.CurrentBranch, branchStats.CompareMode)
	title := TitleStyle.Render(titleText)
	sections = append(sections, title)

	// Branch Changes Section
	branchSection := RenderStatsSection(
		"Changes in Branch",
		branchStats.BranchChanges,
		width-4, // Account for border padding
	)
	sections = append(sections, branchSection)

	// Separator
	separator := strings.Repeat("═", width-4)
	sections = append(sections, HeaderSeparatorStyle.Render(separator))

	// Branch Commits Section
	if branchStats.Commits.TotalCommits > 0 {
		commitsSection := RenderCommitHistory(branchStats.Commits, width-4)
		sections = append(sections, commitsSection)
	} else {
		noCommitsMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color(currentTheme.ContextFg)).
			Render("No unique commits in this branch")
		sections = append(sections, noCommitsMsg)
	}

	// Footer
	footer := RenderStatsFooter(branchStats.CompareMode, width)
	sections = append(sections, footer)

	// Join all sections
	content := strings.Join(sections, "\n\n")

	// Apply border and center
	bordered := BorderStyleFocused.
		Width(width - 4).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Top, bordered)
}

// RenderStatsSection renders a single stats section with file list and summary.
func RenderStatsSection(title string, s stats.Stats, width int) string {
	var lines []string

	// Section title
	sectionTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.TitleFg)).
		Bold(true).
		Render(title + ":")
	lines = append(lines, sectionTitle)

	// Underline
	underline := strings.Repeat("─", len(title)+1)
	lines = append(lines, HeaderSeparatorStyle.Render(underline))
	lines = append(lines, "") // Empty line

	// Check if there are any files
	if s.TotalFiles == 0 {
		lines = append(lines, "  No changes to display")
		return strings.Join(lines, "\n")
	}

	// Summary line
	summary := fmt.Sprintf("  Files Changed: %d   %s %s",
		s.TotalFiles,
		AdditionsStyle.Render(fmt.Sprintf("+%d", s.TotalAdditions)),
		DeletionsStyle.Render(fmt.Sprintf("-%d", s.TotalDeletions)),
	)
	lines = append(lines, summary)
	lines = append(lines, "") // Empty line

	// File list with bar charts (limit to 20 files for display)
	displayLimit := 20
	filesToDisplay := s.FileStats
	if len(filesToDisplay) > displayLimit {
		filesToDisplay = filesToDisplay[:displayLimit]
	}

	// Bar width (leave room for icon, filename, and stats)
	barWidth := 15
	if width > 80 {
		barWidth = 20
	}

	for _, file := range filesToDisplay {
		fileLine := RenderFileStat(file, s.MaxChanges, barWidth)
		lines = append(lines, "  "+fileLine)
	}

	// Show truncation message if needed
	if len(s.FileStats) > displayLimit {
		remaining := len(s.FileStats) - displayLimit
		truncMsg := fmt.Sprintf("  ... and %d more files", remaining)
		lines = append(lines, "")
		lines = append(lines, FooterStyle.Render(truncMsg))
	}

	return strings.Join(lines, "\n")
}

// RenderFileStat renders a single file with status icon, stats, and bar chart.
func RenderFileStat(file git.FileStat, maxChanges, barWidth int) string {
	// Status icon
	icon := statusIcon(file.Status)

	// Truncate filename if too long
	filename := file.Path
	maxFilenameLen := 30
	if len(filename) > maxFilenameLen {
		filename = filename[:maxFilenameLen-3] + "..."
	}

	// Pad filename for alignment
	filename = lipgloss.NewStyle().Width(maxFilenameLen).Render(filename)

	// Stats
	addStr := AdditionsStyle.Render(fmt.Sprintf("+%-4d", file.Additions))
	delStr := DeletionsStyle.Render(fmt.Sprintf("-%-4d", file.Deletions))

	// Bar chart
	bar := RenderBarChart(file.Additions, file.Deletions, maxChanges, barWidth)

	return fmt.Sprintf("%s %s  %s %s  %s", icon, filename, addStr, delStr, bar)
}

// RenderBarChart creates a visual bar chart using block characters.
func RenderBarChart(additions, deletions, maxChanges, width int) string {
	if maxChanges == 0 || width <= 0 {
		return strings.Repeat("░", width)
	}

	// Calculate proportional widths
	totalChanges := additions + deletions
	addWidth := 0
	delWidth := 0

	if totalChanges > 0 {
		// Scale to maxChanges
		addWidth = (additions * width) / maxChanges
		delWidth = (deletions * width) / maxChanges

		// Ensure we don't exceed total width
		if addWidth+delWidth > width {
			ratio := float64(width) / float64(addWidth+delWidth)
			addWidth = int(float64(addWidth) * ratio)
			delWidth = int(float64(delWidth) * ratio)
		}

		// Ensure at least 1 block if there are changes
		if additions > 0 && addWidth == 0 {
			addWidth = 1
		}
		if deletions > 0 && delWidth == 0 {
			delWidth = 1
		}

		// Final bounds check
		if addWidth+delWidth > width {
			if addWidth > delWidth {
				addWidth = width - delWidth
			} else {
				delWidth = width - addWidth
			}
		}
	}

	// Build bar
	var sb strings.Builder

	// Green blocks for additions
	if addWidth > 0 {
		sb.WriteString(AdditionsStyle.Render(strings.Repeat("█", addWidth)))
	}

	// Red blocks for deletions
	if delWidth > 0 {
		sb.WriteString(DeletionsStyle.Render(strings.Repeat("█", delWidth)))
	}

	// Empty blocks for remaining
	remaining := width - addWidth - delWidth
	if remaining > 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.ContextFg))
		sb.WriteString(emptyStyle.Render(strings.Repeat("░", remaining)))
	}

	return sb.String()
}

// RenderCommitHistory renders the list of recent commits.
func RenderCommitHistory(history stats.CommitHistory, width int) string {
	var lines []string

	// Section title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.TitleFg)).
		Bold(true).
		Render(fmt.Sprintf("Recent Commits (%d):", history.TotalCommits))
	lines = append(lines, title)

	// Underline
	underline := strings.Repeat("─", len(fmt.Sprintf("Recent Commits (%d):", history.TotalCommits)))
	lines = append(lines, HeaderSeparatorStyle.Render(underline))
	lines = append(lines, "") // Empty line

	if len(history.Commits) == 0 {
		lines = append(lines, "  No commits found")
		return strings.Join(lines, "\n")
	}

	// Hash style
	hashStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.ModifiedFg)).
		Bold(true)

	// Render each commit
	for _, commit := range history.Commits {
		// Truncate message if too long
		message := commit.Message
		maxMsgLen := 40
		if len(message) > maxMsgLen {
			message = message[:maxMsgLen-3] + "..."
		}

		// Truncate author if too long
		author := commit.Author
		maxAuthorLen := 15
		if len(author) > maxAuthorLen {
			author = author[:maxAuthorLen-3] + "..."
		}

		// Format: hash  message  author  +X -Y
		hash := hashStyle.Render(commit.Hash)
		msgPadded := lipgloss.NewStyle().Width(maxMsgLen).Render(message)
		authorPadded := lipgloss.NewStyle().Width(maxAuthorLen).Render(author)

		statsStr := ""
		if commit.Additions > 0 || commit.Deletions > 0 {
			statsStr = fmt.Sprintf("%s %s",
				AdditionsStyle.Render(fmt.Sprintf("+%d", commit.Additions)),
				DeletionsStyle.Render(fmt.Sprintf("-%d", commit.Deletions)),
			)
		}

		line := fmt.Sprintf("  %s  %s  %s  %s", hash, msgPadded, authorPadded, statsStr)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// RenderStatsFooter renders footer with keyboard hints for stats view.
func RenderStatsFooter(compareMode string, termWidth int) string {
	hints := []string{
		"r: toggle remote",
		"esc: back to diff",
		"q: quit",
		"t: theme",
	}

	// Add current mode indicator
	modeIndicator := fmt.Sprintf("mode: %s", compareMode)
	hints = append([]string{modeIndicator}, hints...)

	footer := strings.Join(hints, " • ")
	return FooterStyle.Render(footer)
}
