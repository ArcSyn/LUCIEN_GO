# cli/commands/cast.py
import importlib
import typer
from pathlib import Path

from lucien.ui import print_logo, thinking, spell_output, spell_complete, spell_failed, rune_divider

app = typer.Typer()

AGENT_DIR = Path("agents")

@app.command()
def agent(name: str, prompt: str):
    """
    Dynamically cast an agent by name with a prompt.
    Example: python -m cli.main cast agent rewriter "Rewrite this text"
    """
    print_logo()
    rune_divider(f"Casting '{name}' with prompt")
    thinking(f"{name} is thinking...")

    try:
        module_path = f"agents.{name.lower()}"
        module = importlib.import_module(module_path)

        class_name = name.capitalize()
        agent_class = getattr(module, class_name)
        instance = agent_class()

        result = instance.run(prompt)
        spell_output(result)
        spell_complete()

    except Exception as e:
        spell_failed(f"Failed to cast agent '{name}': {e}")