# core/config.py
import json
import os
from pathlib import Path
from typing import Dict, Any, List, Optional
import yaml

class LucienConfig:
    """Lucien configuration system with .lucienrc support"""
    
    def __init__(self):
        self.config_dir = Path.home() / ".lucien"
        self.config_file = self.config_dir / "config.yaml"
        self.history_file = self.config_dir / "history.txt"
        self.memory_dir = self.config_dir / "memory"
        
        self.config = self._load_default_config()
        self._ensure_directories()
        self._load_config()
        
    def _load_default_config(self) -> Dict[str, Any]:
        """Load default configuration"""
        return {
            "shell": {
                "prompt": "lucien> ",
                "history_size": 1000,
                "timeout": 30,
                "auto_cd": True
            },
            "aliases": {
                "ll": "ls -la",
                "la": "ls -A",
                "l": "ls -CF",
                "grep": "grep --color=auto",
                "fgrep": "fgrep --color=auto",
                "egrep": "egrep --color=auto"
            },
            "environment": {
                "LUCIEN_HOME": str(self.config_dir),
                "EDITOR": "nano"
            },
            "agents": {
                "claude_api_key": "",
                "default_model": "claude-3-sonnet-20240229",
                "max_tokens": 4000,
                "temperature": 0.7
            },
            "ui": {
                "theme": "default",
                "show_banner": True,
                "colors": {
                    "primary": "blue",
                    "success": "green", 
                    "error": "red",
                    "warning": "yellow"
                }
            },
            "bmad": {
                "enabled": True,
                "auto_build": False,
                "auto_analyze": True,
                "memory_persistence": True
            },
            "startup_commands": [
                "echo 'Welcome to Lucien CLI'",
                "pwd"
            ],
            "keybindings": {
                "ctrl+r": "history_search",
                "ctrl+l": "clear",
                "tab": "autocomplete"
            }
        }
    
    def _ensure_directories(self):
        """Create necessary directories"""
        self.config_dir.mkdir(exist_ok=True)
        self.memory_dir.mkdir(exist_ok=True)
        
        # Create subdirectories
        (self.memory_dir / "claude").mkdir(exist_ok=True)
        (self.memory_dir / "agents").mkdir(exist_ok=True)
        (self.memory_dir / "sessions").mkdir(exist_ok=True)
    
    def _load_config(self):
        """Load configuration from file"""
        if self.config_file.exists():
            try:
                with open(self.config_file, 'r') as f:
                    user_config = yaml.safe_load(f) or {}
                self._merge_config(self.config, user_config)
            except Exception as e:
                print(f"Warning: Failed to load config: {e}")
    
    def _merge_config(self, base: Dict[str, Any], override: Dict[str, Any]):
        """Recursively merge configuration dictionaries"""
        for key, value in override.items():
            if key in base and isinstance(base[key], dict) and isinstance(value, dict):
                self._merge_config(base[key], value)
            else:
                base[key] = value
    
    def save_config(self):
        """Save current configuration to file"""
        try:
            with open(self.config_file, 'w') as f:
                yaml.dump(self.config, f, default_flow_style=False, indent=2)
        except Exception as e:
            print(f"Error saving config: {e}")
    
    def get(self, key: str, default: Any = None) -> Any:
        """Get configuration value using dot notation"""
        keys = key.split('.')
        value = self.config
        
        for k in keys:
            if isinstance(value, dict) and k in value:
                value = value[k]
            else:
                return default
                
        return value
    
    def set(self, key: str, value: Any):
        """Set configuration value using dot notation"""
        keys = key.split('.')
        config = self.config
        
        for k in keys[:-1]:
            if k not in config:
                config[k] = {}
            config = config[k]
            
        config[keys[-1]] = value
    
    def get_aliases(self) -> Dict[str, str]:
        """Get all command aliases"""
        return self.config.get("aliases", {})
    
    def add_alias(self, name: str, command: str):
        """Add command alias"""
        if "aliases" not in self.config:
            self.config["aliases"] = {}
        self.config["aliases"][name] = command
        
    def get_environment(self) -> Dict[str, str]:
        """Get environment variables"""
        return self.config.get("environment", {})
    
    def set_environment(self, name: str, value: str):
        """Set environment variable"""
        if "environment" not in self.config:
            self.config["environment"] = {}
        self.config["environment"][name] = value
    
    def get_startup_commands(self) -> List[str]:
        """Get commands to run at startup"""
        return self.config.get("startup_commands", [])
    
    def add_startup_command(self, command: str):
        """Add startup command"""
        if "startup_commands" not in self.config:
            self.config["startup_commands"] = []
        self.config["startup_commands"].append(command)
    
    def load_history(self) -> List[str]:
        """Load command history"""
        if self.history_file.exists():
            try:
                with open(self.history_file, 'r') as f:
                    return f.read().splitlines()
            except Exception:
                pass
        return []
    
    def save_history(self, history: List[str]):
        """Save command history"""
        try:
            max_history = self.get("shell.history_size", 1000)
            history = history[-max_history:]  # Keep only recent entries
            
            with open(self.history_file, 'w') as f:
                f.write('\n'.join(history))
        except Exception as e:
            print(f"Error saving history: {e}")
    
    def get_claude_memory_file(self, session_id: str = "default") -> Path:
        """Get Claude memory file path"""
        return self.memory_dir / "claude" / f"{session_id}.md"
    
    def get_agent_memory_file(self, agent_name: str) -> Path:
        """Get agent memory file path"""
        return self.memory_dir / "agents" / f"{agent_name}.json"
    
    def create_default_config(self):
        """Create default configuration file"""
        self.save_config()
        print(f"Created default configuration at {self.config_file}")
        
    def show_config(self) -> str:
        """Show current configuration as YAML"""
        return yaml.dump(self.config, default_flow_style=False, indent=2)
    
    def reset_config(self):
        """Reset to default configuration"""
        self.config = self._load_default_config()
        self.save_config()
        print("Configuration reset to defaults")

# Global config instance
config = LucienConfig()