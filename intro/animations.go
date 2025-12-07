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
	default:
		return renderGlitch(frame, width, height)
	}
}

// renderCatWithScanlines adds subtle scan line effect to cat
func renderCatWithScanlines(cat []string, intense bool) string {
	var lines []string
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	for i, line := range cat {
		if intense && i%2 == 0 {
			lines = append(lines, dimStyle.Render(line))
		} else if intense && i%3 == 0 {
			lines = append(lines, dimStyle.Render(line))
		} else {
			lines = append(lines, normalStyle.Render(line))
		}
	}
	return strings.Join(lines, "\n")
}

// applyDigitalCorruption randomly corrupts characters in the cat
func applyDigitalCorruption(cat []string, seed int) []string {
	corrupted := make([]string, len(cat))
	glitchChars := []rune{'░', '▒', '▓', '█', '▀', '▄'}
	r := rand.New(rand.NewSource(int64(seed * 1000)))

	for i, line := range cat {
		runes := []rune(line)
		// Corrupt 5-10% of characters
		numCorrupt := len(runes) / 15
		for j := 0; j < numCorrupt; j++ {
			pos := r.Intn(len(runes))
			if runes[pos] != ' ' {
				runes[pos] = glitchChars[r.Intn(len(glitchChars))]
			}
		}
		corrupted[i] = string(runes)
	}
	return corrupted
}

// renderWithFlicker applies flicker effect with color alternation
func renderWithFlicker(cat []string, frame int) string {
	var style lipgloss.Style
	// Flicker between cyan and magenta
	if frame%2 == 0 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	}

	var lines []string
	for _, line := range cat {
		lines = append(lines, style.Render(line))
	}
	return strings.Join(lines, "\n")
}

// shiftLines shifts all lines horizontally by amount
func shiftLines(cat []string, amount int, addNoise bool) []string {
	shifted := make([]string, len(cat))
	r := rand.New(rand.NewSource(int64(amount * 42)))

	for i, line := range cat {
		variance := 0
		if addNoise {
			variance = r.Intn(3) - 1 // -1, 0, or 1
		}
		totalShift := amount + variance

		if totalShift > 0 {
			// Shift right
			shifted[i] = strings.Repeat(" ", totalShift) + line
		} else if totalShift < 0 {
			// Shift left (truncate)
			absShift := -totalShift
			if absShift < len(line) {
				shifted[i] = line[absShift:]
			} else {
				shifted[i] = ""
			}
		} else {
			shifted[i] = line
		}
	}
	return shifted
}

// applyLineDisplacement randomly shifts individual lines for tearing effect
func applyLineDisplacement(cat []string, seed int) []string {
	displaced := make([]string, len(cat))
	r := rand.New(rand.NewSource(int64(seed * 777)))

	for i, line := range cat {
		// Randomly displace some lines
		if r.Float32() < 0.3 { // 30% chance
			displacement := r.Intn(8) - 4 // -4 to +3
			if displacement > 0 {
				displaced[i] = strings.Repeat(" ", displacement) + line
			} else if displacement < 0 {
				absDisp := -displacement
				if absDisp < len(line) {
					displaced[i] = line[absDisp:]
				} else {
					displaced[i] = ""
				}
			} else {
				displaced[i] = line
			}
		} else {
			displaced[i] = line
		}
	}
	return displaced
}

// combineRGBChannels overlays three color channels with RGB colors
func combineRGBChannels(red, green, cyan []string) string {
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	maxLen := len(red)
	if len(green) > maxLen {
		maxLen = len(green)
	}
	if len(cyan) > maxLen {
		maxLen = len(cyan)
	}

	var result []string
	for i := 0; i < maxLen; i++ {
		var redLine, greenLine, cyanLine string
		if i < len(red) {
			redLine = red[i]
		}
		if i < len(green) {
			greenLine = green[i]
		}
		if i < len(cyan) {
			cyanLine = cyan[i]
		}

		// Find max line length
		maxLineLen := max(len(redLine), max(len(greenLine), len(cyanLine)))

		// Build combined line by overlaying characters
		combined := make([]rune, maxLineLen)
		for j := 0; j < maxLineLen; j++ {
			// Priority: non-space characters from any channel
			var char rune = ' '

			if j < len(redLine) && redLine[j] != ' ' {
				char = rune(redLine[j])
				combined[j] = rune(redStyle.Render(string(char))[0]) // Simplified
			} else if j < len(greenLine) && greenLine[j] != ' ' {
				char = rune(greenLine[j])
				combined[j] = char
			} else if j < len(cyanLine) && cyanLine[j] != ' ' {
				char = rune(cyanLine[j])
				combined[j] = char
			} else {
				combined[j] = ' '
			}
		}

		// Apply colors properly by rendering each channel separately
		var coloredLine string
		if i < len(red) && redLine != "" {
			coloredLine += redStyle.Render(redLine)
		}
		if i < len(green) && greenLine != "" {
			coloredLine += greenStyle.Render(greenLine)
		}
		if i < len(cyan) && cyanLine != "" {
			coloredLine += cyanStyle.Render(cyanLine)
		}

		result = append(result, coloredLine)
	}
	return strings.Join(result, "\n")
}

// addScanLines adds animated horizontal scan lines
func addScanLines(content string, frame int) string {
	lines := strings.Split(content, "\n")
	scanPos := frame % len(lines) // Scan line position

	for i := range lines {
		if i == scanPos {
			// Add bright scan line
			lines[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Render(lines[i])
		} else if i%2 == 0 {
			// Darken alternating lines
			lines[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render(lines[i])
		}
	}
	return strings.Join(lines, "\n")
}

// glitchText corrupts status text with glitch characters
func glitchText(text string, seed int) string {
	r := rand.New(rand.NewSource(int64(seed * 999)))
	runes := []rune(text)
	glitchChars := []rune{'█', '▓', '▒', '░', '▀', '▄', '◆', '◇', '▪', '▫'}

	// Corrupt 20% of characters
	numCorrupt := len(runes) / 5
	for i := 0; i < numCorrupt; i++ {
		pos := r.Intn(len(runes))
		runes[pos] = glitchChars[r.Intn(len(glitchChars))]
	}
	return string(runes)
}

// renderCatWithGlow renders cat with glow effect
func renderCatWithGlow(cat []string) string {
	glowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	var lines []string
	for _, line := range cat {
		lines = append(lines, glowStyle.Render(line))
	}
	return strings.Join(lines, "\n")
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// abs returns absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// renderGlitch renders the enhanced cyberpunk glitch animation
func renderGlitch(frame, width, height int) string {
	theme := ui.GetTheme()
	_ = theme // Will use theme colors later

	var content string
	var statusText string

	// Phase 1: Stable Display with subtle scanlines (0-10)
	if frame < 11 {
		content = renderCatWithScanlines(PixelCat, false)
		statusText = "INITIALIZING..."
	} else if frame < 16 {
		// Phase 2: Digital Corruption (11-15)
		corruptedCat := applyDigitalCorruption(PixelCat, frame)
		content = renderWithFlicker(corruptedCat, frame)
		statusText = "DETECTING CONFLICTS..."
	} else if frame < 26 {
		// Phase 3: RGB Channel Split - MAXIMUM CHAOS (16-25)
		// Separate into RGB channels with increasing shift
		shiftAmount := 3 + (frame-16)/2 // Increases from 3 to 7

		redChannel := shiftLines(PixelCat, -shiftAmount, true) // left
		greenChannel := shiftLines(PixelCat, 0, false)         // center
		cyanChannel := shiftLines(PixelCat, shiftAmount, true) // right

		// Apply line displacement for tearing effect
		redChannel = applyLineDisplacement(redChannel, frame)
		cyanChannel = applyLineDisplacement(cyanChannel, frame+5)

		// Combine RGB channels
		combined := combineRGBChannels(redChannel, greenChannel, cyanChannel)

		// Add scan lines
		content = addScanLines(combined, frame)

		// Screen shake
		shake := (frame % 4) - 2 // Oscillates: -2, -1, 0, 1
		if shake != 0 {
			content = strings.Repeat(" ", abs(shake)) + content
		}

		// Glitch the status text
		statusText = glitchText("RESOLVING CONFLICTS...", frame)
	} else if frame < 31 {
		// Phase 4: Stabilization (26-30)
		content = renderCatWithScanlines(PixelCat, false)
		statusText = "MERGE SUCCESSFUL ✓"
	} else {
		// Phase 5: Final Display with glow (31-35)
		content = renderCatWithGlow(PixelCat)
		statusText = "MERGE SUCCESSFUL ✓"
	}

	// Center everything
	fullContent := content + "\n\n" + statusText
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, fullContent)
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
	for i, line := range PixelCat {
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
	progress := (scannerPos * 100) / len(PixelCat)
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
