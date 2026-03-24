# Skill: Goca Code Generation Testing

**Domain:** Writing tests for Goca's code generation pipeline. Covers unit tests for validators and template generators, and integration tests that verify generated code compiles and is functionally correct.

---

## Testing Philosophy for Code Generators

Goca generates code, so tests have two distinct layers:

1. **Generator tests** — verify the generator produces the right output strings
2. **Compilation tests** — verify the generated output actually compiles as valid Go

Both layers are mandatory. A generator test alone is insufficient — template changes can produce syntactically correct template output but invalid Go code.

---

## Unit Test Patterns

### Testing `CommandValidator`

```go
func TestCommandValidator_ValidateEntityName(t *testing.T) {
    t.Parallel()
    v := NewCommandValidator()

    cases := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "Product", false},
        {"valid with numbers", "Product1", false},
        {"empty", "", true},
        {"starts with digit", "1Product", true},
        {"path traversal", "../etc", true},
        {"slash", "foo/bar", true},
        {"dot", "foo.bar", true},
        {"underscore", "foo_bar", false}, // or true depending on your convention
        {"very long", strings.Repeat("A", 1000), true},
        {"null byte", "Product\x00Name", true},
        {"shell meta", "Product;rm -rf /", true},
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            err := v.ValidateEntityName(tc.input)
            if tc.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Testing `FieldValidator`

```go
func TestFieldValidator_ParseFields(t *testing.T) {
    t.Parallel()
    v := NewFieldValidator()

    cases := []struct {
        name       string
        input      string
        wantLen    int
        wantFirst  FieldData
        wantErr    bool
    }{
        {
            name:    "single field",
            input:   "Name:string",
            wantLen: 1,
            wantFirst: FieldData{Name: "Name", Type: "string"},
        },
        {
            name:    "multiple fields",
            input:   "Name:string,Price:float64,Active:bool",
            wantLen: 3,
        },
        {
            name:    "pointer type",
            input:   "Parent:*Category",
            wantLen: 1,
        },
        {
            name:    "map type",
            input:   "Meta:map[string]string",
            wantLen: 1,
        },
        {
            name:    "empty",
            input:   "",
            wantLen: 0,
        },
        {
            name:    "missing type",
            input:   "Name:",
            wantErr: true,
        },
    }
    // ...
}
```

### Testing `TemplateGenerator`

```go
func TestTemplateGenerator_EntityTemplate(t *testing.T) {
    t.Parallel()

    data := TemplateData{
        Entity: EntityData{
            Name:      "Product",
            NameLower: "product",
            Package:   "domain",
        },
        Fields: []FieldData{
            {Name: "Name", Type: "string", JSONTag: "name", GormTag: "type:varchar(255)"},
            {Name: "Price", Type: "float64", JSONTag: "price", GormTag: "type:decimal(10,2)"},
        },
        Module: "github.com/test/myapp",
        Features: FeatureFlags{Timestamps: true},
    }

    gen := NewTemplateGenerator()
    output, err := gen.GenerateFromTemplate(entityTemplate, data)
    require.NoError(t, err)

    // Verify key patterns in output
    assert.Contains(t, output, "type Product struct")
    assert.Contains(t, output, "Name  string")
    assert.Contains(t, output, "Price float64")
    assert.Contains(t, output, "CreatedAt time.Time") // timestamps enabled
    assert.Contains(t, output, "package domain")
    assert.NotContains(t, output, "{{") // no unrendered template vars
}
```

### Testing Empty Fields Edge Case

```go
func TestTemplateGenerator_EmptyFields(t *testing.T) {
    t.Parallel()

    data := TemplateData{
        Entity: EntityData{Name: "Empty", NameLower: "empty", Package: "domain"},
        Fields: []FieldData{}, // zero fields
        Module: "github.com/test/myapp",
    }

    gen := NewTemplateGenerator()
    output, err := gen.GenerateFromTemplate(entityTemplate, data)
    require.NoError(t, err)

    assert.Contains(t, output, "type Empty struct {")
    assert.NotContains(t, output, "{{") // no unrendered vars
}
```

---

## Integration Test Patterns

### Minimal Go Module Initialization

Every integration test needs a valid Go module to compile generated code in:

```go
func initTestModule(t *testing.T, dir, moduleName string) {
    t.Helper()
    cmd := exec.Command("go", "mod", "init", moduleName)
    cmd.Dir = dir
    out, err := cmd.CombinedOutput()
    require.NoError(t, err, "go mod init failed: %s", string(out))
}
```

### Compilation Verification Helper

```go
func assertCompilesClean(t *testing.T, dir string) {
    t.Helper()

    // go build
    buildCmd := exec.Command("go", "build", "./...")
    buildCmd.Dir = dir
    buildOut, err := buildCmd.CombinedOutput()
    require.NoError(t, err, "go build failed:\n%s", string(buildOut))

    // go vet
    vetCmd := exec.Command("go", "vet", "./...")
    vetCmd.Dir = dir
    vetOut, err := vetCmd.CombinedOutput()
    require.NoError(t, err, "go vet failed:\n%s", string(vetOut))
}
```

### Full Feature Integration Test

```go
func TestFeatureCommand_GeneratesAndCompiles(t *testing.T) {
    dir := t.TempDir()
    initTestModule(t, dir, "testapp")

    // Run goca feature
    goca := exec.Command("goca", "feature", "Product",
        "--fields", "Name:string,Price:float64",
        "--database", "sqlite",
        "--handlers", "http",
    )
    goca.Dir = dir
    out, err := goca.CombinedOutput()
    require.NoError(t, err, "goca feature failed:\n%s", string(out))

    // Check files exist
    expectedFiles := []string{
        "internal/domain/product.go",
        "internal/usecase/product_service.go",
        "internal/repository/interfaces.go",
        "internal/handler/http/product_handler.go",
    }
    for _, f := range expectedFiles {
        assert.FileExists(t, filepath.Join(dir, f))
    }

    // Verify output contains no template artifacts
    for _, f := range expectedFiles {
        content, err := os.ReadFile(filepath.Join(dir, f))
        require.NoError(t, err)
        assert.NotContains(t, string(content), "{{", "unrendered template in %s", f)
        assert.NotContains(t, string(content), "}}", "unrendered template in %s", f)
    }

    // Run go mod tidy to get dependencies
    tidyCmd := exec.Command("go", "mod", "tidy")
    tidyCmd.Dir = dir
    tidyCmd.CombinedOutput() // ignore errors from missing real deps

    assertCompilesClean(t, dir)
}
```

### Safety Manager Integration Test Pattern

```go
func TestSafetyManager_DryRunPreventsDiskWrite(t *testing.T) {
    t.Parallel()
    dir := t.TempDir()
    sm := NewSafetyManager(true, false, false) // dryRun=true

    filePath := filepath.Join(dir, "output.go")
    err := sm.WriteFile(filePath, "package test\n")
    require.NoError(t, err, "dry-run should not error")

    // File must NOT exist on disk
    _, statErr := os.Stat(filePath)
    assert.True(t, os.IsNotExist(statErr), "file should not exist in dry-run mode")

    // But should be recorded in created files list
    assert.Contains(t, sm.GetCreatedFiles(), filePath)
}
```

---

## Template Content Verification Patterns

Use specific string assertions to verify generated code structure:

```go
// Verify struct fields
assert.Regexp(t, `Name\s+string\s+\x60json:"name"`, output)

// Verify method signatures
assert.Contains(t, output, "func (s *productService) Create(")
assert.Contains(t, output, "func NewProductService(repo repository.ProductRepository)")

// Verify import presence
assert.Regexp(t, `"gorm\.io/gorm"`, output)

// Verify no template artifacts
assert.NotRegexp(t, `\{\{[^}]+\}\}`, output, "unrendered template variables found")
```

---

## Mock Tests

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(entity *domain.Product) (*domain.Product, error) {
    args := m.Called(entity)
    return args.Get(0).(*domain.Product), args.Error(1)
}

func TestProductService_Create_ValidInput(t *testing.T) {
    t.Parallel()

    mockRepo := new(MockRepository)
    svc := NewProductService(mockRepo)

    input := CreateProductInput{Name: "Widget", Price: 9.99}
    expected := &domain.Product{ID: 1, Name: "Widget", Price: 9.99}

    mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(expected, nil)

    result, err := svc.Create(input)
    require.NoError(t, err)
    assert.Equal(t, expected.Name, result.Name)
    mockRepo.AssertExpectations(t)
}
```
