package main

import "charm.land/lipgloss/v2"

var (
	colorHeader  = lipgloss.Color("#89DCEB")
	colorMode    = lipgloss.Color("#CBA6F7")
	colorRunning = lipgloss.Color("#F9E2AF")
	colorDone    = lipgloss.Color("#A6E3A1")
	colorError   = lipgloss.Color("#F38BA8")
	colorPending = lipgloss.Color("#6C7086")
	colorBorder  = lipgloss.Color("#313244")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(colorHeader).
			Bold(true)

	modeBadgeStyle = lipgloss.NewStyle().
			Foreground(colorMode).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(colorMode).
			Bold(true)

	doneStyle = lipgloss.NewStyle().Foreground(colorDone)
	runStyle  = lipgloss.NewStyle().Foreground(colorRunning).Bold(true)
	errStyle  = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	pendStyle = lipgloss.NewStyle().Foreground(colorPending)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	hintStyle = lipgloss.NewStyle().Foreground(colorPending).Italic(true)
)
