#!/usr/bin/env python3
"""
Lucien PlannerAgent CLI Script
Breaks down high-level goals into actionable tasks
"""

import sys
import argparse
import os
from pathlib import Path

# Fix Windows console encoding issues
if sys.platform == 'win32':
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.detach())

# Add current directory to Python path for imports
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

try:
    from agents.planner import PlannerAgent
except ImportError as e:
    print(f"❌ Failed to import PlannerAgent: {e}")
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Break down a goal into actionable tasks using AI planning",
        prog="plan"
    )
    parser.add_argument(
        "goal",
        nargs="*",
        help="The goal to break down into tasks"
    )
    parser.add_argument(
        "--prompt",
        help="Alternative way to specify the goal"
    )
    parser.add_argument(
        "--format",
        choices=["simple", "detailed", "json"],
        default="simple",
        help="Output format"
    )
    
    args = parser.parse_args()
    
    # Get the goal from either positional args or --prompt
    if args.prompt:
        goal = args.prompt
    elif args.goal:
        goal = " ".join(args.goal)
    else:
        print("ERROR: Please provide a goal to plan")
        print("Usage: plan \"build a web app\" or plan --prompt \"create API\"")
        sys.exit(1)
    
    if not goal.strip():
        print("ERROR: Goal cannot be empty")
        sys.exit(1)
    
    try:
        # Initialize the planner agent
        agent = PlannerAgent()
        
        # Generate the plan
        tasks = agent.run(goal.strip())
        
        if not tasks:
            print(f"WARNING: No tasks generated for goal: {goal}")
            sys.exit(1)
        
        # Output results
        if args.format == "json":
            import json
            output = {
                "goal": goal,
                "tasks": tasks,
                "total_tasks": len(tasks)
            }
            print(json.dumps(output, indent=2))
        elif args.format == "detailed":
            print(f">> GOAL: {goal}")
            print("=" * 60)
            print(f">> IMPLEMENTATION PLAN ({len(tasks)} tasks):")
            print()
            for i, task in enumerate(tasks, 1):
                print(f"{i:2d}. {task}")
            print()
            print(">> Ready to start development!")
        else:  # simple format (default)
            print(f">> GOAL: {goal}")
            print(f">> PLAN ({len(tasks)} tasks):")
            for task in tasks:
                print(f"   {task}")
    
    except Exception as e:
        print(f"ERROR: Error generating plan: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()