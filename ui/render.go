package ui

import (
	"fmt"
	"strconv"
	"strings"

	"diffbuble/git"
	"diffbuble/parser"
)

type Side int

const (
	SideLeft Side = iota
	SideRight
)

// RenderSide turns structured diff rows into a string suitable for a viewport.
func RenderSide(rows []parser.DiffRow, side Side, showLineNumbers bool) string {
	var sb strings.Builder
	width := 0
	if showLineNumbers {
		width = lineNumberWidth(rows, side)
	}

	for _, row := range rows {
		line := rowForSide(row, side)
		if line != nil && line.Kind == parser.LineKindHeader {
			sb.WriteString(renderHeader(line.Content))
			continue
		}

		sb.WriteString(renderLine(line, side, width, showLineNumbers))
		sb.WriteByte('\n')
	}

	return sb.String()
}

// ErrorBox renders a stylized error message that can be embedded inside the layout.
func ErrorBox(err error, width int) string {
	message := fmt.Sprintf("Unable to load git diff.\n\n%s", err)

	maxWidth := width - 6
	if maxWidth < 20 {
		maxWidth = 20
	}

	return ErrorBoxStyle.MaxWidth(maxWidth).Render(message)
}

func renderHeader(content string) string {
	separator := HeaderSeparatorStyle.Render(strings.Repeat("─", 30))
	header := HeaderLineStyle.Render(content)
	return separator + "\n" + header + "\n"
}

func renderLine(line *parser.DiffLine, side Side, width int, showLineNumbers bool) string {
	if line == nil {
		if showLineNumbers {
			return strings.Repeat(" ", width+1)
		}
		return ""
	}

	var text string
	if showLineNumbers && line.Number > 0 {
		number := strconv.Itoa(line.Number)
		text = fmt.Sprintf("%*s %s", width, number, line.Content)
	} else {
		text = line.Content
	}

	switch line.Kind {
	case parser.LineKindAddition:
		if side == SideRight {
			return AddStyle.Render(text)
		}
	case parser.LineKindDeletion:
		if side == SideLeft {
			return DelStyle.Render(text)
		}
	}

	return text
}

func lineNumberWidth(rows []parser.DiffRow, side Side) int {
	max := 0
	for _, row := range rows {
		line := rowForSide(row, side)
		if line != nil && line.Number > max {
			max = line.Number
		}
	}

	if max == 0 {
		return 4
	}

	width := len(strconv.Itoa(max))
	if width < 4 {
		return 4
	}

	return width
}

func rowForSide(row parser.DiffRow, side Side) *parser.DiffLine {
	if side == SideLeft {
		return row.Left
	}
	return row.Right
}

// RenderFileList generates the sidebar content showing all modified files.
func RenderFileList(files []git.FileStat, selectedIdx int) string {
	var sb strings.Builder

	if len(files) == 0 {
		sb.WriteString("No modified files")
		return sb.String()
	}

	for i, file := range files {
		isSelected := (i == selectedIdx)
		sb.WriteString(renderFileListItem(file, isSelected))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func renderFileListItem(file git.FileStat, selected bool) string {
	// Status icon with color
	icon := statusIcon(file.Status)

	// Stats: +5 -2
	stats := fmt.Sprintf("+%d -%d", file.Additions, file.Deletions)
	statsStyled := StatsStyle.Render(stats)

	// Filename (truncate if too long)
	filename := truncate(file.Path, 30)

	line := fmt.Sprintf("%s %-32s %s", icon, filename, statsStyled)

	if selected {
		return SelectedFileStyle.Render(line)
	}
	return FileListItemStyle.Render(line)
}

func statusIcon(status git.FileStatus) string {
	switch status {
	case git.StatusModified:
		return StatusModifiedStyle.Render("M")
	case git.StatusAdded:
		return StatusAddedStyle.Render("A")
	case git.StatusDeleted:
		return StatusDeletedStyle.Render("D")
	case git.StatusRenamed:
		return StatusModifiedStyle.Render("R")
	}
	return "?"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// RenderFooter renders the footer with keyboard shortcuts and feature states.
func RenderFooter(showLineNumbers bool, fullContext bool) string {
	lineNumHint := "on"
	if !showLineNumbers {
		lineNumHint = "off"
	}

	contextHint := "focus"
	if fullContext {
		contextHint = "full"
	}

	text := fmt.Sprintf(
		"j/k: scroll/navigate • n: line numbers (%s) • c: context (%s) • q/esc: quit",
		lineNumHint,
		contextHint,
	)

	return FooterStyle.Render(text)
}
