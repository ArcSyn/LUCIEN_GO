package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ArcSyn/LucienCLI/internal/ai"
	"github.com/ArcSyn/LucienCLI/internal/completion"
	"github.com/ArcSyn/LucienCLI/internal/shell"
)


// Model represents the main TUI application state
type Model struct {
	shell              *shell.Shell
	ai                 *ai.Engine
	completion         *completion.Engine
	input              textinput.Model
	viewport           viewport.Model
	output             []string
	currentTheme       Theme
	width              int
	height             int
	ready              bool
	aiThinking         bool
	glitchEffect       bool
	showingSuggestions bool
	suggestions        []completion.Suggestion
	suggestionPage     int
	suggestionsPerPage int
	lastTabPress       time.Time
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

	// Initialize completion engine
	completionEngine := completion.New()

	model := Model{
		shell:              shell,
		ai:                 aiEngine,
		completion:         completionEngine,
		input:              ti,
		viewport:           vp,
		output:             []string{},
		currentTheme:       GetTheme("nexus"), // Load default nexus theme
		width:              80,
		height:             24,
		ready:              true, // Set ready to true so the interface shows immediately
		suggestionPage:     0,
		suggestionsPerPage: 8, // Show 8 suggestions per page
	}
	
	// Animated neural pathways loading sequence
	welcomeMsg := []string{
		"",
		model.currentTheme.SuccessStyle.Render("⚡ NEURAL PATHWAYS LOADING..."),
		model.currentTheme.InfoStyle.Render("▓▒░ [████████████████████] 100% ░▒▓"),
		"",
		model.currentTheme.SuccessStyle.Render("🔴 NEURAL INTERFACE ESTABLISHED"),
		model.currentTheme.SuccessStyle.Render("🔴 QUANTUM ENTANGLEMENT: STABLE"), 
		model.currentTheme.SuccessStyle.Render("🔴 AI SUBSYSTEMS: ONLINE"),
		model.currentTheme.SuccessStyle.Render("🔴 SECURITY PROTOCOLS: MAXIMUM"),
		model.currentTheme.SuccessStyle.Render("🔴 CYBERPUNK THEME: " + strings.ToUpper(model.currentTheme.Name)),
		"",
		model.currentTheme.InfoStyle.Render("▶ Type 'help' for command reference"),
		model.currentTheme.InfoStyle.Render("▶ Type ':theme <name>' to switch visual modes"),
		model.currentTheme.InfoStyle.Render("▶ Type ':ai <query>' for neural consultation"), 
		model.currentTheme.InfoStyle.Render("▶ Type ':hack' to enable glitch mode"),
		model.currentTheme.InfoStyle.Render("▶ Press TAB for intelligent completion"),
		"",
		model.currentTheme.SecondaryStyle.Render("🧠 NEURAL MATRIX SYNCHRONIZED - READY FOR INPUT"),
		"",
	}
	
	model.output = append(model.output, welcomeMsg...)
	model.updateViewport()

	// Initialize completion engine with shell data
	model.updateCompletionData()

	return model
}

// SetHistoryProvider sets the history provider for intelligent completion
func (m *Model) SetHistoryProvider(provider completion.HistoryProvider) {
	if m.completion != nil {
		m.completion.SetHistoryProvider(provider)
	}
}

// updateCompletionData updates the completion engine with current shell state
func (m *Model) updateCompletionData() {
	if m.shell == nil || m.completion == nil {
		return
	}
	
	// TODO: Wire up shell aliases and variables when shell API allows access
	// For now, the completion engine uses its built-in command knowledge
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.showingSuggestions = false
			_, cmdResult := m.handleCommand()
			return m, cmdResult
		case tea.KeyCtrlL:
			m.output = []string{}
			m.showingSuggestions = false
			m.updateViewport()
		case tea.KeyTab:
			_, cmdResult := m.handleTabCompletion()
			return m, cmdResult
		case tea.KeyEsc:
			m.showingSuggestions = false
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

func (m *Model) handleTabCompletion() (*Model, tea.Cmd) {
	now := time.Now()
	currentInput := m.input.Value()
	cursorPos := m.input.Position()
	
	// Check if this is a double-tab within 500ms for paging
	isDoubleTap := m.showingSuggestions && now.Sub(m.lastTabPress) < 500*time.Millisecond
	m.lastTabPress = now
	
	// Get completion suggestions if not already showing or input changed
	if !m.showingSuggestions || !isDoubleTap {
		suggestions := m.completion.Complete(currentInput, cursorPos)
		
		if len(suggestions) == 0 {
			m.showingSuggestions = false
			return m, nil
		}
		
		if len(suggestions) == 1 {
			// Auto-complete with the single suggestion
			text := suggestions[0].Text
			
			// Auto-cd for directory completions
			if suggestions[0].Type == completion.DirectoryCompletion && !strings.HasPrefix(currentInput, "cd ") {
				m.input.SetValue("cd " + text)
				m.input.SetCursor(len("cd " + text))
			} else {
				m.input.SetValue(text)
				m.input.SetCursor(len(text))
			}
			
			m.showingSuggestions = false
			m.suggestionPage = 0
		} else {
			// Show multiple suggestions
			m.suggestions = suggestions
			m.showingSuggestions = true
			m.suggestionPage = 0
			
			// Try to complete common prefix
			bestMatch := m.completion.GetBestMatch(suggestions)
			if len(bestMatch) > len(currentInput) {
				m.input.SetValue(bestMatch)
				m.input.SetCursor(len(bestMatch))
			}
		}
	} else if isDoubleTap {
		// Handle paging through suggestions
		totalPages := (len(m.suggestions) + m.suggestionsPerPage - 1) / m.suggestionsPerPage
		if totalPages > 1 {
			m.suggestionPage = (m.suggestionPage + 1) % totalPages
		}
	}
	
	return m, nil
}

func (m *Model) handleCommand() (*Model, tea.Cmd) {
	command := strings.TrimSpace(m.input.Value())
	if command == "" {
		return m, nil
	}

	// Add command to output without duplicate prompt (styled)
	cmdLine := m.stylePrimary("⚡ ") + command
	m.output = append(m.output, cmdLine)

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
			// Add command output with proper formatting
			if result.Output != "" {
				lines := strings.Split(strings.TrimRight(result.Output, "\n"), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						// Style output lines for better readability
						m.output = append(m.output, "  "+line)
					}
				}
			}
			
			// Show error output if present
			if result.Error != "" {
				errorLines := strings.Split(strings.TrimRight(result.Error, "\n"), "\n") 
				for _, line := range errorLines {
					if strings.TrimSpace(line) != "" {
						m.output = append(m.output, m.styleError("  ⚠️  "+line))
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

func (m *Model) handleSpecialCommand(command string) *Model {
	parts := strings.Fields(command[1:]) // Remove ':'
	if len(parts) == 0 {
		return m
	}

	switch parts[0] {
	case "config":
		if len(parts) >= 3 && parts[1] == "set" {
			// :config set key value
			key := parts[2]
			value := strings.Join(parts[3:], " ")
			m.handleConfigSet(key, value)
		} else if len(parts) == 2 && parts[1] == "show" {
			// :config show
			m.handleConfigShow()
		} else if len(parts) == 2 && parts[1] == "reload" {
			// :config reload
			m.handleConfigReload()
		} else {
			m.output = append(m.output, 
				m.styleError("Usage: :config set <key> <value>, :config show, or :config reload"))
		}
	
	case "theme":
		if len(parts) > 1 {
			themeName := parts[1]
			if IsValidTheme(themeName) {
				m.currentTheme = GetTheme(themeName)
				m.output = append(m.output, 
					m.styleSuccess(fmt.Sprintf("🎨 Theme switched to: %s", m.currentTheme.Name)))
			} else {
				available := GetThemeNames()
				m.output = append(m.output, 
					m.styleError(fmt.Sprintf("❌ Unknown theme. Available: %s", 
					strings.Join(available, ", "))))
			}
		} else {
			// Show current theme and available options
			m.output = append(m.output, m.styleInfo(fmt.Sprintf("Current theme: %s", m.currentTheme.Name)))
			m.output = append(m.output, m.styleSecondary(fmt.Sprintf("Available themes: %s", strings.Join(GetThemeNames(), ", "))))
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

	case "spells":
		m.showSpells()

	case "weather":
		m.showWeather()

	case "help":
		m.showHelp()
	}

	return m
}

func (m *Model) showHelp() {
	helpText := []string{
		"",
		m.stylePrimary("🔴 LUCIEN NEURAL INTERFACE - COMMAND REFERENCE"),
		"",
		"📟 SYSTEM COMMANDS:",
		"  :theme <name>     Switch visual theme (nexus, synthwave, ghost)",
		"  :ai <query>       Consult neural network",
		"  :spells           List all available AI agents",
		"  :weather          Show weather information [WIP]",
		"  :hack             Toggle glitch mode",
		"  :clear            Clear terminal buffer",
		"  :help             Show this reference",
		"",
		"⚡ SHELL OPERATIONS:",
		"  Standard shell commands with pipes, redirects, and variables",
		"  Built-in commands: cd, set, alias, exit",
		"",
		"🤖 AI AGENT COMMANDS:",
		"  plan \"task\"       Break down goals into actionable tasks",
		"  design \"idea\"     Generate UI code from descriptions",
		"  review file.py   Analyze code and suggest improvements",
		"  code \"request\"   Generate, refactor, or explain code",
		"",
		"🧠 AI FEATURES:",
		"  • Predictive command suggestions",
		"  • Context-aware assistance", 
		"  • Neural pattern recognition",
		"  • Intelligent tab completion",
		"",
		"🛡️  SECURITY:",
		"  • OPA policy enforcement",
		"  • Sandboxed plugin execution",
		"  • Safe-mode command filtering",
		"",
	}

	m.output = append(m.output, helpText...)
}

func (m *Model) showSpells() {
	m.output = append(m.output, "")
	m.output = append(m.output, m.stylePrimary("✨ AVAILABLE AI AGENTS (SPELLS)"))
	m.output = append(m.output, m.stylePrimary("================================"))
	m.output = append(m.output, "")
	
	// List the AI agents directly
	agents := []struct {
		name string
		desc string
		example string
	}{
		{"plan", "Break down goals into actionable tasks", `plan "create a web app"`},
		{"design", "Generate UI code from descriptions", `design "login form with styling"`},
		{"review", "Analyze code and suggest improvements", `review main.py`},
		{"code", "Generate, refactor, or explain code", `code "write a fibonacci function"`},
	}
	
	for _, agent := range agents {
		m.output = append(m.output, m.styleSuccess(fmt.Sprintf("🤖 %s", agent.name)))
		m.output = append(m.output, fmt.Sprintf("   %s", agent.desc))
		m.output = append(m.output, m.styleSecondary(fmt.Sprintf("   Example: %s", agent.example)))
		m.output = append(m.output, "")
	}
	
	m.output = append(m.output, m.stylePrimary("💡 Neural Network Integration:"))
	m.output = append(m.output, "   • AI agents process requests through secure Python bridge")
	m.output = append(m.output, "   • Each agent specializes in specific cognitive domains")
	m.output = append(m.output, "   • Responses are styled with cyberpunk aesthetics")
	m.output = append(m.output, "")
}

func (m *Model) showWeather() {
	m.output = append(m.output, "")
	m.output = append(m.output, m.stylePrimary("🌤️  NEURAL WEATHER SYSTEM"))
	m.output = append(m.output, m.stylePrimary("======================="))
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleWarning("⚠️  [WIP] Weather widget integration in progress..."))
	m.output = append(m.output, "")
	m.output = append(m.output, "🔮 Future features:")
	m.output = append(m.output, "   • Real-time weather data")
	m.output = append(m.output, "   • Location-based forecasts")
	m.output = append(m.output, "   • Cyberpunk weather visualization")
	m.output = append(m.output, "   • Integration with system notifications")
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleSecondary("   Run ':theme synthwave' for optimal atmospheric conditions"))
	m.output = append(m.output, "")
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
	
	// Add completion suggestions if showing
	var suggestionSection string
	if m.showingSuggestions && len(m.suggestions) > 0 {
		suggestionSection = m.renderSuggestions()
	}
	
	footer := m.renderFooter()

	viewComponents := []string{header, content, inputSection}
	if suggestionSection != "" {
		viewComponents = append(viewComponents, suggestionSection)
	}
	viewComponents = append(viewComponents, footer)

	view := lipgloss.JoinVertical(lipgloss.Left, viewComponents...)

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

	left := m.currentTheme.HeaderStyle.Render(title)
	right := m.currentTheme.HeaderStyle.Copy().Render(aiStatus)

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, left) + 
		   lipgloss.PlaceHorizontal(m.width, lipgloss.Right, right)
}

func (m Model) renderInput() string {
	prompt := m.formatPrompt()
	return m.currentTheme.InputStyle.Render(prompt + m.input.View())
}

func (m Model) renderSuggestions() string {
	if len(m.suggestions) == 0 {
		return ""
	}
	
	// Calculate pagination
	totalSuggestions := len(m.suggestions)
	totalPages := (totalSuggestions + m.suggestionsPerPage - 1) / m.suggestionsPerPage
	currentPage := m.suggestionPage
	
	if currentPage >= totalPages {
		currentPage = 0
	}
	
	startIdx := currentPage * m.suggestionsPerPage
	endIdx := startIdx + m.suggestionsPerPage
	if endIdx > totalSuggestions {
		endIdx = totalSuggestions
	}
	
	suggestions := m.suggestions[startIdx:endIdx]
	
	var suggestionLines []string
	
	// Header with page info
	if totalPages > 1 {
		header := fmt.Sprintf("🔧 TAB COMPLETION SUGGESTIONS (Page %d/%d):", currentPage+1, totalPages)
		suggestionLines = append(suggestionLines, m.styleSuccess(header))
	} else {
		suggestionLines = append(suggestionLines, m.styleSuccess("🔧 TAB COMPLETION SUGGESTIONS:"))
	}
	
	for _, suggestion := range suggestions {
		var icon string
		switch suggestion.Type {
		case completion.CommandCompletion:
			icon = "⚡"
		case completion.FileCompletion:
			icon = "📄"
		case completion.DirectoryCompletion:
			icon = "📁"
		case completion.VariableCompletion:
			icon = "💲"
		case completion.AliasCompletion:
			icon = "🔗"
		case completion.HistoryCompletion:
			icon = "📜"
		default:
			icon = "▶"
		}
		
		suggestionText := fmt.Sprintf("  %s %s", icon, suggestion.Text)
		if suggestion.Description != "" {
			suggestionText += m.styleSecondary(fmt.Sprintf(" (%s)", suggestion.Description))
		}
		
		suggestionLines = append(suggestionLines, suggestionText)
	}
	
	// Footer with navigation hint
	if totalPages > 1 {
		footer := fmt.Sprintf("  Press TAB again for next page • ESC to dismiss")
		suggestionLines = append(suggestionLines, m.styleSecondary(footer))
	} else if totalSuggestions > 0 {
		footer := "  Press ESC to dismiss"
		suggestionLines = append(suggestionLines, m.styleSecondary(footer))
	}
	
	return strings.Join(suggestionLines, "\n")
}

func (m Model) renderFooter() string {
	shortcuts := "CTRL+C:quit • CTRL+L:clear • TAB:complete • :help for commands"
	return m.currentTheme.FooterStyle.Render(shortcuts)
}

func (m Model) renderGlitchEffect() string {
	// Create cyberpunk glitch overlay
	glitchChars := "▓▒░█▄▀■□▪▫"
	glitchStyle := m.currentTheme.ErrorStyle.Copy().Blink(true)

	return glitchStyle.Render(string(glitchChars[time.Now().Unix()%int64(len(glitchChars))]))
}

func (m Model) formatPrompt() string {
	return m.currentTheme.PromptStyle.Render("lucien@nexus:~$ ")
}

func (m *Model) updateViewport() {
	content := strings.Join(m.output, "\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

// Styling helpers using the production theme system
func (m Model) stylePrimary(text string) string {
	return m.currentTheme.PromptStyle.Render(text)
}

func (m Model) styleSuccess(text string) string {
	return m.currentTheme.SuccessStyle.Render(text)
}

func (m Model) styleError(text string) string {
	return m.currentTheme.ErrorStyle.Render(text)
}

func (m Model) styleWarning(text string) string {
	return m.currentTheme.WarningStyle.Render(text)
}

func (m Model) styleSecondary(text string) string {
	return m.currentTheme.SecondaryStyle.Render(text)
}

func (m Model) styleInfo(text string) string {
	return m.currentTheme.InfoStyle.Render(text)
}

func (m Model) styleCommand(text string) string {
	return m.currentTheme.CommandStyle.Render(text)
}

func (m Model) styleOutput(text string) string {
	return m.currentTheme.OutputStyle.Render(text)
}

// Config handling methods
func (m *Model) handleConfigSet(key, value string) {
	config, err := LoadConfig()
	if err != nil {
		m.output = append(m.output, m.styleError(fmt.Sprintf("❌ Failed to load config: %v", err)))
		return
	}
	
	if err := SetConfigValue(config, key, value); err != nil {
		m.output = append(m.output, m.styleError(fmt.Sprintf("❌ Failed to set config: %v", err)))
		return
	}
	
	if err := SaveConfig(config); err != nil {
		m.output = append(m.output, m.styleError(fmt.Sprintf("❌ Failed to save config: %v", err)))
		return
	}
	
	m.output = append(m.output, m.styleSuccess(fmt.Sprintf("✅ Config set: %s = %s", key, value)))
	m.output = append(m.output, m.styleInfo("💡 Restart shell for changes to take effect"))
}

func (m *Model) handleConfigShow() {
	config, err := LoadConfig()
	if err != nil {
		m.output = append(m.output, m.styleError(fmt.Sprintf("❌ Failed to load config: %v", err)))
		return
	}
	
	configPath, _ := GetConfigPath()
	m.output = append(m.output, m.styleSuccess("📝 CURRENT CONFIGURATION"))
	m.output = append(m.output, m.styleSecondary(fmt.Sprintf("Config file: %s", configPath)))
	m.output = append(m.output, "")
	
	m.output = append(m.output, m.styleInfo("🐚 SHELL"))
	m.output = append(m.output, fmt.Sprintf("  prompt = %s", config.Shell.Prompt))
	m.output = append(m.output, fmt.Sprintf("  safe_mode = %t", config.Shell.SafeMode))
	m.output = append(m.output, fmt.Sprintf("  default_theme = %s", config.Shell.DefaultTheme))
	m.output = append(m.output, fmt.Sprintf("  execution_timeout = %d", config.Shell.ExecutionTimeout))
	
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleInfo("🎨 UI"))
	m.output = append(m.output, fmt.Sprintf("  animated_startup = %t", config.UI.AnimatedStartup))
	m.output = append(m.output, fmt.Sprintf("  glitch_effects = %t", config.UI.GlitchEffects))
	m.output = append(m.output, fmt.Sprintf("  color_support = %s", config.UI.ColorSupport))
	
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleInfo("📜 HISTORY"))
	m.output = append(m.output, fmt.Sprintf("  enabled = %t", config.History.Enabled))
	m.output = append(m.output, fmt.Sprintf("  max_entries = %d", config.History.MaxEntries))
	m.output = append(m.output, fmt.Sprintf("  save_on_exit = %t", config.History.SaveOnExit))
	
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleInfo("🔧 COMPLETION"))
	m.output = append(m.output, fmt.Sprintf("  enabled = %t", config.Completion.Enabled))
	m.output = append(m.output, fmt.Sprintf("  suggestions_per_page = %d", config.Completion.SuggestionsPerPage))
	m.output = append(m.output, fmt.Sprintf("  auto_cd = %t", config.Completion.AutoCD))
	m.output = append(m.output, fmt.Sprintf("  fuzzy_matching = %t", config.Completion.FuzzyMatching))
	
	m.output = append(m.output, "")
	m.output = append(m.output, m.styleInfo("🧠 AI"))
	m.output = append(m.output, fmt.Sprintf("  enabled = %t", config.AI.Enabled))
	m.output = append(m.output, fmt.Sprintf("  suggest_commands = %t", config.AI.SuggestCommands))
	m.output = append(m.output, fmt.Sprintf("  confidence_threshold = %.2f", config.AI.ConfidenceThreshold))
	m.output = append(m.output, fmt.Sprintf("  model_provider = %s", config.AI.ModelProvider))
}

func (m *Model) handleConfigReload() {
	config, err := LoadConfig()
	if err != nil {
		m.output = append(m.output, m.styleError(fmt.Sprintf("❌ Failed to reload config: %v", err)))
		return
	}
	
	// Apply relevant config changes that can be applied immediately
	if config.Completion.SuggestionsPerPage > 0 {
		m.suggestionsPerPage = config.Completion.SuggestionsPerPage
	}
	
	// Switch theme if different
	if IsValidTheme(config.Shell.DefaultTheme) && config.Shell.DefaultTheme != m.currentTheme.Name {
		m.currentTheme = GetTheme(config.Shell.DefaultTheme)
	}
	
	m.output = append(m.output, m.styleSuccess("✅ Configuration reloaded"))
	m.output = append(m.output, m.styleInfo("💡 Some changes require shell restart"))
}