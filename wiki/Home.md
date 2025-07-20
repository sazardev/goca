# Goca Wiki - Go Clean Architecture Code Generator

¡Bienvenido a la documentación oficial de **Goca**! 🎉

Goca es un potente generador de código CLI para Go que te ayuda a crear proyectos siguiendo los principios de **Clean Architecture** de Uncle Bob. Esta wiki te guiará paso a paso para aprovechar al máximo todas las funcionalidades de Goca.

## 📚 Índice de Contenidos

### 🚀 Inicio Rápido
- [**Instalación**](Installation) - Cómo instalar Goca en tu sistema
- [**Primeros Pasos**](Getting-Started) - Tu primer proyecto con Goca
- [**Tutorial Completo**](Complete-Tutorial) - Ejemplo paso a paso de un proyecto real

### 📖 Referencia de Comandos
- [**goca init**](Command-Init) - Inicializar proyectos Clean Architecture
- [**goca feature**](Command-Feature) - Generar features completos
- [**goca entity**](Command-Entity) - Crear entidades de dominio
- [**goca usecase**](Command-UseCase) - Generar casos de uso
- [**goca repository**](Command-Repository) - Crear repositorios
- [**goca handler**](Command-Handler) - Generar adaptadores de entrada
- [**goca di**](Command-DI) - Inyección de dependencias
- [**goca interfaces**](Command-Interfaces) - Generar interfaces para TDD
- [**goca messages**](Command-Messages) - Mensajes y constantes
- [**goca version**](Command-Version) - Información de versión

### 🏗️ Arquitectura y Conceptos
- [**Clean Architecture**](Clean-Architecture) - Principios y estructura
- [**Estructura de Proyecto**](Project-Structure) - Organización de directorios
- [**Patrones Implementados**](Design-Patterns) - Patrones de diseño utilizados
- [**Buenas Prácticas**](Best-Practices) - Recomendaciones y convenciones

### 💡 Ejemplos y Casos de Uso
- [**Proyecto E-commerce**](Example-Ecommerce) - Sistema completo de comercio electrónico
- [**API REST**](Example-REST-API) - API RESTful con múltiples endpoints
- [**Microservicio**](Example-Microservice) - Microservicio con gRPC
- [**CLI Tool**](Example-CLI-Tool) - Herramienta de línea de comandos

### 🔧 Avanzado
- [**Personalización**](Customization) - Adaptar plantillas a tus necesidades
- [**Integración CI/CD**](CICD-Integration) - Automatización y despliegue
- [**Testing**](Testing-Guide) - Estrategias de testing con código generado
- [**Troubleshooting**](Troubleshooting) - Solución de problemas comunes

### 🤝 Contribución
- [**Guía de Contribución**](Contributing) - Cómo contribuir al proyecto
- [**Desarrollo**](Development) - Configurar entorno de desarrollo
- [**Roadmap**](Roadmap) - Funcionalidades futuras

## 🎯 ¿Qué es Clean Architecture?

Clean Architecture es un patrón arquitectónico creado por Robert C. Martin (Uncle Bob) que organiza el código en capas concéntricas, donde las dependencias apuntan hacia el centro del sistema. Esto garantiza:

- ✅ **Independencia de frameworks** - El código de negocio no depende de librerías externas
- ✅ **Testabilidad** - Fácil de probar sin dependencias externas
- ✅ **Independencia de UI** - La interfaz puede cambiar sin afectar la lógica
- ✅ **Independencia de base de datos** - La persistencia es un detalle de implementación
- ✅ **Independencia de agentes externos** - El código de negocio no conoce el mundo exterior

## 🚀 Inicio Rápido

### 1. Instalación
```bash
go install github.com/sazardev/goca@latest
```

### 2. Crear tu primer proyecto
```bash
goca init mi-proyecto --module github.com/usuario/mi-proyecto
cd mi-proyecto
```

### 3. Generar tu primer feature
```bash
goca feature User --fields "name:string,email:string" --validation
```

### 4. Configurar dependencias
```bash
goca di --features "User" --database postgres
```

¡Ya tienes un proyecto completo con Clean Architecture! 🎉

## 📈 Filosofía de Goca

Goca no solo genera código, sino que **enseña** y **hace cumplir** las mejores prácticas de Clean Architecture:

- **🟡 Dominio** → Entidades puras sin dependencias externas
- **🔴 Casos de Uso** → Lógica de aplicación con DTOs claros
- **🟢 Adaptadores** → Interfaces que adaptan entrada/salida
- **🔵 Infraestructura** → Implementaciones específicas de tecnología

## 🛡️ Garantías de Calidad

- ✅ **Dependencias dirigidas hacia el centro**
- ✅ **Interfaces claras entre capas**
- ✅ **Separación de responsabilidades**
- ✅ **Código testeable por diseño**
- ✅ **Patrones probados en producción**

## 🌟 Características Destacadas

- **Multi-protocolo**: HTTP REST, gRPC, CLI, Workers, SOAP
- **Multi-base de datos**: PostgreSQL, MySQL, MongoDB
- **Inyección de dependencias**: Manual y con Wire.dev
- **Validaciones**: En dominio y DTOs
- **Testing**: Interfaces para TDD
- **Documentación**: Swagger automático

## 📞 Soporte y Comunidad

- 🐛 **Issues**: [GitHub Issues](https://github.com/sazardev/goca/issues)
- 💬 **Discusiones**: [GitHub Discussions](https://github.com/sazardev/goca/discussions)
- 📧 **Contacto**: [sazardev@email.com](mailto:sazardev@email.com)

---

**¡Explora la documentación y comienza a crear proyectos increíbles con Clean Architecture!** 🚀
