# core/claude_memory.py
import json
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Any, Optional
import hashlib

from .config import config

class ClaudeMemory:
    """Claude memory and context management system"""
    
    def __init__(self, session_id: str = "default"):
        self.session_id = session_id
        self.memory_file = config.get_claude_memory_file(session_id)
        self.context_file = self.memory_file.parent / f"{session_id}_context.json"
        
        self.memory_file.parent.mkdir(parents=True, exist_ok=True)
        
        self.context = self._load_context()
        self._initialize_memory_file()
        
    def _load_context(self) -> Dict[str, Any]:
        """Load context from JSON file"""
        if self.context_file.exists():
            try:
                with open(self.context_file, 'r') as f:
                    return json.load(f)
            except Exception:
                pass
                
        return {
            "session_id": self.session_id,
            "created": datetime.now().isoformat(),
            "project_context": {},
            "conversation_summary": "",
            "key_topics": [],
            "current_task": "",
            "completed_tasks": [],
            "preferences": {},
            "codebase_knowledge": {},
            "last_updated": datetime.now().isoformat()
        }
    
    def _save_context(self):
        """Save context to JSON file"""
        self.context["last_updated"] = datetime.now().isoformat()
        try:
            with open(self.context_file, 'w') as f:
                json.dump(self.context, f, indent=2)
        except Exception as e:
            print(f"Error saving context: {e}")
    
    def _initialize_memory_file(self):
        """Initialize memory markdown file if it doesn't exist"""
        if not self.memory_file.exists():
            initial_content = f"""# Claude Memory Log - {self.session_id}

Session started: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

## Project Context
- Working Directory: {Path.cwd()}
- Session ID: {self.session_id}

## Conversation History

"""
            with open(self.memory_file, 'w') as f:
                f.write(initial_content)
    
    def log_interaction(self, user_input: str, claude_response: str, command_result: Optional[str] = None):
        """Log user-Claude interaction"""
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        
        entry = f"""
### {timestamp}

**User:** {user_input}

**Claude:** {claude_response}
"""
        
        if command_result:
            entry += f"""
**Command Result:**
```
{command_result}
```
"""
        
        entry += "\n---\n"
        
        # Append to memory file
        with open(self.memory_file, 'a') as f:
            f.write(entry)
            
        # Update context
        self._update_context_from_interaction(user_input, claude_response)
    
    def _update_context_from_interaction(self, user_input: str, claude_response: str):
        """Update context based on interaction"""
        # Extract key topics
        words = user_input.lower().split()
        important_words = [w for w in words if len(w) > 4 and w not in ['that', 'this', 'with', 'from', 'they', 'have', 'will', 'been']]
        
        for word in important_words[:3]:  # Add up to 3 key words
            if word not in self.context["key_topics"]:
                self.context["key_topics"].append(word)
                
        # Keep only recent key topics
        self.context["key_topics"] = self.context["key_topics"][-20:]
        
        # Update conversation summary (keep it brief)
        if len(self.context["conversation_summary"].split()) > 100:
            # Truncate if too long
            words = self.context["conversation_summary"].split()
            self.context["conversation_summary"] = " ".join(words[-50:])
            
        # Add current interaction summary
        interaction_summary = f" User asked about {' '.join(important_words[:2])}."
        self.context["conversation_summary"] += interaction_summary
        
        self._save_context()
    
    def set_project_context(self, context_info: Dict[str, Any]):
        """Set project-specific context"""
        self.context["project_context"].update(context_info)
        self._save_context()
    
    def add_codebase_knowledge(self, file_path: str, analysis: str):
        """Add knowledge about codebase"""
        file_hash = hashlib.md5(file_path.encode()).hexdigest()[:8]
        self.context["codebase_knowledge"][file_path] = {
            "analysis": analysis,
            "hash": file_hash,
            "timestamp": datetime.now().isoformat()
        }
        self._save_context()
    
    def set_current_task(self, task: str):
        """Set current working task"""
        if self.context["current_task"]:
            # Move current task to completed
            self.context["completed_tasks"].append({
                "task": self.context["current_task"],
                "completed": datetime.now().isoformat()
            })
            
        self.context["current_task"] = task
        self._save_context()
    
    def complete_current_task(self, result: str = ""):
        """Mark current task as completed"""
        if self.context["current_task"]:
            self.context["completed_tasks"].append({
                "task": self.context["current_task"],
                "result": result,
                "completed": datetime.now().isoformat()
            })
            self.context["current_task"] = ""
            self._save_context()
    
    def get_context_summary(self) -> str:
        """Get formatted context summary for Claude"""
        summary = f"""# Current Session Context

**Session:** {self.session_id}
**Current Directory:** {Path.cwd()}
**Current Task:** {self.context["current_task"] or "None"}

**Key Topics:** {", ".join(self.context["key_topics"][-10:])}

**Recent Conversation:** {self.context["conversation_summary"][-200:]}

**Project Context:**
"""
        
        for key, value in self.context["project_context"].items():
            summary += f"- {key}: {value}\n"
            
        if self.context["completed_tasks"]:
            summary += "\n**Recently Completed Tasks:**\n"
            for task in self.context["completed_tasks"][-3:]:
                summary += f"- {task['task']} (completed {task['completed']})\n"
                
        return summary
    
    def get_codebase_context(self) -> str:
        """Get codebase knowledge context"""
        if not self.context["codebase_knowledge"]:
            return "No codebase analysis available."
            
        context = "**Codebase Knowledge:**\n"
        for file_path, info in list(self.context["codebase_knowledge"].items())[-5:]:
            context += f"- {file_path}: {info['analysis'][:100]}...\n"
            
        return context
    
    def search_memory(self, query: str) -> List[str]:
        """Search memory for relevant content"""
        if not self.memory_file.exists():
            return []
            
        try:
            with open(self.memory_file, 'r') as f:
                content = f.read()
                
            # Simple search - split into interactions
            interactions = content.split('---')
            relevant = []
            
            query_words = query.lower().split()
            for interaction in interactions:
                interaction_lower = interaction.lower()
                if any(word in interaction_lower for word in query_words):
                    relevant.append(interaction.strip())
                    
            return relevant[-5:]  # Return last 5 relevant interactions
            
        except Exception:
            return []
    
    def export_memory(self) -> str:
        """Export complete memory as string"""
        try:
            with open(self.memory_file, 'r') as f:
                return f.read()
        except Exception:
            return "Memory file not accessible"
    
    def cleanup_old_memory(self, days: int = 30):
        """Clean up old memory entries"""
        # For now, just truncate if file gets too large
        try:
            if self.memory_file.exists() and self.memory_file.stat().st_size > 1024 * 1024:  # 1MB
                with open(self.memory_file, 'r') as f:
                    lines = f.readlines()
                    
                # Keep last half of the file
                with open(self.memory_file, 'w') as f:
                    f.writelines(lines[len(lines)//2:])
                    
        except Exception:
            pass

# Global memory instance
memory = ClaudeMemory()