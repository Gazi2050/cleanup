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

const cardWidth = 44

var (
	cardTitleStyle = lipgloss.NewStyle().
			Foreground(colorHeader).
			Bold(true)

	cardDescStyle = lipgloss.NewStyle().Foreground(colorPending)

	cardMetaStyle = lipgloss.NewStyle().Foreground(colorMode)

	cardSelected = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMode).
			Width(cardWidth).
			Padding(0, 2)

	cardUnselected = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Width(cardWidth).
			Padding(0, 2)

	toastBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDone).
			Padding(0, 2).
			MarginLeft(1)

	sudoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMode).
			Padding(1, 2)

	pwPromptStyle = lipgloss.NewStyle().Foreground(colorHeader).Bold(true)
	pwErrStyle    = lipgloss.NewStyle().Foreground(colorError)
)
