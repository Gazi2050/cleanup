package ui

import "charm.land/lipgloss/v2"

// Style definitions. Colors come from the single theme instance in theme.go.
// Box/card rendering lives in layout.go; this file only holds inline text
// styles used directly inside views.
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(theme.Header).
			Bold(true)

	modeBadgeStyle = lipgloss.NewStyle().
			Foreground(theme.Mode).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(theme.Mode).
			Bold(true)

	doneStyle = lipgloss.NewStyle().Foreground(theme.Done)
	runStyle  = lipgloss.NewStyle().Foreground(theme.Running).Bold(true)
	errStyle  = lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
	pendStyle = lipgloss.NewStyle().Foreground(theme.Pending)

	hintStyle = lipgloss.NewStyle().Foreground(theme.Pending).Italic(true)

	cardTitleStyle = lipgloss.NewStyle().Foreground(theme.Header).Bold(true)
	cardDescStyle  = lipgloss.NewStyle().Foreground(theme.Pending)
	cardMetaStyle  = lipgloss.NewStyle().Foreground(theme.Mode)

	pwErrStyle = lipgloss.NewStyle().Foreground(theme.Error)
)
