import * as vscode from 'vscode';
import { GocaCommandHandler } from './command-handler';
import { GocaExplorer } from './explorer';
import { GocaDocumentation } from './documentation';

export function activate(context: vscode.ExtensionContext) {
    console.log('Goca extension is now active!');

    // Crear instancias de nuestros servicios
    const commandHandler = new GocaCommandHandler(context);
    const explorer = new GocaExplorer(context);
    const documentation = new GocaDocumentation(context);

    // Registrar comandos
    context.subscriptions.push(
        vscode.commands.registerCommand('goca.version', () => commandHandler.executeCommand('version')),
        vscode.commands.registerCommand('goca.init', () => commandHandler.initProject()),
        vscode.commands.registerCommand('goca.feature', () => commandHandler.generateFeature()),
        vscode.commands.registerCommand('goca.entity', () => commandHandler.generateEntity()),
        vscode.commands.registerCommand('goca.usecase', () => commandHandler.generateUseCase()),
        vscode.commands.registerCommand('goca.repository', () => commandHandler.generateRepository()),
        vscode.commands.registerCommand('goca.handler', () => commandHandler.generateHandler()),
        vscode.commands.registerCommand('goca.messages', () => commandHandler.generateMessages()),
        vscode.commands.registerCommand('goca.di', () => commandHandler.generateDI()),
        vscode.commands.registerCommand('goca.interfaces', () => commandHandler.generateInterfaces()),
        vscode.commands.registerCommand('goca.refreshExplorer', () => explorer.refresh()),
        vscode.commands.registerCommand('goca.showDocumentation', () => documentation.showDocumentation())
    );

    // Inicializar vistas
    explorer.initialize();
    documentation.initialize();
}

export function deactivate() {
    console.log('Goca extension is now deactivated!');
}
