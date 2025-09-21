package testing

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CodeValidator validates generated Go code for correctness
type CodeValidator struct {
	suite *TestSuite
}

// NewCodeValidator creates a new code validator
func NewCodeValidator(suite *TestSuite) *CodeValidator {
	return &CodeValidator{suite: suite}
}

// ValidateGoSyntax checks if a Go file has valid syntax
func (v *CodeValidator) ValidateGoSyntax(filePath string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("syntax error in %s: %w", filePath, err)
	}
	return nil
}

// ValidateAllGoFiles validates syntax of all Go files in directory
func (v *CodeValidator) ValidateAllGoFiles(dir string) []error {
	var errors []error

	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") {
			if syntaxErr := v.ValidateGoSyntax(path); syntaxErr != nil {
				errors = append(errors, syntaxErr)
			}
		}
		return nil
	})

	if err != nil {
		errors = append(errors, fmt.Errorf("error walking directory: %w", err))
	}

	return errors
}

// ValidatePackageDeclaration checks if file has correct package declaration
func (v *CodeValidator) ValidatePackageDeclaration(filePath, expectedPackage string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	if node.Name.Name != expectedPackage {
		return fmt.Errorf("expected package %s, got %s", expectedPackage, node.Name.Name)
	}

	return nil
}

// ValidateImports checks if file has required imports
func (v *CodeValidator) ValidateImports(filePath string, requiredImports []string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	importMap := make(map[string]bool)
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		importMap[importPath] = true
	}

	var missing []string
	for _, required := range requiredImports {
		if !importMap[required] {
			missing = append(missing, required)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required imports: %v", missing)
	}

	return nil
}

// ValidateStructFields checks if struct has expected fields
func (v *CodeValidator) ValidateStructFields(filePath, structName string, expectedFields []string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	var structDecl *ast.StructType
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
			if st, ok := typeSpec.Type.(*ast.StructType); ok {
				structDecl = st
				return false
			}
		}
		return true
	})

	if structDecl == nil {
		return fmt.Errorf("struct %s not found in %s", structName, filePath)
	}

	fieldMap := make(map[string]bool)
	for _, field := range structDecl.Fields.List {
		for _, name := range field.Names {
			fieldMap[name.Name] = true
		}
	}

	var missing []string
	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			missing = append(missing, expected)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("struct %s missing fields: %v", structName, missing)
	}

	return nil
}

// ValidateInterfaceMethods checks if interface has expected methods
func (v *CodeValidator) ValidateInterfaceMethods(filePath, interfaceName string, expectedMethods []string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	var interfaceDecl *ast.InterfaceType
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == interfaceName {
			if it, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				interfaceDecl = it
				return false
			}
		}
		return true
	})

	if interfaceDecl == nil {
		return fmt.Errorf("interface %s not found in %s", interfaceName, filePath)
	}

	methodMap := make(map[string]bool)
	for _, method := range interfaceDecl.Methods.List {
		for _, name := range method.Names {
			methodMap[name.Name] = true
		}
	}

	var missing []string
	for _, expected := range expectedMethods {
		if !methodMap[expected] {
			missing = append(missing, expected)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("interface %s missing methods: %v", interfaceName, missing)
	}

	return nil
}

// ValidateNamingConventions checks if code follows Go naming conventions
func (v *CodeValidator) ValidateNamingConventions(filePath string) []error {
	var errors []error

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return []error{fmt.Errorf("failed to parse file: %w", err)}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.IsExported() && !isCapitalized(x.Name.Name) {
				errors = append(errors, fmt.Errorf("exported function %s should start with capital letter", x.Name.Name))
			}
		case *ast.TypeSpec:
			if x.Name.IsExported() && !isCapitalized(x.Name.Name) {
				errors = append(errors, fmt.Errorf("exported type %s should start with capital letter", x.Name.Name))
			}
		case *ast.GenDecl:
			if x.Tok == token.VAR || x.Tok == token.CONST {
				for _, spec := range x.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range vs.Names {
							if name.IsExported() && !isCapitalized(name.Name) {
								errors = append(errors, fmt.Errorf("exported %s %s should start with capital letter",
									x.Tok.String(), name.Name))
							}
						}
					}
				}
			}
		}
		return true
	})

	return errors
}

// ValidateFileHeader checks if file has proper header comments
func (v *CodeValidator) ValidateFileHeader(filePath string, expectedPatterns []string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file %s: %v\n", filePath, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	var headerLines []string

	// Read first 10 lines as header
	for i := 0; i < 10 && scanner.Scan(); i++ {
		headerLines = append(headerLines, scanner.Text())
	}

	header := strings.Join(headerLines, "\n")

	for _, pattern := range expectedPatterns {
		if !strings.Contains(header, pattern) {
			return fmt.Errorf("file header missing pattern: %s", pattern)
		}
	}

	return nil
}

// ValidateModuleReferences checks if generated code uses correct module references
func (v *CodeValidator) ValidateModuleReferences(filePath, expectedModule string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if expectedModule != "" && strings.Contains(string(content), "example.com/project") {
		return fmt.Errorf("file contains hardcoded module reference 'example.com/project', should use %s", expectedModule)
	}

	return nil
}

// isCapitalized checks if string starts with capital letter
func isCapitalized(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}
