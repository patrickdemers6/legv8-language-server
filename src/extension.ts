import { ChildProcess, spawn } from 'mz/child_process'
import * as url from 'url'
import * as vscode from 'vscode'
import {LanguageClient, LanguageClientOptions, RevealOutputChannelOn, StreamInfo } from 'vscode-languageclient/node'

export async function activate(context: vscode.ExtensionContext): Promise<void> {
    const executablePath = __dirname + '/server';

    let client: LanguageClient

    const serverOptions = () =>
        new Promise<ChildProcess | StreamInfo>((resolve, reject) => {
            // vscode.Uri.file()
                const childProcess = spawn(executablePath)
                resolve({writer: childProcess.stdin, reader: childProcess.stdout})
            })
        

    // Options to control the language client
    const clientOptions: LanguageClientOptions = {
        // Register the server for legv8 documents
        documentSelector: [{ scheme: 'file', language: 'legv8' }, { scheme: 'untitled', language: 'legv8' }],
        revealOutputChannelOn: RevealOutputChannelOn.Never,
        uriConverters: {
            // VS Code by default %-encodes even the colon after the drive letter
            // NodeJS handles it much better
            code2Protocol: uri => url.format(url.parse(uri.toString(true))),
            protocol2Code: str => vscode.Uri.parse(str),
        },
        synchronize: {
            // Synchronize the setting section 'legv8' to the server
            configurationSection: 'legv8',
            // Notify the server about changes to legv8 files in the workspace
            fileEvents: [
                vscode.workspace.createFileSystemWatcher('**/*.legv8'),
                vscode.workspace.createFileSystemWatcher('**/*.legv8asm'),
                vscode.workspace.createFileSystemWatcher('**/*.legv8emul'),
        ],
        }
    }
    

    // Create the language client and start the client.
    client = new LanguageClient('Legv8 Language Server', serverOptions, clientOptions)
    client.start()

    // Push the disposable to the context's subscriptions so that the
    // client can be deactivated on extension deactivation
    
    context.subscriptions.push(client)
}