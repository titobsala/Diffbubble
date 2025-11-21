package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	AddStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D"))

	DelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E05252"))

	HeaderSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	HeaderLineStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(4)

	FooterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E05252")).
			Padding(1, 2)

	// File list sidebar styles
	FileListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	FileListItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	SelectedFileStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FAFAFA")).
				Bold(true)

	StatsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	StatusModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5C842")).
				Bold(true)

	StatusAddedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#43BF6D")).
				Bold(true)

	StatusDeletedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E05252")).
				Bold(true)
)
