# Guía para Publicar la Extensión de VS Code para Goca

## Requisitos Previos

1. Node.js y npm instalados
2. Cuenta en Visual Studio Marketplace
3. Personal Access Token (PAT) con permisos para publicar

## Pasos para Publicar

### 1. Instalar VSCE

```bash
npm install -g @vscode/vsce
```

### 2. Configurar package.json

Asegúrate de que tu `package.json` tenga la siguiente información:

```json
{
  "name": "goca-extension",
  "displayName": "Goca - Go Clean Architecture Assistant",
  "description": "Asistente oficial para Goca - Generador de código para Go Clean Architecture",
  "version": "0.1.0",
  "publisher": "tu-nombre-de-usuario",
  "repository": {
    "type": "git",
    "url": "https://github.com/tu-usuario/goca-vscode-extension"
  }
}
```

### 3. Inicializar proyecto e instalar dependencias

```bash
cd vscode-extension
npm install
```

### 4. Compilar la extensión

```bash
npm run compile
```

### 5. Empaquetar la extensión

```bash
vsce package
```

Esto creará un archivo `.vsix` que puedes instalar manualmente o publicar.

### 6. Publicar la extensión

```bash
vsce publish
```

Te pedirá tu Personal Access Token la primera vez.

## Actualizaciones

Para actualizar la extensión:

1. Incrementa la versión en `package.json`
2. Ejecuta `vsce publish`

## Instalar manualmente para pruebas

Para instalar la extensión manualmente:

1. En VS Code, abre la paleta de comandos (Ctrl+Shift+P)
2. Escribe "Extensions: Install from VSIX..."
3. Selecciona el archivo .vsix generado

## Notas

- Asegúrate de probar exhaustivamente la extensión antes de publicarla
- Actualiza el README.md con capturas de pantalla y ejemplos
- Añade etiquetas relevantes en el package.json para mejorar la visibilidad
