package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ArcSyn/LucienCLI/internal/ai"
	"github.com/ArcSyn/LucienCLI/internal/env"
	"github.com/ArcSyn/LucienCLI/internal/history"
	"github.com/ArcSyn/LucienCLI/internal/jobs"
	"github.com/ArcSyn/LucienCLI/internal/policy"
	"github.com/ArcSyn/LucienCLI/internal/plugin"
	"github.com/ArcSyn/LucienCLI/internal/sandbox"
	"github.com/ArcSyn/LucienCLI/internal/shell"
	"github.com/ArcSyn/LucienCLI/internal/ui"
)

var (
	sshFlag     = flag.Bool("ssh", false, "Start SSH server for remote access")
	configFlag  = flag.String("config", "", "Path to config file (default: ~/.lucien/config.toml)")
	safeMode    = flag.Bool("safe-mode", true, "Enable OPA policy enforcement (default: true)")
	unsafeMode  = flag.Bool("unsafe-mode", false, "Disable security validation (NOT RECOMMENDED)")
	port        = flag.String("port", "2222", "SSH server port")
	versionFlag = flag.Bool("version", false, "Show version information")
	batchFlag   = flag.Bool("batch", false, "Run in batch mode (non-interactive)")
	
	// Special Visual Bliss behaviors
	awakenFlag   = flag.Bool("awaken", false, "Animated boot with sound + 'brain loading' bar")
	prophecyFlag = flag.Bool("prophecy", false, "Procedural poetic oracle")
	listenFlag   = flag.Bool("listen", false, "Ambient generative music")
	vanishFlag   = flag.Bool("vanish", false, "Stylish shutdown with fade")
	
	version     = "1.0.0-nexus7"
	commit      = "unknown"
	buildTime   = "unknown"
)

// Core represents the Lucien shell core with all subsystems
type Core struct {
	Shell   *shell.Shell
	Policy  *policy.Engine
	Plugin  *plugin.Manager
	Sandbox *sandbox.Manager
	AI      *ai.Engine
	History *history.Manager
	Jobs    *jobs.Manager
	Env     *env.Manager
}

func main() {
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("Lucien CLI version %s\nCommit: %s\nBuild time: %s\n", version, commit, buildTime)
		os.Exit(0)
	}

	// Handle special behavior flags
	if *awakenFlag || *prophecyFlag || *listenFlag || *vanishFlag {
		runSpecialBehavior()
		return
	}

	// Check batch mode flag or piped input
	if *batchFlag {
		runBatch()
		return
	}
	
	// Check if input is piped
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		// Input is piped, run in non-interactive mode
		runBatch()
		return
	}

	// Run in interactive mode
	runInteractive()
}

func runSpecialBehavior() {
	// Initialize minimal core systems for special behaviors
	core, err := initCore()
	if err != nil {
		log.Fatalf("❌ Core initialization failed: %v", err)
	}

	// Create the Enhanced Visual Bliss UI model
	model := ui.NewEnhancedModel(core.Shell, core.AI)

	// Determine which special behavior to activate
	var behavior ui.SpecialBehavior
	if *awakenFlag {
		behavior = ui.AwakenMode
	} else if *prophecyFlag {
		behavior = ui.ProphecyMode
	} else if *listenFlag {
		behavior = ui.ListenMode
	} else if *vanishFlag {
		behavior = ui.VanishMode
	}

	// Activate the special behavior
	cmd := model.GetSpecialBehaviors().ActivateBehavior(behavior)

	// Create a program that will run the special behavior
	program := tea.NewProgram(model, tea.WithAltScreen())
	
	// Send the special behavior activation command
	if cmd != nil {
		program.Send(cmd())
	}

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatalf("❌ Special behavior failed: %v", err)
	}
}

func runBatch() {
	// Initialize core systems quietly
	core, err := initCore()
	if err != nil {
		log.Fatalf("Core initialization failed: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		result, err := core.Shell.Execute(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		if result.Output != "" {
			fmt.Print(result.Output)
		}
		if result.Error != "" {
			fmt.Fprint(os.Stderr, result.Error)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func runInteractive() {
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

	// Initialize plugin bridge for Python agents
	// pluginBridge, err := plugin.NewBridge(filepath.Join(lucienDir, "plugins"))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize plugin bridge: %w", err)
	// }

	// Initialize AI engine
	aiEngine, err := ai.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI engine: %w", err)
	}

	// Initialize history manager
	historyMgr, err := history.New(&history.Config{
		HistoryFile: filepath.Join(lucienDir, "history.jsonl"),
		MaxEntries:  10000,
		AutoSave:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize history manager: %w", err)
	}

	// Initialize job manager
	jobMgr := jobs.New()

	// Initialize environment manager
	envMgr, err := env.New(&env.Config{
		PersistFile: filepath.Join(lucienDir, "environment.json"),
		AutoSave:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize environment manager: %w", err)
	}

	// Determine final security mode - safe by default unless explicitly disabled
	finalSafeMode := *safeMode && !*unsafeMode
	
	if *unsafeMode {
		fmt.Println("⚠️  WARNING: Running in UNSAFE MODE - Security validation disabled!")
		fmt.Println("⚠️  This mode should ONLY be used in secure, controlled environments")
	} else {
		fmt.Println("🛡️  Security validation enabled - Safe mode active")
	}

	// Initialize shell
	shellEngine := shell.New(&shell.Config{
		SafeMode: finalSafeMode,
	})

	return &Core{
		Shell:   shellEngine,
		Policy:  policyEngine,
		Plugin:  pluginMgr,
		Sandbox: sandboxMgr,
		AI:      aiEngine,
		History: historyMgr,
		Jobs:    jobMgr,
		Env:     envMgr,
	}, nil
}

func startLocal(core *Core) {
	fmt.Println("⚡ Initializing NEXUS-7 neural pathways...")
	fmt.Println("🧠 Loading cyberpunk interface...")
	fmt.Println()
	
	// Create the Enhanced Visual Bliss UI model 
	model := ui.NewEnhancedModel(core.Shell, core.AI)
	
	// Launch the Enhanced Visual Bliss TUI
	program := tea.NewProgram(model, tea.WithAltScreen())
	
	if _, err := program.Run(); err != nil {
		log.Fatalf("❌ Neural interface initialization failed: %v", err)
	}
	
	fmt.Println("👋 Neural connection terminated. Goodbye!")
}

// SSH functionality temporarily removed for build compatibility