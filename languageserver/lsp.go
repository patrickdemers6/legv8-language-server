package languageserver

import (
	"context"
	"encoding/json"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type Server struct {
	conn      jsonrpc2.Conn
	workspace string
	handlers  handlers
}

// handler is a jsonrpc2.Handler with a custom logger.
type handler = func(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error

type handlers map[string]handler

// NewServer creates a new language server.
func NewServer(conn jsonrpc2.Conn) *Server {
	s := &Server{
		conn: conn,
	}
	s.buildHandlers()
	return s
}

func (s *Server) buildHandlers() {
	s.handlers = map[string]handler{
		lsp.MethodInitialize:                     s.handleInitialize,
		lsp.MethodTextDocumentDidOpen:            s.handleDocumentOpen,
		lsp.MethodWorkspaceDidChangeWatchedFiles: s.handleWatchedFileChange,
		lsp.MethodTextDocumentDidChange:          s.handleDocumentChange,
		lsp.MethodTextDocumentDidSave:            s.handleDocumentSave,
	}
}

// Handler handles the client requests.
func (s *Server) Handler(ctx context.Context, reply jsonrpc2.Replier, r jsonrpc2.Request) error {
	if handler, ok := s.handlers[r.Method()]; ok {
		return handler(ctx, reply, r)
	}

	return reply(ctx, nil, jsonrpc2.ErrMethodNotFound)
}

func (s *Server) handleInitialize(
	ctx context.Context,
	reply jsonrpc2.Replier,
	r jsonrpc2.Request,
) error {
	type initParams struct {
		ProcessID int    `json:"processId,omitempty"`
		RootURI   string `json:"rootUri,omitempty"`
	}

	var params initParams
	if err := json.Unmarshal(r.Params(), &params); err != nil {
		return jsonrpc2.ErrInvalidParams
	}

	s.workspace = string(uri.New(params.RootURI).Filename())
	reply(ctx, lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			// if we support `goto` definition.
			DefinitionProvider: false,

			// If we support `hover` info.
			HoverProvider: false,

			TextDocumentSync: lsp.TextDocumentSyncOptions{
				// Send all file content on every change (can be optimized later).
				Change: lsp.TextDocumentSyncKindFull,

				// if we want to be notified about open/close of Terramate files.
				OpenClose: true,
				Save: &lsp.SaveOptions{
					// If we want the file content on save,
					IncludeText: false,
				},
			},
		},
	}, nil)

	s.conn.Notify(ctx, lsp.MethodWindowShowMessage, lsp.ShowMessageParams{
		Message: "connected to legv8",
		Type:    lsp.MessageTypeInfo,
	})

	return nil
}

func (s *Server) handleWatchedFileChange(
	ctx context.Context,
	reply jsonrpc2.Replier,
	r jsonrpc2.Request,
) error {
	var params lsp.DidChangeWatchedFilesParams
	if err := json.Unmarshal(r.Params(), &params); err != nil {
		return err
	}

	diagnose(params.Changes[0].URI, ctx, s)

	return nil
}

func (s *Server) handleDocumentOpen(
	ctx context.Context,
	reply jsonrpc2.Replier,
	r jsonrpc2.Request,
) error {
	var params lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(r.Params(), &params); err != nil {
		return err
	}

	diagnose(params.TextDocument.URI, ctx, s)

	return nil
}

func (s *Server) handleDocumentChange(
	ctx context.Context,
	reply jsonrpc2.Replier,
	r jsonrpc2.Request,
) error {
	var params lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(r.Params(), &params); err != nil {
		return err
	}

	diagnose(params.TextDocument.URI, ctx, s)

	return nil
}

func (s *Server) handleDocumentSave(
	ctx context.Context,
	reply jsonrpc2.Replier,
	r jsonrpc2.Request,
) error {
	var params lsp.DidSaveTextDocumentParams
	if err := json.Unmarshal(r.Params(), &params); err != nil {
		return err
	}

	diagnose(params.TextDocument.URI, ctx, s)

	return nil
}

func diagnose(uri uri.URI, ctx context.Context, server *Server) {
	tokenizedLines := TokenizeFile(uri)
	diagnostics := Parse(tokenizedLines)

	server.conn.Notify(ctx, lsp.MethodTextDocumentPublishDiagnostics, lsp.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: *diagnostics,
	})
}
