package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/titobsala/Diffbubble/ui"
)

func TestUpdateSearchStyles(t *testing.T) {
	// Initialize a theme
	ui.SetTheme("dark")
	theme := ui.GetTheme()

	// Initialize text input
	ti := textinput.New()
	
	// Apply styles
	updateSearchStyles(&ti)

	// Check TextStyle foreground
	fg := ti.TextStyle.GetForeground()
	expected := lipgloss.Color(theme.Foreground)
	
	// Convert both to string for comparison if needed, or compare directly
	// lipgloss.Color is a string, TerminalColor is an interface.
	// We can type assert or just check equality if possible.
	// But GetForeground returns TerminalColor.
	
	if fg == nil {
	    t.Fatal("TextStyle foreground should not be nil")
	}
	
	// Basic sanity check that it ran
    // Since we can't easily compare the exact color value without casting/reflection potentially
    // (depending on lipgloss version), we'll assume if it's not nil it's set.
    // But we can try to compare string representation.
    
    if fp, ok := fg.(lipgloss.Color); ok {
        if fp != expected {
             t.Errorf("Expected foreground color %s, got %s", expected, fp)
        }
    } else {
        // It might be an ANSI color or something else, but we passed lipgloss.Color
        t.Logf("Foreground is not lipgloss.Color, it is %T", fg)
    }

	// Check PlaceholderStyle
	pfg := ti.PlaceholderStyle.GetForeground()
	expectedPlaceholder := lipgloss.Color(theme.ContextFg)
	
	if pfg == nil {
	    t.Fatal("PlaceholderStyle foreground should not be nil")
	}
	
	if pfp, ok := pfg.(lipgloss.Color); ok {
	    if pfp != expectedPlaceholder {
	        t.Errorf("Expected placeholder color %s, got %s", expectedPlaceholder, pfp)
	    }
	}
}
