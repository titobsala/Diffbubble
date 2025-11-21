//go:build ignore

package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Simulate the bug
	winWidth := 100
	halfWidth := winWidth / 2 // 50

	// Viewport initialized with width 48
	vpWidth := halfWidth - 2 // 48

	// Box rendered with width 47
	boxWidth := halfWidth - 3 // 47

	fmt.Printf("Viewport Width: %d\n", vpWidth)
	fmt.Printf("Box Width: %d\n", boxWidth)

	// Create a line exactly 48 chars long
	line := strings.Repeat("x", vpWidth)

	// Render it inside a box of width 47
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(boxWidth)

	out := style.Render(line)

	// Count lines in output
	lines := strings.Split(out, "\n")
	fmt.Printf("Output has %d lines (expected 3: top border, content, bottom border)\n", len(lines))
	for _, l := range lines {
		fmt.Println(l)
	}
}
