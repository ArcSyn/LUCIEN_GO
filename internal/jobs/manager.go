package jobs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// JobState represents the state of a background job
type JobState int

const (
	JobRunning JobState = iota
	JobStopped
	JobCompleted
	JobFailed
	JobKilled
)

func (s JobState) String() string {
	switch s {
	case JobRunning:
		return "Running"
	case JobStopped:
		return "Stopped"
	case JobCompleted:
		return "Completed"
	case JobFailed:
		return "Failed"
	case JobKilled:
		return "Killed"
	default:
		return "Unknown"
	}
}

// Job represents a background job with process tracking
type Job struct {
	ID          int                `json:"id"`
	Command     string             `json:"command"`
	Args        []string           `json:"args"`
	PID         int                `json:"pid"`
	State       JobState           `json:"state"`
	StartTime   time.Time          `json:"start_time"`
	EndTime     *time.Time         `json:"end_time,omitempty"`
	ExitCode    int                `json:"exit_code"`
	Output      string             `json:"output"`
	Error       string             `json:"error"`
	WorkingDir  string             `json:"working_dir"`
	Environment map[string]string  `json:"environment"`
	
	// Internal process management
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
	mu     sync.RWMutex
}

// Manager handles background job lifecycle and monitoring
type Manager struct {
	jobs       map[int]*Job
	nextJobID  int
	mu         sync.RWMutex
	monitoring bool
	stopChan   chan struct{}
}

// New creates a new job manager
func New() *Manager {
	manager := &Manager{
		jobs:      make(map[int]*Job),
		nextJobID: 1,
		stopChan:  make(chan struct{}),
	}
	
	// Start background monitoring
	go manager.monitorJobs()
	manager.monitoring = true
	
	return manager
}

// StartJob starts a new background job
func (m *Manager) StartJob(command string, args []string, workingDir string, env map[string]string) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create context for job cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create the job
	job := &Job{
		ID:          m.nextJobID,
		Command:     command,
		Args:        args,
		State:       JobRunning,
		StartTime:   time.Now(),
		WorkingDir:  workingDir,
		Environment: env,
		ctx:         ctx,
		cancel:      cancel,
		done:        make(chan struct{}),
	}
	
	m.nextJobID++
	
	// Create the command
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = workingDir
	
	// Set environment
	if env != nil {
		envSlice := os.Environ()
		for key, value := range env {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = envSlice
	}
	
	// Start the process
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start job: %w", err)
	}
	
	job.cmd = cmd
	job.PID = cmd.Process.Pid
	m.jobs[job.ID] = job
	
	// Start monitoring this job
	go m.monitorJob(job)
	
	return job, nil
}

// monitorJob monitors a single job's lifecycle
func (m *Manager) monitorJob(job *Job) {
	defer close(job.done)
	
	// Wait for the process to complete
	err := job.cmd.Wait()
	
	job.mu.Lock()
	defer job.mu.Unlock()
	
	endTime := time.Now()
	job.EndTime = &endTime
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			job.ExitCode = exitError.ExitCode()
			if job.ExitCode != 0 {
				job.State = JobFailed
			} else {
				job.State = JobCompleted
			}
		} else {
			job.State = JobFailed
			job.Error = err.Error()
		}
	} else {
		job.State = JobCompleted
		job.ExitCode = 0
	}
	
	// Handle context cancellation
	select {
	case <-job.ctx.Done():
		if job.State == JobRunning {
			job.State = JobKilled
		}
	default:
	}
}

// monitorJobs runs the background monitoring loop
func (m *Manager) monitorJobs() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.cleanupCompletedJobs()
		case <-m.stopChan:
			return
		}
	}
}

// cleanupCompletedJobs removes old completed jobs to prevent memory leaks
func (m *Manager) cleanupCompletedJobs() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	maxAge := 24 * time.Hour
	now := time.Now()
	
	for id, job := range m.jobs {
		job.mu.RLock()
		shouldCleanup := job.State != JobRunning && 
			job.EndTime != nil && 
			now.Sub(*job.EndTime) > maxAge
		job.mu.RUnlock()
		
		if shouldCleanup {
			delete(m.jobs, id)
		}
	}
}

// GetJob returns a job by ID
func (m *Manager) GetJob(id int) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	job, exists := m.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job %d not found", id)
	}
	
	return job, nil
}

// ListJobs returns all active jobs
func (m *Manager) ListJobs() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	jobs := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	
	return jobs
}

// KillJob kills a job by sending SIGTERM, then SIGKILL if necessary
func (m *Manager) KillJob(id int) error {
	m.mu.RLock()
	job, exists := m.jobs[id]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("job %d not found", id)
	}
	
	job.mu.Lock()
	defer job.mu.Unlock()
	
	if job.State != JobRunning {
		return fmt.Errorf("job %d is not running (state: %s)", id, job.State)
	}
	
	// First try SIGTERM for graceful shutdown
	if err := job.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		// If SIGTERM fails, try SIGKILL
		if killErr := job.cmd.Process.Kill(); killErr != nil {
			return fmt.Errorf("failed to kill job %d: %w", id, killErr)
		}
	} else {
		// Give the process a chance to terminate gracefully
		select {
		case <-job.done:
			// Process terminated gracefully
		case <-time.After(5 * time.Second):
			// Force kill after timeout
			job.cmd.Process.Kill()
		}
	}
	
	// Cancel the context to stop monitoring
	job.cancel()
	job.State = JobKilled
	
	return nil
}

// SuspendJob suspends a job by sending SIGSTOP
func (m *Manager) SuspendJob(id int) error {
	m.mu.RLock()
	job, exists := m.jobs[id]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("job %d not found", id)
	}
	
	job.mu.Lock()
	defer job.mu.Unlock()
	
	if job.State != JobRunning {
		return fmt.Errorf("job %d is not running (state: %s)", id, job.State)
	}
	
	// On Windows, we can't suspend processes with signals like Unix
	if runtime.GOOS == "windows" {
		// Windows doesn't have SIGSTOP/SIGCONT, so we simulate suspension
		// by keeping track of the state without actually suspending
		job.State = JobStopped
		return nil
	} else {
		// Unix-like systems can use SIGSTOP (signal 19)
		if err := job.cmd.Process.Signal(syscall.Signal(19)); err != nil {
			return fmt.Errorf("failed to suspend job %d: %w", id, err)
		}
	}
	
	job.State = JobStopped
	return nil
}

// ResumeJob resumes a suspended job by sending SIGCONT
func (m *Manager) ResumeJob(id int) error {
	m.mu.RLock()
	job, exists := m.jobs[id]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("job %d not found", id)
	}
	
	job.mu.Lock()
	defer job.mu.Unlock()
	
	if job.State != JobStopped {
		return fmt.Errorf("job %d is not stopped (state: %s)", id, job.State)
	}
	
	// On Windows, we can't resume processes with signals like Unix
	if runtime.GOOS == "windows" {
		// Windows doesn't have SIGSTOP/SIGCONT, so we simulate resumption
		job.State = JobRunning
		return nil
	} else {
		// Unix-like systems can use SIGCONT (signal 18)
		if err := job.cmd.Process.Signal(syscall.Signal(18)); err != nil {
			return fmt.Errorf("failed to resume job %d: %w", id, err)
		}
	}
	
	job.State = JobRunning
	return nil
}

// WaitForJob waits for a job to complete with timeout
func (m *Manager) WaitForJob(id int, timeout time.Duration) error {
	job, err := m.GetJob(id)
	if err != nil {
		return err
	}
	
	if timeout > 0 {
		select {
		case <-job.done:
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("job %d did not complete within %v", id, timeout)
		}
	} else {
		<-job.done
		return nil
	}
}

// GetJobStats returns statistics about jobs
func (m *Manager) GetJobStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total_jobs": len(m.jobs),
		"running":    0,
		"stopped":    0,
		"completed":  0,
		"failed":     0,
		"killed":     0,
	}
	
	for _, job := range m.jobs {
		job.mu.RLock()
		switch job.State {
		case JobRunning:
			stats["running"] = stats["running"].(int) + 1
		case JobStopped:
			stats["stopped"] = stats["stopped"].(int) + 1
		case JobCompleted:
			stats["completed"] = stats["completed"].(int) + 1
		case JobFailed:
			stats["failed"] = stats["failed"].(int) + 1
		case JobKilled:
			stats["killed"] = stats["killed"].(int) + 1
		}
		job.mu.RUnlock()
	}
	
	return stats
}

// Shutdown gracefully shuts down the job manager
func (m *Manager) Shutdown() {
	if m.monitoring {
		close(m.stopChan)
		m.monitoring = false
	}
	
	// Kill all running jobs
	m.mu.RLock()
	runningJobs := make([]int, 0)
	for id, job := range m.jobs {
		job.mu.RLock()
		if job.State == JobRunning {
			runningJobs = append(runningJobs, id)
		}
		job.mu.RUnlock()
	}
	m.mu.RUnlock()
	
	for _, id := range runningJobs {
		m.KillJob(id)
	}
}

// IsRunning returns whether a job is currently running
func (j *Job) IsRunning() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.State == JobRunning
}

// Duration returns how long the job has been running or ran for
func (j *Job) Duration() time.Duration {
	j.mu.RLock()
	defer j.mu.RUnlock()
	
	if j.EndTime != nil {
		return j.EndTime.Sub(j.StartTime)
	}
	return time.Since(j.StartTime)
}

// String returns a string representation of the job
func (j *Job) String() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	
	duration := j.Duration()
	return fmt.Sprintf("[%d] %s %s (%s, %v)", 
		j.ID, j.State, j.Command, 
		duration.Round(time.Second), j.PID)
}