#!/usr/bin/env python3
"""
Lucien DesignerAgent CLI Script
Generates UI code from natural language descriptions
"""

import sys
import argparse
import os
from pathlib import Path

# Add current directory to Python path for imports
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

try:
    from agents.designer import DesignerAgent
except ImportError as e:
    print(f"❌ Failed to import DesignerAgent: {e}")
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Generate UI code from natural language descriptions",
        prog="design"
    )
    parser.add_argument(
        "description",
        nargs="*",
        help="Description of the UI component to generate"
    )
    parser.add_argument(
        "--prompt",
        help="Alternative way to specify the UI description"
    )
    parser.add_argument(
        "--framework",
        choices=["react", "vue", "angular"],
        default="react",
        help="UI framework to use (default: react)"
    )
    parser.add_argument(
        "--output-dir",
        default="snapmethod/exports",
        help="Output directory for generated code"
    )
    parser.add_argument(
        "--preview",
        action="store_true",
        help="Show generated code without saving to file"
    )
    
    args = parser.parse_args()
    
    # Get the description from either positional args or --prompt
    if args.prompt:
        description = args.prompt
    elif args.description:
        description = " ".join(args.description)
    else:
        print("❌ Error: Please provide a UI description")
        print("Usage: design \"dark login page with neon glow\"")
        sys.exit(1)
    
    if not description.strip():
        print("❌ Error: UI description cannot be empty")
        sys.exit(1)
    
    try:
        # Initialize the designer agent
        agent = DesignerAgent()
        
        # Override output directory if specified
        if args.output_dir != "snapmethod/exports":
            agent.output_dir = Path(args.output_dir)
        
        # Generate the UI code
        print(f"🎨 Generating UI component: {description}")
        print(f"📋 Framework: {args.framework.capitalize()}")
        
        code = agent.generate(description.strip())
        
        if not code:
            print(f"⚠️  Failed to generate code for: {description}")
            sys.exit(1)
        
        if args.preview:
            print("=" * 60)
            print("📄 GENERATED CODE:")
            print("=" * 60)
            print(code)
        else:
            # Code is already saved by the agent, just show success message
            output_file = agent._save_code(code, agent._parse_prompt(description)['component_name'], args.framework)
            print(f"✅ UI component generated successfully!")
            print(f"📁 Saved to: {output_file}")
            print(f"🚀 Component ready for use in your {args.framework} project")
    
    except Exception as e:
        print(f"❌ Error generating UI: {e}")
        import traceback
        if os.getenv('DEBUG'):
            print("Debug traceback:")
            traceback.print_exc()
        sys.exit(1)

if __name__ == "__main__":
    main()