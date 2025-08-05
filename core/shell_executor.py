# core/shell_executor.py
import subprocess
import os
import sys
from pathlib import Path
from typing import Optional, Dict, Any, List, Tuple
import threading
import time

from .shell_parser import Pipeline, Command

class ShellState:
    """Persistent shell state across commands"""
    def __init__(self):
        self.current_dir = Path.cwd()
        self.environment = os.environ.copy()
        self.aliases = {}
        self.history = []
        
    def change_directory(self, path: str) -> bool:
        """Change current directory with error handling"""
        try:
            new_path = Path(path).expanduser().resolve()
            if new_path.exists() and new_path.is_dir():
                self.current_dir = new_path
                os.chdir(new_path)
                return True
            else:
                return False
        except Exception:
            return False
            
    def set_env(self, name: str, value: str):
        """Set environment variable"""
        self.environment[name] = value
        os.environ[name] = value
        
    def get_env(self, name: str) -> Optional[str]:
        """Get environment variable"""
        return self.environment.get(name)
        
    def add_alias(self, name: str, command: str):
        """Add command alias"""
        self.aliases[name] = command
        
    def resolve_alias(self, command: str) -> str:
        """Resolve alias if exists"""
        return self.aliases.get(command, command)

class ShellExecutor:
    """Production shell executor with pipes, redirects, and persistent state"""
    
    def __init__(self):
        self.state = ShellState()
        self.builtin_commands = {
            'cd': self._builtin_cd,
            'pwd': self._builtin_pwd,
            'set': self._builtin_set,
            'export': self._builtin_export,
            'alias': self._builtin_alias,
            'env': self._builtin_env,
            'exit': self._builtin_exit,
            'echo': self._builtin_echo,
        }
        
    def execute_pipeline(self, pipeline: Pipeline) -> Tuple[int, str, str]:
        """Execute a complete pipeline with pipes and redirects"""
        if not pipeline.commands:
            return 0, "", ""
            
        try:
            if len(pipeline.commands) == 1:
                # Single command - handle redirects
                return self._execute_single_command(pipeline.commands[0])
            else:
                # Pipeline - connect with pipes
                return self._execute_pipeline(pipeline.commands)
                
        except Exception as e:
            return 1, "", f"Execution error: {e}"
    
    def _execute_single_command(self, command: Command) -> Tuple[int, str, str]:
        """Execute single command with redirects"""
        # Check if it's a builtin command
        if command.name in self.builtin_commands:
            return self.builtin_commands[command.name](command.args)
            
        # Resolve alias
        resolved_name = self.state.resolve_alias(command.name)
        
        try:
            # Prepare environment
            env = self.state.environment.copy()
            env.update(command.env_vars)
            
            # Handle input redirect
            stdin_input = None
            if command.input_redirect:
                try:
                    with open(command.input_redirect, 'r') as f:
                        stdin_input = f.read()
                except Exception as e:
                    return 1, "", f"Input redirect error: {e}"
            
            # Execute command
            process = subprocess.Popen(
                [resolved_name] + command.args,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                stdin=subprocess.PIPE if stdin_input else None,
                text=True,
                env=env,
                cwd=self.state.current_dir
            )
            
            stdout, stderr = process.communicate(input=stdin_input, timeout=30)
            
            # Handle output redirects
            if command.output_redirect:
                try:
                    with open(command.output_redirect, 'w') as f:
                        f.write(stdout)
                    stdout = ""  # Don't show output since it's redirected
                except Exception as e:
                    return 1, "", f"Output redirect error: {e}"
                    
            elif command.append_redirect:
                try:
                    with open(command.append_redirect, 'a') as f:
                        f.write(stdout)
                    stdout = ""  # Don't show output since it's redirected
                except Exception as e:
                    return 1, "", f"Append redirect error: {e}"
            
            return process.returncode, stdout, stderr
            
        except subprocess.TimeoutExpired:
            return 1, "", "Command timed out"
        except FileNotFoundError:
            return 1, "", f"Command not found: {command.name}"
        except Exception as e:
            return 1, "", f"Command execution error: {e}"
    
    def _execute_pipeline(self, commands: List[Command]) -> Tuple[int, str, str]:
        """Execute pipeline of commands connected by pipes"""
        processes = []
        
        try:
            for i, command in enumerate(commands):
                # Prepare environment
                env = self.state.environment.copy()
                env.update(command.env_vars)
                
                # Setup stdin
                if i == 0:
                    # First command
                    if command.input_redirect:
                        stdin = open(command.input_redirect, 'r')
                    else:
                        stdin = None
                else:
                    # Use output from previous command
                    stdin = processes[i-1].stdout
                
                # Setup stdout
                if i == len(commands) - 1:
                    # Last command
                    if command.output_redirect:
                        stdout = open(command.output_redirect, 'w')
                    elif command.append_redirect:
                        stdout = open(command.append_redirect, 'a')
                    else:
                        stdout = subprocess.PIPE
                else:
                    # Pipe to next command
                    stdout = subprocess.PIPE
                
                # Resolve alias
                resolved_name = self.state.resolve_alias(command.name)
                
                # Start process
                process = subprocess.Popen(
                    [resolved_name] + command.args,
                    stdin=stdin,
                    stdout=stdout,
                    stderr=subprocess.PIPE,
                    text=True,
                    env=env,
                    cwd=self.state.current_dir
                )
                
                processes.append(process)
                
                # Close stdin for previous process to allow pipeline flow
                if i > 0 and processes[i-1].stdout:
                    processes[i-1].stdout.close()
            
            # Wait for all processes and collect output
            final_stdout = ""
            final_stderr = ""
            final_returncode = 0
            
            for i, process in enumerate(processes):
                if i == len(processes) - 1:
                    # Last process - get its output
                    stdout, stderr = process.communicate(timeout=30)
                    final_stdout = stdout or ""
                    final_stderr = stderr or ""
                    final_returncode = process.returncode
                else:
                    # Wait for intermediate processes
                    process.wait(timeout=30)
                    if process.returncode != 0:
                        final_returncode = process.returncode
            
            return final_returncode, final_stdout, final_stderr
            
        except subprocess.TimeoutExpired:
            # Kill all processes on timeout
            for process in processes:
                try:
                    process.kill()
                except:
                    pass
            return 1, "", "Pipeline timed out"
            
        except Exception as e:
            # Kill all processes on error
            for process in processes:
                try:
                    process.kill()
                except:
                    pass
            return 1, "", f"Pipeline execution error: {e}"
    
    # Builtin commands
    def _builtin_cd(self, args: List[str]) -> Tuple[int, str, str]:
        """Change directory builtin"""
        if not args:
            target = Path.home()
        else:
            target = args[0]
            
        if self.state.change_directory(target):
            return 0, "", ""
        else:
            return 1, "", f"cd: {target}: No such file or directory"
    
    def _builtin_pwd(self, args: List[str]) -> Tuple[int, str, str]:
        """Print working directory"""
        return 0, str(self.state.current_dir), ""
    
    def _builtin_set(self, args: List[str]) -> Tuple[int, str, str]:
        """Set shell variable (not environment)"""
        if len(args) < 2:
            return 1, "", "Usage: set VAR value"
        return 0, "", ""
    
    def _builtin_export(self, args: List[str]) -> Tuple[int, str, str]:
        """Export environment variable"""
        if not args:
            # Show all env vars
            result = "\n".join(f"{k}={v}" for k, v in self.state.environment.items())
            return 0, result, ""
        
        for arg in args:
            if '=' in arg:
                name, value = arg.split('=', 1)
                self.state.set_env(name, value)
            else:
                # Export existing variable
                if arg in self.state.environment:
                    os.environ[arg] = self.state.environment[arg]
                    
        return 0, "", ""
    
    def _builtin_alias(self, args: List[str]) -> Tuple[int, str, str]:
        """Create command alias"""
        if not args:
            # Show all aliases
            result = "\n".join(f"{k}={v}" for k, v in self.state.aliases.items())
            return 0, result, ""
            
        for arg in args:
            if '=' in arg:
                name, command = arg.split('=', 1)
                self.state.add_alias(name, command)
                
        return 0, "", ""
    
    def _builtin_env(self, args: List[str]) -> Tuple[int, str, str]:
        """Show environment variables"""
        result = "\n".join(f"{k}={v}" for k, v in sorted(self.state.environment.items()))
        return 0, result, ""
    
    def _builtin_exit(self, args: List[str]) -> Tuple[int, str, str]:
        """Exit shell"""
        exit_code = 0
        if args:
            try:
                exit_code = int(args[0])
            except ValueError:
                pass
        sys.exit(exit_code)
    
    def _builtin_echo(self, args: List[str]) -> Tuple[int, str, str]:
        """Echo text"""
        return 0, " ".join(args), ""