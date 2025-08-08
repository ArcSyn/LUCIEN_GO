#!/usr/bin/env python3
"""
Lucien ReviewAgent CLI Script  
Analyzes code files and provides improvement suggestions
"""

import sys
import argparse
import os
from pathlib import Path

# Add current directory to Python path for imports
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

try:
    from agents.reviewer import ReviewAgent
except ImportError as e:
    print(f"❌ Failed to import ReviewAgent: {e}")
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Analyze code files and provide improvement suggestions",
        prog="review"
    )
    parser.add_argument(
        "file",
        nargs="?",
        help="Path to the file to review"
    )
    parser.add_argument(
        "--file",
        dest="file_path",
        help="Alternative way to specify the file path"
    )
    parser.add_argument(
        "--format",
        choices=["markdown", "plain", "json"],
        default="plain",
        help="Output format (default: plain)"
    )
    parser.add_argument(
        "--severity",
        choices=["all", "high", "medium", "low"],
        default="all",
        help="Minimum severity level to show"
    )
    parser.add_argument(
        "--category",
        choices=["all", "security", "performance", "maintainability", "style", "bug"],
        default="all",
        help="Filter by issue category"
    )
    
    args = parser.parse_args()
    
    # Get the file path from either positional arg or --file
    file_path = args.file or args.file_path
    
    if not file_path:
        print("❌ Error: Please provide a file to review")
        print("Usage: review <file_path> or review --file <file_path>")
        sys.exit(1)
    
    # Convert to absolute path if relative
    file_path = os.path.abspath(file_path)
    
    if not os.path.exists(file_path):
        print(f"❌ Error: File not found: {file_path}")
        sys.exit(1)
    
    if not os.path.isfile(file_path):
        print(f"❌ Error: Path is not a file: {file_path}")
        sys.exit(1)
    
    try:
        # Initialize the review agent
        agent = ReviewAgent()
        
        print(f"🔍 Reviewing file: {os.path.basename(file_path)}")
        print(f"📁 Path: {file_path}")
        print()
        
        # Analyze the file
        review_result = agent.analyze(file_path)
        
        if not review_result:
            print("⚠️  No review results generated")
            sys.exit(1)
        
        # Output the results
        if args.format == "json":
            # For JSON output, we'd need to modify the agent to return structured data
            # For now, wrap the text result
            import json
            output = {
                "file": file_path,
                "review": review_result,
                "format": "text"
            }
            print(json.dumps(output, indent=2))
        elif args.format == "markdown":
            # Convert plain text markers to markdown
            markdown_result = review_result
            markdown_result = markdown_result.replace("**", "**")  # Keep bold
            markdown_result = markdown_result.replace("✅", "✅")   # Keep emojis
            markdown_result = markdown_result.replace("❌", "❌")
            markdown_result = markdown_result.replace("⚠️", "⚠️")
            markdown_result = markdown_result.replace("💡", "💡")
            markdown_result = markdown_result.replace("🎨", "🎨")
            print(markdown_result)
        else:  # plain format (default)
            print(review_result)
    
    except Exception as e:
        print(f"❌ Error reviewing file: {e}")
        import traceback
        if os.getenv('DEBUG'):
            print("Debug traceback:")
            traceback.print_exc()
        sys.exit(1)

if __name__ == "__main__":
    main()