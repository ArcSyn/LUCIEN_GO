from agents.agent_base import Agent
import os
import subprocess
from pathlib import Path
import json

class AnalyzeAgent(Agent):
    """BMAD Analyze Agent - Code analysis, metrics, project insights"""
    
    def run(self, input_text: str) -> str:
        """
        Analysis operations:
        - "analyze codebase"
        - "count lines" 
        - "find todos"
        - "security scan"
        - "dependency check"
        """
        input_lower = input_text.lower().strip()
        
        if "codebase" in input_lower or "analyze" in input_lower:
            return self._analyze_codebase()
        elif "lines" in input_lower or "loc" in input_lower:
            return self._count_lines()
        elif "todo" in input_lower or "fixme" in input_lower:
            return self._find_todos()
        elif "security" in input_lower or "vuln" in input_lower:
            return self._security_scan()
        elif "dep" in input_lower and ("check" in input_lower or "audit" in input_lower):
            return self._dependency_audit()
        elif "git" in input_lower:
            return self._git_analysis()
        else:
            return self._general_analysis(input_text)
    
    def _analyze_codebase(self) -> str:
        """Comprehensive codebase analysis"""
        cwd = Path.cwd()
        results = []
        
        # Basic project info
        results.append(f"📁 Project: {cwd.name}")
        results.append(f"📍 Path: {cwd}")
        
        # File counts by extension
        file_stats = {}
        total_files = 0
        
        for file_path in cwd.rglob("*"):
            if file_path.is_file() and not any(part.startswith('.') for part in file_path.parts):
                total_files += 1
                ext = file_path.suffix or "no_extension"
                file_stats[ext] = file_stats.get(ext, 0) + 1
        
        results.append(f"📄 Total files: {total_files}")
        
        # Top file types
        if file_stats:
            results.append("📊 File types:")
            sorted_stats = sorted(file_stats.items(), key=lambda x: x[1], reverse=True)[:10]
            for ext, count in sorted_stats:
                results.append(f"   {ext}: {count} files")
        
        # Project type detection
        project_types = []
        if (cwd / "package.json").exists():
            project_types.append("Node.js")
        if (cwd / "requirements.txt").exists() or (cwd / "setup.py").exists():
            project_types.append("Python")
        if (cwd / "go.mod").exists():
            project_types.append("Go")
        if (cwd / "Cargo.toml").exists():
            project_types.append("Rust")
        if (cwd / ".git").exists():
            project_types.append("Git")
            
        if project_types:
            results.append(f"🔧 Project types: {', '.join(project_types)}")
        
        # Add line count
        loc_result = self._count_lines()
        results.append(loc_result.split('\n', 1)[1] if '\n' in loc_result else loc_result)
        
        return f"[AnalyzeAgent] Codebase Analysis:\n" + "\n".join(results)
    
    def _count_lines(self) -> str:
        """Count lines of code"""
        cwd = Path.cwd()
        code_extensions = {'.py', '.js', '.ts', '.go', '.rs', '.java', '.cpp', '.c', '.h', '.css', '.html'}
        
        total_lines = 0
        code_lines = 0
        files_counted = 0
        
        for file_path in cwd.rglob("*"):
            if (file_path.is_file() and 
                file_path.suffix in code_extensions and 
                not any(part.startswith('.') for part in file_path.parts)):
                
                try:
                    with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                        lines = f.readlines()
                        total_lines += len(lines)
                        # Count non-empty, non-comment lines as code
                        code_lines += sum(1 for line in lines if line.strip() and not line.strip().startswith('#'))
                        files_counted += 1
                except Exception:
                    continue
        
        return f"[AnalyzeAgent] Line Count:\n📄 Files: {files_counted}\n📏 Total lines: {total_lines}\n💻 Code lines: {code_lines}"
    
    def _find_todos(self) -> str:
        """Find TODO, FIXME, HACK comments"""
        cwd = Path.cwd()
        keywords = ['TODO', 'FIXME', 'HACK', 'BUG', 'NOTE']
        findings = []
        
        for file_path in cwd.rglob("*"):
            if (file_path.is_file() and 
                file_path.suffix in {'.py', '.js', '.ts', '.go', '.rs', '.java', '.cpp', '.c'} and
                not any(part.startswith('.') for part in file_path.parts)):
                
                try:
                    with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                        for line_num, line in enumerate(f, 1):
                            line_upper = line.upper()
                            for keyword in keywords:
                                if keyword in line_upper:
                                    findings.append(f"{file_path.name}:{line_num} [{keyword}] {line.strip()}")
                                    break
                except Exception:
                    continue
        
        if findings:
            result = f"[AnalyzeAgent] Found {len(findings)} TODO items:\n"
            result += "\n".join(findings[:20])  # Limit to first 20
            if len(findings) > 20:
                result += f"\n... and {len(findings) - 20} more"
        else:
            result = "[AnalyzeAgent] No TODO items found"
            
        return result
    
    def _security_scan(self) -> str:
        """Basic security scanning"""
        cwd = Path.cwd()
        issues = []
        
        # Check for common security issues
        dangerous_patterns = [
            ('password', 'Hardcoded password'),
            ('api_key', 'Hardcoded API key'), 
            ('secret', 'Hardcoded secret'),
            ('token', 'Hardcoded token'),
            ('eval(', 'Dangerous eval() usage'),
            ('exec(', 'Dangerous exec() usage'),
            ('os.system(', 'Dangerous system call'),
            ('shell=True', 'Shell injection risk')
        ]
        
        for file_path in cwd.rglob("*.py"):
            if not any(part.startswith('.') for part in file_path.parts):
                try:
                    with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                        content = f.read().lower()
                        for pattern, description in dangerous_patterns:
                            if pattern in content:
                                issues.append(f"{file_path.name}: {description}")
                except Exception:
                    continue
        
        if issues:
            result = f"[AnalyzeAgent] Security Issues Found ({len(issues)}):\n"
            result += "\n".join(issues[:10])
        else:
            result = "[AnalyzeAgent] No obvious security issues found"
            
        return result
    
    def _dependency_audit(self) -> str:
        """Audit dependencies for vulnerabilities"""
        cwd = Path.cwd()
        results = []
        
        # Python dependencies
        if (cwd / "requirements.txt").exists():
            try:
                result = subprocess.run(
                    ["pip", "audit"], 
                    capture_output=True, text=True, timeout=30
                )
                if result.returncode == 0:
                    results.append("✅ Python dependencies: No vulnerabilities found")
                else:
                    results.append(f"⚠️ Python dependencies: {result.stdout}")
            except Exception:
                results.append("⚠️ Python audit failed (pip audit not available)")
        
        # Node.js dependencies  
        if (cwd / "package.json").exists():
            try:
                result = subprocess.run(
                    ["npm", "audit"], 
                    capture_output=True, text=True, timeout=30
                )
                if "0 vulnerabilities" in result.stdout:
                    results.append("✅ Node.js dependencies: No vulnerabilities found")
                else:
                    results.append(f"⚠️ Node.js dependencies: Check required")
            except Exception:
                results.append("⚠️ Node.js audit failed")
        
        if not results:
            results.append("🔍 No dependency files found for audit")
            
        return f"[AnalyzeAgent] Dependency Audit:\n" + "\n".join(results)
    
    def _git_analysis(self) -> str:
        """Git repository analysis"""
        try:
            # Get basic git info
            branch_result = subprocess.run(
                ["git", "branch", "--show-current"], 
                capture_output=True, text=True, timeout=10
            )
            
            status_result = subprocess.run(
                ["git", "status", "--porcelain"], 
                capture_output=True, text=True, timeout=10
            )
            
            log_result = subprocess.run(
                ["git", "log", "--oneline", "-10"], 
                capture_output=True, text=True, timeout=10
            )
            
            results = []
            if branch_result.returncode == 0:
                results.append(f"🌿 Current branch: {branch_result.stdout.strip()}")
            
            if status_result.returncode == 0:
                if status_result.stdout.strip():
                    modified_count = len(status_result.stdout.strip().split('\n'))
                    results.append(f"📝 Modified files: {modified_count}")
                else:
                    results.append("✅ Working tree clean")
            
            if log_result.returncode == 0:
                results.append("📋 Recent commits:")
                for line in log_result.stdout.strip().split('\n')[:5]:
                    results.append(f"   {line}")
            
            return f"[AnalyzeAgent] Git Analysis:\n" + "\n".join(results)
            
        except Exception as e:
            return f"[AnalyzeAgent] Git analysis failed: {e}"
    
    def _general_analysis(self, query: str) -> str:
        """Handle general analysis queries"""
        return f"[AnalyzeAgent] Analysis for '{query}':\n🔍 Use specific commands like 'analyze codebase', 'count lines', 'find todos', 'security scan', or 'dependency check'"