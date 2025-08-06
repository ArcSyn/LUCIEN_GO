package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luciendev/lucien-core/internal/ai"
	"github.com/luciendev/lucien-core/internal/shell"
)

// Theme configuration for that retro hacker aesthetic
type Theme struct {
	Name      string
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Error     lipgloss.Color
	Warning   lipgloss.Color
	Background lipgloss.Color
	Text       lipgloss.Color
}

var themes = map[string]Theme{
	"nexus": {
		Name:       "nexus",
		Primary:    lipgloss.Color("#00ff41"), // Matrix green
		Secondary:  lipgloss.Color("#0066cc"), // Cyber blue
		Success:    lipgloss.Color("#00ff00"), // Bright green
		Error:      lipgloss.Color("#ff0066"), // Neon pink
		Warning:    lipgloss.Color("#ffaa00"), // Amber
		Background: lipgloss.Color("#000000"), // Pure black
		Text:       lipgloss.Color("#ffffff"), // White
	},
	"synthwave": {
		Name:       "synthwave",
		Primary:    lipgloss.Color("#ff00ff"), // Hot pink
		Secondary:  lipgloss.Color("#00ffff"), // Cyan
		Success:    lipgloss.Color("#39ff14"), // Electric lime
		Error:      lipgloss.Color("#ff073a"), // Red
		Warning:    lipgloss.Color("#ffb000"), // Orange
		Background: lipgloss.Color("#1a0033"), // Dark purple
		Text:       lipgloss.Color("#ffffff"), // White
	},
	"ghost": {
		Name:       "ghost",
		Primary:    lipgloss.Color("#ffffff"), // White
		Secondary:  lipgloss.Color("#aaaaaa"), // Gray
		Success:    lipgloss.Color("#00cc00"), // Green
		Error:      lipgloss.Color("#cc0000"), // Red
		Warning:    lipgloss.Color("#ccaa00"), // Yellow
		Background: lipgloss.Color("#000000"), // Black
		Text:       lipgloss.Color("#cccccc"), // Light gray
	},
}

// Model represents the main TUI application state
type Model struct {
	shell        *shell.Shell
	ai           *ai.Engine
	input        textinput.Model
	viewport     viewport.Model
	output       []string
	theme        Theme
	width        int
	height       int
	ready        bool
	aiThinking   bool
	glitchEffect bool
}

// Mind-blowing feature 1: AI predictive suggestions with neural network visualization
type aiSuggestion struct {
	command    string
	confidence float64
	reasoning  string
}

// Mind-blowing feature 2: Glitch effect for system alerts/hacks
type glitchMsg struct{}

func (m Model) glitchTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return glitchMsg{}
	})
}

// NewModel creates a new UI model with cyberpunk aesthetics
func NewModel(shell *shell.Shell, aiEngine *ai.Engine) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter command... Neural interface ready █"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	vp := viewport.New(80, 20)
	vp.SetContent("")

	model := Model{
		shell:    shell,
		ai:       aiEngine,
		input:    ti,
		viewport: vp,
		output:   []string{},
		theme:    themes["nexus"], // Default to Matrix theme
	}

	// Welcome message with hacker aesthetic
	welcomeMsg := []string{
		"",
		"🔴 NEURAL INTERFACE ESTABLISHED",
		"🔴 QUANTUM ENTANGLEMENT: STABLE", 
		"🔴 AI SUBSYSTEMS: ONLINE",
		"🔴 SECURITY PROTOCOLS: MAXIMUM",
		"",
		"▶ Type 'help' for command reference",
		"▶ Type ':theme <name>' to switch visual modes",
		"▶ Type ':ai <query>' for neural consultation", 
		"▶ Type ':hack' to enable glitch mode",
		"",
	}
	
	model.output = append(model.output, welcomeMsg...)
	model.updateViewport()

	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8
		m.input.Width = msg.Width - 20
		m.ready = true

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleCommand()
		case tea.KeyCtrlL:
			m.output = []string{}
			m.updateViewport()
		}

	case glitchMsg:
		if m.glitchEffect {
			return m, m.glitchTick()
		}
	}

	// Update input
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleCommand() (Model, tea.Cmd) {
	command := strings.TrimSpace(m.input.Value())
	if command == "" {
		return m, nil
	}

	// Add command to output with cyberpunk prompt
	prompt := m.formatPrompt() + command
	m.output = append(m.output, prompt)

	// Handle special commands
	if strings.HasPrefix(command, ":") {
		m = m.handleSpecialCommand(command)
	} else {
		// Execute through shell
		result, err := m.shell.Execute(command)
		if err != nil {
			errorMsg := m.styleError(fmt.Sprintf("❌ ERROR: %v", err))
			m.output = append(m.output, errorMsg)
		} else {
			// Add command output
			if result.Output != "" {
				lines := strings.Split(result.Output, "\n")
				for _, line := range lines {
					if line != "" {
						m.output = append(m.output, "  "+line)
					}
				}
			}
			
			// Show AI suggestions if available
			if suggestions := m.getAISuggestions(command); len(suggestions) > 0 {
				m.output = append(m.output, "")
				m.output = append(m.output, m.styleSuccess("🧠 NEURAL SUGGESTIONS:"))
				for _, suggestion := range suggestions {
					confidence := fmt.Sprintf("%.0f%%", suggestion.confidence*100)
					suggestionLine := fmt.Sprintf("  ▶ %s [%s confidence]", 
						suggestion.command, confidence)
					m.output = append(m.output, m.stylePrimary(suggestionLine))
				}
			}
		}
	}

	m.input.SetValue("")
	m.updateViewport()
	
	return m, nil
}

func (m Model) handleSpecialCommand(command string) Model {
	parts := strings.Fields(command[1:]) // Remove ':'
	if len(parts) == 0 {
		return m
	}

	switch parts[0] {
	case "theme":
		if len(parts) > 1 {
			if theme, exists := themes[parts[1]]; exists {
				m.theme = theme
				m.output = append(m.output, 
					m.styleSuccess(fmt.Sprintf("🎨 Theme switched to: %s", theme.Name)))
			} else {
				available := make([]string, 0, len(themes))
				for name := range themes {
					available = append(available, name)
				}
				m.output = append(m.output, 
					m.styleError(fmt.Sprintf("❌ Unknown theme. Available: %s", 
					strings.Join(available, ", "))))
			}
		}

	case "ai":
		if len(parts) > 1 {
			query := strings.Join(parts[1:], " ")
			m.aiThinking = true
			response, err := m.ai.Query(query)
			m.aiThinking = false
			
			if err != nil {
				m.output = append(m.output, m.styleError(fmt.Sprintf("🧠 AI ERROR: %v", err)))
			} else {
				m.output = append(m.output, m.styleSuccess("🧠 AI RESPONSE:"))
				lines := strings.Split(response, "\n")
				for _, line := range lines {
					m.output = append(m.output, "  "+line)
				}
			}
		}

	case "hack":
		m.glitchEffect = !m.glitchEffect
		if m.glitchEffect {
			m.output = append(m.output, m.styleError("🔥 GLITCH MODE ACTIVATED"))
			m.output = append(m.output, m.styleError("▓▒░ REALITY.EXE HAS STOPPED WORKING ░▒▓"))
		} else {
			m.output = append(m.output, m.styleSuccess("✅ SYSTEMS STABILIZED"))
		}

	case "clear":
		m.output = []string{}

	case "help":
		m.showHelp()
	}

	return m
}

func (m Model) showHelp() {
	helpText := []string{
		"",
		m.stylePrimary("🔴 LUCIEN NEURAL INTERFACE - COMMAND REFERENCE"),
		"",
		"📟 SYSTEM COMMANDS:",
		"  :theme <name>     Switch visual theme (nexus, synthwave, ghost)",
		"  :ai <query>       Consult neural network",
		"  :hack             Toggle glitch mode",
		"  :clear            Clear terminal buffer",
		"  :help             Show this reference",
		"",
		"⚡ SHELL OPERATIONS:",
		"  Standard shell commands with pipes, redirects, and variables",
		"  Built-in commands: cd, set, alias, exit",
		"",
		"🧠 AI FEATURES:",
		"  • Predictive command suggestions",
		"  • Context-aware assistance", 
		"  • Neural pattern recognition",
		"",
		"🛡️  SECURITY:",
		"  • OPA policy enforcement",
		"  • Sandboxed plugin execution",
		"  • Safe-mode command filtering",
		"",
	}

	m.output = append(m.output, helpText...)
}

func (m Model) getAISuggestions(command string) []aiSuggestion {
	// This would integrate with the actual AI engine
	// For now, return some smart suggestions based on patterns
	suggestions := []aiSuggestion{}

	if strings.Contains(command, "git") && strings.Contains(command, "status") {
		suggestions = append(suggestions, aiSuggestion{
			command:    "git add .",
			confidence: 0.89,
			reasoning:  "Common workflow after git status",
		})
	}

	if strings.Contains(command, "ls") {
		suggestions = append(suggestions, aiSuggestion{
			command:    "cd <directory>",
			confidence: 0.76,
			reasoning:  "Navigation after listing",
		})
	}

	return suggestions
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing neural pathways...\n"
	}

	// Apply glitch effect if enabled
	var glitchOverlay string
	if m.glitchEffect {
		glitchOverlay = m.renderGlitchEffect()
	}

	// Main view components
	header := m.renderHeader()
	content := m.viewport.View()
	inputSection := m.renderInput()
	footer := m.renderFooter()

	view := lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		inputSection,
		footer,
	)

	if glitchOverlay != "" {
		// Layer glitch effect over the main view
		return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, view) + glitchOverlay
	}

	return view
}

func (m Model) renderHeader() string {
	title := "LUCIEN NEXUS-7 TERMINAL"
	aiStatus := "🧠 AI:READY"
	if m.aiThinking {
		aiStatus = "🧠 AI:THINKING..."
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Background(m.theme.Background).
		Bold(true).
		Padding(0, 2)

	left := headerStyle.Render(title)
	right := headerStyle.Copy().Foreground(m.theme.Secondary).Render(aiStatus)

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, left) + 
		   lipgloss.PlaceHorizontal(m.width, lipgloss.Right, right)
}

func (m Model) renderInput() string {
	prompt := m.formatPrompt()
	inputStyle := lipgloss.NewStyle().
		Foreground(m.theme.Text).
		Background(m.theme.Background)

	return inputStyle.Render(prompt + m.input.View())
}

func (m Model) renderFooter() string {
	shortcuts := "CTRL+C:quit • CTRL+L:clear • :help for commands"
	footerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Secondary).
		Background(m.theme.Background).
		Italic(true)

	return footerStyle.Render(shortcuts)
}

func (m Model) renderGlitchEffect() string {
	// Create cyberpunk glitch overlay
	glitchChars := "▓▒░█▄▀■□▪▫"
	glitchStyle := lipgloss.NewStyle().
		Foreground(m.theme.Error).
		Background(m.theme.Background).
		Blink(true)

	return glitchStyle.Render(string(glitchChars[time.Now().Unix()%int64(len(glitchChars))]))
}

func (m Model) formatPrompt() string {
	promptStyle := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true)

	return promptStyle.Render("lucien@nexus:~$ ")
}

func (m Model) updateViewport() {
	content := strings.Join(m.output, "\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

// Styling helpers for that cyberpunk aesthetic
func (m Model) stylePrimary(text string) string {
	return lipgloss.NewStyle().Foreground(m.theme.Primary).Render(text)
}

func (m Model) styleSuccess(text string) string {
	return lipgloss.NewStyle().Foreground(m.theme.Success).Bold(true).Render(text)
}

func (m Model) styleError(text string) string {
	return lipgloss.NewStyle().Foreground(m.theme.Error).Bold(true).Render(text)
}

func (m Model) styleWarning(text string) string {
	return lipgloss.NewStyle().Foreground(m.theme.Warning).Render(text)
}