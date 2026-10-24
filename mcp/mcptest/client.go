// Package mcptest contains an in-process MCP client for testing servers.
package mcptest

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/retroenv/retrogolib/mcp"
)

// Method names used by the client.
const (
	methodDiscover      = "server/discover"
	methodPromptsGet    = "prompts/get"
	methodResourcesRead = "resources/read"
	methodToolsCall     = "tools/call"
	methodToolsList     = "tools/list"
)

// DiscoverResult is the parsed result of server/discover.
type DiscoverResult struct {
	ResultType        string                     `json:"resultType"`
	SupportedVersions []string                   `json:"supportedVersions"`
	Capabilities      map[string]json.RawMessage `json:"capabilities"`
	Instructions      string                     `json:"instructions"`
}

type response struct {
	result json.RawMessage
	rpcErr *mcp.Error
}

// Client is an in-process MCP client connected to a server over io.Pipe. Its
// methods are safe for concurrent use.
type Client struct {
	stdin  *io.PipeWriter
	nextID atomic.Int64
	cancel context.CancelFunc
	done   chan struct{}

	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[string]chan response
	closeOnce sync.Once
}

// New connects a client to the server over an in-process pipe.
func New(ctx context.Context, server *mcp.Server) *Client {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	serveCtx, cancel := context.WithCancel(ctx)
	client := &Client{
		stdin:   stdinWriter,
		cancel:  cancel,
		done:    make(chan struct{}),
		pending: map[string]chan response{},
	}

	go func() {
		defer cancel()
		defer close(client.done)
		defer func() { _ = stdoutWriter.Close() }()
		_ = server.Serve(serveCtx, stdinReader, stdoutWriter)
	}()
	go client.readResponses(stdoutReader)

	return client
}

// Close shuts down the client and waits for the server to stop.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		_ = c.stdin.Close()
		<-c.done
	})
}

// Discover calls server/discover.
func (c *Client) Discover(ctx context.Context) (*DiscoverResult, error) {
	raw, err := c.request(ctx, methodDiscover, nil)
	if err != nil {
		return nil, err
	}

	var result DiscoverResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decoding discover result: %w", err)
	}
	return &result, nil
}

// ListTools calls tools/list.
func (c *Client) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	raw, err := c.request(ctx, methodToolsList, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tools []mcp.Tool `json:"tools"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decoding tool list: %w", err)
	}
	return result.Tools, nil
}

// CallTool calls tools/call with the given arguments.
func (c *Client) CallTool(ctx context.Context, name string, args any) (*mcp.ToolResult, error) {
	params := map[string]any{"name": name}
	if args != nil {
		params["arguments"] = args
	}

	raw, err := c.request(ctx, methodToolsCall, params)
	if err != nil {
		return nil, err
	}

	var result mcp.ToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decoding tool result: %w", err)
	}
	return &result, nil
}

// CallToolInto calls tools/call and decodes the structured content into out.
func (c *Client) CallToolInto(ctx context.Context, name string, args, out any) error {
	result, err := c.CallTool(ctx, name, args)
	if err != nil {
		return err
	}
	if result.IsError {
		text := ""
		if len(result.Content) > 0 {
			text = result.Content[0].Text
		}

		return fmt.Errorf("tool execution failed: %s", text)
	}
	if result.StructuredContent == nil {
		return errors.New("tool returned no structured content")
	}

	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return fmt.Errorf("encoding structured content: %w", err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decoding structured content: %w", err)
	}
	return nil
}

// ReadResource calls resources/read.
func (c *Client) ReadResource(ctx context.Context, uri string) ([]mcp.ResourceContents, error) {
	raw, err := c.request(ctx, methodResourcesRead, map[string]any{"uri": uri})
	if err != nil {
		return nil, err
	}

	var result struct {
		Contents []mcp.ResourceContents `json:"contents"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decoding resource result: %w", err)
	}
	return result.Contents, nil
}

// GetPrompt calls prompts/get.
func (c *Client) GetPrompt(ctx context.Context, name string, args map[string]string) (*mcp.PromptResult, error) {
	params := map[string]any{"name": name}
	if args != nil {
		params["arguments"] = args
	}

	raw, err := c.request(ctx, methodPromptsGet, params)
	if err != nil {
		return nil, err
	}

	var result mcp.PromptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decoding prompt result: %w", err)
	}
	return &result, nil
}

// request sends a request and waits for its response.
func (c *Client) request(ctx context.Context, method string, params any) (json.RawMessage, error) {
	id := strconv.FormatInt(c.nextID.Add(1), 10)

	paramsRaw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("encoding params: %w", err)
	}

	ch := make(chan response, 1)
	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	envelope := map[string]any{
		"jsonrpc": "2.0",
		"id":      json.Number(id),
		"method":  method,
		"params":  withMeta(paramsRaw),
	}
	if err := c.send(envelope); err != nil {
		return nil, err
	}

	select {
	case res := <-ch:
		if res.rpcErr != nil {
			return nil, res.rpcErr
		}
		return res.result, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("waiting for response: %w", ctx.Err())

	case <-c.done:
		return nil, errors.New("server closed the connection")
	}
}

// send writes one message line to the server.
func (c *Client) send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("writing request: %w", err)
	}
	return nil
}

// readResponses reads response lines and routes them to the matching request.
func (c *Client) readResponses(r *io.PipeReader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return // EOF: the server stopped
		}

		var envelope struct {
			ID     json.RawMessage `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  *mcp.Error      `json:"error"`
		}
		if err := json.Unmarshal(line, &envelope); err != nil {
			continue
		}

		c.pendingMu.Lock()
		ch, ok := c.pending[string(envelope.ID)]
		c.pendingMu.Unlock()
		if !ok {
			continue // notification or stale response
		}

		ch <- response{
			result: envelope.Result,
			rpcErr: envelope.Error,
		}
	}
}

// meta is the protocol metadata attached to every request.
var meta = func() json.RawMessage {
	data, err := json.Marshal(map[string]any{
		"io.modelcontextprotocol/protocolVersion": mcp.ProtocolVersion,
		"io.modelcontextprotocol/clientInfo": mcp.Implementation{
			Name:    "mcptest",
			Version: "0.0.0",
		},
		"io.modelcontextprotocol/clientCapabilities": map[string]any{},
	})
	if err != nil {
		panic(err)
	}
	return data
}()

// withMeta injects the required protocol metadata into request params.
func withMeta(params json.RawMessage) json.RawMessage {
	merged := map[string]json.RawMessage{}
	if len(params) > 0 && string(params) != "null" {
		_ = json.Unmarshal(params, &merged)
	}

	merged["_meta"] = meta

	data, _ := json.Marshal(merged)

	return data
}
