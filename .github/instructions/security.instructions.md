---
applyTo: "**/*.go"
---

# Security Rules for Goca

## Input Validation — OWASP A03 (Injection Prevention)

### Entity and Field Name Validation

ALL user-supplied names (entity names, field names, project names) MUST pass through `CommandValidator` before ANY use in:

- File path construction
- Template rendering
- Shell commands

```go
// REQUIRED — validate before use
validator := NewCommandValidator()
if err := validator.ValidateEntityName(entityName); err != nil {
    validator.errorHandler.HandleError(err, "entity name")
}
// Only proceed after validation passes

// FORBIDDEN — using raw user input directly
filePath := filepath.Join(outputDir, args[0]+".go") // ← path traversal risk
```

### Allowed Name Pattern

Names MUST match `^[A-Za-z][A-Za-z0-9]*$`:

- Must start with a letter
- Only alphanumeric characters
- No dots, slashes, backslashes, underscores in entity names
- No null bytes, shell metacharacters, Unicode tricks

### Path Construction Safety

```go
// REQUIRED — use filepath.Join with validated components only
outputPath := filepath.Join(projectRoot, "internal", "domain", validated+".go")

// FORBIDDEN — string concatenation with user input
outputPath := projectRoot + "/" + userInput + ".go"

// VERIFY — constructed path stays within project root
if !strings.HasPrefix(filepath.Clean(outputPath), filepath.Clean(projectRoot)) {
    return fmt.Errorf("security: path escapes project root: %s", outputPath)
}
```

## File System Safety

### Never Write Outside Project Root

```go
// SafetyManager already enforces this, but verify in any direct writes:
func writeFileSecurely(projectRoot, relativePath, content string) error {
    absPath := filepath.Join(projectRoot, relativePath)
    // Verify no escape
    if !strings.HasPrefix(filepath.Clean(absPath), filepath.Clean(projectRoot)) {
        return fmt.Errorf("path traversal attempt blocked: %s", relativePath)
    }
    return os.WriteFile(absPath, []byte(content), 0644)
}
```

### Required File Permissions

- Source files: `0644` (owner rw, group r, other r)
- Directories: `0755`
- Never `0777`, `0666`, or executable bits on data files

### No File Deletion

`SafetyManager` only creates backups — never deletes. **Never add delete operations to safety manager or command code.**

## Shell Command Security

### Use `exec.Command` With Separate Arguments

NEVER construct shell commands by concatenating strings with user input:

```go
// FORBIDDEN — shell injection vulnerability
cmd := exec.Command("sh", "-c", "go get "+userInput)

// CORRECT — separate arguments, no shell interpretation
cmd := exec.Command("go", "get", validatedModulePath)
cmd.Dir = projectRoot
```

### Module Path Validation Before `go get`

```go
// Validate module path matches safe pattern before go get
var modulePathPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9./\-]*(@[a-z0-9.\-]+)?$`)

func validateModulePath(path string) error {
    if !modulePathPattern.MatchString(path) {
        return fmt.Errorf("invalid module path: %s", path)
    }
    return nil
}
```

## Template Security

### Use `text/template`, Not `html/template`

Goca generates Go source code, not HTML. Use `text/template`. The data passed to templates MUST be pre-validated author-controlled data (from `TemplateData`) — never raw user strings.

```go
// CORRECT — TemplateData fields are all validated before template execution
tmpl, err := template.New("entity").Parse(entityTemplate)
if err != nil {
    return fmt.Errorf("parse entity template: %w", err)
}
err = tmpl.Execute(buf, data) // data is TemplateData with validated fields
```

### No Dynamic Template Construction

Templates MUST be static string constants. Never build template strings by concatenating user input:

```go
// FORBIDDEN
templateStr := "type " + userEntityName + " struct {" // ← template injection

// CORRECT — use TemplateData to parameterize static templates
const entityTemplate = `type {{.Entity.Name}} struct { ... }`
```

## Dependency Security

### Verify Go Module Path Format

Before executing `go get` or adding to `go.mod`, validate:

- Path does not contain shell metacharacters
- Path matches expected module path format
- Version specifier (if provided) follows semver: `@v1.2.3` or `@latest`

## Sensitive Data Handling

- Never log user file contents to stdout/stderr
- Never print template output in error messages (may contain user data structure)
- Config files (`.goca.yaml`) may contain database credentials — never surface their contents in CLI output

## OWASP Top 10 Checklist for CLI Tools

| Risk                          | Mitigation in Goca                                                    |
| ----------------------------- | --------------------------------------------------------------------- |
| A01 Broken Access Control     | `SafetyManager` prevents writes outside project root                  |
| A02 Cryptographic Failures    | No crypto in scope; file perms `0644`                                 |
| A03 Injection                 | `CommandValidator` on all user inputs; `exec.Command` with arg arrays |
| A04 Insecure Design           | `SafetyManager` dry-run/backup; no destructive ops                    |
| A05 Security Misconfiguration | `0644`/`0755` file perms; no world-writable files                     |
| A06 Vulnerable Components     | Validate module paths before `go get`                                 |
| A07 Auth Failures             | CLI tool — N/A                                                        |
| A08 Software Integrity        | `SafetyManager` backup before overwrite                               |
| A09 Security Logging          | Errors logged to stderr; no sensitive data in logs                    |
| A10 SSRF                      | No HTTP calls in CLI core; `go get` uses validated paths              |
