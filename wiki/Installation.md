# Instalación

Esta página te guiará a través de los diferentes métodos para instalar Goca en tu sistema.

## 📋 Requisitos Previos

- **Go 1.21+** - [Descargar Go](https://golang.org/dl/)
- **Git** - Para clonar repositorios y gestión de versiones
- **Terminal/PowerShell** - Para ejecutar comandos

## 🚀 Métodos de Instalación

### 1. Instalación con go install (Recomendado)

Este es el método más rápido y siempre te dará la última versión estable:

```bash
go install github.com/sazardev/goca@latest
```

**Verificar instalación:**
```bash
goca version
```

**Salida esperada:**
```
Goca v1.0.5
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### 2. Descarga de Binarios

Descarga el binario pre-compilado para tu sistema operativo desde [GitHub Releases](https://github.com/sazardev/goca/releases).

#### Para Windows:
```powershell
# Descargar la última versión
Invoke-WebRequest -Uri "https://github.com/sazardev/goca/releases/latest/download/goca-windows-amd64.exe" -OutFile "goca.exe"

# Mover a una ubicación en el PATH
Move-Item goca.exe C:\Windows\System32\goca.exe
```

#### Para Linux:
```bash
# Descargar la última versión
wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64

# Hacer ejecutable y mover al PATH
chmod +x goca-linux-amd64
sudo mv goca-linux-amd64 /usr/local/bin/goca
```

#### Para macOS (Intel):
```bash
# Descargar la última versión
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64 -o goca

# Hacer ejecutable y mover al PATH
chmod +x goca
sudo mv goca /usr/local/bin/goca
```

#### Para macOS (Apple Silicon):
```bash
# Descargar la última versión
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-arm64 -o goca

# Hacer ejecutable y mover al PATH
chmod +x goca
sudo mv goca /usr/local/bin/goca
```

### 3. Instalación con Homebrew (macOS)

Si tienes Homebrew instalado:

```bash
# Agregar el tap
brew tap sazardev/tools

# Instalar goca
brew install goca
```

### 4. Compilación desde Código Fuente

Para desarrolladores que quieren la última versión de desarrollo:

```bash
# Clonar el repositorio
git clone https://github.com/sazardev/goca.git
cd goca

# Compilar
go build -o goca

# Instalar globalmente (opcional)
go install
```

## 🔧 Configuración del PATH

Si instalaste manualmente el binario, asegúrate de que esté en tu PATH:

### Windows:
1. Abre "Variables de entorno del sistema"
2. Haz clic en "Variables de entorno"
3. En "Variables del sistema", busca "Path" y haz clic en "Editar"
4. Haz clic en "Nuevo" y agrega la ruta donde guardaste `goca.exe`

### Linux/macOS:
Agrega esta línea a tu `~/.bashrc`, `~/.zshrc` o `~/.profile`:

```bash
export PATH=$PATH:/ruta/donde/guardaste/goca
```

Luego recarga tu shell:
```bash
source ~/.bashrc  # o ~/.zshrc
```

## ✅ Verificación de Instalación

Una vez instalado, verifica que todo funcione correctamente:

```bash
# Verificar versión
goca version

# Mostrar ayuda
goca help

# Probar comando básico
goca init test-project --module test
```

Si ves la información de versión y la ayuda, ¡la instalación fue exitosa! 🎉

## 🆙 Actualización

### Con go install:
```bash
go install github.com/sazardev/goca@latest
```

### Con Homebrew:
```bash
brew upgrade goca
```

### Con binarios:
Descarga la nueva versión siguiendo los pasos de instalación con binarios.

## 🐛 Solución de Problemas

### Error: "goca: command not found"
- ✅ Verifica que Goca esté en tu PATH
- ✅ Reinicia tu terminal después de la instalación
- ✅ En Windows, asegúrate de usar PowerShell o CMD como administrador

### Error: "permission denied"
```bash
# Linux/macOS - Agregar permisos de ejecución
chmod +x goca

# Windows - Ejecutar como administrador
```

### Error de certificados SSL (Windows)
```powershell
# Usar TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
```

### Versión de Go antigua
Goca requiere Go 1.21+. Actualiza Go desde [golang.org](https://golang.org/dl/).

## 🔄 Desinstalación

### Si instalaste con go install:
```bash
# Encontrar la ubicación
which goca  # Linux/macOS
where goca  # Windows

# Eliminar el binario
rm $(which goca)  # Linux/macOS
del (where goca)  # Windows
```

### Con Homebrew:
```bash
brew uninstall goca
brew untap sazardev/tools
```

## 📞 Soporte

Si tienes problemas con la instalación:

1. 🔍 Revisa los [Issues conocidos](https://github.com/sazardev/goca/issues)
2. 💬 Pregunta en [GitHub Discussions](https://github.com/sazardev/goca/discussions)
3. 🐛 Reporta un nuevo [Issue](https://github.com/sazardev/goca/issues/new)

---

**¡Siguiente paso: [Primeros Pasos](Getting-Started) →**
