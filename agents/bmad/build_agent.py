from agents.agent_base import Agent
import subprocess
import os
import json
from pathlib import Path

class BuildAgent(Agent):
    """BMAD Build Agent - Handles project building, compilation, and dependencies"""
    
    def run(self, input_text: str) -> str:
        """
        Build operations: compile, install deps, setup env
        Examples:
        - "build python project"  
        - "install dependencies"
        - "setup node environment"
        """
        input_lower = input_text.lower().strip()
        
        # Detect project type
        cwd = Path.cwd()
        
        if "python" in input_lower or (cwd / "requirements.txt").exists() or (cwd / "setup.py").exists():
            return self._build_python_project()
        elif "node" in input_lower or "npm" in input_lower or (cwd / "package.json").exists():
            return self._build_node_project()
        elif "go" in input_lower or (cwd / "go.mod").exists():
            return self._build_go_project()
        elif "rust" in input_lower or (cwd / "Cargo.toml").exists():
            return self._build_rust_project()
        elif "install" in input_lower and "dep" in input_lower:
            return self._install_dependencies()
        else:
            return self._generic_build(input_text)
    
    def _build_python_project(self) -> str:
        """Build Python project"""
        results = []
        cwd = Path.cwd()
        
        # Check for requirements.txt
        if (cwd / "requirements.txt").exists():
            try:
                result = subprocess.run(
                    ["pip", "install", "-r", "requirements.txt"], 
                    capture_output=True, text=True, timeout=120
                )
                if result.returncode == 0:
                    results.append("[OK] Dependencies installed successfully")
                else:
                    results.append(f"[ERROR] Dependency installation failed: {result.stderr}")
            except Exception as e:
                results.append(f"[ERROR] Error installing dependencies: {e}")
        
        # Check for setup.py
        if (cwd / "setup.py").exists():
            try:
                result = subprocess.run(
                    ["python", "setup.py", "build"], 
                    capture_output=True, text=True, timeout=60
                )
                if result.returncode == 0:
                    results.append("[OK] Python package built successfully")
                else:
                    results.append(f"[ERROR] Build failed: {result.stderr}")
            except Exception as e:
                results.append(f"[ERROR] Build error: {e}")
        
        if not results:
            results.append("[INFO] No Python build files found (requirements.txt, setup.py)")
            
        return f"[BuildAgent] Python Build Results:\n" + "\n".join(results)
    
    def _build_node_project(self) -> str:
        """Build Node.js project"""
        results = []
        cwd = Path.cwd()
        
        if (cwd / "package.json").exists():
            # Install dependencies
            try:
                result = subprocess.run(
                    ["npm", "install"], 
                    capture_output=True, text=True, timeout=180
                )
                if result.returncode == 0:
                    results.append("[OK] npm install completed")
                else:
                    results.append(f"[ERROR] npm install failed: {result.stderr}")
            except Exception as e:
                results.append(f"[ERROR] npm install error: {e}")
            
            # Try to build if build script exists
            try:
                with open(cwd / "package.json") as f:
                    package_data = json.load(f)
                    if "scripts" in package_data and "build" in package_data["scripts"]:
                        result = subprocess.run(
                            ["npm", "run", "build"],
                            capture_output=True, text=True, timeout=120
                        )
                        if result.returncode == 0:
                            results.append("✅ npm run build completed")
                        else:
                            results.append(f"❌ npm run build failed: {result.stderr}")
            except Exception as e:
                results.append(f"❌ Build script error: {e}")
        else:
            results.append("🔍 No package.json found")
            
        return f"[BuildAgent] Node.js Build Results:\n" + "\n".join(results)
    
    def _build_go_project(self) -> str:
        """Build Go project"""
        results = []
        
        try:
            # Get dependencies
            result = subprocess.run(
                ["go", "mod", "tidy"], 
                capture_output=True, text=True, timeout=60
            )
            if result.returncode == 0:
                results.append("✅ go mod tidy completed")
            else:
                results.append(f"❌ go mod tidy failed: {result.stderr}")
                
            # Build project
            result = subprocess.run(
                ["go", "build", "."], 
                capture_output=True, text=True, timeout=60
            )
            if result.returncode == 0:
                results.append("✅ Go build completed")
            else:
                results.append(f"❌ Go build failed: {result.stderr}")
                
        except Exception as e:
            results.append(f"❌ Go build error: {e}")
            
        return f"[BuildAgent] Go Build Results:\n" + "\n".join(results)
    
    def _build_rust_project(self) -> str:
        """Build Rust project"""
        results = []
        
        try:
            result = subprocess.run(
                ["cargo", "build"], 
                capture_output=True, text=True, timeout=120
            )
            if result.returncode == 0:
                results.append("✅ Cargo build completed")
            else:
                results.append(f"❌ Cargo build failed: {result.stderr}")
                
        except Exception as e:
            results.append(f"❌ Cargo build error: {e}")
            
        return f"[BuildAgent] Rust Build Results:\n" + "\n".join(results)
    
    def _install_dependencies(self) -> str:
        """Generic dependency installation"""
        cwd = Path.cwd()
        results = []
        
        # Try different package managers
        if (cwd / "requirements.txt").exists():
            results.append(self._build_python_project())
        if (cwd / "package.json").exists():
            results.append(self._build_node_project())
        if (cwd / "go.mod").exists():
            results.append(self._build_go_project()) 
        if (cwd / "Cargo.toml").exists():
            results.append(self._build_rust_project())
            
        if not results:
            results.append("🔍 No recognized dependency files found")
            
        return f"[BuildAgent] Dependency Installation:\n" + "\n".join(results)
    
    def _generic_build(self, command: str) -> str:
        """Handle generic build commands"""
        try:
            result = subprocess.run(
                command, shell=True, 
                capture_output=True, text=True, timeout=60
            )
            if result.returncode == 0:
                return f"[BuildAgent] ✅ Command succeeded:\n{result.stdout}"
            else:
                return f"[BuildAgent] ❌ Command failed:\n{result.stderr}"
        except Exception as e:
            return f"[BuildAgent] ❌ Error executing '{command}': {e}"