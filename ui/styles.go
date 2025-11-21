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
)
