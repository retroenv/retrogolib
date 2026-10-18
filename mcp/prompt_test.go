package mcp

import (
	"context"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestPromptListAndGet(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddPrompt(Prompt{
		Name:        "greet",
		Description: "Greet a person.",
		Arguments: []PromptArgument{
			{Name: "name", Description: "Person name.", Required: true},
		},
	}, func(_ context.Context, args map[string]string) (*PromptResult, error) {
		return &PromptResult{
			Description: "Greeting prompt.",
			Messages:    []PromptMessage{{Role: "user", Content: Text("Hello, " + args["name"] + "!")}},
		}, nil
	}))

	input := request(`1`, methodPromptsList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","prompts":[{"name":"greet",` +
		`"description":"Greet a person.","arguments":[{"name":"name","description":"Person name.",` +
		`"required":true}]}],"ttlMs":3600000,"cacheScope":"public","_meta":` +
		`{"io.modelcontextprotocol/serverInfo":{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)

	input = request(`2`, methodPromptsGet, metaParams(`"name":"greet","arguments":{"name":"Ada"}`))
	lines = runServe(t, s, input+"\n")

	want = `{"jsonrpc":"2.0","id":2,"result":{"resultType":"complete","description":"Greeting prompt.",` +
		`"messages":[{"role":"user","content":{"type":"text","text":"Hello, Ada!"}}],"_meta":` +
		`{"io.modelcontextprotocol/serverInfo":{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestPromptMissingRequiredArgument(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddPrompt(Prompt{
		Name:      "greet",
		Arguments: []PromptArgument{{Name: "name", Required: true}},
	}, promptHandler))

	input := request(`1`, methodPromptsGet, metaParams(`"name":"greet","arguments":{}`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Invalid params: missing required argument: name"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestPromptUnknown(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, methodPromptsGet, metaParams(`"name":"nope"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Unknown prompt: nope"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestPromptHandlerError(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddPrompt(Prompt{Name: "broken"}, func(_ context.Context, _ map[string]string) (*PromptResult, error) {
		return nil, errStorage
	}))

	input := request(`1`, methodPromptsGet, metaParams(`"name":"broken"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error","data":{"error":"storage broken"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestAddPromptDuplicate(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddPrompt(Prompt{Name: "p"}, promptHandler))
	assert.ErrorIs(t, s.AddPrompt(Prompt{Name: "p"}, promptHandler), ErrDuplicatePrompt)
}
