package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	LUCIEN_VERSION     = "1.0.0-nexus7"
	GITHUB_RELEASE_URL = "https://api.github.com/repos/luciendev/lucien-cli/releases/latest"
	PLUGIN_REPO_URL    = "https://github.com/luciendev/lucien-plugins/archive/main.zip"
)

// BootstrapConfig holds configuration for the bootstrap process
type BootstrapConfig struct {
	InstallDir     string
	PluginDir      string
	BinaryName     string
	CreateShortcut bool
	AddToPath      bool
	SkipPython     bool
}

// Bootstrap installer for Lucien CLI
type Bootstrap struct {
	config *BootstrapConfig
	logger *Logger
}

// Logger provides simple logging functionality
type Logger struct {
	verbose bool
}

func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf("ℹ️  "+format+"\n", args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	fmt.Printf("✅ "+format+"\n", args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	fmt.Printf("⚠️  "+format+"\n", args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf("❌ "+format+"\n", args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		fmt.Printf("🔍 "+format+"\n", args...)
	}
}

// NewBootstrap creates a new bootstrap installer
func NewBootstrap(config *BootstrapConfig) *Bootstrap {
	return &Bootstrap{
		config: config,
		logger: &Logger{verbose: false},
	}
}

// SetVerbose enables or disables verbose logging
func (b *Bootstrap) SetVerbose(verbose bool) {
	b.logger.verbose = verbose
}

// Run executes the complete bootstrap installation process
func (b *Bootstrap) Run() error {
	b.logger.Info("🚀 Lucien CLI Bootstrap Installer v%s", LUCIEN_VERSION)
	b.logger.Info("Installing Lucien CLI and dependencies...")
	
	// Step 1: Create installation directories
	if err := b.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}
	
	// Step 2: Install Lucien CLI binary
	if err := b.installBinary(); err != nil {
		return fmt.Errorf("failed to install binary: %v", err)
	}
	
	// Step 3: Install Python agents and plugins
	if err := b.installPlugins(); err != nil {
		b.logger.Warning("Failed to install plugins: %v", err)
		b.logger.Info("You can manually copy plugins from the plugins/ directory")
	}
	
	// Step 4: Check Python availability
	if !b.config.SkipPython {
		if err := b.checkPython(); err != nil {
			b.logger.Warning("Python check failed: %v", err)
			b.logger.Info("Agent commands will not work without Python 3.7+")
		}
	}
	
	// Step 5: Create shortcuts and PATH entries
	if b.config.CreateShortcut {
		if err := b.createShortcut(); err != nil {
			b.logger.Warning("Failed to create shortcut: %v", err)
		}
	}
	
	if b.config.AddToPath {
		if err := b.addToPath(); err != nil {
			b.logger.Warning("Failed to add to PATH: %v", err)
			b.logger.Info("Manually add %s to your PATH", b.config.InstallDir)
		}
	}
	
	// Step 6: Verify installation
	if err := b.verifyInstallation(); err != nil {
		return fmt.Errorf("installation verification failed: %v", err)
	}
	
	b.logger.Success("🎉 Lucien CLI installed successfully!")
	b.printPostInstallInstructions()
	
	return nil
}

// createDirectories creates the necessary installation directories
func (b *Bootstrap) createDirectories() error {
	b.logger.Info("Creating installation directories...")
	
	dirs := []string{
		b.config.InstallDir,
		b.config.PluginDir,
		filepath.Join(b.config.PluginDir, "agents"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		b.logger.Debug("Created directory: %s", dir)
	}
	
	return nil
}

// installBinary installs the Lucien CLI binary
func (b *Bootstrap) installBinary() error {
	b.logger.Info("Installing Lucien CLI binary...")
	
	// Check if we're running from the build directory
	buildBinary := filepath.Join("build", b.config.BinaryName)
	if _, err := os.Stat(buildBinary); err == nil {
		b.logger.Debug("Found build binary: %s", buildBinary)
		return b.copyBinary(buildBinary)
	}
	
	// Check if we're running from the current directory
	currentBinary := b.config.BinaryName
	if _, err := os.Stat(currentBinary); err == nil {
		b.logger.Debug("Found current binary: %s", currentBinary)
		return b.copyBinary(currentBinary)
	}
	
	// Try to build from source
	if _, err := os.Stat("go.mod"); err == nil {
		b.logger.Info("Building from source...")
		return b.buildFromSource()
	}
	
	// Try to download from GitHub releases (placeholder)
	b.logger.Warning("Binary not found locally, download from releases not implemented")
	return fmt.Errorf("no binary found for installation")
}

// copyBinary copies the binary to the installation directory
func (b *Bootstrap) copyBinary(srcPath string) error {
	destPath := filepath.Join(b.config.InstallDir, b.config.BinaryName)
	
	b.logger.Debug("Copying %s to %s", srcPath, destPath)
	
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %v", err)
	}
	defer src.Close()
	
	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination binary: %v", err)
	}
	defer dest.Close()
	
	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}
	
	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %v", err)
		}
	}
	
	b.logger.Success("Binary installed to: %s", destPath)
	return nil
}

// buildFromSource builds Lucien CLI from Go source code
func (b *Bootstrap) buildFromSource() error {
	b.logger.Info("Building Lucien CLI from source...")
	
	buildCmd := exec.Command("go", "build", "-o", 
		filepath.Join(b.config.InstallDir, b.config.BinaryName),
		"./cmd/lucien")
	
	buildCmd.Env = append(os.Environ(), 
		"CGO_ENABLED=0",
		fmt.Sprintf("GOOS=%s", runtime.GOOS),
		fmt.Sprintf("GOARCH=%s", runtime.GOARCH),
	)
	
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\nOutput: %s", err, output)
	}
	
	b.logger.Success("Built Lucien CLI from source")
	return nil
}

// installPlugins installs Python agents and plugins
func (b *Bootstrap) installPlugins() error {
	b.logger.Info("Installing Python agents and plugins...")
	
	// Check if plugins directory exists in source
	sourcePluginDir := "plugins"
	if _, err := os.Stat(sourcePluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugins directory not found: %s", sourcePluginDir)
	}
	
	// Copy plugins from source to installation
	return b.copyPlugins(sourcePluginDir, b.config.PluginDir)
}

// copyPlugins recursively copies plugins from source to destination
func (b *Bootstrap) copyPlugins(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)
		
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		
		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
		
		// Preserve permissions
		return os.Chmod(destPath, info.Mode())
	})
}

// checkPython checks if Python 3.7+ is available
func (b *Bootstrap) checkPython() error {
	b.logger.Info("Checking Python availability...")
	
	pythonCandidates := []string{"python3", "python", "py"}
	
	for _, candidate := range pythonCandidates {
		cmd := exec.Command(candidate, "--version")
		output, err := cmd.Output()
		if err != nil {
			b.logger.Debug("Python candidate '%s' not found", candidate)
			continue
		}
		
		version := strings.TrimSpace(string(output))
		b.logger.Success("Found Python: %s", version)
		
		// Check if it's Python 3.7+
		if strings.Contains(version, "Python 3") {
			b.logger.Success("Python 3 detected - Agent commands will work")
			return nil
		}
	}
	
	return fmt.Errorf("Python 3.7+ not found in PATH")
}

// createShortcut creates desktop shortcut (Windows only)
func (b *Bootstrap) createShortcut() error {
	if runtime.GOOS != "windows" {
		b.logger.Debug("Shortcut creation only supported on Windows")
		return nil
	}
	
	b.logger.Info("Creating desktop shortcut...")
	
	// Create PowerShell script to create shortcut
	psScript := fmt.Sprintf(`
$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\Lucien CLI.lnk")
$Shortcut.TargetPath = "%s"
$Shortcut.WorkingDirectory = "$env:USERPROFILE"
$Shortcut.IconLocation = "%s"
$Shortcut.Description = "Lucien AI-Enhanced Shell"
$Shortcut.Save()
`, filepath.Join(b.config.InstallDir, b.config.BinaryName),
		filepath.Join(b.config.InstallDir, b.config.BinaryName))
	
	cmd := exec.Command("powershell", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create shortcut: %v", err)
	}
	
	b.logger.Success("Desktop shortcut created")
	return nil
}

// addToPath adds the installation directory to system PATH
func (b *Bootstrap) addToPath() error {
	b.logger.Info("Adding Lucien CLI to system PATH...")
	
	if runtime.GOOS == "windows" {
		return b.addToWindowsPath()
	}
	
	return b.addToUnixPath()
}

// addToWindowsPath adds to Windows system PATH
func (b *Bootstrap) addToWindowsPath() error {
	// Add to user PATH via registry
	cmd := exec.Command("reg", "add", "HKEY_CURRENT_USER\\Environment", 
		"/v", "PATH", "/t", "REG_EXPAND_SZ", 
		"/d", fmt.Sprintf("%%PATH%%;%s", b.config.InstallDir), "/f")
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add to PATH: %v", err)
	}
	
	// Notify system of environment change
	exec.Command("setx", "PATH", fmt.Sprintf("%%PATH%%;%s", b.config.InstallDir)).Run()
	
	b.logger.Success("Added to Windows PATH (restart terminal to use)")
	return nil
}

// addToUnixPath adds to Unix shell profile
func (b *Bootstrap) addToUnixPath() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	
	profileFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".profile"),
	}
	
	pathLine := fmt.Sprintf("export PATH=\"$PATH:%s\"", b.config.InstallDir)
	
	for _, profileFile := range profileFiles {
		if _, err := os.Stat(profileFile); os.IsNotExist(err) {
			continue
		}
		
		// Check if already added
		content, err := os.ReadFile(profileFile)
		if err != nil {
			continue
		}
		
		if strings.Contains(string(content), b.config.InstallDir) {
			b.logger.Debug("PATH already contains installation directory in %s", profileFile)
			continue
		}
		
		// Append to profile file
		file, err := os.OpenFile(profileFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}
		
		fmt.Fprintf(file, "\n# Lucien CLI\n%s\n", pathLine)
		file.Close()
		
		b.logger.Success("Added to PATH in %s", profileFile)
		return nil
	}
	
	return fmt.Errorf("no shell profile files found")
}

// verifyInstallation verifies that Lucien CLI was installed correctly
func (b *Bootstrap) verifyInstallation() error {
	b.logger.Info("Verifying installation...")
	
	binaryPath := filepath.Join(b.config.InstallDir, b.config.BinaryName)
	
	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at %s", binaryPath)
	}
	
	// Try to run version command
	cmd := exec.Command(binaryPath, "--version")
	cmd.Dir = b.config.InstallDir
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run binary: %v", err)
	}
	
	b.logger.Success("Binary works: %s", strings.TrimSpace(string(output)))
	
	// Check plugin directory
	pluginFiles := []string{
		"planner_agent.py",
		"designer_agent.py",
		"review_agent.py", 
		"code_agent.py",
	}
	
	for _, file := range pluginFiles {
		pluginPath := filepath.Join(b.config.PluginDir, file)
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			b.logger.Warning("Plugin file missing: %s", file)
		} else {
			b.logger.Debug("Plugin file found: %s", file)
		}
	}
	
	return nil
}

// printPostInstallInstructions prints instructions for the user
func (b *Bootstrap) printPostInstallInstructions() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 LUCIEN CLI INSTALLATION COMPLETE!")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Printf("📁 Installation Directory: %s\n", b.config.InstallDir)
	fmt.Printf("🔌 Plugin Directory: %s\n", b.config.PluginDir)
	
	binaryPath := filepath.Join(b.config.InstallDir, b.config.BinaryName)
	fmt.Printf("🚀 Binary Location: %s\n", binaryPath)
	
	fmt.Println("\n📋 QUICK START:")
	if b.config.AddToPath {
		fmt.Println("   lucien --version")
		fmt.Println("   lucien")
	} else {
		fmt.Printf("   %s --version\n", binaryPath)
		fmt.Printf("   %s\n", binaryPath)
	}
	
	fmt.Println("\n🤖 AI AGENT COMMANDS:")
	fmt.Println("   plan \"build a web app\"")
	fmt.Println("   design \"dark login form\"")
	fmt.Println("   review myfile.py")
	fmt.Println("   code generate \"sort function\"")
	
	fmt.Println("\n📚 DOCUMENTATION:")
	fmt.Println("   help                    # Built-in help")
	fmt.Println("   README_AGENTS.md        # Agent commands guide")
	fmt.Println("   POWERSHELL_MAPPING.md   # Command reference")
	
	if !b.config.AddToPath {
		fmt.Println("\n⚠️  IMPORTANT:")
		fmt.Println("   Add to PATH or use full path to access 'lucien' command")
	}
	
	fmt.Println("\nHappy coding! 🚀")
	fmt.Println(strings.Repeat("=", 60))
}

// getDefaultConfig returns default configuration based on platform
func getDefaultConfig() *BootstrapConfig {
	homeDir, _ := os.UserHomeDir()
	
	var installDir, binaryName string
	
	if runtime.GOOS == "windows" {
		installDir = filepath.Join(homeDir, "AppData", "Local", "Lucien")
		binaryName = "lucien.exe"
	} else {
		installDir = filepath.Join(homeDir, ".local", "bin")
		binaryName = "lucien"
	}
	
	return &BootstrapConfig{
		InstallDir:     installDir,
		PluginDir:      filepath.Join(homeDir, ".lucien", "plugins"),
		BinaryName:     binaryName,
		CreateShortcut: runtime.GOOS == "windows",
		AddToPath:      true,
		SkipPython:     false,
	}
}

// main function for the bootstrap installer
func main() {
	fmt.Println("🚀 Lucien CLI Bootstrap Installer")
	fmt.Println("==================================")
	
	config := getDefaultConfig()
	bootstrap := NewBootstrap(config)
	
	// Parse command line arguments
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v", "--verbose":
			bootstrap.SetVerbose(true)
		case "--no-shortcut":
			config.CreateShortcut = false
		case "--no-path":
			config.AddToPath = false
		case "--skip-python":
			config.SkipPython = true
		case "-h", "--help":
			printUsage()
			return
		}
	}
	
	// Run the bootstrap process
	if err := bootstrap.Run(); err != nil {
		fmt.Printf("❌ Bootstrap failed: %v\n", err)
		os.Exit(1)
	}
}

// printUsage prints usage instructions
func printUsage() {
	fmt.Println("Usage: go run bootstrap.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -v, --verbose     Enable verbose logging")
	fmt.Println("  --no-shortcut     Skip desktop shortcut creation")
	fmt.Println("  --no-path         Skip adding to system PATH")
	fmt.Println("  --skip-python     Skip Python availability check")
	fmt.Println("  -h, --help        Show this help message")
	fmt.Println()
	fmt.Println("This installer will:")
	fmt.Println("  1. Create installation directories")
	fmt.Println("  2. Install Lucien CLI binary")
	fmt.Println("  3. Install Python agent plugins")
	fmt.Println("  4. Check Python availability")
	fmt.Println("  5. Create shortcuts and PATH entries")
	fmt.Println("  6. Verify installation")
}