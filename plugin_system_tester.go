package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/luciendev/lucien-core/internal/plugin"
)

func main() {
	fmt.Println("=== LUCIEN PLUGIN SYSTEM TEST ===")
	
	// Initialize plugin manager
	pluginDir := filepath.Join(".", "plugins")
	manager := plugin.New(pluginDir)
	
	fmt.Printf("Plugin directory: %s\n", pluginDir)
	
	// Test plugin discovery and loading
	fmt.Println("\nDiscovering and loading plugins...")
	err := manager.LoadPlugins()
	if err != nil {
		fmt.Printf("Error loading plugins: %v\n", err)
	}
	
	// List loaded plugins
	fmt.Println("\nLoaded plugins:")
	plugins := manager.ListPlugins()
	if len(plugins) == 0 {
		fmt.Println("No plugins loaded")
	} else {
		for name, info := range plugins {
			fmt.Printf("- %s v%s: %s\n", name, info.Version, info.Description)
			fmt.Printf("  Author: %s\n", info.Author)
			fmt.Printf("  Capabilities: %v\n", info.Capabilities)
		}
	}
	
	// Test plugin execution if any are loaded
	if len(plugins) > 0 {
		fmt.Println("\nTesting plugin execution...")
		for name := range plugins {
			fmt.Printf("Testing plugin: %s\n", name)
			result, err := manager.ExecutePlugin(name, "test", []string{"arg1", "arg2"})
			if err != nil {
				fmt.Printf("  Error executing plugin: %v\n", err)
			} else {
				fmt.Printf("  Output: %s\n", result.Output)
				fmt.Printf("  Exit Code: %d\n", result.ExitCode)
				if result.Error != "" {
					fmt.Printf("  Error: %s\n", result.Error)
				}
			}
		}
	}
	
	// Test plugin directory structure
	fmt.Println("\nTesting plugin directory structure...")
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		fmt.Printf("Error reading plugin directory: %v\n", err)
	} else {
		for _, entry := range entries {
			if entry.IsDir() {
				pluginPath := filepath.Join(pluginDir, entry.Name())
				fmt.Printf("Plugin directory: %s\n", entry.Name())
				
				// Check for manifest
				manifestPath := filepath.Join(pluginPath, "manifest.json")
				if _, err := os.Stat(manifestPath); err == nil {
					fmt.Printf("  ✅ manifest.json found\n")
				} else {
					fmt.Printf("  ❌ manifest.json missing\n")
				}
				
				// List files in plugin directory
				pluginEntries, err := os.ReadDir(pluginPath)
				if err == nil {
					for _, pluginEntry := range pluginEntries {
						fmt.Printf("  - %s\n", pluginEntry.Name())
					}
				}
			}
		}
	}
	
	// Test error conditions
	fmt.Println("\nTesting error conditions...")
	
	// Try to execute non-existent plugin
	_, err = manager.ExecutePlugin("nonexistent", "test", []string{})
	if err != nil {
		fmt.Printf("✅ Correctly handled non-existent plugin: %v\n", err)
	} else {
		fmt.Printf("❌ Should have errored for non-existent plugin\n")
	}
	
	// Try to get info for non-existent plugin
	_, err = manager.GetPlugin("nonexistent")
	if err != nil {
		fmt.Printf("✅ Correctly handled non-existent plugin info: %v\n", err)
	} else {
		fmt.Printf("❌ Should have errored for non-existent plugin info\n")
	}
	
	// Clean up
	fmt.Println("\nCleaning up...")
	manager.UnloadAllPlugins()
	
	fmt.Println("Plugin system test complete!")
}