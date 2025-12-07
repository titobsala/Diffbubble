package intro

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/titobsala/diffbubble/ui"
)

// RenderAnimation renders the appropriate animation based on type and frame
func RenderAnimation(animType AnimationType, frame, width, height int) string {
	switch animType {
	case Glitch:
		return renderGlitch(frame, width, height)
	case Scan:
		return renderScan(frame, width, height)
	case Stream:
		return renderStream(frame, width, height)
	default:
		return renderGlitch(frame, width, height)
	}
}

// renderGlitch renders the "Merge Conflict" chromatic aberration animation
func renderGlitch(frame, width, height int) string {
	theme := ui.GetTheme()

	// Define colors
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.DeletionFg))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AdditionFg))

	var content string

	// Status text based on frame
	var statusText string
	if frame < 16 {
		statusText = "RESOLVING CONFLICTS..."
	} else {
		statusText = "MERGE SUCCESSFUL ✓"
	}

	// Phase 1: Solid white cat (frames 0-10)
	if frame < 11 {
		var catLines []string
		for _, line := range OctocatASCII {
			catLines = append(catLines, whiteStyle.Render(line))
		}
		catArt := strings.Join(catLines, "\n")
		content = fmt.Sprintf("%s\n\n%s", catArt, statusText)
	} else if frame < 21 {
		// Phase 2: Chromatic aberration (frames 11-20)
		// Calculate shift amount with vibration
		shift := 2
		vibration := 0
		if frame%2 == 0 {
			vibration = rand.Intn(3) - 1 // -1, 0, or 1
		}

		var lines []string
		for i, line := range OctocatASCII {
			// Create three layers
			redLine := redStyle.Render(line)
			whiteLine := whiteStyle.Render(line)
			greenLine := greenStyle.Render(line)

			// Shift red left, green right with vibration
			leftShift := shift + vibration
			rightShift := shift - vibration

			// Combine layers (simplified overlay)
			var combinedLine string
			if i%2 == 0 {
				combinedLine = strings.Repeat(" ", leftShift) + redLine + whiteLine + strings.Repeat(" ", rightShift) + greenLine
			} else {
				combinedLine = redLine + whiteLine + greenLine
			}

lines = append(lines, combinedLine)
		}
		catArt := strings.Join(lines, "\n")
		content = fmt.Sprintf("%s\n\n%s", catArt, statusText)
	} else {
		// Phase 3: Snap back together (frames 21-30)
		var catLines []string
		for _, line := range OctocatASCII {
			catLines = append(catLines, whiteStyle.Render(line))
		}
		catArt := strings.Join(catLines, "\n")
		content = fmt.Sprintf("%s\n\n%s", catArt, statusText)
	}

	// Center the content
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// renderScan renders the "Diff Scan" animation
func renderScan(frame, width, height int) string {
	theme := ui.GetTheme()

	// Define colors
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.DeletionFg))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AdditionFg))
	scannerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	// Scanner position (moves 1 line per 2 frames)
	scannerPos := frame / 2

	var lines []string
	for i, line := range OctocatASCII {
		if i < scannerPos {
			// Scanned: green
			lines = append(lines, greenStyle.Render(line))
		} else if i == scannerPos {
			// Scanner bar
			scanBar := strings.Repeat("█", len(line))
			lines = append(lines, scannerStyle.Render(scanBar))
		} else {
			// Not scanned: red
			lines = append(lines, redStyle.Render(line))
		}
	}

	catArt := strings.Join(lines, "\n")

	// Progress bar
	progress := (scannerPos * 100) / len(OctocatASCII)
	if progress > 100 {
		progress = 100
	}
	filled := progress / 10
	empty := 10 - filled
	progressBar := fmt.Sprintf("[%s%s] %d%%",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		progress)

	statusText := "COMPARING HEAD vs REMOTE..."
	content := fmt.Sprintf("%s\n\n%s\n%s", catArt, statusText, progressBar)

	// Center the content
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// renderStream renders the "Commit Stream" Matrix-style animation
func renderStream(frame, width, height int) string {
	theme := ui.GetTheme()

	// Initialize streams on first frame (simulated with deterministic positions)
	catHeight := len(OctocatLookingUp)
	catY := height - catHeight - 2 // Position cat near bottom

	// Create symbol streams
	symbols := []string{"+", "-", "~", "a1b", "2c3", "d4e"}
	colors := []string{theme.AdditionFg, theme.DeletionFg, "#FFFF00", "#00FFFF", "#00FFFF", "#00FFFF"}

	// Calculate stream positions (deterministic based on frame)
	numStreams := 12
	streams := make([]SymbolStream, numStreams)
	for i := 0; i < numStreams; i++ {
		streams[i] = SymbolStream{
			X:      (i * (width / numStreams)) + (frame%3)*2 - 2, // Slight horizontal movement
			Y:      -i*3 + (frame / 2),                            // Different start times
			Speed:  2,
			Symbol: symbols[i%len(symbols)],
			Color:  colors[i%len(colors)],
		}
	}

	// Create a grid to place symbols
	grid := make([][]string, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]string, width)
		for x := 0; x < width; x++ {
			grid[y][x] = " "
		}
	}

	// Draw streams
	for _, stream := range streams {
		if stream.Y >= 0 && stream.Y < height && stream.X >= 0 && stream.X < width {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(stream.Color))
			grid[stream.Y][stream.X] = style.Render(stream.Symbol)
		}
	}

	// Draw Octocat at bottom
	catStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	// Glow effect: brightness increases with frame
	if frame > 20 {
		catStyle = catStyle.Bold(true)
	}

	catX := (width - len(OctocatLookingUp[0])) / 2
	for i, line := range OctocatLookingUp {
		y := catY + i
		if y >= 0 && y < height {
			for j, char := range line {
				x := catX + j
				if x >= 0 && x < width {
					grid[y][x] = catStyle.Render(string(char))
				}
			}
		}
	}

	// Convert grid to string
	var lines []string
	for _, row := range grid {
		lines = append(lines, strings.Join(row, ""))
	}
	gridStr := strings.Join(lines, "\n")

	// Status text
	statusText := "LOADING COMMIT HISTORY..."
	if frame%20 == 10 {
		statusText = "⚠ force-push detected!"
	}

	// Screen shake for force-push
	shake := 0
	if frame%20 >= 10 && frame%20 <= 12 {
		shake = 2
	}

	content := strings.Repeat(" ", shake) + gridStr + "\n\n" + statusText

	return content
}
