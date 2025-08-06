package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luciendev/lucien-core/internal/ai"
	"github.com/luciendev/lucien-core/internal/policy"
	"github.com/luciendev/lucien-core/internal/plugin"
	"github.com/luciendev/lucien-core/internal/sandbox"
	"github.com/luciendev/lucien-core/internal/shell"
	"github.com/luciendev/lucien-core/internal/ui"
)

var (
	sshFlag    = flag.Bool("ssh", false, "Start SSH server for remote access")
	configFlag = flag.String("config", "", "Path to config file (default: ~/.lucien/config.toml)")
	safeMode   = flag.Bool("safe-mode", false, "Enable OPA policy enforcement")
	port       = flag.String("port", "2222", "SSH server port")
	version    = "1.0.0-alpha"
)

// Core represents the Lucien shell core with all subsystems
type Core struct {
	Shell   *shell.Shell
	Policy  *policy.Engine
	Plugin  *plugin.Manager
	Sandbox *sandbox.Manager
	AI      *ai.Engine
}

func main() {
	flag.Parse()

	// ASCII art banner with retro vibe
	banner := `
██╗     ██╗   ██╗ ██████╗██╗███████╗███╗   ██╗    ███████╗██╗  ██╗███████╗██╗     ██╗     
██║     ██║   ██║██╔════╝██║██╔════╝████╗  ██║    ██╔════╝██║  ██║██╔════╝██║     ██║     
██║     ██║   ██║██║     ██║█████╗  ██╔██╗ ██║    ███████╗███████║█████╗  ██║     ██║     
██║     ██║   ██║██║     ██║██╔══╝  ██║╚██╗██║    ╚════██║██╔══██║██╔══╝  ██║     ██║     
███████╗╚██████╔╝╚██████╗██║███████╗██║ ╚████║    ███████║██║  ██║███████╗███████╗███████╗
╚══════╝ ╚═════╝  ╚═════╝╚═╝╚══════╝╚═╝  ╚═══╝    ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝
                                      [ NEXUS-7 TERMINAL INTERFACE v%s ]
`

	fmt.Printf(banner, version)
	fmt.Println("⚡ Initializing neural pathways...")

	// Initialize core systems
	core, err := initCore()
	if err != nil {
		log.Fatalf("❌ Core initialization failed: %v", err)
	}

	fmt.Println("✅ Core systems online")
	fmt.Println("🧠 AI subsystems loaded")
	fmt.Println("🛡️  Security protocols active")

	if *sshFlag {
		fmt.Println("❌ SSH functionality temporarily disabled")
		fmt.Println("🚀 Starting in local mode instead...")
	}
	startLocal(core)
}

func initCore() (*Core, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	lucienDir := filepath.Join(homeDir, ".lucien")
	if err := os.MkdirAll(lucienDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create lucien directory: %w", err)
	}

	// Initialize policy engine
	policyEngine, err := policy.New(filepath.Join(lucienDir, "policies"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize policy engine: %w", err)
	}

	// Initialize sandbox
	sandboxMgr, err := sandbox.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sandbox: %w", err)
	}

	// Initialize plugin manager
	pluginMgr := plugin.New(filepath.Join(lucienDir, "plugins"))

	// Initialize AI engine
	aiEngine, err := ai.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI engine: %w", err)
	}

	// Initialize shell
	shellEngine := shell.New(&shell.Config{
		PolicyEngine: policyEngine,
		PluginMgr:   pluginMgr,
		SandboxMgr:  sandboxMgr,
		AIEngine:    aiEngine,
		SafeMode:    *safeMode,
	})

	return &Core{
		Shell:   shellEngine,
		Policy:  policyEngine,
		Plugin:  pluginMgr,
		Sandbox: sandboxMgr,
		AI:      aiEngine,
	}, nil
}

func startLocal(core *Core) {
	model := ui.NewModel(core.Shell, core.AI)
	program := tea.NewProgram(model, tea.WithAltScreen())

	fmt.Println("🚀 Launching local interface...")
	
	if _, err := program.Run(); err != nil {
		log.Fatalf("❌ TUI crashed: %v", err)
	}
}

// SSH functionality temporarily removed for build compatibility