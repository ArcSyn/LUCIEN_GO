# Lucien CLI PowerShell Bootstrap Installer
# Safe version - Skips shortcut on failure

param(
    [switch]$Verbose,
    [switch]$NoShortcut,
    [switch]$NoPath,
    [switch]$SkipPython,
    [string]$InstallDir,
    [switch]$Help
)

function Write-Info    { Write-Host "[INFO] $args" -ForegroundColor Cyan }
function Write-Success { Write-Host "[SUCCESS] $args" -ForegroundColor Green }
function Write-Warning { Write-Host "[WARNING] $args" -ForegroundColor Yellow }
function Write-ErrorMsg { Write-Host "[ERROR] $args" -ForegroundColor Red }

Write-Host "`nLucien CLI PowerShell Bootstrap Installer"
Write-Host "============================================="

# Resolve directories
$LocalAppData = $env:LOCALAPPDATA
$UserProfile  = $env:USERPROFILE
$InstallDir   = "$LocalAppData\Lucien"
$PluginDir    = "$UserProfile\.lucien\plugins"
$BinaryName   = "lucien.exe"

Write-Info "Starting installation with configuration:"
Write-Info "InstallDir: $InstallDir"
Write-Info "PluginDir : $PluginDir"

# 1. Create directories
Write-Info "Creating installation directories..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path $PluginDir | Out-Null
New-Item -ItemType Directory -Force -Path "$PluginDir\agents" | Out-Null
Write-Success "Directories created successfully"

# 2. Copy binary
Write-Info "Installing Lucien CLI binary..."
$sourceBinary = ".\build\$BinaryName"
$destBinary = Join-Path $InstallDir $BinaryName

if (Test-Path $sourceBinary) {
    Copy-Item -Force $sourceBinary -Destination $destBinary
    Write-Success "Binary installed from local build: $destBinary"
} else {
    Write-ErrorMsg "Lucien binary not found in ./build folder"
    Write-ErrorMsg "Run 'make build' first to create the binary"
    exit 1
}

# 3. Install Python plugins
Write-Info "Installing Python agents and plugins..."
if (Test-Path ".\plugins") {
    Copy-Item -Recurse -Force ".\plugins\*" -Destination $PluginDir
    Write-Success "Plugins installed to: $PluginDir"
} else {
    Write-Warning "Plugins directory not found - agent commands may not work"
}

# 4. Check Python
if (-not $SkipPython) {
    Write-Info "Checking Python availability..."
    $pythonCandidates = @("python3", "python", "py")
    $foundPython = $false
    
    foreach ($candidate in $pythonCandidates) {
        try {
            $version = & $candidate --version 2>$null
            if ($version) {
                Write-Success "Found Python: $version"
                $foundPython = $true
                break
            }
        } catch {
            # Continue to next candidate
        }
    }
    
    if (-not $foundPython) {
        Write-Warning "Python not found. Agent commands will not work without Python 3.7+"
    }
} else {
    Write-Info "Skipping Python check as requested"
}

# 5. Add to PATH
if (-not $NoPath) {
    Write-Info "Adding Lucien CLI to system PATH..."
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        $newPath = if ($currentPath) { "$currentPath;$InstallDir" } else { $InstallDir }
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Lucien added to user PATH (restart terminal to use 'lucien' command)"
    } else {
        Write-Info "Lucien already in PATH"
    }
}

# 6. Create desktop shortcut
if (-not $NoShortcut) {
    Write-Info "Creating desktop shortcut..."
    try {
        $WScriptShell = New-Object -ComObject WScript.Shell
        $Shortcut = $WScriptShell.CreateShortcut("$env:USERPROFILE\Desktop\Lucien CLI.lnk")
        $Shortcut.TargetPath = "$InstallDir\$BinaryName"
        $Shortcut.WorkingDirectory = $env:USERPROFILE
        $Shortcut.Description = "Lucien AI-Enhanced Shell"
        $Shortcut.Save()
        Write-Success "Desktop shortcut created"
    } catch {
        Write-Warning "Could not create desktop shortcut: $_"
    }
}

# 7. Verify installation
Write-Info "Verifying installation..."
try {
    $version = & $destBinary --version 2>$null
    Write-Success "Binary verified: $version"
} catch {
    Write-Warning "Could not verify binary execution"
}

Write-Host ""
Write-Host "=============================================" -ForegroundColor Green
Write-Host "LUCIEN CLI INSTALLATION COMPLETE!" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green
Write-Host ""
Write-Host "Binary Location: $destBinary" -ForegroundColor Cyan
Write-Host "Plugin Directory: $PluginDir" -ForegroundColor Cyan
Write-Host ""
Write-Host "QUICK START:" -ForegroundColor Yellow
if (-not $NoPath) {
    Write-Host "   lucien --version" -ForegroundColor White
    Write-Host "   lucien" -ForegroundColor White
} else {
    Write-Host "   `"$destBinary`" --version" -ForegroundColor White
    Write-Host "   `"$destBinary`"" -ForegroundColor White
}
Write-Host ""
Write-Host "AI AGENT COMMANDS:" -ForegroundColor Yellow
Write-Host "   plan `"build a web app`"" -ForegroundColor White
Write-Host "   design `"dark login form`"" -ForegroundColor White
Write-Host "   review myfile.py" -ForegroundColor White
Write-Host "   code generate `"sort function`"" -ForegroundColor White
Write-Host ""
Write-Host "Happy coding!" -ForegroundColor Green