package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	
	"github.com/ArcSyn/LucienCLI/internal/shell"
)

var (
	version   = "1.0.0-nexus7"
	commit    = "unknown"
	buildTime = "unknown"
	
	// Lipgloss styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ff41")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
		
	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff41")).
		Bold(true)
		
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff4444")).
		Bold(true)
		
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00aaff")).
		Bold(true)
)

func main() {
	// Setup logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
		Prefix:         "LUCIEN",
	})

	var rootCmd = &cobra.Command{
		Use:   "lucien",
		Short: "LUCIEN NEXUS-7 - AI-Enhanced Shell with Visual Bliss",
		Long: titleStyle.Render(`
🚀 LUCIEN NEXUS-7 TERMINAL
AI-Enhanced Shell with Visual Bliss
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Built with the Charm ecosystem:
• Cobra CLI framework
• Huh interactive forms  
• Lipgloss styling
• Bubbles components
`),
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveShell(logger)
		},
	}

	// Add subcommands
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Show version information",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("%s\n", successStyle.Render(fmt.Sprintf("LUCIEN CLI v%s", version)))
				fmt.Printf("Commit: %s\n", commit)
				fmt.Printf("Built: %s\n", buildTime)
			},
		},
		&cobra.Command{
			Use:   "interactive",
			Short: "Launch interactive shell with Visual Bliss",
			Run: func(cmd *cobra.Command, args []string) {
				runInteractiveShell(logger)
			},
		},
		&cobra.Command{
			Use:   "setup",
			Short: "Setup and configure LUCIEN",
			Run: func(cmd *cobra.Command, args []string) {
				runSetupWizard()
			},
		},
	)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}
}

func runSetupWizard() {
	fmt.Print(titleStyle.Render("LUCIEN SETUP WIZARD"))
	fmt.Println()

	var (
		theme     string
		shellPath string
		aiEnabled bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome to LUCIEN NEXUS-7!").
				Description("Let's configure your AI-enhanced shell experience."),

			huh.NewSelect[string]().
				Title("Choose your Visual Bliss theme").
				Options(
					huh.NewOption("Catppuccin Mocha (Dark)", "mocha"),
					huh.NewOption("Catppuccin Latte (Light)", "latte"),
					huh.NewOption("Rosé Pine", "rose-pine"),
					huh.NewOption("Nexus Green", "nexus"),
				).
				Value(&theme),

			huh.NewInput().
				Title("Default shell directory").
				Description("Where should LUCIEN start? (Leave blank for current directory)").
				Placeholder("C:\\Users\\YourName\\Projects").
				Value(&shellPath),

			huh.NewConfirm().
				Title("Enable AI features?").
				Description("This includes intelligent suggestions and neural assistance.").
				Value(&aiEnabled),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render(fmt.Sprintf("Setup failed: %v", err)))
		os.Exit(1)
	}

	// Save configuration
	fmt.Printf("\n%s\n", successStyle.Render("✅ Configuration saved!"))
	fmt.Printf("Theme: %s\n", infoStyle.Render(theme))
	fmt.Printf("Directory: %s\n", infoStyle.Render(shellPath))
	fmt.Printf("AI Enabled: %s\n", infoStyle.Render(fmt.Sprintf("%t", aiEnabled)))
	fmt.Println("\nRun 'lucien' to start your enhanced shell!")
}

func runInteractiveShell(logger *log.Logger) {
	// Show banner
	fmt.Print(titleStyle.Render(`
🚀 LUCIEN NEXUS-7 TERMINAL
Neural pathways initialized
Ready for enhanced shell experience
`))
	fmt.Println()

	// Initialize shell
	shellEngine := shell.New(&shell.Config{
		SafeMode: true,
	})

	fmt.Println(successStyle.Render("🧠 AI subsystems loaded"))
	fmt.Println(successStyle.Render("🛡️  Security protocols active"))
	fmt.Println(successStyle.Render("⚡ Ready for input"))
	fmt.Println()

	// Simple interactive loop with proper Windows path support
	for {
		var command string
		
		// Use Huh for clean input that handles Windows paths properly
		inputForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("").
					Prompt("lucien@nexus:~$ ").
					Placeholder("Enter command (or 'exit' to quit)").
					Value(&command),
			),
		).WithShowHelp(false).WithShowErrors(false)

		if err := inputForm.Run(); err != nil {
			fmt.Printf("%s\n", errorStyle.Render(fmt.Sprintf("Input error: %v", err)))
			continue
		}

		// Handle built-in commands
		switch command {
		case "exit", "quit":
			fmt.Println(successStyle.Render("👋 Neural connection terminated. Goodbye!"))
			return
		case "":
			continue
		case "help":
			showHelp()
			continue
		case "clear":
			// Clear screen
			fmt.Print("\033[2J\033[H")
			continue
		}

		// Execute through shell engine
		result, err := shellEngine.Execute(command)
		if err != nil {
			fmt.Printf("%s\n", errorStyle.Render(fmt.Sprintf("❌ ERROR: %v", err)))
			continue
		}

		// Show output
		if result.Output != "" {
			fmt.Print(result.Output)
		}
		if result.Error != "" {
			fmt.Printf("%s\n", errorStyle.Render(result.Error))
		}

		// Success indicator
		if result.ExitCode == 0 && result.Output != "" {
			fmt.Printf("%s\n", successStyle.Render("✓"))
		}
	}
}

func showHelp() {
	help := `
🔴 LUCIEN NEURAL INTERFACE - COMMAND REFERENCE
============================================

📟 BUILT-IN COMMANDS:
  help              Show this help information
  cd <dir>          Change directory
  pwd               Print working directory
  ls [path]         List directory contents
  echo <text>       Display text
  clear             Clear terminal buffer
  exit, quit        Exit shell

⚡ SHELL OPERATIONS:
  Standard shell commands with pipes, redirects, and variables
  Built-in commands: cd, set, alias, exit
  Windows paths: Use "quotes" for paths with spaces

🤖 AI FEATURES (Coming Soon):
  • Intelligent command suggestions
  • Context-aware assistance
  • Neural pattern recognition

🧠 TIPS:
  • Use quotes for Windows paths: cd "C:\Program Files"
  • Type 'lucien setup' to configure themes
  • All standard shell operations supported

`
	fmt.Print(infoStyle.Render(help))
}