// Package mcp provides a Model Context Protocol (MCP) server toolkit.
// It lets a Go program expose tools, resources, and prompts to AI agents.
// The package uses only the Go standard library and implements protocol
// revision 2026-07-28.
//
// The stdio transport runs with Serve. The client sends one JSON-RPC message
// per line on standard input, and the server answers on standard output.
//
//	srv := mcp.New(logger, mcp.Options{
//		Name:    "example",
//		Version: "1.0.0",
//	})
//
//	type greetArgs struct {
//		Name string `json:"name" description:"Name of the person."`
//	}
//	err := mcp.AddTool(srv, mcp.Tool{
//		Name:        "greet",
//		Description: "Greet a person by name.",
//	}, func(ctx context.Context, req mcp.ToolRequest, in greetArgs) (string, error) {
//		return "Hello, " + in.Name + "!", nil
//	})
//	if err != nil {
//		return err
//	}
//	return srv.Serve(ctx, os.Stdin, os.Stdout)
//
// The generic AddTool generates the input schema from the argument type. The
// Server.AddTool method registers a tool with a raw JSON handler.
//
// The streamable HTTP transport runs with Handler. Both transports share the
// same Server, so one registration serves both.
//
// Every handler must be safe for concurrent use. The server dispatches each
// request on its own goroutine and serializes the writes.
//
// The package does not implement the legacy initialize handshake, protocol
// sessions, ping, logging/setLevel, the HTTP+SSE transport, multi round-trip
// requests, subscriptions, tasks, SSE responses, or authentication. Clients
// that open with initialize cannot use the server. The server answers
// initialize with an UnsupportedProtocolVersion error that names the
// supported version. Servers log to stderr, because stdout carries MCP
// messages only.
package mcp
