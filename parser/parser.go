package parser

import (
	"bufio"
	"io"
	"strings"
)

// LineKind represents the semantic meaning of a diff line.
type LineKind int

const (
	LineKindUnknown LineKind = iota
	LineKindContext
	LineKindAddition
	LineKindDeletion
	LineKindHeader
)

// DiffLine represents a single diff line that belongs to either the left or right side.
type DiffLine struct {
	Number  int
	Content string
	Kind    LineKind
}

// DiffRow represents two aligned lines (left/right) in a diff hunk.
type DiffRow struct {
	Left  *DiffLine
	Right *DiffLine
}

// Parse consumes unified diff text from r and returns aligned rows suitable for rendering.
func Parse(r io.Reader) ([]DiffRow, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var rows []DiffRow
	var pendingMinus []DiffLine
	var pendingPlus []DiffLine
	leftLineNum := 1
	rightLineNum := 1

	flush := func() {
		maxLen := len(pendingMinus)
		if len(pendingPlus) > maxLen {
			maxLen = len(pendingPlus)
		}

		for i := 0; i < maxLen; i++ {
			var leftLine *DiffLine
			if i < len(pendingMinus) {
				line := pendingMinus[i]
				line.Number = leftLineNum
				leftLineNum++
				leftLine = &line
			}

			var rightLine *DiffLine
			if i < len(pendingPlus) {
				line := pendingPlus[i]
				line.Number = rightLineNum
				rightLineNum++
				rightLine = &line
			}

			rows = append(rows, DiffRow{
				Left:  leftLine,
				Right: rightLine,
			})
		}

		pendingMinus = nil
		pendingPlus = nil
	}

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "diff"),
			strings.HasPrefix(line, "index"),
			strings.HasPrefix(line, "---"),
			strings.HasPrefix(line, "+++"):
			continue
		case strings.HasPrefix(line, "@@"):
			flush()
			headerLeft := &DiffLine{Content: line, Kind: LineKindHeader}
			headerRight := &DiffLine{Content: line, Kind: LineKindHeader}
			rows = append(rows, DiffRow{
				Left:  headerLeft,
				Right: headerRight,
			})
			continue
		}

		if len(line) == 0 {
			flush()
			left := &DiffLine{
				Number:  leftLineNum,
				Content: "",
				Kind:    LineKindContext,
			}
			right := &DiffLine{
				Number:  rightLineNum,
				Content: "",
				Kind:    LineKindContext,
			}
			leftLineNum++
			rightLineNum++
			rows = append(rows, DiffRow{Left: left, Right: right})
			continue
		}

		switch line[0] {
		case '-':
			pendingMinus = append(pendingMinus, DiffLine{
				Content: line,
				Kind:    LineKindDeletion,
			})
		case '+':
			pendingPlus = append(pendingPlus, DiffLine{
				Content: line,
				Kind:    LineKindAddition,
			})
		case ' ':
			flush()
			left := &DiffLine{
				Number:  leftLineNum,
				Content: line,
				Kind:    LineKindContext,
			}
			right := &DiffLine{
				Number:  rightLineNum,
				Content: line,
				Kind:    LineKindContext,
			}
			leftLineNum++
			rightLineNum++
			rows = append(rows, DiffRow{Left: left, Right: right})
		default:
			// Ignore anything else (e.g. "\ No newline at end of file")
		}
	}

	flush()

	if err := scanner.Err(); err != nil {
		return rows, err
	}

	return rows, nil
}
