package ui

import "github.com/charmbracelet/lipgloss"

var (
	boxStyle      = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.NormalBorder())
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
)
