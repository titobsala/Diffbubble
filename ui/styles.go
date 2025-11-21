package ui

import "github.com/charmbracelet/lipgloss"

// All styles use the current theme colors
var (
	TitleStyle           lipgloss.Style
	BorderStyle          lipgloss.Style
	AddStyle             lipgloss.Style
	DelStyle             lipgloss.Style
	HeaderSeparatorStyle lipgloss.Style
	HeaderLineStyle      lipgloss.Style
	FooterStyle          lipgloss.Style
	ErrorBoxStyle        lipgloss.Style

	// File list sidebar styles
	FileListStyle        lipgloss.Style
	FileListStyleFocused lipgloss.Style
	FileListItemStyle    lipgloss.Style
	SelectedFileStyle    lipgloss.Style

	// Stats styles
	AdditionsStyle      lipgloss.Style
	DeletionsStyle      lipgloss.Style
	DeltaStyle          lipgloss.Style
	StatusModifiedStyle lipgloss.Style
	StatusAddedStyle    lipgloss.Style
	StatusDeletedStyle  lipgloss.Style

	// Border styles for focused/unfocused panes
	BorderStyleFocused   lipgloss.Style
	BorderStyleUnfocused lipgloss.Style
)

// updateStyles applies the current theme to all styles
func updateStyles() {
	theme := currentTheme

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.TitleFg)).
		Background(lipgloss.Color(theme.FocusedBorderColor)).
		Padding(0, 1).
		MarginBottom(1)

	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor))

	AddStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.AdditionFg)).
		Background(lipgloss.Color(theme.AdditionBg))

	DelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.DeletionFg)).
		Background(lipgloss.Color(theme.DeletionBg))

	HeaderSeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.HeaderFg))

	HeaderLineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.HeaderFg)).
		PaddingLeft(4)

	FooterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ContextFg))

	ErrorBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.DeletionFg)).
		Padding(1, 2)

	// File list sidebar styles
	FileListStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		Padding(0, 1)

	FileListStyleFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.FocusedBorderColor)).
		Padding(0, 1)

	FileListItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground))

	SelectedFileStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FocusedBorderColor)).
		Foreground(lipgloss.Color(theme.TitleFg)).
		Bold(true)

	// Stats styles with theme colors
	AdditionsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.AddedFg)).
		Bold(true)

	DeletionsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.DeletedFg)).
		Bold(true)

	DeltaStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ModifiedFg)).
		Bold(true)

	StatusModifiedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ModifiedFg)).
		Bold(true)

	StatusAddedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.AddedFg)).
		Bold(true)

	StatusDeletedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.DeletedFg)).
		Bold(true)

	// Border styles for focused/unfocused panes
	BorderStyleFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.FocusedBorderColor))

	BorderStyleUnfocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor))
}
