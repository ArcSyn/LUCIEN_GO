"""
PlannerAgent - AI-powered task planning and breakdown
"""

import re
import json
from typing import List, Dict, Any
from dataclasses import dataclass
import os

@dataclass
class Task:
    id: int
    title: str
    description: str
    priority: str  # high, medium, low
    estimated_time: str
    dependencies: List[int]
    category: str

class PlannerAgent:
    """AI-powered planning agent that breaks down goals into actionable tasks"""
    
    def __init__(self):
        self.task_patterns = {
            'setup': ['install', 'configure', 'initialize', 'setup', 'create environment'],
            'design': ['design', 'mockup', 'wireframe', 'prototype', 'UI', 'UX'],
            'development': ['code', 'implement', 'build', 'develop', 'write'],
            'testing': ['test', 'validate', 'verify', 'debug', 'QA'],
            'deployment': ['deploy', 'release', 'publish', 'ship', 'launch'],
            'documentation': ['document', 'readme', 'guide', 'manual', 'docs'],
            'maintenance': ['maintain', 'update', 'patch', 'fix', 'optimize']
        }
        
    def run(self, goal: str) -> List[str]:
        """
        Break down a high-level goal into actionable tasks
        """
        # Normalize goal
        goal = goal.strip().lower()
        
        # Determine project type and generate specialized plan
        if self._is_web_project(goal):
            return self._plan_web_project(goal)
        elif self._is_api_project(goal):
            return self._plan_api_project(goal)
        elif self._is_game_project(goal):
            return self._plan_game_project(goal)
        elif self._is_cli_project(goal):
            return self._plan_cli_project(goal)
        elif self._is_data_project(goal):
            return self._plan_data_project(goal)
        else:
            return self._plan_generic_project(goal)
    
    def _is_web_project(self, goal: str) -> bool:
        web_keywords = ['website', 'web app', 'frontend', 'backend', 'full stack', 'react', 'vue', 'angular']
        return any(keyword in goal for keyword in web_keywords)
    
    def _is_api_project(self, goal: str) -> bool:
        api_keywords = ['api', 'rest', 'graphql', 'server', 'endpoint', 'microservice']
        return any(keyword in goal for keyword in api_keywords)
    
    def _is_game_project(self, goal: str) -> bool:
        game_keywords = ['game', 'gaming', 'unity', 'unreal', '2d', '3d', 'player', 'level']
        return any(keyword in goal for keyword in game_keywords)
    
    def _is_cli_project(self, goal: str) -> bool:
        cli_keywords = ['cli', 'command line', 'terminal', 'shell', 'script']
        return any(keyword in goal for keyword in cli_keywords)
    
    def _is_data_project(self, goal: str) -> bool:
        data_keywords = ['data', 'analytics', 'machine learning', 'ml', 'ai', 'analysis', 'visualization']
        return any(keyword in goal for keyword in data_keywords)
    
    def _plan_web_project(self, goal: str) -> List[str]:
        """Plan for web development projects"""
        tasks = [
            "📋 Define project requirements and scope",
            "🎨 Create wireframes and UI mockups", 
            "🏗️ Setup development environment (Node.js, package manager)",
            "⚙️ Initialize project structure and configuration",
            "🎯 Setup build tools and development workflow",
            "💾 Design database schema and data models",
            "🔐 Implement authentication and user management",
            "🌐 Build core frontend components and routing",
            "🔗 Develop backend API endpoints", 
            "🧪 Write comprehensive tests (unit, integration, e2e)",
            "🎨 Apply styling and responsive design",
            "⚡ Optimize performance and loading times",
            "🔒 Implement security measures and validation",
            "📱 Test cross-browser compatibility",
            "🚀 Setup CI/CD pipeline for deployment",
            "📊 Add analytics and monitoring",
            "📚 Write user documentation and guides"
        ]
        return tasks
    
    def _plan_api_project(self, goal: str) -> List[str]:
        """Plan for API development projects"""
        tasks = [
            "📋 Define API specification and endpoints",
            "🏗️ Setup development environment and framework",
            "📊 Design database schema and relationships",
            "🔐 Implement authentication and authorization",
            "🌐 Build core API endpoints with CRUD operations",
            "✅ Add input validation and error handling", 
            "📝 Implement request/response serialization",
            "🧪 Write comprehensive API tests",
            "📊 Add logging and monitoring",
            "⚡ Implement caching and optimization",
            "🔒 Add security headers and CORS configuration",
            "📈 Setup rate limiting and throttling",
            "🔧 Create database migrations and seeders",
            "📚 Generate API documentation (Swagger/OpenAPI)",
            "🚀 Setup deployment and environment configuration",
            "🔍 Add health checks and status endpoints"
        ]
        return tasks
        
    def _plan_game_project(self, goal: str) -> List[str]:
        """Plan for game development projects"""
        tasks = [
            "🎮 Define game concept and core mechanics",
            "🎨 Create art style guide and asset pipeline",
            "🏗️ Setup game development environment",
            "⚙️ Initialize project structure and version control",
            "🎯 Implement core game loop and state management",
            "👤 Create player character and controls",
            "🌍 Build level/world generation system",
            "🎵 Add audio system and sound effects",
            "💥 Implement game physics and collision detection",
            "🎪 Create game UI and menu systems",
            "💾 Add save/load game functionality", 
            "🧪 Playtesting and balancing",
            "🎨 Polish graphics and animations",
            "🔧 Optimize performance for target platforms",
            "📱 Test on different devices/platforms",
            "🚀 Build and package for distribution",
            "📚 Create player documentation and tutorials"
        ]
        return tasks
    
    def _plan_cli_project(self, goal: str) -> List[str]:
        """Plan for CLI tool development"""
        tasks = [
            "📋 Define CLI interface and command structure",
            "🏗️ Setup development environment and dependencies",
            "⚙️ Initialize project with CLI framework",
            "🔧 Implement core command parsing and routing",
            "📝 Add configuration file support",
            "💾 Implement data storage and persistence",
            "✅ Add input validation and error handling",
            "📊 Implement logging and verbose modes",
            "🧪 Write comprehensive unit tests",
            "📚 Create help system and documentation",
            "🎨 Add colored output and progress indicators",
            "🔌 Support plugins or extensions",
            "⚡ Optimize performance for large inputs",
            "📦 Setup packaging and distribution",
            "🚀 Create installation scripts",
            "🔧 Add auto-completion support"
        ]
        return tasks
    
    def _plan_data_project(self, goal: str) -> List[str]:
        """Plan for data analysis/ML projects"""
        tasks = [
            "📋 Define project objectives and success metrics",
            "📊 Collect and explore available datasets",
            "🧹 Clean and preprocess data",
            "🔍 Perform exploratory data analysis (EDA)",
            "🎯 Feature engineering and selection", 
            "🤖 Choose and implement ML algorithms/models",
            "🧪 Split data and setup validation strategy",
            "⚙️ Train and tune model hyperparameters",
            "📈 Evaluate model performance and metrics",
            "🔧 Implement model versioning and tracking",
            "🚀 Deploy model to production environment",
            "📊 Create monitoring and alerting systems",
            "📚 Document methodology and findings",
            "🎨 Build visualization dashboards",
            "🔄 Setup automated retraining pipelines",
            "✅ Validate results with domain experts"
        ]
        return tasks
    
    def _plan_generic_project(self, goal: str) -> List[str]:
        """Fallback planning for generic projects"""
        tasks = [
            "📋 Define project scope and requirements",
            "🎯 Set clear objectives and success criteria",
            "🏗️ Setup development environment",
            "⚙️ Initialize project structure",
            "🔧 Implement core functionality",
            "🧪 Add comprehensive testing",
            "📚 Write documentation",
            "⚡ Optimize and refactor code",
            "🔒 Implement security measures",
            "🚀 Deploy to production environment",
            "📊 Monitor and maintain solution"
        ]
        
        # Try to add domain-specific tasks based on keywords
        if 'database' in goal or 'data' in goal:
            tasks.insert(3, "💾 Design database schema")
            tasks.insert(4, "🔗 Setup database connections")
        
        if 'user' in goal or 'auth' in goal:
            tasks.insert(5, "🔐 Implement user authentication")
        
        if 'api' in goal:
            tasks.insert(6, "🌐 Create API endpoints")
        
        return tasks

    def get_task_details(self, task_title: str) -> Dict[str, Any]:
        """Get detailed information about a specific task"""
        # This could be expanded to provide more detailed task breakdowns
        return {
            'title': task_title,
            'estimated_time': '2-4 hours',
            'priority': 'medium',
            'resources': [],
            'tips': []
        }