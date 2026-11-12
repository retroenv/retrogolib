// Package jsonschema validates a subset of JSON Schema 2020-12.
//
// Schemas can be objects or booleans. Supported validation keywords are type,
// properties, required, additionalProperties, enum, minimum, maximum, items,
// minItems, maxItems, uniqueItems, and oneOf. Other keywords are rejected,
// except for the annotations below. Schemas can have at most 64 levels.
//
// The annotations $schema, $comment, title, description, default, examples,
// readOnly, writeOnly, deprecated, contentEncoding, contentMediaType, and
// x-mcp-header are accepted. They do not add validation. If present, $schema
// must select the 2020-12 dialect.
//
// Numbers are compared without float64 rounding. Missing fields and explicit
// null values are distinct. Compiled schemas are safe for concurrent use.
package jsonschema
