# cli/commands/validate.py
import typer
import sys
from pathlib import Path

# Add parent directory for imports
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from lucien.ui import Console, spell_complete, spell_failed, rune_divider
from core.config import config
from core.claude_memory import memory

app = typer.Typer()
console = Console()

@app.command()
def system():
    """Validate entire Lucien system"""
    rune_divider("System Validation")
    
    results = []
    
    # Test core components
    results.append(_test_config_system())
    results.append(_test_memory_system())
    results.append(_test_agent_system())
    results.append(_test_shell_system())
    
    # Summary
    passed = sum(1 for r in results if r)
    total = len(results)
    
    if passed == total:
        spell_complete(f"All {total} system tests passed!")
    else:
        spell_failed(f"{passed}/{total} tests passed")

def _test_config_system() -> bool:
    """Test configuration system"""
    console.print("[cyan]Testing configuration system...[/cyan]")
    
    try:
        # Test config loading
        test_value = config.get("shell.prompt", "default")
        config.set("test.validation", "success")
        retrieved = config.get("test.validation")
        
        if retrieved == "success":
            console.print("  [green]Config system: PASS[/green]")
            return True
        else:
            console.print("  [red]Config system: FAIL - Value mismatch[/red]")
            return False
            
    except Exception as e:
        console.print(f"  [red]Config system: FAIL - {e}[/red]")
        return False

def _test_memory_system() -> bool:
    """Test Claude memory system"""
    console.print("[cyan]Testing memory system...[/cyan]")
    
    try:
        # Test memory logging
        memory.log_interaction("test input", "test response", "test result")
        context = memory.get_context_summary()
        
        if "test input" in memory.export_memory() and context:
            console.print("  [green]Memory system: PASS[/green]")
            return True
        else:
            console.print("  [red]Memory system: FAIL - Memory not persisted[/red]")
            return False
            
    except Exception as e:
        console.print(f"  [red]Memory system: FAIL - {e}[/red]")
        return False

def _test_agent_system() -> bool:
    """Test agent casting system"""
    console.print("[cyan]Testing agent system...[/cyan]")
    
    try:
        # Test basic agent
        from agents.build import Build
        agent = Build()
        result = agent.run("test")
        
        if result and "Build" in result:
            console.print("  [green]Agent system: PASS[/green]")
            return True
        else:
            console.print("  [red]Agent system: FAIL - Agent not responding[/red]")
            return False
            
    except Exception as e:
        console.print(f"  [red]Agent system: FAIL - {e}[/red]")
        return False

def _test_shell_system() -> bool:
    """Test shell execution system"""
    console.print("[cyan]Testing shell system...[/cyan]")
    
    try:
        from core.shell_parser import ShellParser
        from core.shell_executor import ShellExecutor
        
        parser = ShellParser()
        executor = ShellExecutor()
        
        # Test simple command
        pipeline = parser.parse("echo test")
        returncode, stdout, stderr = executor.execute_pipeline(pipeline)
        
        if returncode == 0 and "test" in stdout:
            console.print("  [green]Shell system: PASS[/green]")
            return True
        else:
            console.print("  [red]Shell system: FAIL - Command execution failed[/red]")
            return False
            
    except Exception as e:
        console.print(f"  [red]Shell system: FAIL - {e}[/red]")
        return False

@app.command()
def agents():
    """Validate all agents"""
    rune_divider("Agent Validation")
    
    agent_files = [
        ("build", "Build Agent"),
        ("analyze", "Analyze Agent"), 
        ("manage", "Manage Agent"),
        ("rewriter", "Rewriter Agent")
    ]
    
    passed = 0
    total = len(agent_files)
    
    for agent_name, description in agent_files:
        console.print(f"[cyan]Testing {description}...[/cyan]")
        
        try:
            module = __import__(f"agents.{agent_name}", fromlist=[agent_name.capitalize()])
            agent_class = getattr(module, agent_name.capitalize())
            agent = agent_class()
            result = agent.run("test validation")
            
            if result and agent_name.capitalize() in result:
                console.print(f"  [green]{description}: PASS[/green]")
                passed += 1
            else:
                console.print(f"  [red]{description}: FAIL - No response[/red]")
                
        except Exception as e:
            console.print(f"  [red]{description}: FAIL - {e}[/red]")
    
    if passed == total:
        spell_complete(f"All {total} agents validated!")
    else:
        spell_failed(f"{passed}/{total} agents working")

@app.command()
def shell():
    """Validate shell functionality"""
    rune_divider("Shell Validation")
    
    from core.shell_parser import ShellParser
    from core.shell_executor import ShellExecutor
    
    parser = ShellParser()
    executor = ShellExecutor()
    
    test_commands = [
        ("echo hello", "Basic command"),
        ("echo hello > test.txt", "Output redirect"),
        ("pwd", "Builtin command"),
        ("set TEST=value", "Variable setting")
    ]
    
    passed = 0
    total = len(test_commands)
    
    for command, description in test_commands:
        console.print(f"[cyan]Testing {description}: {command}[/cyan]")
        
        try:
            pipeline = parser.parse(command)
            returncode, stdout, stderr = executor.execute_pipeline(pipeline)
            
            if returncode == 0:
                console.print(f"  [green]{description}: PASS[/green]")
                passed += 1
            else:
                console.print(f"  [red]{description}: FAIL - Return code {returncode}[/red]")
                
        except Exception as e:
            console.print(f"  [red]{description}: FAIL - {e}[/red]")
    
    # Cleanup
    test_file = Path("test.txt")
    if test_file.exists():
        test_file.unlink()
    
    if passed == total:
        spell_complete(f"All {total} shell tests passed!")
    else:
        spell_failed(f"{passed}/{total} shell tests working")

@app.command()
def config_test():
    """Validate configuration system"""
    rune_divider("Configuration Validation")
    
    # Test config file creation
    console.print("[cyan]Testing config initialization...[/cyan]")
    
    config_dir = Path.home() / ".lucien"
    if config_dir.exists():
        console.print(f"  [green]Config directory exists: {config_dir}[/green]")
    else:
        console.print(f"  [yellow]Creating config directory: {config_dir}[/yellow]")
        config_dir.mkdir(exist_ok=True)
    
    # Test config values
    console.print("[cyan]Testing config operations...[/cyan]")
    
    config.set("validation.test", "success")
    value = config.get("validation.test")
    
    if value == "success":
        console.print("  [green]Config read/write: PASS[/green]")
    else:
        console.print("  [red]Config read/write: FAIL[/red]")
    
    # Show current config
    console.print("\n[dim]Current configuration:[/dim]")
    console.print(config.show_config())

if __name__ == "__main__":
    app()