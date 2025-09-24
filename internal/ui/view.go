package ui

import (
	"strings"
	"github.com/chetanjangir0/blueboy/internal/network"
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var title, main string

	if m.password.isAsking {
		title = "Enter Password"
		main = m.password.View()
	} else {
		switch m.CurrentMenu {
		case MainMenu:
			title = "Blueboy"
			main = renderList(m.MainOptions, m.cursor)
		case ScanMenu:
			title = "Scan Results"
			main = renderList(m.ScanResults, m.cursor)
		case PairedMenu:
			title = "Paired Devices"
			main = renderList(m.PairedDevices, m.cursor)
		}
	}

	return layoutBox(title, main, m.status, m.width, m.height)
}

func renderList[T string | network.Device](list []T, cursor int) string {
	var out []string
	for i, item := range list {
		var line string
		switch v := any(item).(type) {
		case string:
			line = v
		case network.Device:
			line = v.Name
		}
		if i == cursor {
			line = selectedStyle.Render("> " + line)
		} else {
			line = "  " + line
		}
		out = append(out, line)
	}
	return lipgloss.JoinVertical(lipgloss.Left, out...)
}

func layoutBox(title, main string, status string, width, height int) string {
	// Title styled
	titleLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Bold(true).
		Render(title)

	// Separator as full-width horizontal line
	separatorWidth := lipgloss.Width(titleLine)
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(strings.Repeat("─", separatorWidth))

	// Body composed with precise spacing
	body := lipgloss.JoinVertical(lipgloss.Left,
		titleLine,
		separator,
		main,
		lipgloss.NewStyle().MarginTop(1).Render(status),
		lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color("8")).Render("q: Quit, b/esc: Main menu."),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, boxStyle.Render(body))
}
