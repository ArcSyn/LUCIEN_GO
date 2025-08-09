"""
ReviewAgent - AI-powered code review and analysis
"""

import re
import os
import ast
import json
from pathlib import Path
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass

@dataclass
class Issue:
    type: str  # 'error', 'warning', 'suggestion', 'style'
    line: Optional[int]
    column: Optional[int]
    message: str
    severity: str  # 'high', 'medium', 'low'
    category: str  # 'security', 'performance', 'maintainability', 'style', 'bug'
    suggestion: Optional[str] = None

class ReviewAgent:
    """AI-powered code reviewer that analyzes files for issues and improvements"""
    
    def __init__(self):
        self.supported_extensions = {
            '.py': self._review_python,
            '.js': self._review_javascript,
            '.ts': self._review_typescript,
            '.jsx': self._review_react,
            '.tsx': self._review_react,
            '.go': self._review_go,
            '.java': self._review_java,
            '.cpp': self._review_cpp,
            '.c': self._review_cpp,
            '.cs': self._review_csharp,
            '.php': self._review_php,
            '.rb': self._review_ruby,
            '.rs': self._review_rust,
            '.sql': self._review_sql,
            '.yml': self._review_yaml,
            '.yaml': self._review_yaml,
            '.json': self._review_json,
            '.md': self._review_markdown,
            '.html': self._review_html,
            '.css': self._review_css,
        }
        
    def analyze(self, file_path: str) -> str:
        """Analyze a file and return comprehensive review results"""
        
        if not os.path.exists(file_path):
            return f"❌ File not found: {file_path}"
        
        file_path_obj = Path(file_path)
        extension = file_path_obj.suffix.lower()
        
        if extension not in self.supported_extensions:
            return f"⚠️ Unsupported file type: {extension}"
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except UnicodeDecodeError:
            try:
                with open(file_path, 'r', encoding='latin-1') as f:
                    content = f.read()
            except Exception as e:
                return f"❌ Error reading file: {e}"
        
        # Get language-specific review
        review_func = self.supported_extensions[extension]
        issues = review_func(content, file_path)
        
        # Add general file analysis
        issues.extend(self._analyze_general_issues(content, file_path))
        
        # Format and return results
        return self._format_review_results(issues, file_path, content)
    
    def _review_python(self, content: str, file_path: str) -> List[Issue]:
        """Review Python code"""
        issues = []
        lines = content.splitlines()
        
        # Parse AST for advanced analysis
        try:
            tree = ast.parse(content)
            issues.extend(self._analyze_python_ast(tree))
        except SyntaxError as e:
            issues.append(Issue(
                type='error',
                line=e.lineno,
                column=e.offset,
                message=f"Syntax error: {e.msg}",
                severity='high',
                category='bug'
            ))
        
        # Line-by-line analysis
        for i, line in enumerate(lines, 1):
            # Security issues
            if re.search(r'eval\(|exec\(', line):
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Avoid using eval() or exec() - security risk",
                    severity='high',
                    category='security',
                    suggestion="Use safer alternatives like literal_eval() or specific parsing"
                ))
            
            # Performance issues
            if re.search(r'\.append\(.*\)\s*$', line) and 'for' in lines[max(0, i-2):i]:
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Consider using list comprehension instead of append in loop",
                    severity='medium',
                    category='performance'
                ))
            
            # Style issues
            if len(line) > 88:  # PEP 8 recommends 79, but 88 is common with black
                issues.append(Issue(
                    type='style',
                    line=i,
                    column=None,
                    message="Line too long (consider breaking it up)",
                    severity='low',
                    category='style'
                ))
            
            # Common anti-patterns
            if re.search(r'except:\s*$', line):
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Bare except clause - specify exception types",
                    severity='medium',
                    category='maintainability'
                ))
        
        return issues
    
    def _analyze_python_ast(self, tree: ast.AST) -> List[Issue]:
        """Analyze Python AST for advanced issues"""
        issues = []
        
        for node in ast.walk(tree):
            # Check for unused variables
            if isinstance(node, ast.Name) and isinstance(node.ctx, ast.Store):
                # This would need more sophisticated analysis
                pass
            
            # Check for missing docstrings
            if isinstance(node, (ast.FunctionDef, ast.ClassDef)):
                if not ast.get_docstring(node):
                    issues.append(Issue(
                        type='suggestion',
                        line=node.lineno,
                        column=node.col_offset,
                        message=f"{type(node).__name__} '{node.name}' missing docstring",
                        severity='low',
                        category='maintainability'
                    ))
            
            # Check for complex functions
            if isinstance(node, ast.FunctionDef):
                complexity = self._calculate_cyclomatic_complexity(node)
                if complexity > 10:
                    issues.append(Issue(
                        type='warning',
                        line=node.lineno,
                        column=node.col_offset,
                        message=f"Function '{node.name}' has high complexity ({complexity})",
                        severity='medium',
                        category='maintainability',
                        suggestion="Consider breaking this function into smaller functions"
                    ))
        
        return issues
    
    def _calculate_cyclomatic_complexity(self, node: ast.FunctionDef) -> int:
        """Calculate cyclomatic complexity of a function"""
        complexity = 1  # Base complexity
        
        for child in ast.walk(node):
            if isinstance(child, (ast.If, ast.While, ast.For, ast.ExceptHandler)):
                complexity += 1
            elif isinstance(child, ast.BoolOp):
                complexity += len(child.values) - 1
        
        return complexity
    
    def _review_javascript(self, content: str, file_path: str) -> List[Issue]:
        """Review JavaScript code"""
        issues = []
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            # Security issues
            if re.search(r'eval\(|innerHTML\s*=', line):
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Potential security risk - avoid eval() or innerHTML",
                    severity='high',
                    category='security'
                ))
            
            # Modern JS suggestions
            if re.search(r'var\s+\w+', line):
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Consider using 'let' or 'const' instead of 'var'",
                    severity='low',
                    category='style'
                ))
            
            # Performance issues
            if re.search(r'document\.getElementById|document\.querySelector', line):
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Consider caching DOM queries for better performance",
                    severity='low',
                    category='performance'
                ))
        
        return issues
    
    def _review_typescript(self, content: str, file_path: str) -> List[Issue]:
        """Review TypeScript code"""
        issues = self._review_javascript(content, file_path)  # Inherit JS rules
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            # Type safety
            if re.search(r':\s*any\b', line):
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Avoid 'any' type - use specific types for better type safety",
                    severity='medium',
                    category='maintainability'
                ))
            
            # Missing type annotations
            if re.search(r'function\s+\w+\s*\([^)]*\)\s*{', line):
                if not re.search(r':\s*\w+\s*=>', line):
                    issues.append(Issue(
                        type='suggestion',
                        line=i,
                        column=None,
                        message="Consider adding return type annotation",
                        severity='low',
                        category='maintainability'
                    ))
        
        return issues
    
    def _review_react(self, content: str, file_path: str) -> List[Issue]:
        """Review React/JSX code"""
        issues = self._review_typescript(content, file_path)
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            # React-specific issues
            if 'useState' in line and not re.search(r'const\s*\[.*,\s*set\w+\]\s*=', line):
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Follow destructuring pattern for useState",
                    severity='low',
                    category='style'
                ))
            
            # Missing key prop in lists
            if re.search(r'\.map\s*\(.*=>', line) and 'key=' not in line:
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Missing 'key' prop in list items",
                    severity='medium',
                    category='bug'
                ))
            
            # Accessibility
            if re.search(r'<img\s+', line) and 'alt=' not in line:
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Missing 'alt' attribute for accessibility",
                    severity='medium',
                    category='maintainability'
                ))
        
        return issues
    
    def _review_go(self, content: str, file_path: str) -> List[Issue]:
        """Review Go code"""
        issues = []
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            # Error handling
            if 'err != nil' in line and i < len(lines) and 'return' not in lines[i]:
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Error not handled or returned",
                    severity='high',
                    category='bug'
                ))
            
            # Goroutine without proper synchronization
            if re.search(r'go\s+\w+\(', line):
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Ensure proper synchronization for goroutines",
                    severity='medium',
                    category='maintainability'
                ))
        
        return issues
    
    def _review_java(self, content: str, file_path: str) -> List[Issue]:
        """Review Java code"""
        issues = []
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            # Exception handling
            if 'catch' in line and 'printStackTrace' in lines[min(i, len(lines)-1)]:
                issues.append(Issue(
                    type='warning',
                    line=i+1,
                    column=None,
                    message="Avoid printStackTrace() in production code",
                    severity='medium',
                    category='maintainability'
                ))
        
        return issues
    
    def _review_sql(self, content: str, file_path: str) -> List[Issue]:
        """Review SQL code"""
        issues = []
        lines = content.splitlines()
        
        for i, line in enumerate(lines, 1):
            line_upper = line.upper()
            
            # SQL injection risks
            if re.search(r'["\'][^"\']*\+.*["\']', line):
                issues.append(Issue(
                    type='warning',
                    line=i,
                    column=None,
                    message="Potential SQL injection - use parameterized queries",
                    severity='high',
                    category='security'
                ))
            
            # Performance issues
            if 'SELECT *' in line_upper:
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="Avoid SELECT * - specify needed columns",
                    severity='medium',
                    category='performance'
                ))
        
        return issues
    
    def _review_generic(self, content: str, file_path: str) -> List[Issue]:
        """Generic review for unsupported file types"""
        return []
    
    # Placeholder methods for other languages
    def _review_cpp(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_csharp(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_php(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_ruby(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_rust(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_yaml(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_json(self, content: str, file_path: str) -> List[Issue]:
        issues = []
        try:
            json.loads(content)
        except json.JSONDecodeError as e:
            issues.append(Issue(
                type='error',
                line=e.lineno,
                column=e.colno,
                message=f"JSON syntax error: {e.msg}",
                severity='high',
                category='bug'
            ))
        return issues
    
    def _review_markdown(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_html(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _review_css(self, content: str, file_path: str) -> List[Issue]:
        return self._review_generic(content, file_path)
    
    def _analyze_general_issues(self, content: str, file_path: str) -> List[Issue]:
        """Analyze general file issues"""
        issues = []
        lines = content.splitlines()
        
        # File size check
        if len(content) > 50000:  # 50KB
            issues.append(Issue(
                type='warning',
                line=None,
                column=None,
                message="Large file size - consider breaking into smaller files",
                severity='low',
                category='maintainability'
            ))
        
        # Long lines in any language
        for i, line in enumerate(lines, 1):
            if len(line) > 120:
                issues.append(Issue(
                    type='style',
                    line=i,
                    column=None,
                    message="Line exceeds 120 characters",
                    severity='low',
                    category='style'
                ))
        
        # TODO/FIXME comments
        for i, line in enumerate(lines, 1):
            if re.search(r'(TODO|FIXME|HACK|XXX)', line, re.IGNORECASE):
                issues.append(Issue(
                    type='suggestion',
                    line=i,
                    column=None,
                    message="TODO/FIXME comment found - consider addressing",
                    severity='low',
                    category='maintainability'
                ))
        
        return issues
    
    def _format_review_results(self, issues: List[Issue], file_path: str, content: str) -> str:
        """Format the review results into a readable report"""
        if not issues:
            return f"✅ **{os.path.basename(file_path)}** - No issues found! Code looks good."
        
        # Group issues by severity
        errors = [i for i in issues if i.type == 'error']
        warnings = [i for i in issues if i.type == 'warning']
        suggestions = [i for i in issues if i.type == 'suggestion']
        style_issues = [i for i in issues if i.type == 'style']
        
        lines = content.splitlines()
        result = [f"📋 **Code Review: {os.path.basename(file_path)}**"]
        result.append("=" * 50)
        
        # Summary
        total_issues = len(issues)
        result.append(f"**Summary:** {total_issues} issue{'s' if total_issues != 1 else ''} found")
        
        if errors:
            result.append(f"🚫 {len(errors)} error{'s' if len(errors) != 1 else ''}")
        if warnings:
            result.append(f"⚠️  {len(warnings)} warning{'s' if len(warnings) != 1 else ''}")
        if suggestions:
            result.append(f"💡 {len(suggestions)} suggestion{'s' if len(suggestions) != 1 else ''}")
        if style_issues:
            result.append(f"🎨 {len(style_issues)} style issue{'s' if len(style_issues) != 1 else ''}")
        
        result.append("")
        
        # Detailed issues
        for category, icon, items in [
            ("🚫 ERRORS", "🚫", errors),
            ("⚠️ WARNINGS", "⚠️", warnings),  
            ("💡 SUGGESTIONS", "💡", suggestions),
            ("🎨 STYLE ISSUES", "🎨", style_issues)
        ]:
            if items:
                result.append(f"## {category}")
                result.append("")
                for issue in sorted(items, key=lambda x: x.line or 0):
                    location = f"Line {issue.line}" if issue.line else "File level"
                    result.append(f"{icon} **{location}:** {issue.message}")
                    
                    if issue.line and issue.line <= len(lines):
                        # Show the problematic line
                        line_content = lines[issue.line - 1].strip()
                        if line_content:
                            result.append(f"   `{line_content}`")
                    
                    if issue.suggestion:
                        result.append(f"   💡 *Suggestion: {issue.suggestion}*")
                    
                    result.append("")
        
        # Overall assessment
        result.append("## 🎯 **Overall Assessment**")
        
        if errors:
            result.append("❌ **Not ready for production** - Fix errors first")
        elif warnings:
            result.append("🟡 **Needs attention** - Address warnings before deployment")  
        elif suggestions or style_issues:
            result.append("🟢 **Good code quality** - Minor improvements suggested")
        else:
            result.append("✨ **Excellent** - No issues found!")
        
        return "\\n".join(result)