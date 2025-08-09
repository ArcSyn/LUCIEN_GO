package shell

import "testing"

func TestCommandChainAccessors(t *testing.T) {
	t.Run("Len method", func(t *testing.T) {
		// Test nil chain
		var chain *CommandChain
		if chain.Len() != 0 {
			t.Errorf("Expected Len() to return 0 for nil chain, got %d", chain.Len())
		}

		// Test empty chain
		emptyChain := &CommandChain{}
		if emptyChain.Len() != 0 {
			t.Errorf("Expected Len() to return 0 for empty chain, got %d", emptyChain.Len())
		}

		// Test chain with commands
		commands := []Command{
			{Name: "echo", Args: []string{"hello"}},
			{Name: "ls", Args: []string{"-la"}},
		}
		chain = &CommandChain{Commands: commands}
		if chain.Len() != 2 {
			t.Errorf("Expected Len() to return 2, got %d", chain.Len())
		}
	})

	t.Run("At method", func(t *testing.T) {
		commands := []Command{
			{Name: "echo", Args: []string{"hello"}},
			{Name: "ls", Args: []string{"-la"}},
		}
		chain := &CommandChain{Commands: commands}

		// Test valid indices
		cmd0 := chain.At(0)
		if cmd0.Name != "echo" {
			t.Errorf("Expected At(0).Name to be 'echo', got '%s'", cmd0.Name)
		}
		if len(cmd0.Args) != 1 || cmd0.Args[0] != "hello" {
			t.Errorf("Expected At(0).Args to be ['hello'], got %v", cmd0.Args)
		}

		cmd1 := chain.At(1)
		if cmd1.Name != "ls" {
			t.Errorf("Expected At(1).Name to be 'ls', got '%s'", cmd1.Name)
		}
		if len(cmd1.Args) != 1 || cmd1.Args[0] != "-la" {
			t.Errorf("Expected At(1).Args to be ['-la'], got %v", cmd1.Args)
		}
	})

	t.Run("Slice method", func(t *testing.T) {
		// Test nil chain
		var chain *CommandChain
		slice := chain.Slice()
		if slice != nil {
			t.Errorf("Expected Slice() to return nil for nil chain, got %v", slice)
		}

		// Test empty chain
		emptyChain := &CommandChain{}
		emptySlice := emptyChain.Slice()
		if len(emptySlice) != 0 {
			t.Errorf("Expected Slice() to return empty slice for empty chain, got %v", emptySlice)
		}

		// Test chain with commands
		commands := []Command{
			{Name: "echo", Args: []string{"hello"}},
			{Name: "ls", Args: []string{"-la"}},
		}
		chain = &CommandChain{Commands: commands}
		slice = chain.Slice()

		if len(slice) != 2 {
			t.Errorf("Expected Slice() to return 2 commands, got %d", len(slice))
		}

		// Verify it's a copy (modifying the returned slice shouldn't affect original)
		slice[0].Name = "modified"
		if chain.Commands[0].Name != "echo" {
			t.Error("Slice() should return a copy, but modifying it affected the original")
		}

		// Verify content is correct
		if slice[0].Name != "modified" { // This is the modified copy
			t.Error("Failed to modify the copy - this indicates a shallow copy issue")
		}
		if slice[1].Name != "ls" {
			t.Errorf("Expected slice[1].Name to be 'ls', got '%s'", slice[1].Name)
		}
	})
}