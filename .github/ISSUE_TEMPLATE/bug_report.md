---
name: Bug Report
about: Report a bug to help us improve Goca
title: '[BUG] '
labels: 'bug'
assignees: ''

---

## Bug Description

A clear and concise description of what the bug is.

## Steps to Reproduce

Please provide detailed steps to reproduce the behavior:

1. Run command: `goca ...`
2. With configuration: `...`
3. In directory structure: `...`
4. Observe error: `...`

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

A clear and concise description of what actually happened.

## Environment Information

**Goca Version:**
```
Output of: goca version
```

**Go Version:**
```
Output of: go version
```

**Operating System:**
- OS: [e.g., Ubuntu 22.04, macOS 14.0, Windows 11]
- Architecture: [e.g., amd64, arm64]

**Installation Method:**
- [ ] Binary from GitHub Releases
- [ ] Built from source
- [ ] Go install

## Generated Code Context

**Command Used:**
```bash
goca feature User --fields "name:string,email:string"
```

**Project Structure (if relevant):**
```
project/
├── cmd/
├── internal/
└── ...
```

**Configuration Files:**
If you have custom configuration, please share relevant parts of `goca.yaml` or similar:
```yaml
# Configuration here
```

## Error Output

**Terminal Output:**
```
Paste the complete error output here
```

**Log Files:**
If applicable, include relevant log output.

## Screenshots

If applicable, add screenshots to help explain your problem.

## Code Snippets

If relevant, include snippets of generated code or configuration:

```go
// Generated code that shows the issue
```

## Attempted Solutions

Have you tried any solutions? If so, what were the results?

- [ ] Searched existing issues
- [ ] Checked documentation
- [ ] Tried with --dry-run flag
- [ ] Verified installation
- [ ] Updated to latest version

## Impact Assessment

How does this bug affect your workflow?

- [ ] Blocks project setup
- [ ] Prevents feature generation
- [ ] Causes incorrect code generation
- [ ] Minor inconvenience
- [ ] Other: [please describe]

## Possible Solution

If you have suggestions on how to fix the bug, please describe them here.

## Additional Context

Add any other context about the problem here. This might include:

- When did this issue start occurring?
- Does it happen consistently or intermittently?
- Are there any workarounds you've found?
- Any related issues or pull requests?

## Checklist

Before submitting, please ensure:

- [ ] I have searched for existing issues
- [ ] I have provided all required environment information
- [ ] I have included steps to reproduce the issue
- [ ] I have described the expected vs actual behavior
- [ ] I have included relevant code/configuration snippets
- [ ] I have checked this occurs on the latest version
