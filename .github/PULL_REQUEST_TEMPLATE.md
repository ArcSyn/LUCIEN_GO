## Description

Brief description of the changes in this PR.

Fixes #(issue number)

## Type of Change

Please delete options that are not relevant:

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test coverage improvement

## Changes Made

- [ ] List specific changes
- [ ] Use bullet points
- [ ] Be clear and concise

## Testing

### Test Coverage
- [ ] Tests pass locally (`go test ./...`)
- [ ] New tests added for new functionality
- [ ] Coverage maintained or improved

### Manual Testing
- [ ] Tested in interactive mode
- [ ] Tested in batch mode (`--batch` flag)
- [ ] Tested across different operators (&&, ||, ;)
- [ ] Tested security modes (strict/permissive)
- [ ] Cross-platform testing (if applicable)

### Test Commands
```bash
# Provide specific test commands used
echo "pwd && echo test" | ./lucien --batch
# Add more test examples
```

## Security Review

- [ ] No sensitive information exposed
- [ ] Command injection protection maintained
- [ ] Security modes work correctly
- [ ] Input validation implemented

## Documentation

- [ ] Code is self-documenting
- [ ] README updated (if needed)
- [ ] User manual updated (if needed)
- [ ] Comments added for complex logic

## Checklist

- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Screenshots/Examples

If applicable, add screenshots or command examples to help explain your changes:

```bash
# Before
lucien> old behavior

# After  
lucien> new behavior
```

## Additional Notes

Any additional information that reviewers should know about this PR.