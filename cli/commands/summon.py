import typer
app = typer.Typer()

@app.command()
def agent(name: str):
    typer.echo(f"Summoning agent: {name}")