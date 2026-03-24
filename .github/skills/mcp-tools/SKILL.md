# Skill: MCP Tools for Goca

**Domain:** Maximizing the use of Serena MCP and Context7 MCP to work efficiently with the Goca codebase — minimal token usage, maximum accuracy.

---

## Serena MCP — Symbol-First Code Navigation

Serena provides symbolic understanding of Go code. Always use symbolic tools BEFORE reading raw files.

### The Serena Workflow

```
1. Get overview → 2. Find specific symbol → 3. Read body only if editing → 4. Edit symbolically
```

#### 1. Get structural overview of a file

```
mcp_oraios_serena_get_symbols_overview(relative_path="cmd/entity.go")
```

Returns: all top-level symbols (types, functions, vars) with line numbers — no code bodies.

#### 2. Find a specific symbol across the codebase

```
mcp_oraios_serena_find_symbol(name_path="TemplateData", include_body=false)
mcp_oraios_serena_find_symbol(name_path="SafetyManager/WriteFile", include_body=false)
```

- `name_path` supports `/` for nested symbols: `ClassName/MethodName`
- Set `include_body=true` only when you need to edit or understand the implementation

#### 3. Find all usages of a symbol before refactoring

```
mcp_oraios_serena_find_referencing_symbols(name_path="TemplateData", relative_path="cmd/")
```

Returns: all locations that reference `TemplateData` — safe to rename or restructure.

#### 4. Edit a complete function/method body

```
mcp_oraios_serena_replace_symbol_body(
    name_path="SafetyManager/WriteFile",
    relative_path="cmd/safety.go",
    new_body="..."
)
```

Replaces the entire body of `WriteFile` in `SafetyManager` — no need to read the full file.

#### 5. Insert new code at end of file

```
mcp_oraios_serena_insert_after_symbol(
    name_path="<last top-level symbol>",
    relative_path="cmd/entity.go",
    new_code="func newHelper() { ... }"
)
```

#### 6. Find symbol in specific file

```
mcp_oraios_serena_find_symbol(
    name_path="SafetyManager",
    relative_path="cmd/safety.go",
    include_body=false,
    depth=1  // show all direct children (methods)
)
```

### Serena for Goca — Common Patterns

| Task                                  | Serena Command                                                                      |
| ------------------------------------- | ----------------------------------------------------------------------------------- |
| What's in `cmd/safety.go`?            | `get_symbols_overview("cmd/safety.go")`                                             |
| What fields does `TemplateData` have? | `find_symbol("TemplateData", body=false)`                                           |
| How is `ValidateEntityName` called?   | `find_referencing_symbols("ValidateEntityName")`                                    |
| Edit `SafetyManager.WriteFile`        | `find_symbol("SafetyManager/WriteFile", body=true)` then `replace_symbol_body(...)` |
| What template constants exist?        | `get_symbols_overview("cmd/templates.go")`                                          |
| Add a new helper at end of file       | `get_symbols_overview` → get last symbol → `insert_after_symbol`                    |
| Rename a function safely              | `find_referencing_symbols` → update all, then `replace_symbol_body`                 |

### Serena Name Path Format

```
# Top-level symbol
"SafetyManager"           → type or function named SafetyManager
"entityCmd"               → variable named entityCmd

# Method of struct
"SafetyManager/WriteFile"      → method WriteFile on SafetyManager
"SafetyManager/CheckFileConflict"

# Nested (Go doesn't have deep nesting, but for receivers)
"CommandValidator/ValidateEntityName"
```

### When NOT to Use Serena

- For non-Go files: `.yaml`, `.md`, `.json`, `.toml` → use `read_file` directly
- For very small helper files (< 30 lines) where reading the whole file is faster
- For searching by pattern across many files → use `search_for_pattern` or `grep_search`

---

## Context7 MCP — Up-to-Date Library Documentation

### When to Use Context7

Use Context7 whenever working with:

- `github.com/spf13/cobra` — CLI flags, subcommands, completion
- `github.com/stretchr/testify` — assert/require/mock usage
- `gorm.io/gorm` — GORM v2 queries, associations, migrations
- `github.com/gorilla/mux` — router configuration
- `gopkg.in/yaml.v3` — YAML marshal/unmarshal
- `golang.org/x/text` — text transformation (used in utils.go)

### Context7 Workflow

```
1. Resolve library ID (do once per session)
2. Fetch specific topic docs
3. Apply to implementation
```

#### Step 1 — Resolve library ID

```
mcp_context7_resolve-library-id(libraryName="cobra")
→ returns: "/spf13/cobra"

mcp_context7_resolve-library-id(libraryName="testify")
→ returns: "/stretchr/testify"
```

#### Step 2 — Fetch targeted docs

```
mcp_context7_get-library-docs(
    context7LibraryId="/spf13/cobra",
    topic="persistent flags and flag groups",
    tokens=3000
)
```

Keep `tokens` ≤ 4000. Smaller = faster and cheaper. Be specific in `topic`.

### Context7 Topics for Goca

| Scenario                       | Library     | Topic                                   |
| ------------------------------ | ----------- | --------------------------------------- |
| Adding new flag type           | cobra       | `flags: string, bool, int, stringSlice` |
| Writing mock expectations      | testify     | `mock: On, Return, AssertExpectations`  |
| GORM associations in templates | gorm        | `associations: has many, belongs to`    |
| GORM auto-migrate              | gorm        | `auto migration`                        |
| Router pattern                 | gorilla/mux | `route variables, subrouters`           |
| YAML struct tags               | yaml.v3     | `struct tags, omitempty, inline`        |
| String case conversion         | x/text      | `cases: Title, Lower, Upper`            |

### Context7 Best Practices

- Resolve library IDs ONCE per session and reuse
- Request ONLY the relevant topic — don't fetch entire library docs at once
- If the first fetch doesn't have the answer, refetch with a more specific `topic`
- After getting docs, apply directly — don't re-fetch the same topic

---

## Combined MCP Strategy — Token Efficiency

### Starting a new task

```
1. Serena: get_symbols_overview(relevant file)  → understand structure
2. Serena: find_symbol(target, body=false)       → get signature
3. Context7: get-library-docs(relevant lib, specific topic)  → only if needed
4. Serena: find_symbol(target, body=true)        → read body only for editing
5. Serena: replace_symbol_body(...)              → make targeted edit
```

### Discovering where to make a change

```
1. Serena: find_symbol("FunctionName")                 → find location
2. Serena: find_referencing_symbols("FunctionName")    → understand impact
3. Edit with replace_symbol_body or insert_after_symbol
```

### What NOT to do (wastes tokens)

- ❌ `read_file(entire 500-line file)` — use Serena instead
- ❌ `mcp_context7_get-library-docs(tokens=10000)` — use ≤ 4000
- ❌ Fetching Context7 docs when the question is about Goca internals (use Serena)
- ❌ `get_symbols_overview` on every file — only the ones you need
- ❌ Re-fetching Context7 docs for the same library in the same session
