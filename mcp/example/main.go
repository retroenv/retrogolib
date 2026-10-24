package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retroenv/retrogolib/log"
	"github.com/retroenv/retrogolib/mcp"
)

type greetArgs struct {
	Name string `json:"name" description:"Name of the person to greet."`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	srv := mcp.New(log.NewNop(), mcp.Options{
		Name:    "retrogolib-example",
		Version: "0.1.0",
	})

	err := mcp.AddTool(srv, mcp.Tool{
		Name:        "greet",
		Description: "Greet a person by name.",
	}, func(_ context.Context, _ mcp.ToolRequest, in greetArgs) (string, error) {
		return "Hello, " + in.Name + "!", nil
	})
	if err != nil {
		return fmt.Errorf("registering greet tool: %w", err)
	}

	if err := srv.Serve(context.Background(), os.Stdin, os.Stdout); err != nil {
		return fmt.Errorf("serving MCP requests: %w", err)
	}

	return nil
}
