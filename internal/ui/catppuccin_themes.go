package ui

import "github.com/charmbracelet/lipgloss"

// Catppuccin color palettes - production-ready themes
var (
	// Catppuccin Mocha (dark theme)
	catppuccinMocha = map[string]string{
		"rosewater": "#f5e0dc",
		"flamingo":  "#f2cdcd", 
		"pink":      "#f5c2e7",
		"mauve":     "#cba6f7",
		"red":       "#f38ba8",
		"maroon":    "#eba0ac",
		"peach":     "#fab387",
		"yellow":    "#f9e2af",
		"green":     "#a6e3a1",
		"teal":      "#94e2d5",
		"sky":       "#89dceb",
		"sapphire":  "#74c7ec",
		"blue":      "#89b4fa",
		"lavender":  "#b4befe",
		"text":      "#cdd6f4",
		"subtext1":  "#bac2de",
		"subtext0":  "#a6adc8",
		"overlay2":  "#9399b2",
		"overlay1":  "#7f849c",
		"overlay0":  "#6c7086",
		"surface2":  "#585b70",
		"surface1":  "#45475a",
		"surface0":  "#313244",
		"base":      "#1e1e2e",
		"mantle":    "#181825",
		"crust":     "#11111b",
	}

	// Catppuccin Latte (light theme)
	catppuccinLatte = map[string]string{
		"rosewater": "#dc8a78",
		"flamingo":  "#dd7878",
		"pink":      "#ea76cb",
		"mauve":     "#8839ef",
		"red":       "#d20f39",
		"maroon":    "#e64553",
		"peach":     "#fe640b",
		"yellow":    "#df8e1d",
		"green":     "#40a02b",
		"teal":      "#179299",
		"sky":       "#04a5e5",
		"sapphire":  "#209fb5",
		"blue":      "#1e66f5",
		"lavender":  "#7287fd",
		"text":      "#4c4f69",
		"subtext1":  "#5c5f77",
		"subtext0":  "#6c6f85",
		"overlay2":  "#7c7f93",
		"overlay1":  "#8c8fa1",
		"overlay0":  "#9ca0b0",
		"surface2":  "#acb0be",
		"surface1":  "#bcc0cc",
		"surface0":  "#ccd0da",
		"base":      "#eff1f5",
		"mantle":    "#e6e9ef",
		"crust":     "#dce0e8",
	}

	// Catppuccin Frappé (warm dark theme)
	catppuccinFrappe = map[string]string{
		"rosewater": "#f2d5cf",
		"flamingo":  "#eebebe",
		"pink":      "#f4b8e4",
		"mauve":     "#ca9ee6",
		"red":       "#e78284",
		"maroon":    "#ea999c",
		"peach":     "#ef9f76",
		"yellow":    "#e5c890",
		"green":     "#a6d189",
		"teal":      "#81c8be",
		"sky":       "#99d1db",
		"sapphire":  "#85c1dc",
		"blue":      "#8caaee",
		"lavender":  "#babbf1",
		"text":      "#c6d0f5",
		"subtext1":  "#b5bfe2",
		"subtext0":  "#a5adce",
		"overlay2":  "#949cbb",
		"overlay1":  "#838ba7",
		"overlay0":  "#737994",
		"surface2":  "#626880",
		"surface1":  "#51576d",
		"surface0":  "#414559",
		"base":      "#303446",
		"mantle":    "#292c3c",
		"crust":     "#232634",
	}

	// Catppuccin Macchiato (cool dark theme)
	catppuccinMacchiato = map[string]string{
		"rosewater": "#f4dbd6",
		"flamingo":  "#f0c6c6",
		"pink":      "#f5bde6",
		"mauve":     "#c6a0f6",
		"red":       "#ed8796",
		"maroon":    "#ee99a0",
		"peach":     "#f5a97f",
		"yellow":    "#eed49f",
		"green":     "#a6da95",
		"teal":      "#8bd5ca",
		"sky":       "#91d7e3",
		"sapphire":  "#7dc4e4",
		"blue":      "#8aadf4",
		"lavender":  "#b7bdf8",
		"text":      "#cad3f5",
		"subtext1":  "#b8c0e0",
		"subtext0":  "#a5adcb",
		"overlay2":  "#939ab7",
		"overlay1":  "#8087a2",
		"overlay0":  "#6e738d",
		"surface2":  "#5b6078",
		"surface1":  "#494d64",
		"surface0":  "#363a4f",
		"base":      "#24273a",
		"mantle":    "#1e2030",
		"crust":     "#181926",
	}
)

// GetCatppuccinTheme returns a fully configured Catppuccin theme
func GetCatppuccinTheme(variant string) Theme {
	var palette map[string]string
	var name string
	
	switch variant {
	case "mocha":
		palette = catppuccinMocha
		name = "Catppuccin Mocha"
	case "latte":
		palette = catppuccinLatte
		name = "Catppuccin Latte"
	case "frappe":
		palette = catppuccinFrappe
		name = "Catppuccin Frappé"
	case "macchiato":
		palette = catppuccinMacchiato
		name = "Catppuccin Macchiato"
	default:
		palette = catppuccinMocha
		name = "Catppuccin Mocha"
	}
	
	return Theme{
		Name: name,
		PromptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])).
			Bold(true),
		OutputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["text"])),
		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["red"])).
			Bold(true),
		InfoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["blue"])),
		SuccessStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])).
			Bold(true),
		WarningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["yellow"])),
		SecondaryStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["subtext0"])),
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color(palette["surface2"])),
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color(palette["base"])),
		HeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["mauve"])).
			Bold(true).
			Background(lipgloss.Color(palette["mantle"])).
			Padding(0, 2),
		FooterStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["overlay1"])).
			Background(lipgloss.Color(palette["mantle"])).
			Padding(0, 2),
		InputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["text"])).
			Background(lipgloss.Color(palette["surface0"])).
			Padding(0, 1),
		CommandStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["peach"])).
			Bold(true),
	}
}

// Enhanced theme with additional styling for Visual Bliss components
type EnhancedTheme struct {
	Theme
	
	// Tab styling
	TabActiveStyle   lipgloss.Style
	TabInactiveStyle lipgloss.Style
	TabSeparator     lipgloss.Style
	
	// Animation styling
	PulseStyle       lipgloss.Style
	GlitchStyle      lipgloss.Style
	MatrixStyle      lipgloss.Style
	
	// Status and presence
	PresenceActive   lipgloss.Style
	PresenceIdle     lipgloss.Style
	PresenceAway     lipgloss.Style
	
	// AI Chat styling
	AIUserStyle      lipgloss.Style
	AIBotStyle       lipgloss.Style
	AIThinkingStyle  lipgloss.Style
	
	// System monitoring
	SystemMetricStyle lipgloss.Style
	SystemAlertStyle  lipgloss.Style
	SystemGoodStyle   lipgloss.Style
	
	// Network monitoring
	NetworkConnectedStyle    lipgloss.Style
	NetworkDisconnectedStyle lipgloss.Style
	NetworkTransferStyle     lipgloss.Style
}

// GetEnhancedCatppuccinTheme returns an enhanced theme with all Visual Bliss components
func GetEnhancedCatppuccinTheme(variant string) EnhancedTheme {
	baseTheme := GetCatppuccinTheme(variant)
	
	var palette map[string]string
	switch variant {
	case "mocha":
		palette = catppuccinMocha
	case "latte":
		palette = catppuccinLatte
	case "frappe":
		palette = catppuccinFrappe
	case "macchiato":
		palette = catppuccinMacchiato
	default:
		palette = catppuccinMocha
	}
	
	return EnhancedTheme{
		Theme: baseTheme,
		
		// Tab styling
		TabActiveStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["base"])).
			Background(lipgloss.Color(palette["mauve"])).
			Bold(true).
			Padding(0, 2).
			MarginRight(1),
			
		TabInactiveStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["subtext0"])).
			Background(lipgloss.Color(palette["surface0"])).
			Padding(0, 2).
			MarginRight(1),
			
		TabSeparator: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["surface2"])),
		
		// Animation styling
		PulseStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])).
			Bold(true),
			
		GlitchStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["red"])).
			Background(lipgloss.Color(palette["surface0"])).
			Blink(true),
			
		MatrixStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])).
			Faint(true),
		
		// Presence styling
		PresenceActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])).
			Bold(true),
			
		PresenceIdle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["yellow"])),
			
		PresenceAway: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["overlay1"])),
		
		// AI Chat styling
		AIUserStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["blue"])).
			Bold(true),
			
		AIBotStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["mauve"])).
			Bold(true),
			
		AIThinkingStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["overlay1"])).
			Italic(true),
		
		// System monitoring
		SystemMetricStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["text"])),
			
		SystemAlertStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["red"])).
			Bold(true),
			
		SystemGoodStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])),
		
		// Network monitoring
		NetworkConnectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["green"])),
			
		NetworkDisconnectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["red"])),
			
		NetworkTransferStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["blue"])),
	}
}

// GetCatppuccinThemeNames returns all available Catppuccin theme names
func GetCatppuccinThemeNames() []string {
	return []string{"mocha", "latte", "frappe", "macchiato"}
}

// IsValidCatppuccinTheme checks if a Catppuccin theme name is valid
func IsValidCatppuccinTheme(name string) bool {
	for _, validName := range GetCatppuccinThemeNames() {
		if name == validName {
			return true
		}
	}
	return false
}

// GetRoséPineTheme returns a Rosé Pine theme (alternative to Catppuccin)
func GetRoséPineTheme(variant string) Theme {
	var palette map[string]string
	var name string
	
	switch variant {
	case "dawn": // Light theme
		palette = map[string]string{
			"base":      "#faf4ed",
			"surface":   "#fffaf3",
			"overlay":   "#f2e9e1",
			"muted":     "#9893a5",
			"subtle":    "#797593",
			"text":      "#575279",
			"love":      "#b4637a",
			"gold":      "#ea9d34",
			"rose":      "#d7827e",
			"pine":      "#286983",
			"foam":      "#56949f",
			"iris":      "#907aa9",
			"highlight": "#eee9e6",
		}
		name = "Rosé Pine Dawn"
		
	case "moon": // Dark theme
		palette = map[string]string{
			"base":      "#232136",
			"surface":   "#2a273f",
			"overlay":   "#393552",
			"muted":     "#6e6a86",
			"subtle":    "#908caa",
			"text":      "#e0def4",
			"love":      "#eb6f92",
			"gold":      "#f6c177",
			"rose":      "#ea9a97",
			"pine":      "#3e8fb0",
			"foam":      "#9ccfd8",
			"iris":      "#c4a7e7",
			"highlight": "#2a283e",
		}
		name = "Rosé Pine Moon"
		
	default: // Main (dark theme)
		palette = map[string]string{
			"base":      "#191724",
			"surface":   "#1f1d2e",
			"overlay":   "#26233a",
			"muted":     "#6e6a86",
			"subtle":    "#908caa",
			"text":      "#e0def4",
			"love":      "#eb6f92",
			"gold":      "#f6c177",
			"rose":      "#ebbcba",
			"pine":      "#31748f",
			"foam":      "#9ccfd8",
			"iris":      "#c4a7e7",
			"highlight": "#2a283e",
		}
		name = "Rosé Pine"
	}
	
	return Theme{
		Name: name,
		PromptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["pine"])).
			Bold(true),
		OutputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["text"])),
		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["love"])).
			Bold(true),
		InfoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["foam"])),
		SuccessStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["pine"])).
			Bold(true),
		WarningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["gold"])),
		SecondaryStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["subtle"])),
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color(palette["overlay"])),
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color(palette["base"])),
		HeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["iris"])).
			Bold(true).
			Background(lipgloss.Color(palette["surface"])).
			Padding(0, 2),
		FooterStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["muted"])).
			Background(lipgloss.Color(palette["surface"])).
			Padding(0, 2),
		InputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["text"])).
			Background(lipgloss.Color(palette["surface"])).
			Padding(0, 1),
		CommandStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette["rose"])).
			Bold(true),
	}
}