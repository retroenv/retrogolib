package mcptest

import (
	"context"
	"sync"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
	"github.com/retroenv/retrogolib/mcp"
)

type addArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func newTestServer(t *testing.T) *mcp.Server {
	t.Helper()

	server := mcp.New(log.NewTestLogger(t), mcp.Options{
		Name:    "test-server",
		Version: "1.2.3",
	})

	assert.NoError(t, mcp.AddTool(server, mcp.Tool{
		Name:        "add",
		Description: "Add two integers.",
	},
		func(_ context.Context, _ mcp.ToolRequest, in addArgs) (int, error) {
			return in.A + in.B, nil
		}))
	assert.NoError(t, server.AddResource(mcp.Resource{
		URI:      "mem://data",
		Name:     "data",
		MimeType: "text/plain",
	},
		func(_ context.Context, uri string) ([]mcp.ResourceContents, error) {
			return []mcp.ResourceContents{{
				URI:  uri,
				Text: "contents",
			}}, nil
		}))
	assert.NoError(t, server.AddPrompt(mcp.Prompt{
		Name:        "greet",
		Description: "Greet a person.",
	},
		func(_ context.Context, args map[string]string) (*mcp.PromptResult, error) {
			return &mcp.PromptResult{
				Messages: []mcp.PromptMessage{{
					Role:    "user",
					Content: mcp.Text("Hello, " + args["name"] + "!"),
				}},
			}, nil
		}))

	return server
}

func TestClientDiscover(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	result, err := client.Discover(t.Context())
	assert.NoError(t, err)

	assert.Equal(t, "complete", result.ResultType)
	assert.Equal(t, []string{"2026-07-28"}, result.SupportedVersions)
	_, ok := result.Capabilities["tools"]
	assert.True(t, ok)
	_, ok = result.Capabilities["resources"]
	assert.True(t, ok)
	_, ok = result.Capabilities["prompts"]
	assert.True(t, ok)
}

func TestClientListTools(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	tools, err := client.ListTools(t.Context())
	assert.NoError(t, err)

	assert.Len(t, tools, 1)
	assert.Equal(t, "add", tools[0].Name)
}

func TestClientCallToolInto(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	var sum int
	assert.NoError(t, client.CallToolInto(t.Context(), "add", addArgs{
		A: 2,
		B: 3,
	}, &sum))
	assert.Equal(t, 5, sum)
}

func TestClientCallToolErrorResult(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	result, err := client.CallTool(t.Context(), "add", map[string]any{"a": 1, "b": "text"})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestClientCallToolProtocolError(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	_, err := client.CallTool(t.Context(), "missing", nil)

	var rpcErr *mcp.Error
	assert.ErrorAs(t, err, &rpcErr)
	assert.Equal(t, mcp.CodeInvalidParams, rpcErr.Code)
	assert.Contains(t, rpcErr.Message, "Unknown tool: missing")
}

func TestClientReadResource(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	contents, err := client.ReadResource(t.Context(), "mem://data")
	assert.NoError(t, err)

	assert.Len(t, contents, 1)
	assert.Equal(t, "mem://data", contents[0].URI)
	assert.Equal(t, "contents", contents[0].Text)
}

func TestClientGetPrompt(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	result, err := client.GetPrompt(t.Context(), "greet", map[string]string{"name": "Ada"})
	assert.NoError(t, err)

	assert.Len(t, result.Messages, 1)
	assert.Equal(t, "user", result.Messages[0].Role)
	assert.Equal(t, "Hello, Ada!", result.Messages[0].Content.Text)
}

func TestClientConcurrentCalls(t *testing.T) {
	client := New(t.Context(), newTestServer(t))
	defer client.Close()

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			var sum int
			assert.NoError(t, client.CallToolInto(t.Context(), "add", addArgs{
				A: i,
				B: 1,
			}, &sum))
			assert.Equal(t, i+1, sum)
		})
	}
	wg.Wait()
}
