import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';

/**
 * Clase para el ítem de documentación
 */
class DocItem extends vscode.TreeItem {
    constructor(
        public readonly label: string,
        public readonly filePath: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(label, collapsibleState);

        this.tooltip = label;
        this.description = '';

        this.command = {
            command: 'goca.openDocumentation',
            title: 'Abrir documentación',
            arguments: [filePath]
        };

        this.iconPath = new vscode.ThemeIcon('book');
    }
}

/**
 * Proveedor de datos para la documentación
 */
class DocProvider implements vscode.TreeDataProvider<DocItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<DocItem | undefined | null | void> = new vscode.EventEmitter<DocItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<DocItem | undefined | null | void> = this._onDidChangeTreeData.event;

    private docsPath: string;

    constructor(private context: vscode.ExtensionContext) {
        // La documentación estará en la carpeta docs de la extensión
        this.docsPath = path.join(context.extensionPath, 'docs');
    }

    refresh(): void {
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: DocItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: DocItem): Thenable<DocItem[]> {
        if (element) {
            return Promise.resolve([]);
        }

        // Lista de documentación embebida
        return Promise.resolve([
            new DocItem('Introducción', 'intro', vscode.TreeItemCollapsibleState.None),
            new DocItem('Guía de inicio rápido', 'quickstart', vscode.TreeItemCollapsibleState.None),
            new DocItem('Estructura del proyecto', 'structure', vscode.TreeItemCollapsibleState.None),
            new DocItem('Comandos', 'commands', vscode.TreeItemCollapsibleState.Collapsed),
            new DocItem('Tutoriales', 'tutorials', vscode.TreeItemCollapsibleState.Collapsed),
            new DocItem('Clean Architecture', 'clean-architecture', vscode.TreeItemCollapsibleState.None),
            new DocItem('FAQ', 'faq', vscode.TreeItemCollapsibleState.None)
        ]);
    }
}

/**
 * Clase principal para la documentación de Goca
 */
export class GocaDocumentation {
    private docProvider: DocProvider;
    private treeView: vscode.TreeView<DocItem>;
    private docs: Map<string, string> = new Map();

    constructor(private context: vscode.ExtensionContext) {
        this.docProvider = new DocProvider(context);
        this.loadDocumentation();
    }

    private loadDocumentation(): void {
        // Cargar documentación desde los archivos markdown de la wiki
        this.docs.set('intro', `# Goca - Go Clean Architecture Generator

Goca es un generador de código CLI para Go que te ayuda a crear proyectos siguiendo los principios de Clean Architecture.

## Características principales

- Generación rápida de código para todas las capas
- Plantillas personalizables
- Validación integrada
- Soporte para múltiples bases de datos
- Soporte para múltiples tipos de API
- Enfoque en buenas prácticas

## Instalación

\`\`\`bash
go install github.com/usuario/goca@latest
\`\`\`
        `);

        this.docs.set('quickstart', `# Guía de inicio rápido

## Crear un nuevo proyecto

\`\`\`bash
goca init mi-proyecto --module=github.com/usuario/mi-proyecto --database=postgres --api=rest
\`\`\`

## Generar un feature completo

\`\`\`bash
goca feature User --fields="name:string,email:string,age:int" --validation
\`\`\`

## Generar una entidad

\`\`\`bash
goca entity Product --fields="name:string,price:float,description:string" --validation --timestamps
\`\`\`

## Generar un caso de uso

\`\`\`bash
goca usecase User --operation=create
\`\`\`

## Generar un repositorio

\`\`\`bash
goca repository User --database=postgres
\`\`\`

## Generar un handler

\`\`\`bash
goca handler User --type=rest
\`\`\`
        `);

        this.docs.set('structure', `# Estructura del proyecto

Un proyecto generado por Goca sigue la siguiente estructura:

\`\`\`
mi-proyecto/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── constants/
│   ├── di/
│   ├── domain/
│   ├── handler/
│   ├── messages/
│   ├── repository/
│   └── usecase/
├── migrations/
├── pkg/
│   ├── config/
│   └── logger/
├── go.mod
├── Makefile
└── README.md
\`\`\`

## Capas de Clean Architecture

1. **Domain (internal/domain)**: Entidades y reglas de negocio
2. **Use Cases (internal/usecase)**: Lógica de aplicación
3. **Interface Adapters (internal/repository, internal/handler)**: Adaptadores para DB y API
4. **Frameworks & Drivers (pkg)**: Componentes externos
        `);

        this.docs.set('commands', `# Comandos de Goca

## init
Inicializa un nuevo proyecto con estructura Clean Architecture.

\`\`\`bash
goca init <project-name> --module=<module-name> [--database=<db>] [--api=<api>] [--auth]
\`\`\`

## feature
Genera un feature completo con todas las capas.

\`\`\`bash
goca feature <name> --fields="<fields>" [--validation] [--timestamps]
\`\`\`

## entity
Genera entidades de dominio puras.

\`\`\`bash
goca entity <name> --fields="<fields>" [--validation] [--business-rules] [--timestamps] [--soft-delete]
\`\`\`

## usecase
Genera casos de uso con DTOs.

\`\`\`bash
goca usecase <entity> --operation=<operation>
\`\`\`

## repository
Genera repositorios con interfaces.

\`\`\`bash
goca repository <entity> --database=<database>
\`\`\`

## handler
Genera handlers para diferentes protocolos.

\`\`\`bash
goca handler <entity> --type=<type>
\`\`\`

## messages
Genera mensajes y constantes.

\`\`\`bash
goca messages <entity>
\`\`\`

## di
Genera contenedor de inyección de dependencias.

\`\`\`bash
goca di <entity>
\`\`\`

## interfaces
Genera solo interfaces para TDD.

\`\`\`bash
goca interfaces <entity> --layer=<layer>
\`\`\`

## version
Muestra la versión de Goca.

\`\`\`bash
goca version [--short]
\`\`\`
        `);

        this.docs.set('clean-architecture', `# Clean Architecture

La Clean Architecture, propuesta por Robert C. Martin (Uncle Bob), es un enfoque de diseño de software que separa las preocupaciones en capas concéntricas.

## Principios

1. **Independencia de frameworks**: La arquitectura no depende de la existencia de bibliotecas.
2. **Testeabilidad**: Las reglas de negocio pueden ser probadas sin UI, DB, servidor web, etc.
3. **Independencia de UI**: La UI puede cambiar fácilmente sin cambiar el resto del sistema.
4. **Independencia de base de datos**: Las reglas de negocio no están vinculadas a la base de datos.
5. **Independencia de cualquier agencia externa**: Las reglas de negocio no saben nada sobre el mundo exterior.

## Capas

1. **Entities (Entidades)**: Encapsulan las reglas de negocio críticas de la empresa.
2. **Use Cases (Casos de uso)**: Contienen reglas de negocio específicas de la aplicación.
3. **Interface Adapters (Adaptadores de interfaz)**: Convierten datos entre casos de uso y entidades.
4. **Frameworks & Drivers (Frameworks y controladores)**: Componentes que interactúan con el mundo exterior.

## Regla de dependencia

Las dependencias siempre apuntan hacia adentro. Las capas externas pueden depender de las capas internas, pero no al revés.
        `);

        this.docs.set('tutorials', `# Tutoriales

## Tutorial básico: Creación de una API REST

1. Inicializa un nuevo proyecto:
   \`\`\`bash
   goca init api-demo --module=github.com/usuario/api-demo --database=postgres --api=rest
   \`\`\`

2. Genera un feature completo:
   \`\`\`bash
   goca feature User --fields="name:string,email:string,password:string,age:int" --validation --timestamps
   \`\`\`

3. Navega al directorio del proyecto:
   \`\`\`bash
   cd api-demo
   \`\`\`

4. Compila dependencias:
   \`\`\`bash
   go mod tidy
   \`\`\`

5. Ejecuta migraciones:
   \`\`\`bash
   make migrate-up
   \`\`\`

6. Inicia el servidor:
   \`\`\`bash
   make run
   \`\`\`

## Tutorial avanzado: Implementación de autenticación

1. Inicializa un proyecto con autenticación:
   \`\`\`bash
   goca init auth-demo --module=github.com/usuario/auth-demo --database=postgres --api=rest --auth
   \`\`\`

2. Genera el feature de usuario:
   \`\`\`bash
   goca feature User --fields="name:string,email:string,password:string" --validation --timestamps
   \`\`\`

3. Personaliza la lógica de autenticación en \`internal/usecase/auth_usecase.go\`

4. Implementa middlewares de autenticación
        `);

        this.docs.set('faq', `# Preguntas frecuentes (FAQ)

## ¿Qué bases de datos soporta Goca?

Goca soporta las siguientes bases de datos:
- PostgreSQL
- MySQL
- SQLite
- MongoDB

## ¿Qué tipos de API puedo generar?

Puedes generar:
- REST
- gRPC
- GraphQL

## ¿Puedo personalizar las plantillas?

Sí, puedes personalizar las plantillas editando los archivos en la carpeta \`templates\` de tu instalación de Goca.

## ¿Cómo añado validación a mis entidades?

Usa el flag \`--validation\` al generar entidades o features:

\`\`\`bash
goca entity User --fields="name:string,email:string" --validation
\`\`\`

## ¿Cómo genero código para tests?

Puedes generar interfaces para facilitar el TDD:

\`\`\`bash
goca interfaces User --layer=all
\`\`\`

## ¿Goca genera código para autenticación?

Sí, utiliza el flag \`--auth\` al inicializar un proyecto:

\`\`\`bash
goca init mi-proyecto --module=github.com/usuario/mi-proyecto --auth
\`\`\`
        `);
    }

    public initialize(): void {
        this.treeView = vscode.window.createTreeView('gocaDocumentation', {
            treeDataProvider: this.docProvider
        });

        this.context.subscriptions.push(this.treeView);

        // Registrar comando para abrir documentación
        this.context.subscriptions.push(
            vscode.commands.registerCommand('goca.openDocumentation', (docId) => {
                this.openDocumentation(docId);
            })
        );
    }

    public showDocumentation(): void {
        const items = Array.from(this.docs.keys()).map(key => ({
            label: key.charAt(0).toUpperCase() + key.slice(1).replace(/-/g, ' '),
            id: key
        }));

        vscode.window.showQuickPick(items, {
            placeHolder: 'Selecciona documentación'
        }).then(selected => {
            if (selected) {
                this.openDocumentation(selected.id);
            }
        });
    }

    private openDocumentation(docId: string): void {
        const content = this.docs.get(docId);

        if (!content) {
            vscode.window.showErrorMessage(`Documentación '${docId}' no encontrada`);
            return;
        }

        // Crear un archivo temporal con el contenido de la documentación
        const tempFilePath = path.join(this.context.extensionPath, `${docId}.md`);

        fs.writeFileSync(tempFilePath, content);

        // Abrir el archivo en el editor
        vscode.workspace.openTextDocument(tempFilePath).then(doc => {
            vscode.window.showTextDocument(doc);
        });
    }
}
