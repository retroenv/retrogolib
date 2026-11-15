package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Resource describes a readable server resource.
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourceTemplate describes a URI template that matches many resources.
// Templates support {name} segments only (RFC 6570 level 1).
type ResourceTemplate struct {
	URITemplate string `json:"uriTemplate"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourceHandler reads the contents of one resource. It returns
// ErrResourceNotFound when no resource exists for the URI.
type ResourceHandler func(ctx context.Context, uri string) ([]ResourceContents, error)

// ResourceContents is one content item of a resources/read result or of an
// embedded resource. A non-nil Blob makes it binary contents; otherwise it is
// text contents.
type ResourceContents struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     []byte `json:"blob,omitempty"`
}

// AddResource registers a resource with its read handler.
func (s *Server) AddResource(r Resource, h ResourceHandler) error {
	if r.URI == "" {
		return ErrEmptyName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.resources {
		if entry.resource.URI == r.URI {
			return ErrDuplicateResource
		}
	}

	s.resources = append(s.resources, resourceEntry{
		handler:  h,
		resource: r,
	})

	return nil
}

// AddResourceTemplate registers a URI template with its read handler.
// Templates support {name} segments only (RFC 6570 level 1).
func (s *Server) AddResourceTemplate(t ResourceTemplate, h ResourceHandler) error {
	if t.URITemplate == "" {
		return ErrEmptyName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.templates {
		if entry.template.URITemplate == t.URITemplate {
			return ErrDuplicateResource
		}
	}

	s.templates = append(s.templates, templateEntry{
		handler:  h,
		template: t,
	})

	return nil
}

// MarshalJSON encodes the contents as text contents or as blob contents. The
// text field stays in the output when it is empty, because the protocol
// schema requires it.
func (r ResourceContents) MarshalJSON() ([]byte, error) {
	var v any
	if r.Blob != nil {
		v = struct {
			URI      string `json:"uri"`
			MimeType string `json:"mimeType,omitempty"`
			Blob     []byte `json:"blob"`
		}{
			URI:      r.URI,
			MimeType: r.MimeType,
			Blob:     r.Blob,
		}
	} else {
		v = struct {
			URI      string `json:"uri"`
			MimeType string `json:"mimeType,omitempty"`
			Text     string `json:"text"`
		}{
			URI:      r.URI,
			MimeType: r.MimeType,
			Text:     r.Text,
		}
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encoding resource contents: %w", err)
	}
	return data, nil
}

// resourceListResult builds the resources/list result in registration order.
func (s *Server) resourceListResult() *resourceListResult {
	s.mu.RLock()
	resources := make([]Resource, 0, len(s.resources))
	for _, entry := range s.resources {
		resources = append(resources, entry.resource)
	}

	s.mu.RUnlock()

	return &resourceListResult{
		ResultType: resultTypeComplete,
		Resources:  resources,
		TTLMs:      s.opts.ListTTL.Milliseconds(),
		CacheScope: cacheScopePublic,
		Meta:       s.serverInfo,
	}
}

// templateListResult builds the resources/templates/list result in registration order.
func (s *Server) templateListResult() *templateListResult {
	s.mu.RLock()
	templates := make([]ResourceTemplate, 0, len(s.templates))
	for _, entry := range s.templates {
		templates = append(templates, entry.template)
	}

	s.mu.RUnlock()

	return &templateListResult{
		ResultType:        resultTypeComplete,
		ResourceTemplates: templates,
		TTLMs:             s.opts.ListTTL.Milliseconds(),
		CacheScope:        cacheScopePublic,
		Meta:              s.serverInfo,
	}
}

// readResourceResult reads one resource by URI or template match.
func (s *Server) readResourceResult(ctx context.Context, params json.RawMessage) (any, *Error) {
	var body struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &body); err != nil {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params",
		}
	}
	if body.URI == "" {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: uri missing",
		}
	}

	handler, ok := s.findResource(body.URI)
	if !ok {
		return nil, resourceNotFoundError(body.URI)
	}

	contents, err := handler(ctx, body.URI)
	if errors.Is(err, ErrResourceNotFound) {
		return nil, resourceNotFoundError(body.URI)
	}
	if err != nil {
		return nil, &Error{
			Code:    CodeInternal,
			Message: "Internal error",
			Data:    map[string]any{"error": err.Error()},
		}
	}
	if contents == nil {
		contents = []ResourceContents{}
	}

	return &readResourceResult{
		ResultType: resultTypeComplete,
		Contents:   contents,
		TTLMs:      s.opts.ListTTL.Milliseconds(),
		CacheScope: cacheScopePublic,
		Meta:       s.serverInfo,
	}, nil
}

// findResource locates the handler for a resource URI, checking static
// resources first and then templates in registration order.
func (s *Server) findResource(uri string) (ResourceHandler, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entry := range s.resources {
		if entry.resource.URI == uri {
			return entry.handler, true
		}
	}

	for _, entry := range s.templates {
		if matchURITemplate(entry.template.URITemplate, uri) {
			return entry.handler, true
		}
	}

	return nil, false
}

type resourceEntry struct {
	handler  ResourceHandler
	resource Resource
}

type templateEntry struct {
	handler  ResourceHandler
	template ResourceTemplate
}

// resourceListResult is the wire form of a resources/list result.
type resourceListResult struct {
	ResultType string          `json:"resultType"`
	Resources  []Resource      `json:"resources"`
	TTLMs      int64           `json:"ttlMs"`
	CacheScope string          `json:"cacheScope"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// templateListResult is the wire form of a resources/templates/list result.
type templateListResult struct {
	ResultType        string             `json:"resultType"`
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
	TTLMs             int64              `json:"ttlMs"`
	CacheScope        string             `json:"cacheScope"`
	Meta              json.RawMessage    `json:"_meta,omitempty"`
}

// readResourceResult is the wire form of a resources/read result.
type readResourceResult struct {
	ResultType string             `json:"resultType"`
	Contents   []ResourceContents `json:"contents"`
	TTLMs      int64              `json:"ttlMs"`
	CacheScope string             `json:"cacheScope"`
	Meta       json.RawMessage    `json:"_meta,omitempty"`
}

// resourceNotFoundError builds the protocol error for a missing resource.
func resourceNotFoundError(uri string) *Error {
	return &Error{
		Code:    CodeInvalidParams,
		Message: "Resource not found",
		Data:    map[string]any{"uri": uri},
	}
}

// matchURITemplate reports whether a URI matches an RFC 6570 level 1
// template. Variable segments in {braces} match any segment that is not empty.
func matchURITemplate(template, uri string) bool {
	templateParts := strings.Split(template, "/")
	uriParts := strings.Split(uri, "/")
	if len(templateParts) != len(uriParts) {
		return false
	}

	for i, part := range templateParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			if uriParts[i] == "" {
				return false
			}
			continue
		}
		if part != uriParts[i] {
			return false
		}
	}

	return true
}
