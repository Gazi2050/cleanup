package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme holds the single Catppuccin Mocha palette used by cleanup.
// No runtime switching — palette is hardcoded.
type Theme struct {
	Header  color.Color
	Mode    color.Color
	Running color.Color
	Done    color.Color
	Error   color.Color
	Pending color.Color
	Border  color.Color

	IconSuccess string
	IconWarning string
	IconError   string
}

var theme = Theme{
	Header:      lipgloss.Color("#89DCEB"),
	Mode:        lipgloss.Color("#CBA6F7"),
	Running:     lipgloss.Color("#F9E2AF"),
	Done:        lipgloss.Color("#A6E3A1"),
	Error:       lipgloss.Color("#F38BA8"),
	Pending:     lipgloss.Color("#6C7086"),
	Border:      lipgloss.Color("#313244"),
	IconSuccess: "✔",
	IconWarning: "⚠",
	IconError:   "✖",
}

// borderColorForVariant maps a BoxVariant to the theme border color.
func borderColorForVariant(v BoxVariant) color.Color {
	switch v {
	case BoxSuccess:
		return theme.Done
	case BoxWarning:
		return theme.Mode
	case BoxError:
		return theme.Error
	default:
		return theme.Border
	}
}
