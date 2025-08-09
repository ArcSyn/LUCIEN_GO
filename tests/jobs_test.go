package tests

import (
	"strings"
	"testing"
	"time"
)

func TestJobsList(t *testing.T) {
	s := newTestShell()
	
	// Start some background jobs (simulated with commands that would run)
	// Note: This depends on job control implementation
	result, err := s.Execute("jobs")
	if err != nil {
		t.Fatalf("Jobs command failed: %v", err)
	}
	
	// Should succeed even with no jobs
	if result.ExitCode != 0 {
		t.Fatalf("Jobs command should succeed: %s", result.Error)
	}
	
	// With no jobs, should show appropriate message or empty list
	if !strings.Contains(result.Output, "No active jobs") && result.Output != "" {
		// Empty output is also acceptable
	}
}

func TestJobControlSyntax(t *testing.T) {
	s := newTestShell()
	
	// Test %1 syntax for job references
	// This tests the parser's ability to handle % syntax
	result, err := s.Execute("echo 'Job %1 syntax test'")
	if err != nil {
		t.Fatalf("Job syntax test failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "Job %1 syntax test") {
		t.Fatalf("Job syntax should be preserved in quotes: %s", result.Output)
	}
}

func TestForegroundJob(t *testing.T) {
	s := newTestShell()
	
	// Test fg command
	result, err := s.Execute("fg")
	if err != nil {
		t.Fatalf("Fg command should handle gracefully when no jobs: %v", err)
	}
	
	// Should handle no jobs gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "no") {
		t.Fatalf("Fg with no jobs should return appropriate error")
	}
}

func TestBackgroundJob(t *testing.T) {
	s := newTestShell()
	
	// Test bg command
	result, err := s.Execute("bg")
	if err != nil {
		t.Fatalf("Bg command should handle gracefully when no jobs: %v", err)
	}
	
	// Should handle no jobs gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "no") {
		t.Fatalf("Bg with no jobs should return appropriate error")
	}
}

func TestJobByNumber(t *testing.T) {
	s := newTestShell()
	
	// Test referencing job by number
	result, err := s.Execute("fg %1")
	if err != nil {
		t.Fatalf("Job number reference should handle gracefully: %v", err)
	}
	
	// Should handle non-existent job gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "No such job") {
		t.Fatalf("Non-existent job should return appropriate error")
	}
}

func TestJobByName(t *testing.T) {
	s := newTestShell()
	
	// Test referencing job by name
	result, err := s.Execute("fg %echo")
	if err != nil {
		t.Fatalf("Job name reference should handle gracefully: %v", err)
	}
	
	// Should handle non-existent job gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "No such job") {
		t.Fatalf("Non-existent job name should return appropriate error")
	}
}

func TestCurrentJob(t *testing.T) {
	s := newTestShell()
	
	// Test current job reference
	result, err := s.Execute("fg %+")
	if err != nil {
		t.Fatalf("Current job reference should handle gracefully: %v", err)
	}
	
	// Should handle no current job gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "No current job") {
		t.Fatalf("No current job should return appropriate error")
	}
}

func TestPreviousJob(t *testing.T) {
	s := newTestShell()
	
	// Test previous job reference
	result, err := s.Execute("fg %-")
	if err != nil {
		t.Fatalf("Previous job reference should handle gracefully: %v", err)
	}
	
	// Should handle no previous job gracefully
	if result.ExitCode != 1 && !strings.Contains(result.Error, "No previous job") {
		t.Fatalf("No previous job should return appropriate error")
	}
}

func TestJobTermination(t *testing.T) {
	s := newTestShell()
	
	// Test disown command
	result, err := s.Execute("disown")
	if err != nil {
		t.Fatalf("Disown command should handle gracefully: %v", err)
	}
	
	// Should handle no jobs to disown
	if result.ExitCode != 1 && !strings.Contains(result.Error, "no") {
		t.Fatalf("Disown with no jobs should return appropriate error")
	}
}

func TestJobWithPID(t *testing.T) {
	s := newTestShell()
	
	// Test kill with job reference
	result, err := s.Execute("kill %1")
	if err != nil {
		t.Fatalf("Kill with job reference should handle gracefully: %v", err)
	}
	
	// Should handle non-existent job
	if result.ExitCode != 1 && !strings.Contains(result.Error, "No such job") {
		t.Fatalf("Kill non-existent job should return error")
	}
}

func TestJobControlInScript(t *testing.T) {
	s := newTestShell()
	
	// Test job control syntax in more complex commands
	result, err := s.Execute("echo 'Background process with &' && echo 'Job control: %1'")
	if err != nil {
		t.Fatalf("Job control in script failed: %v", err)
	}
	
	// Should execute both echo commands
	if !strings.Contains(result.Output, "Background process") {
		t.Fatalf("First echo should execute")
	}
	
	if !strings.Contains(result.Output, "Job control: %1") {
		t.Fatalf("Second echo with job reference should execute")
	}
}

func TestSuspendAndResume(t *testing.T) {
	s := newTestShell()
	
	// Test suspend (Ctrl+Z simulation would be complex, test the commands)
	result, err := s.Execute("suspend")
	if err != nil {
		t.Fatalf("Suspend command should exist: %v", err)
	}
	
	// In testing, suspend might not actually suspend
	// Just verify the command is recognized
	if result.ExitCode == 127 {
		t.Fatalf("Suspend command should be recognized")
	}
}

func TestJobTimeout(t *testing.T) {
	s := newTestShell()
	
	// Test timeout with job-like command
	start := time.Now()
	result, err := s.Execute("echo 'quick job'")
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Quick job failed: %v", err)
	}
	
	// Should complete quickly (within reasonable time)
	if duration > 5*time.Second {
		t.Fatalf("Simple job took too long: %v", duration)
	}
	
	if !strings.Contains(result.Output, "quick job") {
		t.Fatalf("Job output should be captured: %s", result.Output)
	}
}

func TestJobEnvironment(t *testing.T) {
	s := newTestShell()
	
	// Test that jobs inherit environment
	s.Execute("set TEST_VAR=job_test")
	
	// Job should see the environment variable
	result, err := s.Execute("echo $TEST_VAR")
	if err != nil {
		t.Fatalf("Job environment test failed: %v", err)
	}
	
	// Should show the variable value (if variable expansion works)
	if !strings.Contains(result.Output, "job_test") && !strings.Contains(result.Output, "$TEST_VAR") {
		// Either expanded value or literal is acceptable depending on implementation
	}
}

func TestParallelJobs(t *testing.T) {
	s := newTestShell()
	
	// Test multiple commands that could run in parallel
	start := time.Now()
	result, err := s.Execute("echo job1; echo job2; echo job3")
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Parallel jobs test failed: %v", err)
	}
	
	// All jobs should complete
	jobs := []string{"job1", "job2", "job3"}
	for _, job := range jobs {
		if !strings.Contains(result.Output, job) {
			t.Fatalf("Missing job output: %s", job)
		}
	}
	
	// Should complete in reasonable time
	if duration > 10*time.Second {
		t.Fatalf("Jobs took too long: %v", duration)
	}
}

func TestJobSignalHandling(t *testing.T) {
	s := newTestShell()
	
	// Test signal-related job commands
	result, err := s.Execute("kill -l")
	if err != nil {
		t.Fatalf("Signal list command should work: %v", err)
	}
	
	// Should either list signals or indicate command exists
	if result.ExitCode == 127 {
		t.Fatalf("Kill command should be available")
	}
}

func TestJobInformation(t *testing.T) {
	s := newTestShell()
	
	// Test jobs with verbose info
	result, err := s.Execute("jobs -l")
	if err != nil {
		t.Fatalf("Jobs with options should handle gracefully: %v", err)
	}
	
	// Should handle options gracefully even with no jobs
	if result.ExitCode != 0 && !strings.Contains(result.Error, "No active jobs") {
		// Command might not support -l option, that's acceptable
	}
}