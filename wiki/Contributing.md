# Guía de Contribución

¡Gracias por tu interés en contribuir a Goca! Esta guía te ayudará a entender cómo puedes participar en el desarrollo del proyecto.

## 🎯 Formas de Contribuir

### 🐛 Reportar Bugs
- Usa el [template de bug report](https://github.com/sazardev/goca/issues/new?template=bug_report.md)
- Incluye información de versión (`goca version`)
- Proporciona pasos para reproducir el problema
- Incluye ejemplos de código si es relevante

### 💡 Sugerir Features
- Usa el [template de feature request](https://github.com/sazardev/goca/issues/new?template=feature_request.md)
- Explica el caso de uso y beneficios
- Considera la compatibilidad con Clean Architecture
- Discute la implementación en issues antes de codificar

### 📖 Mejorar Documentación
- Corregir errores tipográficos
- Agregar ejemplos y casos de uso
- Traducir documentación
- Mejorar la claridad de explicaciones

### 🔧 Contribuir Código
- Implementar nuevas características
- Corregir bugs existentes
- Optimizar performance
- Agregar tests

## 🚀 Configurar Entorno de Desarrollo

### Prerrequisitos
- **Go 1.21+**
- **Git**
- **Make** (opcional)

### Setup Inicial
```bash
# 1. Fork el repositorio en GitHub
# 2. Clonar tu fork
git clone https://github.com/tu-usuario/goca.git
cd goca

# 3. Agregar remote upstream
git remote add upstream https://github.com/sazardev/goca.git

# 4. Instalar dependencias
go mod tidy

# 5. Verificar que todo funciona
go build
./goca version
```

### Estructura del Proyecto de Desarrollo
```
goca/
├── cmd/                     # Comandos CLI
│   ├── di.go               # Comando di
│   ├── entity.go           # Comando entity
│   ├── feature.go          # Comando feature
│   ├── handler.go          # Comando handler
│   ├── init.go             # Comando init
│   ├── repository.go       # Comando repository
│   ├── usecase.go          # Comando usecase
│   ├── version.go          # Comando version
│   └── utils.go            # Utilidades comunes
├── examples/               # Ejemplos y demos
├── scripts/                # Scripts de automatización
├── wiki/                   # Documentación wiki
├── .github/workflows/      # CI/CD
├── go.mod
├── main.go
└── README.md
```

## 📝 Proceso de Desarrollo

### 1. Crear Branch
```bash
# Actualizar main
git checkout main
git pull upstream main

# Crear branch para feature/fix
git checkout -b feature/nueva-funcionalidad
# o
git checkout -b fix/descripcion-del-bug
```

### 2. Desarrollo
```bash
# Hacer cambios
# Ejecutar tests
go test ./...

# Verificar que compila
go build

# Probar manualmente
./goca help
```

### 3. Commit Guidelines
Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Features
git commit -m "feat: agregar soporte para MongoDB en repositorios"

# Bug fixes
git commit -m "fix: corregir validación de email en entidades"

# Documentation
git commit -m "docs: actualizar ejemplos en README"

# Tests
git commit -m "test: agregar tests para comando feature"

# Refactoring
git commit -m "refactor: simplificar generación de DTOs"
```

### 4. Push y Pull Request
```bash
# Push branch
git push origin feature/nueva-funcionalidad

# Crear Pull Request en GitHub
# Usar el template proporcionado
# Incluir descripción detallada
# Referenciar issues relacionados
```

## 🧪 Testing

### Ejecutar Tests
```bash
# Todos los tests
go test ./...

# Tests con coverage
go test -cover ./...

# Tests verbosos
go test -v ./...

# Tests específicos
go test ./cmd -run TestEntityGeneration
```

### Escribir Tests
```go
func TestGenerateEntity(t *testing.T) {
    tests := []struct {
        name     string
        entity   string
        fields   string
        expected string
    }{
        {
            name:     "basic entity",
            entity:   "User",
            fields:   "name:string,email:string",
            expected: "package domain",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := generateEntity(tt.entity, tt.fields, false, false, false, false)
            if !strings.Contains(result, tt.expected) {
                t.Errorf("Expected %s to contain %s", result, tt.expected)
            }
        })
    }
}
```

### Tests de Integración
```bash
# Crear proyecto de prueba
mkdir test-project
cd test-project

# Probar comando init
../goca init test --module github.com/test/test

# Verificar estructura
ls -la test/

# Probar generación de features
../goca feature User --fields "name:string,email:string"

# Verificar que compila
cd test && go mod tidy && go build
```

## 📚 Agregar Nueva Funcionalidad

### 1. Nuevo Comando
Para agregar un nuevo comando (ej: `goca migrate`):

```go
// cmd/migrate.go
package cmd

import (
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Generar migraciones de base de datos",
    Long:  `Descripción larga del comando...`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementación
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
    
    // Flags
    migrateCmd.Flags().StringP("database", "d", "postgres", "Tipo de base de datos")
}
```

### 2. Nueva Funcionalidad en Comando Existente
Para agregar un flag o modificar comportamiento:

```go
// En el comando existente
func init() {
    // Nuevo flag
    featureCmd.Flags().BoolP("swagger", "s", false, "Generar documentación Swagger")
}

// En la función principal
swagger, _ := cmd.Flags().GetBool("swagger")
if swagger {
    generateSwaggerDocs(featureName)
}
```

### 3. Nuevas Plantillas
Para agregar soporte para nuevas tecnologías:

```go
// cmd/repository.go
func generateRedisRepository(dir, entity string) {
    content := `package redis

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type %sRepository struct {
    client *redis.Client
}

func New%sRepository(client *redis.Client) *%sRepository {
    return &%sRepository{
        client: client,
    }
}
`
    content = fmt.Sprintf(content, entity, entity, entity, entity)
    writeFile(filepath.Join(dir, strings.ToLower(entity)+"_repository.go"), content)
}
```

## 🎨 Estándares de Código

### Formateo
```bash
# Formatear código
go fmt ./...

# Imports organizados
goimports -w .

# Linting
golangci-lint run
```

### Convenciones
- **Funciones públicas**: PascalCase con comentarios
- **Variables**: camelCase descriptivas
- **Constantes**: UPPER_SNAKE_CASE
- **Archivos**: snake_case.go
- **Packages**: lowercase, singular

### Comentarios
```go
// generateEntity crea una nueva entidad de dominio con los campos especificados.
// Parámetros:
//   - entityName: nombre de la entidad (ej: "User")
//   - fields: campos separados por coma (ej: "name:string,email:string")
//   - validation: si incluir validaciones automáticas
//   - businessRules: si generar métodos de reglas de negocio
func generateEntity(entityName, fields string, validation, businessRules bool) string {
    // Implementación...
}
```

## 🚀 Release Process

### Versionado
Seguimos [Semantic Versioning](https://semver.org/):
- **MAJOR**: Cambios incompatibles en API
- **MINOR**: Nuevas funcionalidades compatibles
- **PATCH**: Correcciones de bugs compatibles

### Proceso de Release
```bash
# 1. Actualizar version.go
# cmd/version.go
var Version = "1.1.0"

# 2. Actualizar CHANGELOG.md
# Agregar nueva sección con cambios

# 3. Commit y tag
git commit -m "release: v1.1.0"
git tag v1.1.0
git push origin main --tags

# 4. GitHub Actions automáticamente:
# - Ejecuta tests
# - Compila binarios
# - Crea release en GitHub
# - Publica en repositorios
```

## 📖 Documentación

### Wiki
La documentación está en el directorio `wiki/`:

```bash
# Editar documentación
vim wiki/Command-Entity.md

# Verificar markdown
markdownlint wiki/*.md

# Previsualizar localmente
cd wiki && python -m http.server 8000
```

### README
- Mantener ejemplos actualizados
- Incluir casos de uso comunes
- Verificar que enlaces funcionen

### Comentarios en Código
- Documentar funciones públicas
- Explicar algoritmos complejos
- Incluir ejemplos de uso

## 🤝 Community Guidelines

### Comunicación
- **Ser respetuoso** y constructivo
- **Ayudar a newcomers** con paciencia
- **Discutir ideas** antes de implementar
- **Dar feedback** útil en code reviews

### Code Review
- **Revisar lógica** y arquitectura
- **Verificar tests** están incluidos
- **Comprobar documentación** está actualizada
- **Sugerir mejoras** constructivamente

### Issues y Discussions
- **Buscar duplicados** antes de crear
- **Usar templates** apropiados
- **Proporcionar contexto** completo
- **Seguir up** en conversaciones

## 🏆 Reconocimiento

### Contributors
Todos los contributors son reconocidos en:
- README.md
- Release notes
- Contributors page

### Tipos de Contribución
- 💻 **Code**: Implementación de features y fixes
- 📖 **Documentation**: Mejoras en docs y ejemplos
- 🐛 **Bug Reports**: Identificación y reporte de issues
- 💡 **Ideas**: Sugerencias y discusiones de features
- 🎨 **Design**: UX/UI y arquitectura
- 🔍 **Testing**: Escritura y mejora de tests

## 📞 Contacto

### Canales de Comunicación
- **GitHub Issues**: Para bugs y feature requests
- **GitHub Discussions**: Para preguntas y discusiones
- **Email**: sazardev@example.com (mantenedor principal)

### Respuesta Esperada
- **Issues**: 24-48 horas
- **Pull Requests**: 2-7 días
- **Discussions**: 1-3 días

## 📋 Checklist para Contributors

### Antes de Enviar PR
- [ ] Tests pasan (`go test ./...`)
- [ ] Código formateado (`go fmt ./...`)
- [ ] Documentación actualizada
- [ ] CHANGELOG.md actualizado (para features)
- [ ] Commits siguen convenciones
- [ ] Branch está actualizado con main

### Para Maintainers
- [ ] Code review completo
- [ ] Tests de integración pasan
- [ ] Documentación revisada
- [ ] Breaking changes documentados
- [ ] Release notes preparadas

---

**¡Gracias por contribuir a Goca! Tu participación hace que este proyecto sea mejor para toda la comunidad. 🙏**

**← [Troubleshooting](Troubleshooting) | [Development](Development) →**
