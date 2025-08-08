package shell

import (
	"fmt"
	"regexp"
	"strings"
)

// AgentRouter provides intelligent routing for unknown commands to appropriate agents
type AgentRouter struct {
	// Command patterns that suggest specific agents
	patterns map[string][]CommandPattern
}

// CommandPattern represents a pattern that matches to an agent suggestion
type CommandPattern struct {
	Pattern     *regexp.Regexp
	Agent       string
	Confidence  float64
	Description string
}

// AgentSuggestion represents a suggested agent command for an unknown command
type AgentSuggestion struct {
	Agent       string
	Command     string
	Confidence  float64
	Reasoning   string
}

// NewAgentRouter creates a new intelligent agent router
func NewAgentRouter() *AgentRouter {
	router := &AgentRouter{
		patterns: make(map[string][]CommandPattern),
	}
	
	router.initializePatterns()
	return router
}

// initializePatterns sets up the pattern matching rules for agent routing
func (ar *AgentRouter) initializePatterns() {
	// Planning-related patterns
	planningPatterns := []CommandPattern{
		{
			Pattern:     regexp.MustCompile(`(?i)^(plan|create|build|make|develop|implement|setup|configure).*`),
			Agent:       "plan",
			Confidence:  0.8,
			Description: "Command suggests project planning or task breakdown",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(roadmap|strategy|steps|tasks|workflow|process).*`),
			Agent:       "plan",
			Confidence:  0.7,
			Description: "Keywords suggest planning or task organization",
		},
		{
			Pattern:     regexp.MustCompile(`(?i)^(how to|explain how|steps to).*`),
			Agent:       "plan", 
			Confidence:  0.6,
			Description: "Question format suggests need for step-by-step planning",
		},
	}
	
	// Design-related patterns
	designPatterns := []CommandPattern{
		{
			Pattern:     regexp.MustCompile(`(?i)^(design|generate|create).*(ui|interface|component|page|form|button).*`),
			Agent:       "design",
			Confidence:  0.9,
			Description: "Command explicitly mentions UI design or component creation",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(react|vue|angular|html|css|bootstrap|tailwind).*`),
			Agent:       "design",
			Confidence:  0.7,
			Description: "Web framework or styling keywords suggest UI design task",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(login|signup|dashboard|card|modal|navbar|sidebar).*`),
			Agent:       "design",
			Confidence:  0.8,
			Description: "Common UI component names suggest design task",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(dark|light|theme|responsive|mobile|layout).*`),
			Agent:       "design",
			Confidence:  0.6,
			Description: "Design-related styling keywords",
		},
	}
	
	// Code review patterns
	reviewPatterns := []CommandPattern{
		{
			Pattern:     regexp.MustCompile(`(?i)^(review|check|analyze|examine|audit|inspect).*\.(py|js|go|java|cpp|c|rb|php|ts|tsx|jsx)$`),
			Agent:       "review",
			Confidence:  0.9,
			Description: "Command explicitly mentions code review with file extension",
		},
		{
			Pattern:     regexp.MustCompile(`(?i)^(review|check|analyze|examine|audit|inspect).*`),
			Agent:       "review",
			Confidence:  0.7,
			Description: "General review/analysis command that could apply to code",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(security|vulnerability|bug|error|lint|quality).*`),
			Agent:       "review",
			Confidence:  0.6,
			Description: "Keywords suggest code quality or security review",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(optimize|performance|refactor|improve).*`),
			Agent:       "review",
			Confidence:  0.5,
			Description: "Code improvement keywords suggest review first",
		},
	}
	
	// Code generation patterns
	codePatterns := []CommandPattern{
		{
			Pattern:     regexp.MustCompile(`(?i)^(code|generate|write|create).*(function|class|method|script).*`),
			Agent:       "code",
			Confidence:  0.9,
			Description: "Explicit code generation request",
		},
		{
			Pattern:     regexp.MustCompile(`(?i)^(implement|write|create).*(algorithm|function|class|api|endpoint).*`),
			Agent:       "code",
			Confidence:  0.8,
			Description: "Implementation request suggests code generation",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(refactor|optimize|fix|debug|modify).*\.(py|js|go|java|cpp|c|rb|php|ts|tsx|jsx)$`),
			Agent:       "code",
			Confidence:  0.7,
			Description: "Code modification task with file extension",
		},
		{
			Pattern:     regexp.MustCompile(`(?i).*(python|javascript|golang|go|java|c\+\+|ruby|php|typescript).*`),
			Agent:       "code",
			Confidence:  0.6,
			Description: "Programming language mentioned suggests code task",
		},
	}
	
	ar.patterns["plan"] = planningPatterns
	ar.patterns["design"] = designPatterns
	ar.patterns["review"] = reviewPatterns
	ar.patterns["code"] = codePatterns
}

// SuggestAgent analyzes an unknown command and suggests the most appropriate agent
func (ar *AgentRouter) SuggestAgent(command string, args []string) *AgentSuggestion {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}
	
	bestSuggestion := &AgentSuggestion{
		Confidence: 0.0,
	}
	
	// Test against all patterns
	for agent, patterns := range ar.patterns {
		for _, pattern := range patterns {
			if pattern.Pattern.MatchString(fullCommand) {
				if pattern.Confidence > bestSuggestion.Confidence {
					bestSuggestion = &AgentSuggestion{
						Agent:      agent,
						Command:    ar.buildAgentCommand(agent, fullCommand),
						Confidence: pattern.Confidence,
						Reasoning:  pattern.Description,
					}
				}
			}
		}
	}
	
	// Apply additional heuristics
	ar.applyContextualHeuristics(fullCommand, bestSuggestion)
	
	// Only suggest if confidence is above threshold
	if bestSuggestion.Confidence < 0.4 {
		return nil
	}
	
	return bestSuggestion
}

// applyContextualHeuristics applies additional context-based analysis
func (ar *AgentRouter) applyContextualHeuristics(fullCommand string, suggestion *AgentSuggestion) {
	command := strings.ToLower(fullCommand)
	
	// Boost confidence for explicit requests
	if strings.Contains(command, "help me") || strings.Contains(command, "how do i") {
		suggestion.Confidence += 0.1
	}
	
	// Boost confidence for file extensions
	if regexp.MustCompile(`\.(py|js|go|java|cpp|c|rb|php|ts|tsx|jsx|html|css)$`).MatchString(command) {
		if suggestion.Agent == "review" || suggestion.Agent == "code" {
			suggestion.Confidence += 0.2
		}
	}
	
	// Boost confidence for project-related keywords
	projectKeywords := []string{"project", "application", "app", "system", "platform", "service"}
	for _, keyword := range projectKeywords {
		if strings.Contains(command, keyword) && suggestion.Agent == "plan" {
			suggestion.Confidence += 0.1
		}
	}
	
	// Reduce confidence for very short commands
	if len(fullCommand) < 10 {
		suggestion.Confidence *= 0.8
	}
	
	// Increase confidence for longer, more descriptive commands
	if len(fullCommand) > 50 {
		suggestion.Confidence *= 1.1
	}
	
	// Cap confidence at 1.0
	if suggestion.Confidence > 1.0 {
		suggestion.Confidence = 1.0
	}
}

// buildAgentCommand constructs the appropriate agent command from the user input
func (ar *AgentRouter) buildAgentCommand(agent, fullCommand string) string {
	switch agent {
	case "plan":
		// Extract the main goal/objective
		goal := ar.extractPlanningGoal(fullCommand)
		return fmt.Sprintf("plan \"%s\"", goal)
		
	case "design":
		// Extract the UI component description
		description := ar.extractDesignDescription(fullCommand)
		return fmt.Sprintf("design \"%s\"", description)
		
	case "review":
		// Extract filename if present, otherwise use the command as-is
		filename := ar.extractFilename(fullCommand)
		if filename != "" {
			return fmt.Sprintf("review \"%s\"", filename)
		}
		return fmt.Sprintf("review \"%s\"", fullCommand)
		
	case "code":
		// Extract the code generation request
		request := ar.extractCodeRequest(fullCommand)
		return fmt.Sprintf("code generate \"%s\"", request)
		
	default:
		return fullCommand
	}
}

// extractPlanningGoal extracts the main goal from a planning-related command
func (ar *AgentRouter) extractPlanningGoal(command string) string {
	// Remove common prefixes
	prefixes := []string{"plan", "create", "build", "make", "develop", "implement", "setup", "configure", "how to", "steps to"}
	
	goal := strings.ToLower(command)
	for _, prefix := range prefixes {
		if strings.HasPrefix(goal, prefix) {
			goal = strings.TrimSpace(goal[len(prefix):])
			break
		}
	}
	
	// Clean up the goal
	goal = strings.Trim(goal, ".,!?")
	if goal == "" {
		return command
	}
	
	return goal
}

// extractDesignDescription extracts the UI component description
func (ar *AgentRouter) extractDesignDescription(command string) string {
	// Remove common prefixes
	prefixes := []string{"design", "generate", "create"}
	
	description := strings.ToLower(command)
	for _, prefix := range prefixes {
		if strings.HasPrefix(description, prefix) {
			description = strings.TrimSpace(description[len(prefix):])
			break
		}
	}
	
	// Clean up the description
	description = strings.Trim(description, ".,!?")
	if description == "" {
		return command
	}
	
	return description
}

// extractFilename attempts to extract a filename from the command
func (ar *AgentRouter) extractFilename(command string) string {
	// Look for file extensions
	filePattern := regexp.MustCompile(`(\S+\.(py|js|go|java|cpp|c|rb|php|ts|tsx|jsx|html|css|json|xml|yaml|yml))`)
	matches := filePattern.FindStringSubmatch(command)
	
	if len(matches) > 1 {
		return matches[1]
	}
	
	return ""
}

// extractCodeRequest extracts the code generation request
func (ar *AgentRouter) extractCodeRequest(command string) string {
	// Remove common prefixes
	prefixes := []string{"code", "generate", "write", "create", "implement"}
	
	request := strings.ToLower(command)
	for _, prefix := range prefixes {
		if strings.HasPrefix(request, prefix) {
			request = strings.TrimSpace(request[len(prefix):])
			break
		}
	}
	
	// Clean up the request
	request = strings.Trim(request, ".,!?")
	if request == "" {
		return command
	}
	
	return request
}

// GetAgentHelp returns help text for using agents
func (ar *AgentRouter) GetAgentHelp() string {
	return `
🤖 AI AGENT COMMANDS

Available agents:
  plan    - Break down goals into actionable tasks
  design  - Generate UI components from descriptions  
  review  - Analyze code files for improvements
  code    - Generate, refactor, or explain code

Examples:
  plan "build a REST API"
  design "dark login form with validation"
  review myfile.py
  code generate "function that sorts a list"

💡 SMART ROUTING

Lucien can automatically suggest the right agent:
  "create a web app" → plan "create a web app"
  "generate login form" → design "generate login form" 
  "check main.py" → review "main.py"
  "implement sorting algorithm" → code generate "implement sorting algorithm"

Type any natural language command and Lucien will suggest the appropriate agent!
`
}

// FormatSuggestion formats an agent suggestion for display
func (ar *AgentRouter) FormatSuggestion(suggestion *AgentSuggestion) string {
	confidence := int(suggestion.Confidence * 100)
	
	return fmt.Sprintf(`
💡 Did you mean: %s
   Confidence: %d%%
   Reasoning: %s
   
   Run this command? (y/n): `, 
		suggestion.Command, confidence, suggestion.Reasoning)
}