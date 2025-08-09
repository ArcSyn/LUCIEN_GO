package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// System command implementations

func (s *Shell) processListing(args []string) (*ExecutionResult, error) {
	var output strings.Builder
	
	// Cross-platform process listing
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FO", "CSV")
	} else {
		cmd = exec.Command("ps", "aux")
	}
	
	result, err := cmd.Output()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("ps: failed to list processes: %v", err),
			ExitCode: 1,
		}, err
	}
	
	if runtime.GOOS == "windows" {
		// Parse CSV output from tasklist
		lines := strings.Split(string(result), "\n")
		output.WriteString("PID     NAME                     MEM\n")
		for _, line := range lines[1:] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Split(line, ",")
			if len(fields) >= 5 {
				// Remove quotes from fields
				for i := range fields {
					fields[i] = strings.Trim(fields[i], "\"")
				}
				fmt.Fprintf(&output, "%-8s %-24s %s\n", fields[1], fields[0], fields[4])
			}
		}
	} else {
		output.Write(result)
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) killProcess(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "kill: missing process ID",
			ExitCode: 1,
		}, fmt.Errorf("kill: missing process ID")
	}
	
	force := false
	var pids []int
	
	// Parse arguments
	for _, arg := range args {
		if arg == "-9" || arg == "-KILL" {
			force = true
		} else {
			pid, err := strconv.Atoi(arg)
			if err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("kill: invalid process ID: %s", arg),
					ExitCode: 1,
				}, err
			}
			pids = append(pids, pid)
		}
	}
	
	if len(pids) == 0 {
		return &ExecutionResult{
			Error:    "kill: no process ID specified",
			ExitCode: 1,
		}, fmt.Errorf("kill: no process ID specified")
	}
	
	for _, pid := range pids {
		if runtime.GOOS == "windows" {
			var cmd *exec.Cmd
			if force {
				cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
			} else {
				cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid))
			}
			
			if err := cmd.Run(); err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("kill: failed to kill process %d: %v", pid, err),
					ExitCode: 1,
				}, err
			}
		} else {
			// Unix-like systems
			process, err := os.FindProcess(pid)
			if err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("kill: process %d not found", pid),
					ExitCode: 1,
				}, err
			}
			
			signal := syscall.SIGTERM
			if force {
				signal = syscall.SIGKILL
			}
			
			if err := process.Signal(signal); err != nil {
				return &ExecutionResult{
					Error:    fmt.Sprintf("kill: failed to kill process %d: %v", pid, err),
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

func (s *Shell) diskUsage(args []string) (*ExecutionResult, error) {
	var output strings.Builder
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption", "/format:csv")
	} else {
		cmd = exec.Command("df", "-h")
	}
	
	result, err := cmd.Output()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("df: failed to get disk usage: %v", err),
			ExitCode: 1,
		}, err
	}
	
	if runtime.GOOS == "windows" {
		// Parse CSV output from wmic
		lines := strings.Split(string(result), "\n")
		output.WriteString("Drive   Size        Free        Used%\n")
		for _, line := range lines[1:] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Split(line, ",")
			if len(fields) >= 4 && fields[1] != "" {
				caption := strings.TrimSpace(fields[1])
				freeBytes := strings.TrimSpace(fields[2])
				sizeBytes := strings.TrimSpace(fields[3])
				
				if freeBytes != "" && sizeBytes != "" {
					size, _ := strconv.ParseInt(sizeBytes, 10, 64)
					free, _ := strconv.ParseInt(freeBytes, 10, 64)
					used := size - free
					usedPercent := float64(used) / float64(size) * 100
					
					sizeGB := float64(size) / (1024 * 1024 * 1024)
					freeGB := float64(free) / (1024 * 1024 * 1024)
					
					fmt.Fprintf(&output, "%-8s %8.1fG %8.1fG %6.1f%%\n", 
						caption, sizeGB, freeGB, usedPercent)
				}
			}
		}
	} else {
		output.Write(result)
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) systemUptime(args []string) (*ExecutionResult, error) {
	var output strings.Builder
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "os", "get", "lastbootuptime", "/format:csv")
	} else {
		cmd = exec.Command("uptime")
	}
	
	result, err := cmd.Output()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("uptime: failed to get system uptime: %v", err),
			ExitCode: 1,
		}, err
	}
	
	if runtime.GOOS == "windows" {
		// Parse boot time from wmic and calculate uptime
		lines := strings.Split(string(result), "\n")
		for _, line := range lines[1:] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Split(line, ",")
			if len(fields) >= 2 && fields[1] != "" {
				bootTime := strings.TrimSpace(fields[1])
				if len(bootTime) >= 14 {
					// Parse WMI date format: YYYYMMDDHHMMSS
					year := bootTime[0:4]
					month := bootTime[4:6]
					day := bootTime[6:8]
					hour := bootTime[8:10]
					minute := bootTime[10:12]
					second := bootTime[12:14]
					
					bootTimeStr := fmt.Sprintf("%s-%s-%s %s:%s:%s", 
						year, month, day, hour, minute, second)
					
					bootTimeFormatted, parseErr := time.Parse("2006-01-02 15:04:05", bootTimeStr)
					if parseErr == nil {
						uptime := time.Since(bootTimeFormatted)
						days := int(uptime.Hours() / 24)
						hours := int(uptime.Hours()) % 24
						minutes := int(uptime.Minutes()) % 60
						
						fmt.Fprintf(&output, "System up %d days, %d hours, %d minutes\n", 
							days, hours, minutes)
					} else {
						fmt.Fprintf(&output, "Boot time: %s\n", bootTime)
					}
				}
				break
			}
		}
	} else {
		output.Write(result)
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) currentUser(args []string) (*ExecutionResult, error) {
	var output strings.Builder
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("whoami")
	} else {
		cmd = exec.Command("whoami")
	}
	
	result, err := cmd.Output()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("whoami: failed to get current user: %v", err),
			ExitCode: 1,
		}, err
	}
	
	output.Write(result)
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) systemInfo(args []string) (*ExecutionResult, error) {
	var output strings.Builder
	
	// Basic system information
	fmt.Fprintf(&output, "Operating System: %s\n", runtime.GOOS)
	fmt.Fprintf(&output, "Architecture: %s\n", runtime.GOARCH)
	fmt.Fprintf(&output, "Go Version: %s\n", runtime.Version())
	fmt.Fprintf(&output, "CPU Cores: %d\n", runtime.NumCPU())
	
	// Get hostname
	hostname, err := os.Hostname()
	if err == nil {
		fmt.Fprintf(&output, "Hostname: %s\n", hostname)
	}
	
	// Get working directory
	wd, err := os.Getwd()
	if err == nil {
		fmt.Fprintf(&output, "Working Directory: %s\n", wd)
	}
	
	// Get user info
	homeDir, err := os.UserHomeDir()
	if err == nil {
		fmt.Fprintf(&output, "Home Directory: %s\n", homeDir)
	}
	
	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) networkPing(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "ping: missing target host",
			ExitCode: 1,
		}, fmt.Errorf("ping: missing target host")
	}
	
	host := args[0]
	count := 4 // default ping count
	
	// Parse arguments
	for i, arg := range args {
		if arg == "-c" && i+1 < len(args) {
			if c, err := strconv.Atoi(args[i+1]); err == nil {
				count = c
			}
		}
	}
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), host)
	} else {
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), host)
	}
	
	result, err := cmd.Output()
	if err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("ping: failed to ping %s: %v", host, err),
			ExitCode: 1,
		}, err
	}
	
	return &ExecutionResult{
		Output:   string(result),
		ExitCode: 0,
	}, nil
}

// Register system command built-ins
func (s *Shell) registerSystemBuiltins() {
	s.builtins["ps"] = s.processListing
	s.builtins["tasklist"] = s.processListing
	s.builtins["kill"] = s.killProcess  
	s.builtins["taskkill"] = s.killProcess
	s.builtins["df"] = s.diskUsage
	s.builtins["uptime"] = s.systemUptime
	s.builtins["whoami"] = s.currentUser
	s.builtins["systeminfo"] = s.systemInfo
	s.builtins["ping"] = s.networkPing
}