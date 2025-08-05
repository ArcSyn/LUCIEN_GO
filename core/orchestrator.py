# core/orchestrator.py
"""
REVOLUTIONARY FEATURE 3: Adaptive Workflow Orchestration & Project Consciousness
Understands your entire project ecosystem and orchestrates complex multi-step workflows automatically
"""
import json
import time
import asyncio
import threading
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Set, Tuple, Any, Callable
import subprocess
import hashlib
from dataclasses import dataclass
from enum import Enum

from .config import config
from .claude_memory import memory
from .intelligence import predictor

class WorkflowState(Enum):
    IDLE = "idle"
    PLANNING = "planning"
    EXECUTING = "executing"
    MONITORING = "monitoring"
    ADAPTING = "adapting"
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class ProjectContext:
    """Complete project understanding"""
    root_path: Path
    project_type: str
    dependencies: Dict[str, Any]
    build_system: str
    test_framework: str
    deployment_target: str
    git_info: Dict[str, Any]
    team_patterns: Dict[str, Any]
    performance_metrics: Dict[str, Any]
    last_analyzed: datetime

@dataclass 
class WorkflowStep:
    """Individual workflow step"""
    id: str
    name: str
    command: str
    dependencies: List[str]
    success_criteria: List[Callable]
    failure_recovery: List[str]
    estimated_duration: int
    priority: int

class AdaptiveOrchestrator:
    """AI-powered workflow orchestration system"""
    
    def __init__(self):
        self.orchestrator_dir = config.config_dir / "orchestrator"
        self.orchestrator_dir.mkdir(exist_ok=True)
        
        self.projects = {}  # project_path -> ProjectContext
        self.active_workflows = {}  # workflow_id -> workflow_state
        self.workflow_templates = {}  # workflow_type -> template
        self.performance_history = {}  # command -> performance_data
        
        self.state = WorkflowState.IDLE
        self.current_workflow = None
        self.monitoring_thread = None
        
        self._load_workflow_templates()
        self._load_performance_history()
        
    def _load_workflow_templates(self):
        """Load pre-built workflow templates"""
        self.workflow_templates = {
            "python_setup": [
                WorkflowStep("venv", "Create virtual environment", "python -m venv venv", [], [], [], 10, 1),
                WorkflowStep("activate", "Activate virtual environment", "source venv/bin/activate", ["venv"], [], [], 2, 1),
                WorkflowStep("install", "Install dependencies", "pip install -r requirements.txt", ["activate"], [], [], 30, 2),
                WorkflowStep("test", "Run tests", "pytest", ["install"], [], [], 15, 3)
            ],
            "node_setup": [
                WorkflowStep("install", "Install dependencies", "npm install", [], [], [], 45, 1),
                WorkflowStep("build", "Build project", "npm run build", ["install"], [], [], 20, 2),
                WorkflowStep("test", "Run tests", "npm test", ["install"], [], [], 10, 3)
            ],
            "git_workflow": [
                WorkflowStep("status", "Check git status", "git status", [], [], [], 2, 1),
                WorkflowStep("pull", "Pull latest changes", "git pull", ["status"], [], [], 5, 2),
                WorkflowStep("install_deps", "Update dependencies", "auto", ["pull"], [], [], 30, 3),
                WorkflowStep("test", "Run tests", "auto", ["install_deps"], [], [], 20, 4)
            ],
            "deployment": [
                WorkflowStep("test", "Run full test suite", "auto", [], [], [], 60, 1),
                WorkflowStep("build", "Build for production", "auto", ["test"], [], [], 45, 2),
                WorkflowStep("package", "Package application", "auto", ["build"], [], [], 15, 3),
                WorkflowStep("deploy", "Deploy to target", "auto", ["package"], [], [], 30, 4)
            ]
        }
    
    def _load_performance_history(self):
        """Load performance history for optimization"""
        history_file = self.orchestrator_dir / "performance_history.json"
        if history_file.exists():
            try:
                with open(history_file, 'r') as f:
                    self.performance_history = json.load(f)
            except:
                pass
    
    def _save_performance_history(self):
        """Save performance metrics for future optimization"""
        history_file = self.orchestrator_dir / "performance_history.json"
        try:
            with open(history_file, 'w') as f:
                json.dump(self.performance_history, f, indent=2)
        except:
            pass
    
    def analyze_project(self, project_path: Path = None) -> ProjectContext:
        """Deep analysis of project structure and requirements"""
        if not project_path:
            project_path = Path.cwd()
            
        # Check if we have recent analysis
        if str(project_path) in self.projects:
            context = self.projects[str(project_path)]
            if (datetime.now() - context.last_analyzed).seconds < 3600:  # 1 hour cache
                return context
        
        # Perform deep project analysis
        context = ProjectContext(
            root_path=project_path,
            project_type=self._detect_project_type(project_path),
            dependencies=self._analyze_dependencies(project_path),
            build_system=self._detect_build_system(project_path),
            test_framework=self._detect_test_framework(project_path),
            deployment_target=self._detect_deployment_target(project_path),
            git_info=self._analyze_git_info(project_path),
            team_patterns=self._analyze_team_patterns(project_path),
            performance_metrics=self._analyze_performance_metrics(project_path),
            last_analyzed=datetime.now()
        )
        
        self.projects[str(project_path)] = context
        return context
    
    def _detect_project_type(self, path: Path) -> str:
        """Detect comprehensive project type"""
        indicators = {
            "python_webapp": ["app.py", "manage.py", "requirements.txt", "templates/"],
            "python_package": ["setup.py", "src/", "__init__.py"],
            "python_data": ["jupyter notebooks", "data/", "models/", ".ipynb"],
            "node_webapp": ["package.json", "src/", "public/", "index.html"],
            "node_api": ["package.json", "routes/", "controllers/", "middleware/"],
            "react_app": ["package.json", "src/", "public/", "App.js"],
            "go_service": ["go.mod", "main.go", "cmd/", "internal/"],
            "rust_project": ["Cargo.toml", "src/main.rs"],
            "docker_project": ["Dockerfile", "docker-compose.yml"],
            "kubernetes": ["*.yaml", "*.yml", "kustomization.yaml"]
        }
        
        files = set(str(f.name) for f in path.rglob("*") if f.is_file())
        dirs = set(str(d.name) for d in path.rglob("*") if d.is_dir())
        all_items = files | dirs
        
        scores = {}
        for project_type, indicators_list in indicators.items():
            score = sum(1 for indicator in indicators_list if 
                       any(item for item in all_items if indicator.replace("/", "") in item))
            scores[project_type] = score
        
        return max(scores, key=scores.get) if scores else "unknown"
    
    def _analyze_dependencies(self, path: Path) -> Dict[str, Any]:
        """Analyze project dependencies and their health"""
        deps = {"files": [], "outdated": [], "security_issues": [], "conflicts": []}
        
        # Check different dependency files
        dep_files = {
            "requirements.txt": "python",
            "package.json": "node",
            "go.mod": "go",
            "Cargo.toml": "rust",
            "pom.xml": "java"
        }
        
        for dep_file, lang in dep_files.items():
            if (path / dep_file).exists():
                deps["files"].append(dep_file)
                
                # Analyze dependency health (simplified)
                if lang == "python" and dep_file == "requirements.txt":
                    try:
                        result = subprocess.run(
                            ["pip", "list", "--outdated", "--format=json"],
                            capture_output=True, text=True, timeout=10
                        )
                        if result.returncode == 0:
                            outdated = json.loads(result.stdout)
                            deps["outdated"] = [pkg["name"] for pkg in outdated]
                    except:
                        pass
                
                elif lang == "node" and dep_file == "package.json":
                    try:
                        result = subprocess.run(
                            ["npm", "outdated", "--json"],
                            capture_output=True, text=True, timeout=10
                        )
                        if result.returncode == 0:
                            outdated = json.loads(result.stdout)
                            deps["outdated"] = list(outdated.keys())
                    except:
                        pass
        
        return deps
    
    def _detect_build_system(self, path: Path) -> str:
        """Detect build system used"""
        build_indicators = {
            "make": ["Makefile", "makefile"],
            "cmake": ["CMakeLists.txt"],
            "gradle": ["build.gradle", "gradle.properties"],
            "maven": ["pom.xml"],
            "setuptools": ["setup.py"],
            "poetry": ["pyproject.toml"],
            "npm": ["package.json"],
            "cargo": ["Cargo.toml"],
            "go": ["go.mod"]
        }
        
        for build_system, files in build_indicators.items():
            if any((path / f).exists() for f in files):
                return build_system
        
        return "unknown"
    
    def _detect_test_framework(self, path: Path) -> str:
        """Detect testing framework"""
        test_indicators = {
            "pytest": ["pytest.ini", "tests/", "*_test.py"],
            "unittest": ["test_*.py", "*_test.py"],
            "jest": ["jest.config.js", "*.test.js"],
            "mocha": ["mocha.opts", "*.spec.js"],
            "go_test": ["*_test.go"],
            "cargo_test": ["tests/", "#[test]"],
        }
        
        files = list(path.rglob("*"))
        
        for framework, indicators in test_indicators.items():
            if any(any(f.match(pattern) for f in files) for pattern in indicators):
                return framework
        
        return "unknown"
    
    def _detect_deployment_target(self, path: Path) -> str:
        """Detect deployment target/platform"""
        deployment_indicators = {
            "docker": ["Dockerfile", "docker-compose.yml"],
            "kubernetes": ["*.yaml", "kustomization.yaml"],
            "heroku": ["Procfile", "runtime.txt"],
            "aws": [".aws/", "serverless.yml", "sam.yaml"],
            "vercel": ["vercel.json", ".vercel/"],
            "netlify": ["netlify.toml", "_redirects"]
        }
        
        for target, indicators in deployment_indicators.items():
            if any((path / indicator).exists() or 
                   any(path.glob(indicator)) for indicator in indicators):
                return target
        
        return "local"
    
    def _analyze_git_info(self, path: Path) -> Dict[str, Any]:
        """Analyze git repository information"""
        git_info = {"is_repo": False, "branch": None, "remote": None, "status": "clean"}
        
        if (path / ".git").exists():
            git_info["is_repo"] = True
            
            try:
                # Get current branch
                result = subprocess.run(
                    ["git", "branch", "--show-current"],
                    cwd=path, capture_output=True, text=True, timeout=5
                )
                if result.returncode == 0:
                    git_info["branch"] = result.stdout.strip()
                
                # Get remote info
                result = subprocess.run(
                    ["git", "remote", "get-url", "origin"],
                    cwd=path, capture_output=True, text=True, timeout=5
                )
                if result.returncode == 0:
                    git_info["remote"] = result.stdout.strip()
                
                # Get status
                result = subprocess.run(
                    ["git", "status", "--porcelain"],
                    cwd=path, capture_output=True, text=True, timeout=5
                )
                if result.returncode == 0:
                    if result.stdout.strip():
                        git_info["status"] = "dirty"
                    else:
                        git_info["status"] = "clean"
                        
            except:
                pass
        
        return git_info
    
    def _analyze_team_patterns(self, path: Path) -> Dict[str, Any]:
        """Analyze team development patterns"""
        patterns = {"commit_frequency": 0, "contributors": [], "common_commands": []}
        
        if (path / ".git").exists():
            try:
                # Analyze commit patterns
                result = subprocess.run(
                    ["git", "log", "--oneline", "-30"],
                    cwd=path, capture_output=True, text=True, timeout=10
                )
                if result.returncode == 0:
                    commits = result.stdout.strip().split('\n')
                    patterns["commit_frequency"] = len(commits)
                
                # Get contributors
                result = subprocess.run(
                    ["git", "shortlog", "-sn"],
                    cwd=path, capture_output=True, text=True, timeout=10
                )
                if result.returncode == 0:
                    contributors = []
                    for line in result.stdout.strip().split('\n')[:5]:
                        if '\t' in line:
                            count, name = line.split('\t', 1)
                            contributors.append({"name": name, "commits": int(count)})
                    patterns["contributors"] = contributors
                    
            except:
                pass
        
        return patterns
    
    def _analyze_performance_metrics(self, path: Path) -> Dict[str, Any]:
        """Analyze project performance characteristics"""
        metrics = {
            "file_count": 0,
            "loc": 0,
            "complexity_estimate": "low",
            "build_time_estimate": 0
        }
        
        # Count files and estimate lines of code
        code_extensions = {'.py', '.js', '.ts', '.go', '.rs', '.java', '.cpp', '.c'}
        
        for file_path in path.rglob("*"):
            if file_path.is_file():
                metrics["file_count"] += 1
                
                if file_path.suffix in code_extensions:
                    try:
                        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                            lines = len(f.readlines())
                            metrics["loc"] += lines
                    except:
                        pass
        
        # Estimate complexity
        if metrics["loc"] > 10000:
            metrics["complexity_estimate"] = "high"
        elif metrics["loc"] > 1000:
            metrics["complexity_estimate"] = "medium"
        
        # Estimate build time based on project characteristics
        base_time = 5  # Base 5 seconds
        if metrics["file_count"] > 100:
            base_time += 10
        if metrics["loc"] > 5000:
            base_time += 15
            
        metrics["build_time_estimate"] = base_time
        
        return metrics
    
    def create_adaptive_workflow(self, goal: str, context: ProjectContext = None) -> str:
        """Create adaptive workflow based on goal and project context"""
        if not context:
            context = self.analyze_project()
        
        workflow_id = f"workflow_{int(time.time())}"
        
        # Generate workflow steps based on goal and context
        steps = self._generate_workflow_steps(goal, context)
        
        # Optimize step order based on performance history
        optimized_steps = self._optimize_step_order(steps, context)
        
        # Create workflow
        workflow = {
            "id": workflow_id,
            "goal": goal,
            "context": context,
            "steps": optimized_steps,
            "state": WorkflowState.PLANNING,
            "created": datetime.now().isoformat(),
            "progress": 0,
            "estimated_duration": sum(step.estimated_duration for step in optimized_steps)
        }
        
        self.active_workflows[workflow_id] = workflow
        return workflow_id
    
    def _generate_workflow_steps(self, goal: str, context: ProjectContext) -> List[WorkflowStep]:
        """Generate workflow steps based on goal and context"""
        steps = []
        
        goal_lower = goal.lower()
        
        if "setup" in goal_lower or "init" in goal_lower:
            # Project setup workflow
            if context.project_type.startswith("python"):
                steps.extend(self.workflow_templates["python_setup"])
            elif context.project_type.startswith("node"):
                steps.extend(self.workflow_templates["node_setup"])
        
        elif "deploy" in goal_lower:
            steps.extend(self.workflow_templates["deployment"])
        
        elif "update" in goal_lower or "sync" in goal_lower:
            steps.extend(self.workflow_templates["git_workflow"])
        
        elif "test" in goal_lower:
            # Add testing steps based on detected framework
            if context.test_framework == "pytest":
                steps.append(WorkflowStep("test", "Run pytest", "pytest", [], [], [], 15, 1))
            elif context.test_framework == "jest":
                steps.append(WorkflowStep("test", "Run jest", "npm test", [], [], [], 10, 1))
        
        elif "build" in goal_lower:
            # Add build steps based on build system
            if context.build_system == "npm":
                steps.append(WorkflowStep("build", "Build with npm", "npm run build", [], [], [], 20, 1))
            elif context.build_system == "make":
                steps.append(WorkflowStep("build", "Build with make", "make", [], [], [], 30, 1))
        
        # Add intelligent steps based on project state
        if context.git_info["status"] == "dirty":
            steps.insert(0, WorkflowStep("commit", "Commit changes", "git add . && git commit -m 'Auto commit'", [], [], [], 5, 0))
        
        if context.dependencies["outdated"]:
            steps.insert(-1, WorkflowStep("update_deps", "Update dependencies", "auto", [], [], [], 30, 2))
        
        return steps
    
    def _optimize_step_order(self, steps: List[WorkflowStep], context: ProjectContext) -> List[WorkflowStep]:
        """Optimize step execution order based on performance history"""
        # Sort by priority and dependencies
        optimized = sorted(steps, key=lambda s: (s.priority, len(s.dependencies)))
        
        # Apply performance-based optimizations
        for step in optimized:
            if step.command in self.performance_history:
                perf_data = self.performance_history[step.command]
                # Adjust estimated duration based on historical data
                avg_duration = perf_data.get("avg_duration", step.estimated_duration)
                step.estimated_duration = int(avg_duration * 1.1)  # Add 10% buffer
        
        return optimized
    
    async def execute_workflow(self, workflow_id: str) -> Dict[str, Any]:
        """Execute workflow with adaptive monitoring"""
        if workflow_id not in self.active_workflows:
            return {"error": "Workflow not found"}
        
        workflow = self.active_workflows[workflow_id]
        workflow["state"] = WorkflowState.EXECUTING
        
        results = {"steps": [], "success": True, "duration": 0}
        start_time = time.time()
        
        for step in workflow["steps"]:
            step_start = time.time()
            
            try:
                # Execute step
                step_result = await self._execute_step(step, workflow["context"])
                step_duration = time.time() - step_start
                
                # Record performance
                self._record_step_performance(step.command, step_duration, step_result["success"])
                
                results["steps"].append({
                    "step": step.name,
                    "success": step_result["success"],
                    "duration": step_duration,
                    "output": step_result.get("output", "")
                })
                
                if not step_result["success"]:
                    # Attempt recovery
                    recovery_result = await self._attempt_recovery(step, step_result)
                    if not recovery_result["success"]:
                        results["success"] = False
                        break
                
                workflow["progress"] = len(results["steps"]) / len(workflow["steps"]) * 100
                
            except Exception as e:
                results["steps"].append({
                    "step": step.name,
                    "success": False,
                    "error": str(e),
                    "duration": time.time() - step_start
                })
                results["success"] = False
                break
        
        results["duration"] = time.time() - start_time
        workflow["state"] = WorkflowState.COMPLETED if results["success"] else WorkflowState.FAILED
        
        # Save performance history
        self._save_performance_history()
        
        return results
    
    async def _execute_step(self, step: WorkflowStep, context: ProjectContext) -> Dict[str, Any]:
        """Execute individual workflow step"""
        command = step.command
        
        # Handle "auto" commands
        if command == "auto":
            command = self._resolve_auto_command(step.name, context)
        
        try:
            result = subprocess.run(
                command,
                shell=True,
                capture_output=True,
                text=True,
                timeout=step.estimated_duration + 30,
                cwd=context.root_path
            )
            
            return {
                "success": result.returncode == 0,
                "output": result.stdout,
                "error": result.stderr,
                "returncode": result.returncode
            }
            
        except subprocess.TimeoutExpired:
            return {
                "success": False,
                "error": f"Step timed out after {step.estimated_duration + 30} seconds"
            }
        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }
    
    def _resolve_auto_command(self, step_name: str, context: ProjectContext) -> str:
        """Resolve 'auto' command based on context"""
        step_lower = step_name.lower()
        
        if "install" in step_lower or "dependencies" in step_lower:
            if context.build_system == "npm":
                return "npm install"
            elif context.build_system == "poetry":
                return "poetry install"
            else:
                return "pip install -r requirements.txt"
        
        elif "test" in step_lower:
            if context.test_framework == "pytest":
                return "pytest"
            elif context.test_framework == "jest":
                return "npm test"
            else:
                return "python -m unittest"
        
        elif "build" in step_lower:
            if context.build_system == "npm":
                return "npm run build"
            elif context.build_system == "make":
                return "make build"
            else:
                return "python setup.py build"
        
        return "echo 'Auto command not resolved'"
    
    async def _attempt_recovery(self, step: WorkflowStep, failure_result: Dict) -> Dict[str, Any]:
        """Attempt to recover from step failure"""
        for recovery_cmd in step.failure_recovery:
            try:
                result = subprocess.run(
                    recovery_cmd,
                    shell=True,
                    capture_output=True,
                    text=True,
                    timeout=30
                )
                
                if result.returncode == 0:
                    # Retry original step
                    retry_result = await self._execute_step(step, None)
                    if retry_result["success"]:
                        return {"success": True, "recovered": True}
            except:
                continue
        
        return {"success": False, "recovered": False}
    
    def _record_step_performance(self, command: str, duration: float, success: bool):
        """Record step performance for future optimization"""
        if command not in self.performance_history:
            self.performance_history[command] = {
                "durations": [],
                "success_rate": 0,
                "total_runs": 0
            }
        
        perf = self.performance_history[command]
        perf["durations"].append(duration)
        perf["total_runs"] += 1
        
        if success:
            perf["success_rate"] = (perf["success_rate"] * (perf["total_runs"] - 1) + 1) / perf["total_runs"]
        else:
            perf["success_rate"] = (perf["success_rate"] * (perf["total_runs"] - 1)) / perf["total_runs"]
        
        # Keep only last 50 durations
        if len(perf["durations"]) > 50:
            perf["durations"] = perf["durations"][-50:]
        
        # Calculate average duration
        perf["avg_duration"] = sum(perf["durations"]) / len(perf["durations"])
    
    def get_project_insights(self, project_path: Path = None) -> Dict[str, Any]:
        """Get comprehensive project insights and recommendations"""
        context = self.analyze_project(project_path)
        
        insights = {
            "project_health": "good",
            "recommendations": [],
            "optimization_opportunities": [],
            "risk_factors": []
        }
        
        # Analyze project health
        health_score = 100
        
        if context.dependencies["outdated"]:
            health_score -= len(context.dependencies["outdated"]) * 5
            insights["recommendations"].append(f"Update {len(context.dependencies['outdated'])} outdated dependencies")
        
        if context.git_info["status"] == "dirty":
            health_score -= 10
            insights["recommendations"].append("Commit pending changes")
        
        if context.test_framework == "unknown":
            health_score -= 20
            insights["risk_factors"].append("No testing framework detected")
        
        # Set health level
        if health_score >= 80:
            insights["project_health"] = "excellent"
        elif health_score >= 60:
            insights["project_health"] = "good"
        elif health_score >= 40:
            insights["project_health"] = "fair"
        else:
            insights["project_health"] = "poor"
        
        # Optimization opportunities
        if context.performance_metrics["complexity_estimate"] == "high":
            insights["optimization_opportunities"].append("Consider code refactoring for better maintainability")
        
        if context.build_system == "unknown":
            insights["optimization_opportunities"].append("Add build automation")
        
        return insights

# Global orchestrator instance
orchestrator = AdaptiveOrchestrator()