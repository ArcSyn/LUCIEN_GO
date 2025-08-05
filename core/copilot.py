# core/copilot.py
"""
REVOLUTIONARY FEATURE 2: Real-time AI Pair Programming Copilot
Works alongside you as you code, suggests fixes, optimizations, and complete workflows
"""
import json
import time
import threading
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional, Set, Tuple, Any
import subprocess
import difflib
import re

from .config import config
from .claude_memory import memory

class AIPairProgrammer:
    """Real-time AI copilot for collaborative coding"""
    
    def __init__(self):
        self.copilot_dir = config.config_dir / "copilot"
        self.copilot_dir.mkdir(exist_ok=True)
        
        self.active_files = {}  # file_path -> file_info
        self.code_suggestions = {}  # file_path -> suggestions
        self.error_patterns = {}  # error -> suggested_fixes
        self.workflow_state = "idle"  # idle, coding, testing, debugging
        
        self.monitoring = False
        self.monitor_thread = None
        
        self._load_knowledge_base()
        
    def _load_knowledge_base(self):
        """Load AI knowledge base for common patterns and fixes"""
        self.knowledge_base = {
            "python": {
                "common_errors": {
                    "ModuleNotFoundError": [
                        "pip install {module}",
                        "Check if module is in requirements.txt",
                        "Verify virtual environment activation"
                    ],
                    "IndentationError": [
                        "Fix indentation with consistent spaces/tabs",
                        "Use 4 spaces for Python indentation",
                        "Check for mixed tabs and spaces"
                    ],
                    "ImportError": [
                        "Check module installation",
                        "Verify PYTHONPATH",
                        "Check for circular imports"
                    ]
                },
                "best_practices": [
                    "Add type hints to function parameters",
                    "Use descriptive variable names",
                    "Add docstrings to functions",
                    "Handle exceptions properly",
                    "Use list comprehensions for simple loops"
                ],
                "code_smells": {
                    r"print\(.*\)": "Consider using logging instead of print",
                    r"except:": "Use specific exception types instead of bare except",
                    r"== True|== False": "Use 'if condition:' instead of '== True'",
                    r"len\([^)]+\) == 0": "Use 'if not sequence:' instead of 'len(sequence) == 0'"
                }
            },
            "javascript": {
                "common_errors": {
                    "ReferenceError": [
                        "Check variable declaration",
                        "Verify variable scope",
                        "Check for typos in variable names"
                    ],
                    "TypeError": [
                        "Check data types",
                        "Verify object properties exist",
                        "Add null/undefined checks"
                    ]
                },
                "best_practices": [
                    "Use const/let instead of var",
                    "Add JSDoc comments",
                    "Use strict equality (===)",
                    "Handle async operations properly",
                    "Use destructuring for cleaner code"
                ]
            }
        }
    
    def start_monitoring(self, directory: Path = None):
        """Start real-time file monitoring and AI assistance"""
        if self.monitoring:
            return
            
        self.monitoring = True
        target_dir = directory or Path.cwd()
        
        self.monitor_thread = threading.Thread(
            target=self._monitor_files, 
            args=(target_dir,), 
            daemon=True
        )
        self.monitor_thread.start()
        
        return f"AI Copilot started monitoring: {target_dir}"
    
    def stop_monitoring(self):
        """Stop file monitoring"""
        self.monitoring = False
        if self.monitor_thread:
            self.monitor_thread.join(timeout=1)
        return "AI Copilot monitoring stopped"
    
    def _monitor_files(self, directory: Path):
        """Monitor files for changes and provide real-time assistance"""
        code_extensions = {'.py', '.js', '.ts', '.go', '.rs', '.java', '.cpp', '.c', '.h'}
        last_check = {}
        
        while self.monitoring:
            try:
                for file_path in directory.rglob("*"):
                    if (file_path.is_file() and 
                        file_path.suffix in code_extensions and
                        not any(part.startswith('.') for part in file_path.parts)):
                        
                        try:
                            current_mtime = file_path.stat().st_mtime
                            if str(file_path) not in last_check or last_check[str(file_path)] < current_mtime:
                                last_check[str(file_path)] = current_mtime
                                self._analyze_file_change(file_path)
                        except:
                            continue
                            
                time.sleep(2)  # Check every 2 seconds
            except Exception:
                time.sleep(5)  # Slower retry on error
    
    def _analyze_file_change(self, file_path: Path):
        """Analyze file change and provide suggestions"""
        try:
            with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
            
            # Detect language
            language = self._detect_language(file_path)
            
            # Analyze for issues and suggestions
            suggestions = []
            
            # Check for common patterns and issues
            if language in self.knowledge_base:
                kb = self.knowledge_base[language]
                
                # Check for code smells
                if "code_smells" in kb:
                    for pattern, suggestion in kb["code_smells"].items():
                        if re.search(pattern, content):
                            suggestions.append({
                                "type": "code_smell",
                                "message": suggestion,
                                "file": str(file_path),
                                "priority": "medium"
                            })
                
                # Check for incomplete patterns
                incomplete_patterns = self._detect_incomplete_code(content, language)
                suggestions.extend(incomplete_patterns)
            
            # Store suggestions
            if suggestions:
                self.code_suggestions[str(file_path)] = {
                    "timestamp": datetime.now().isoformat(),
                    "suggestions": suggestions
                }
                
                # Auto-save suggestions for later review
                self._save_suggestions_to_file(file_path, suggestions)
                
        except Exception as e:
            pass  # Silently handle file access errors
    
    def _detect_incomplete_code(self, content: str, language: str) -> List[Dict]:
        """Detect incomplete code patterns and suggest completions"""
        suggestions = []
        lines = content.split('\n')
        
        if language == "python":
            for i, line in enumerate(lines):
                stripped = line.strip()
                
                # Detect incomplete function definitions
                if stripped.startswith('def ') and not stripped.endswith(':'):
                    suggestions.append({
                        "type": "incomplete",
                        "message": f"Complete function definition on line {i+1}",
                        "line": i+1,
                        "suggestion": "Add colon (:) at end of function definition",
                        "priority": "high"
                    })
                
                # Detect incomplete try blocks
                if stripped == 'try:' and i+1 < len(lines):
                    next_lines = lines[i+1:i+5]
                    if not any('except' in l for l in next_lines):
                        suggestions.append({
                            "type": "incomplete",
                            "message": f"Try block on line {i+1} needs except clause",
                            "line": i+1,
                            "suggestion": "Add except clause to handle exceptions",
                            "priority": "high"
                        })
                
                # Detect TODO comments that could be auto-implemented
                if 'TODO' in stripped.upper():
                    todo_text = stripped.split('TODO')[-1].strip(':# ')
                    if any(keyword in todo_text.lower() for keyword in ['add', 'implement', 'create']):
                        suggestions.append({
                            "type": "todo_implementation",
                            "message": f"TODO on line {i+1} could be auto-implemented",
                            "line": i+1,
                            "todo": todo_text,
                            "priority": "low"
                        })
        
        return suggestions
    
    def _detect_language(self, file_path: Path) -> str:
        """Detect programming language from file extension"""
        ext = file_path.suffix.lower()
        language_map = {
            '.py': 'python',
            '.js': 'javascript',
            '.ts': 'typescript',
            '.go': 'go',
            '.rs': 'rust',
            '.java': 'java',
            '.cpp': 'cpp',
            '.c': 'c',
            '.h': 'c'
        }
        return language_map.get(ext, 'unknown')
    
    def get_contextual_suggestions(self, current_file: Path, cursor_line: int = None) -> List[Dict]:
        """Get AI suggestions based on current context"""
        if not current_file.exists():
            return []
        
        try:
            with open(current_file, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
            
            language = self._detect_language(current_file)
            suggestions = []
            
            # Context-aware suggestions based on cursor position
            if cursor_line:
                lines = content.split('\n')
                if cursor_line <= len(lines):
                    current_line = lines[cursor_line - 1]
                    
                    # Suggest imports for common patterns
                    if language == "python":
                        if 'requests.' in current_line and 'import requests' not in content:
                            suggestions.append({
                                "type": "missing_import",
                                "message": "Add missing import: requests",
                                "suggestion": "import requests",
                                "priority": "high"
                            })
                        
                        if 'json.' in current_line and 'import json' not in content:
                            suggestions.append({
                                "type": "missing_import", 
                                "message": "Add missing import: json",
                                "suggestion": "import json",
                                "priority": "high"
                            })
            
            # Project-wide suggestions
            project_suggestions = self._get_project_suggestions(current_file.parent)
            suggestions.extend(project_suggestions)
            
            return suggestions
            
        except Exception:
            return []
    
    def _get_project_suggestions(self, project_dir: Path) -> List[Dict]:
        """Get project-wide improvement suggestions"""
        suggestions = []
        
        # Check for missing common files
        if not (project_dir / "README.md").exists():
            suggestions.append({
                "type": "project_structure",
                "message": "Consider adding a README.md file",
                "suggestion": "Create README.md with project description",
                "priority": "medium"
            })
        
        if not (project_dir / ".gitignore").exists() and (project_dir / ".git").exists():
            suggestions.append({
                "type": "project_structure",
                "message": "Consider adding a .gitignore file",
                "suggestion": "Create .gitignore to exclude build files",
                "priority": "medium"
            })
        
        # Language-specific suggestions
        if (project_dir / "requirements.txt").exists():
            if not (project_dir / "setup.py").exists():
                suggestions.append({
                    "type": "python_project",
                    "message": "Consider adding setup.py for proper package structure",
                    "suggestion": "Create setup.py for installable package",
                    "priority": "low"
                })
        
        return suggestions
    
    def analyze_error_output(self, error_text: str, context: Dict = None) -> Dict:
        """Analyze error output and suggest fixes"""
        language = context.get('language', 'unknown') if context else 'unknown'
        
        analysis = {
            "error_type": "unknown",
            "suggested_fixes": [],
            "confidence": 0.0
        }
        
        # Parse common error patterns
        if language in self.knowledge_base:
            kb = self.knowledge_base[language]
            
            for error_type, fixes in kb.get("common_errors", {}).items():
                if error_type in error_text:
                    analysis["error_type"] = error_type
                    analysis["suggested_fixes"] = fixes
                    analysis["confidence"] = 0.8
                    break
        
        # Extract specific information from error
        if "ModuleNotFoundError" in error_text:
            # Extract module name
            module_match = re.search(r"No module named '([^']+)'", error_text)
            if module_match:
                module_name = module_match.group(1)
                analysis["suggested_fixes"] = [
                    f"pip install {module_name}",
                    f"Add {module_name} to requirements.txt",
                    "Check if you're in the correct virtual environment"
                ]
                analysis["confidence"] = 0.9
        
        return analysis
    
    def suggest_workflow_next_step(self, current_context: Dict) -> Optional[str]:
        """Suggest next step in development workflow"""
        cwd = Path.cwd()
        
        # Detect current workflow state
        if (cwd / ".git").exists():
            # Check git status
            try:
                result = subprocess.run(
                    ["git", "status", "--porcelain"], 
                    capture_output=True, text=True, timeout=5
                )
                if result.returncode == 0:
                    if result.stdout.strip():
                        return "You have uncommitted changes. Consider: git add . && git commit -m 'Your message'"
                    else:
                        return "Working tree clean. Consider: git pull to sync with remote"
            except:
                pass
        
        # Check for testing
        if any((cwd / f).exists() for f in ["test", "tests", "spec"]):
            return "Run tests to ensure code quality: pytest or npm test"
        
        # Check for build requirements
        if (cwd / "package.json").exists():
            return "Consider running: npm install && npm run build"
        elif (cwd / "requirements.txt").exists():
            return "Consider running: pip install -r requirements.txt"
        
        return None
    
    def _save_suggestions_to_file(self, file_path: Path, suggestions: List[Dict]):
        """Save suggestions to copilot suggestions file"""
        suggestions_file = self.copilot_dir / f"{file_path.name}_suggestions.json"
        
        try:
            with open(suggestions_file, 'w') as f:
                json.dump({
                    "file": str(file_path),
                    "timestamp": datetime.now().isoformat(),
                    "suggestions": suggestions
                }, f, indent=2)
        except:
            pass
    
    def get_all_suggestions(self) -> Dict[str, List[Dict]]:
        """Get all current suggestions"""
        return self.code_suggestions.copy()
    
    def clear_suggestions(self, file_path: str = None):
        """Clear suggestions for specific file or all files"""
        if file_path:
            self.code_suggestions.pop(file_path, None)
        else:
            self.code_suggestions.clear()

# Global copilot instance
copilot = AIPairProgrammer()