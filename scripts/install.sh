#!/bin/bash

# Lucien CLI Bash Bootstrap Installer
# Automatically downloads and installs Lucien CLI with all dependencies

set -e

# Configuration
LUCIEN_VERSION="1.0.0-nexus7"
GITHUB_REPO="luciendev/lucien-cli"
RELEASE_URL="https://api.github.com/repos/$GITHUB_REPO/releases/latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# Default configuration
INSTALL_DIR=""
PLUGIN_DIR=""
BINARY_NAME="lucien"
CREATE_SHORTCUT=false
ADD_TO_PATH=true
SKIP_PYTHON=false
VERBOSE=false

# Logging functions
log_info() { echo -e "ℹ️  $1" >&2; }
log_success() { echo -e "✅ $1" >&2; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}" >&2; }
log_error() { echo -e "${RED}❌ $1${NC}" >&2; }
log_debug() { [[ "$VERBOSE" == "true" ]] && echo -e "${GRAY}🔍 $1${NC}" >&2; }

show_help() {
    cat << 'EOF'
🚀 Lucien CLI Bash Bootstrap Installer

USAGE:
    install.sh [options]

OPTIONS:
    -v, --verbose           Enable verbose logging
    --no-path              Skip adding to system PATH
    --skip-python          Skip Python availability check
    --install-dir DIR      Custom installation directory
    -h, --help             Show this help message

EXAMPLES:
    # Standard installation
    ./install.sh

    # Verbose installation without PATH modification
    ./install.sh -v --no-path

    # Custom installation directory
    ./install.sh --install-dir "/opt/lucien"

DESCRIPTION:
This installer will:
  1. Create installation directories
  2. Install Lucien CLI binary
  3. Install Python agent plugins
  4. Check Python availability
  5. Update shell profiles (if requested)
  6. Verify installation

EOF
}

initialize_config() {
    local home_dir="$HOME"
    
    if [[ -z "$INSTALL_DIR" ]]; then
        INSTALL_DIR="$home_dir/.local/bin"
    fi
    
    PLUGIN_DIR="$home_dir/.lucien/plugins"
    
    log_debug "Configuration initialized:"
    log_debug "  Install Directory: $INSTALL_DIR"
    log_debug "  Plugin Directory: $PLUGIN_DIR"
    log_debug "  Add to PATH: $ADD_TO_PATH"
    log_debug "  Skip Python: $SKIP_PYTHON"
}

create_directories() {
    log_info "Creating installation directories..."
    
    local directories=(
        "$INSTALL_DIR"
        "$PLUGIN_DIR"
        "$PLUGIN_DIR/agents"
    )
    
    for dir in "${directories[@]}"; do
        if [[ ! -d "$dir" ]]; then
            mkdir -p "$dir"
            log_debug "Created directory: $dir"
        fi
    done
    
    log_success "Directories created successfully"
}

install_binary() {
    log_info "Installing Lucien CLI binary..."
    
    local binary_path="$INSTALL_DIR/$BINARY_NAME"
    
    # Try local build first
    local local_paths=(
        "build/lucien"
        "lucien"
        "../build/lucien"
    )
    
    for path in "${local_paths[@]}"; do
        if [[ -f "$path" ]]; then
            log_debug "Found local binary: $path"
            cp "$path" "$binary_path"
            chmod +x "$binary_path"
            log_success "Binary installed from local build: $binary_path"
            return 0
        fi
    done
    
    # Try to build from source
    if [[ -f "go.mod" ]]; then
        log_info "Building from source..."
        
        export CGO_ENABLED=0
        export GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
        export GOARCH=$(uname -m)
        
        # Map common architectures
        case "$GOARCH" in
            x86_64) GOARCH="amd64" ;;
            aarch64) GOARCH="arm64" ;;
            armv7l) GOARCH="arm" ;;
        esac
        
        go build -o "$binary_path" \
            -ldflags "-X main.version=$LUCIEN_VERSION" \
            ./cmd/lucien
        
        if [[ $? -eq 0 ]]; then
            chmod +x "$binary_path"
            log_success "Built Lucien CLI from source"
            return 0
        else
            log_error "Failed to build from source"
            return 1
        fi
    fi
    
    # Download from GitHub releases (future implementation)
    log_warning "Local binary not found and GitHub download not implemented"
    log_error "No binary source available for installation"
    return 1
}

install_plugins() {
    log_info "Installing Python agents and plugins..."
    
    local source_plugin_dir="plugins"
    if [[ ! -d "$source_plugin_dir" ]]; then
        log_warning "Plugins directory not found: $source_plugin_dir"
        log_info "Skipping plugin installation - copy manually if needed"
        return 0
    fi
    
    # Copy plugins recursively
    cp -r "$source_plugin_dir"/* "$PLUGIN_DIR/"
    
    # Ensure agent scripts are executable
    find "$PLUGIN_DIR" -name "*.py" -exec chmod +x {} \;
    
    log_success "Plugins installed to: $PLUGIN_DIR"
}

check_python() {
    if [[ "$SKIP_PYTHON" == "true" ]]; then
        log_info "Skipping Python check as requested"
        return 0
    fi
    
    log_info "Checking Python availability..."
    
    local python_candidates=("python3" "python" "py")
    
    for candidate in "${python_candidates[@]}"; do
        if command -v "$candidate" &> /dev/null; then
            local version=$("$candidate" --version 2>&1)
            log_success "Found Python: $version"
            if [[ "$version" == *"Python 3"* ]]; then
                log_success "Python 3 detected - Agent commands will work"
                return 0
            fi
        else
            log_debug "Python candidate '$candidate' not found"
        fi
    done
    
    log_warning "Python 3.7+ not found in PATH"
    log_warning "Agent commands will not work without Python"
    return 1
}

add_to_path() {
    if [[ "$ADD_TO_PATH" != "true" ]]; then
        return 0
    fi
    
    log_info "Adding Lucien CLI to shell profiles..."
    
    local profile_files=(
        "$HOME/.bashrc"
        "$HOME/.zshrc" 
        "$HOME/.profile"
    )
    
    local path_line="export PATH=\"\$PATH:$INSTALL_DIR\""
    local comment_line="# Lucien CLI"
    
    local added_to_any=false
    
    for profile_file in "${profile_files[@]}"; do
        if [[ -f "$profile_file" ]]; then
            # Check if already added
            if grep -q "$INSTALL_DIR" "$profile_file" 2>/dev/null; then
                log_debug "PATH already contains installation directory in $profile_file"
                continue
            fi
            
            # Add to profile file
            echo "" >> "$profile_file"
            echo "$comment_line" >> "$profile_file"  
            echo "$path_line" >> "$profile_file"
            
            log_success "Added to PATH in $profile_file"
            added_to_any=true
        fi
    done
    
    if [[ "$added_to_any" == "true" ]]; then
        log_info "Restart terminal or run 'source ~/.bashrc' to use 'lucien' command"
    else
        log_warning "No shell profile files found to update"
        log_info "Manually add $INSTALL_DIR to your PATH"
    fi
}

verify_installation() {
    log_info "Verifying installation..."
    
    local binary_path="$INSTALL_DIR/$BINARY_NAME"
    
    # Check binary exists
    if [[ ! -f "$binary_path" ]]; then
        log_error "Binary not found at: $binary_path"
        return 1
    fi
    
    # Test binary execution
    if ! "$binary_path" --version >/dev/null 2>&1; then
        log_error "Failed to execute binary"
        return 1
    fi
    
    local version=$("$binary_path" --version 2>/dev/null || echo "unknown")
    log_success "Binary works: $version"
    
    # Check plugin files
    local plugin_files=(
        "planner_agent.py"
        "designer_agent.py"
        "review_agent.py" 
        "code_agent.py"
    )
    
    local missing_plugins=()
    for file in "${plugin_files[@]}"; do
        local plugin_path="$PLUGIN_DIR/$file"
        if [[ ! -f "$plugin_path" ]]; then
            missing_plugins+=("$file")
            log_warning "Plugin file missing: $file"
        else
            log_debug "Plugin file found: $file"
        fi
    done
    
    if [[ ${#missing_plugins[@]} -eq 0 ]]; then
        log_success "All plugin files verified"
    fi
    
    return 0
}

show_post_install() {
    local separator=$(printf '=%.0s' {1..60})
    local binary_path="$INSTALL_DIR/$BINARY_NAME"
    
    echo ""
    echo -e "${GREEN}$separator${NC}"
    echo -e "${GREEN}🎉 LUCIEN CLI INSTALLATION COMPLETE!${NC}"
    echo -e "${GREEN}$separator${NC}"
    
    echo -e "${CYAN}📁 Installation Directory: $INSTALL_DIR${NC}"
    echo -e "${CYAN}🔌 Plugin Directory: $PLUGIN_DIR${NC}"
    echo -e "${CYAN}🚀 Binary Location: $binary_path${NC}"
    
    echo -e "\n${YELLOW}📋 QUICK START:${NC}"
    if [[ "$ADD_TO_PATH" == "true" ]]; then
        echo "   lucien --version"
        echo "   lucien"
    else
        echo "   $binary_path --version"
        echo "   $binary_path"
    fi
    
    echo -e "\n${YELLOW}🤖 AI AGENT COMMANDS:${NC}"
    echo '   plan "build a web app"'
    echo '   design "dark login form"'
    echo '   review myfile.py'
    echo '   code generate "sort function"'
    
    echo -e "\n${YELLOW}📚 DOCUMENTATION:${NC}"
    echo "   help                    # Built-in help"
    echo "   README_AGENTS.md        # Agent commands guide"
    echo "   POWERSHELL_MAPPING.md   # Command reference"
    
    if [[ "$ADD_TO_PATH" != "true" ]]; then
        echo -e "\n${RED}⚠️  IMPORTANT:${NC}"
        echo -e "${RED}   Use full path to access 'lucien' command${NC}"
    fi
    
    echo -e "\n${GREEN}Happy coding! 🚀${NC}"
    echo -e "${GREEN}$separator${NC}"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --no-path)
                ADD_TO_PATH=false
                shift
                ;;
            --skip-python)
                SKIP_PYTHON=true
                shift
                ;;
            --install-dir)
                if [[ -n $2 && $2 != -* ]]; then
                    INSTALL_DIR="$2"
                    shift 2
                else
                    log_error "Option --install-dir requires a directory path"
                    exit 1
                fi
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                log_info "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

# Main execution function
main() {
    echo -e "${CYAN}🚀 Lucien CLI Bash Bootstrap Installer${NC}"
    echo -e "${CYAN}=====================================${NC}"
    
    parse_args "$@"
    
    initialize_config
    create_directories
    install_binary
    install_plugins
    check_python
    add_to_path
    verify_installation
    
    show_post_install
}

# Run main function with all arguments
main "$@"