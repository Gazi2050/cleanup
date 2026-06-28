package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

// BoxVariant selects the border color for a box or card.
type BoxVariant int

const (
	BoxPrimary BoxVariant = iota
	BoxSuccess
	BoxWarning
	BoxError
)

// terminalWidth is updated by the Bubble Tea Update loop on WindowSizeMsg.
// Defaults to 80 until the first size report arrives.
var terminalWidth = 80

// SetTerminalWidth is called from Update on tea.WindowSizeMsg.
func SetTerminalWidth(w int) {
	if w > 0 {
		terminalWidth = w
	}
}

func TerminalWidth() int { return terminalWidth }

// BoxWidthForTerminalWidth clamps a target box width to safe bounds.
func BoxWidthForTerminalWidth(tw int) int {
	if tw <= 0 {
		return 80
	}
	w := tw - 2
	if w < 20 {
		w = tw
	}
	if w > 100 {
		w = 100
	}
	if w > tw {
		w = tw
	}
	if w < 10 {
		w = 10
	}
	return w
}

func boxWidth() int { return BoxWidthForTerminalWidth(TerminalWidth()) }

// BoxInnerWidth is the usable content width inside a padded box.
func BoxInnerWidth() int { return boxWidth() - 6 }

// CardOptions configures RenderCard.
type CardOptions struct {
	Variant   BoxVariant
	Title     string
	Content   string
	FullWidth bool
	MaxWidth  int
}

// variantTitleStyle returns the title style matching a variant's border color.
func variantTitleStyle(v BoxVariant) lipgloss.Style {
	switch v {
	case BoxSuccess:
		return doneStyle
	case BoxError:
		return errStyle
	case BoxWarning:
		return modeBadgeStyle
	default:
		return titleStyle
	}
}

// RenderBox renders a full-width bordered box with an optional bold title
// prepended inside the box.
func RenderBox(variant BoxVariant, title, content string) string {
	border := borderColorForVariant(variant)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(1, 2).
		Width(boxWidth())

	if strings.TrimSpace(title) != "" {
		content = variantTitleStyle(variant).Render(title) + "\n\n" + content
	}
	return box.Render(content)
}

// RenderCard renders an auto-sized bordered card. Width fits content unless
// FullWidth is set; always capped by MaxWidth and terminal width.
func RenderCard(opts CardOptions) string {
	border := borderColorForVariant(opts.Variant)
	maxW := opts.MaxWidth
	if maxW <= 0 {
		maxW = 80
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(0, 1)

	if opts.FullWidth {
		cardStyle = cardStyle.Width(boxWidth())
	} else {
		contentWidth := measureContentWidth(opts.Content, opts.Title)
		w := minInt(maxW, maxInt(contentWidth+6, 30))
		tw := TerminalWidth() - 4
		if w > tw {
			w = tw
		}
		cardStyle = cardStyle.Width(w)
	}

	content := opts.Content
	if strings.TrimSpace(opts.Title) != "" {
		content = variantTitleStyle(opts.Variant).Render(opts.Title) + "\n" + content
	}
	return cardStyle.Render(content)
}

// SuccessCard renders a compact green-bordered card with the success icon.
func SuccessCard(title, body string) string {
	return RenderCard(CardOptions{
		Variant: BoxSuccess,
		Title:   theme.IconSuccess + " " + title,
		Content: body,
	})
}

// ErrorCard renders a compact red-bordered card with the error icon.
func ErrorCard(title, body string) string {
	return RenderCard(CardOptions{
		Variant: BoxError,
		Title:   theme.IconError + " " + title,
		Content: body,
	})
}

// WarningCard renders a compact purple-bordered card with the warning icon.
func WarningCard(title, body string) string {
	return RenderCard(CardOptions{
		Variant: BoxWarning,
		Title:   theme.IconWarning + " " + title,
		Content: body,
	})
}

// measureContentWidth returns the widest visible line in content/title,
// stripping ANSI codes and respecting double-width runes via runewidth.
func measureContentWidth(content, title string) int {
	maxW := 0
	if title != "" {
		for _, line := range strings.Split(runewidth.Truncate(title, 200, ""), "\n") {
			if w := runewidth.StringWidth(line); w > maxW {
				maxW = w
			}
		}
	}
	for _, line := range strings.Split(content, "\n") {
		if w := runewidth.StringWidth(runewidth.Truncate(stripAnsi(line), 200, "")); w > maxW {
			maxW = w
		}
	}
	return maxW
}

// stripAnsi removes ANSI escape sequences so visible width can be measured.
func stripAnsi(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' {
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) {
					if (s[i] >= '0' && s[i] <= '9') || s[i] == ';' || s[i] == '?' {
						i++
					} else if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') {
						i++
						break
					} else {
						break
					}
				}
			}
			continue
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
