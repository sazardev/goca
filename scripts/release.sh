#!/bin/bash

# Script de release inteligente para Goca CLI
# Uso: ./scripts/release.sh [major|minor|patch|auto]

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar ayuda
show_help() {
    echo "Uso: $0 [major|minor|patch|auto]"
    echo ""
    echo "Tipos de release:"
    echo "  major  - Incrementa versión mayor (x.0.0) - Cambios incompatibles"
    echo "  minor  - Incrementa versión menor (x.y.0) - Nuevas funcionalidades"
    echo "  patch  - Incrementa versión de parche (x.y.z) - Correcciones"
    echo "  auto   - Detecta automáticamente basado en commits"
    echo ""
    echo "Si no se especifica tipo, usa 'auto' por defecto"
}

# Funciones de logging
log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Validaciones iniciales
if [ ! -f "go.mod" ] || [ ! -d "cmd" ]; then
    error "Este script debe ejecutarse desde la raíz del proyecto Goca"
fi

# Obtener parámetros
RELEASE_TYPE=${1:-auto}

if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

if [[ ! "$RELEASE_TYPE" =~ ^(major|minor|patch|auto)$ ]]; then
    error "Tipo de release inválido: $RELEASE_TYPE"
fi

# Validar estado del repositorio
if [ -n "$(git status --porcelain)" ]; then
    warn "Hay cambios sin confirmar en el repositorio"
    git status --short
    echo ""
    read -p "¿Continuar de todos modos? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log "Release cancelado"
        exit 0
    fi
fi

# Obtener versión actual y commits
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
CURRENT_VERSION=${CURRENT_TAG#v}

log "Versión actual: $CURRENT_VERSION"

# Si es auto, detectar tipo de release basado en commits
if [ "$RELEASE_TYPE" = "auto" ]; then
    log "Analizando commits para detectar tipo de release..."
    
    COMMITS=$(git log "$CURRENT_TAG..HEAD" --oneline 2>/dev/null || git log --oneline)
    
    if [ -z "$COMMITS" ]; then
        warn "No hay commits nuevos desde el último tag"
        exit 0
    fi
    
    echo "Commits desde $CURRENT_TAG:"
    echo "$COMMITS"
    echo ""
    
    # Detectar tipo basado en conventional commits
    if echo "$COMMITS" | grep -qE "^[a-f0-9]+ (feat|feature)(\(.+\))?!:"; then
        RELEASE_TYPE="major"
        log "Detectado: MAJOR release (breaking changes)"
    elif echo "$COMMITS" | grep -qE "^[a-f0-9]+ (feat|feature)(\(.+\))?:"; then
        RELEASE_TYPE="minor"
        log "Detectado: MINOR release (nuevas funcionalidades)"
    elif echo "$COMMITS" | grep -qE "^[a-f0-9]+ (fix|bugfix|hotfix)(\(.+\))?:"; then
        RELEASE_TYPE="patch"
        log "Detectado: PATCH release (correcciones)"
    else
        RELEASE_TYPE="patch"
        log "Detectado: PATCH release (cambios menores)"
    fi
fi

# Parsear versión actual
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Incrementar versión según el tipo
case "$RELEASE_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
NEW_TAG="v$NEW_VERSION"

# Mostrar resumen
echo ""
echo "🚀 Release Summary:"
echo "   Tipo: $RELEASE_TYPE"
echo "   Versión actual: $CURRENT_VERSION"
echo "   Nueva versión: $NEW_VERSION"
echo "   Nuevo tag: $NEW_TAG"
echo ""

read -p "¿Proceder con el release? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log "Release cancelado"
    exit 0
fi

# Ejecutar tests
log "Ejecutando tests..."
if ! make test-cli > /dev/null 2>&1; then
    error "Los tests han fallado. No se puede proceder con el release."
fi
success "Tests pasaron correctamente"

# Actualizar CHANGELOG si existe
if [ -f "CHANGELOG.md" ]; then
    log "Actualizando CHANGELOG.md..."
    
    TEMP_FILE=$(mktemp)
    echo "## [$NEW_VERSION] - $(date +%Y-%m-%d)" > "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    
    # Agregar commits desde último tag
    if [ "$CURRENT_TAG" != "v0.0.0" ]; then
        echo "### Changes" >> "$TEMP_FILE"
        git log "$CURRENT_TAG..HEAD" --pretty=format:"- %s" >> "$TEMP_FILE"
    else
        echo "### Added" >> "$TEMP_FILE"
        echo "- Initial release" >> "$TEMP_FILE"
    fi
    
    echo "" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    cat CHANGELOG.md >> "$TEMP_FILE"
    mv "$TEMP_FILE" CHANGELOG.md
    
    git add CHANGELOG.md
fi

# Build de prueba para verificar que compila
log "Verificando que el código compila..."
go build -o /tmp/goca-test . > /dev/null 2>&1
rm -f /tmp/goca-test

# Commit cambios si los hay
if [ -n "$(git status --porcelain)" ]; then
    log "Confirmando cambios..."
    git commit -m "chore: prepare release $NEW_VERSION"
fi

# Crear tag
log "Creando tag $NEW_TAG..."
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD)

git tag -a "$NEW_TAG" -m "Release $NEW_VERSION

🎉 Goca CLI $NEW_VERSION ($RELEASE_TYPE release)

Build: $BUILD_TIME
Commit: $GIT_COMMIT

## Installation
\`\`\`bash
go install github.com/sazardev/goca@$NEW_TAG
\`\`\`

## Quick Start
\`\`\`bash
goca init myproject --module github.com/user/myproject
goca feature User --fields \"name:string,email:string\"
goca version
\`\`\`
"

# Push cambios y tag
log "Enviando cambios a GitHub..."
git push origin master
git push origin "$NEW_TAG"

success "🎉 Release $NEW_VERSION creado exitosamente!"
echo ""
echo "📋 Próximos pasos:"
echo "   1. GitHub Actions construirá automáticamente el release"
echo "   2. Los binarios estarán disponibles en:"
echo "      https://github.com/sazardev/goca/releases/tag/$NEW_TAG"
echo "   3. Verificar progreso en:"
echo "      https://github.com/sazardev/goca/actions"
echo ""
echo "⏳ El proceso de build puede tomar unos minutos..."
echo ""
echo "🧪 Para probar localmente:"
echo "   go install github.com/sazardev/goca@$NEW_TAG"
