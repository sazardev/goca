import * as vscode from 'vscode';
import * as child_process from 'child_process';
import * as path from 'path';

export class GocaCommandHandler {
    constructor(private context: vscode.ExtensionContext) { }

    /**
     * Ejecuta un comando de Goca directamente
     */
    public executeCommand(command: string, args: string[] = []): void {
        // Obtener la configuración del usuario
        const gocaPath = vscode.workspace.getConfiguration('goca').get('path', 'goca');

        // Crear terminal para ejecutar el comando
        const terminal = vscode.window.createTerminal('Goca CLI');
        terminal.show();

        // Construir el comando completo
        const fullCommand = `${gocaPath} ${command} ${args.join(' ')}`;
        terminal.sendText(fullCommand);
    }

    /**
     * Inicializa un nuevo proyecto Goca
     */
    public async initProject(): Promise<void> {
        // Solicitar nombre del proyecto
        const projectName = await vscode.window.showInputBox({
            prompt: 'Nombre del proyecto',
            placeHolder: 'mi-proyecto',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre del proyecto es requerido';
            }
        });

        if (!projectName) {
            return; // Usuario canceló
        }

        // Solicitar el módulo
        const moduleName = await vscode.window.showInputBox({
            prompt: 'Nombre del módulo Go',
            placeHolder: 'github.com/usuario/proyecto',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre del módulo es requerido';
            }
        });

        if (!moduleName) {
            return; // Usuario canceló
        }

        // Opciones para base de datos
        const dbOptions = ['postgres', 'mysql', 'sqlite', 'mongodb', 'ninguna'];
        const database = await vscode.window.showQuickPick(dbOptions, {
            placeHolder: 'Selecciona una base de datos'
        });

        if (!database) {
            return; // Usuario canceló
        }

        // Opciones para API
        const apiOptions = ['rest', 'grpc', 'graphql', 'ninguna'];
        const api = await vscode.window.showQuickPick(apiOptions, {
            placeHolder: 'Selecciona un tipo de API'
        });

        if (!api) {
            return; // Usuario canceló
        }

        // Preguntar si incluir autenticación
        const includeAuth = await vscode.window.showQuickPick(['Sí', 'No'], {
            placeHolder: '¿Incluir autenticación?'
        });

        if (!includeAuth) {
            return; // Usuario canceló
        }

        const auth = includeAuth === 'Sí';

        // Construir argumentos
        const args = [
            projectName,
            `--module=${moduleName}`,
            `--database=${database}`,
            `--api=${api}`
        ];

        if (auth) {
            args.push('--auth');
        }

        // Ejecutar comando
        this.executeCommand('init', args);
    }

    /**
     * Genera un feature completo
     */
    public async generateFeature(): Promise<void> {
        const featureName = await vscode.window.showInputBox({
            prompt: 'Nombre del feature',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre del feature es requerido';
            }
        });

        if (!featureName) {
            return; // Usuario canceló
        }

        const fields = await vscode.window.showInputBox({
            prompt: 'Campos (formato: nombre:tipo,edad:int)',
            placeHolder: 'name:string,email:string,age:int',
            validateInput: (input) => {
                return input.trim() ? null : 'Al menos un campo es requerido';
            }
        });

        if (!fields) {
            return; // Usuario canceló
        }

        // Opciones adicionales
        const withValidation = await vscode.window.showQuickPick(['Sí', 'No'], {
            placeHolder: '¿Incluir validaciones?'
        });

        if (!withValidation) {
            return;
        }

        const args = [
            featureName,
            `--fields="${fields}"`
        ];

        if (withValidation === 'Sí') {
            args.push('--validation');
        }

        this.executeCommand('feature', args);
    }

    /**
     * Genera una entidad
     */
    public async generateEntity(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        const fields = await vscode.window.showInputBox({
            prompt: 'Campos (formato: nombre:tipo,edad:int)',
            placeHolder: 'name:string,email:string,age:int',
            validateInput: (input) => {
                return input.trim() ? null : 'Al menos un campo es requerido';
            }
        });

        if (!fields) {
            return; // Usuario canceló
        }

        // Opciones adicionales
        const withValidation = await vscode.window.showQuickPick(['Sí', 'No'], {
            placeHolder: '¿Incluir validaciones?'
        });

        if (!withValidation) {
            return;
        }

        const withTimestamps = await vscode.window.showQuickPick(['Sí', 'No'], {
            placeHolder: '¿Incluir timestamps?'
        });

        if (!withTimestamps) {
            return;
        }

        const args = [
            entityName,
            `--fields="${fields}"`
        ];

        if (withValidation === 'Sí') {
            args.push('--validation');
        }

        if (withTimestamps === 'Sí') {
            args.push('--timestamps');
        }

        this.executeCommand('entity', args);
    }

    /**
     * Genera un caso de uso
     */
    public async generateUseCase(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        const operationOptions = ['create', 'read', 'update', 'delete', 'list', 'search'];
        const operation = await vscode.window.showQuickPick(operationOptions, {
            placeHolder: 'Selecciona una operación'
        });

        if (!operation) {
            return; // Usuario canceló
        }

        const args = [
            entityName,
            `--operation=${operation}`
        ];

        this.executeCommand('usecase', args);
    }

    /**
     * Genera un repositorio
     */
    public async generateRepository(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        const dbOptions = ['postgres', 'mysql', 'sqlite', 'mongodb', 'memory'];
        const database = await vscode.window.showQuickPick(dbOptions, {
            placeHolder: 'Selecciona una base de datos'
        });

        if (!database) {
            return; // Usuario canceló
        }

        const args = [
            entityName,
            `--database=${database}`
        ];

        this.executeCommand('repository', args);
    }

    /**
     * Genera un handler
     */
    public async generateHandler(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        const typeOptions = ['rest', 'grpc', 'graphql'];
        const handlerType = await vscode.window.showQuickPick(typeOptions, {
            placeHolder: 'Selecciona un tipo de handler'
        });

        if (!handlerType) {
            return; // Usuario canceló
        }

        const args = [
            entityName,
            `--type=${handlerType}`
        ];

        this.executeCommand('handler', args);
    }

    /**
     * Genera mensajes
     */
    public async generateMessages(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        this.executeCommand('messages', [entityName]);
    }

    /**
     * Genera inyección de dependencias
     */
    public async generateDI(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        this.executeCommand('di', [entityName]);
    }

    /**
     * Genera interfaces
     */
    public async generateInterfaces(): Promise<void> {
        const entityName = await vscode.window.showInputBox({
            prompt: 'Nombre de la entidad',
            placeHolder: 'User, Product, Order, etc.',
            validateInput: (input) => {
                return input.trim() ? null : 'El nombre de la entidad es requerido';
            }
        });

        if (!entityName) {
            return; // Usuario canceló
        }

        const layerOptions = ['repository', 'usecase', 'handler', 'all'];
        const layer = await vscode.window.showQuickPick(layerOptions, {
            placeHolder: 'Selecciona una capa'
        });

        if (!layer) {
            return; // Usuario canceló
        }

        const args = [
            entityName,
            `--layer=${layer}`
        ];

        this.executeCommand('interfaces', args);
    }
}
