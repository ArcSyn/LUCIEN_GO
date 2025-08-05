# core/intelligence.py
"""
REVOLUTIONARY FEATURE 1: Intelligent Command Prediction & Auto-Execution
Learns your patterns and predicts/executes commands before you finish typing
"""
import json
import time
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from collections import defaultdict, Counter
import difflib

from .config import config
from .claude_memory import memory

class IntelligentPredictor:
    """AI-powered command prediction and auto-execution system"""
    
    def __init__(self):
        self.patterns_file = config.config_dir / "intelligence" / "patterns.json"
        self.patterns_file.parent.mkdir(exist_ok=True)
        
        self.command_patterns = self._load_patterns()
        self.session_commands = []
        self.context_stack = []  # Track current context (directory, project type, etc.)
        
    def _load_patterns(self) -> Dict:
        """Load learned command patterns"""
        if self.patterns_file.exists():
            try:
                with open(self.patterns_file, 'r') as f:
                    return json.load(f)
            except:
                pass
        
        return {
            "sequences": defaultdict(lambda: defaultdict(int)),  # cmd1 -> {cmd2: count}
            "context_commands": defaultdict(lambda: defaultdict(int)),  # context -> {cmd: count}
            "time_patterns": defaultdict(lambda: defaultdict(int)),  # hour -> {cmd: count}
            "project_patterns": defaultdict(lambda: defaultdict(int)),  # project_type -> {cmd: count}
            "frequency": defaultdict(int),  # cmd -> total_count
            "recent_boosts": {},  # cmd -> boost_factor (for recently successful predictions)
        }
    
    def _save_patterns(self):
        """Save learned patterns"""
        try:
            # Convert defaultdicts to regular dicts for JSON serialization
            save_data = {}
            for key, value in self.command_patterns.items():
                if isinstance(value, defaultdict):
                    save_data[key] = {k: dict(v) if isinstance(v, defaultdict) else v 
                                    for k, v in value.items()}
                else:
                    save_data[key] = dict(value) if isinstance(value, defaultdict) else value
            
            with open(self.patterns_file, 'w') as f:
                json.dump(save_data, f, indent=2)
        except Exception as e:
            print(f"Warning: Could not save intelligence patterns: {e}")
    
    def learn_command(self, command: str, context: Dict = None):
        """Learn from executed command"""
        if not command.strip() or command.startswith('#'):
            return
            
        cmd_clean = command.strip().split()[0]  # Get base command
        current_time = datetime.now()
        
        # Update frequency
        self.command_patterns["frequency"][cmd_clean] += 1
        
        # Learn sequences (what command typically follows what)
        if self.session_commands:
            prev_cmd = self.session_commands[-1].split()[0]
            self.command_patterns["sequences"][prev_cmd][cmd_clean] += 1
        
        # Learn time patterns
        hour = current_time.hour
        self.command_patterns["time_patterns"][hour][cmd_clean] += 1
        
        # Learn context patterns
        if context:
            for ctx_key, ctx_value in context.items():
                self.command_patterns["context_commands"][f"{ctx_key}:{ctx_value}"][cmd_clean] += 1
        
        # Learn project patterns
        cwd = Path.cwd()
        project_type = self._detect_project_type(cwd)
        if project_type:
            self.command_patterns["project_patterns"][project_type][cmd_clean] += 1
        
        # Track session
        self.session_commands.append(command)
        if len(self.session_commands) > 50:  # Keep last 50 commands
            self.session_commands = self.session_commands[-50:]
        
        # Save every 10 commands
        if len(self.session_commands) % 10 == 0:
            self._save_patterns()
    
    def predict_next_commands(self, partial_input: str = "", context: Dict = None, limit: int = 5) -> List[Tuple[str, float]]:
        """Predict next commands with confidence scores"""
        predictions = {}
        
        # If we have partial input, find similar commands
        if partial_input.strip():
            for cmd in self.command_patterns["frequency"]:
                if cmd.startswith(partial_input.strip()):
                    score = self.command_patterns["frequency"][cmd] * 0.8
                    predictions[cmd] = predictions.get(cmd, 0) + score
        
        # Sequence-based predictions
        if self.session_commands:
            last_cmd = self.session_commands[-1].split()[0]
            for next_cmd, count in self.command_patterns["sequences"][last_cmd].items():
                score = count * 0.6
                predictions[next_cmd] = predictions.get(next_cmd, 0) + score
        
        # Time-based predictions
        current_hour = datetime.now().hour
        for cmd, count in self.command_patterns["time_patterns"][current_hour].items():
            score = count * 0.3
            predictions[cmd] = predictions.get(cmd, 0) + score
        
        # Project-based predictions
        cwd = Path.cwd()
        project_type = self._detect_project_type(cwd)
        if project_type:
            for cmd, count in self.command_patterns["project_patterns"][project_type].items():
                score = count * 0.5
                predictions[cmd] = predictions.get(cmd, 0) + score
        
        # Context-based predictions
        if context:
            for ctx_key, ctx_value in context.items():
                ctx_str = f"{ctx_key}:{ctx_value}"
                for cmd, count in self.command_patterns["context_commands"][ctx_str].items():
                    score = count * 0.4
                    predictions[cmd] = predictions.get(cmd, 0) + score
        
        # Apply recent success boosts
        for cmd in predictions:
            if cmd in self.command_patterns["recent_boosts"]:
                predictions[cmd] *= self.command_patterns["recent_boosts"][cmd]
        
        # Sort by confidence and return top predictions
        sorted_predictions = sorted(predictions.items(), key=lambda x: x[1], reverse=True)
        
        # Normalize scores to 0-1 range
        if sorted_predictions:
            max_score = sorted_predictions[0][1]
            normalized = [(cmd, score/max_score) for cmd, score in sorted_predictions[:limit]]
            return normalized
        
        return []
    
    def get_intelligent_suggestions(self, current_dir: Path, recent_commands: List[str]) -> List[str]:
        """Get contextually intelligent command suggestions"""
        suggestions = []
        
        # Analyze current situation
        files = list(current_dir.glob("*"))
        project_type = self._detect_project_type(current_dir)
        
        # Git repository suggestions
        if (current_dir / ".git").exists():
            if not any("git status" in cmd for cmd in recent_commands[-5:]):
                suggestions.append("git status")
            if not any("git pull" in cmd for cmd in recent_commands[-10:]):
                suggestions.append("git pull")
        
        # Project-specific suggestions
        if project_type == "python":
            if (current_dir / "requirements.txt").exists():
                suggestions.append("pip install -r requirements.txt")
            if any(f.name == "setup.py" for f in files):
                suggestions.append("python setup.py develop")
        
        elif project_type == "node":
            if (current_dir / "package.json").exists():
                suggestions.append("npm install")
                suggestions.append("npm start")
                suggestions.append("npm run dev")
        
        # File-based suggestions
        if any(f.suffix == ".py" for f in files):
            suggestions.append("python main.py")
            suggestions.append("pytest")
        
        if any(f.name == "Makefile" for f in files):
            suggestions.append("make")
            suggestions.append("make build")
        
        if any(f.name == "docker-compose.yml" for f in files):
            suggestions.append("docker-compose up")
        
        return suggestions[:5]
    
    def auto_execute_if_confident(self, prediction: str, confidence: float, threshold: float = 0.9) -> bool:
        """Auto-execute command if confidence is high enough"""
        if confidence >= threshold:
            # Additional safety checks
            safe_commands = {
                'ls', 'dir', 'pwd', 'git status', 'git log', 'cat', 'head', 'tail',
                'ps', 'df', 'free', 'uptime', 'whoami', 'date', 'which', 'type'
            }
            
            base_cmd = prediction.split()[0]
            if base_cmd in safe_commands:
                return True
        
        return False
    
    def boost_successful_prediction(self, command: str, boost_factor: float = 1.5):
        """Boost confidence for successful predictions"""
        base_cmd = command.split()[0]
        self.command_patterns["recent_boosts"][base_cmd] = boost_factor
        
        # Decay boosts over time
        current_time = time.time()
        expired_boosts = []
        for cmd, (boost, timestamp) in self.command_patterns.get("boost_timestamps", {}).items():
            if current_time - timestamp > 3600:  # 1 hour expiry
                expired_boosts.append(cmd)
        
        for cmd in expired_boosts:
            self.command_patterns["recent_boosts"].pop(cmd, None)
    
    def _detect_project_type(self, path: Path) -> Optional[str]:
        """Detect project type from directory contents"""
        if (path / "package.json").exists():
            return "node"
        elif (path / "requirements.txt").exists() or (path / "setup.py").exists():
            return "python"
        elif (path / "go.mod").exists():
            return "go"
        elif (path / "Cargo.toml").exists():
            return "rust"
        elif (path / "pom.xml").exists():
            return "java"
        elif (path / "Makefile").exists():
            return "c/cpp"
        return None
    
    def get_smart_completions(self, partial: str) -> List[str]:
        """Get smart tab completions based on learning"""
        completions = []
        
        # Command completions
        for cmd in self.command_patterns["frequency"]:
            if cmd.startswith(partial):
                completions.append(cmd)
        
        # Add file completions
        try:
            current_dir = Path.cwd()
            for item in current_dir.glob(f"{partial}*"):
                completions.append(str(item.name))
        except:
            pass
        
        # Sort by frequency/relevance
        def sort_key(item):
            base_cmd = item.split()[0] if ' ' in item else item
            return self.command_patterns["frequency"].get(base_cmd, 0)
        
        return sorted(set(completions), key=sort_key, reverse=True)[:10]

# Global predictor instance
predictor = IntelligentPredictor()