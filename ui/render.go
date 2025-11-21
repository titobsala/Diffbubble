package ui

import (
	"fmt"
	"strconv"
	"strings"

	"diffbuble/parser"
)

type Side int

const (
	SideLeft Side = iota
	SideRight
)

// RenderSide turns structured diff rows into a string suitable for a viewport.
func RenderSide(rows []parser.DiffRow, side Side) string {
	var sb strings.Builder
	width := lineNumberWidth(rows, side)

	for _, row := range rows {
		line := rowForSide(row, side)
		if line != nil && line.Kind == parser.LineKindHeader {
			sb.WriteString(renderHeader(line.Content))
			continue
		}

		sb.WriteString(renderLine(line, side, width))
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

func renderLine(line *parser.DiffLine, side Side, width int) string {
	if line == nil {
		return strings.Repeat(" ", width+1)
	}

	number := ""
	if line.Number > 0 {
		number = strconv.Itoa(line.Number)
	}

	text := fmt.Sprintf("%*s %s", width, number, line.Content)

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
