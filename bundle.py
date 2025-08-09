#!/usr/bin/env python3
"""
📦 LUCIEN CLI BUNDLER
Creates production-ready LUCIEN_TOOLKIT.zip for deployment and sharing
"""

import zipfile
import shutil
from pathlib import Path
from datetime import datetime
import json
import os

def create_lucien_bundle():
    """Create comprehensive Lucien CLI bundle"""
    
    project_root = Path.cwd()
    bundle_name = f"LUCIEN_TOOLKIT_{datetime.now().strftime('%Y%m%d_%H%M%S')}.zip"
    bundle_path = project_root / bundle_name
    
    print(f"Creating Lucien CLI Bundle: {bundle_name}")
    
    with zipfile.ZipFile(bundle_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
        
        # Core CLI source
        print("Adding core CLI source...")
        add_directory_to_zip(zipf, project_root / "cli", "lucien-cli/cli")
        add_directory_to_zip(zipf, project_root / "lucien", "lucien-cli/lucien")
        
        # Configuration files
        print("Adding configuration...")
        essential_files = [
            "requirements.txt",
        ]
        
        for file_name in essential_files:
            file_path = project_root / file_name
            if file_path.exists():
                zipf.write(file_path, f"lucien-cli/{file_name}")
        
        # Create startup README
        print("Creating startup README...")
        startup_readme = create_startup_readme()
        zipf.writestr("lucien-cli/QUICKSTART.md", startup_readme)
        
        # Create installation script
        print("Creating installation script...")
        install_script = create_install_script()
        zipf.writestr("lucien-cli/install.py", install_script)
        
        # Bundle metadata
        print("Adding bundle metadata...")
        metadata = {
            "name": "Lucien CLI Toolkit",
            "version": "1.0.0",
            "created": datetime.now().isoformat(),
            "description": "AI-Enhanced Shell Replacement with Revolutionary Features",
            "features": [
                "AI-powered shell replacement",
                "Complete shell functionality", 
                "Production safety features",
                "Cross-platform compatibility"
            ],
            "requirements": ["Python 3.11+", "Rich", "Typer", "PyYAML", "TOML"],
            "platform": "Cross-platform (Windows, macOS, Linux)"
        }
        
        zipf.writestr("lucien-cli/BUNDLE_INFO.json", json.dumps(metadata, indent=2))
    
    print(f"Bundle created successfully: {bundle_path}")
    print(f"Size: {bundle_path.stat().st_size / (1024*1024):.2f} MB")
    
    return bundle_path

def add_directory_to_zip(zipf, source_dir, archive_dir):
    """Add directory to zip file"""
    source_path = Path(source_dir)
    if not source_path.exists():
        return
    
    for file_path in source_path.rglob("*"):
        if file_path.is_file() and not should_exclude(file_path):
            relative_path = file_path.relative_to(source_path)
            archive_path = f"{archive_dir}/{relative_path}"
            zipf.write(file_path, archive_path)

def should_exclude(file_path):
    """Check if file should be excluded from bundle"""
    exclude_patterns = [
        "__pycache__",
        "*.pyc",
        "*.pyo", 
        ".git",
        ".venv",
        "venv",
        "node_modules",
        "*.log",
        "*.tmp",
        ".DS_Store",
        "Thumbs.db",
        "*.swp",
        "*.swo"
    ]
    
    file_str = str(file_path)
    return any(pattern in file_str or file_str.endswith(pattern.replace("*", "")) 
              for pattern in exclude_patterns)

def create_startup_readme():
    """Create comprehensive startup README"""
    return """# 🚀 LUCIEN CLI QUICKSTART

Welcome to Lucien CLI - The AI-enhanced command-line interface!

## ⚡ Quick Installation

### Prerequisites
- Python 3.11 or higher

### Install Steps

1. **Extract the bundle**
   ```bash
   # Extract LUCIEN_TOOLKIT.zip to your desired location
   unzip LUCIEN_TOOLKIT.zip
   cd lucien-cli
   ```

2. **Install dependencies**
   ```bash
   pip install -r requirements.txt
   ```

3. **Test installation**
   ```bash
   python cli/main.py --help
   ```

4. **Run system tests**
   ```bash
   python cli/main.py --test interactive
   ```

## 🎮 First Run Experience

### Interactive Shell Mode
```bash
# Start Lucien in interactive shell mode
python cli/main.py interactive

# Test safe mode
python cli/main.py --safe interactive
```

## 🔒 Production Safety

### Safe Mode
```bash
# Run in safe mode (blocks dangerous commands)
python cli/main.py --safe interactive
```

### Safety Features
- Dangerous command detection and blocking
- File system operation monitoring  
- Comprehensive logging
- Confirmation prompts for destructive actions

## 🆘 Troubleshooting

### Common Issues

**"Module not found" errors**
```bash
pip install -r requirements.txt
```

**Permission errors**
```bash
# On Unix systems
chmod +x cli/main.py
```

### Get Help
```bash
python cli/main.py --help           # General help
python cli/main.py interactive      # Start interactive mode
```

## 🚀 Ready to Start?

```bash
# Start your journey
python cli/main.py interactive

# Experience the AI-enhanced command line
```

---

**Lucien CLI - Where Intelligence Meets Command Line**

*The command-line interface that adapts to you.*
"""

def create_install_script():
    """Create automated installation script"""
    return """#!/usr/bin/env python3
\"\"\"
🔧 LUCIEN CLI INSTALLER
Automated installation and setup script
\"\"\"

import subprocess
import sys
import os
from pathlib import Path

def main():
    print("🚀 Lucien CLI Installer")
    print("=" * 50)
    
    # Check Python version
    if sys.version_info < (3, 11):
        print("❌ Python 3.11+ required")
        print(f"Current version: {sys.version}")
        sys.exit(1)
    
    print("✅ Python version check passed")
    
    # Install requirements
    print("📦 Installing dependencies...")
    try:
        subprocess.run([sys.executable, "-m", "pip", "install", "-r", "requirements.txt"], 
                      check=True, capture_output=True)
        print("✅ Dependencies installed successfully")
    except subprocess.CalledProcessError as e:
        print(f"❌ Failed to install dependencies: {e}")
        sys.exit(1)
    
    # Run system tests
    print("🧪 Running system tests...")
    try:
        result = subprocess.run([sys.executable, "cli/main.py", "--test", "interactive"], 
                              capture_output=True, text=True)
        if result.returncode == 0:
            print("✅ System tests passed")
        else:
            print("⚠️ Some tests failed, but installation can continue")
            print(result.stderr)
    except Exception as e:
        print(f"⚠️ Could not run tests: {e}")
    
    # Setup complete
    print("🎉 Installation Complete!")
    print()
    print("🎯 Quick Start:")
    print(f"   python cli/main.py interactive")
    print()
    print("📖 Full documentation: QUICKSTART.md")
    print("🆘 Need help? Run: python cli/main.py --help")

if __name__ == "__main__":
    main()
"""

if __name__ == "__main__":
    bundle_path = create_lucien_bundle()
    print(f"\nLUCIEN_TOOLKIT.zip ready for deployment!")
    print(f"Location: {bundle_path}")