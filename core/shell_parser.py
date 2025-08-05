# core/shell_parser.py
import re
import shlex
from typing import List, Dict, Any, Optional, Tuple
from dataclasses import dataclass
from enum import Enum

class TokenType(Enum):
    COMMAND = "command"
    ARGUMENT = "argument"
    PIPE = "pipe"
    REDIRECT_OUT = "redirect_out"
    REDIRECT_APPEND = "redirect_append"
    REDIRECT_IN = "redirect_in"
    VARIABLE = "variable"
    STRING = "string"
    BACKGROUND = "background"

@dataclass
class Token:
    type: TokenType
    value: str
    position: int

@dataclass
class Command:
    name: str
    args: List[str]
    env_vars: Dict[str, str]
    input_redirect: Optional[str] = None
    output_redirect: Optional[str] = None
    append_redirect: Optional[str] = None
    is_background: bool = False

@dataclass
class Pipeline:
    commands: List[Command]
    is_background: bool = False

class ShellParser:
    """Production-grade shell command parser supporting pipes, redirects, variables"""
    
    def __init__(self):
        self.variables = {}
        
    def tokenize(self, command_line: str) -> List[Token]:
        """Tokenize shell command line into structured tokens"""
        tokens = []
        position = 0
        
        # Handle quotes and escaping properly
        try:
            # Split respecting quotes but keep track of operators
            parts = []
            current = ""
            in_quotes = False
            quote_char = None
            i = 0
            
            while i < len(command_line):
                char = command_line[i]
                
                if char in ['"', "'"] and not in_quotes:
                    in_quotes = True
                    quote_char = char
                    current += char
                elif char == quote_char and in_quotes:
                    in_quotes = False
                    quote_char = None
                    current += char
                elif not in_quotes and char in ['|', '>', '<', '&']:
                    if current.strip():
                        parts.append(current.strip())
                        current = ""
                    
                    # Handle multi-character operators
                    if char == '>' and i + 1 < len(command_line) and command_line[i + 1] == '>':
                        parts.append('>>')
                        i += 1
                    elif char == '&' and i + 1 < len(command_line) and command_line[i + 1] == '&':
                        parts.append('&&')
                        i += 1
                    elif char == '|' and i + 1 < len(command_line) and command_line[i + 1] == '|':
                        parts.append('||')
                        i += 1
                    else:
                        parts.append(char)
                        
                elif not in_quotes and char.isspace():
                    if current.strip():
                        parts.append(current.strip())
                        current = ""
                else:
                    current += char
                    
                i += 1
                
            if current.strip():
                parts.append(current.strip())
                
        except Exception:
            # Fallback to simple split if complex parsing fails
            parts = command_line.split()
        
        # Convert parts to tokens
        for i, part in enumerate(parts):
            if part == '|':
                tokens.append(Token(TokenType.PIPE, part, position))
            elif part == '>':
                tokens.append(Token(TokenType.REDIRECT_OUT, part, position))
            elif part == '>>':
                tokens.append(Token(TokenType.REDIRECT_APPEND, part, position))
            elif part == '<':
                tokens.append(Token(TokenType.REDIRECT_IN, part, position))
            elif part == '&':
                tokens.append(Token(TokenType.BACKGROUND, part, position))
            elif part.startswith('$'):
                tokens.append(Token(TokenType.VARIABLE, part, position))
            elif i == 0 or (i > 0 and parts[i-1] in ['|', '&&', '||']):
                tokens.append(Token(TokenType.COMMAND, part, position))
            else:
                tokens.append(Token(TokenType.ARGUMENT, part, position))
                
            position += len(part) + 1
            
        return tokens
    
    def parse(self, command_line: str) -> Pipeline:
        """Parse command line into structured pipeline"""
        if not command_line.strip():
            return Pipeline(commands=[], is_background=False)
            
        tokens = self.tokenize(command_line)
        commands = []
        current_command = None
        
        i = 0
        while i < len(tokens):
            token = tokens[i]
            
            if token.type == TokenType.COMMAND:
                # Start new command
                if current_command:
                    commands.append(current_command)
                current_command = Command(
                    name=self._expand_variables(token.value),
                    args=[],
                    env_vars={}
                )
                
            elif token.type == TokenType.ARGUMENT and current_command:
                current_command.args.append(self._expand_variables(token.value))
                
            elif token.type == TokenType.PIPE:
                if current_command:
                    commands.append(current_command)
                    current_command = None
                    
            elif token.type == TokenType.REDIRECT_OUT and current_command:
                if i + 1 < len(tokens):
                    current_command.output_redirect = self._expand_variables(tokens[i + 1].value)
                    i += 1
                    
            elif token.type == TokenType.REDIRECT_APPEND and current_command:
                if i + 1 < len(tokens):
                    current_command.append_redirect = self._expand_variables(tokens[i + 1].value)
                    i += 1
                    
            elif token.type == TokenType.REDIRECT_IN and current_command:
                if i + 1 < len(tokens):
                    current_command.input_redirect = self._expand_variables(tokens[i + 1].value)
                    i += 1
                    
            elif token.type == TokenType.BACKGROUND:
                if current_command:
                    current_command.is_background = True
                    
            elif token.type == TokenType.VARIABLE:
                # Handle variable assignments like VAR=value
                if '=' in token.value:
                    var_name, var_value = token.value.split('=', 1)
                    var_name = var_name.lstrip('$')
                    if current_command:
                        current_command.env_vars[var_name] = self._expand_variables(var_value)
                    else:
                        self.variables[var_name] = self._expand_variables(var_value)
                        
            i += 1
            
        # Add final command
        if current_command:
            commands.append(current_command)
            
        return Pipeline(
            commands=commands,
            is_background=any(cmd.is_background for cmd in commands)
        )
    
    def _expand_variables(self, text: str) -> str:
        """Expand shell variables like $VAR and ${VAR}"""
        import os
        
        # Handle ${VAR} format
        pattern = r'\$\{([^}]+)\}'
        for match in re.finditer(pattern, text):
            var_name = match.group(1)
            var_value = self.variables.get(var_name, os.environ.get(var_name, ''))
            text = text.replace(match.group(0), var_value)
            
        # Handle $VAR format
        pattern = r'\$([A-Za-z_][A-Za-z0-9_]*)'
        for match in re.finditer(pattern, text):
            var_name = match.group(1)
            var_value = self.variables.get(var_name, os.environ.get(var_name, ''))
            text = text.replace(match.group(0), var_value)
            
        return text
    
    def set_variable(self, name: str, value: str):
        """Set shell variable"""
        self.variables[name] = value
        
    def get_variable(self, name: str) -> Optional[str]:
        """Get shell variable"""
        return self.variables.get(name) or os.environ.get(name)
    
    def get_all_variables(self) -> Dict[str, str]:
        """Get all variables (shell + environment)"""
        import os
        result = os.environ.copy()
        result.update(self.variables)
        return result