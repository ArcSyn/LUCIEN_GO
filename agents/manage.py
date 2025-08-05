from agents.agent_base import Agent
import os
import json
from pathlib import Path

class Manage(Agent):
    """BMAD Manage Agent - Environment and dependency management"""
    
    def run(self, input_text: str) -> str:
        """
        Manage operations: env vars, dependencies, configs
        Examples:
        - "show environment"
        - "manage dependencies"
        - "check system info"
        """
        input_lower = input_text.lower().strip()
        
        if "environment" in input_lower or "env" in input_lower:
            return self._show_environment()
        elif "dependencies" in input_lower or "deps" in input_lower:
            return self._manage_dependencies()
        elif "system" in input_lower or "info" in input_lower:
            return self._show_system_info()
        elif "config" in input_lower:
            return self._show_config()
        else:
            return self._general_management(input_text)
    
    def _show_environment(self) -> str:
        """Show environment variables and system info"""
        important_vars = ['PATH', 'PYTHONPATH', 'HOME', 'USER', 'SHELL', 'EDITOR']
        result = "[Manage] Environment Variables:\n"
        
        for var in important_vars:
            value = os.environ.get(var, "Not set")
            if len(value) > 100:
                value = value[:100] + "..."
            result += f"  {var} = {value}\n"
            
        # Current directory
        result += f"\nCurrent Directory: {Path.cwd()}\n"
        result += f"Python Version: {os.sys.version.split()[0]}\n"
        result += f"Platform: {os.name}\n"
        
        return result
    
    def _manage_dependencies(self) -> str:
        """Check and manage project dependencies"""
        cwd = Path.cwd()
        result = "[Manage] Dependency Status:\n"
        
        # Check different project types
        if (cwd / "requirements.txt").exists():
            result += "  Python project detected (requirements.txt)\n"
            try:
                with open(cwd / "requirements.txt") as f:
                    deps = f.read().splitlines()
                result += f"  Found {len(deps)} Python dependencies\n"
            except:
                result += "  Could not read requirements.txt\n"
                
        if (cwd / "package.json").exists():
            result += "  Node.js project detected (package.json)\n"
            try:
                with open(cwd / "package.json") as f:
                    package_data = json.load(f)
                deps = package_data.get("dependencies", {})
                dev_deps = package_data.get("devDependencies", {})
                result += f"  Found {len(deps)} dependencies, {len(dev_deps)} dev dependencies\n"
            except:
                result += "  Could not read package.json\n"
                
        if (cwd / "go.mod").exists():
            result += "  Go project detected (go.mod)\n"
            
        if (cwd / "Cargo.toml").exists():
            result += "  Rust project detected (Cargo.toml)\n"
            
        if not any((cwd / f).exists() for f in ["requirements.txt", "package.json", "go.mod", "Cargo.toml"]):
            result += "  No recognized dependency files found\n"
            
        return result
    
    def _show_system_info(self) -> str:
        """Show system information"""
        import platform
        import shutil
        
        result = "[Manage] System Information:\n"
        result += f"  OS: {platform.platform()}\n"
        result += f"  Python: {platform.python_version()}\n"
        result += f"  Architecture: {platform.machine()}\n"
        result += f"  Processor: {platform.processor()}\n"
        
        # Check available tools
        tools = ['git', 'npm', 'pip', 'go', 'cargo', 'docker']
        result += "\nAvailable Tools:\n"
        for tool in tools:
            if shutil.which(tool):
                result += f"  {tool}: Available\n"
            else:
                result += f"  {tool}: Not found\n"
                
        return result
    
    def _show_config(self) -> str:
        """Show configuration information"""
        result = "[Manage] Configuration Status:\n"
        
        # Check for common config files
        config_files = ['.gitconfig', '.bashrc', '.zshrc', '.vimrc', 'config.yaml']
        
        for config_file in config_files:
            if Path.home().joinpath(config_file).exists():
                result += f"  ~/{config_file}: Exists\n"
            else:
                result += f"  ~/{config_file}: Not found\n"
                
        # Check Lucien config
        lucien_config = Path.home() / ".lucien" / "config.yaml"
        if lucien_config.exists():
            result += f"  Lucien config: {lucien_config}\n"
        else:
            result += "  Lucien config: Not initialized\n"
            
        return result
    
    def _general_management(self, query: str) -> str:
        """Handle general management queries"""
        return f"[Manage] Management for '{query}':\nUse specific commands like 'show environment', 'manage dependencies', 'check system info', or 'show config'"