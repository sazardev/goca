# Test-Project

Proyecto generado con Goca - Go Clean Architecture Code Generator

## 🏗️ Arquitectura

Este proyecto sigue los principios de Clean Architecture:

- **Domain**: Entidades y reglas de negocio
- **Use Cases**: Lógica de aplicación
- **Repository**: Abstracción de datos
- **Handler**: Adaptadores de entrega

## 🚀 Inicio Rápido

1. Instalar dependencias:
```bash
   go mod tidy
```


2. Configurar variables de entorno:
```bash
   cp .env.example .env
```


3. Ejecutar la aplicación:
```bash
   go run cmd/server/main.go
```


## 📁 Estructura del Proyecto

```
test-project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── handler/
├── pkg/
│   ├── config/
│   └── logger/
├── go.mod
└── README.md
```


## 🔧 Comandos Útiles

Generar un nuevo feature:
```bash
goca feature User --fields "name:string,email:string"
```


Generar solo una entidad:
```bash
goca entity Product --fields "name:string,price:float64"
```

