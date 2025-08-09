package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/glamour"
	
	"github.com/ArcSyn/LucienCLI/internal/ai"
	"github.com/ArcSyn/LucienCLI/internal/completion"
	"github.com/ArcSyn/LucienCLI/internal/shell"
)

// ViewMode represents different UI view modes
type ViewMode int

const (
	TerminalView ViewMode = iota
	SystemView
	AIView
	NetworkView
	SettingsView
)

// EnhancedModel represents the advanced Visual Bliss TUI with tabs, panes, and animations
type EnhancedModel struct {
	// Core components
	shell              *shell.Shell
	ai                 *ai.Engine
	completion         *completion.Engine
	currentTheme       Theme
	
	// Layout and dimensions
	width              int
	height             int
	ready              bool
	
	// Tab system
	tabNames           []string
	activeTab          int
	tabContent         []string
	
	// Terminal pane
	input              CustomInput
	viewport           viewport.Model
	output             []string
	suggestions        []completion.Suggestion
	showingSuggestions bool
	suggestionPage     int
	suggestionsPerPage int
	lastTabPress       time.Time
	
	// System monitoring pane
	systemStats        SystemStats
	systemSpinner      spinner.Model
	
	// AI pane with Glamour markdown rendering
	aiHistory          []AIInteraction
	aiInput            CustomInput
	aiViewport         viewport.Model
	aiThinking         bool
	glamourRenderer    *glamour.TermRenderer
	
	// Network monitoring pane
	networkStats       NetworkStats
	networkPaginator   paginator.Model
	
	// Settings pane
	settingsViewport   viewport.Model
	settingsInput      CustomInput
	editingSettings    bool
	
	// Visual effects
	glitchEffect       bool
	animations         map[string]Animation
	pulseFrame         int
	
	// Presence system
	presenceState      PresenceState
	activityLevel      float64
	lastActivity       time.Time
	
	// Key bindings
	keyMap             KeyMap
	
	// Sound system integration
	soundSystem        *SoundSystem
	webOverlay         *WebOverlay
	
	// ASCII candy system
	asciiCandy         *ASCIICandy
	
	// Special behaviors manager
	specialBehaviors   *SpecialBehaviorManager
}

// SystemStats holds real system monitoring data
type SystemStats struct {
	CPUUsage     float64
	MemoryUsage  float64
	DiskUsage    float64
	ActiveJobs   int
	Uptime       time.Duration
	LastUpdate   time.Time
}

// NetworkStats holds network monitoring data
type NetworkStats struct {
	Connections []NetworkConnection
	Bandwidth   BandwidthStats
	LastUpdate  time.Time
}

type NetworkConnection struct {
	LocalAddr  string
	RemoteAddr string
	State      string
	Process    string
}

type BandwidthStats struct {
	BytesIn  uint64
	BytesOut uint64
	PacketsIn uint64
	PacketsOut uint64
}

// AIInteraction represents a conversation with the AI
type AIInteraction struct {
	Query     string
	Response  string
	Timestamp time.Time
	Markdown  bool
}

// Animation represents visual animations
type Animation struct {
	Name      string
	Frame     int
	MaxFrames int
	Duration  time.Duration
	LastTick  time.Time
}

// PresenceState represents user presence and activity
type PresenceState struct {
	Mode         string  // "focus", "idle", "away"
	GoalProgress float64 // 0.0 to 1.0
	LastCommand  string
	SessionTime  time.Duration
}

// KeyMap defines the application key bindings
type KeyMap struct {
	NextTab    key.Binding
	PrevTab    key.Binding
	Quit       key.Binding
	Clear      key.Binding
	Help       key.Binding
	Settings   key.Binding
	AI         key.Binding
	Network    key.Binding
	System     key.Binding
}

// NewEnhancedModel creates a new enhanced Visual Bliss model
func NewEnhancedModel(shell *shell.Shell, aiEngine *ai.Engine) *EnhancedModel {
	// Initialize tabs
	tabNames := []string{
		"🖥️  Terminal",
		"📊 System", 
		"🧠 AI Chat",
		"🌐 Network",
		"⚙️  Settings",
	}
	
	// Initialize terminal components
	ti := NewCustomInput()
	ti.SetPlaceholder("Neural interface ready... Enter command █")
	ti.Focus()
	ti.SetMaxLength(512)
	ti.SetWidth(80)
	
	vp := viewport.New(100, 30)
	vp.SetContent("")
	
	// Initialize AI components
	aiInput := NewCustomInput()
	aiInput.SetPlaceholder("Ask the AI anything...")
	aiInput.SetMaxLength(512)
	aiInput.SetWidth(80)
	
	aiVP := viewport.New(100, 30)
	aiVP.SetContent("")
	
	// Initialize Glamour renderer for markdown
	glamourRenderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	
	// Initialize system spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff41"))
	
	// Initialize paginator for network view
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff41")).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render("•")
	
	// Initialize settings viewport
	settingsVP := viewport.New(100, 30)
	settingsVP.SetContent("")
	
	settingsInput := NewCustomInput()
	settingsInput.SetPlaceholder("config key=value")
	settingsInput.SetMaxLength(256)
	settingsInput.SetWidth(60)
	
	// Initialize key map
	keyMap := KeyMap{
		NextTab:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
		PrevTab:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
		Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		Clear:      key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "clear")),
		Help:       key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "help")),
		Settings:   key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "settings")),
		AI:         key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "ai chat")),
		Network:    key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "network")),
		System:     key.NewBinding(key.WithKeys("ctrl+y"), key.WithHelp("ctrl+y", "system")),
	}
	
	model := &EnhancedModel{
		shell:              shell,
		ai:                 aiEngine,
		completion:         completion.New(),
		currentTheme:       GetCatppuccinTheme("mocha"), // Default to Catppuccin
		width:              120,
		height:             40,
		ready:              true,
		tabNames:           tabNames,
		activeTab:          0,
		tabContent:         make([]string, len(tabNames)),
		input:              ti,
		viewport:           vp,
		output:             []string{},
		suggestions:        []completion.Suggestion{},
		showingSuggestions: false,
		suggestionsPerPage: 8,
		aiHistory:          []AIInteraction{},
		aiInput:            aiInput,
		aiViewport:         aiVP,
		aiThinking:         false,
		glamourRenderer:    glamourRenderer,
		systemSpinner:      s,
		networkPaginator:   p,
		settingsViewport:   settingsVP,
		settingsInput:      settingsInput,
		editingSettings:    false,
		animations:         make(map[string]Animation),
		pulseFrame:         0,
		presenceState: PresenceState{
			Mode:         "focus",
			GoalProgress: 0.0,
			SessionTime:  0,
		},
		activityLevel: 1.0,
		lastActivity:  time.Now(),
		keyMap:        keyMap,
		soundSystem:   NewSoundSystem(),
		webOverlay:    NewWebOverlay(shell, 8080),
		asciiCandy:    NewASCIICandy(),
	}
	
	// Initialize special behaviors (after sound system and ASCII candy)
	model.specialBehaviors = NewSpecialBehaviorManager(model.soundSystem, model.asciiCandy)
	
	// Initialize animations
	model.initAnimations()
	
	// Show enhanced startup sequence
	model.showEnhancedStartup()
	
	// Initialize completion engine
	model.updateCompletionData()
	
	// Set up shell message dispatcher
	model.shell.SetDispatcher(nil)
	
	// Ensure the terminal input is focused by default
	model.focusActiveTabInput()
	
	return model
}

// initAnimations sets up visual animations
func (m *EnhancedModel) initAnimations() {
	m.animations["pulse"] = Animation{
		Name:      "pulse",
		Frame:     0,
		MaxFrames: 30,
		Duration:  time.Millisecond * 100,
		LastTick:  time.Now(),
	}
	
	m.animations["matrix"] = Animation{
		Name:      "matrix",
		Frame:     0,
		MaxFrames: 60,
		Duration:  time.Millisecond * 150,
		LastTick:  time.Now(),
	}
}

// showEnhancedStartup displays the Visual Bliss startup sequence
func (m *EnhancedModel) showEnhancedStartup() {
	// Generate enhanced startup banner using ASCII candy
	banner := m.asciiCandy.GenerateStartupBanner()
	
	startupSequence := []string{
		"",
		banner, // Use the ASCII candy generated banner
		"",
		m.currentTheme.SuccessStyle.Render("⚡ INITIALIZING VISUAL BLISS SYSTEMS..."),
		"",
		m.currentTheme.SuccessStyle.Render("🎨 Layer 1 - Bubbletea TUI Engine: ") + m.currentTheme.SuccessStyle.Render("ACTIVE"),
		m.currentTheme.SuccessStyle.Render("✨ Layer 2 - Lip Gloss + Glamour Styling: ") + m.currentTheme.SuccessStyle.Render("LOADED"),
		m.currentTheme.SuccessStyle.Render("📝 Layer 3 - Markdown & Syntax Rendering: ") + m.currentTheme.SuccessStyle.Render("READY"),
		m.currentTheme.SuccessStyle.Render("🌐 Layer 4 - WebSocket Terminal Bridge: ") + m.currentTheme.InfoStyle.Render("STANDBY"),
		m.currentTheme.SuccessStyle.Render("🎭 Layer 5 - Catppuccin Theme System: ") + m.currentTheme.SuccessStyle.Render("ENABLED"),
		m.currentTheme.SuccessStyle.Render("🔊 Layer 6 - Sound Feedback: ") + m.currentTheme.InfoStyle.Render("READY"),
		m.currentTheme.SuccessStyle.Render("🎪 Layer 7 - ASCII & Visual Effects: ") + m.currentTheme.SuccessStyle.Render("LOADED"),
		"",
		m.currentTheme.InfoStyle.Render("🧠 AI SUBSYSTEMS: Neural pathways synchronized"),
		m.currentTheme.InfoStyle.Render("🛡️  SECURITY: OPA policies active"),
		m.currentTheme.InfoStyle.Render("📊 MONITORING: System vitals tracking"),
		m.currentTheme.InfoStyle.Render("🌐 NETWORK: Connection monitoring active"),
		"",
		m.currentTheme.SuccessStyle.Render("▶ TAB/SHIFT+TAB: Navigate between panes"),
		m.currentTheme.SuccessStyle.Render("▶ CTRL+A: AI Chat | CTRL+Y: System | CTRL+N: Network"),
		m.currentTheme.SuccessStyle.Render("▶ CTRL+H: Help | CTRL+S: Settings | CTRL+C: Quit"),
		"",
		m.currentTheme.SecondaryStyle.Render("🚀 VISUAL BLISS INTERFACE - READY FOR NEURAL INPUT"),
		"",
	}
	
	m.output = append(m.output, startupSequence...)
	m.updateViewport()
}

// updateCompletionData updates the completion engine with current shell state
func (m *EnhancedModel) updateCompletionData() {
	if m.shell == nil || m.completion == nil {
		return
	}
	// Wire up completion data as needed
}

func (m *EnhancedModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.systemSpinner.Tick,
		m.tickAnimations(),
		m.tickPresence(),
	)
}

// tickAnimations creates commands to update animations
func (m *EnhancedModel) tickAnimations() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return AnimationTickMsg{Time: t}
	})
}

// tickPresence creates commands to update presence system
func (m *EnhancedModel) tickPresence() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return PresenceTickMsg{Time: t}
	})
}

// Message types for the enhanced model
type AnimationTickMsg struct{ Time time.Time }
type PresenceTickMsg struct{ Time time.Time }
type SystemStatsMsg struct{ Stats SystemStats }
type NetworkStatsMsg struct{ Stats NetworkStats }

func (m *EnhancedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeComponents()
		m.ready = true
		
	case AnimationTickMsg:
		m.updateAnimations(msg.Time)
		cmds = append(cmds, m.tickAnimations())
		
	case PresenceTickMsg:
		m.updatePresence(msg.Time)
		cmds = append(cmds, m.tickPresence())
		
	case SystemStatsMsg:
		m.systemStats = msg.Stats
		
	case NetworkStatsMsg:
		m.networkStats = msg.Stats
		
	// Special behavior messages
	case AwakenSequenceMsg:
		return m.HandleAwakenSequence(msg)
		
	case ProphecyMsg:
		return m.HandleProphecy(msg)
		
	case ListenTickMsg:
		return m.HandleListenTick(msg)
		
	case VanishSequenceMsg:
		return m.HandleVanishSequence(msg)
		
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.systemSpinner, cmd = m.systemSpinner.Update(msg)
		cmds = append(cmds, cmd)
		
	case tea.KeyMsg:
		// Update activity tracking
		m.lastActivity = time.Now()
		m.activityLevel = 1.0
		
		// Handle global key bindings FIRST - but only for specific keys that should override input
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
			
		case key.Matches(msg, m.keyMap.NextTab):
			m.activeTab = (m.activeTab + 1) % len(m.tabNames)
			m.soundSystem.PlaySound("notification", "tab_next")
			// Focus the appropriate input for the new tab
			m.focusActiveTabInput()
			return m, nil
			
		case key.Matches(msg, m.keyMap.PrevTab):
			m.activeTab = (m.activeTab - 1 + len(m.tabNames)) % len(m.tabNames)
			m.soundSystem.PlaySound("notification", "tab_prev")
			// Focus the appropriate input for the new tab
			m.focusActiveTabInput()
			return m, nil
			
		case key.Matches(msg, m.keyMap.AI):
			m.activeTab = 2 // AI tab
			m.aiInput.Focus()
			return m, nil
			
		case key.Matches(msg, m.keyMap.System):
			m.activeTab = 1 // System tab
			return m, nil
			
		case key.Matches(msg, m.keyMap.Network):
			m.activeTab = 3 // Network tab
			return m, nil
			
		case key.Matches(msg, m.keyMap.Settings):
			m.activeTab = 4 // Settings tab
			m.settingsInput.Focus()
			return m, nil
		}
		
		// Handle tab-specific key events (for special keys like Enter, Tab, Ctrl+L)
		cmd := m.handleTabKeyMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	// CRITICAL FIX: Always update the active input components for ALL message types
	// This ensures typing, pasting, and other input events work properly
	var cmd tea.Cmd
	switch m.activeTab {
	case 0: // Terminal
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		
	case 2: // AI Chat
		m.aiInput, cmd = m.aiInput.Update(msg)
		cmds = append(cmds, cmd)
		m.aiViewport, cmd = m.aiViewport.Update(msg)
		cmds = append(cmds, cmd)
		
	case 4: // Settings
		m.settingsInput, cmd = m.settingsInput.Update(msg)
		cmds = append(cmds, cmd)
		m.settingsViewport, cmd = m.settingsViewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

func (m *EnhancedModel) resizeComponents() {
	contentWidth := m.width - 4
	contentHeight := m.height - 8
	
	// Resize viewports
	m.viewport.Width = contentWidth
	m.viewport.Height = contentHeight
	m.aiViewport.Width = contentWidth
	m.aiViewport.Height = contentHeight
	m.settingsViewport.Width = contentWidth
	m.settingsViewport.Height = contentHeight
	
	// Resize inputs
	m.input.SetWidth(contentWidth - 20)
	m.aiInput.SetWidth(contentWidth - 20)
	m.settingsInput.SetWidth(contentWidth - 20)
}

// focusActiveTabInput focuses the input field for the currently active tab
func (m *EnhancedModel) focusActiveTabInput() {
	switch m.activeTab {
	case 0: // Terminal
		m.input.Focus()
		m.aiInput.Blur()
		m.settingsInput.Blur()
	case 2: // AI Chat
		m.aiInput.Focus()
		m.input.Blur()
		m.settingsInput.Blur()
	case 4: // Settings
		m.settingsInput.Focus()
		m.input.Blur()
		m.aiInput.Blur()
	default:
		// For System and Network tabs, blur all inputs since they don't have input fields
		m.input.Blur()
		m.aiInput.Blur()
		m.settingsInput.Blur()
	}
}

func (m *EnhancedModel) updateAnimations(t time.Time) {
	for name, anim := range m.animations {
		if t.Sub(anim.LastTick) >= anim.Duration {
			anim.Frame = (anim.Frame + 1) % anim.MaxFrames
			anim.LastTick = t
			m.animations[name] = anim
		}
	}
	
	m.pulseFrame = (m.pulseFrame + 1) % 60
}

func (m *EnhancedModel) updatePresence(t time.Time) {
	// Update presence based on activity
	timeSinceActivity := t.Sub(m.lastActivity)
	previousMode := m.presenceState.Mode
	
	if timeSinceActivity > 5*time.Minute {
		m.presenceState.Mode = "away"
		m.activityLevel = 0.1
		
		// React to going away - dim theme, play ambient sound
		if previousMode != "away" {
			m.soundSystem.PlaySound("idle", "presence_away")
		}
	} else if timeSinceActivity > 1*time.Minute {
		m.presenceState.Mode = "idle"
		m.activityLevel = 0.3
		
		// React to going idle
		if previousMode != "idle" {
			m.soundSystem.PlaySound("idle", "presence_idle") 
		}
	} else {
		m.presenceState.Mode = "focus"
		m.activityLevel = 1.0
		
		// React to becoming active again
		if previousMode == "away" || previousMode == "idle" {
			m.soundSystem.PlaySound("notification", "presence_active")
		}
	}
	
	// Dynamic theme adjustment based on presence
	m.adjustThemeForPresence()
	
	m.presenceState.SessionTime += 5 * time.Second
}

// adjustThemeForPresence modifies theme intensity based on presence state
func (m *EnhancedModel) adjustThemeForPresence() {
	switch m.presenceState.Mode {
	case "focus":
		// Full brightness, vibrant colors - already the default
		break
	case "idle":
		// Slightly dimmed but still visible - could implement theme variations
		break
	case "away":
		// Dimmed colors, minimal visual activity
		// Could switch to a darker theme variant
		break
	}
}

func (m *EnhancedModel) handleTabKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch m.activeTab {
	case 0: // Terminal tab
		return m.handleTerminalKeys(msg)
	case 2: // AI Chat tab
		return m.handleAIKeys(msg)
	case 4: // Settings tab
		return m.handleSettingsKeys(msg)
	}
	return nil
}

func (m *EnhancedModel) handleTerminalKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		m.showingSuggestions = false
		return m.handleCommand()
	case tea.KeyTab:
		return m.handleTabCompletion()
	case tea.KeyEsc:
		m.showingSuggestions = false
		return nil
	case tea.KeyCtrlL:
		m.output = []string{}
		m.showingSuggestions = false
		m.updateViewport()
		return nil
	}
	// For all other keys (typing, Ctrl+V paste, etc.), return nil to let the
	// textinput component handle them in the main Update method
	return nil
}

func (m *EnhancedModel) handleAIKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		return m.handleAIQuery()
	}
	// For all other keys (typing, Ctrl+V paste, etc.), return nil to let the
	// textinput component handle them in the main Update method
	return nil
}

func (m *EnhancedModel) handleSettingsKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		return m.handleSettingsCommand()
	}
	// For all other keys (typing, Ctrl+V paste, etc.), return nil to let the
	// textinput component handle them in the main Update method
	return nil
}

func (m *EnhancedModel) handleCommand() tea.Cmd {
	command := strings.TrimSpace(m.input.Value())
	if command == "" {
		return nil
	}
	
	// Track command for presence system
	m.presenceState.LastCommand = command
	
	// Add command to output with enhanced styling
	cmdLine := m.renderCommandLine(command)
	m.output = append(m.output, cmdLine)
	
	// Handle special commands first
	if strings.HasPrefix(command, ":") {
		m.handleSpecialCommand(command)
	} else {
		// Execute through shell
		result, err := m.shell.Execute(command)
		if err != nil {
			errorMsg := m.currentTheme.ErrorStyle.Render(fmt.Sprintf("❌ ERROR: %v", err))
			m.output = append(m.output, errorMsg)
			m.soundSystem.PlaySound("error", "command_error")
		} else {
			// Add output with enhanced formatting
			if result.Output != "" {
				lines := strings.Split(strings.TrimRight(result.Output, "\n"), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						m.output = append(m.output, "  "+line)
					}
				}
			}
			
			if result.Error != "" {
				errorLines := strings.Split(strings.TrimRight(result.Error, "\n"), "\n") 
				for _, line := range errorLines {
					if strings.TrimSpace(line) != "" {
						m.output = append(m.output, m.currentTheme.ErrorStyle.Render("  ⚠️  "+line))
					}
				}
				m.soundSystem.PlaySound("error", "command_error")
			} else if result.ExitCode == 0 {
				m.soundSystem.PlaySound("success", "command_success")
			}
		}
	}
	
	m.input.SetValue("")
	m.updateViewport()
	return nil
}

func (m *EnhancedModel) handleAIQuery() tea.Cmd {
	query := strings.TrimSpace(m.aiInput.Value())
	if query == "" {
		return nil
	}
	
	// Play AI thinking sound
	m.soundSystem.PlaySound("ai_thinking", "ai_query")
	
	// Add query to AI history
	interaction := AIInteraction{
		Query:     query,
		Response:  "",
		Timestamp: time.Now(),
		Markdown:  false,
	}
	m.aiHistory = append(m.aiHistory, interaction)
	
	m.aiThinking = true
	m.aiInput.SetValue("")
	
	// Update AI viewport
	m.updateAIViewport()
	
	// Simulate AI response (replace with actual AI call)
	go func() {
		time.Sleep(2 * time.Second)
		
		response := fmt.Sprintf("I understand you want to know about '%s'. This is a simulated response that would come from the AI engine. The response could include **markdown formatting** and `code blocks` for better readability.", query)
		
		// Update the last interaction
		if len(m.aiHistory) > 0 {
			m.aiHistory[len(m.aiHistory)-1].Response = response
			m.aiHistory[len(m.aiHistory)-1].Markdown = true
		}
		
		m.aiThinking = false
		m.updateAIViewport()
	}()
	
	return nil
}

func (m *EnhancedModel) handleSettingsCommand() tea.Cmd {
	command := strings.TrimSpace(m.settingsInput.Value())
	if command == "" {
		return nil
	}
	
	// Handle settings commands
	parts := strings.SplitN(command, "=", 2)
	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Update setting (simplified)
		settingsOutput := fmt.Sprintf("✅ Setting updated: %s = %s", key, value)
		m.settingsViewport.SetContent(m.settingsViewport.View() + "\n" + settingsOutput)
	}
	
	m.settingsInput.SetValue("")
	return nil
}

func (m *EnhancedModel) handleTabCompletion() tea.Cmd {
	now := time.Now()
	currentInput := m.input.Value()
	cursorPos := m.input.Position()
	
	// Check if this is a double-tab within 500ms for paging
	isDoubleTap := m.showingSuggestions && now.Sub(m.lastTabPress) < 500*time.Millisecond
	m.lastTabPress = now
	
	if !m.showingSuggestions || !isDoubleTap {
		suggestions := m.completion.Complete(currentInput, cursorPos)
		
		if len(suggestions) == 0 {
			m.showingSuggestions = false
			return nil
		}
		
		if len(suggestions) == 1 {
			// Auto-complete with the single suggestion
			text := suggestions[0].Text
			m.input.SetValue(text)
			m.input.SetCursor(len(text))
			m.showingSuggestions = false
		} else {
			// Show multiple suggestions
			m.suggestions = suggestions
			m.showingSuggestions = true
		}
	} else if isDoubleTap {
		// Handle paging through suggestions
		totalPages := (len(m.suggestions) + m.suggestionsPerPage - 1) / m.suggestionsPerPage
		if totalPages > 1 {
			m.suggestionPage = (m.suggestionPage + 1) % totalPages
		}
	}
	
	return nil
}

func (m *EnhancedModel) handleSpecialCommand(command string) {
	parts := strings.Fields(command[1:]) // Remove ':'
	if len(parts) == 0 {
		return
	}
	
	switch parts[0] {
	case "theme":
		if len(parts) > 1 {
			themeName := parts[1]
			if IsValidCatppuccinTheme(themeName) {
				m.currentTheme = GetCatppuccinTheme(themeName)
				m.output = append(m.output, 
					m.currentTheme.SuccessStyle.Render(fmt.Sprintf("🎨 Theme switched to: %s", themeName)))
			} else {
				available := GetCatppuccinThemeNames()
				m.output = append(m.output, 
					m.currentTheme.ErrorStyle.Render(fmt.Sprintf("❌ Unknown theme. Available: %s", 
					strings.Join(available, ", "))))
			}
		}
		
	case "awaken":
		m.showAwakeningSequence()
		
	case "prophecy":
		m.showProphecy()
		
	case "listen":
		m.startAmbientMode()
		
	case "vanish":
		m.showVanishSequence()
	}
}

func (m *EnhancedModel) updateViewport() {
	content := strings.Join(m.output, "\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

func (m *EnhancedModel) updateAIViewport() {
	var content []string
	
	for _, interaction := range m.aiHistory {
		content = append(content, m.currentTheme.SuccessStyle.Render(fmt.Sprintf("🧠 [%s] You:", 
			interaction.Timestamp.Format("15:04"))))
		content = append(content, "  "+interaction.Query)
		content = append(content, "")
		
		if interaction.Response != "" {
			content = append(content, m.currentTheme.InfoStyle.Render("🤖 AI:"))
			if interaction.Markdown && m.glamourRenderer != nil {
				rendered, err := m.glamourRenderer.Render(interaction.Response)
				if err == nil {
					content = append(content, rendered)
				} else {
					content = append(content, "  "+interaction.Response)
				}
			} else {
				content = append(content, "  "+interaction.Response)
			}
		} else if m.aiThinking {
			content = append(content, m.currentTheme.InfoStyle.Render("🤖 AI: Thinking..."))
		}
		
		content = append(content, "")
		content = append(content, strings.Repeat("─", 80))
		content = append(content, "")
	}
	
	m.aiViewport.SetContent(strings.Join(content, "\n"))
	m.aiViewport.GotoBottom()
}

func (m *EnhancedModel) renderCommandLine(command string) string {
	pulse := m.renderPulse()
	return m.currentTheme.PromptStyle.Render(pulse + " ") + 
		   m.currentTheme.CommandStyle.Render(command)
}

func (m *EnhancedModel) renderPulse() string {
	intensity := math.Sin(float64(m.pulseFrame) * 0.2)
	if intensity > 0.5 {
		return "⚡"
	}
	return "▶"
}

// Visual Bliss special sequences
func (m *EnhancedModel) showAwakeningSequence() {
	sequence := []string{
		"",
		m.currentTheme.SuccessStyle.Render("╔════════════════════════════════════════════════════════════════╗"),
		m.currentTheme.SuccessStyle.Render("║                    NEURAL AWAKENING PROTOCOL                  ║"),
		m.currentTheme.SuccessStyle.Render("╚════════════════════════════════════════════════════════════════╝"),
		"",
		m.currentTheme.InfoStyle.Render("🧠 Initializing consciousness matrix..."),
		m.currentTheme.InfoStyle.Render("⚡ Quantum entanglement established..."),
		m.currentTheme.InfoStyle.Render("🌌 Accessing universal knowledge base..."),
		m.currentTheme.SuccessStyle.Render("✨ AWAKENING COMPLETE - Full cognitive abilities online"),
		"",
	}
	m.output = append(m.output, sequence...)
}

func (m *EnhancedModel) showProphecy() {
	prophecies := []string{
		"The shell shall become more than the sum of its commands...",
		"In the matrix of data, patterns emerge that reveal the future...",
		"The neural pathways converge on a solution not yet conceived...",
		"Through the chaos of information, order shall manifest...",
	}
	
	prophecy := prophecies[time.Now().Unix()%int64(len(prophecies))]
	
	sequence := []string{
		"",
		m.currentTheme.SuccessStyle.Render("╔════════════════════════════════════════════════════════════════╗"),
		m.currentTheme.SuccessStyle.Render("║                      NEURAL PROPHECY                          ║"),
		m.currentTheme.SuccessStyle.Render("╚════════════════════════════════════════════════════════════════╝"),
		"",
		m.currentTheme.InfoStyle.Render("🔮 Consulting the quantum oracle..."),
		"",
		m.currentTheme.SecondaryStyle.Render("    \"" + prophecy + "\""),
		"",
		m.currentTheme.SuccessStyle.Render("✨ The prophecy has been spoken"),
		"",
	}
	m.output = append(m.output, sequence...)
}

func (m *EnhancedModel) startAmbientMode() {
	sequence := []string{
		"",
		m.currentTheme.SuccessStyle.Render("🎵 AMBIENT MODE ACTIVATED"),
		m.currentTheme.InfoStyle.Render("♫ Generating neural harmonics..."),
		m.currentTheme.InfoStyle.Render("♪ Binaural focus frequencies enabled"),
		m.currentTheme.SecondaryStyle.Render("   Close your eyes and feel the digital zen..."),
		"",
	}
	m.output = append(m.output, sequence...)
}

func (m *EnhancedModel) showVanishSequence() {
	sequence := []string{
		"",
		m.currentTheme.InfoStyle.Render("👻 Initiating phantom protocol..."),
		m.currentTheme.InfoStyle.Render("🌫️  Activating stealth mode..."),
		m.currentTheme.SecondaryStyle.Render("   Fading into the digital mist..."),
		m.currentTheme.SuccessStyle.Render("💨 Vanish complete - you are one with the void"),
		"",
	}
	m.output = append(m.output, sequence...)
}

func (m *EnhancedModel) View() string {
	if !m.ready {
		return "\n  🧠 Initializing Visual Bliss neural pathways...\n"
	}
	
	// Render tabs
	tabsView := m.renderTabBar()
	
	// Render current tab content
	var contentView string
	switch m.activeTab {
	case 0: // Terminal
		contentView = m.renderTerminalView()
	case 1: // System
		contentView = m.renderSystemView()
	case 2: // AI Chat
		contentView = m.renderAIView()
	case 3: // Network
		contentView = m.renderNetworkView()
	case 4: // Settings
		contentView = m.renderSettingsView()
	}
	
	// Render status bar with presence info
	statusBar := m.renderStatusBar()
	
	// Apply glitch effect if enabled
	if m.glitchEffect {
		contentView = m.applyGlitchEffect(contentView)
	}
	
	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabsView,
		contentView,
		statusBar,
	)
}

// renderTabBar creates a custom tab bar
func (m *EnhancedModel) renderTabBar() string {
	var tabs []string
	
	for i, tabName := range m.tabNames {
		if i == m.activeTab {
			// Active tab
			style := lipgloss.NewStyle().
				Background(lipgloss.Color("#cba6f7")).
				Foreground(lipgloss.Color("#1e1e2e")).
				Bold(true).
				Padding(0, 2)
			tabs = append(tabs, style.Render(tabName))
		} else {
			// Inactive tab
			style := lipgloss.NewStyle().
				Background(lipgloss.Color("#313244")).
				Foreground(lipgloss.Color("#a6adc8")).
				Padding(0, 2)
			tabs = append(tabs, style.Render(tabName))
		}
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
}

func (m *EnhancedModel) renderTerminalView() string {
	content := m.viewport.View()
	inputSection := m.renderInput()
	
	var suggestionSection string
	if m.showingSuggestions && len(m.suggestions) > 0 {
		suggestionSection = m.renderSuggestions()
	}
	
	components := []string{content, inputSection}
	if suggestionSection != "" {
		components = append(components, suggestionSection)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, components...)
}

func (m *EnhancedModel) renderSystemView() string {
	stats := []string{
		m.currentTheme.SuccessStyle.Render("📊 SYSTEM VITALS") + " " + m.systemSpinner.View(),
		"",
		fmt.Sprintf("🖥️  CPU Usage: %.1f%%", m.systemStats.CPUUsage),
		fmt.Sprintf("💾 Memory: %.1f%%", m.systemStats.MemoryUsage),
		fmt.Sprintf("💽 Disk: %.1f%%", m.systemStats.DiskUsage),
		fmt.Sprintf("⚡ Active Jobs: %d", m.systemStats.ActiveJobs),
		fmt.Sprintf("⏱️  Uptime: %v", m.systemStats.Uptime),
		"",
		m.currentTheme.InfoStyle.Render("Press CTRL+Y to refresh system stats"),
	}
	
	return strings.Join(stats, "\n")
}

func (m *EnhancedModel) renderAIView() string {
	content := m.aiViewport.View()
	input := m.renderAIInput()
	
	return lipgloss.JoinVertical(lipgloss.Left, content, input)
}

func (m *EnhancedModel) renderNetworkView() string {
	networkInfo := []string{
		m.currentTheme.SuccessStyle.Render("🌐 NETWORK MONITORING"),
		"",
		fmt.Sprintf("📊 Active Connections: %d", len(m.networkStats.Connections)),
		fmt.Sprintf("📈 Bytes In: %d", m.networkStats.Bandwidth.BytesIn),
		fmt.Sprintf("📉 Bytes Out: %d", m.networkStats.Bandwidth.BytesOut),
		"",
	}
	
	// Add connection details
	for i, conn := range m.networkStats.Connections {
		if i >= 10 { // Limit display
			break
		}
		networkInfo = append(networkInfo, fmt.Sprintf("%s → %s [%s]", 
			conn.LocalAddr, conn.RemoteAddr, conn.State))
	}
	
	return strings.Join(networkInfo, "\n")
}

func (m *EnhancedModel) renderSettingsView() string {
	content := m.settingsViewport.View()
	input := m.renderSettingsInput()
	
	return lipgloss.JoinVertical(lipgloss.Left, content, input)
}

func (m *EnhancedModel) renderInput() string {
	prompt := m.currentTheme.PromptStyle.Render("lucien@nexus:~$ ")
	return prompt + m.input.View()
}

func (m *EnhancedModel) renderAIInput() string {
	prompt := m.currentTheme.PromptStyle.Render("🧠 AI> ")
	return prompt + m.aiInput.View()
}

func (m *EnhancedModel) renderSettingsInput() string {
	prompt := m.currentTheme.PromptStyle.Render("⚙️  > ")
	return prompt + m.settingsInput.View()
}

func (m *EnhancedModel) renderSuggestions() string {
	if len(m.suggestions) == 0 {
		return ""
	}
	
	// Implement suggestion rendering similar to the original model
	var suggestionLines []string
	suggestionLines = append(suggestionLines, m.currentTheme.SuccessStyle.Render("🔧 TAB COMPLETION:"))
	
	for i, suggestion := range m.suggestions {
		if i >= m.suggestionsPerPage {
			break
		}
		icon := "▶"
		switch suggestion.Type {
		case completion.CommandCompletion:
			icon = "⚡"
		case completion.FileCompletion:
			icon = "📄"
		case completion.DirectoryCompletion:
			icon = "📁"
		}
		suggestionLines = append(suggestionLines, fmt.Sprintf("  %s %s", icon, suggestion.Text))
	}
	
	return strings.Join(suggestionLines, "\n")
}

func (m *EnhancedModel) renderStatusBar() string {
	// Presence indicator
	var presenceIcon string
	switch m.presenceState.Mode {
	case "focus":
		presenceIcon = "🔴"
	case "idle":
		presenceIcon = "🟡"
	case "away":
		presenceIcon = "⚫"
	}
	
	left := fmt.Sprintf("%s %s", presenceIcon, m.presenceState.Mode)
	right := fmt.Sprintf("Session: %v | Theme: %s", 
		m.presenceState.SessionTime.Truncate(time.Second), m.currentTheme.Name)
	
	leftStyled := m.currentTheme.FooterStyle.Render(left)
	rightStyled := m.currentTheme.FooterStyle.Render(right)
	
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, leftStyled) +
		   lipgloss.PlaceHorizontal(m.width, lipgloss.Right, rightStyled)
}

func (m *EnhancedModel) applyGlitchEffect(content string) string {
	// Simple glitch effect - add random characters
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if time.Now().Unix()%10 == 0 && i%5 == 0 { // Random glitch
			glitchChars := "▓▒░█▄▀■□▪▫"
			glitchChar := string(glitchChars[i%len(glitchChars)])
			lines[i] = line + m.currentTheme.ErrorStyle.Render(glitchChar)
		}
	}
	return strings.Join(lines, "\n")
}