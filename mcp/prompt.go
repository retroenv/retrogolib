package mcp

import (
	"context"
	"encoding/json"
)

// Prompt describes a prompt template the server exposes.
type Prompt struct {
	Name        string           `json:"name"`
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument describes one argument of a prompt. The server rejects a
// prompts/get request that does not supply a required argument.
type PromptArgument struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptHandler renders a prompt with the given arguments.
type PromptHandler func(ctx context.Context, args map[string]string) (*PromptResult, error)

// PromptResult is a rendered prompt.
type PromptResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// PromptMessage is one message of a rendered prompt.
type PromptMessage struct {
	Role    string  `json:"role"`
	Content Content `json:"content"`
}

// promptListResult builds the prompts/list result in registration order.
func (s *Server) promptListResult() *promptListResult {
	s.mu.RLock()
	prompts := make([]Prompt, 0, len(s.prompts))
	for _, entry := range s.prompts {
		prompts = append(prompts, entry.prompt)
	}

	s.mu.RUnlock()

	return &promptListResult{
		ResultType: resultTypeComplete,
		Prompts:    prompts,
		TTLMs:      s.opts.ListTTL.Milliseconds(),
		CacheScope: cacheScopePublic,
		Meta:       s.serverInfo,
	}
}

// getPromptResult renders the named prompt with the given arguments.
func (s *Server) getPromptResult(ctx context.Context, params json.RawMessage) (any, *Error) {
	var body struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}
	if err := json.Unmarshal(params, &body); err != nil {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params",
		}
	}
	if body.Name == "" {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: name missing",
		}
	}

	entry, ok := s.findPrompt(body.Name)
	if !ok {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Unknown prompt: " + body.Name,
		}
	}

	for _, arg := range entry.prompt.Arguments {
		if _, ok := body.Arguments[arg.Name]; arg.Required && !ok {
			return nil, &Error{
				Code:    CodeInvalidParams,
				Message: "Invalid params: missing required argument: " + arg.Name,
			}
		}
	}

	result, err := entry.handler(ctx, body.Arguments)
	if err != nil {
		return nil, &Error{
			Code:    CodeInternal,
			Message: "Internal error",
			Data:    map[string]any{"error": err.Error()},
		}
	}
	if result == nil {
		result = &PromptResult{}
	}

	if result.Messages == nil {
		result.Messages = []PromptMessage{}
	}

	return &promptGetResult{
		ResultType:  resultTypeComplete,
		Description: result.Description,
		Messages:    result.Messages,
		Meta:        s.serverInfo,
	}, nil
}

// promptListResult is the wire form of a prompts/list result.
type promptListResult struct {
	ResultType string          `json:"resultType"`
	Prompts    []Prompt        `json:"prompts"`
	TTLMs      int64           `json:"ttlMs"`
	CacheScope string          `json:"cacheScope"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// promptGetResult is the wire form of a prompts/get result.
type promptGetResult struct {
	ResultType  string          `json:"resultType"`
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}
