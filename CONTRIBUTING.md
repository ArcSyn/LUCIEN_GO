# Contributing to Lucien CLI

Thank you for your interest in contributing to Lucien CLI! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork locally
3. **Create** a feature branch
4. **Make** your changes
5. **Test** thoroughly
6. **Submit** a pull request

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- **Go 1.21+**: [Download here](https://golang.org/dl/)
- **Git**: [Download here](https://git-scm.com/)
- **Make** (optional but recommended): For using Makefile commands

### Setting Up Development Environment

```bash
# Fork the repository on GitHub first, then:

# Clone your fork
git clone https://github.com/YOUR_USERNAME/LUCIEN_GO.git
cd LUCIEN_GO

# Add upstream remote
git remote add upstream https://github.com/ArcSyn/LUCIEN_GO.git

# Install dependencies
go mod tidy

# Build the project
go build -o lucien ./cmd/lucien

# Run tests
go test ./...
```

### Project Structure

```
LUCIEN_GO/
├── cmd/lucien/           # Main application entry point
├── internal/             # Internal packages
│   ├── shell/           # Core shell functionality
│   ├── ui/              # User interface components
│   ├── completion/      # Tab completion
│   ├── history/         # Command history
│   ├── jobs/            # Job control
│   ├── plugin/          # Plugin system
│   ├── policy/          # Security policies
│   └── sandbox/         # Security sandbox
├── plugins/             # Plugin implementations
├── policies/            # OPA security policies
├── scripts/             # Build and utility scripts
├── tests/               # Integration tests
└── docs/                # Documentation
```

## Development Workflow

### 1. Create Feature Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write clean, readable code
- Follow existing patterns and conventions
- Add tests for new functionality
- Update documentation as needed

### 3. Test Changes

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Test specific packages
go test ./internal/shell/
go test ./tests/

# Manual testing
go build -o lucien ./cmd/lucien
./lucien --batch < test_commands.txt
```

### 4. Commit Changes

Follow conventional commit format:

```bash
git add .
git commit -m "feat: add new security validation feature"
```

## Code Style

### Go Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format code
- Follow Go naming conventions
- Write clear, self-documenting code
- Add comments for exported functions and complex logic

### Code Formatting

```bash
# Format all Go files
go fmt ./...

# Run linter (if available)
golangci-lint run
```

### Naming Conventions

```go
// Good
func validateCommand(cmd string) error { }
type SecurityGuard struct { }
const MaxHistorySize = 1000

// Avoid
func Validate_Command(cmd string) error { }
type securityguard struct { }
const MAX_HISTORY_SIZE = 1000
```

### Error Handling

```go
// Good - descriptive errors
if err := validateInput(input); err != nil {
    return fmt.Errorf("failed to validate input %q: %w", input, err)
}

// Good - context-aware errors
func (s *Shell) executeCommand(cmd string) error {
    if cmd == "" {
        return errors.New("command cannot be empty")
    }
    // ...
}
```

## Testing

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete user workflows

### Writing Tests

```go
func TestParseCommand(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected CommandType
        wantErr  bool
    }{
        {
            name:     "simple command",
            input:    "echo hello",
            expected: CommandSimple,
            wantErr:  false,
        },
        {
            name:     "and operator",
            input:    "echo first && echo second",
            expected: CommandAnd,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parseCommand(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseCommand() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result.Type != tt.expected {
                t.Errorf("parseCommand() = %v, expected %v", result.Type, tt.expected)
            }
        })
    }
}
```

### Test Requirements

- All new features must have tests
- Maintain or improve code coverage
- Tests should be fast and reliable
- Use table-driven tests for multiple scenarios

## Submitting Changes

### Pull Request Process

1. **Update Documentation**: Ensure README and docs reflect changes
2. **Add Tests**: Include comprehensive tests for new features
3. **Check CI**: Ensure all automated checks pass
4. **Write Description**: Clearly explain your changes

### Pull Request Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] New tests added for new functionality
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Code style changes (formatting, missing semi-colons, etc.)
- `refactor`: Code changes that neither fix a bug nor add a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to build process or auxiliary tools

**Examples:**
```
feat(security): add strict mode validation
fix(parser): handle quoted operators correctly
docs(readme): update installation instructions
test(shell): add operator chaining tests
```

## Reporting Issues

### Before Reporting

1. **Search existing issues** to avoid duplicates
2. **Check documentation** for known limitations
3. **Test with latest version** to ensure issue still exists

### Issue Templates

Use the appropriate issue template:

- **Bug Report**: For functional problems
- **Feature Request**: For new functionality
- **Documentation**: For documentation improvements

### Good Bug Reports Include

- **Environment**: OS, Go version, Lucien version
- **Steps to Reproduce**: Exact commands and inputs
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Additional Context**: Logs, screenshots, etc.

### Example Bug Report

```markdown
**Environment:**
- OS: Ubuntu 22.04
- Go: 1.21.5
- Lucien: v1.0.0

**Steps to Reproduce:**
1. Run `lucien --batch`
2. Enter command: `echo "test" && pwd`
3. Observe output

**Expected Behavior:**
Should show "test" followed by current directory

**Actual Behavior:**
Only shows "test", pwd doesn't execute

**Additional Context:**
Works fine in interactive mode, only fails in batch mode.
```

## Feature Requests

### Before Requesting

- Check if feature aligns with project goals
- Consider if it can be implemented as a plugin
- Search for existing requests

### Good Feature Requests

- **Clear Use Case**: Why is this needed?
- **Proposed Solution**: How should it work?
- **Alternatives**: What other approaches were considered?
- **Implementation**: Any implementation details or suggestions

## Development Tips

### Local Testing

```bash
# Quick build and test
go build -o lucien ./cmd/lucien && echo "pwd && echo test" | ./lucien --batch

# Test specific functionality
go test -run TestParseCommand ./internal/shell/

# Benchmark performance
go test -bench=. ./internal/shell/
```

### Debugging

```bash
# Enable debug output
LUCIEN_DEBUG=1 ./lucien

# Use Go debugger
dlv debug ./cmd/lucien
```

### Performance Considerations

- Profile memory usage with `go tool pprof`
- Benchmark critical paths
- Consider memory allocations in hot paths
- Test with large inputs and long-running sessions

## Security Considerations

- **Input Validation**: Always validate user input
- **Command Injection**: Be extremely careful with command execution
- **File Access**: Respect sandbox boundaries
- **Secrets**: Never log or expose sensitive information

### Security Review Process

All security-related changes undergo additional review:

1. **Security Impact Assessment**: Document potential security implications
2. **Threat Modeling**: Consider attack vectors
3. **Penetration Testing**: Test against common attacks
4. **Code Review**: Multiple reviewers for security changes

## Getting Help

- **Documentation**: Check [User Manual](docs/USER_MANUAL.md)
- **Discussions**: Use [GitHub Discussions](https://github.com/ArcSyn/LUCIEN_GO/discussions)
- **Issues**: Create issue for bugs or questions
- **Email**: Contact maintainers at hello@arcsyn.dev

## Recognition

Contributors will be recognized in:
- Project README
- Release notes
- Annual contributor report

Thank you for contributing to Lucien CLI! 🚀