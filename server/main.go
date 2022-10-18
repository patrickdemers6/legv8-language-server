package main

import (
	"context"
	"go.lsp.dev/jsonrpc2"
	"io"
	"os"
	"os/signal"
	"server/languageserver"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	bufStream := jsonrpc2.NewStream(&readWriter{os.Stdin, os.Stdout, nil})
	rootConn := jsonrpc2.NewConn(bufStream)
	server := languageserver.NewServer(rootConn)

	rootConn.Go(ctx, server.Handler)
	<-rootConn.Done()
}

type readWriter struct {
	io.Reader
	io.Writer
	io.Closer
}
