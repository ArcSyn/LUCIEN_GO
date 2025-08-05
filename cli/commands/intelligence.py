# cli/commands/intelligence.py
import typer
import sys
from pathlib import Path

# Add parent directory for imports
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from lucien.ui import console, spell_complete, spell_failed, rune_divider
from core.intelligence import predictor
from core.copilot import copilot
from core.orchestrator import orchestrator

app = typer.Typer()

@app.command()
def predict(input_text: str = ""):
    """Get intelligent command predictions"""
    rune_divider("Command Prediction")
    
    predictions = predictor.predict_next_commands(input_text, limit=10)
    
    if predictions:
        console.print("[cyan]Predicted commands (with confidence):[/cyan]")
        for i, (cmd, confidence) in enumerate(predictions, 1):
            confidence_bar = "█" * int(confidence * 10)
            console.print(f"{i:2d}. [bold]{cmd}[/bold] {confidence_bar} ({confidence:.1%})")
    else:
        console.print("[yellow]No predictions available. Use more commands to train the system.[/yellow]")

@app.command()
def suggestions():
    """Get intelligent suggestions for current context"""
    rune_divider("Context Suggestions")
    
    current_dir = Path.cwd()
    recent_commands = []  # Could load from history
    
    suggestions = predictor.get_intelligent_suggestions(current_dir, recent_commands)
    
    if suggestions:
        console.print("[cyan]Intelligent suggestions for current context:[/cyan]")
        for i, suggestion in enumerate(suggestions, 1):
            console.print(f"{i}. [green]{suggestion}[/green]")
    else:
        console.print("[yellow]No contextual suggestions available.[/yellow]")

@app.command()
def completions(partial: str):
    """Get smart tab completions"""
    completions = predictor.get_smart_completions(partial)
    
    if completions:
        console.print("[cyan]Smart completions:[/cyan]")
        for completion in completions:
            console.print(f"  {completion}")
    else:
        console.print("[yellow]No completions found.[/yellow]")

@app.command()
def start_copilot(directory: str = "."):
    """Start AI pair programming copilot"""
    target_dir = Path(directory).resolve()
    
    if not target_dir.exists():
        spell_failed(f"Directory not found: {directory}")
        return
    
    result = copilot.start_monitoring(target_dir)
    spell_complete(result)

@app.command()
def stop_copilot():
    """Stop AI pair programming copilot"""
    result = copilot.stop_monitoring()
    spell_complete(result)

@app.command()
def copilot_suggestions(file_path: str = None, line: int = None):
    """Get copilot suggestions for current file"""
    rune_divider("AI Copilot Suggestions")
    
    if file_path:
        target_file = Path(file_path)
        if not target_file.exists():
            spell_failed(f"File not found: {file_path}")
            return
    else:
        # Try to find recently modified files
        current_dir = Path.cwd()
        code_files = list(current_dir.glob("*.py")) + list(current_dir.glob("*.js")) + list(current_dir.glob("*.ts"))
        if not code_files:
            spell_failed("No code files found in current directory")
            return
        target_file = max(code_files, key=lambda f: f.stat().st_mtime)
    
    suggestions = copilot.get_contextual_suggestions(target_file, line)
    
    if suggestions:
        console.print(f"[cyan]AI suggestions for {target_file.name}:[/cyan]")
        for suggestion in suggestions:
            priority_color = {"high": "red", "medium": "yellow", "low": "green"}.get(suggestion["priority"], "white")
            console.print(f"[{priority_color}]{suggestion['type'].upper()}[/{priority_color}]: {suggestion['message']}")
            if "suggestion" in suggestion:
                console.print(f"  → {suggestion['suggestion']}")
    else:
        console.print("[green]No issues found - code looks good![/green]")

@app.command()
def analyze_error(error_text: str, language: str = "python"):
    """Analyze error and get AI suggestions"""
    rune_divider("Error Analysis")
    
    analysis = copilot.analyze_error_output(error_text, {"language": language})
    
    console.print(f"[red]Error Type:[/red] {analysis['error_type']}")
    console.print(f"[cyan]Confidence:[/cyan] {analysis['confidence']:.1%}")
    
    if analysis["suggested_fixes"]:
        console.print("[green]Suggested Fixes:[/green]")
        for i, fix in enumerate(analysis["suggested_fixes"], 1):
            console.print(f"  {i}. {fix}")
    else:
        console.print("[yellow]No specific fixes suggested.[/yellow]")

@app.command()
def next_step():
    """Get AI suggestion for next workflow step"""
    suggestion = copilot.suggest_workflow_next_step({})
    
    if suggestion:
        console.print(f"[cyan]Suggested next step:[/cyan] {suggestion}")
    else:
        console.print("[green]All caught up! No immediate actions needed.[/green]")

@app.command()
def analyze_project(path: str = "."):
    """Deep project analysis with AI insights"""
    rune_divider("Project Analysis")
    
    project_path = Path(path).resolve()
    context = orchestrator.analyze_project(project_path)
    
    console.print(f"[bold]Project:[/bold] {context.root_path.name}")
    console.print(f"[bold]Type:[/bold] {context.project_type}")
    console.print(f"[bold]Build System:[/bold] {context.build_system}")
    console.print(f"[bold]Test Framework:[/bold] {context.test_framework}")
    console.print(f"[bold]Deployment:[/bold] {context.deployment_target}")
    
    if context.git_info["is_repo"]:
        console.print(f"[bold]Git Branch:[/bold] {context.git_info['branch']}")
        console.print(f"[bold]Git Status:[/bold] {context.git_info['status']}")
    
    console.print(f"[bold]Files:[/bold] {context.performance_metrics['file_count']}")
    console.print(f"[bold]Lines of Code:[/bold] {context.performance_metrics['loc']:,}")
    console.print(f"[bold]Complexity:[/bold] {context.performance_metrics['complexity_estimate']}")

@app.command()
def project_insights(path: str = "."):
    """Get comprehensive project insights and recommendations"""
    rune_divider("Project Insights")
    
    project_path = Path(path).resolve()
    insights = orchestrator.get_project_insights(project_path)
    
    # Health status
    health_color = {
        "excellent": "green",
        "good": "cyan", 
        "fair": "yellow",
        "poor": "red"
    }.get(insights["project_health"], "white")
    
    console.print(f"[bold]Project Health:[/bold] [{health_color}]{insights['project_health'].upper()}[/{health_color}]")
    
    # Recommendations
    if insights["recommendations"]:
        console.print("\n[cyan]Recommendations:[/cyan]")
        for i, rec in enumerate(insights["recommendations"], 1):
            console.print(f"  {i}. {rec}")
    
    # Optimization opportunities
    if insights["optimization_opportunities"]:
        console.print("\n[yellow]Optimization Opportunities:[/yellow]")
        for i, opt in enumerate(insights["optimization_opportunities"], 1):
            console.print(f"  {i}. {opt}")
    
    # Risk factors
    if insights["risk_factors"]:
        console.print("\n[red]Risk Factors:[/red]")
        for i, risk in enumerate(insights["risk_factors"], 1):
            console.print(f"  {i}. {risk}")

@app.command()
def create_workflow(goal: str):
    """Create adaptive workflow for specific goal"""
    rune_divider(f"Creating Workflow: {goal}")
    
    workflow_id = orchestrator.create_adaptive_workflow(goal)
    
    workflow = orchestrator.active_workflows[workflow_id]
    
    console.print(f"[green]Created workflow:[/green] {workflow_id}")
    console.print(f"[cyan]Goal:[/cyan] {goal}")
    console.print(f"[cyan]Steps:[/cyan] {len(workflow['steps'])}")
    console.print(f"[cyan]Estimated Duration:[/cyan] {workflow['estimated_duration']} seconds")
    
    console.print("\n[bold]Workflow Steps:[/bold]")
    for i, step in enumerate(workflow['steps'], 1):
        console.print(f"  {i}. {step.name} ({step.estimated_duration}s)")
    
    console.print(f"\n[yellow]Use 'lucien intelligence execute-workflow {workflow_id}' to run[/yellow]")

@app.command()
def execute_workflow(workflow_id: str):
    """Execute adaptive workflow"""
    import asyncio
    
    if workflow_id not in orchestrator.active_workflows:
        spell_failed(f"Workflow not found: {workflow_id}")
        return
    
    rune_divider(f"Executing Workflow: {workflow_id}")
    
    # Run the async workflow
    results = asyncio.run(orchestrator.execute_workflow(workflow_id))
    
    if results["success"]:
        spell_complete(f"Workflow completed successfully in {results['duration']:.1f} seconds")
    else:
        spell_failed(f"Workflow failed after {results['duration']:.1f} seconds")
    
    # Show step results
    console.print("\n[bold]Step Results:[/bold]")
    for step_result in results["steps"]:
        status_color = "green" if step_result["success"] else "red"
        status_icon = "✓" if step_result["success"] else "✗"
        console.print(f"  [{status_color}]{status_icon}[/{status_color}] {step_result['step']} ({step_result['duration']:.1f}s)")
        
        if step_result.get("output"):
            console.print(f"    Output: {step_result['output'][:100]}...")

@app.command()
def list_workflows():
    """List active workflows"""
    if not orchestrator.active_workflows:
        console.print("[yellow]No active workflows[/yellow]")
        return
    
    rune_divider("Active Workflows")
    
    for workflow_id, workflow in orchestrator.active_workflows.items():
        state_color = {
            "planning": "cyan",
            "executing": "yellow", 
            "completed": "green",
            "failed": "red"
        }.get(workflow["state"].value, "white")
        
        console.print(f"[bold]{workflow_id}[/bold]")
        console.print(f"  Goal: {workflow['goal']}")
        console.print(f"  State: [{state_color}]{workflow['state'].value}[/{state_color}]")
        console.print(f"  Progress: {workflow.get('progress', 0):.1f}%")
        console.print()

if __name__ == "__main__":
    app()