import * as vscode from 'vscode';
import * as path from 'path';

/**
 * Clase que representa un comando en el explorador
 */
class CommandTreeItem extends vscode.TreeItem {
    constructor(
        public readonly label: string,
        public readonly description: string,
        public readonly command: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(label, collapsibleState);

        this.tooltip = description;
        this.description = command;

        this.command = {
            command: `goca.${command}`,
            title: `Ejecutar goca ${command}`,
            arguments: []
        };

        this.iconPath = {
            light: path.join(__filename, '..', '..', 'resources', 'light', `${command}.svg`),
            dark: path.join(__filename, '..', '..', 'resources', 'dark', `${command}.svg`)
        };
    }
}

/**
 * Proveedor de datos para el explorador
 */
class GocaCommandProvider implements vscode.TreeDataProvider<CommandTreeItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<CommandTreeItem | undefined | null | void> = new vscode.EventEmitter<CommandTreeItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<CommandTreeItem | undefined | null | void> = this._onDidChangeTreeData.event;

    refresh(): void {
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: CommandTreeItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: CommandTreeItem): Thenable<CommandTreeItem[]> {
        if (element) {
            return Promise.resolve([]);
        }

        return Promise.resolve([
            new CommandTreeItem('Inicializar proyecto', 'Crear nuevo proyecto con Clean Architecture', 'init', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar feature completo', 'Crear feature con todas las capas', 'feature', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar entidad', 'Crear entidad de dominio', 'entity', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar caso de uso', 'Crear caso de uso con lógica de negocio', 'usecase', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar repositorio', 'Crear capa de acceso a datos', 'repository', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar handler', 'Crear handler para API', 'handler', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar mensajes', 'Crear mensajes y constantes', 'messages', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar DI', 'Crear inyección de dependencias', 'di', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Generar interfaces', 'Crear interfaces para TDD', 'interfaces', vscode.TreeItemCollapsibleState.None),
            new CommandTreeItem('Mostrar versión', 'Ver versión actual de Goca', 'version', vscode.TreeItemCollapsibleState.None)
        ]);
    }
}

/**
 * Clase principal para el explorador de Goca
 */
export class GocaExplorer {
    private commandProvider: GocaCommandProvider;
    private treeView: vscode.TreeView<CommandTreeItem>;

    constructor(private context: vscode.ExtensionContext) {
        this.commandProvider = new GocaCommandProvider();
    }

    public initialize(): void {
        this.treeView = vscode.window.createTreeView('gocaExplorer', {
            treeDataProvider: this.commandProvider
        });

        this.context.subscriptions.push(this.treeView);
    }

    public refresh(): void {
        this.commandProvider.refresh();
    }
}
