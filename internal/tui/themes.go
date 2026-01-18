package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines a color scheme
type Theme struct {
	Name        string
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Accent      lipgloss.Color
	Muted       lipgloss.Color
	Text        lipgloss.Color
	Warning     lipgloss.Color
	Success     lipgloss.Color
	Background  lipgloss.Color
	SelectedBg  lipgloss.Color
}

// Available themes
var Themes = map[string]Theme{
	"default": {
		Name:       "Default",
		Primary:    lipgloss.Color("#7C3AED"),
		Secondary:  lipgloss.Color("#06B6D4"),
		Accent:     lipgloss.Color("#10B981"),
		Muted:      lipgloss.Color("#6B7280"),
		Text:       lipgloss.Color("#F3F4F6"),
		Warning:    lipgloss.Color("#F59E0B"),
		Success:    lipgloss.Color("#10B981"),
		SelectedBg: lipgloss.Color("#374151"),
	},
	"catppuccin": {
		Name:       "Catppuccin",
		Primary:    lipgloss.Color("#CBA6F7"), // Mauve
		Secondary:  lipgloss.Color("#89DCEB"), // Sky
		Accent:     lipgloss.Color("#A6E3A1"), // Green
		Muted:      lipgloss.Color("#6C7086"), // Overlay0
		Text:       lipgloss.Color("#CDD6F4"), // Text
		Warning:    lipgloss.Color("#F9E2AF"), // Yellow
		Success:    lipgloss.Color("#A6E3A1"), // Green
		SelectedBg: lipgloss.Color("#45475A"), // Surface0
	},
	"dracula": {
		Name:       "Dracula",
		Primary:    lipgloss.Color("#BD93F9"), // Purple
		Secondary:  lipgloss.Color("#8BE9FD"), // Cyan
		Accent:     lipgloss.Color("#50FA7B"), // Green
		Muted:      lipgloss.Color("#6272A4"), // Comment
		Text:       lipgloss.Color("#F8F8F2"), // Foreground
		Warning:    lipgloss.Color("#FFB86C"), // Orange
		Success:    lipgloss.Color("#50FA7B"), // Green
		SelectedBg: lipgloss.Color("#44475A"), // Current Line
	},
	"nord": {
		Name:       "Nord",
		Primary:    lipgloss.Color("#81A1C1"), // Nord9
		Secondary:  lipgloss.Color("#88C0D0"), // Nord8
		Accent:     lipgloss.Color("#A3BE8C"), // Nord14
		Muted:      lipgloss.Color("#4C566A"), // Nord3
		Text:       lipgloss.Color("#ECEFF4"), // Nord6
		Warning:    lipgloss.Color("#EBCB8B"), // Nord13
		Success:    lipgloss.Color("#A3BE8C"), // Nord14
		SelectedBg: lipgloss.Color("#3B4252"), // Nord1
	},
	"tokyonight": {
		Name:       "Tokyo Night",
		Primary:    lipgloss.Color("#BB9AF7"), // Purple
		Secondary:  lipgloss.Color("#7DCFFF"), // Cyan
		Accent:     lipgloss.Color("#9ECE6A"), // Green
		Muted:      lipgloss.Color("#565F89"), // Comment
		Text:       lipgloss.Color("#C0CAF5"), // Foreground
		Warning:    lipgloss.Color("#E0AF68"), // Yellow
		Success:    lipgloss.Color("#9ECE6A"), // Green
		SelectedBg: lipgloss.Color("#292E42"), // BG highlight
	},
	"gruvbox": {
		Name:       "Gruvbox",
		Primary:    lipgloss.Color("#D3869B"), // Purple
		Secondary:  lipgloss.Color("#83A598"), // Aqua
		Accent:     lipgloss.Color("#B8BB26"), // Green
		Muted:      lipgloss.Color("#928374"), // Gray
		Text:       lipgloss.Color("#EBDBB2"), // FG
		Warning:    lipgloss.Color("#FABD2F"), // Yellow
		Success:    lipgloss.Color("#B8BB26"), // Green
		SelectedBg: lipgloss.Color("#3C3836"), // BG1
	},
}

// CurrentTheme holds the active theme
var CurrentTheme = Themes["default"]

// SetTheme changes the active theme and updates all styles
func SetTheme(name string) bool {
	theme, ok := Themes[name]
	if !ok {
		return false
	}
	CurrentTheme = theme
	updateStyles()
	return true
}

// GetThemeNames returns a list of available theme names
func GetThemeNames() []string {
	names := make([]string, 0, len(Themes))
	for name := range Themes {
		names = append(names, name)
	}
	return names
}

// updateStyles refreshes all styles with current theme colors
func updateStyles() {
	PrimaryColor = CurrentTheme.Primary
	SecondaryColor = CurrentTheme.Secondary
	MutedColor = CurrentTheme.Muted
	TextColor = CurrentTheme.Text
	WarningColor = CurrentTheme.Warning
	SuccessColor = CurrentTheme.Success
	SelectedBg = CurrentTheme.SelectedBg

	// Update all styles
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
}
