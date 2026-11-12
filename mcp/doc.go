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
// AddTool generates an input schema when none is supplied and checks arguments
// against that schema. Server.AddTool registers handlers that receive raw JSON.
// Use CompileSchema to check manually supplied schemas and their values.
// See package jsonschema for supported keywords and limits.
//
// The streamable HTTP transport runs with Handler. Both transports share the
// same Server, so one registration serves both.
//
// Handlers must be safe for concurrent use. MaxConcurrentRequests limits
// active requests across both transports.
//
// The package does not implement the legacy initialize handshake, protocol
// sessions, ping, logging/setLevel, the HTTP+SSE transport, multi round-trip
// requests, subscriptions, tasks, SSE responses, or authentication. Clients
// that send initialize receive an UnsupportedProtocolVersion error.
// Send logs to stderr; stdout carries MCP messages only.
package mcp
