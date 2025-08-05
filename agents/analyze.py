from agents.agent_base import Agent
import os
from pathlib import Path
import subprocess

class Analyze(Agent):
    """BMAD Analyze Agent - Code and project analysis"""
    
    def run(self, input_text: str) -> str:
        """
        Analyze operations: code, files, projects
        Examples:
        - "analyze project"
        - "count code lines"  
        - "find issues"
        """
        input_lower = input_text.lower().strip()
        
        if "project" in input_lower:
            return self._analyze_project()
        elif "lines" in input_lower or "loc" in input_lower:
            return self._count_lines()
        elif "files" in input_lower:
            return self._analyze_files()
        elif "issues" in input_lower or "problems" in input_lower:
            return self._find_issues()
        else:
            return self._general_analysis(input_text)
    
    def _analyze_project(self) -> str:
        """Comprehensive project analysis"""
        cwd = Path.cwd()
        result = f"[Analyze] Project Analysis: {cwd.name}\n"
        result += f"Location: {cwd}\n\n"
        
        # File statistics
        total_files = 0
        file_types = {}
        
        for file_path in cwd.rglob("*"):
            if file_path.is_file() and not any(part.startswith('.') for part in file_path.parts):
                total_files += 1
                ext = file_path.suffix or "no_extension"
                file_types[ext] = file_types.get(ext, 0) + 1
        
        result += f"Files: {total_files} total\n"
        
        # Top file types
        if file_types:
            sorted_types = sorted(file_types.items(), key=lambda x: x[1], reverse=True)[:10]
            result += "File types:\n"
            for ext, count in sorted_types:
                result += f"  {ext}: {count}\n"
        
        # Project type detection
        project_types = []
        if (cwd / "package.json").exists(): project_types.append("Node.js")
        if (cwd / "requirements.txt").exists(): project_types.append("Python")
        if (cwd / "go.mod").exists(): project_types.append("Go") 
        if (cwd / "Cargo.toml").exists(): project_types.append("Rust")
        if (cwd / ".git").exists(): project_types.append("Git")
        
        if project_types:
            result += f"\nProject types: {', '.join(project_types)}\n"
        
        return result
    
    def _count_lines(self) -> str:
        """Count lines of code"""
        cwd = Path.cwd()
        code_extensions = {'.py', '.js', '.ts', '.go', '.rs', '.java', '.cpp', '.c', '.h'}
        
        total_lines = 0
        code_files = 0
        
        for file_path in cwd.rglob("*"):
            if (file_path.is_file() and 
                file_path.suffix in code_extensions and
                not any(part.startswith('.') for part in file_path.parts)):
                
                try:
                    with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                        lines = len(f.readlines())
                        total_lines += lines
                        code_files += 1
                except:
                    continue
        
        result = f"[Analyze] Line Count:\n"
        result += f"  Code files: {code_files}\n"
        result += f"  Total lines: {total_lines}\n"
        
        if code_files > 0:
            result += f"  Average: {total_lines // code_files} lines per file\n"
            
        return result
    
    def _analyze_files(self) -> str:
        """Analyze file structure"""
        cwd = Path.cwd()
        result = "[Analyze] File Structure:\n"
        
        # Directory structure
        dirs = [d for d in cwd.iterdir() if d.is_dir() and not d.name.startswith('.')]
        files = [f for f in cwd.iterdir() if f.is_file() and not f.name.startswith('.')]
        
        result += f"  Directories: {len(dirs)}\n"
        result += f"  Files: {len(files)}\n"
        
        # Show main directories
        if dirs:
            result += "\nMain directories:\n"
            for dir_path in sorted(dirs)[:10]:
                try:
                    file_count = len(list(dir_path.rglob("*")))
                    result += f"  {dir_path.name}/ ({file_count} items)\n"
                except:
                    result += f"  {dir_path.name}/\n"
        
        # Show main files
        if files:
            result += "\nMain files:\n"
            for file_path in sorted(files)[:10]:
                try:
                    size = file_path.stat().st_size
                    if size > 1024:
                        size_str = f"{size // 1024}KB"
                    else:
                        size_str = f"{size}B"
                    result += f"  {file_path.name} ({size_str})\n"
                except:
                    result += f"  {file_path.name}\n"
        
        return result
    
    def _find_issues(self) -> str:
        """Find potential issues in code"""
        cwd = Path.cwd()
        result = "[Analyze] Issue Detection:\n"
        
        issues = []
        
        # Check for common issues
        for file_path in cwd.rglob("*.py"):
            if any(part.startswith('.') for part in file_path.parts):
                continue
                
            try:
                with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                    content = f.read()
                    lines = content.split('\n')
                    
                for line_num, line in enumerate(lines, 1):
                    line_lower = line.lower()
                    
                    # Check for common issues
                    if 'todo' in line_lower or 'fixme' in line_lower:
                        issues.append(f"{file_path.name}:{line_num} - TODO/FIXME found")
                    elif 'password' in line_lower and '=' in line:
                        issues.append(f"{file_path.name}:{line_num} - Possible hardcoded password")
                    elif 'eval(' in line:
                        issues.append(f"{file_path.name}:{line_num} - Dangerous eval() usage")
                        
            except:
                continue
                
        if issues:
            result += f"Found {len(issues)} potential issues:\n"
            for issue in issues[:10]:  # Show first 10
                result += f"  {issue}\n"
            if len(issues) > 10:
                result += f"  ... and {len(issues) - 10} more\n"
        else:
            result += "No obvious issues found\n"
            
        return result
    
    def _general_analysis(self, query: str) -> str:
        """Handle general analysis queries"""
        return f"[Analyze] Analysis for '{query}':\nUse specific commands like 'analyze project', 'count lines', 'analyze files', or 'find issues'"