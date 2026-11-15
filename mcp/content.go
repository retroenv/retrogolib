package mcp

import (
	"encoding/json"
	"fmt"
	"slices"
)

// Content types.
const (
	contentTypeAudio    = "audio"
	contentTypeImage    = "image"
	contentTypeResource = "resource"
	contentTypeText     = "text"
)

// Content is one content item in a tool result or prompt message.
type Content struct {
	Type     string            `json:"type"`
	Text     string            `json:"text,omitempty"`
	Data     []byte            `json:"data,omitempty"`
	MimeType string            `json:"mimeType,omitempty"`
	Resource *ResourceContents `json:"resource,omitempty"`
}

// MarshalJSON encodes the content item with the fields that its type
// requires. An empty text or empty data stays in the output, because the
// protocol schema requires these fields.
func (c Content) MarshalJSON() ([]byte, error) {
	var v any

	switch c.Type {
	case contentTypeText:
		v = struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			Type: c.Type,
			Text: c.Text,
		}

	case contentTypeImage, contentTypeAudio:
		data := c.Data
		if data == nil {
			data = []byte{}
		}

		v = struct {
			Type     string `json:"type"`
			Data     []byte `json:"data"`
			MimeType string `json:"mimeType"`
		}{
			Type:     c.Type,
			Data:     data,
			MimeType: c.MimeType,
		}

	case contentTypeResource:
		v = struct {
			Type     string            `json:"type"`
			Resource *ResourceContents `json:"resource"`
		}{
			Type:     c.Type,
			Resource: c.Resource,
		}

	default:
		type plain Content
		v = plain(c)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encoding content: %w", err)
	}
	return data, nil
}

func (c Content) clone() Content {
	c.Data = slices.Clone(c.Data)
	if c.Resource != nil {
		resource := *c.Resource
		resource.Blob = slices.Clone(resource.Blob)
		c.Resource = &resource
	}

	return c
}

// Text builds a text content item.
func Text(s string) Content {
	return Content{
		Type: contentTypeText,
		Text: s,
	}
}

// Image builds an image content item from raw image bytes.
func Image(data []byte, mimeType string) Content {
	return Content{
		Type:     contentTypeImage,
		Data:     data,
		MimeType: mimeType,
	}
}

// Audio builds an audio content item from raw audio bytes.
func Audio(data []byte, mimeType string) Content {
	return Content{
		Type:     contentTypeAudio,
		Data:     data,
		MimeType: mimeType,
	}
}

// EmbeddedText builds an embedded text resource content item.
func EmbeddedText(uri, mimeType, text string) Content {
	return Content{
		Type: contentTypeResource,
		Resource: &ResourceContents{
			URI:      uri,
			MimeType: mimeType,
			Text:     text,
		},
	}
}
