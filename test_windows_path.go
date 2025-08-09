package main

import (
	"fmt"
	"github.com/ArcSyn/LucienCLI/internal/shell"
)

func main() {
	// Test Windows path handling directly
	shellEngine := shell.New(&shell.Config{
		SafeMode: true,
	})

	// Test commands with Windows paths
	testCommands := []string{
		`cd "C:\Users"`,
		`cd C:\Windows`,
		`echo "C:\Program Files\Test"`,
		`pwd`,
	}

	for _, cmd := range testCommands {
		fmt.Printf("Testing: %s\n", cmd)
		result, err := shellEngine.Execute(cmd)
		if err != nil {
			fmt.Printf("❌ ERROR: %v\n", err)
		} else {
			fmt.Printf("✅ SUCCESS: %s\n", result.Output)
		}
		fmt.Println()
	}
}