# 🤖 LUCIEN AGENT COMMANDS

Lucien Shell now supports Python-powered AI agents that extend the shell with intelligent capabilities.

## ✅ AVAILABLE AGENT COMMANDS

| Command | Purpose | Usage |
|---------|---------|--------|
| `plan` | AI task planning and breakdown | `plan "build a web app"` |
| `design` | UI code generation | `design "dark login page with neon glow"` |
| `review` | Code analysis and suggestions | `review myfile.py` |
| `code` | Code generation and refactoring | `code generate "sort function"` |

## 🚀 USAGE EXAMPLES

### Planning Projects
```bash
lucien> plan "create a REST API for user management"
>> GOAL: create a REST API for user management
>> PLAN (16 tasks):
   📋 Define API specification and endpoints
   🏗️ Setup development environment and framework
   📊 Design database schema and relationships
   🔐 Implement authentication and authorization
   ...
```

### Generating UI Code
```bash
lucien> design "responsive card component with dark theme"
🎨 Generating UI component: responsive card component with dark theme
📋 Framework: React
✅ UI component generated successfully!
📁 Saved to: snapmethod/exports/CardComponent.tsx
🚀 Component ready for use in your react project
```

### Reviewing Code
```bash
lucien> review src/main.py
📋 Code Review: main.py
==================================================
Summary: 3 issues found
⚠️  1 warning
💡 2 suggestions

## ⚠️ WARNINGS
⚠️ Line 45: Avoid using eval() or exec() - security risk
   Suggestion: Use safer alternatives like literal_eval()
...
```

### Generating Code
```bash
lucien> code generate "function that validates email addresses"
🤖 Generating code: function that validates email addresses
📄 GENERATED CODE:
def validate_email(email):
    """
    Validate email addresses using regex pattern
    """
    import re
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None
```

## 🔧 TECHNICAL DETAILS

### Architecture
- **Go Shell**: Core command routing and execution
- **Python Agents**: AI-powered plugins for specialized tasks  
- **Plugin Bridge**: Secure subprocess execution with timeout
- **Agent Framework**: Modular agent system with standardized interfaces

### Agent Locations
- **Scripts**: `~/.lucien/plugins/*.py` (CLI entry points)
- **Agents**: `~/.lucien/plugins/agents/*.py` (Core agent classes)
- **Dependencies**: Minimal - uses Python standard library

### Security Features
- **Sandboxed Execution**: Agents run in separate Python processes
- **Timeout Protection**: 30-second execution limit per command
- **Input Validation**: Command arguments are sanitized
- **Safe Mode Integration**: Agent commands respect shell security policies

## 🛠️ EXTENDING WITH CUSTOM AGENTS

Create new agents by:

1. **Add Agent Class**: Create `~/.lucien/plugins/agents/my_agent.py`
2. **Add CLI Script**: Create `~/.lucien/plugins/my_agent.py` 
3. **Register Command**: Update `internal/plugin/bridge.go`
4. **Rebuild Shell**: `go build cmd/lucien/main.go`

### Example Custom Agent
```python
# ~/.lucien/plugins/agents/my_agent.py
class MyAgent:
    def run(self, prompt):
        return f"Processed: {prompt}"

# ~/.lucien/plugins/my_agent.py  
#!/usr/bin/env python3
import sys
from agents.my_agent import MyAgent

agent = MyAgent()
result = agent.run(" ".join(sys.argv[1:]))
print(result)
```

## 📊 PERFORMANCE

- **Startup Time**: < 1 second per agent command
- **Memory Usage**: ~20MB per Python subprocess  
- **Concurrency**: Agents run independently, no conflicts
- **Caching**: Agent classes cached between invocations

## 🚨 TROUBLESHOOTING

**Python Not Found**: Ensure Python 3.7+ is in PATH
**Import Errors**: Check agent files exist in `~/.lucien/plugins/`
**Permission Denied**: Verify plugins directory is readable
**Timeout Errors**: Increase timeout in `bridge.go` if needed

The agent system provides powerful AI capabilities while maintaining the security and performance of the core Go shell.