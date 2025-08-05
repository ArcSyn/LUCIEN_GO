# cli/commands/shell.py
import typer
from pathlib import Path
import sys
import os

# Add parent directory to path to import core modules
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from core.shell_parser import ShellParser
from core.shell_executor import ShellExecutor
from core.config import config
from core.claude_memory import memory
from lucien.ui import console, spell_complete, spell_failed

app = typer.Typer()

# Global shell instances
parser = ShellParser()
executor = ShellExecutor()

@app.command()
def run(command: str):
    """Execute shell command with full parsing, pipes, and redirects"""
    try:
        # Parse the command
        pipeline = parser.parse(command)
        
        if not pipeline.commands:
            return
            
        # Execute the pipeline
        returncode, stdout, stderr = executor.execute_pipeline(pipeline)
        
        # Display output
        if stdout:
            console.print(stdout)
        if stderr:
            console.print(f"[red]{stderr}[/red]")
            
        # Log to memory if it's an important command
        if any(cmd.name in ['git', 'npm', 'pip', 'cargo'] for cmd in pipeline.commands):
            memory.log_interaction(
                f"shell: {command}",
                f"Executed with return code {returncode}",
                stdout + stderr if stdout or stderr else None
            )
            
        return returncode
        
    except Exception as e:
        error_msg = f"Shell execution error: {e}"
        console.print(f"[red]{error_msg}[/red]")
        return 1

@app.command()
def cd(path: str = typer.Argument(default="~")):
    """Change directory with persistent state"""
    if executor.state.change_directory(path):
        spell_complete(f"Changed to: {executor.state.current_dir}")
    else:
        spell_failed(f"Cannot change to directory: {path}")

@app.command()
def pwd():
    """Print current working directory"""
    console.print(str(executor.state.current_dir))

@app.command()
def env(variable: str = None, value: str = None):
    """Manage environment variables"""
    if variable and value:
        # Set environment variable
        executor.state.set_env(variable, value)
        spell_complete(f"Set {variable}={value}")
    elif variable:
        # Get specific variable
        val = executor.state.get_env(variable)
        if val:
            console.print(f"{variable}={val}")
        else:
            spell_failed(f"Variable {variable} not found")
    else:
        # Show all environment variables
        for key, value in sorted(executor.state.environment.items()):
            console.print(f"{key}={value}")

@app.command()
def alias(name: str = None, command: str = None):
    """Manage command aliases"""
    if name and command:
        # Set alias
        executor.state.add_alias(name, command)
        config.add_alias(name, command)
        spell_complete(f"Created alias: {name}={command}")
    elif name:
        # Show specific alias
        cmd = executor.state.aliases.get(name)
        if cmd:
            console.print(f"{name}={cmd}")
        else:
            spell_failed(f"Alias {name} not found")
    else:
        # Show all aliases
        for alias_name, alias_cmd in executor.state.aliases.items():
            console.print(f"{alias_name}={alias_cmd}")

@app.command()
def history():
    """Show command history"""
    hist = config.load_history()
    for i, cmd in enumerate(hist[-20:], start=max(1, len(hist)-19)):
        console.print(f"{i:4d}  {cmd}")

@app.command()
def script(file_path: str):
    """Execute script file line by line"""
    script_file = Path(file_path)
    
    if not script_file.exists():
        spell_failed(f"Script file not found: {file_path}")
        return
        
    try:
        with open(script_file, 'r') as f:
            lines = f.readlines()
            
        executed = 0
        failed = 0
        
        for line_num, line in enumerate(lines, 1):
            line = line.strip()
            
            # Skip empty lines and comments
            if not line or line.startswith('#'):
                continue
                
            console.print(f"[cyan]Line {line_num}:[/cyan] {line}")
            
            # Parse and execute
            pipeline = parser.parse(line)
            returncode, stdout, stderr = executor.execute_pipeline(pipeline)
            
            if stdout:
                console.print(stdout)
            if stderr:
                console.print(f"[red]{stderr}[/red]")
                
            if returncode == 0:
                executed += 1
            else:
                failed += 1
                console.print(f"[red]Line {line_num} failed with code {returncode}[/red]")
                
        spell_complete(f"Script completed: {executed} successful, {failed} failed")
        
    except Exception as e:
        spell_failed(f"Script execution error: {e}")

@app.command() 
def pipe(command: str):
    """Execute command with advanced pipe parsing"""
    # Just use the regular run command - it handles pipes
    return run.callback(command)

@app.command()
def interactive():
    """Enter interactive shell mode"""
    console.print("[cyan]Entering Lucien interactive shell mode[/cyan]")
    console.print("[dim]Type 'exit' to return to Lucien CLI[/dim]")
    
    # Load aliases and environment from config
    for alias_name, alias_cmd in config.get_aliases().items():
        executor.state.add_alias(alias_name, alias_cmd)
        
    for env_name, env_value in config.get_environment().items():
        executor.state.set_env(env_name, env_value)
    
    try:
        while True:
            # Show prompt with current directory
            current_dir = Path.cwd().name
            prompt_text = f"[bold blue]{current_dir}[/bold blue] > "
            
            try:
                command = typer.prompt("", prompt_suffix=prompt_text, show_default=False)
            except (KeyboardInterrupt, EOFError):
                break
                
            if command.lower().strip() in ['exit', 'quit']:
                break
                
            if command.strip():
                # Add to history
                history_list = config.load_history()
                history_list.append(command)
                config.save_history(history_list)
                
                # Execute command
                run.callback(command)
                
    except Exception as e:
        spell_failed(f"Interactive shell error: {e}")
    
    console.print("[cyan]Exited interactive shell mode[/cyan]")

if __name__ == "__main__":
    app()