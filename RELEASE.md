# Automated Release System

This document describes the fully automated release system for Goca CLI.

## 🎯 System Features 

### 1. **Dynamic Versions**
- ✅ Version injected at compile time from git tags
- ✅ No more hardcoded versions in code
- ✅ Automatic build information (timestamp, commit)

### 2. **Automatic Release by Commits**
- ✅ Automatic release type detection based on [Conventional Commits](https://www.conventionalcommits.org/)
- ✅ Automatic semantic version increment
- ✅ Automatic tag and release generation

### 3. **GitHub Actions Workflows**
- ✅ Automatic tests before each release
- ✅ Automatic multi-platform build
- ✅ Automatic publishing on GitHub Releases

## 📋 Release Types

### Automatic Detection by Commits

| Commit Pattern              | Release Type | Increment |
| --------------------------- | ------------ | --------- |
| `feat!:`                    | **major**    | `x.0.0`   |
| `feat:`                     | **minor**    | `x.y.0`   |
| `fix:`                      | **patch**    | `x.y.z`   |
| `chore:`, `docs:`, `style:` | **patch**    | `x.y.z`   |

### Commit Examples

```bash
# Major release (breaking changes)
git commit -m "feat!: change API interface for better performance"

# Minor release (new features)
git commit -m "feat: add support for MongoDB integration"

# Patch release (bug fixes)
git commit -m "fix: resolve template generation error in handlers"

# Patch release (maintenance)
git commit -m "chore: update dependencies and improve documentation"
```

## 🚀 Métodos de Release

### 1. **Release Automático (Recomendado)**

Solo haz push de tus commits y el sistema detectará automáticamente si necesita hacer un release:

```bash
git add .
git commit -m "feat: add new entity validation system"
git push origin master
```

**¿Qué pasa automáticamente?**
1. GitHub Actions analiza los commits desde el último tag
2. Determina el tipo de release basado en los mensajes
3. Incrementa la versión automáticamente
4. Crea el tag y release
5. Compila binarios para todas las plataformas
6. Publica el release en GitHub

### 2. **Release Manual con Script**

Si necesitas control manual, usa el script de release:

```bash
# Auto-detecta el tipo basado en commits
./scripts/release.sh auto

# Especifica el tipo manualmente
./scripts/release.sh patch   # x.y.Z
./scripts/release.sh minor   # x.Y.0
./scripts/release.sh major   # X.0.0
```

### 3. **Release con Makefile**

```bash
# Release automático (detecta tipo por commits)
make release

# Releases específicos
make release-patch
make release-minor  
make release-major
```

## 🔧 Configuración del Entorno

### Para Desarrolladores

```bash
# 1. Setup inicial
make dev-setup

# 2. Verificar que todo está listo
make pre-release-check

# 3. Ver versión actual
make version
```

### Variables de Build

El sistema inyecta estas variables automáticamente:

```go
// cmd/version.go
var (
    Version   = "dev"        // Se inyecta desde git tags
    BuildTime = "unknown"    // Timestamp de compilación
    GitCommit = "unknown"    // Hash del commit
    GoVersion = runtime.Version() // Versión de Go
)
```

## 📦 Compilación con Versiones

### Desarrollo Local

```bash
# Build con versión automática
make build

# Ver información de versión
./goca version
```

### Build Manual con Versión Específica

```bash
go build -ldflags "-X github.com/sazardev/goca/cmd.Version=1.2.5 \
                   -X github.com/sazardev/goca/cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
                   -X github.com/sazardev/goca/cmd.GitCommit=$(git rev-parse --short HEAD)" \
         -o goca .
```

## 🎛️ Configuración de Workflows

### Auto Release (.github/workflows/auto-release.yml)

**Trigger:** Push a `master` branch  
**Funcionalidad:**
- Analiza commits desde último tag
- Determina tipo de release
- Crea tag automáticamente
- Dispara el workflow de release

### Release (.github/workflows/release.yml)

**Trigger:** Creación de tags `v*`  
**Funcionalidad:**
- Compila para múltiples plataformas
- Crea checksums
- Publica release en GitHub
- Genera instalación automática

### Test (.github/workflows/test.yml)

**Trigger:** Push y Pull Requests  
**Funcionalidad:**
- Ejecuta suite de tests optimizada
- Validación de código generado
- Verificación de compilación

## 🚦 Flujo de Trabajo Completo

### 1. **Desarrollo**
```bash
# Hacer cambios
git add .
git commit -m "feat: add new authentication system"
```

### 2. **Push**
```bash
git push origin master
```

### 3. **Automático**
- ✅ Tests automáticos
- ✅ Detección de tipo: `minor` (por `feat:`)
- ✅ Nueva versión: `v1.3.0`
- ✅ Tag automático
- ✅ Build multiplataforma
- ✅ Release publicado

### 4. **Instalación para Usuarios**
```bash
go install github.com/sazardev/goca@v1.3.0
```

## 🛠️ Comandos Útiles

```bash
# Ver estado actual
make status

# Verificar antes de release
make pre-release-check

# Build para todas las plataformas
make build-all

# Ver logs de commits para próximo release
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Ver último tag
git describe --tags --abbrev=0
```

## ⚙️ Configuración Avanzada

### Saltarse Auto-Release

Si no quieres que ciertos commits disparen releases:

```bash
# Solo commits de mantenimiento (no disparan release automático)
git commit -m "chore: update documentation"
git commit -m "docs: improve README"
git commit -m "style: fix formatting"

# Para forzar release en commits de mantenimiento
git commit -m "chore: update dependencies [release]"
```

### Rollback de Release

```bash
# Eliminar tag local y remoto
git tag -d v1.2.5
git push origin :refs/tags/v1.2.5

# Eliminar release en GitHub manualmente
# https://github.com/sazardev/goca/releases
```

## 🎉 Beneficios del Sistema

1. **🤖 Totalmente Automatizado:** Solo push y listo
2. **📏 Versionado Semántico:** Cumple estrictamente con SemVer
3. **🔄 Consistente:** Misma versión en código, tags y releases
4. **⚡ Rápido:** Proceso completo en ~2-3 minutos
5. **🛡️ Seguro:** Tests automáticos antes de cada release
6. **📱 Multiplataforma:** Binarios para Windows, Linux, macOS
7. **📋 Trazabilidad:** Historial completo de cambios

## 🔗 Enlaces Importantes

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://github.com/sazardev/goca/releases)
- [GitHub Actions](https://github.com/sazardev/goca/actions)