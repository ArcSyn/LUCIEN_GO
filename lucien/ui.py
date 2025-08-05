# lucien/ui.py
from rich.console import Console
from rich.panel import Panel
from rich.text import Text
from rich.progress import Progress
from rich.status import Status
from time import sleep
from sys import stdout

console = Console()

def print_logo():
    ascii_art = r"""
 _____       __    __     ____    _____    _____      __      _  
(_   _)      ) )  ( (    / ___)  (_   _)  / ___/     /  \    / ) 
  | |       ( (    ) )  / /        | |   ( (__      / /\ \  / /  
  | |        ) )  ( (  ( (         | |    ) __)     ) ) ) ) ) )  
  | |   __  ( (    ) ) ( (         | |   ( (       ( ( ( ( ( (   
__| |___) )  ) \__/ (   \ \___    _| |__  \ \___   / /  \ \/ /   
\________/   \______/    \____)  /_____(   \____\ (_/    \__/    
"""
    styled_text = Text(ascii_art, style="bold blue", justify="center")
    panel = Panel(
        styled_text,
        title="[bold magenta]Lucien CLI[/bold magenta]",
        subtitle="[cyan]Modular Spellcasting Shell",
        border_style="bright_blue"
    )
    console.print(panel)

def thinking(message="Casting spell..."):
    with console.status(f"[bold cyan]{message}[/bold cyan]", spinner="bouncingBar"):
        for _ in range(30):
            sleep(0.05)

def spell_complete(message="Spell complete!"):
    console.print(f"[bold green]{message}[/bold green]")

def spell_failed(message="Spell failed.") -> None:
    console.print(f"[bold red]{message}[/bold red]")

def typewriter(text: str, delay=0.02, style="white"):
    """Print text like a typewriter effect with optional color."""
    for char in text:
        console.print(char, style=style, end="")
        stdout.flush()
        sleep(delay)
    print()

def spell_output(text: str) -> None:
    """Show the result of the spell with a typewriter effect inside a panel."""
    output = Text()
    for char in text:
        output.append(char, style="white")
        sleep(0.01)

    panel = Panel.fit(
        output,
        title="[bold green]Spell Output[/bold green]",
        border_style="bright_cyan"
    )
    console.print(panel)

def rune_divider(label="Arcane Divider") -> None:
    console.rule(f"[bold blue]{label}")

def loading_bar(task_desc="Preparing incantation...") -> None:
    with Progress() as progress:
        task = progress.add_task(f"[cyan]{task_desc}", total=100)
        for _ in range(100):
            sleep(0.01)
            progress.update(task, advance=1)