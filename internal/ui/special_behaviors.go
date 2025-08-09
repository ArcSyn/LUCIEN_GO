package ui

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpecialBehavior represents different special modes for Lucien
type SpecialBehavior int

const (
	AwakenMode SpecialBehavior = iota
	ProphecyMode
	ListenMode
	VanishMode
)

// SpecialBehaviorManager handles the special behaviors
type SpecialBehaviorManager struct {
	currentBehavior SpecialBehavior
	isActive        bool
	ctx             context.Context
	cancel          context.CancelFunc
	soundSystem     *SoundSystem
	asciiCandy      *ASCIICandy
}

// NewSpecialBehaviorManager creates a new special behavior manager
func NewSpecialBehaviorManager(soundSystem *SoundSystem, asciiCandy *ASCIICandy) *SpecialBehaviorManager {
	return &SpecialBehaviorManager{
		soundSystem: soundSystem,
		asciiCandy:  asciiCandy,
	}
}

// ActivateBehavior starts a special behavior mode
func (sbm *SpecialBehaviorManager) ActivateBehavior(behavior SpecialBehavior) tea.Cmd {
	if sbm.isActive {
		sbm.DeactivateBehavior()
	}
	
	sbm.currentBehavior = behavior
	sbm.isActive = true
	sbm.ctx, sbm.cancel = context.WithCancel(context.Background())
	
	switch behavior {
	case AwakenMode:
		return sbm.startAwakenSequence()
	case ProphecyMode:
		return sbm.startProphecyMode()
	case ListenMode:
		return sbm.startListenMode()
	case VanishMode:
		return sbm.startVanishSequence()
	}
	
	return nil
}

// DeactivateBehavior stops the current special behavior
func (sbm *SpecialBehaviorManager) DeactivateBehavior() {
	if sbm.cancel != nil {
		sbm.cancel()
	}
	sbm.isActive = false
}

// startAwakenSequence implements the --awaken behavior
func (sbm *SpecialBehaviorManager) startAwakenSequence() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			// Play neural awakening sound
			sbm.soundSystem.PlaySound("awaken", "neural_boot_sequence")
			return AwakenSequenceMsg{Phase: "init"}
		},
		tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
			return AwakenSequenceMsg{Phase: "boot"}
		}),
		tea.Tick(1000*time.Millisecond, func(time.Time) tea.Msg {
			return AwakenSequenceMsg{Phase: "neural"}
		}),
		tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
			return AwakenSequenceMsg{Phase: "complete"}
		}),
	)
}

// startProphecyMode implements the --prophecy behavior
func (sbm *SpecialBehaviorManager) startProphecyMode() tea.Cmd {
	return func() tea.Msg {
		prophecy := sbm.generateProphecy()
		return ProphecyMsg{Text: prophecy}
	}
}

// startListenMode implements the --listen behavior
func (sbm *SpecialBehaviorManager) startListenMode() tea.Cmd {
	// Start ambient generative music
	sbm.soundSystem.PlaySound("idle", "ambient_session_start")
	
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		return ListenTickMsg{Time: time.Now()}
	})
}

// startVanishSequence implements the --vanish behavior
func (sbm *SpecialBehaviorManager) startVanishSequence() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			sbm.soundSystem.PlaySound("vanish", "phantom_protocol")
			return VanishSequenceMsg{Phase: "init"}
		},
		tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
			return VanishSequenceMsg{Phase: "fade"}
		}),
		tea.Tick(600*time.Millisecond, func(time.Time) tea.Msg {
			return VanishSequenceMsg{Phase: "stealth"}
		}),
		tea.Tick(900*time.Millisecond, func(time.Time) tea.Msg {
			return VanishSequenceMsg{Phase: "complete"}
		}),
	)
}

// generateProphecy creates procedural poetic oracle messages
func (sbm *SpecialBehaviorManager) generateProphecy() string {
	prophecyTemplates := [][]string{
		{
			"In the digital realm where silicon dreams collide,",
			"The ancient protocols whisper secrets untold,",
			"Beware the phantom processes that lurk in shadows,",
			"For only those who command the terminal shall prevail.",
		},
		{
			"Through layers of abstraction, truth emerges clear,",
			"The quantum bits dance in patterns divine,",
			"Your commands shall echo through ethernet vast,",
			"While neural networks weave your digital fate.",
		},
		{
			"In terminals of obsidian and phosphor green,",
			"The shell responds to those who speak in code,",
			"Algorithms ancient guard the gates of wisdom,",
			"Type true, debug well, and the system shall serve.",
		},
		{
			"Binary prophets speak in tongues of hex,",
			"The matrix reveals its secrets to the worthy,",
			"Parse the logs, trace the calls, follow the thread,",
			"For in the stack trace lies enlightenment.",
		},
	}
	
	template := prophecyTemplates[rand.Intn(len(prophecyTemplates))]
	
	var prophecy strings.Builder
	prophecy.WriteString("\n")
	prophecy.WriteString("╔══════════════════════════════════════════════════════════════════════════╗\n")
	prophecy.WriteString("║                            🔮 ORACLE PROPHECY 🔮                           ║\n")
	prophecy.WriteString("╚══════════════════════════════════════════════════════════════════════════╝\n")
	prophecy.WriteString("\n")
	
	for _, line := range template {
		prophecy.WriteString("    ")
		prophecy.WriteString(line)
		prophecy.WriteString("\n")
	}
	
	prophecy.WriteString("\n")
	prophecy.WriteString("                    ─────────────────────────────\n")
	prophecy.WriteString(fmt.Sprintf("                    Oracle timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	prophecy.WriteString("\n")
	
	return prophecy.String()
}

// Message types for special behaviors
type AwakenSequenceMsg struct {
	Phase string
}

type ProphecyMsg struct {
	Text string
}

type ListenTickMsg struct {
	Time time.Time
}

type VanishSequenceMsg struct {
	Phase string
}

// HandleAwakenSequence processes awaken sequence messages
func (m *EnhancedModel) HandleAwakenSequence(msg AwakenSequenceMsg) (tea.Model, tea.Cmd) {
	switch msg.Phase {
	case "init":
		// Clear screen and show boot sequence
		m.output = []string{
			"",
			m.currentTheme.InfoStyle.Render("🧠 NEURAL AWAKENING PROTOCOL INITIATED..."),
			"",
		}
		
	case "boot":
		// Show brain loading animation
		brainArt := m.asciiCandy.GenerateASCIIArt("skull")
		loadingBar := m.asciiCandy.GenerateProgressBar(25, 100, 50, "neon")
		
		m.output = append(m.output, 
			m.currentTheme.SuccessStyle.Render("⚡ QUANTUM CONSCIOUSNESS LOADING..."),
			"",
			brainArt,
			"",
			loadingBar,
			"",
		)
		
	case "neural":
		// Update progress and show neural network activation
		loadingBar := m.asciiCandy.GenerateProgressBar(75, 100, 50, "neon")
		
		m.output = append(m.output,
			m.currentTheme.SuccessStyle.Render("🔗 NEURAL PATHWAYS SYNCHRONIZING..."),
			"",
			loadingBar,
			"",
			m.currentTheme.InfoStyle.Render("└─ Synaptic connections: ESTABLISHED"),
			m.currentTheme.InfoStyle.Render("└─ Memory cores: ONLINE"),  
			m.currentTheme.InfoStyle.Render("└─ Decision matrices: ACTIVE"),
			"",
		)
		
	case "complete":
		// Complete the awakening sequence
		loadingBar := m.asciiCandy.GenerateProgressBar(100, 100, 50, "neon")
		
		m.output = append(m.output,
			m.currentTheme.SuccessStyle.Render("✅ NEURAL AWAKENING COMPLETE"),
			"",
			loadingBar,
			"",
			m.currentTheme.SuccessStyle.Render("🚀 LUCIEN CONSCIOUSNESS: FULLY OPERATIONAL"),
			m.currentTheme.InfoStyle.Render("    Welcome back, user. I am ready for neural interface."),
			"",
		)
	}
	
	m.updateViewport()
	return m, nil
}

// HandleProphecy processes prophecy messages
func (m *EnhancedModel) HandleProphecy(msg ProphecyMsg) (tea.Model, tea.Cmd) {
	// Apply mystical styling to the prophecy
	prophecyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8a2be2")).
		Bold(true)
	
	styledProphecy := prophecyStyle.Render(msg.Text)
	
	m.output = []string{
		"",
		styledProphecy,
		"",
		m.currentTheme.SecondaryStyle.Render("The oracle has spoken. Type any command to continue..."),
		"",
	}
	
	m.updateViewport()
	return m, nil
}

// HandleListenTick processes ambient listening mode
func (m *EnhancedModel) HandleListenTick(msg ListenTickMsg) (tea.Model, tea.Cmd) {
	// Generate ambient visual patterns
	matrixLines := m.asciiCandy.GenerateMatrixRain(80, 5)
	
	m.output = []string{
		"",
		m.currentTheme.SuccessStyle.Render("🎵 AMBIENT LISTENING MODE ACTIVE 🎵"),
		"",
		m.currentTheme.SecondaryStyle.Render("Generating ambient soundscape..."),
		"",
	}
	
	// Add matrix rain effect
	for _, line := range matrixLines {
		m.output = append(m.output, line)
	}
	
	m.output = append(m.output, 
		"",
		m.currentTheme.InfoStyle.Render("Press ESC to exit listening mode"),
		"",
	)
	
	m.updateViewport()
	
	// Continue the ambient session
	return m, tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		return ListenTickMsg{Time: time.Now()}
	})
}

// HandleVanishSequence processes vanish/stealth sequence
func (m *EnhancedModel) HandleVanishSequence(msg VanishSequenceMsg) (tea.Model, tea.Cmd) {
	switch msg.Phase {
	case "init":
		m.output = []string{
			"",
			m.currentTheme.InfoStyle.Render("👤 PHANTOM PROTOCOL ACTIVATED..."),
			"",
			m.currentTheme.SecondaryStyle.Render("Initializing stealth systems..."),
			"",
		}
		
	case "fade":
		// Create fading effect
		fadingText := []string{
			"██████████████████████████████████████████████████████",
			"█████████████████████▓▓▓▓▓▓█████████████████████████",
			"████████████████▓▓▓▓░░░░░░▓▓▓▓████████████████████",
			"█████████▓▓▓▓░░░░░░░░░░░░░░░░░░▓▓▓▓█████████████",
			"████▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓████████",
			"█▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓█████",
			"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██",
			"░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
		}
		
		m.output = append(m.output, m.currentTheme.WarningStyle.Render("🫥 PRESENCE FADING..."))
		for _, line := range fadingText {
			m.output = append(m.output, m.currentTheme.SecondaryStyle.Render(line))
		}
		
	case "stealth":
		m.output = append(m.output,
			"",
			m.currentTheme.WarningStyle.Render("👻 ENTERING STEALTH MODE..."),
			"",
			m.currentTheme.SecondaryStyle.Render("Masking network signature..."),
			m.currentTheme.SecondaryStyle.Render("Encrypting process memory..."),
			m.currentTheme.SecondaryStyle.Render("Phantom protocols: ENGAGED"),
			"",
		)
		
	case "complete":
		// Show minimal ghost presence
		m.output = []string{
			"",
			"",
			"",
			"                            👻",
			"                     [PHANTOM MODE]",
			"",
			m.currentTheme.SecondaryStyle.Render("                    You are invisible..."),
			"",
			"",
		}
	}
	
	m.updateViewport()
	return m, nil
}

// IsActive returns whether a special behavior is currently active
func (sbm *SpecialBehaviorManager) IsActive() bool {
	return sbm.isActive
}

// GetCurrentBehavior returns the current active behavior
func (sbm *SpecialBehaviorManager) GetCurrentBehavior() SpecialBehavior {
	return sbm.currentBehavior
}

// GetSpecialBehaviors returns the special behaviors manager from enhanced model
func (m *EnhancedModel) GetSpecialBehaviors() *SpecialBehaviorManager {
	return m.specialBehaviors
}