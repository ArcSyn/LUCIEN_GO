# 🚀 LUCIEN CLI INSTALLATION GUIDE

Complete installation and uninstallation instructions for all platforms.

---

## ⚡ **QUICK INSTALL (RECOMMENDED)**

### **Option 1: Make Install (Easiest)**
```bash
# Windows (PowerShell)
cd "C:\Users\lcolo\OneDrive\Desktop\LUCIEN CLI"
make install

# Unix/Linux/macOS (Terminal) 
cd "/path/to/LUCIEN CLI"
make install
```

### **Option 2: Direct Script**
```powershell
# Windows
cd "C:\Users\lcolo\OneDrive\Desktop\LUCIEN CLI"
.\scripts\install.ps1

# Unix/Linux/macOS
cd "/path/to/LUCIEN CLI"
bash scripts/install.sh
```

---

## 🔧 **DETAILED INSTALLATION OPTIONS**

### **Windows Installation**

#### **Method 1: Full Installation (Recommended)**
```powershell
# Navigate to project directory
cd "C:\Users\lcolo\OneDrive\Desktop\LUCIEN CLI"

# Build and install in one command
make install
```

#### **Method 2: Manual PowerShell Script**
```powershell
# Navigate to project directory
cd "C:\Users\lcolo\OneDrive\Desktop\LUCIEN CLI"

# If you get execution policy error, run this first:
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# Run the installer
.\scripts\install.ps1
```

#### **Method 3: Development Install (No Plugins)**
```powershell
make install-dev
```

#### **Method 4: Force Reinstall**
```powershell
make install-force  # Removes old version, then installs fresh
```

### **Unix/Linux/macOS Installation**

#### **Method 1: Full Installation (Recommended)**
```bash
cd "/path/to/LUCIEN CLI"
make install
```

#### **Method 2: Manual Bash Script**
```bash
cd "/path/to/LUCIEN CLI"
bash scripts/install.sh
```

#### **Method 3: Development Install (No Plugins)**
```bash
make install-dev
```

---

## 📋 **WHAT GETS INSTALLED**

### **Full Installation Includes:**
- ✅ **Lucien CLI Binary** → Installed to system directories
- ✅ **Python Agent Plugins** → 4 AI agents (plan, design, review, code)  
- ✅ **System PATH Integration** → Run `lucien` from anywhere
- ✅ **Desktop Shortcut** → Windows only
- ✅ **Shell Profile Updates** → Unix/Linux/macOS
- ✅ **Plugin Directory Setup** → `~/.lucien/plugins/`

### **Installation Locations:**
- **Windows**: `%LOCALAPPDATA%\Lucien\` + `%USERPROFILE%\.lucien\plugins\`
- **Unix/Linux/macOS**: `~/.local/bin/` + `~/.lucien/plugins/`

---

## 🗑️ **UNINSTALLATION**

### **Complete Removal**
```bash
# Remove everything
make uninstall
```

### **Keep Plugins (Remove Binary Only)**
```powershell
# Windows
.\scripts\uninstall.ps1 -KeepPlugins

# Unix/Linux/macOS  
bash scripts/uninstall.sh --keep-plugins
```

### **What Gets Removed:**
- ❌ Lucien CLI binary and installation directory
- ❌ System PATH entries  
- ❌ Desktop shortcuts
- ❌ Plugin directory and agent files (unless `-KeepPlugins`)
- ❌ Shell profile entries
- ❌ Configuration files

---

## 🛠️ **ADVANCED INSTALLATION OPTIONS**

### **Custom Installation Directory**
```powershell
# Windows
.\scripts\install.ps1 -InstallDir "C:\Tools\Lucien"

# Unix/Linux/macOS
bash scripts/install.sh --install-dir "/opt/lucien"
```

### **Skip Components**
```powershell
# Skip desktop shortcut
.\scripts\install.ps1 -NoShortcut

# Skip PATH modification
.\scripts\install.ps1 -NoPath  

# Skip Python check
.\scripts\install.ps1 -SkipPython

# Combine options
.\scripts\install.ps1 -NoShortcut -SkipPython -Verbose
```

### **Advanced Bootstrap Installer**
```bash
# Full-featured Go installer with more options
make bootstrap
```

---

## 🔍 **TROUBLESHOOTING**

### **Common Issues & Fixes**

#### **PowerShell Execution Policy Error**
```powershell
# Error: "execution of scripts is disabled on this system"
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
```

#### **Binary Not Found Error**
```bash
# Error: "Lucien binary not found in ./build folder"
make build          # Build the binary first
make install        # Then install
```

#### **Python Not Found Warning**
```bash
# Install Python 3.7+ from python.org
# Or skip Python check:
make install ARGS="-SkipPython"
```

#### **Permission Denied (Unix/Linux)**
```bash
# Make sure install script is executable
chmod +x scripts/install.sh
bash scripts/install.sh
```

### **Verification Commands**
```bash
# Check if installed correctly
lucien --version
lucien help

# Check installation locations
ls ~/.local/bin/lucien          # Unix/Linux/macOS
dir "%LOCALAPPDATA%\Lucien\"    # Windows
```

---

## 🎯 **QUICK REFERENCE**

| **Command** | **Description** |
|-------------|----------------|
| `make install` | Full installation (recommended) |
| `make uninstall` | Complete removal |
| `make install-force` | Force reinstall |
| `make install-dev` | Development install only |
| `make build` | Build binary only |
| `lucien --version` | Verify installation |

---

## ✅ **POST-INSTALLATION**

After successful installation:

1. **Restart your terminal** for PATH changes to take effect
2. **Verify installation**: `lucien --version`
3. **Start using Lucien**: `lucien`
4. **Try AI agents**: 
   - `plan "build a web app"`
   - `design "dark login form"`
   - `review myfile.py`
   - `code generate "sort function"`

### **Getting Help**
```bash
lucien help                    # Built-in help
cat README_AGENTS.md          # Agent commands guide  
cat POWERSHELL_MAPPING.md     # Command reference
```

---

**🎉 You're all set! Welcome to Lucien CLI!** 🚀