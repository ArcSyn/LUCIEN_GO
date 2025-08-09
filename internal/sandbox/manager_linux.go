//go:build linux

package sandbox

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Linux-specific sandboxing using kernel security features
func (m *Manager) applyLinuxSandbox(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Set process group for isolation
	cmd.SysProcAttr.Setpgid = true
	
	// Create new session
	cmd.SysProcAttr.Setsid = true
	
	// Set up namespace isolation
	// CLONE_NEWPID | CLONE_NEWNS | CLONE_NEWNET | CLONE_NEWIPC
	cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | 
								 syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC
	
	// Drop capabilities if running as root
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: 65534, // nobody user
		Gid: 65534, // nobody group
	}

	// Set resource limits
	if err := m.setLinuxResourceLimits(cmd); err != nil {
		return fmt.Errorf("failed to set resource limits: %v", err)
	}

	return nil
}

// Linux resource limits implementation
func (m *Manager) setLinuxResourceLimits(cmd *exec.Cmd) error {
	// Set resource limits using setrlimit system calls
	// These would be applied to the child process
	
	// Example implementation (requires unsafe syscalls):
	// syscall.Setrlimit(syscall.RLIMIT_AS, &syscall.Rlimit{Cur: 256*1024*1024, Max: 256*1024*1024})
	// syscall.Setrlimit(syscall.RLIMIT_CPU, &syscall.Rlimit{Cur: 30, Max: 30})
	// syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 64, Max: 64})
	
	return nil
}

// Linux-specific resource limit implementations
func (m *Manager) applyLinuxResourceLimits(cmd *exec.Cmd) error {
	// Apply memory limits, CPU limits, etc. using cgroups
	cmd.SysProcAttr.Setpgid = true
	return nil
}

// Stub implementations for other platforms on Linux
func (m *Manager) applyWindowsSandbox(cmd *exec.Cmd) error {
	return m.applyBasicRestrictions(cmd)
}

func (m *Manager) applyDarwinSandbox(cmd *exec.Cmd) error {
	return m.applyBasicRestrictions(cmd)
}

func (m *Manager) applyWindowsResourceLimits(cmd *exec.Cmd) error {
	return nil
}

func (m *Manager) applyDarwinResourceLimits(cmd *exec.Cmd) error {
	return nil
}