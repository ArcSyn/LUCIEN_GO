from rich.console import Console
from rich.panel import Panel
from rich.text import Text

console = Console()

def spell_thinking(message: str):
    """Show thinking indicator"""
    console.print(f"[cyan]Thinking...[/cyan] {message}")

def spell_complete(message: str):
    """Show completion indicator"""
    console.print(f"[green]Complete:[/green] {message}")

def spell_failed(message: str):
    """Show failure indicator"""
    console.print(f"[red]Failed:[/red] {message}")