package mcp_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/retroenv/retrogolib/log"
	"github.com/retroenv/retrogolib/mcp"
)

// ExampleAddTool registers a typed tool and answers one tools/call request
// over the stdio transport.
func ExampleAddTool() {
	srv := mcp.New(log.NewNop(), mcp.Options{
		Name:    "example",
		Version: "1.0.0",
	})

	type addArgs struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	err := mcp.AddTool(srv, mcp.Tool{
		Name:        "add",
		Description: "Add two integers.",
	}, func(_ context.Context, _ mcp.ToolRequest, in addArgs) (int, error) {
		return in.A + in.B, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	input := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"add",` +
		`"arguments":{"a":2,"b":3},"_meta":{"io.modelcontextprotocol/protocolVersion":` +
		`"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}}}` + "\n"

	if err := srv.Serve(context.Background(), strings.NewReader(input), os.Stdout); err != nil {
		fmt.Println(err)
	}
	// Output:
	// {"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","content":[{"type":"text","text":"5"}],"structuredContent":5,"isError":false,"_meta":{"io.modelcontextprotocol/serverInfo":{"name":"example","version":"1.0.0"}}}}
}
