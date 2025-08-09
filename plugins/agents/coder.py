"""
CoderAgent - AI-powered code generation and manipulation
"""

import re
import os
import ast
import json
from pathlib import Path
from typing import Dict, List, Optional, Any

class CoderAgent:
    """AI-powered coding agent for generation, refactoring, and explanation"""
    
    def __init__(self):
        self.language_extensions = {
            'python': ['.py'],
            'javascript': ['.js'],
            'typescript': ['.ts'],
            'react': ['.jsx', '.tsx'],
            'go': ['.go'],
            'java': ['.java'],
            'cpp': ['.cpp', '.cc', '.cxx'],
            'c': ['.c'],
            'csharp': ['.cs'],
            'rust': ['.rs'],
            'php': ['.php'],
            'ruby': ['.rb'],
            'sql': ['.sql'],
            'html': ['.html', '.htm'],
            'css': ['.css'],
            'json': ['.json'],
            'yaml': ['.yml', '.yaml'],
            'markdown': ['.md'],
        }
        
        self.code_templates = {
            'python': self._get_python_templates(),
            'javascript': self._get_javascript_templates(),
            'go': self._get_go_templates(),
            'react': self._get_react_templates()
        }
    
    def generate(self, prompt: str, language: str = None, output_file: str = None) -> str:
        """Generate code based on natural language prompt"""
        
        # Parse prompt to understand what to generate
        request = self._parse_code_request(prompt)
        
        # Determine language if not specified
        if not language:
            language = self._detect_language_from_prompt(prompt)
        
        # Generate code based on request type
        if request['type'] == 'function':
            code = self._generate_function(request, language)
        elif request['type'] == 'class':
            code = self._generate_class(request, language)
        elif request['type'] == 'api':
            code = self._generate_api(request, language)
        elif request['type'] == 'script':
            code = self._generate_script(request, language)
        elif request['type'] == 'component':
            code = self._generate_component(request, language)
        else:
            code = self._generate_generic_code(request, language)
        
        # Save to file if specified
        if output_file:
            self._save_code_to_file(code, output_file)
        
        return code
    
    def refactor(self, file_path: str, refactor_type: str = "improve") -> str:
        """Refactor existing code"""
        
        if not os.path.exists(file_path):
            return f"❌ File not found: {file_path}"
        
        with open(file_path, 'r', encoding='utf-8') as f:
            original_code = f.read()
        
        language = self._detect_language_from_file(file_path)
        
        if refactor_type == "improve":
            refactored = self._improve_code(original_code, language)
        elif refactor_type == "optimize":
            refactored = self._optimize_code(original_code, language)
        elif refactor_type == "modernize":
            refactored = self._modernize_code(original_code, language)
        else:
            refactored = self._improve_code(original_code, language)
        
        return refactored
    
    def explain(self, code_or_file: str) -> str:
        """Explain what code does"""
        
        # Check if it's a file path or code snippet
        if os.path.exists(code_or_file):
            with open(code_or_file, 'r', encoding='utf-8') as f:
                code = f.read()
            file_name = os.path.basename(code_or_file)
            language = self._detect_language_from_file(code_or_file)
        else:
            code = code_or_file
            file_name = "code snippet"
            language = self._detect_language_from_code(code)
        
        return self._generate_explanation(code, language, file_name)
    
    def _parse_code_request(self, prompt: str) -> Dict[str, Any]:
        """Parse the code generation request"""
        prompt_lower = prompt.lower()
        
        request = {
            'type': 'function',  # default
            'name': '',
            'description': prompt,
            'features': [],
            'parameters': [],
            'return_type': None
        }
        
        # Determine request type
        if any(word in prompt_lower for word in ['class', 'object', 'struct']):
            request['type'] = 'class'
        elif any(word in prompt_lower for word in ['api', 'server', 'endpoint', 'rest']):
            request['type'] = 'api'
        elif any(word in prompt_lower for word in ['script', 'program', 'app']):
            request['type'] = 'script'
        elif any(word in prompt_lower for word in ['component', 'ui', 'interface']):
            request['type'] = 'component'
        
        # Extract function/class name
        name_patterns = [
            r'(?:function|def|class)\s+(\w+)',
            r'create\s+(?:a\s+)?(\w+)',
            r'make\s+(?:a\s+)?(\w+)',
            r'build\s+(?:a\s+)?(\w+)'
        ]
        
        for pattern in name_patterns:
            match = re.search(pattern, prompt_lower)
            if match:
                request['name'] = match.group(1)
                break
        
        # Extract features
        if 'validation' in prompt_lower:
            request['features'].append('validation')
        if 'error handling' in prompt_lower:
            request['features'].append('error_handling')
        if 'logging' in prompt_lower:
            request['features'].append('logging')
        if 'async' in prompt_lower or 'asynchronous' in prompt_lower:
            request['features'].append('async')
        
        return request
    
    def _detect_language_from_prompt(self, prompt: str) -> str:
        """Detect programming language from prompt"""
        prompt_lower = prompt.lower()
        
        language_keywords = {
            'python': ['python', 'django', 'flask', 'pandas', 'numpy'],
            'javascript': ['javascript', 'js', 'node', 'express', 'npm'],
            'typescript': ['typescript', 'ts'],
            'react': ['react', 'jsx', 'component'],
            'go': ['go', 'golang'],
            'java': ['java', 'spring', 'maven'],
            'cpp': ['c++', 'cpp'],
            'rust': ['rust', 'cargo'],
            'php': ['php', 'laravel'],
            'ruby': ['ruby', 'rails']
        }
        
        for lang, keywords in language_keywords.items():
            if any(keyword in prompt_lower for keyword in keywords):
                return lang
        
        return 'python'  # default
    
    def _detect_language_from_file(self, file_path: str) -> str:
        """Detect language from file extension"""
        ext = Path(file_path).suffix.lower()
        
        for lang, extensions in self.language_extensions.items():
            if ext in extensions:
                return lang
        
        return 'text'
    
    def _detect_language_from_code(self, code: str) -> str:
        """Detect language from code content"""
        if 'def ' in code or 'import ' in code:
            return 'python'
        elif 'function' in code or 'const' in code or 'let' in code:
            return 'javascript'
        elif 'func' in code and 'package' in code:
            return 'go'
        elif 'public class' in code or 'import java' in code:
            return 'java'
        else:
            return 'text'
    
    def _generate_function(self, request: Dict, language: str) -> str:
        """Generate a function based on request"""
        
        if language == 'python':
            return self._generate_python_function(request)
        elif language == 'javascript':
            return self._generate_javascript_function(request)
        elif language == 'go':
            return self._generate_go_function(request)
        else:
            return self._generate_python_function(request)  # fallback
    
    def _generate_python_function(self, request: Dict) -> str:
        """Generate Python function"""
        name = request.get('name', 'generated_function')
        description = request.get('description', 'Generated function')
        features = request.get('features', [])
        
        # Build function signature
        params = "data"  # default parameter
        if 'validation' in features:
            params = "data: Any"
        
        # Build function body
        body_lines = [
            f'    """',
            f'    {description}',
            f'    """'
        ]
        
        if 'validation' in features:
            body_lines.extend([
                '',
                '    if not data:',
                '        raise ValueError("Data cannot be empty")'
            ])
        
        if 'logging' in features:
            body_lines.extend([
                '',
                '    import logging',
                '    logging.info(f"Processing {name}")'
            ])
        
        if 'error_handling' in features:
            body_lines.extend([
                '',
                '    try:',
                '        # Process the data',
                '        result = process_data(data)',
                '        return result',
                '    except Exception as e:',
                '        logging.error(f"Error in {name}: {e}")',
                '        raise'
            ])
        else:
            body_lines.extend([
                '',
                '    # TODO: Implement function logic',
                '    result = None',
                '',
                '    return result'
            ])
        
        # Combine everything
        imports = []
        if 'validation' in features:
            imports.append('from typing import Any')
        if 'logging' in features:
            imports.append('import logging')
        
        function_code = []
        if imports:
            function_code.extend(imports)
            function_code.append('')
        
        if 'async' in features:
            function_code.append(f'async def {name}({params}):')
        else:
            function_code.append(f'def {name}({params}):')
        
        function_code.extend(body_lines)
        
        return '\\n'.join(function_code)
    
    def _generate_javascript_function(self, request: Dict) -> str:
        """Generate JavaScript function"""
        name = request.get('name', 'generatedFunction')
        description = request.get('description', 'Generated function')
        features = request.get('features', [])
        
        # Build function
        if 'async' in features:
            func_def = f'async function {name}(data) {{'
        else:
            func_def = f'function {name}(data) {{'
        
        body_lines = [
            f'  /**',
            f'   * {description}',
            f'   */',
        ]
        
        if 'validation' in features:
            body_lines.extend([
                '',
                '  if (!data) {',
                '    throw new Error("Data cannot be empty");',
                '  }'
            ])
        
        if 'logging' in features:
            body_lines.extend([
                '',
                f'  console.log("Processing {name}");'
            ])
        
        body_lines.extend([
            '',
            '  // TODO: Implement function logic',
            '  const result = null;',
            '',
            '  return result;'
        ])
        
        function_code = [func_def] + body_lines + ['}']
        
        return '\\n'.join(function_code)
    
    def _generate_go_function(self, request: Dict) -> str:
        """Generate Go function"""
        name = request.get('name', 'GeneratedFunction')
        description = request.get('description', 'Generated function')
        
        # Capitalize function name for Go conventions
        if name and name[0].islower():
            name = name[0].upper() + name[1:]
        
        function_code = [
            '// ' + description,
            f'func {name}(data interface{{}}) (interface{{}}, error) {{',
            '    // TODO: Implement function logic',
            '    return nil, nil',
            '}'
        ]
        
        return '\\n'.join(function_code)
    
    def _generate_class(self, request: Dict, language: str) -> str:
        """Generate a class"""
        if language == 'python':
            return self._generate_python_class(request)
        elif language == 'javascript':
            return self._generate_javascript_class(request)
        else:
            return self._generate_python_class(request)
    
    def _generate_python_class(self, request: Dict) -> str:
        """Generate Python class"""
        name = request.get('name', 'GeneratedClass')
        description = request.get('description', 'Generated class')
        
        class_code = [
            f'class {name}:',
            f'    """',
            f'    {description}',
            f'    """',
            '',
            '    def __init__(self):',
            '        """Initialize the class"""',
            '        pass',
            '',
            '    def process(self, data):',
            '        """Process data"""',
            '        # TODO: Implement processing logic',
            '        return data'
        ]
        
        return '\\n'.join(class_code)
    
    def _generate_javascript_class(self, request: Dict) -> str:
        """Generate JavaScript class"""
        name = request.get('name', 'GeneratedClass')
        description = request.get('description', 'Generated class')
        
        class_code = [
            f'/**',
            f' * {description}',
            f' */',
            f'class {name} {{',
            '',
            '  constructor() {',
            '    // Initialize the class',
            '  }',
            '',
            '  process(data) {',
            '    // TODO: Implement processing logic',
            '    return data;',
            '  }',
            '}'
        ]
        
        return '\\n'.join(class_code)
    
    def _generate_generic_code(self, request: Dict, language: str) -> str:
        """Generate generic code when specific type is unknown"""
        return self._generate_function(request, language)
    
    def _generate_api(self, request: Dict, language: str) -> str:
        """Generate API code"""
        # Implementation would depend on framework
        return "# API generation not implemented yet"
    
    def _generate_script(self, request: Dict, language: str) -> str:
        """Generate script code"""
        return "# Script generation not implemented yet"
    
    def _generate_component(self, request: Dict, language: str) -> str:
        """Generate UI component code"""
        return "// Component generation not implemented yet"
    
    def _improve_code(self, code: str, language: str) -> str:
        """Improve existing code"""
        # This is a simplified version - real implementation would be much more complex
        lines = code.splitlines()
        improved_lines = []
        
        for line in lines:
            # Remove trailing whitespace
            line = line.rstrip()
            
            # Basic improvements for Python
            if language == 'python':
                # Replace print with logging where appropriate
                if 'print(' in line and 'debug' in line.lower():
                    line = line.replace('print(', 'logging.debug(')
            
            improved_lines.append(line)
        
        return '\\n'.join(improved_lines)
    
    def _optimize_code(self, code: str, language: str) -> str:
        """Optimize code for performance"""
        return self._improve_code(code, language)  # Placeholder
    
    def _modernize_code(self, code: str, language: str) -> str:
        """Modernize code to use current best practices"""
        return self._improve_code(code, language)  # Placeholder
    
    def _generate_explanation(self, code: str, language: str, file_name: str) -> str:
        """Generate explanation of what the code does"""
        lines = code.splitlines()
        
        explanation = [
            f"📄 **Code Explanation: {file_name}**",
            "=" * 50,
            f"**Language:** {language.capitalize()}",
            f"**Lines of code:** {len(lines)}",
            "",
            "## 🎯 **Purpose**"
        ]
        
        # Analyze code structure
        if language == 'python':
            explanation.extend(self._explain_python_code(code))
        elif language == 'javascript':
            explanation.extend(self._explain_javascript_code(code))
        else:
            explanation.extend(self._explain_generic_code(code))
        
        return '\\n'.join(explanation)
    
    def _explain_python_code(self, code: str) -> List[str]:
        """Explain Python code"""
        explanation = []
        
        # Look for imports
        imports = re.findall(r'^import\s+(\w+)', code, re.MULTILINE)
        from_imports = re.findall(r'^from\s+(\w+)\s+import', code, re.MULTILINE)
        
        if imports or from_imports:
            explanation.append("**Dependencies:**")
            for imp in imports:
                explanation.append(f"- `{imp}` - External library")
            for imp in from_imports:
                explanation.append(f"- `{imp}` - Specific imports")
            explanation.append("")
        
        # Look for functions and classes
        functions = re.findall(r'^def\s+(\w+)', code, re.MULTILINE)
        classes = re.findall(r'^class\s+(\w+)', code, re.MULTILINE)
        
        if classes:
            explanation.append("**Classes:**")
            for cls in classes:
                explanation.append(f"- `{cls}` - Class definition")
            explanation.append("")
        
        if functions:
            explanation.append("**Functions:**")
            for func in functions:
                explanation.append(f"- `{func}()` - Function definition")
            explanation.append("")
        
        explanation.append("**Overall Structure:** Python script/module with standard structure")
        
        return explanation
    
    def _explain_javascript_code(self, code: str) -> List[str]:
        """Explain JavaScript code"""
        explanation = []
        
        # Look for functions
        functions = re.findall(r'function\s+(\w+)', code)
        arrow_functions = re.findall(r'const\s+(\w+)\s*=\s*\(', code)
        
        if functions or arrow_functions:
            explanation.append("**Functions:**")
            for func in functions:
                explanation.append(f"- `{func}()` - Function declaration")
            for func in arrow_functions:
                explanation.append(f"- `{func}()` - Arrow function")
            explanation.append("")
        
        explanation.append("**Overall Structure:** JavaScript code with modern ES6+ features")
        
        return explanation
    
    def _explain_generic_code(self, code: str) -> List[str]:
        """Explain generic code"""
        return [
            "This code file contains programming logic.",
            "Use more specific language detection for detailed analysis."
        ]
    
    def _save_code_to_file(self, code: str, file_path: str):
        """Save generated code to file"""
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(code)
    
    def _get_python_templates(self) -> Dict:
        """Get Python code templates"""
        return {}  # Could be expanded
    
    def _get_javascript_templates(self) -> Dict:
        """Get JavaScript code templates"""
        return {}  # Could be expanded
    
    def _get_go_templates(self) -> Dict:
        """Get Go code templates"""
        return {}  # Could be expanded
    
    def _get_react_templates(self) -> Dict:
        """Get React code templates"""
        return {}  # Could be expanded