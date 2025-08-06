package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-plugin"
)

// Handshake configuration for plugins
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "LUCIEN_PLUGIN",
	MagicCookieValue: "lucien_neural_interface",
}

// PluginInterface defines the interface that all Lucien plugins must implement
type PluginInterface interface {
	Execute(ctx context.Context, command string, args []string) (*Result, error)
	GetInfo() (*Info, error)
	Initialize(config map[string]interface{}) error
}

// Info contains plugin metadata
type Info struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Capabilities []string         `json:"capabilities"`
	Config      map[string]string `json:"config"`
}

// Result contains plugin execution results
type Result struct {
	Output   string            `json:"output"`
	Error    string            `json:"error,omitempty"`
	ExitCode int               `json:"exit_code"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// Manifest describes a plugin's configuration
type Manifest struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Executable   string            `json:"executable"`
	Capabilities []string          `json:"capabilities"`
	Config       map[string]string `json:"config"`
	Dependencies []string          `json:"dependencies"`
}

// Manager handles plugin lifecycle and execution
type Manager struct {
	pluginDir   string
	plugins     map[string]*LoadedPlugin
	manifests   map[string]*Manifest
	mu          sync.RWMutex
}

// LoadedPlugin represents a loaded plugin instance
type LoadedPlugin struct {
	manifest *Manifest
	client   *plugin.Client
	plugin   PluginInterface
	info     *Info
}

// Plugin RPC implementation
type PluginRPC struct {
	client *rpc.Client
}

type PluginRPCServer struct {
	// Simplified server implementation
}

// New creates a new plugin manager
func New(pluginDir string) *Manager {
	return &Manager{
		pluginDir: pluginDir,
		plugins:   make(map[string]*LoadedPlugin),
		manifests: make(map[string]*Manifest),
	}
}

// LoadPlugins discovers and loads all plugins from the plugin directory
func (m *Manager) LoadPlugins() error {
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Discover plugin manifests
	if err := m.discoverPlugins(); err != nil {
		return fmt.Errorf("failed to discover plugins: %w", err)
	}

	// Load each discovered plugin
	for name, manifest := range m.manifests {
		if err := m.loadPlugin(name, manifest); err != nil {
			fmt.Printf("Warning: failed to load plugin %s: %v\n", name, err)
			continue
		}
		fmt.Printf("✅ Loaded plugin: %s v%s\n", name, manifest.Version)
	}

	return nil
}

func (m *Manager) discoverPlugins() error {
	return filepath.WalkDir(m.pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || d.Name() != "manifest.json" {
			return nil
		}

		// Read manifest
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read manifest %s: %w", path, err)
		}

		var manifest Manifest
		if err := json.Unmarshal(content, &manifest); err != nil {
			return fmt.Errorf("failed to parse manifest %s: %w", path, err)
		}

		// Resolve executable path relative to manifest directory
		manifestDir := filepath.Dir(path)
		if !filepath.IsAbs(manifest.Executable) {
			manifest.Executable = filepath.Join(manifestDir, manifest.Executable)
		}

		m.manifests[manifest.Name] = &manifest
		return nil
	})
}

func (m *Manager) loadPlugin(name string, manifest *Manifest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin executable exists
	if _, err := os.Stat(manifest.Executable); os.IsNotExist(err) {
		return fmt.Errorf("plugin executable not found: %s", manifest.Executable)
	}

	// Create plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"plugin": &PluginRPCPlugin{},
		},
		Cmd:              exec.Command(manifest.Executable),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC},
	})

	// Connect to plugin
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to connect to plugin: %w", err)
	}

	// Get plugin instance
	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin: %w", err)
	}

	pluginInstance := raw.(PluginInterface)

	// Initialize plugin
	if err := pluginInstance.Initialize(convertToInterface(manifest.Config)); err != nil {
		client.Kill()
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}

	// Get plugin info
	info, err := pluginInstance.GetInfo()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	loadedPlugin := &LoadedPlugin{
		manifest: manifest,
		client:   client,
		plugin:   pluginInstance,
		info:     info,
	}

	m.plugins[name] = loadedPlugin
	return nil
}

// ExecutePlugin executes a command in the specified plugin
func (m *Manager) ExecutePlugin(pluginName, command string, args []string) (*Result, error) {
	m.mu.RLock()
	loadedPlugin, exists := m.plugins[pluginName]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginName)
	}

	ctx := context.Background()
	return loadedPlugin.plugin.Execute(ctx, command, args)
}

// ListPlugins returns information about all loaded plugins
func (m *Manager) ListPlugins() map[string]*Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Info)
	for name, plugin := range m.plugins {
		result[name] = plugin.info
	}
	return result
}

// GetPlugin returns information about a specific plugin
func (m *Manager) GetPlugin(name string) (*Info, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if plugin, exists := m.plugins[name]; exists {
		return plugin.info, nil
	}
	return nil, fmt.Errorf("plugin not found: %s", name)
}

// UnloadPlugin unloads a specific plugin
func (m *Manager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if plugin, exists := m.plugins[name]; exists {
		plugin.client.Kill()
		delete(m.plugins, name)
		return nil
	}
	return fmt.Errorf("plugin not found: %s", name)
}

// UnloadAllPlugins unloads all plugins
func (m *Manager) UnloadAllPlugins() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, plugin := range m.plugins {
		plugin.client.Kill()
		delete(m.plugins, name)
	}
}

// RPC Plugin implementation
type PluginRPCPlugin struct{}

func (p *PluginRPCPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &PluginRPCServer{}, nil
}

func (p *PluginRPCPlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginRPC{client: c}, nil
}

// RPC method implementations
func (p *PluginRPC) Execute(ctx context.Context, command string, args []string) (*Result, error) {
	var resp Result
	req := map[string]interface{}{
		"command": command,
		"args":    args,
	}
	
	err := p.client.Call("Plugin.Execute", req, &resp)
	return &resp, err
}

func (p *PluginRPC) GetInfo() (*Info, error) {
	var resp Info
	err := p.client.Call("Plugin.GetInfo", new(interface{}), &resp)
	return &resp, err
}

func (p *PluginRPC) Initialize(config map[string]interface{}) error {
	return p.client.Call("Plugin.Initialize", config, new(interface{}))
}

// Server RPC methods - simplified for demo
func (s *PluginRPCServer) Execute(req map[string]interface{}, resp *Result) error {
	// Simplified implementation for demo
	*resp = Result{
		Output:   "Plugin execution not fully implemented",
		ExitCode: 0,
	}
	return nil
}

func (s *PluginRPCServer) GetInfo(req interface{}, resp *Info) error {
	// Simplified implementation for demo
	*resp = Info{
		Name:        "demo-plugin",
		Version:     "1.0.0",
		Description: "Demo plugin implementation",
	}
	return nil
}

func (s *PluginRPCServer) Initialize(config map[string]interface{}, resp *interface{}) error {
	// Simplified implementation for demo
	return nil
}

// Helper functions
func convertToInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

// CreatePluginTemplate creates a basic plugin template
func CreatePluginTemplate(pluginDir, name string) error {
	pluginPath := filepath.Join(pluginDir, name)
	if err := os.MkdirAll(pluginPath, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Create manifest.json
	manifest := Manifest{
		Name:        name,
		Version:     "1.0.0",
		Description: fmt.Sprintf("%s plugin for Lucien CLI", strings.Title(name)),
		Author:      "Lucien Developer",
		Executable:  name,
		Capabilities: []string{"execute"},
		Config:      make(map[string]string),
	}

	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	manifestPath := filepath.Join(pluginPath, "manifest.json")
	
	return os.WriteFile(manifestPath, manifestData, 0644)
}