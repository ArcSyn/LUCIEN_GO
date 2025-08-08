package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the complete visual styling for the Lucien CLI interface
type Theme struct {
	Name           string
	PromptStyle    lipgloss.Style
	OutputStyle    lipgloss.Style
	ErrorStyle     lipgloss.Style
	InfoStyle      lipgloss.Style
	SuccessStyle   lipgloss.Style
	WarningStyle   lipgloss.Style
	SecondaryStyle lipgloss.Style
	BorderStyle    lipgloss.Style
	Background     lipgloss.Style
	HeaderStyle    lipgloss.Style
	FooterStyle    lipgloss.Style
	InputStyle     lipgloss.Style
	CommandStyle   lipgloss.Style
}

// GetTheme returns a fully configured theme by name
func GetTheme(name string) Theme {
	switch name {
	case "nexus":
		return Theme{
			Name: "nexus",
			PromptStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff41")).
				Bold(true),
			OutputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")),
			ErrorStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0066")).
				Bold(true),
			InfoStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffaa00")),
			SuccessStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff00")).
				Bold(true),
			WarningStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffaa00")),
			SecondaryStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0066cc")),
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), true).
				BorderForeground(lipgloss.Color("#00ff41")),
			Background: lipgloss.NewStyle().
				Background(lipgloss.Color("#000000")),
			HeaderStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff41")).
				Bold(true).
				Background(lipgloss.Color("#000000")),
			FooterStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0066cc")).
				Italic(true),
			InputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#000000")),
			CommandStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff41")),
		}

	case "synthwave":
		return Theme{
			Name: "synthwave",
			PromptStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff5fd2")).
				Bold(true),
			OutputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0affef")),
			ErrorStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff005b")).
				Bold(true),
			InfoStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffe762")),
			SuccessStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#39ff14")).
				Bold(true),
			WarningStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffb000")),
			SecondaryStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")),
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true).
				BorderForeground(lipgloss.Color("#ff00ff")),
			Background: lipgloss.NewStyle().
				Background(lipgloss.Color("#1a0033")),
			HeaderStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff5fd2")).
				Bold(true).
				Background(lipgloss.Color("#1a0033")),
			FooterStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Italic(true),
			InputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#1a0033")),
			CommandStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff5fd2")),
		}

	case "ghost":
		return Theme{
			Name: "ghost",
			PromptStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Italic(true),
			OutputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cccccc")),
			ErrorStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cc0000")).
				Bold(true),
			InfoStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ccaa00")),
			SuccessStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00cc00")).
				Bold(true),
			WarningStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ccaa00")),
			SecondaryStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#aaaaaa")),
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(lipgloss.Color("#666666")),
			Background: lipgloss.NewStyle().
				Background(lipgloss.Color("#000000")),
			HeaderStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true).
				Background(lipgloss.Color("#000000")),
			FooterStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Italic(true),
			InputStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cccccc")).
				Background(lipgloss.Color("#000000")),
			CommandStyle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")),
		}

	default:
		return GetTheme("nexus") // Default to nexus theme
	}
}

// GetThemeNames returns all available theme names
func GetThemeNames() []string {
	return []string{"nexus", "synthwave", "ghost"}
}

// IsValidTheme checks if a theme name is valid
func IsValidTheme(name string) bool {
	for _, validName := range GetThemeNames() {
		if name == validName {
			return true
		}
	}
	return false
}