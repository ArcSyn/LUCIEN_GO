package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CustomInput is a Windows-compatible input field that doesn't mess with backslashes
type CustomInput struct {
	value       string
	cursor      int
	placeholder string
	focused     bool
	width       int
	maxLength   int
}

// NewCustomInput creates a new custom input field
func NewCustomInput() CustomInput {
	return CustomInput{
		value:       "",
		cursor:      0,
		placeholder: "Enter command...",
		focused:     true,
		width:       80,
		maxLength:   512,
	}
}

// Focus sets the input as focused
func (ci *CustomInput) Focus() {
	ci.focused = true
}

// Blur removes focus from the input
func (ci *CustomInput) Blur() {
	ci.focused = false
}

// SetValue sets the input value
func (ci *CustomInput) SetValue(value string) {
	ci.value = value
	if ci.cursor > len(ci.value) {
		ci.cursor = len(ci.value)
	}
}

// SetCursor sets the cursor position
func (ci *CustomInput) SetCursor(pos int) {
	if pos < 0 {
		ci.cursor = 0
	} else if pos > len(ci.value) {
		ci.cursor = len(ci.value)
	} else {
		ci.cursor = pos
	}
}

// SetPlaceholder sets the placeholder text
func (ci *CustomInput) SetPlaceholder(placeholder string) {
	ci.placeholder = placeholder
}

// SetWidth sets the input width
func (ci *CustomInput) SetWidth(width int) {
	ci.width = width
}

// SetMaxLength sets the maximum input length
func (ci *CustomInput) SetMaxLength(length int) {
	ci.maxLength = length
}

// Value returns the current input value
func (ci *CustomInput) Value() string {
	return ci.value
}

// Position returns the cursor position
func (ci *CustomInput) Position() int {
	return ci.cursor
}

// Update handles key events - THIS IS THE KEY FIX - no backslash processing!
func (ci *CustomInput) Update(msg tea.Msg) (CustomInput, tea.Cmd) {
	if !ci.focused {
		return *ci, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// Let parent handle quit
			return *ci, nil
		case tea.KeyEnter:
			// Let parent handle command execution
			return *ci, nil
		case tea.KeyTab:
			// Let parent handle tab completion
			return *ci, nil
		case tea.KeyEsc:
			// Let parent handle escape
			return *ci, nil
		case tea.KeyCtrlL:
			// Let parent handle clear
			return *ci, nil
		case tea.KeyBackspace:
			if ci.cursor > 0 && len(ci.value) > 0 {
				ci.value = ci.value[:ci.cursor-1] + ci.value[ci.cursor:]
				ci.cursor--
			}
		case tea.KeyDelete:
			if ci.cursor < len(ci.value) {
				ci.value = ci.value[:ci.cursor] + ci.value[ci.cursor+1:]
			}
		case tea.KeyLeft:
			if ci.cursor > 0 {
				ci.cursor--
			}
		case tea.KeyRight:
			if ci.cursor < len(ci.value) {
				ci.cursor++
			}
		case tea.KeyHome:
			ci.cursor = 0
		case tea.KeyEnd:
			ci.cursor = len(ci.value)
		case tea.KeyCtrlA:
			ci.cursor = 0
		case tea.KeyCtrlE:
			ci.cursor = len(ci.value)
		case tea.KeyCtrlU:
			// Clear line
			ci.value = ""
			ci.cursor = 0
		case tea.KeyCtrlK:
			// Clear from cursor to end
			ci.value = ci.value[:ci.cursor]
		case tea.KeyCtrlW:
			// Delete word backwards
			if ci.cursor > 0 {
				// Find the start of the current word
				start := ci.cursor - 1
				for start > 0 && ci.value[start] != ' ' {
					start--
				}
				if ci.value[start] == ' ' {
					start++
				}
				ci.value = ci.value[:start] + ci.value[ci.cursor:]
				ci.cursor = start
			}
		case tea.KeyCtrlV:
			// Handle paste - this is crucial for Windows paths!
			// Note: In a real implementation, you'd get clipboard contents
			// For now, we just let the regular character input handle it
			return *ci, nil
		case tea.KeyRunes:
			// CRITICAL: Handle regular character input INCLUDING BACKSLASHES
			// This is where Bubbletea's textinput was processing backslashes as escapes
			for _, r := range msg.Runes {
				if len(ci.value) < ci.maxLength {
					// Insert character at cursor position - NO PROCESSING OF BACKSLASHES!
					ci.value = ci.value[:ci.cursor] + string(r) + ci.value[ci.cursor:]
					ci.cursor++
				}
			}
		}
	}

	return *ci, nil
}

// View renders the input field
func (ci *CustomInput) View() string {
	var b strings.Builder
	
	if len(ci.value) == 0 && ci.focused {
		// Show placeholder
		placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
		return placeholderStyle.Render(ci.placeholder)
	}
	
	// Show actual input with cursor
	value := ci.value
	if ci.focused && ci.cursor <= len(value) {
		// Insert cursor
		cursor := "█"
		if ci.cursor == len(value) {
			value += cursor
		} else {
			value = value[:ci.cursor] + cursor + value[ci.cursor:]
		}
	}
	
	// Truncate if too long for display
	if len(value) > ci.width-2 {
		start := len(value) - ci.width + 2
		value = "…" + value[start:]
	}
	
	b.WriteString(value)
	
	return b.String()
}