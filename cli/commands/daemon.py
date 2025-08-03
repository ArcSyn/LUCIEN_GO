# cli/commands/daemon.py

import typer

app = typer.Typer()

@app.command()
def start():
    typer.echo("🔧 Starting Lucien daemon...")

@app.command()
def stop():
    typer.echo("🛑 Stopping Lucien daemon...")