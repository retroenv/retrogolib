package mcp

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
)

// EOF must cancel a request before Serve waits for that request to finish.
func TestEOFCancelsInFlightRequest(t *testing.T) {
	server := newTestServer(t)
	server.opts.CancelOnEOF = true
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	err := server.AddTool(Tool{
		Name:        "wait",
		Description: "Wait for cancellation.",
		InputSchema: []byte(`{"type":"object"}`),
	},
		func(ctx context.Context, _ ToolRequest) (*ToolResult, error) {
			close(entered)
			<-ctx.Done()
			close(cancelled)
			return Errorf("cancelled"), nil
		})
	assert.NoError(t, err)
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	defer func() { _ = writer.Close() }()
	var output bytes.Buffer
	done := make(chan error, 1)
	go func() { done <- server.Serve(t.Context(), reader, &output) }()
	_, err = io.WriteString(writer, request("1", "tools/call", metaParams(`"name":"wait","arguments":{}`))+"\n")
	assert.NoError(t, err)
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("request did not start")
	}
	assert.NoError(t, writer.Close())
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("EOF did not cancel the request")
	}
	<-cancelled
}
