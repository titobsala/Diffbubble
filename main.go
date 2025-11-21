package main

import (
	"bytes"
	"fmt"

	"diffbuble/git"
	"diffbuble/parser"
	"diffbuble/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const appTitle = "Git Diff Side-by-Side"

type model struct {
	ready     bool
	leftView  viewport.Model
	rightView viewport.Model
	rows      []parser.DiffRow
	err       error
	winWidth  int
	winHeight int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height

		// Calculate dimensions
		headerHeight := 2
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.ready = true

			diffOutput, err := git.Diff()
			if err != nil {
				m.err = err
			} else {
				rows, parseErr := parser.Parse(bytes.NewReader(diffOutput))
				if parseErr != nil {
					m.err = parseErr
				} else {
					m.rows = rows
				}
			}

			m.leftView = viewport.New(msg.Width/2-2, msg.Height-verticalMarginHeight)
			m.rightView = viewport.New(msg.Width/2-2, msg.Height-verticalMarginHeight)

			if m.err == nil {
				m.leftView.SetContent(ui.RenderSide(m.rows, ui.SideLeft))
				m.rightView.SetContent(ui.RenderSide(m.rows, ui.SideRight))
			}
		} else {
			// Handle resize
			m.leftView.Width = msg.Width/2 - 2
			m.leftView.Height = msg.Height - verticalMarginHeight
			m.rightView.Width = msg.Width/2 - 2
			m.rightView.Height = msg.Height - verticalMarginHeight
		}
	}

	// Sync Scrolling: Update both viewports with the same message
	// This ensures if you press "down" on one, both move.
	m.leftView, cmd = m.leftView.Update(msg)
	cmds = append(cmds, cmd)

	m.rightView, _ = m.rightView.Update(msg)
	// We don't append the right view's command to avoid duplicate key handling
	// artifacts if they both processed the same key, though for viewports it's usually fine.
	// To keep them perfectly synced, we force the Y offset to match.
	m.rightView.YOffset = m.leftView.YOffset

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := ui.TitleStyle.Render(appTitle)
	footer := ui.FooterStyle.Render("Scroll with j/k/arrows • q to quit")

	if m.err != nil {
		errorBox := ui.ErrorBox(m.err, m.winWidth)
		return lipgloss.JoinVertical(lipgloss.Top, header, errorBox, footer)
	}

	leftBox := ui.BorderStyle.Width(m.winWidth/2 - 3).Render(m.leftView.View())
	rightBox := ui.BorderStyle.Width(m.winWidth/2 - 3).Render(m.rightView.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	return lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
}

func main() {
	p := tea.NewProgram(
		model{},
		tea.WithAltScreen(),       // Use full screen
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
