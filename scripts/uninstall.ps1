# Lucien CLI PowerShell Uninstaller
# Completely removes Lucien CLI from the system

param(
    [switch]$Verbose,
    [switch]$KeepPlugins,
    [switch]$Help
)

function Write-Info    { Write-Host "[INFO] $args" -ForegroundColor Cyan }
function Write-Success { Write-Host "[SUCCESS] $args" -ForegroundColor Green }
function Write-Warning { Write-Host "[WARNING] $args" -ForegroundColor Yellow }
function Write-ErrorMsg { Write-Host "[ERROR] $args" -ForegroundColor Red }

if ($Help) {
    Write-Host @"
Lucien CLI PowerShell Uninstaller

USAGE:
    uninstall.ps1 [options]

OPTIONS:
    -Verbose        Enable verbose logging
    -KeepPlugins    Keep plugin directory and agent files
    -Help           Show this help message

DESCRIPTION:
This uninstaller will:
  1. Remove Lucien CLI binary and installation directory
  2. Remove from system PATH
  3. Remove desktop shortcut
  4. Remove plugin directory (unless -KeepPlugins specified)
  5. Clean up registry entries

"@
    exit 0
}

Write-Host "`nLucien CLI PowerShell Uninstaller"
Write-Host "=================================="

# Resolve directories
$LocalAppData = $env:LOCALAPPDATA
$UserProfile  = $env:USERPROFILE
$InstallDir   = "$LocalAppData\Lucien"
$PluginDir    = "$UserProfile\.lucien"
$BinaryName   = "lucien.exe"
$ShortcutPath = "$UserProfile\Desktop\Lucien CLI.lnk"

Write-Info "Starting uninstallation process..."
Write-Info "InstallDir: $InstallDir"
Write-Info "PluginDir : $PluginDir"

# 1. Remove from PATH
Write-Info "Removing Lucien CLI from system PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -like "*$InstallDir*") {
    $newPath = $currentPath -split ";" | Where-Object { $_ -ne $InstallDir } | Join-String -Separator ";"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Success "Removed from system PATH"
} else {
    Write-Info "Lucien not found in PATH"
}

# 2. Remove desktop shortcut
if (Test-Path $ShortcutPath) {
    Write-Info "Removing desktop shortcut..."
    Remove-Item -Force $ShortcutPath
    Write-Success "Desktop shortcut removed"
} else {
    Write-Info "No desktop shortcut found"
}

# 3. Stop any running processes
Write-Info "Checking for running Lucien processes..."
$processes = Get-Process -Name "lucien" -ErrorAction SilentlyContinue
if ($processes) {
    Write-Warning "Found running Lucien processes. Terminating..."
    $processes | Stop-Process -Force
    Write-Success "Processes terminated"
    Start-Sleep -Seconds 2
}

# 4. Remove installation directory
if (Test-Path $InstallDir) {
    Write-Info "Removing installation directory..."
    try {
        Remove-Item -Recurse -Force $InstallDir
        Write-Success "Installation directory removed: $InstallDir"
    } catch {
        Write-ErrorMsg "Failed to remove installation directory: $_"
        Write-Warning "You may need to manually delete: $InstallDir"
    }
} else {
    Write-Info "Installation directory not found"
}

# 5. Remove plugin directory (optional)
if (Test-Path $PluginDir) {
    if ($KeepPlugins) {
        Write-Info "Keeping plugin directory as requested: $PluginDir"
    } else {
        Write-Info "Removing plugin directory..."
        try {
            Remove-Item -Recurse -Force $PluginDir
            Write-Success "Plugin directory removed: $PluginDir"
        } catch {
            Write-ErrorMsg "Failed to remove plugin directory: $_"
            Write-Warning "You may need to manually delete: $PluginDir"
        }
    }
} else {
    Write-Info "Plugin directory not found"
}

# 6. Clean up any remaining registry entries
Write-Info "Cleaning up registry entries..."
try {
    # Remove any Windows uninstall entries (if they exist)
    $uninstallKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\Lucien CLI"
    if (Test-Path $uninstallKey) {
        Remove-Item -Recurse -Force $uninstallKey
        Write-Success "Registry entries cleaned"
    }
} catch {
    Write-Warning "Could not clean registry entries: $_"
}

Write-Host ""
Write-Host "=================================" -ForegroundColor Green
Write-Host "LUCIEN CLI UNINSTALL COMPLETE!" -ForegroundColor Green
Write-Host "=================================" -ForegroundColor Green
Write-Host ""
Write-Host "What was removed:" -ForegroundColor Yellow
Write-Host "  - Lucien CLI binary and installation directory" -ForegroundColor White
Write-Host "  - System PATH entries" -ForegroundColor White
Write-Host "  - Desktop shortcut" -ForegroundColor White
if (-not $KeepPlugins) {
    Write-Host "  - Plugin directory and agent files" -ForegroundColor White
}
Write-Host ""
Write-Host "Note: Restart your terminal for PATH changes to take effect" -ForegroundColor Cyan
Write-Host "Thank you for using Lucien CLI!" -ForegroundColor Green