# Goca VS Code Extension

Esta extensión proporciona soporte para el generador de código de Clean Architecture para Go (Goca CLI).

## Características

- Integración con comandos de Goca directamente en VS Code
- Snippets para generar código rápidamente
- Vista de explorador con comandos disponibles
- Documentación integrada
- Ayuda contextual y autocompletado

## Requisitos

- Goca CLI instalado (`go install github.com/usuario/goca@latest`)
- Go 1.16 o superior
- Visual Studio Code 1.60.0 o superior

## Instalación

1. Abrir VS Code
2. Presionar `Ctrl+Shift+X` para abrir la vista de extensiones
3. Buscar "Goca"
4. Hacer clic en "Instalar"

## Uso

### Comandos

La extensión añade los siguientes comandos a la paleta de comandos (`Ctrl+Shift+P`):

- `Goca: Mostrar versión`
- `Goca: Inicializar proyecto`
- `Goca: Generar feature completo`
- `Goca: Generar entidad`
- `Goca: Generar caso de uso`
- `Goca: Generar repositorio`
- `Goca: Generar handler`
- `Goca: Generar mensajes`
- `Goca: Generar inyección de dependencias`
- `Goca: Generar interfaces`

### Explorador de Goca

La extensión añade una vista de explorador en la barra de actividad que muestra todos los comandos disponibles de Goca.

### Snippets

La extensión proporciona los siguientes snippets para Go:

- `goca-entity`: Crear una entidad de dominio
- `goca-repo-interface`: Crear una interfaz de repositorio
- `goca-usecase-interface`: Crear una interfaz de caso de uso
- `goca-rest-handler`: Crear un handler REST
- `goca-di-container`: Crear un contenedor de inyección de dependencias

## Configuración

La extensión proporciona las siguientes configuraciones:

- `goca.path`: Ruta al ejecutable de Goca (por defecto: "goca")
- `goca.enableSnippets`: Habilitar snippets de Goca (por defecto: true)
- `goca.defaultProjectStructure`: Estructura de proyecto predeterminada (por defecto: "standard")

## Contribuir

Si quieres contribuir a esta extensión, por favor visita el repositorio en GitHub.

## Licencia

Esta extensión está licenciada bajo la Licencia MIT.
