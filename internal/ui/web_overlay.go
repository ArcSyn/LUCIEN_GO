package ui

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/browser"

	"github.com/ArcSyn/LucienCLI/internal/shell"
)

// WebOverlay provides a browser-based terminal interface using xterm.js
type WebOverlay struct {
	shell       *shell.Shell
	port        int
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]*WebClient
	clientsMux  sync.RWMutex
	broadcast   chan WebMessage
	running     bool
	server      *http.Server
}

// WebClient represents a connected web client
type WebClient struct {
	conn   *websocket.Conn
	send   chan WebMessage
	id     string
	theme  string
	active bool
}

// WebMessage represents messages sent between web client and backend
type WebMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ClientID  string      `json:"clientId,omitempty"`
}

// CommandMessage represents command execution messages
type CommandMessage struct {
	Command string `json:"command"`
	Path    string `json:"path"`
}

// OutputMessage represents command output messages
type OutputMessage struct {
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exitCode"`
}

// ThemeMessage represents theme change messages
type ThemeMessage struct {
	Theme string `json:"theme"`
}

// StatusMessage represents status updates
type StatusMessage struct {
	Status    string                 `json:"status"`
	SystemStats map[string]interface{} `json:"systemStats,omitempty"`
}

// NewWebOverlay creates a new web overlay instance
func NewWebOverlay(shell *shell.Shell, port int) *WebOverlay {
	return &WebOverlay{
		shell: shell,
		port:  port,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development - restrict in production
				return true
			},
		},
		clients:   make(map[*websocket.Conn]*WebClient),
		broadcast: make(chan WebMessage, 256),
		running:   false,
	}
}

// Start launches the web overlay server
func (w *WebOverlay) Start() error {
	if w.running {
		return fmt.Errorf("web overlay already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", w.handleIndex)
	mux.HandleFunc("/ws", w.handleWebSocket)
	mux.HandleFunc("/static/", w.handleStatic)

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
		Handler: mux,
	}

	// Start the message broadcaster
	go w.broadcaster()

	// Start the HTTP server
	go func() {
		log.Printf("🌐 Web overlay starting on http://localhost:%d", w.port)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Web overlay server error: %v", err)
		}
	}()

	w.running = true

	// Open browser automatically
	time.Sleep(500 * time.Millisecond) // Give server time to start
	url := fmt.Sprintf("http://localhost:%d", w.port)
	if err := browser.OpenURL(url); err != nil {
		log.Printf("⚠️  Could not open browser: %v", err)
		log.Printf("🌐 Please open: %s", url)
	}

	return nil
}

// Stop shuts down the web overlay server
func (w *WebOverlay) Stop() error {
	if !w.running {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close all client connections
	w.clientsMux.Lock()
	for conn, client := range w.clients {
		close(client.send)
		conn.Close()
	}
	w.clients = make(map[*websocket.Conn]*WebClient)
	w.clientsMux.Unlock()

	// Shutdown server
	err := w.server.Shutdown(ctx)
	w.running = false
	close(w.broadcast)

	return err
}

// SendOutput sends command output to all connected clients
func (w *WebOverlay) SendOutput(output, error string, exitCode int) {
	if !w.running {
		return
	}

	msg := WebMessage{
		Type: "output",
		Data: OutputMessage{
			Output:   output,
			Error:    error,
			ExitCode: exitCode,
		},
		Timestamp: time.Now(),
	}

	select {
	case w.broadcast <- msg:
	default:
		// Channel full, skip message
	}
}

// SendStatus sends status updates to all connected clients
func (w *WebOverlay) SendStatus(status string, stats map[string]interface{}) {
	if !w.running {
		return
	}

	msg := WebMessage{
		Type: "status",
		Data: StatusMessage{
			Status:      status,
			SystemStats: stats,
		},
		Timestamp: time.Now(),
	}

	select {
	case w.broadcast <- msg:
	default:
	}
}

// handleIndex serves the main web terminal page
func (w *WebOverlay) handleIndex(wr http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	data := struct {
		Port int
		Title string
	}{
		Port:  w.port,
		Title: "Lucien Visual Bliss Terminal",
	}
	tmpl.Execute(wr, data)
}

// handleStatic serves static files (CSS, JS)
func (w *WebOverlay) handleStatic(wr http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[8:] // Remove "/static/" prefix
	
	switch path {
	case "xterm.css":
		wr.Header().Set("Content-Type", "text/css")
		wr.Write([]byte(xtermCSS))
	case "xterm.js":
		wr.Header().Set("Content-Type", "application/javascript")
		wr.Write([]byte(xtermJS))
	case "app.js":
		wr.Header().Set("Content-Type", "application/javascript")
		wr.Write([]byte(appJS))
	default:
		http.NotFound(wr, r)
	}
}

// handleWebSocket handles WebSocket connections
func (w *WebOverlay) handleWebSocket(wr http.ResponseWriter, r *http.Request) {
	conn, err := w.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WebClient{
		conn:   conn,
		send:   make(chan WebMessage, 256),
		id:     fmt.Sprintf("client_%d", time.Now().Unix()),
		theme:  "mocha", // Default Catppuccin theme
		active: true,
	}

	w.clientsMux.Lock()
	w.clients[conn] = client
	w.clientsMux.Unlock()

	log.Printf("🌐 New web client connected: %s", client.id)

	// Start client handlers
	go w.handleClientRead(client)
	go w.handleClientWrite(client)

	// Send welcome message
	welcomeMsg := WebMessage{
		Type: "welcome",
		Data: map[string]interface{}{
			"clientId": client.id,
			"theme":    client.theme,
			"message":  "Connected to Lucien Visual Bliss Terminal",
		},
		Timestamp: time.Now(),
	}
	client.send <- welcomeMsg
}

// handleClientRead handles messages from web client
func (w *WebOverlay) handleClientRead(client *WebClient) {
	defer func() {
		w.clientsMux.Lock()
		delete(w.clients, client.conn)
		w.clientsMux.Unlock()
		client.conn.Close()
		close(client.send)
		log.Printf("🌐 Web client disconnected: %s", client.id)
	}()

	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WebMessage
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.ClientID = client.id
		msg.Timestamp = time.Now()

		w.handleClientMessage(client, msg)
	}
}

// handleClientWrite handles messages to web client
func (w *WebOverlay) handleClientWrite(client *WebClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteJSON(msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleClientMessage processes messages from web clients
func (w *WebOverlay) handleClientMessage(client *WebClient, msg WebMessage) {
	switch msg.Type {
	case "command":
		w.handleCommandMessage(client, msg)
	case "theme":
		w.handleThemeMessage(client, msg)
	case "resize":
		w.handleResizeMessage(client, msg)
	}
}

// handleCommandMessage executes commands from web clients
func (w *WebOverlay) handleCommandMessage(client *WebClient, msg WebMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	command, ok := data["command"].(string)
	if !ok {
		return
	}

	log.Printf("🌐 Web command from %s: %s", client.id, command)

	// Execute command through shell
	result, err := w.shell.Execute(command)
	
	var outputMsg OutputMessage
	if err != nil {
		outputMsg = OutputMessage{
			Output:   "",
			Error:    err.Error(),
			ExitCode: 1,
		}
	} else {
		outputMsg = OutputMessage{
			Output:   result.Output,
			Error:    result.Error,
			ExitCode: result.ExitCode,
		}
	}

	// Send result back to client
	responseMsg := WebMessage{
		Type:      "output",
		Data:      outputMsg,
		Timestamp: time.Now(),
		ClientID:  client.id,
	}

	client.send <- responseMsg
}

// handleThemeMessage changes the theme for a client
func (w *WebOverlay) handleThemeMessage(client *WebClient, msg WebMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	theme, ok := data["theme"].(string)
	if !ok {
		return
	}

	client.theme = theme
	log.Printf("🌐 Client %s changed theme to: %s", client.id, theme)

	// Send theme confirmation
	responseMsg := WebMessage{
		Type: "theme_changed",
		Data: ThemeMessage{
			Theme: theme,
		},
		Timestamp: time.Now(),
		ClientID:  client.id,
	}

	client.send <- responseMsg
}

// handleResizeMessage handles terminal resize events
func (w *WebOverlay) handleResizeMessage(client *WebClient, msg WebMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	cols, _ := data["cols"].(float64)
	rows, _ := data["rows"].(float64)

	log.Printf("🌐 Client %s resized terminal: %dx%d", client.id, int(cols), int(rows))
}

// broadcaster sends messages to all connected clients
func (w *WebOverlay) broadcaster() {
	for msg := range w.broadcast {
		w.clientsMux.RLock()
		for conn, client := range w.clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(w.clients, conn)
			}
		}
		w.clientsMux.RUnlock()
	}
}

// HTML template for the web terminal interface
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/xterm.css">
    <style>
        body {
            margin: 0;
            padding: 0;
            background: #1e1e2e;
            font-family: 'Fira Code', 'JetBrains Mono', 'Cascadia Code', monospace;
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #cba6f7, #89b4fa);
            color: #1e1e2e;
            padding: 10px 20px;
            font-weight: bold;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .title {
            font-size: 18px;
        }
        
        .status {
            font-size: 14px;
        }
        
        .terminal-container {
            height: calc(100vh - 60px);
            padding: 10px;
        }
        
        #terminal {
            height: 100%;
            width: 100%;
        }
        
        .theme-selector {
            position: absolute;
            top: 15px;
            right: 20px;
            z-index: 1000;
        }
        
        .theme-selector select {
            background: #313244;
            color: #cdd6f4;
            border: 1px solid #45475a;
            border-radius: 4px;
            padding: 5px;
        }
        
        .connection-indicator {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background: #a6e3a1;
            display: inline-block;
            margin-right: 5px;
        }
        
        .connection-indicator.disconnected {
            background: #f38ba8;
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="title">
            🚀 Lucien Visual Bliss Terminal - Web Interface
        </div>
        <div class="status">
            <span class="connection-indicator" id="connectionStatus"></span>
            <span id="statusText">Connecting...</span>
        </div>
        <div class="theme-selector">
            <select id="themeSelect">
                <option value="mocha">Catppuccin Mocha</option>
                <option value="latte">Catppuccin Latte</option>
                <option value="frappe">Catppuccin Frappé</option>
                <option value="macchiato">Catppuccin Macchiato</option>
            </select>
        </div>
    </div>
    <div class="terminal-container">
        <div id="terminal"></div>
    </div>
    
    <script src="/static/xterm.js"></script>
    <script src="/static/app.js"></script>
</body>
</html>`

// Minimal xterm.js CSS (in production, load from CDN or bundle properly)
const xtermCSS = `
.xterm {
    font-feature-settings: "liga" 0;
    position: relative;
    user-select: none;
    -ms-user-select: none;
    -webkit-user-select: none;
}

.xterm.focus,
.xterm:focus {
    outline: none;
}

.xterm .xterm-helpers {
    position: absolute;
    top: 0;
    z-index: 5;
}

.xterm .xterm-helper-textarea {
    position: absolute;
    opacity: 0;
    left: -9999em;
    top: 0;
    width: 0;
    height: 0;
    z-index: -5;
    white-space: nowrap;
    overflow: hidden;
    resize: none;
}

.xterm .composition-view {
    background: #000;
    color: #FFF;
    display: none;
    position: absolute;
    white-space: nowrap;
    z-index: 1;
}

.xterm .composition-view.active {
    display: block;
}

.xterm .xterm-viewport {
    background-color: #000;
    overflow-y: scroll;
    cursor: default;
    position: absolute;
    right: 0;
    left: 0;
    top: 0;
    bottom: 0;
}

.xterm .xterm-screen {
    position: relative;
}

.xterm .xterm-screen canvas {
    position: absolute;
    left: 0;
    top: 0;
}

.xterm .xterm-scroll-area {
    visibility: hidden;
}

.xterm-char-measure-element {
    display: inline-block;
    visibility: hidden;
    position: absolute;
    top: 0;
    left: -9999em;
    line-height: normal;
}

.xterm.enable-mouse-events {
    cursor: default;
}

.xterm.xterm-cursor-pointer {
    cursor: pointer;
}

.xterm.column-select.focus {
    cursor: crosshair;
}

.xterm .xterm-accessibility,
.xterm .xterm-message {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    z-index: 10;
    color: transparent;
}

.xterm .live-region {
    position: absolute;
    left: -9999px;
    width: 1px;
    height: 1px;
    overflow: hidden;
}

.xterm-dim {
    opacity: 0.5;
}

.xterm-underline-1 { text-decoration: underline; }
.xterm-underline-2 { text-decoration: double underline; }
.xterm-underline-3 { text-decoration: wavy underline; }
.xterm-underline-4 { text-decoration: dotted underline; }
.xterm-underline-5 { text-decoration: dashed underline; }

.xterm-overline {
    text-decoration: overline;
}

.xterm-overline.xterm-underline-1 { text-decoration: overline underline; }
.xterm-overline.xterm-underline-2 { text-decoration: overline double underline; }
.xterm-overline.xterm-underline-3 { text-decoration: overline wavy underline; }
.xterm-overline.xterm-underline-4 { text-decoration: overline dotted underline; }
.xterm-overline.xterm-underline-5 { text-decoration: overline dashed underline; }

.xterm-strikethrough {
    text-decoration: line-through;
}

.xterm-screen .xterm-decoration-container .xterm-decoration {
	z-index: 6;
	position: absolute;
}

.xterm-screen .xterm-decoration-container .xterm-decoration.xterm-decoration-top-layer {
	z-index: 7;
}

.xterm-decoration-overview-ruler {
    z-index: 8;
    position: absolute;
    top: 0;
    right: 0;
    pointer-events: none;
}

.xterm-decoration-top {
    z-index: 2;
    position: relative;
}
`

// Minimal xterm.js JavaScript (in production, load from CDN or bundle properly)
const xtermJS = `
// This would contain the xterm.js library code
// For production, load from CDN: https://cdn.jsdelivr.net/npm/xterm@5.1.0/lib/xterm.min.js
console.log("xterm.js placeholder - load from CDN in production");
`

// Application JavaScript for the web terminal
const appJS = `
class LucienWebTerminal {
    constructor() {
        this.term = null;
        this.ws = null;
        this.currentTheme = 'mocha';
        this.connected = false;
        this.init();
    }

    init() {
        // Initialize xterm.js terminal (placeholder - load real xterm.js)
        if (typeof Terminal !== 'undefined') {
            this.term = new Terminal({
                cursorBlink: true,
                theme: this.getThemeColors('mocha'),
                fontFamily: 'Fira Code, JetBrains Mono, Cascadia Code, monospace',
                fontSize: 14,
            });
            
            this.term.open(document.getElementById('terminal'));
        } else {
            // Fallback for when xterm.js isn't loaded
            document.getElementById('terminal').innerHTML = 
                '<div style="color: #cdd6f4; padding: 20px; font-family: monospace;">' +
                '⚠️ Loading terminal interface...<br><br>' +
                'Note: Load xterm.js from CDN for full functionality<br>' +
                'WebSocket connection will still work for basic terminal emulation.' +
                '</div>';
        }

        this.connectWebSocket();
        this.setupEventListeners();
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = protocol + '//' + window.location.host + '/ws';
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            this.connected = true;
            this.updateConnectionStatus('Connected', true);
            console.log('🌐 WebSocket connected');
        };
        
        this.ws.onclose = () => {
            this.connected = false;
            this.updateConnectionStatus('Disconnected', false);
            console.log('🌐 WebSocket disconnected');
            
            // Attempt to reconnect after 3 seconds
            setTimeout(() => {
                if (!this.connected) {
                    this.connectWebSocket();
                }
            }, 3000);
        };
        
        this.ws.onerror = (error) => {
            console.error('🌐 WebSocket error:', error);
            this.updateConnectionStatus('Error', false);
        };
        
        this.ws.onmessage = (event) => {
            this.handleMessage(JSON.parse(event.data));
        };
    }

    setupEventListeners() {
        // Theme selector
        const themeSelect = document.getElementById('themeSelect');
        themeSelect.addEventListener('change', (e) => {
            this.changeTheme(e.target.value);
        });

        // Terminal input (if xterm.js is loaded)
        if (this.term) {
            let currentLine = '';
            
            this.term.onKey((e) => {
                const { key, domEvent } = e;
                
                if (domEvent.key === 'Enter') {
                    this.sendCommand(currentLine);
                    currentLine = '';
                    this.term.write('\r\n');
                } else if (domEvent.key === 'Backspace') {
                    if (currentLine.length > 0) {
                        currentLine = currentLine.slice(0, -1);
                        this.term.write('\b \b');
                    }
                } else if (domEvent.key.length === 1) {
                    currentLine += domEvent.key;
                    this.term.write(domEvent.key);
                }
            });
            
            // Initial prompt
            this.term.write('lucien@nexus:~$ ');
        }

        // Window resize
        window.addEventListener('resize', () => {
            if (this.term) {
                this.term.fit();
            }
        });
    }

    handleMessage(message) {
        console.log('🌐 Received message:', message);
        
        switch (message.type) {
            case 'welcome':
                this.handleWelcomeMessage(message);
                break;
            case 'output':
                this.handleOutputMessage(message);
                break;
            case 'theme_changed':
                this.handleThemeChanged(message);
                break;
            case 'status':
                this.handleStatusMessage(message);
                break;
        }
    }

    handleWelcomeMessage(message) {
        const welcomeText = message.data.message + '\r\n\r\n';
        if (this.term) {
            this.term.write(welcomeText);
        }
    }

    handleOutputMessage(message) {
        const { output, error, exitCode } = message.data;
        
        if (this.term) {
            if (output) {
                this.term.write(output);
            }
            if (error) {
                this.term.write('\x1b[31m' + error + '\x1b[0m'); // Red text
            }
            this.term.write('\r\nlucien@nexus:~$ ');
        } else {
            // Fallback display
            const terminalDiv = document.getElementById('terminal');
            terminalDiv.innerHTML += 
                '<div style="color: #cdd6f4;">' + (output || '') + '</div>' +
                (error ? '<div style="color: #f38ba8;">' + error + '</div>' : '') +
                '<div style="color: #a6e3a1;">lucien@nexus:~$ </div>';
            terminalDiv.scrollTop = terminalDiv.scrollHeight;
        }
    }

    handleThemeChanged(message) {
        const theme = message.data.theme;
        this.currentTheme = theme;
        
        if (this.term) {
            this.term.options.theme = this.getThemeColors(theme);
        }
        
        // Update document theme
        this.updateDocumentTheme(theme);
    }

    handleStatusMessage(message) {
        console.log('📊 Status update:', message.data);
    }

    sendCommand(command) {
        if (!this.connected || !this.ws) {
            console.error('⚠️  Not connected to server');
            return;
        }

        const message = {
            type: 'command',
            data: {
                command: command.trim()
            },
            timestamp: new Date().toISOString()
        };

        this.ws.send(JSON.stringify(message));
    }

    changeTheme(theme) {
        if (!this.connected || !this.ws) {
            return;
        }

        const message = {
            type: 'theme',
            data: {
                theme: theme
            },
            timestamp: new Date().toISOString()
        };

        this.ws.send(JSON.stringify(message));
    }

    updateConnectionStatus(status, connected) {
        const statusIndicator = document.getElementById('connectionStatus');
        const statusText = document.getElementById('statusText');
        
        statusIndicator.className = connected ? 'connection-indicator' : 'connection-indicator disconnected';
        statusText.textContent = status;
    }

    getThemeColors(theme) {
        const themes = {
            mocha: {
                background: '#1e1e2e',
                foreground: '#cdd6f4',
                cursor: '#f5e0dc',
                selection: '#45475a'
            },
            latte: {
                background: '#eff1f5',
                foreground: '#4c4f69',
                cursor: '#dc8a78',
                selection: '#ccd0da'
            },
            frappe: {
                background: '#303446',
                foreground: '#c6d0f5',
                cursor: '#f2d5cf',
                selection: '#414559'
            },
            macchiato: {
                background: '#24273a',
                foreground: '#cad3f5',
                cursor: '#f4dbd6',
                selection: '#363a4f'
            }
        };
        
        return themes[theme] || themes.mocha;
    }

    updateDocumentTheme(theme) {
        const colors = this.getThemeColors(theme);
        document.body.style.backgroundColor = colors.background;
        document.body.style.color = colors.foreground;
    }
}

// Initialize the terminal when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new LucienWebTerminal();
});
`