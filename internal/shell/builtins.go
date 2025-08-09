package shell

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Built-in command implementations

func (s *Shell) changeDirectory(args []string) (*ExecutionResult, error) {
	var dir string
	var err error
	
	if len(args) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("cd: %v", err),
				ExitCode: 1,
			}, nil
		}
		dir = homeDir
	} else {
		// Use FirstArgAsPath to properly handle quoted paths
		dir, err = FirstArgAsPath(args)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("cd: %v", err),
				ExitCode: 1,
			}, nil
		}
	}
	
	// Handle special directory shortcuts
	if dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("cd: %v", err),
				ExitCode: 1,
			}, nil
		}
		dir = homeDir
	} else if strings.HasPrefix(dir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("cd: %v", err),
				ExitCode: 1,
			}, nil
		}
		dir = filepath.Join(homeDir, dir[2:])
	}
	
	// Clean the path (handles .., ., removes duplicate slashes, etc.)
	dir = filepath.Clean(dir)
	
	// Change directory using os.Chdir
	if err := os.Chdir(dir); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cd: %v", err),
			ExitCode: 1,
		}, nil
	}
	
	// Update shell current directory to match the actual cwd
	if newDir, err := os.Getwd(); err == nil {
		s.currentDir = newDir
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) printWorkingDirectory(args []string) (*ExecutionResult, error) {
	wd, err := os.Getwd()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to get working directory: %v", err),
			ExitCode: 1,
		}, err
	}
	
	return &ExecutionResult{
		Output:   wd + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) echo(args []string) (*ExecutionResult, error) {
	output := strings.Join(args, " ") + "\n"
	
	return &ExecutionResult{
		Output:   output,
		ExitCode: 0,
	}, nil
}

func (s *Shell) listDirectory(args []string) (*ExecutionResult, error) {
	var targetDir string
	showAll := false
	showLong := false
	
	// Parse arguments
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "a") {
				showAll = true
			}
			if strings.Contains(arg, "l") {
				showLong = true
			}
		} else {
			targetDir = arg
		}
	}
	
	if targetDir == "" {
		targetDir = "."
	}
	
	// Read directory entries
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to read directory: %v", err),
			ExitCode: 1,
		}, err
	}
	
	var output strings.Builder
	
	if showLong {
		// Long format listing
		for _, entry := range entries {
			if !showAll && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			// Format: permissions size date name
			mode := info.Mode().String()
			size := info.Size()
			modTime := info.ModTime().Format("Jan 02 15:04")
			
			fmt.Fprintf(&output, "%s %8d %s %s\n", mode, size, modTime, entry.Name())
		}
	} else {
		// Simple listing
		for _, entry := range entries {
			if !showAll && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			fmt.Fprintf(&output, "%s\n", entry.Name())
		}
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) makeDirectory(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "mkdir: missing directory name",
			ExitCode: 1,
		}, fmt.Errorf("mkdir: missing directory name")
	}
	
	createParents := false
	var dirs []string
	
	// Parse arguments
	for _, arg := range args {
		if arg == "-p" {
			createParents = true
		} else {
			dirs = append(dirs, arg)
		}
	}
	
	if len(dirs) == 0 {
		return &ExecutionResult{
			Error:    "mkdir: missing directory name",
			ExitCode: 1,
		}, fmt.Errorf("mkdir: missing directory name")
	}
	
	for _, dir := range dirs {
		var err error
		if createParents {
			err = os.MkdirAll(dir, 0755)
		} else {
			err = os.Mkdir(dir, 0755)
		}
		
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("mkdir: %v", err),
				ExitCode: 1,
			}, err
		}
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) removeFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "rm: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("rm: missing file name")
	}
	
	recursive := false
	force := false
	var files []string
	
	// Parse arguments
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "r") || strings.Contains(arg, "R") {
				recursive = true
			}
			if strings.Contains(arg, "f") {
				force = true
			}
		} else {
			files = append(files, arg)
		}
	}
	
	if len(files) == 0 {
		return &ExecutionResult{
			Error:    "rm: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("rm: missing file name")
	}
	
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			if !force {
				return &ExecutionResult{
					Error:    fmt.Sprintf("rm: %v", err),
					ExitCode: 1,
				}, err
			}
			continue
		}
		
		if info.IsDir() && !recursive {
			return &ExecutionResult{
				Error:    fmt.Sprintf("rm: %s is a directory (use -r to remove)", file),
				ExitCode: 1,
			}, fmt.Errorf("rm: %s is a directory", file)
		}
		
		if recursive && info.IsDir() {
			err = os.RemoveAll(file)
		} else {
			err = os.Remove(file)
		}
		
		if err != nil && !force {
			return &ExecutionResult{
				Error:    fmt.Sprintf("rm: %v", err),
				ExitCode: 1,
			}, err
		}
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) copyFile(args []string) (*ExecutionResult, error) {
	if len(args) < 2 {
		return &ExecutionResult{
			Error:    "cp: missing source or destination",
			ExitCode: 1,
		}, fmt.Errorf("cp: missing source or destination")
	}
	
	source := args[0]
	destination := args[1]
	
	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cp: %v", err),
			ExitCode: 1,
		}, err
	}
	
	if sourceInfo.IsDir() {
		return &ExecutionResult{
			Error:    "cp: directory copying not implemented (use cp -r)",
			ExitCode: 1,
		}, fmt.Errorf("cp: directory copying not implemented")
	}
	
	// Read source file
	data, err := os.ReadFile(source)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cp: failed to read source: %v", err),
			ExitCode: 1,
		}, err
	}
	
	// Write to destination
	err = os.WriteFile(destination, data, sourceInfo.Mode())
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cp: failed to write destination: %v", err),
			ExitCode: 1,
		}, err
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) moveFile(args []string) (*ExecutionResult, error) {
	if len(args) < 2 {
		return &ExecutionResult{
			Error:    "mv: missing source or destination",
			ExitCode: 1,
		}, fmt.Errorf("mv: missing source or destination")
	}
	
	source := args[0]
	destination := args[1]
	
	// Check if source exists
	if _, err := os.Stat(source); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("mv: %v", err),
			ExitCode: 1,
		}, err
	}
	
	// Rename/move the file
	err := os.Rename(source, destination)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("mv: %v", err),
			ExitCode: 1,
		}, err
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) catFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "cat: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("cat: missing file name")
	}
	
	var output strings.Builder
	
	for _, filename := range args {
		data, err := os.ReadFile(filename)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("cat: %v", err),
				ExitCode: 1,
			}, err
		}
		output.Write(data)
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) touchFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "touch: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("touch: missing file name")
	}
	
	for _, filename := range args {
		// Check if file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// Create new file
			file, err := os.Create(filename)
			if err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("touch: %v", err),
					ExitCode: 1,
				}, err
			}
			file.Close()
		} else {
			// Update modification time
			now := time.Now()
			err := os.Chtimes(filename, now, now)
			if err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("touch: %v", err),
					ExitCode: 1,
				}, err
			}
		}
	}
	
	return &ExecutionResult{
		Output:   "",
		ExitCode: 0,
	}, nil
}

func (s *Shell) grepFile(args []string) (*ExecutionResult, error) {
	if len(args) < 2 {
		return &ExecutionResult{
			Error:    "grep: usage: grep pattern file [file...]",
			ExitCode: 1,
		}, fmt.Errorf("grep: insufficient arguments")
	}
	
	pattern := args[0]
	files := args[1:]
	
	// Compile regex
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("grep: invalid regex: %v", err),
			ExitCode: 1,
		}, err
	}
	
	var output strings.Builder
	
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("grep: %v", err),
				ExitCode: 1,
			}, err
		}
		
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if regex.MatchString(line) {
				if len(files) > 1 {
					fmt.Fprintf(&output, "%s:%s\n", filename, line)
				} else {
					fmt.Fprintf(&output, "%s\n", line)
				}
			}
		}
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) sortFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "sort: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("sort: missing file name")
	}
	
	filename := args[0]
	data, err := os.ReadFile(filename)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("sort: %v", err),
			ExitCode: 1,
		}, err
	}
	
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	sort.Strings(lines)
	
	output := strings.Join(lines, "\n") + "\n"
	
	return &ExecutionResult{
		Output:   output,
		ExitCode: 0,
	}, nil
}

func (s *Shell) wordCount(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "wc: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("wc: missing file name")
	}
	
	countLines := true
	countWords := false
	countChars := false
	var files []string
	
	// Parse arguments
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			countLines = false // Reset defaults when options specified
			if strings.Contains(arg, "l") {
				countLines = true
			}
			if strings.Contains(arg, "w") {
				countWords = true
			}
			if strings.Contains(arg, "c") {
				countChars = true
			}
		} else {
			files = append(files, arg)
		}
	}
	
	// Default to line count if no options specified
	if !countLines && !countWords && !countChars {
		countLines = true
	}
	
	var output strings.Builder
	
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("wc: %v", err),
				ExitCode: 1,
			}, err
		}
		
		content := string(data)
		
		if countLines {
			lines := strings.Count(content, "\n")
			if !strings.HasSuffix(content, "\n") && len(content) > 0 {
				lines++ // Count last line if no trailing newline
			}
			fmt.Fprintf(&output, "%d", lines)
		}
		
		if countWords {
			words := len(strings.Fields(content))
			if countLines {
				fmt.Fprintf(&output, " %d", words)
			} else {
				fmt.Fprintf(&output, "%d", words)
			}
		}
		
		if countChars {
			chars := len(content)
			if countLines || countWords {
				fmt.Fprintf(&output, " %d", chars)
			} else {
				fmt.Fprintf(&output, "%d", chars)
			}
		}
		
		if len(files) > 1 {
			fmt.Fprintf(&output, " %s", filename)
		}
		
		fmt.Fprintf(&output, "\n")
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) headFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "head: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("head: missing file name")
	}
	
	lines := 10 // default
	var filename string
	
	// Parse arguments
	for i, arg := range args {
		if arg == "-n" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				lines = n
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if n, err := strconv.Atoi(arg[1:]); err == nil {
				lines = n
			}
		} else if !strings.HasPrefix(arg, "-") && filename == "" {
			filename = arg
		}
	}
	
	if filename == "" {
		return &ExecutionResult{
			Error:    "head: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("head: missing file name")
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("head: %v", err),
			ExitCode: 1,
		}, err
	}
	
	fileLines := strings.Split(string(data), "\n")
	if lines > len(fileLines) {
		lines = len(fileLines)
	}
	
	output := strings.Join(fileLines[:lines], "\n")
	if lines > 0 && lines <= len(fileLines) {
		output += "\n"
	}
	
	return &ExecutionResult{
		Output:   output,
		ExitCode: 0,
	}, nil
}

func (s *Shell) tailFile(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "tail: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("tail: missing file name")
	}
	
	lines := 10 // default
	var filename string
	
	// Parse arguments
	for i, arg := range args {
		if arg == "-n" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				lines = n
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if n, err := strconv.Atoi(arg[1:]); err == nil {
				lines = n
			}
		} else if !strings.HasPrefix(arg, "-") && filename == "" {
			filename = arg
		}
	}
	
	if filename == "" {
		return &ExecutionResult{
			Error:    "tail: missing file name",
			ExitCode: 1,
		}, fmt.Errorf("tail: missing file name")
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("tail: %v", err),
			ExitCode: 1,
		}, err
	}
	
	fileLines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	start := len(fileLines) - lines
	if start < 0 {
		start = 0
	}
	
	output := strings.Join(fileLines[start:], "\n") + "\n"
	
	return &ExecutionResult{
		Output:   output,
		ExitCode: 0,
	}, nil
}

func (s *Shell) findFiles(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "find: missing path",
			ExitCode: 1,
		}, fmt.Errorf("find: missing path")
	}
	
	startPath := args[0]
	pattern := ""
	
	// Parse arguments for -name pattern
	for i, arg := range args {
		if arg == "-name" && i+1 < len(args) {
			pattern = args[i+1]
			break
		}
	}
	
	var output strings.Builder
	
	err := filepath.WalkDir(startPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue on errors
		}
		
		if pattern == "" {
			fmt.Fprintf(&output, "%s\n", path)
		} else {
			matched, _ := filepath.Match(pattern, d.Name())
			if matched {
				fmt.Fprintf(&output, "%s\n", path)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("find: %v", err),
			ExitCode: 1,
		}, err
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

// Register additional built-ins
func (s *Shell) registerAdditionalBuiltins() {
	s.builtins["ls"] = s.listDirectory
	s.builtins["dir"] = s.listDirectory
	s.builtins["mkdir"] = s.makeDirectory
	s.builtins["rm"] = s.removeFile
	s.builtins["del"] = s.removeFile
	s.builtins["cp"] = s.copyFile
	s.builtins["copy"] = s.copyFile
	s.builtins["mv"] = s.moveFile
	s.builtins["move"] = s.moveFile
	s.builtins["cat"] = s.catFile
	s.builtins["type"] = s.catFile
	s.builtins["touch"] = s.touchFile
	s.builtins["grep"] = s.grepFile
	s.builtins["findstr"] = s.grepFile
	s.builtins["sort"] = s.sortFile
	s.builtins["wc"] = s.wordCount
	s.builtins["head"] = s.headFile
	s.builtins["tail"] = s.tailFile
	s.builtins["find"] = s.findFiles
}

// showHelp displays comprehensive help information for all available commands
func (s *Shell) showHelp(args []string) (*ExecutionResult, error) {
	var helpText strings.Builder
	
	// Main header
	helpText.WriteString("🔴 LUCIEN NEURAL INTERFACE - COMMAND REFERENCE\n")
	helpText.WriteString("============================================\n\n")
	
	// Built-in commands section
	helpText.WriteString("📟 BUILT-IN COMMANDS:\n")
	helpText.WriteString("  help              Show this help information\n")
	helpText.WriteString("  cd <dir>          Change directory\n") 
	helpText.WriteString("  pwd               Print working directory\n")
	helpText.WriteString("  ls [path]         List directory contents\n")
	helpText.WriteString("  echo <text>       Display text\n")
	helpText.WriteString("  cat <file>        Display file contents\n")
	helpText.WriteString("  mkdir <dir>       Create directory\n")
	helpText.WriteString("  rm <file>         Remove file\n")
	helpText.WriteString("  cp <src> <dst>    Copy file\n")
	helpText.WriteString("  mv <src> <dst>    Move/rename file\n")
	helpText.WriteString("  set <var=value>   Set environment variable\n")
	helpText.WriteString("  export <var>      Export environment variable\n")
	helpText.WriteString("  env               List environment variables\n")
	helpText.WriteString("  alias <name=cmd>  Create command alias\n")
	helpText.WriteString("  history           Show command history\n")
	helpText.WriteString("  jobs              List background jobs\n")
	helpText.WriteString("  exit              Exit the shell\n\n")
	
	// AI Agent commands section
	helpText.WriteString("🤖 AI AGENT COMMANDS:\n")
	helpText.WriteString("  plan \"task\"       Break down goals into actionable tasks\n")
	helpText.WriteString("  design \"idea\"     Generate UI code from descriptions\n")
	helpText.WriteString("  review <file>     Analyze code and suggest improvements\n")
	helpText.WriteString("  code \"request\"    Generate, refactor, or explain code\n\n")
	
	// Special UI commands section
	helpText.WriteString("⚡ SPECIAL COMMANDS:\n")
	helpText.WriteString("  :theme <name>     Switch visual theme (nexus, synthwave, ghost)\n")
	helpText.WriteString("  :ai <query>       Consult neural network\n")
	helpText.WriteString("  :spells           List all available AI agents\n")
	helpText.WriteString("  :weather          Show weather information [WIP]\n")
	helpText.WriteString("  :hack             Toggle glitch mode\n")
	helpText.WriteString("  :clear            Clear terminal buffer\n")
	helpText.WriteString("  :help             Show UI help reference\n\n")
	
	// External commands section
	helpText.WriteString("🚀 EXTERNAL COMMANDS:\n")
	helpText.WriteString("  Any command not found in built-ins will be executed through:\n")
	if runtime.GOOS == "windows" {
		helpText.WriteString("  • PowerShell (Windows)\n")
	} else {
		helpText.WriteString("  • Bash (Unix/Linux)\n")
	}
	helpText.WriteString("  Examples: git status, python script.py, npm install\n\n")
	
	// Features section
	helpText.WriteString("🧠 AI FEATURES:\n")
	helpText.WriteString("  • Predictive command suggestions\n")
	helpText.WriteString("  • Context-aware assistance\n") 
	helpText.WriteString("  • Neural pattern recognition\n")
	helpText.WriteString("  • Intelligent tab completion\n\n")
	
	// Security section
	helpText.WriteString("🛡️  SECURITY:\n")
	helpText.WriteString("  • OPA policy enforcement\n")
	helpText.WriteString("  • Sandboxed plugin execution\n")
	helpText.WriteString("  • Safe-mode command filtering\n\n")
	
	// Footer
	helpText.WriteString("💡 TIPS:\n")
	helpText.WriteString("  • Use TAB for intelligent completion\n")
	helpText.WriteString("  • Press CTRL+L to clear screen\n")
	helpText.WriteString("  • Press CTRL+C to quit\n")
	helpText.WriteString("  • Pipe commands with | (example: ls | grep .go)\n")
	
	return &ExecutionResult{
		Output:   helpText.String(),
		ExitCode: 0,
	}, nil
}