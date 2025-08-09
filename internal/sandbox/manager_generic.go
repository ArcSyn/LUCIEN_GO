//go:build !linux && !windows && !darwin

package sandbox

import (
	"os/exec"
)

// Generic implementations for unsupported platforms
func (m *Manager) applyLinuxSandbox(cmd *exec.Cmd) error {
	return m.applyBasicRestrictions(cmd)
}

func (m *Manager) applyWindowsSandbox(cmd *exec.Cmd) error {
	return m.applyBasicRestrictions(cmd)
}

func (m *Manager) applyDarwinSandbox(cmd *exec.Cmd) error {
	return m.applyBasicRestrictions(cmd)
}

func (m *Manager) applyLinuxResourceLimits(cmd *exec.Cmd) error {
	return nil
}

func (m *Manager) applyWindowsResourceLimits(cmd *exec.Cmd) error {
	return nil
}

func (m *Manager) applyDarwinResourceLimits(cmd *exec.Cmd) error {
	return nil
}