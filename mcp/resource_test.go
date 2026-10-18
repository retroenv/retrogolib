package mcp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// errStorage is the error returned by failing test handlers.
var errStorage = errors.New("storage broken")

func TestResourceListAndRead(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResource(Resource{
		URI:      "mem://config",
		Name:     "config",
		MimeType: "application/json",
	}, func(_ context.Context, uri string) ([]ResourceContents, error) {
		return []ResourceContents{{URI: uri, MimeType: "application/json", Text: `{"a":1}`}}, nil
	}))

	input := request(`1`, methodResourcesList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","resources":[{"uri":"mem://config",` +
		`"name":"config","mimeType":"application/json"}],"ttlMs":3600000,"cacheScope":"public",` +
		`"_meta":{"io.modelcontextprotocol/serverInfo":{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)

	input = request(`2`, methodResourcesRead, metaParams(`"uri":"mem://config"`))
	lines = runServe(t, s, input+"\n")

	want = `{"jsonrpc":"2.0","id":2,"result":{"resultType":"complete","contents":[{"uri":"mem://config",` +
		`"mimeType":"application/json","text":"{\"a\":1}"}],"ttlMs":3600000,"cacheScope":"public",` +
		`"_meta":{"io.modelcontextprotocol/serverInfo":{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceTemplateMatch(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResourceTemplate(ResourceTemplate{
		URITemplate: "cart://{id}/bank",
		Name:        "cart-bank",
	}, func(_ context.Context, uri string) ([]ResourceContents, error) {
		return []ResourceContents{{
			URI:  uri,
			Text: "bank data",
		}}, nil
	}))

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"cart://42/bank"`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"uri":"cart://42/bank"`)

	input = request(`2`, methodResourcesRead, metaParams(`"uri":"cart://42/other"`))
	lines = runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":2,"error":{"code":-32602,"message":"Resource not found",` +
		`"data":{"uri":"cart://42/other"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceTemplateList(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResourceTemplate(ResourceTemplate{
		URITemplate: "cart://{id}/bank",
		Name:        "cart-bank",
	}, readHandler))

	input := request(`1`, methodResourcesTemplatesList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","resourceTemplates":[` +
		`{"uriTemplate":"cart://{id}/bank","name":"cart-bank"}],"ttlMs":3600000,` +
		`"cacheScope":"public","_meta":{"io.modelcontextprotocol/serverInfo":` +
		`{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceNotFound(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"mem://missing"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Resource not found",` +
		`"data":{"uri":"mem://missing"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceHandlerError(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResource(Resource{
		URI:  "mem://broken",
		Name: "broken",
	},
		func(_ context.Context, _ string) ([]ResourceContents, error) {
			return nil, errStorage
		}))

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"mem://broken"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error","data":{"error":"storage broken"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceHandlerNotFound(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResourceTemplate(ResourceTemplate{
		URITemplate: "cart://{id}",
		Name:        "cart",
	}, func(_ context.Context, _ string) ([]ResourceContents, error) {
		return nil, fmt.Errorf("loading cart: %w", ErrResourceNotFound)
	}))

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"cart://7"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Resource not found",` +
		`"data":{"uri":"cart://7"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestResourceHandlerPanic(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResource(Resource{
		URI:  "mem://panic",
		Name: "panic",
	}, func(_ context.Context, _ string) ([]ResourceContents, error) {
		panic("kaboom")
	}))

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"mem://panic"`)) + "\n" +
		request(`2`, methodResourcesList, metaParams("")) + "\n"
	lines := runServe(t, s, input)

	assert.Len(t, lines, 2)
	output := strings.Join(lines, "\n")
	assert.Contains(t, output, `"id":1,"error":{"code":-32603,"message":"Internal error","data":{"error":"kaboom"}}`)
	assert.Contains(t, output, `"id":2,"result"`)
}

func TestResourceEmptyText(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResource(Resource{
		URI:  "mem://empty",
		Name: "empty",
	}, func(_ context.Context, uri string) ([]ResourceContents, error) {
		return []ResourceContents{{URI: uri}}, nil
	}))

	input := request(`1`, methodResourcesRead, metaParams(`"uri":"mem://empty"`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"contents":[{"uri":"mem://empty","text":""}]`)
}

func TestAddResourceDuplicate(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResource(Resource{
		URI:  "mem://a",
		Name: "a",
	}, readHandler))
	assert.ErrorIs(t, s.AddResource(Resource{
		URI:  "mem://a",
		Name: "a",
	}, readHandler), ErrDuplicateResource)
}

func TestAddResourceTemplateDuplicate(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddResourceTemplate(ResourceTemplate{
		URITemplate: "mem://{id}",
		Name:        "a",
	}, readHandler))
	assert.ErrorIs(t, s.AddResourceTemplate(ResourceTemplate{
		URITemplate: "mem://{id}",
		Name:        "a",
	}, readHandler), ErrDuplicateResource)
}

func TestMatchURITemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		uri      string
		want     bool
	}{
		{name: "exact", template: "mem://config", uri: "mem://config", want: true},
		{name: "exact mismatch", template: "mem://config", uri: "mem://other", want: false},
		{name: "variable match", template: "cart://{id}/bank", uri: "cart://12/bank", want: true},
		{name: "variable mismatch", template: "cart://{id}/bank", uri: "cart://12/banks", want: false},
		{name: "segment count", template: "cart://{id}", uri: "cart://12/extra", want: false},
		{name: "empty variable", template: "cart://{id}/bank", uri: "cart:///bank", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, matchURITemplate(tt.template, tt.uri))
		})
	}
}
