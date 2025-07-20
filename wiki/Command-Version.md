# Comando goca version

El comando `goca version` proporciona información detallada sobre la versión instalada de Goca, incluyendo metadatos de compilación y compatibilidad.

## 📋 Sintaxis

```bash
goca version [flags]
```

## 🎯 Propósito

Muestra información completa sobre la instalación actual de Goca:

- 🏷️ **Número de versión** - Versión semántica actual
- 📅 **Fecha de compilación** - Cuándo fue compilada esta versión
- 🔧 **Versión de Go** - Versión de Go utilizada para compilar
- 📦 **Información de build** - Metadatos adicionales de compilación

## 🚩 Flags Disponibles

| Flag             | Tipo   | Requerido | Valor por Defecto | Descripción                       |
| ---------------- | ------ | --------- | ----------------- | --------------------------------- |
| `--short` / `-s` | `bool` | ❌ No      | `false`           | Muestra solo el número de versión |

## 📖 Ejemplos de Uso

### Información Completa
```bash
goca version
```

**Salida:**
```
Goca v1.0.5
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### Versión Corta
```bash
goca version --short
# o
goca version -s
```

**Salida:**
```
1.0.5
```

## 🔍 Información Detallada

### Número de Versión
- Sigue el formato **Semantic Versioning (SemVer)**
- Formato: `MAJOR.MINOR.PATCH`
- Ejemplo: `1.0.5` significa:
  - **Major (1)**: Cambios incompatibles en la API
  - **Minor (0)**: Nuevas funcionalidades compatibles
  - **Patch (5)**: Correcciones de bugs compatibles

### Fecha de Compilación
- Formato **ISO 8601**: `YYYY-MM-DDTHH:MM:SSZ`
- Siempre en **UTC**
- Indica cuándo se compiló el binario específico

### Versión de Go
- Muestra la versión exacta de Go utilizada
- Importante para **compatibilidad** y **debugging**
- Formato: `go1.XX.Y`

## 🛠️ Casos de Uso

### Verificar Instalación
```bash
# Comprobar que Goca está instalado correctamente
goca version
```

### Scripts de Automatización
```bash
#!/bin/bash

# Obtener solo la versión para scripts
VERSION=$(goca version --short)
echo "Usando Goca v$VERSION"

# Verificar versión mínima requerida
REQUIRED="1.0.0"
if [[ "$(printf '%s\n' "$REQUIRED" "$VERSION" | sort -V | head -n1)" != "$REQUIRED" ]]; then
    echo "Error: Se requiere Goca v$REQUIRED o superior"
    exit 1
fi
```

### CI/CD Integration
```yaml
# GitHub Actions
- name: Check Goca Version
  run: |
    goca version
    GOCA_VERSION=$(goca version --short)
    echo "GOCA_VERSION=$GOCA_VERSION" >> $GITHUB_ENV
```

### Debugging
```bash
# Información completa para reportes de bugs
goca version > goca-version.txt
echo "Sistema: $(uname -a)" >> goca-version.txt
echo "Go instalado: $(go version)" >> goca-version.txt
```

## 📊 Análisis de Versiones

### Versiones de Desarrollo
```bash
# Versiones de desarrollo pueden incluir sufijos
goca version
# Output: Goca v1.1.0-dev
```

### Versiones Release Candidate
```bash
# Versiones candidatas a release
goca version
# Output: Goca v1.1.0-rc.1
```

### Versiones Estables
```bash
# Versiones finales sin sufijos
goca version
# Output: Goca v1.0.5
```

## 🔄 Compatibilidad

### Compatibilidad con Go
| Versión Goca | Go Mínimo | Go Recomendado | Notas           |
| ------------ | --------- | -------------- | --------------- |
| v1.0.x       | Go 1.21   | Go 1.24+       | Versión estable |
| v1.1.x       | Go 1.22   | Go 1.24+       | Próxima versión |

### Compatibilidad de Features
```bash
# Verificar si tu versión soporta una característica
goca version

# Comparar con documentación de features:
# v1.0.0: Funcionalidades básicas
# v1.0.1: Correcciones de bugs
# v1.0.5: Mejoras en gRPC y validaciones
```

## 🚀 Actualizaciones

### Verificar si Hay Actualizaciones
```bash
# Versión actual
CURRENT=$(goca version --short)
echo "Versión actual: v$CURRENT"

# Verificar última versión en GitHub (requiere curl/jq)
LATEST=$(curl -s https://api.github.com/repos/sazardev/goca/releases/latest | jq -r .tag_name)
echo "Última versión: $LATEST"

if [ "v$CURRENT" != "$LATEST" ]; then
    echo "¡Actualización disponible!"
    echo "Ejecuta: go install github.com/sazardev/goca@latest"
fi
```

### Actualizar a Última Versión
```bash
# Actualizar usando go install
go install github.com/sazardev/goca@latest

# Verificar actualización
goca version
```

### Instalar Versión Específica
```bash
# Instalar versión específica
go install github.com/sazardev/goca@v1.0.5

# Verificar versión instalada
goca version
```

## 🔍 Información de Build Detallada

### Variables de Build
El comando `version` muestra información que se compila en tiempo de build:

```go
// Definidas en cmd/version.go
var (
    Version   = "1.0.5"                    // Versión del software
    BuildTime = "2025-07-19T15:00:00Z"     // Timestamp de compilación
    GoVersion = runtime.Version()          // Versión de Go runtime
)
```

### Build Tags y Flags
```bash
# Información de compilación (si está disponible)
goca version --verbose  # (si se implementa en futuras versiones)
```

## 📝 Formato de Salida

### Formato Normal
```
Goca v1.0.5
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### Formato Corto
```
1.0.5
```

### Formato JSON (Futuro)
```bash
# Posible implementación futura
goca version --json
```

```json
{
  "version": "1.0.5",
  "build_time": "2025-07-19T15:00:00Z",
  "go_version": "go1.24.5",
  "git_commit": "abc123def",
  "build_user": "github-actions"
}
```

## 🐛 Troubleshooting

### Comando No Encontrado
```bash
# Error: command not found
which goca          # Linux/macOS
where goca          # Windows

# Verificar PATH
echo $PATH          # Linux/macOS
echo $env:PATH      # PowerShell
```

### Versión Antigua
```bash
# Verificar múltiples instalaciones
which -a goca       # Linux/macOS

# Limpiar cache de Go
go clean -modcache

# Reinstalar
go install github.com/sazardev/goca@latest
```

### Información Inconsistente
```bash
# Verificar integridad
goca version

# Comparar con archivo go.mod del proyecto
cat go.mod | grep goca

# Verificar en GitHub
curl -s https://api.github.com/repos/sazardev/goca/releases/latest
```

## 📞 Soporte y Reportes

### Incluir en Reportes de Bugs
Siempre incluye la salida de `goca version` en reportes de bugs:

```bash
# Información para reportes
echo "=== GOCA VERSION INFO ===" > bug-report.txt
goca version >> bug-report.txt
echo "=== SYSTEM INFO ===" >> bug-report.txt
uname -a >> bug-report.txt
go version >> bug-report.txt
```

### Links Útiles
- 🐛 **Issues**: [GitHub Issues](https://github.com/sazardev/goca/issues)
- 📋 **Releases**: [GitHub Releases](https://github.com/sazardev/goca/releases)
- 📖 **Changelog**: [CHANGELOG.md](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)

## 🔄 Historial de Versiones

### Versiones Importantes

#### v1.0.5 (Actual)
- ✅ Mejoras en generación de gRPC
- ✅ Validaciones mejoradas
- ✅ Correcciones de bugs

#### v1.0.0 (Release Inicial)
- 🎉 Lanzamiento inicial
- ✅ Funcionalidades básicas de Clean Architecture
- ✅ Soporte para múltiples bases de datos
- ✅ Handlers HTTP y gRPC

### Próximas Versiones
- 🔮 **v1.1.0**: Soporte para microservicios
- 🔮 **v1.2.0**: Templates personalizables
- 🔮 **v2.0.0**: Rewrite con mejoras de performance

---

**← [Comando goca messages](Command-Messages) | [Primeros Pasos](Getting-Started) →**
