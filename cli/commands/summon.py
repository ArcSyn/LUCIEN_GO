import typer
from pathlib import Path

app = typer.Typer()

@app.command("agent")
def summon_agent(name: str):
    """
    Summon a new agent by name. Generates a starter agent Python file.
    """
    agent_file = Path(f"agents/{name.lower().replace(' ', '_')}.py")
    agent_file.parent.mkdir(parents=True, exist_ok=True)

    if agent_file.exists():
        typer.echo(f"⚠️ Agent '{name}' already exists at {agent_file}")
        return

    agent_file.write_text(f'''"""
Lucien Agent: {name}
"""

def run(message: str):
    return f"[{name}] I received: '{{message}}'"
'''.strip())

    typer.echo(f"✅ Summoned agent '{name}' at {agent_file}")
