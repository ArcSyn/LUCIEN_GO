import typer
from cli.commands import summon, cast, daemon

app = typer.Typer()
app.add_typer(summon.app, name="summon")
app.add_typer(cast.app, name="cast")
app.add_typer(daemon.app, name="daemon")
app.add_typer(cast.app, name="cast")

if __name__ == "__main__":
    app()