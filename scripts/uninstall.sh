#!/bin/bash

# Lucien CLI Bash Uninstaller
# Completely removes Lucien CLI from Unix/Linux systems

set -e

# Configuration
VERBOSE=false
KEEP_PLUGINS=false
HELP=false

# Default directories
HOME_DIR="$HOME"
INSTALL_DIR="$HOME_DIR/.local/bin"
PLUGIN_DIR="$HOME_DIR/.lucien"
BINARY_NAME="lucien"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# Logging functions
log_info() { echo -e "${CYAN}[INFO] $1${NC}"; }
log_success() { echo -e "${GREEN}[SUCCESS] $1${NC}"; }
log_warning() { echo -e "${YELLOW}[WARNING] $1${NC}"; }
log_error() { echo -e "${RED}[ERROR] $1${NC}"; }

show_help() {
    cat << 'EOF'
Lucien CLI Bash Uninstaller

USAGE:
    uninstall.sh [options]

OPTIONS:
    -v, --verbose           Enable verbose logging
    --keep-plugins         Keep plugin directory and agent files
    -h, --help             Show this help message

DESCRIPTION:
This uninstaller will:
  1. Remove Lucien CLI binary
  2. Remove from shell profiles (PATH)
  3. Remove plugin directory (unless --keep-plugins specified)
  4. Clean up configuration files

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --keep-plugins)
                KEEP_PLUGINS=true
                shift
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

# Remove binary from installation directory
remove_binary() {
    log_info "Removing Lucien CLI binary..."
    
    local binary_path="$INSTALL_DIR/$BINARY_NAME"
    
    if [[ -f "$binary_path" ]]; then
        rm -f "$binary_path"
        log_success "Binary removed: $binary_path"
    else
        log_info "Binary not found at: $binary_path"
    fi
}

# Remove from shell profiles
remove_from_path() {
    log_info "Removing Lucien CLI from shell profiles..."
    
    local profile_files=(
        "$HOME/.bashrc"
        "$HOME/.zshrc" 
        "$HOME/.profile"
    )
    
    local removed_from_any=false
    
    for profile_file in "${profile_files[@]}"; do
        if [[ -f "$profile_file" ]]; then
            # Check if Lucien CLI is in the profile
            if grep -q "$INSTALL_DIR" "$profile_file" 2>/dev/null; then
                # Create backup
                cp "$profile_file" "$profile_file.backup"
                
                # Remove Lucien CLI entries
                grep -v "$INSTALL_DIR" "$profile_file.backup" > "$profile_file"
                
                # Remove Lucien CLI comment lines
                sed -i '/# Lucien CLI/d' "$profile_file" 2>/dev/null || true
                
                log_success "Removed from PATH in $profile_file"
                removed_from_any=true
            fi
        fi
    done
    
    if [[ "$removed_from_any" == "true" ]]; then
        log_info "Restart terminal or run 'source ~/.bashrc' for PATH changes to take effect"
    else
        log_info "Lucien CLI not found in any shell profiles"
    fi
}

# Remove plugin directory
remove_plugins() {
    if [[ "$KEEP_PLUGINS" == "true" ]]; then
        log_info "Keeping plugin directory as requested: $PLUGIN_DIR"
        return
    fi
    
    if [[ -d "$PLUGIN_DIR" ]]; then
        log_info "Removing plugin directory..."
        rm -rf "$PLUGIN_DIR"
        log_success "Plugin directory removed: $PLUGIN_DIR"
    else
        log_info "Plugin directory not found: $PLUGIN_DIR"
    fi
}

# Clean up configuration files
cleanup_config() {
    log_info "Cleaning up configuration files..."
    
    # Remove any Lucien-specific config files
    local config_files=(
        "$HOME/.lucien_config"
        "$HOME/.config/lucien"
    )
    
    for config_file in "${config_files[@]}"; do
        if [[ -e "$config_file" ]]; then
            rm -rf "$config_file"
            log_success "Removed config: $config_file"
        fi
    done
}

# Kill any running processes
kill_processes() {
    log_info "Checking for running Lucien processes..."
    
    local pids=$(pgrep -f "$BINARY_NAME" 2>/dev/null || true)
    
    if [[ -n "$pids" ]]; then
        log_warning "Found running Lucien processes. Terminating..."
        echo "$pids" | xargs kill -TERM 2>/dev/null || true
        sleep 2
        
        # Force kill if still running
        pids=$(pgrep -f "$BINARY_NAME" 2>/dev/null || true)
        if [[ -n "$pids" ]]; then
            echo "$pids" | xargs kill -KILL 2>/dev/null || true
        fi
        
        log_success "Processes terminated"
    else
        log_info "No running Lucien processes found"
    fi
}

# Main execution function
main() {
    echo -e "${CYAN}Lucien CLI Bash Uninstaller${NC}"
    echo -e "${CYAN}===========================${NC}"
    
    parse_args "$@"
    
    log_info "Starting uninstallation process..."
    log_info "Install Directory: $INSTALL_DIR"
    log_info "Plugin Directory: $PLUGIN_DIR"
    
    kill_processes
    remove_from_path
    remove_binary
    remove_plugins
    cleanup_config
    
    echo ""
    echo -e "${GREEN}=================================${NC}"
    echo -e "${GREEN}LUCIEN CLI UNINSTALL COMPLETE!${NC}"
    echo -e "${GREEN}=================================${NC}"
    echo ""
    echo -e "${YELLOW}What was removed:${NC}"
    echo -e "  - Lucien CLI binary"
    echo -e "  - Shell profile entries (PATH)"
    if [[ "$KEEP_PLUGINS" != "true" ]]; then
        echo -e "  - Plugin directory and agent files"
    fi
    echo -e "  - Configuration files"
    echo ""
    echo -e "${CYAN}Note: Restart your terminal for PATH changes to take effect${NC}"
    echo -e "${GREEN}Thank you for using Lucien CLI!${NC}"
}

# Run main function with all arguments
main "$@"