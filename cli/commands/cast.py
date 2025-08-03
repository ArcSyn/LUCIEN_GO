ChatGPT said:
🔥 Perfect. Lucien just leveled up.

Now let’s build the cast command — the spellcasting layer that lets you invoke any agent from the command line like this:

bash
Copy
Edit
python -m cli.main cast rewriter "Make this message sound elegant."
✅ Step-by-Step: Build cast.py
1. Open the file:
powershell
Copy
Edit
notepad cli/commands/cast.py
2. Paste this code:
python
Copy
Edit
import importlib
import typer
from pathlib import Path

app = typer.Typer()

AGENT_DIR = Path("agents")

@app.command()
def agent(name: str, prompt: str):
    """
    Dynamically cast an agent by name with a prompt.
    Example: python -m cli.main cast agent rewriter "Rewrite this text"
    """
    try:
        module_path = f"agents.{name.lower()}"
        module = importlib.import_module(module_path)

        class_name = name.capitalize()
        agent_class = getattr(module, class_name)
        instance = agent_class()

        result = instance.run(prompt)
        typer.echo(result)

    except Exception as e:
        typer.secho(f"⚠️ Failed to cast agent '{name}': {e}", fg=typer.colors.RED)