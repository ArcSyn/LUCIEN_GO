package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// SoundSystem provides audio feedback for the Visual Bliss interface
type SoundSystem struct {
	enabled     bool
	volume      float64
	sounds      map[string]SoundEffect
	playQueue   chan SoundEvent
	mutex       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// SoundEffect represents a sound that can be played
type SoundEffect struct {
	Name        string
	Type        SoundType
	Frequency   float64
	Duration    time.Duration
	Volume      float64
	Pattern     []Note
	Description string
}

// SoundType represents different types of sound generation
type SoundType int

const (
	BeepSound SoundType = iota
	ToneSound
	ChordSound
	SequenceSound
	NoiseSound
	GenerativeSound
)

// Note represents a musical note for complex sounds
type Note struct {
	Frequency float64
	Duration  time.Duration
	Volume    float64
}

// SoundEvent represents a sound that should be played
type SoundEvent struct {
	SoundName string
	Context   string
	Timestamp time.Time
}

// NewSoundSystem creates a new sound system
func NewSoundSystem() *SoundSystem {
	ctx, cancel := context.WithCancel(context.Background())
	
	ss := &SoundSystem{
		enabled:   true,
		volume:    0.5,
		sounds:    make(map[string]SoundEffect),
		playQueue: make(chan SoundEvent, 100),
		ctx:       ctx,
		cancel:    cancel,
	}
	
	ss.initializeSounds()
	go ss.soundProcessor()
	
	return ss
}

// initializeSounds sets up all the sound effects for Visual Bliss
func (ss *SoundSystem) initializeSounds() {
	// Success sounds
	ss.sounds["success"] = SoundEffect{
		Name:        "success",
		Type:        ChordSound,
		Pattern: []Note{
			{Frequency: 523.25, Duration: 100 * time.Millisecond, Volume: 0.3}, // C5
			{Frequency: 659.25, Duration: 100 * time.Millisecond, Volume: 0.3}, // E5
			{Frequency: 783.99, Duration: 200 * time.Millisecond, Volume: 0.4}, // G5
		},
		Description: "Command completed successfully",
	}
	
	// Error sounds
	ss.sounds["error"] = SoundEffect{
		Name:        "error",
		Type:        SequenceSound,
		Pattern: []Note{
			{Frequency: 200.0, Duration: 150 * time.Millisecond, Volume: 0.4},
			{Frequency: 180.0, Duration: 150 * time.Millisecond, Volume: 0.3},
			{Frequency: 160.0, Duration: 200 * time.Millisecond, Volume: 0.2},
		},
		Description: "Command failed or error occurred",
	}
	
	// Startup/awakening sounds
	ss.sounds["awaken"] = SoundEffect{
		Name:        "awaken",
		Type:        SequenceSound,
		Pattern: []Note{
			{Frequency: 220.0, Duration: 300 * time.Millisecond, Volume: 0.2},  // A3
			{Frequency: 293.66, Duration: 300 * time.Millisecond, Volume: 0.3}, // D4
			{Frequency: 369.99, Duration: 300 * time.Millisecond, Volume: 0.4}, // F#4
			{Frequency: 440.0, Duration: 400 * time.Millisecond, Volume: 0.5},  // A4
			{Frequency: 554.37, Duration: 500 * time.Millisecond, Volume: 0.4}, // C#5
		},
		Description: "Neural awakening sequence",
	}
	
	// Notification sounds
	ss.sounds["notification"] = SoundEffect{
		Name:        "notification",
		Type:        ToneSound,
		Frequency:   800.0,
		Duration:    200 * time.Millisecond,
		Volume:      0.3,
		Description: "General notification",
	}
	
	// Idle/ambient sounds
	ss.sounds["idle"] = SoundEffect{
		Name:        "idle",
		Type:        GenerativeSound,
		Frequency:   60.0, // Base frequency for generative music
		Duration:    5 * time.Second,
		Volume:      0.1,
		Description: "Ambient idle sounds",
	}
	
	// Command input sounds
	ss.sounds["input"] = SoundEffect{
		Name:        "input",
		Type:        BeepSound,
		Frequency:   1000.0,
		Duration:    50 * time.Millisecond,
		Volume:      0.1,
		Description: "Key input feedback",
	}
	
	// Tab completion sound
	ss.sounds["complete"] = SoundEffect{
		Name:        "complete",
		Type:        ToneSound,
		Frequency:   1200.0,
		Duration:    100 * time.Millisecond,
		Volume:      0.2,
		Description: "Tab completion",
	}
	
	// Theme change sound
	ss.sounds["theme_change"] = SoundEffect{
		Name:        "theme_change",
		Type:        ChordSound,
		Pattern: []Note{
			{Frequency: 440.0, Duration: 150 * time.Millisecond, Volume: 0.3},  // A4
			{Frequency: 554.37, Duration: 150 * time.Millisecond, Volume: 0.3}, // C#5
			{Frequency: 659.25, Duration: 200 * time.Millisecond, Volume: 0.4}, // E5
		},
		Description: "Visual theme changed",
	}
	
	// AI thinking sound
	ss.sounds["ai_thinking"] = SoundEffect{
		Name:        "ai_thinking",
		Type:        GenerativeSound,
		Frequency:   300.0,
		Duration:    2 * time.Second,
		Volume:      0.15,
		Description: "AI processing query",
	}
	
	// Vanish/stealth sound
	ss.sounds["vanish"] = SoundEffect{
		Name:        "vanish",
		Type:        SequenceSound,
		Pattern: []Note{
			{Frequency: 800.0, Duration: 100 * time.Millisecond, Volume: 0.4},
			{Frequency: 600.0, Duration: 100 * time.Millisecond, Volume: 0.3},
			{Frequency: 400.0, Duration: 100 * time.Millisecond, Volume: 0.2},
			{Frequency: 200.0, Duration: 200 * time.Millisecond, Volume: 0.1},
		},
		Description: "Phantom protocol activation",
	}
	
	// Process start sound
	ss.sounds["process_start"] = SoundEffect{
		Name:        "process_start",
		Type:        ToneSound,
		Frequency:   600.0,
		Duration:    150 * time.Millisecond,
		Volume:      0.25,
		Description: "Background process started",
	}
	
	// Glitch mode sound
	ss.sounds["glitch"] = SoundEffect{
		Name:        "glitch",
		Type:        NoiseSound,
		Frequency:   0.0, // Random noise
		Duration:    300 * time.Millisecond,
		Volume:      0.2,
		Description: "Reality glitch effect",
	}
}

// PlaySound queues a sound to be played
func (ss *SoundSystem) PlaySound(soundName, context string) {
	if !ss.enabled {
		return
	}
	
	event := SoundEvent{
		SoundName: soundName,
		Context:   context,
		Timestamp: time.Now(),
	}
	
	select {
	case ss.playQueue <- event:
	default:
		// Queue full, skip this sound
		log.Printf("🔊 Sound queue full, skipping: %s", soundName)
	}
}

// soundProcessor processes the sound queue
func (ss *SoundSystem) soundProcessor() {
	for {
		select {
		case <-ss.ctx.Done():
			return
		case event := <-ss.playQueue:
			ss.playEventSound(event)
		}
	}
}

// playEventSound plays a sound effect
func (ss *SoundSystem) playEventSound(event SoundEvent) {
	ss.mutex.RLock()
	sound, exists := ss.sounds[event.SoundName]
	ss.mutex.RUnlock()
	
	if !exists {
		log.Printf("🔊 Unknown sound: %s", event.SoundName)
		return
	}
	
	log.Printf("🔊 Playing sound: %s (context: %s)", sound.Name, event.Context)
	
	switch sound.Type {
	case BeepSound:
		ss.playBeep(sound)
	case ToneSound:
		ss.playTone(sound)
	case ChordSound:
		ss.playChord(sound)
	case SequenceSound:
		ss.playSequence(sound)
	case NoiseSound:
		ss.playNoise(sound)
	case GenerativeSound:
		ss.playGenerative(sound)
	}
}

// Platform-specific sound generation methods

// playBeep plays a simple system beep
func (ss *SoundSystem) playBeep(sound SoundEffect) {
	switch runtime.GOOS {
	case "windows":
		ss.playWindowsBeep(sound)
	case "darwin":
		ss.playMacBeep(sound)
	case "linux":
		ss.playLinuxBeep(sound)
	default:
		fmt.Print("\a") // Basic ASCII bell
	}
}

// playTone plays a pure tone
func (ss *SoundSystem) playTone(sound SoundEffect) {
	switch runtime.GOOS {
	case "windows":
		ss.playWindowsTone(sound)
	case "darwin":
		ss.playMacTone(sound)
	case "linux":
		ss.playLinuxTone(sound)
	default:
		fmt.Print("\a")
	}
}

// playChord plays multiple notes simultaneously
func (ss *SoundSystem) playChord(sound SoundEffect) {
	// Play all notes in the pattern at once
	for _, note := range sound.Pattern {
		go ss.playNote(note)
	}
}

// playSequence plays notes in sequence
func (ss *SoundSystem) playSequence(sound SoundEffect) {
	for _, note := range sound.Pattern {
		ss.playNote(note)
		time.Sleep(note.Duration)
	}
}

// playNoise plays random noise
func (ss *SoundSystem) playNoise(sound SoundEffect) {
	// Generate glitch-like noise effect
	ss.playGlitchNoise(sound)
}

// playGenerative plays generative ambient music
func (ss *SoundSystem) playGenerative(sound SoundEffect) {
	ss.playAmbientGenerative(sound)
}

// playNote plays a single musical note
func (ss *SoundSystem) playNote(note Note) {
	adjustedVolume := note.Volume * ss.volume
	
	switch runtime.GOOS {
	case "windows":
		ss.playWindowsNote(note.Frequency, note.Duration, adjustedVolume)
	case "darwin":
		ss.playMacNote(note.Frequency, note.Duration, adjustedVolume)
	case "linux":
		ss.playLinuxNote(note.Frequency, note.Duration, adjustedVolume)
	default:
		fmt.Print("\a")
	}
}

// Windows-specific sound implementations
func (ss *SoundSystem) playWindowsBeep(sound SoundEffect) {
	freq := int(sound.Frequency)
	duration := int(sound.Duration.Milliseconds())
	
	// Use PowerShell to generate sound
	cmd := exec.Command("powershell", "-c", 
		fmt.Sprintf("[console]::beep(%d,%d)", freq, duration))
	cmd.Run()
}

func (ss *SoundSystem) playWindowsTone(sound SoundEffect) {
	ss.playWindowsBeep(sound)
}

func (ss *SoundSystem) playWindowsNote(freq float64, duration time.Duration, volume float64) {
	cmd := exec.Command("powershell", "-c",
		fmt.Sprintf("[console]::beep(%d,%d)", int(freq), int(duration.Milliseconds())))
	cmd.Run()
}

// macOS-specific sound implementations
func (ss *SoundSystem) playMacBeep(sound SoundEffect) {
	// Use afplay or osascript for sound generation
	script := fmt.Sprintf(`
		set freq to %f
		set duration to %f
		do shell script "( speaker-test -t sine -f " & freq & " )& pid=$!; sleep " & duration & "; kill -9 $pid"
	`, sound.Frequency, sound.Duration.Seconds())
	
	cmd := exec.Command("osascript", "-e", script)
	cmd.Run()
}

func (ss *SoundSystem) playMacTone(sound SoundEffect) {
	ss.playMacBeep(sound)
}

func (ss *SoundSystem) playMacNote(freq float64, duration time.Duration, volume float64) {
	// Generate sine wave using Sox or similar tool if available
	if ss.hasSox() {
		cmd := exec.Command("sox", "-n", "-t", "alsa", "default", "synth",
			fmt.Sprintf("%.2f", duration.Seconds()), "sine", fmt.Sprintf("%.2f", freq),
			"vol", fmt.Sprintf("%.2f", volume))
		cmd.Run()
	} else {
		// Fallback to system beep
		cmd := exec.Command("osascript", "-e", "beep")
		cmd.Run()
	}
}

// Linux-specific sound implementations
func (ss *SoundSystem) playLinuxBeep(sound SoundEffect) {
	// Try different methods based on what's available
	if ss.hasSox() {
		ss.playLinuxSoxTone(sound)
	} else if ss.hasBeep() {
		ss.playLinuxSystemBeep(sound)
	} else if ss.hasSpeakerTest() {
		ss.playSpeakerTest(sound)
	} else {
		// Fallback to ASCII bell
		fmt.Print("\a")
	}
}

func (ss *SoundSystem) playLinuxTone(sound SoundEffect) {
	ss.playLinuxBeep(sound)
}

func (ss *SoundSystem) playLinuxNote(freq float64, duration time.Duration, volume float64) {
	if ss.hasSox() {
		cmd := exec.Command("sox", "-n", "-t", "alsa", "default", "synth",
			fmt.Sprintf("%.2f", duration.Seconds()), "sine", fmt.Sprintf("%.2f", freq),
			"vol", fmt.Sprintf("%.2f", volume))
		cmd.Run()
	} else {
		fmt.Print("\a")
	}
}

func (ss *SoundSystem) playLinuxSoxTone(sound SoundEffect) {
	cmd := exec.Command("sox", "-n", "-t", "alsa", "default", "synth",
		fmt.Sprintf("%.2f", sound.Duration.Seconds()), "sine", fmt.Sprintf("%.2f", sound.Frequency),
		"vol", fmt.Sprintf("%.2f", sound.Volume*ss.volume))
	cmd.Run()
}

func (ss *SoundSystem) playLinuxSystemBeep(sound SoundEffect) {
	cmd := exec.Command("beep", "-f", fmt.Sprintf("%.0f", sound.Frequency),
		"-l", fmt.Sprintf("%.0f", sound.Duration.Milliseconds()))
	cmd.Run()
}

func (ss *SoundSystem) playSpeakerTest(sound SoundEffect) {
	// Use speaker-test for tone generation
	ctx, cancel := context.WithTimeout(context.Background(), sound.Duration)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "speaker-test", "-t", "sine", "-f", 
		fmt.Sprintf("%.0f", sound.Frequency))
	cmd.Run()
}

// Special effects
func (ss *SoundSystem) playGlitchNoise(sound SoundEffect) {
	// Generate glitch-like sound effect
	switch runtime.GOOS {
	case "linux":
		if ss.hasSox() {
			// Generate noise with random modulation
			cmd := exec.Command("sox", "-n", "-t", "alsa", "default", "synth",
				fmt.Sprintf("%.2f", sound.Duration.Seconds()), "noise",
				"vol", fmt.Sprintf("%.2f", sound.Volume*ss.volume),
				"tremolo", "20", "0.5")
			cmd.Run()
		}
	default:
		// Fallback: rapid beeps
		for i := 0; i < 5; i++ {
			freq := 200.0 + float64(i*100)
			ss.playNote(Note{
				Frequency: freq,
				Duration:  50 * time.Millisecond,
				Volume:    sound.Volume,
			})
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (ss *SoundSystem) playAmbientGenerative(sound SoundEffect) {
	// Generate ambient music using algorithmic composition
	baseFreq := sound.Frequency
	duration := sound.Duration
	volume := sound.Volume * ss.volume
	
	// Create a simple ambient pattern
	notes := ss.generateAmbientSequence(baseFreq, duration)
	
	for _, note := range notes {
		note.Volume = volume
		go ss.playNote(note)
		time.Sleep(note.Duration / 3) // Overlap notes for ambient effect
	}
}

// generateAmbientSequence creates a sequence of notes for ambient music
func (ss *SoundSystem) generateAmbientSequence(baseFreq float64, totalDuration time.Duration) []Note {
	var notes []Note
	
	// Use pentatonic scale for pleasant ambient sounds
	pentatonicRatios := []float64{1.0, 9.0/8.0, 5.0/4.0, 3.0/2.0, 5.0/3.0}
	
	noteDuration := 800 * time.Millisecond
	noteCount := int(totalDuration / noteDuration)
	
	for i := 0; i < noteCount; i++ {
		// Select note from pentatonic scale
		ratio := pentatonicRatios[i%len(pentatonicRatios)]
		freq := baseFreq * ratio
		
		// Add some octave variation
		if i%3 == 0 {
			freq *= 2.0 // One octave up
		}
		
		notes = append(notes, Note{
			Frequency: freq,
			Duration:  noteDuration,
			Volume:    0.1 + 0.05*math.Sin(float64(i)*0.5), // Gentle volume variation
		})
	}
	
	return notes
}

// Utility methods to check for available sound tools
func (ss *SoundSystem) hasSox() bool {
	_, err := exec.LookPath("sox")
	return err == nil
}

func (ss *SoundSystem) hasBeep() bool {
	_, err := exec.LookPath("beep")
	return err == nil
}

func (ss *SoundSystem) hasSpeakerTest() bool {
	_, err := exec.LookPath("speaker-test")
	return err == nil
}

// Configuration methods
func (ss *SoundSystem) SetEnabled(enabled bool) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.enabled = enabled
}

func (ss *SoundSystem) IsEnabled() bool {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	return ss.enabled
}

func (ss *SoundSystem) SetVolume(volume float64) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	ss.volume = volume
}

func (ss *SoundSystem) GetVolume() float64 {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	return ss.volume
}

// GetAvailableSounds returns a list of all available sound effects
func (ss *SoundSystem) GetAvailableSounds() []string {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	
	var sounds []string
	for name := range ss.sounds {
		sounds = append(sounds, name)
	}
	return sounds
}

// GetSoundInfo returns information about a specific sound effect
func (ss *SoundSystem) GetSoundInfo(name string) (SoundEffect, bool) {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	
	sound, exists := ss.sounds[name]
	return sound, exists
}

// SaveSoundConfig saves the current sound configuration to a file
func (ss *SoundSystem) SaveSoundConfig(configPath string) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	
	config := struct {
		Enabled bool                       `json:"enabled"`
		Volume  float64                    `json:"volume"`
		Sounds  map[string]SoundEffect     `json:"sounds"`
	}{
		Enabled: ss.enabled,
		Volume:  ss.volume,
		Sounds:  ss.sounds,
	}
	
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Write config file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sound config: %v", err)
	}
	
	return os.WriteFile(configPath, data, 0644)
}

// LoadSoundConfig loads sound configuration from a file
func (ss *SoundSystem) LoadSoundConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read sound config: %v", err)
	}
	
	var config struct {
		Enabled bool                       `json:"enabled"`
		Volume  float64                    `json:"volume"`
		Sounds  map[string]SoundEffect     `json:"sounds"`
	}
	
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal sound config: %v", err)
	}
	
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	
	ss.enabled = config.Enabled
	ss.volume = config.Volume
	
	// Merge custom sounds with default sounds
	for name, sound := range config.Sounds {
		ss.sounds[name] = sound
	}
	
	return nil
}

// Close shuts down the sound system
func (ss *SoundSystem) Close() {
	ss.cancel()
	close(ss.playQueue)
}

// TestSound plays a test sound to verify the system is working
func (ss *SoundSystem) TestSound() {
	ss.PlaySound("notification", "sound system test")
}