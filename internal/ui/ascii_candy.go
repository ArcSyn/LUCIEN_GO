package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

// ASCIICandy provides ASCII art, figlet text, cowsay, and other visual enhancements
type ASCIICandy struct {
	themes map[string]ASCIITheme
	config ASCIICandyConfig
}

// ASCIITheme defines colors and styles for ASCII art
type ASCIITheme struct {
	Name        string
	PrimaryColor   lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	GradientColors []lipgloss.Color
}

// ASCIICandyConfig holds configuration for ASCII effects
type ASCIICandyConfig struct {
	EnableAnimations bool
	AnimationSpeed   time.Duration
	DefaultFont      string
	EnableGradients  bool
	EnableShadows    bool
	MaxWidth         int
}

// NewASCIICandy creates a new ASCII candy system
func NewASCIICandy() *ASCIICandy {
	ac := &ASCIICandy{
		themes: make(map[string]ASCIITheme),
		config: ASCIICandyConfig{
			EnableAnimations: true,
			AnimationSpeed:   100 * time.Millisecond,
			DefaultFont:      "doom",
			EnableGradients:  true,
			EnableShadows:    true,
			MaxWidth:         120,
		},
	}
	
	ac.initializeThemes()
	return ac
}

// initializeThemes sets up predefined ASCII themes
func (ac *ASCIICandy) initializeThemes() {
	// Neon theme
	ac.themes["neon"] = ASCIITheme{
		Name:           "Neon",
		PrimaryColor:   lipgloss.Color("#00ff41"),
		SecondaryColor: lipgloss.Color("#0080ff"),
		AccentColor:    lipgloss.Color("#ff0080"),
		GradientColors: []lipgloss.Color{
			lipgloss.Color("#00ff41"),
			lipgloss.Color("#40ff80"),
			lipgloss.Color("#80ffbf"),
			lipgloss.Color("#bfffff"),
		},
	}
	
	// Fire theme
	ac.themes["fire"] = ASCIITheme{
		Name:           "Fire",
		PrimaryColor:   lipgloss.Color("#ff4500"),
		SecondaryColor: lipgloss.Color("#ff8c00"),
		AccentColor:    lipgloss.Color("#ffd700"),
		GradientColors: []lipgloss.Color{
			lipgloss.Color("#ff0000"),
			lipgloss.Color("#ff4500"),
			lipgloss.Color("#ff8c00"),
			lipgloss.Color("#ffd700"),
		},
	}
	
	// Ice theme
	ac.themes["ice"] = ASCIITheme{
		Name:           "Ice",
		PrimaryColor:   lipgloss.Color("#87ceeb"),
		SecondaryColor: lipgloss.Color("#b0e0e6"),
		AccentColor:    lipgloss.Color("#e6f3ff"),
		GradientColors: []lipgloss.Color{
			lipgloss.Color("#4169e1"),
			lipgloss.Color("#6495ed"),
			lipgloss.Color("#87ceeb"),
			lipgloss.Color("#b0e0e6"),
		},
	}
	
	// Purple haze theme
	ac.themes["purple"] = ASCIITheme{
		Name:           "Purple Haze",
		PrimaryColor:   lipgloss.Color("#8a2be2"),
		SecondaryColor: lipgloss.Color("#9370db"),
		AccentColor:    lipgloss.Color("#ba55d3"),
		GradientColors: []lipgloss.Color{
			lipgloss.Color("#4b0082"),
			lipgloss.Color("#8a2be2"),
			lipgloss.Color("#9370db"),
			lipgloss.Color("#ba55d3"),
		},
	}
}

// GenerateFigletText creates styled figlet text
func (ac *ASCIICandy) GenerateFigletText(text, font, theme string) string {
	if font == "" {
		font = ac.config.DefaultFont
	}
	
	// Generate figlet text
	myFigure := figure.NewFigure(text, font, true)
	asciiText := myFigure.String()
	
	// Apply theme styling
	return ac.ApplyThemeToText(asciiText, theme)
}

// ApplyThemeToText applies color themes to ASCII text
func (ac *ASCIICandy) ApplyThemeToText(text, themeName string) string {
	theme, exists := ac.themes[themeName]
	if !exists {
		theme = ac.themes["neon"] // Default theme
	}
	
	lines := strings.Split(text, "\n")
	var styledLines []string
	
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			styledLines = append(styledLines, line)
			continue
		}
		
		if ac.config.EnableGradients && len(theme.GradientColors) > 1 {
			// Apply gradient effect
			colorIndex := i % len(theme.GradientColors)
			style := lipgloss.NewStyle().Foreground(theme.GradientColors[colorIndex])
			styledLines = append(styledLines, style.Render(line))
		} else {
			// Apply single color
			style := lipgloss.NewStyle().Foreground(theme.PrimaryColor)
			styledLines = append(styledLines, style.Render(line))
		}
	}
	
	return strings.Join(styledLines, "\n")
}

// GenerateStartupBanner creates an enhanced startup banner
func (ac *ASCIICandy) GenerateStartupBanner() string {
	var banner strings.Builder
	
	// Main title
	title := ac.GenerateFigletText("LUCIEN", "doom", "neon")
	banner.WriteString(title)
	banner.WriteString("\n")
	
	// Subtitle
	subtitle := ac.GenerateFigletText("Neural Shell", "small", "purple")
	banner.WriteString(subtitle)
	banner.WriteString("\n")
	
	// Version info with styling
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#40ff80")).
		Bold(true).
		Padding(0, 2)
	
	banner.WriteString(versionStyle.Render("🚀 v1.0.0-nexus7 - Visual Bliss Enhanced"))
	banner.WriteString("\n\n")
	
	// Status indicators
	indicators := []string{
		"🧠 Neural Networks: ONLINE",
		"🎨 Visual Bliss: ACTIVE", 
		"🔊 Audio Feedback: ENABLED",
		"🌐 Web Interface: READY",
		"⚡ Quantum Core: CHARGED",
	}
	
	indicatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff41")).
		Padding(0, 1)
	
	for _, indicator := range indicators {
		banner.WriteString(indicatorStyle.Render(indicator))
		banner.WriteString("\n")
	}
	
	return banner.String()
}

// GenerateMatrixRain creates animated matrix-style text effect
func (ac *ASCIICandy) GenerateMatrixRain(width, height int) []string {
	chars := []rune("01アイウエオカキクケコサシスセソタチツテト")
	var lines []string
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff41"))
	
	for i := 0; i < height; i++ {
		var line strings.Builder
		for j := 0; j < width; j++ {
			if rand.Float64() < 0.1 { // 10% chance for character
				char := chars[rand.Intn(len(chars))]
				line.WriteRune(char)
			} else {
				line.WriteString(" ")
			}
		}
		lines = append(lines, style.Render(line.String()))
	}
	
	return lines
}

// GenerateProgressBar creates an enhanced progress bar
func (ac *ASCIICandy) GenerateProgressBar(current, max int, width int, theme string) string {
	if width <= 0 {
		width = 50
	}
	
	percentage := float64(current) / float64(max)
	filled := int(percentage * float64(width))
	
	asciiTheme, exists := ac.themes[theme]
	if !exists {
		asciiTheme = ac.themes["neon"]
	}
	
	var bar strings.Builder
	
	// Progress bar styling
	filledStyle := lipgloss.NewStyle().
		Foreground(asciiTheme.PrimaryColor).
		Bold(true)
	
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333333"))
	
	bar.WriteString("[")
	
	// Filled portion
	for i := 0; i < filled; i++ {
		bar.WriteString(filledStyle.Render("█"))
	}
	
	// Empty portion  
	for i := filled; i < width; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}
	
	bar.WriteString("]")
	
	// Percentage display
	percentStyle := lipgloss.NewStyle().
		Foreground(asciiTheme.AccentColor).
		Bold(true)
	
	bar.WriteString(" ")
	bar.WriteString(percentStyle.Render(fmt.Sprintf("%.1f%%", percentage*100)))
	
	return bar.String()
}

// GenerateGlitchText creates glitch effect text
func (ac *ASCIICandy) GenerateGlitchText(text string) string {
	glitchChars := []rune("!@#$%^&*()_+-=[]{}|;:,.<>?")
	lines := strings.Split(text, "\n")
	var glitchedLines []string
	
	for _, line := range lines {
		var glitchedLine strings.Builder
		runes := []rune(line)
		
		for _, r := range runes {
			if rand.Float64() < 0.05 { // 5% chance of glitch
				glitchChar := glitchChars[rand.Intn(len(glitchChars))]
				
				// Apply glitch styling
				glitchStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#ff0080")).
					Background(lipgloss.Color("#000080")).
					Blink(true)
				
				glitchedLine.WriteString(glitchStyle.Render(string(glitchChar)))
			} else {
				glitchedLine.WriteRune(r)
			}
		}
		
		glitchedLines = append(glitchedLines, glitchedLine.String())
	}
	
	return strings.Join(glitchedLines, "\n")
}

// GenerateCowsay creates enhanced cowsay-style messages
func (ac *ASCIICandy) GenerateCowsay(message, mood string) string {
	// Word wrap the message
	words := strings.Fields(message)
	var lines []string
	var currentLine strings.Builder
	maxLineLength := 40
	
	for _, word := range words {
		if currentLine.Len()+len(word)+1 > maxLineLength {
			if currentLine.Len() > 0 {
				lines = append(lines, strings.TrimSpace(currentLine.String()))
				currentLine.Reset()
			}
		}
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	if currentLine.Len() > 0 {
		lines = append(lines, strings.TrimSpace(currentLine.String()))
	}
	
	// Calculate bubble width
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	if maxLen < 10 {
		maxLen = 10
	}
	
	var cow strings.Builder
	
	// Top border
	cow.WriteString(" ")
	for i := 0; i < maxLen+2; i++ {
		cow.WriteString("_")
	}
	cow.WriteString("\n")
	
	// Message lines
	for i, line := range lines {
		if len(lines) == 1 {
			cow.WriteString(fmt.Sprintf("< %s >\n", line))
		} else if i == 0 {
			cow.WriteString(fmt.Sprintf("/ %s \\\n", line))
		} else if i == len(lines)-1 {
			cow.WriteString(fmt.Sprintf("\\ %s /\n", line))
		} else {
			cow.WriteString(fmt.Sprintf("| %s |\n", line))
		}
	}
	
	// Bottom border
	cow.WriteString(" ")
	for i := 0; i < maxLen+2; i++ {
		cow.WriteString("-")
	}
	cow.WriteString("\n")
	
	// Cow figure based on mood
	cowFigure := ac.getCowFigure(mood)
	cow.WriteString(cowFigure)
	
	// Apply styling
	cowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#40ff80")).
		Bold(true)
	
	return cowStyle.Render(cow.String())
}

// getCowFigure returns different cow figures based on mood
func (ac *ASCIICandy) getCowFigure(mood string) string {
	switch mood {
	case "happy":
		return `        \   ^__^
         \  (^^)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`
	case "sad":
		return `        \   ^__^
         \  (TT)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`
	case "cool":
		return `        \   ^__^
         \  (==)\_______
            (__)\       )\/\
             U  ||----w |
                ||     ||`
	case "ninja":
		return `        \   ^__^
         \  (@*)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`
	default: // normal
		return `        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`
	}
}

// GenerateASCIIArt creates complex ASCII art from simple patterns
func (ac *ASCIICandy) GenerateASCIIArt(artType string) string {
	switch artType {
	case "dragon":
		dragon := `                 \                /
                  \\             //
                   \\           //
                    >\         /<
                   /  \       /  \
                  |    \     /    |
                   \    \   /    /
                    \    \_/    /
                     \         /
                      |       |
                      |  _|_  |
                      | / | \ |
                      |/  |  \|
                     /   |   \
                    <    |    >
                     \   |   /
                      \  |  /
                       \_|_/`
		return ac.ApplyThemeToText(dragon, "fire")
		
	case "skull":
		skull := `                      uuuuuuu
                  uu$$$$$$$$$$$uu
               uu$$$$$$$$$$$$$$$$$uu
              u$$$$$$$$$$$$$$$$$$$$$u
             u$$$$$$$$$$$$$$$$$$$$$$$u
            u$$$$$$$$$$$$$$$$$$$$$$$$$u
            u$$$$$$$$$$$$$$$$$$$$$$$$$u
            u$$$$$$'   '$$$'   '$$$$$$u
            '$$$$'      u$u       $$$$'
             $$$u       u$u       u$$$
             $$$u      u$$$u      u$$$
              '$$$$uu$$$   $$$uu$$$$'
               '$$$$$$$'   '$$$$$$$'
                 u$$$$$$$u$$$$$$$u
                  u$'$'$'$'$'$'$u
       uuu        $$u$ $ $ $ $u$$       uuu
      u$$$$        $$$$$u$u$u$$$       u$$$$
       $$$$$uu      '$$$$$$$$$'     uu$$$$$$
     u$$$$$$$$$$$uu    '''''    uuuu$$$$$$$$$$
     $$$$'''$$$$$$$$$$uuu   uu$$$$$$$$$'''$$$'
      '''      ''$$$$$$$$$$$uu '''
                uuuu ''$$$$$$$$$$uuu
       u$$$uuu$$$$$$$$$uu ''$$$$$$$$$$$uuu$$$
       $$$$$$$$$$''''           ''$$$$$$$$$$$'
        '$$$$$'                      ''$$$$''
          $$$'                         $$$$'`
		return ac.ApplyThemeToText(skull, "purple")
		
	case "rocket":
		rocket := `                 /\
                /  \
               /    \
              /______\
             /        \
            /    /\    \
           /    /  \    \
          |    |    |    |
          |    |    |    |
          |    |    |    |
          |    |    |    |
           \   |    |   /
            \  |    |  /
             \ |____| /
              \|    |/
               |    |
               |    |
              /|    |\
             / |____| \
            /__|    |__\
               |____|`
		return ac.ApplyThemeToText(rocket, "neon")
		
	default:
		return ac.GenerateStartupBanner()
	}
}

// AnimateText creates animated text effects
func (ac *ASCIICandy) AnimateText(text string, effect string) chan string {
	frames := make(chan string, 100)
	
	go func() {
		defer close(frames)
		
		switch effect {
		case "typewriter":
			ac.typewriterEffect(text, frames)
		case "wave":
			ac.waveEffect(text, frames)
		case "pulse":
			ac.pulseEffect(text, frames)
		case "rainbow":
			ac.rainbowEffect(text, frames)
		default:
			frames <- text
		}
	}()
	
	return frames
}

// typewriterEffect simulates typing animation
func (ac *ASCIICandy) typewriterEffect(text string, frames chan string) {
	runes := []rune(text)
	for i := 0; i <= len(runes); i++ {
		frame := string(runes[:i])
		if i < len(runes) {
			frame += "█" // Cursor
		}
		frames <- frame
		time.Sleep(ac.config.AnimationSpeed)
	}
}

// waveEffect creates a wave-like animation
func (ac *ASCIICandy) waveEffect(text string, frames chan string) {
	colors := []lipgloss.Color{
		lipgloss.Color("#ff0000"),
		lipgloss.Color("#ff8000"), 
		lipgloss.Color("#ffff00"),
		lipgloss.Color("#80ff00"),
		lipgloss.Color("#00ff00"),
		lipgloss.Color("#00ff80"),
		lipgloss.Color("#00ffff"),
		lipgloss.Color("#0080ff"),
		lipgloss.Color("#0000ff"),
		lipgloss.Color("#8000ff"),
		lipgloss.Color("#ff00ff"),
		lipgloss.Color("#ff0080"),
	}
	
	runes := []rune(text)
	for frame := 0; frame < 20; frame++ {
		var styledText strings.Builder
		for i, r := range runes {
			colorIndex := (i + frame) % len(colors)
			style := lipgloss.NewStyle().Foreground(colors[colorIndex])
			styledText.WriteString(style.Render(string(r)))
		}
		frames <- styledText.String()
		time.Sleep(ac.config.AnimationSpeed)
	}
}

// pulseEffect creates a pulsing animation
func (ac *ASCIICandy) pulseEffect(text string, frames chan string) {
	for frame := 0; frame < 10; frame++ {
		intensity := float64(frame) / 10.0
		if frame > 5 {
			intensity = 1.0 - intensity
		}
		
		// Convert intensity to color brightness
		brightness := int(intensity * 255)
		color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", brightness, brightness/2, brightness))
		
		style := lipgloss.NewStyle().
			Foreground(color).
			Bold(intensity > 0.5)
		
		frames <- style.Render(text)
		time.Sleep(ac.config.AnimationSpeed)
	}
}

// rainbowEffect creates rainbow-colored text
func (ac *ASCIICandy) rainbowEffect(text string, frames chan string) {
	colors := []lipgloss.Color{
		lipgloss.Color("#ff0000"), // Red
		lipgloss.Color("#ff8000"), // Orange
		lipgloss.Color("#ffff00"), // Yellow
		lipgloss.Color("#00ff00"), // Green
		lipgloss.Color("#00ffff"), // Cyan
		lipgloss.Color("#0000ff"), // Blue
		lipgloss.Color("#8000ff"), // Purple
	}
	
	runes := []rune(text)
	for shift := 0; shift < len(colors); shift++ {
		var styledText strings.Builder
		for i, r := range runes {
			colorIndex := (i + shift) % len(colors)
			style := lipgloss.NewStyle().Foreground(colors[colorIndex])
			styledText.WriteString(style.Render(string(r)))
		}
		frames <- styledText.String()
		time.Sleep(ac.config.AnimationSpeed * 2)
	}
}

// SetTheme changes the active theme
func (ac *ASCIICandy) SetTheme(themeName string) bool {
	_, exists := ac.themes[themeName]
	return exists
}

// GetAvailableThemes returns list of available themes
func (ac *ASCIICandy) GetAvailableThemes() []string {
	var themes []string
	for name := range ac.themes {
		themes = append(themes, name)
	}
	return themes
}

// SetConfig updates ASCII candy configuration
func (ac *ASCIICandy) SetConfig(config ASCIICandyConfig) {
	ac.config = config
}

// GetConfig returns current configuration
func (ac *ASCIICandy) GetConfig() ASCIICandyConfig {
	return ac.config
}