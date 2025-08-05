from agents.agent_base import Agent
import subprocess
from pathlib import Path

class Build(Agent):
    def run(self, input_text: str) -> str:
        """Simple build agent for testing"""
        cwd = Path.cwd()
        
        if (cwd / "requirements.txt").exists():
            try:
                result = subprocess.run(
                    ["pip", "install", "-r", "requirements.txt"], 
                    capture_output=True, text=True, timeout=60
                )
                if result.returncode == 0:
                    return "[Build] Python dependencies installed successfully"
                else:
                    return f"[Build] Failed to install dependencies: {result.stderr}"
            except Exception as e:
                return f"[Build] Error: {e}"
        
        elif (cwd / "package.json").exists():
            try:
                result = subprocess.run(
                    ["npm", "install"], 
                    capture_output=True, text=True, timeout=60
                )
                if result.returncode == 0:
                    return "[Build] Node.js dependencies installed successfully"
                else:
                    return f"[Build] Failed to install dependencies: {result.stderr}"
            except Exception as e:
                return f"[Build] Error: {e}"
        
        else:
            return "[Build] No recognized project files found (requirements.txt, package.json)"