# 🔧 LUCIEN CLI - TROUBLESHOOTING GUIDE
**Version**: 1.0-alpha  
**Last Updated**: 2025-08-06  
**Support Level**: Development/Community  

## 🎯 OVERVIEW

This guide provides solutions to common issues encountered when building, installing, and running Lucien CLI. Issues are categorized by severity and component area for quick resolution.

## 📊 ISSUE CATEGORIES

| Category | Description | Common Causes |
|----------|-------------|---------------|
| 🔥 Critical | System crashes, security issues | Code bugs, memory issues |
| ⚠️ High | Core functionality broken | Configuration errors, dependencies |
| 🛠️ Medium | Features not working | Environment setup, permissions |
| 💡 Low | Minor annoyances, UX issues | Settings, preferences |

## 🔥 CRITICAL ISSUES

### Issue: Variable Expansion Returns Literal Text
**Severity**: CRITICAL  
**Component**: Shell Engine  
**Affects**: All variable operations

#### Symptoms
```bash
set TEST "hello"
echo $UNDEFINED_VAR
# Returns: $UNDEFINED_VAR (should return empty string)
```

#### Root Cause
Bug in `internal/shell/shell.go:expandVariables()` function doesn't handle undefined variables correctly.

#### Workaround
```bash
# Always ensure variables are defined before use
set DEFAULT_VALUE ""
set MY_VAR $DEFAULT_VALUE  # Set default first
set MY_VAR "actual_value"   # Then set real value
```

#### Permanent Fix (For Developers)
```go
// In internal/shell/shell.go, modify expandVariables()
func (s *Shell) expandVariables(cmdStr string) string {
    result := cmdStr
    for key, value := range s.env {
        result = strings.ReplaceAll(result, "$"+key, value)
    }
    // Fix: Handle undefined variables
    undefinedVarRegex := regexp.MustCompile(`\$[A-Za-z_][A-Za-z0-9_]*`)
    result = undefinedVarRegex.ReplaceAllString(result, "")
    return result
}
```

### Issue: Duration Measurement Missing
**Severity**: CRITICAL  
**Component**: Shell Engine  
**Affects**: Performance monitoring, benchmarking

#### Symptoms
```bash
# Builtin commands don't show execution time
echo "test"
# Duration field is 0 or missing in ExecutionResult
```

#### Root Cause
Builtin command functions don't implement duration tracking.

#### Workaround
```bash
# Use external time command for timing
time echo "test"  # If available on system
```

#### Permanent Fix (For Developers)
```go
// Add to each builtin command function
func (s *Shell) echo(args []string) (*ExecutionResult, error) {
    start := time.Now()
    // ... existing code ...
    return &ExecutionResult{
        Output:   output,
        Duration: time.Since(start),  // Add this line
    }, nil
}
```

### Issue: Empty Command Throws Error
**Severity**: CRITICAL  
**Component**: Command Parser  
**Affects**: User experience

#### Symptoms
```bash
# Pressing Enter on empty command line
> 
Error: invalid command
```

#### Root Cause
Parser doesn't handle empty input gracefully.

#### Workaround
```bash
# Always type something, even spaces
>  # (space then enter)
```

#### Permanent Fix (For Developers)
```go
// In parseCommand(), add empty string check
func (s *Shell) parseCommand(cmdStr string) (*Command, error) {
    if cmdStr == "" || strings.TrimSpace(cmdStr) == "" {
        return nil, nil // Return nil without error
    }
    // ... rest of function
}
```

## ⚠️ HIGH PRIORITY ISSUES

### Issue: Policy Engine Compilation Error
**Severity**: HIGH  
**Component**: Security Policy Engine  
**Affects**: Safe mode operation, security features

#### Symptoms
```bash
go test ./internal/policy/
# Error: syntax error in engine_test.go:427
```

#### Root Cause
Syntax error in test file prevents policy engine testing.

#### Workaround
```bash
# Disable safe mode temporarily
lucien --no-safe-mode
```

#### Solution
1. Open `internal/policy/engine_test.go`
2. Fix syntax error on line 427
3. Run tests: `go test ./internal/policy/`

### Issue: Plugin Loading Fails
**Severity**: HIGH  
**Component**: Plugin System  
**Affects**: All plugin functionality

#### Symptoms
```bash
lucien plugin list
# Error: failed to load plugin: manifest.json not found
```

#### Troubleshooting Steps
1. **Verify Plugin Directory**
   ```bash
   ls -la plugins/
   # Should show plugin directories with manifest.json
   ```

2. **Check Plugin Manifest**
   ```bash
   cat plugins/example-bmad/manifest.json
   # Verify JSON is valid
   ```

3. **Verify Plugin Permissions**
   ```bash
   # Windows
   dir plugins\example-bmad
   
   # Linux/Mac
   chmod +x plugins/example-bmad/bmad.exe
   ```

4. **Test Plugin Manually**
   ```bash
   cd plugins/example-bmad
   ./bmad.exe --test
   ```

#### Common Solutions
```bash
# Rebuild plugins
cd plugins/example-bmad
go build -o bmad.exe .

# Reset plugin directory
rm -rf ~/.lucien/plugins
mkdir -p ~/.lucien/plugins
cp -r plugins/* ~/.lucien/plugins/
```

### Issue: Command Syntax Not Recognized
**Severity**: HIGH  
**Component**: Shell Parser  
**Affects**: Standard shell operations

#### Symptoms
```bash
set VAR=value
# Error: invalid syntax (should work in standard shells)
```

#### Root Cause
Lucien only supports `set VAR value` syntax, not `VAR=value`.

#### Workaround
```bash
# Use Lucien-specific syntax
set VAR value

# For environment export
export VAR
```

#### User Adaptation
```bash
# Instead of:
PATH=/usr/bin:$PATH

# Use:
set PATH "/usr/bin:$PATH"
export PATH
```

## 🛠️ MEDIUM PRIORITY ISSUES

### Issue: Python Dependencies Missing
**Severity**: MEDIUM  
**Component**: Python CLI Layer  
**Affects**: UI and CLI functionality

#### Symptoms
```bash
python cli/main.py
# ModuleNotFoundError: No module named 'rich'
```

#### Solution
```bash
# Install requirements
pip install -r requirements.txt

# Or install specific packages
pip install rich typer pyyaml
```

#### Virtual Environment Setup
```bash
# Create virtual environment
python -m venv .venv

# Activate (Windows)
.venv\Scripts\activate

# Activate (Linux/Mac)
source .venv/bin/activate

# Install dependencies
pip install -r requirements.txt
```

### Issue: Go Build Fails with Dependency Errors
**Severity**: MEDIUM  
**Component**: Build System  
**Affects**: Compilation process

#### Symptoms
```bash
go build cmd/lucien/main.go
# Error: module not found or version conflicts
```

#### Solution Steps
```bash
# Clean module cache
go clean -modcache

# Update dependencies
go mod tidy
go mod download

# Verify modules
go mod verify

# Retry build
go build cmd/lucien/main.go
```

### Issue: Slow Startup Performance
**Severity**: MEDIUM  
**Component**: Initialization  
**Affects**: User experience

#### Symptoms
- Lucien takes >5 seconds to start
- High CPU usage during startup

#### Troubleshooting
```bash
# Enable debug mode to see initialization steps
LUCIEN_DEBUG=1 lucien --debug

# Profile startup
go build -o lucien-debug cmd/lucien/main.go
./lucien-debug --profile
```

#### Solutions
```bash
# Disable auto-loading plugins
lucien --no-plugins

# Use minimal configuration
lucien --config configs/minimal.toml

# Check plugin directory permissions
ls -la ~/.lucien/plugins/
```

### Issue: Cross-Platform Compatibility Problems
**Severity**: MEDIUM  
**Component**: Platform Support  
**Affects**: Non-Windows systems

#### Linux-Specific Issues
```bash
# Sandbox manager issues
Error: gVisor not available

# Solution: Disable sandbox temporarily
lucien --no-sandbox
```

#### macOS-Specific Issues
```bash
# Permission issues
Error: permission denied accessing plugin

# Solution: Fix permissions
chmod +x lucien
xattr -d com.apple.quarantine lucien
```

## 💡 LOW PRIORITY ISSUES

### Issue: Theme Not Applied
**Severity**: LOW  
**Component**: UI System  
**Affects**: Visual appearance

#### Symptoms
```bash
:theme synthwave
# Theme doesn't change or reverts
```

#### Solution
```bash
# Check configuration file
cat ~/.lucien/config.toml

# Set theme in config
[ui]
theme = "synthwave"

# Restart Lucien
```

### Issue: History Not Persistent
**Severity**: LOW  
**Component**: Shell History  
**Affects**: Command history between sessions

#### Symptoms
- Command history lost after restart
- `history` command shows empty

#### Solution
```bash
# Check history file
ls -la ~/.lucien/history

# Set history file location
export LUCIEN_HISTORY_FILE=~/.lucien/history

# Configure history size
set HISTORY_SIZE 1000
```

### Issue: Alias Not Expanding
**Severity**: LOW  
**Component**: Alias System  
**Affects**: Command shortcuts

#### Symptoms
```bash
alias ll='ls -la'
ll
# Alias not recognized
```

#### Troubleshooting
```bash
# List current aliases
alias

# Verify alias syntax
alias test='echo hello'
test

# Check alias expansion in debug mode
LUCIEN_DEBUG=1 lucien --debug
```

## 🐛 DEBUGGING TECHNIQUES

### Enable Debug Mode
```bash
# Environment variable
export LUCIEN_DEBUG=1
lucien

# Command line flag
lucien --debug

# Configuration file
[debug]
enabled = true
log_level = "trace"
```

### Debug Output Analysis
```bash
# Debug log locations
~/.lucien/logs/debug.log
~/.lucien/logs/error.log

# Follow logs in real-time
tail -f ~/.lucien/logs/debug.log
```

### Component-Specific Debugging

#### Shell Engine Debugging
```bash
# Enable shell debug
export LUCIEN_SHELL_DEBUG=1

# Test specific command
echo "set TEST value" | lucien --debug
```

#### Plugin System Debugging
```bash
# Enable plugin debug
export LUCIEN_PLUGIN_DEBUG=1

# Test plugin loading
lucien plugin load --debug ./plugins/example-bmad
```

#### Policy Engine Debugging
```bash
# Enable policy debug
export LUCIEN_POLICY_DEBUG=1

# Test policy evaluation
lucien --safe --debug
```

## 🔍 DIAGNOSTIC TOOLS

### System Information
```bash
# Lucien system info
lucien --version
lucien --system-info

# Environment check
lucien --validate-environment
```

### Health Check
```bash
# Run comprehensive health check
lucien --health-check

# Check specific components
lucien --check-plugins
lucien --check-config
lucien --check-dependencies
```

### Performance Profiling
```bash
# CPU profiling
lucien --cpu-profile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
lucien --mem-profile=mem.prof
go tool pprof mem.prof
```

## 📋 TROUBLESHOOTING CHECKLIST

### Before Reporting Issues

#### ✅ Basic Checks
- [ ] Latest version installed
- [ ] All dependencies installed
- [ ] Configuration file valid
- [ ] Permissions set correctly
- [ ] Debug mode enabled

#### ✅ Environment Verification
- [ ] Go version 1.21+
- [ ] Python version 3.11+
- [ ] System requirements met
- [ ] Network connectivity available
- [ ] Disk space sufficient

#### ✅ Component Testing
- [ ] Shell commands work individually
- [ ] Plugins load without errors
- [ ] Configuration parses correctly
- [ ] Log files accessible
- [ ] Debug output captured

### Information to Collect

#### System Information
```bash
# Operating system
uname -a  # Linux/Mac
systeminfo  # Windows

# Lucien version
lucien --version

# Go environment
go env
```

#### Error Information
```bash
# Full error message
lucien --debug 2>&1 | tee error.log

# Stack trace (if available)
lucien --stack-trace

# Configuration dump
lucien --dump-config
```

#### Reproduction Steps
1. Exact commands that cause the issue
2. Expected behavior
3. Actual behavior
4. Frequency (always, sometimes, once)
5. Environment details

## 🆘 GETTING HELP

### Self-Help Resources

#### Documentation
- Technical Manual: `TECHNICAL_MANUAL.md`
- API Documentation: `API_DOCUMENTATION.md`
- Security Guide: `SECURITY_ASSESSMENT_SUMMARY.md`

#### Debug Commands
```bash
# Comprehensive validation
lucien --validate

# Component health check
lucien --health

# Configuration verification
lucien --check-config
```

### Community Support

#### Issue Reporting
1. Search existing issues first
2. Use issue templates
3. Include reproduction steps
4. Attach debug logs
5. Specify environment details

#### Best Practices for Help Requests
- Be specific about the problem
- Include error messages verbatim
- Provide minimal reproduction case
- Include system information
- Be patient and respectful

### Developer Support

#### For Development Issues
```bash
# Run development tests
make test

# Check build environment
make check-env

# Validate development setup
make validate
```

## 🔄 RECOVERY PROCEDURES

### Reset to Clean State
```bash
# Backup current configuration
cp ~/.lucien/config.toml ~/.lucien/config.toml.backup

# Remove user data
rm -rf ~/.lucien/

# Reinstall fresh
lucien --first-time-setup
```

### Plugin System Recovery
```bash
# Disable all plugins
lucien --no-plugins

# Reset plugin directory
rm -rf ~/.lucien/plugins
mkdir -p ~/.lucien/plugins

# Reinstall default plugins
cp -r plugins/* ~/.lucien/plugins/
```

### Configuration Recovery
```bash
# Reset to default configuration
lucien --reset-config

# Validate configuration
lucien --validate-config

# Apply safe defaults
cp configs/safe-mode.toml ~/.lucien/config.toml
```

---

## 📞 SUPPORT CONTACTS

**For Development Issues**: GitHub Issues  
**For Security Issues**: Report privately first  
**For Documentation**: Update requests welcome  

**Remember**: This is development software with known issues. Production use is not recommended until security vulnerabilities are resolved.

**Status**: 🔶 Development Support Only

*This troubleshooting guide reflects the current state of Lucien CLI v1.0-alpha. As issues are resolved, this guide will be updated accordingly.*