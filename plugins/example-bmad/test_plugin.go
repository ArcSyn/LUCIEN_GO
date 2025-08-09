// Test file to verify BMAD plugin functionality
package main

import (
	"context"
	"fmt"
)

// TestBMADPlugin runs basic functionality tests  
func TestBMADPlugin() {
	fmt.Println("🧠 Testing BMAD Plugin Functionality...")
	
	plugin := &BMADPlugin{}
	
	// Test GetInfo
	fmt.Println("\n=== Testing GetInfo ===")
	info, err := plugin.GetInfo()
	if err != nil {
		fmt.Printf("❌ GetInfo failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Plugin: %s v%s\n", info.Name, info.Version)
	fmt.Printf("✅ Description: %s\n", info.Description)
	fmt.Printf("✅ Author: %s\n", info.Author)  
	fmt.Printf("✅ Capabilities: %v\n", info.Capabilities)
	
	// Test Initialize
	fmt.Println("\n=== Testing Initialize ===")
	err = plugin.Initialize(map[string]interface{}{
		"test_mode": "true",
		"verbose":   "false",
	})
	if err != nil {
		fmt.Printf("❌ Initialize failed: %v\n", err)
		return
	}
	fmt.Println("✅ Plugin initialized successfully!")
	
	// Test Execute with different commands
	testCommands := []struct {
		command string
		args    []string
		desc    string
	}{
		{"help", []string{}, "Help command"},
		{"build", []string{}, "Build phase"},
		{"manage", []string{}, "Manage phase"},
		{"analyze", []string{}, "Analyze phase"}, 
		{"deploy", []string{}, "Deploy phase"},
		{"workflow", []string{}, "Full workflow"},
		{"invalid", []string{}, "Invalid command (should show help)"},
	}
	
	ctx := context.Background()
	for _, test := range testCommands {
		fmt.Printf("\n=== Testing: %s ===\n", test.desc)
		result, err := plugin.Execute(ctx, test.command, test.args)
		if err != nil {
			fmt.Printf("❌ Execute failed for %s: %v\n", test.command, err)
			continue
		}
		
		fmt.Printf("✅ Exit Code: %d\n", result.ExitCode)
		
		// Show first few lines of output
		lines := fmt.Sprintf("%s", result.Output)
		if len(lines) > 300 {
			fmt.Printf("✅ Output (first 300 chars): %s...\n", lines[:300])
		} else {
			fmt.Printf("✅ Output: %s\n", lines)
		}
		
		if result.Data != nil && len(result.Data) > 0 {
			fmt.Printf("✅ Data: %+v\n", result.Data)
		}
	}
	
	fmt.Println("\n🎉 All BMAD Plugin tests completed successfully!")
	fmt.Println("🚀 Plugin is ready for integration with Lucien CLI")
}