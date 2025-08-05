# agents/mcp_agent.py
"""
MCP (Modular Command Pipeline) Agent
Handles complex command pipelines with AI-enhanced processing
"""
import sys
import json
from pathlib import Path
from typing import List, Dict, Any, Optional
import subprocess

# Add parent directory for imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from agents.agent_base import Agent
from core.shell_parser import ShellParser, Pipeline
from core.shell_executor import ShellExecutor
from core.claude_memory import memory

class Mcp_agent(Agent):
    """Modular Command Pipeline Agent - AI-enhanced command processing"""
    
    def __init__(self):
        self.parser = ShellParser()
        self.executor = ShellExecutor()
        
    def run(self, input_text: str) -> str:
        """
        Process commands with AI enhancement
        
        Examples:
        - "find all python files and count lines"
        - "compress all images in current folder"
        - "analyze git history and summarize"
        - "ls | filter .py | count lines"
        """
        input_lower = input_text.lower().strip()
        
        # Check if it's a natural language request
        if self._is_natural_language(input_text):
            return self._process_natural_language(input_text)
        
        # Check if it's a pipeline with AI functions
        elif any(keyword in input_lower for keyword in ['filter', 'summarize', 'analyze', 'enhance']):
            return self._process_ai_pipeline(input_text)
        
        # Regular shell command
        else:
            return self._execute_shell_command(input_text)
    
    def _is_natural_language(self, text: str) -> bool:
        """Detect if input is natural language vs shell command"""
        natural_indicators = [
            'find all', 'show me', 'list all', 'count', 'analyze',
            'summarize', 'compress', 'convert', 'search for'
        ]
        return any(indicator in text.lower() for indicator in natural_indicators)
    
    def _process_natural_language(self, request: str) -> str:
        """Convert natural language to shell commands"""
        request_lower = request.lower()
        
        # Python file operations
        if 'python files' in request_lower and 'count lines' in request_lower:
            commands = [
                'find . -name "*.py" -type f',
                'wc -l'
            ]
            return self._execute_pipeline_commands(commands, "Finding Python files and counting lines")
        
        elif 'find all' in request_lower and 'python' in request_lower:
            command = 'find . -name "*.py" -type f'
            return self._execute_single_command(command, "Finding all Python files")
        
        # Image operations
        elif 'compress' in request_lower and 'image' in request_lower:
            return self._handle_image_compression()
        
        # Git operations  
        elif 'git history' in request_lower:
            commands = [
                'git log --oneline -10',
                'git log --stat --oneline -5'
            ]
            return self._execute_pipeline_commands(commands, "Analyzing git history")
        
        # File system analysis
        elif 'analyze' in request_lower and ('folder' in request_lower or 'directory' in request_lower):
            return self._analyze_directory()
        
        # Fallback - try to interpret
        else:
            return self._attempt_interpretation(request)
    
    def _process_ai_pipeline(self, pipeline_text: str) -> str:
        """Process pipeline with AI functions like filter, summarize, analyze"""
        parts = pipeline_text.split('|')
        
        if len(parts) < 2:
            return self._execute_shell_command(pipeline_text)
        
        # Execute first command
        first_command = parts[0].strip()
        
        try:
            result = subprocess.run(
                first_command, shell=True, capture_output=True, text=True, timeout=30
            )
            
            if result.returncode != 0:
                return f"[MCP] First command failed: {result.stderr}"
                
            output = result.stdout
            
            # Process through AI functions
            for i in range(1, len(parts)):
                function = parts[i].strip().lower()
                output = self._apply_ai_function(function, output)
                
            return f"[MCP] Pipeline Result:\n{output}"
            
        except Exception as e:
            return f"[MCP] Pipeline error: {e}"
    
    def _apply_ai_function(self, function: str, data: str) -> str:
        """Apply AI function to data"""
        if 'filter' in function:
            return self._ai_filter(function, data)
        elif 'summarize' in function:
            return self._ai_summarize(data)
        elif 'analyze' in function:
            return self._ai_analyze(data)
        elif 'count' in function:
            return self._ai_count(function, data)
        else:
            return data
    
    def _ai_filter(self, filter_spec: str, data: str) -> str:
        """AI-powered filtering"""
        lines = data.split('\n')
        
        if '.py' in filter_spec:
            filtered = [line for line in lines if '.py' in line]
        elif '.js' in filter_spec:
            filtered = [line for line in lines if '.js' in line]
        elif 'contains' in filter_spec:
            # Extract what to search for
            search_term = filter_spec.split('contains')[-1].strip(' "\'')
            filtered = [line for line in lines if search_term in line]
        else:
            # Default - remove empty lines
            filtered = [line for line in lines if line.strip()]
            
        return '\n'.join(filtered)
    
    def _ai_summarize(self, data: str) -> str:
        """AI-powered summarization"""
        lines = data.split('\n')
        non_empty_lines = [line for line in lines if line.strip()]
        
        summary = f"Summary: {len(non_empty_lines)} items found\n"
        
        # Group by file extension if they look like files
        extensions = {}
        for line in non_empty_lines:
            if '.' in line:
                parts = line.split('.')
                if len(parts) > 1:
                    ext = parts[-1].split()[0]  # Get extension, ignore extra info
                    extensions[ext] = extensions.get(ext, 0) + 1
        
        if extensions:
            summary += "File types:\n"
            for ext, count in sorted(extensions.items()):
                summary += f"  .{ext}: {count} files\n"
        
        # Show first few items as examples
        if non_empty_lines:
            summary += f"\nFirst few items:\n"
            for line in non_empty_lines[:5]:
                summary += f"  {line}\n"
                
        return summary
    
    def _ai_analyze(self, data: str) -> str:
        """AI-powered analysis"""
        lines = data.split('\n')
        non_empty_lines = [line for line in lines if line.strip()]
        
        analysis = f"Analysis Results:\n"
        analysis += f"- Total items: {len(non_empty_lines)}\n"
        
        # Character analysis
        total_chars = sum(len(line) for line in non_empty_lines)
        avg_length = total_chars / len(non_empty_lines) if non_empty_lines else 0
        analysis += f"- Average line length: {avg_length:.1f} characters\n"
        
        # Pattern detection
        patterns = {}
        for line in non_empty_lines:
            # Look for common patterns
            if line.startswith('/'):
                patterns['paths'] = patterns.get('paths', 0) + 1
            elif '@' in line:
                patterns['emails'] = patterns.get('emails', 0) + 1
            elif line.startswith('http'):
                patterns['urls'] = patterns.get('urls', 0) + 1
            elif any(char.isdigit() for char in line):
                patterns['with_numbers'] = patterns.get('with_numbers', 0) + 1
                
        if patterns:
            analysis += "- Patterns found:\n"
            for pattern, count in patterns.items():
                analysis += f"  {pattern}: {count}\n"
                
        return analysis
    
    def _ai_count(self, count_spec: str, data: str) -> str:
        """AI-powered counting"""
        if 'lines' in count_spec:
            lines = data.split('\n')
            non_empty = [line for line in lines if line.strip()]
            return f"Line count: {len(non_empty)} non-empty lines (total: {len(lines)})"
        
        elif 'words' in count_spec:
            words = data.split()
            return f"Word count: {len(words)} words"
        
        elif 'files' in count_spec:
            lines = data.split('\n')
            files = [line for line in lines if line.strip() and not line.endswith('/')]
            return f"File count: {len(files)} files"
        
        else:
            # Default - count non-empty lines
            lines = [line for line in data.split('\n') if line.strip()]
            return f"Count: {len(lines)} items"
    
    def _execute_pipeline_commands(self, commands: List[str], description: str) -> str:
        """Execute a series of commands in pipeline"""
        result = f"[MCP] {description}\n\n"
        
        for i, command in enumerate(commands):
            try:
                output = subprocess.run(
                    command, shell=True, capture_output=True, text=True, timeout=30
                )
                
                result += f"Command {i+1}: {command}\n"
                if output.stdout:
                    result += f"Output:\n{output.stdout}\n"
                if output.stderr:
                    result += f"Errors:\n{output.stderr}\n"
                result += "\n"
                
            except Exception as e:
                result += f"Command {i+1} failed: {e}\n"
                
        return result
    
    def _execute_single_command(self, command: str, description: str) -> str:
        """Execute single command with description"""
        try:
            result = subprocess.run(
                command, shell=True, capture_output=True, text=True, timeout=30
            )
            
            response = f"[MCP] {description}\n\n"
            if result.stdout:
                response += result.stdout
            if result.stderr:
                response += f"Errors: {result.stderr}"
                
            return response
            
        except Exception as e:
            return f"[MCP] Error executing '{command}': {e}"
    
    def _handle_image_compression(self) -> str:
        """Handle image compression requests"""
        # Check what image files exist
        try:
            result = subprocess.run(
                'find . -name "*.jpg" -o -name "*.png" -o -name "*.jpeg" | head -10',
                shell=True, capture_output=True, text=True
            )
            
            if result.stdout.strip():
                files = result.stdout.strip().split('\n')
                response = f"[MCP] Found {len(files)} image files:\n"
                for file in files:
                    response += f"  {file}\n"
                response += "\nTo compress, you could use:\n"
                response += "  - ImageMagick: convert image.jpg -quality 85 compressed.jpg\n"
                response += "  - Or install a batch compression tool\n"
                return response
            else:
                return "[MCP] No image files found in current directory"
                
        except Exception as e:
            return f"[MCP] Error finding images: {e}"
    
    def _analyze_directory(self) -> str:
        """Analyze current directory structure"""
        commands = [
            'find . -type f | head -20',
            'du -sh * 2>/dev/null | sort -rh | head -10',
            'find . -name "*.py" | wc -l',
            'find . -name "*.js" | wc -l',
            'find . -name "*.md" | wc -l'
        ]
        
        return self._execute_pipeline_commands(commands, "Directory Analysis")
    
    def _attempt_interpretation(self, request: str) -> str:
        """Attempt to interpret unclear natural language"""
        return f"[MCP] Could not interpret: '{request}'\n\nTry being more specific, like:\n- 'find all python files'\n- 'analyze git history'\n- 'compress images'\n- 'ls | filter .py | count lines'"
    
    def _execute_shell_command(self, command: str) -> str:
        """Execute regular shell command"""
        try:
            pipeline = self.parser.parse(command)
            returncode, stdout, stderr = self.executor.execute_pipeline(pipeline)
            
            result = f"[MCP] Executed: {command}\n"
            if stdout:
                result += f"Output:\n{stdout}\n"
            if stderr:
                result += f"Errors:\n{stderr}\n"
            result += f"Exit code: {returncode}"
            
            return result
            
        except Exception as e:
            return f"[MCP] Shell execution error: {e}"