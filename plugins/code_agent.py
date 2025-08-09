#!/usr/bin/env python3
"""
Lucien CoderAgent CLI Script
Generates, refactors, and explains code using AI assistance
"""

import sys
import argparse
import os
from pathlib import Path

# Add current directory to Python path for imports
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

try:
    from agents.coder import CoderAgent
except ImportError as e:
    print(f"❌ Failed to import CoderAgent: {e}")
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Generate, refactor, and explain code using AI assistance",
        prog="code"
    )
    
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Generate command
    gen_parser = subparsers.add_parser("generate", aliases=["gen"], help="Generate code from description")
    gen_parser.add_argument("description", nargs="*", help="Description of code to generate")
    gen_parser.add_argument("--prompt", help="Alternative way to specify description")
    gen_parser.add_argument("--language", "--lang", choices=["python", "javascript", "go", "java", "rust"], help="Programming language")
    gen_parser.add_argument("--output", "-o", help="Output file path")
    gen_parser.add_argument("--type", choices=["function", "class", "script", "api"], help="Type of code to generate")
    
    # Refactor command
    ref_parser = subparsers.add_parser("refactor", aliases=["ref"], help="Refactor existing code")
    ref_parser.add_argument("file", help="File to refactor")
    ref_parser.add_argument("--type", choices=["improve", "optimize", "modernize"], default="improve", help="Type of refactoring")
    ref_parser.add_argument("--output", "-o", help="Output file (default: overwrite original)")
    ref_parser.add_argument("--backup", action="store_true", help="Create backup of original file")
    
    # Explain command
    exp_parser = subparsers.add_parser("explain", aliases=["exp"], help="Explain what code does")
    exp_parser.add_argument("target", help="File path or code snippet to explain")
    exp_parser.add_argument("--format", choices=["markdown", "plain"], default="plain", help="Output format")
    
    # Review command (alias for explain)
    rev_parser = subparsers.add_parser("review", help="Review and explain code")
    rev_parser.add_argument("file", help="File to review")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        sys.exit(1)
    
    try:
        agent = CoderAgent()
        
        if args.command in ["generate", "gen"]:
            handle_generate(agent, args)
        elif args.command in ["refactor", "ref"]:
            handle_refactor(agent, args)
        elif args.command in ["explain", "exp"]:
            handle_explain(agent, args)
        elif args.command == "review":
            handle_review(agent, args)
    
    except Exception as e:
        print(f"❌ Error: {e}")
        if os.getenv('DEBUG'):
            import traceback
            traceback.print_exc()
        sys.exit(1)

def handle_generate(agent: 'CoderAgent', args):
    """Handle code generation"""
    # Get description
    if args.prompt:
        description = args.prompt
    elif args.description:
        description = " ".join(args.description)
    else:
        print("❌ Error: Please provide a code description")
        print("Usage: code generate \"create a function that sorts a list\"")
        sys.exit(1)
    
    if not description.strip():
        print("❌ Error: Code description cannot be empty")
        sys.exit(1)
    
    print(f"🤖 Generating code: {description}")
    if args.language:
        print(f"📋 Language: {args.language}")
    
    # Generate code
    code = agent.generate(
        prompt=description, 
        language=args.language, 
        output_file=args.output
    )
    
    if not code:
        print("⚠️  Failed to generate code")
        sys.exit(1)
    
    if args.output:
        print(f"✅ Code generated successfully!")
        print(f"📁 Saved to: {args.output}")
    else:
        print("=" * 60)
        print("📄 GENERATED CODE:")
        print("=" * 60)
        print(code)

def handle_refactor(agent: 'CoderAgent', args):
    """Handle code refactoring"""
    if not os.path.exists(args.file):
        print(f"❌ Error: File not found: {args.file}")
        sys.exit(1)
    
    print(f"🔧 Refactoring file: {os.path.basename(args.file)}")
    print(f"📋 Type: {args.type}")
    
    # Create backup if requested
    if args.backup:
        backup_path = args.file + ".backup"
        import shutil
        shutil.copy2(args.file, backup_path)
        print(f"💾 Backup created: {backup_path}")
    
    # Refactor code
    refactored_code = agent.refactor(args.file, args.type)
    
    if not refactored_code:
        print("⚠️  Failed to refactor code")
        sys.exit(1)
    
    # Determine output file
    output_file = args.output or args.file
    
    # Save refactored code
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(refactored_code)
    
    print(f"✅ Code refactored successfully!")
    print(f"📁 Saved to: {output_file}")

def handle_explain(agent: 'CoderAgent', args):
    """Handle code explanation"""
    print(f"📖 Analyzing code...")
    
    explanation = agent.explain(args.target)
    
    if not explanation:
        print("⚠️  Failed to generate explanation")
        sys.exit(1)
    
    print(explanation)

def handle_review(agent: 'CoderAgent', args):
    """Handle code review (alias for explain)"""
    if not os.path.exists(args.file):
        print(f"❌ Error: File not found: {args.file}")
        sys.exit(1)
    
    print(f"📖 Reviewing file: {os.path.basename(args.file)}")
    
    explanation = agent.explain(args.file)
    
    if not explanation:
        print("⚠️  Failed to generate review")
        sys.exit(1)
    
    print(explanation)

if __name__ == "__main__":
    main()