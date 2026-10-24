package mcp

import (
	"encoding/json"
)

// ProtocolVersion is the MCP protocol revision this package implements.
const ProtocolVersion = "2026-07-28"

// Reserved _meta key names.
const (
	metaKeyClientCapabilities = "io.modelcontextprotocol/clientCapabilities"
	metaKeyClientInfo         = "io.modelcontextprotocol/clientInfo"
	metaKeyLogLevel           = "io.modelcontextprotocol/logLevel"
	metaKeyProgressToken      = "progressToken"
	metaKeyProtocolVersion    = "io.modelcontextprotocol/protocolVersion"
	metaKeyServerInfo         = "io.modelcontextprotocol/serverInfo"
)

// Implementation identifies a client or server implementation.
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Meta holds the per-request protocol metadata parsed from params._meta.
type Meta struct {
	ProtocolVersion    string
	ClientInfo         *Implementation
	ClientCapabilities json.RawMessage
	LogLevel           string
	ProgressToken      json.RawMessage
}

// versionData is the data member of an unsupported protocol version error.
type versionData struct {
	Supported []string `json:"supported"`
	Requested string   `json:"requested"`
}

// parseRequestMeta validates the protocol metadata of one request. A legacy
// initialize request gets an error that names the supported versions, because
// a legacy client cannot use this server and shows this message to the user.
func parseRequestMeta(msg message) (Meta, *Error) {
	if msg.Method == methodInitialize {
		return Meta{}, legacyInitializeError(msg.Params)
	}
	return parseMeta(msg.Params)
}

// parseMeta extracts and validates the required _meta fields from request
// params. The protocol version check happens here so every request follows
// the same rule.
func parseMeta(params json.RawMessage) (Meta, *Error) {
	var meta Meta

	var body struct {
		Meta json.RawMessage `json:"_meta"`
	}
	if len(params) == 0 {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: _meta missing",
		}
	}
	if err := json.Unmarshal(params, &body); err != nil {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params",
		}
	}
	if len(body.Meta) == 0 {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: _meta missing",
		}
	}

	return parseMetaFields(body.Meta)
}

// parseMetaFields extracts and validates the required _meta fields.
func parseMetaFields(raw json.RawMessage) (Meta, *Error) {
	var meta Meta

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: _meta malformed",
		}
	}

	versionRaw, ok := fields[metaKeyProtocolVersion]
	if !ok {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: protocolVersion missing",
		}
	}
	if err := json.Unmarshal(versionRaw, &meta.ProtocolVersion); err != nil || meta.ProtocolVersion == "" {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: protocolVersion malformed",
		}
	}
	if meta.ProtocolVersion != ProtocolVersion {
		return meta, unsupportedVersionError(meta.ProtocolVersion)
	}

	capabilitiesRaw, ok := fields[metaKeyClientCapabilities]
	if !ok {
		return meta, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: clientCapabilities missing",
		}
	}
	meta.ClientCapabilities = capabilitiesRaw

	if raw, ok := fields[metaKeyClientInfo]; ok {
		var info Implementation
		if err := json.Unmarshal(raw, &info); err == nil {
			meta.ClientInfo = &info
		}
	}

	if raw, ok := fields[metaKeyLogLevel]; ok {
		_ = json.Unmarshal(raw, &meta.LogLevel)
	}

	meta.ProgressToken = fields[metaKeyProgressToken]

	return meta, nil
}

// unsupportedVersionError builds the protocol version negotiation error.
func unsupportedVersionError(requested string) *Error {
	return &Error{
		Code:    CodeUnsupportedProtocolVersion,
		Message: "Unsupported protocol version",
		Data:    unsupportedVersionData(requested),
	}
}

// legacyInitializeError builds the error for a legacy initialize request. The
// message names the supported version, so that a legacy client can show a
// useful diagnostic.
func legacyInitializeError(params json.RawMessage) *Error {
	var body struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	_ = json.Unmarshal(params, &body)

	return &Error{
		Code: CodeUnsupportedProtocolVersion,
		Message: "Unsupported protocol version: this server supports only MCP " + ProtocolVersion +
			" and does not implement the initialize handshake",
		Data: unsupportedVersionData(body.ProtocolVersion),
	}
}

// unsupportedVersionData builds the data member of a version error.
func unsupportedVersionData(requested string) versionData {
	return versionData{
		Supported: []string{ProtocolVersion},
		Requested: requested,
	}
}
