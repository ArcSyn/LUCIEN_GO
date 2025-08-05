import typer
from cli.commands import summon, cast, daemon, shell, validate, intelligence, powershell_parity

app = typer.Typer()
app.add_typer(summon.app, name="summon")
app.add_typer(cast.app, name="cast") 
app.add_typer(daemon.app, name="daemon")
app.add_typer(shell.app, name="shell")
app.add_typer(validate.app, name="validate")
app.add_typer(intelligence.app, name="ai")
app.add_typer(powershell_parity.app, name="ps")

# Add exec command for backward compatibility
@app.command()
def exec(command: str):
    """Execute shell commands (legacy - use 'shell run' instead)"""
    from cli.commands.shell import run
    return run.callback(command)

# Direct shell mode for replacing terminal
@app.command()
def interactive():
    """Start Lucien in interactive shell mode"""
    from cli.commands.shell import interactive
    return interactive.callback()

if __name__ == "__main__":
    app()