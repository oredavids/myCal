package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/savioxavier/termlink"
)

// Theme colors
var (
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#06B6D4") // Cyan
	MutedColor     = lipgloss.Color("#6B7280") // Gray
	TextColor      = lipgloss.Color("#F3F4F6") // Light gray
	WarningColor   = lipgloss.Color("#F59E0B") // Amber
	SuccessColor   = lipgloss.Color("#10B981") // Green
	SelectedBg     = lipgloss.Color("#374151") // Dark gray for selection
)

// Styles
var (
	DateStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	GreetingStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Italic(true)

	SectionTitleStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				MarginTop(1)

	EventBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MutedColor).
			Padding(0, 1)

	EventTitleStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true)

	EventTimeStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	EventAllDayStyle = lipgloss.NewStyle().
				Foreground(WarningColor).
				Italic(true)

	NoEventsStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Padding(0, 1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	CountdownStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	SelectedStyle = lipgloss.NewStyle().
			Background(SelectedBg).
			Foreground(TextColor).
			Bold(true).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Italic(true)
)

// SupportsHyperlinks caches the hyperlink support check
var HyperlinkSupport = termlink.SupportsHyperlinks()

// RenderLink creates a clickable hyperlink if supported
func RenderLink(text string, url string, color string) string {
	if HyperlinkSupport {
		return termlink.ColorLink(text, url, color)
	}
	return ""
}

// RenderFallbackURL returns a styled URL for terminals without hyperlink support
func RenderFallbackURL(url string) string {
	style := lipgloss.NewStyle().Foreground(MutedColor)
	return style.Render("  ↳ " + url)
}
