package jsonschema

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/retroenv/retrogolib/set"
)

func validateSchemaValue(schema, value any, path string) error {
	if allowed, ok := schema.(bool); ok {
		if !allowed {
			return fmt.Errorf("%s is not allowed by the schema", path)
		}

		return nil
	}

	node := schema.(map[string]any)
	if err := validateOneOf(node, value, path); err != nil {
		return err
	}

	if err := validateScalar(node, value, path); err != nil {
		return err
	}

	if err := validateObject(node, value, path); err != nil {
		return err
	}

	return validateArray(node, value, path)
}

func validateScalar(node map[string]any, value any, path string) error {
	if kinds, ok := node[keywordType]; ok && !matchesSchemaTypes(kinds, value) {
		return fmt.Errorf("%s must match type %v", path, kinds)
	}

	if values, ok := node[keywordEnum].([]any); ok && !containsJSONValue(values, value) {
		return fmt.Errorf("%s is not an allowed value", path)
	}

	if number, ok := jsonNumber(value); ok {
		if minimum, ok := jsonNumber(node[keywordMinimum]); ok && number.compare(minimum) < 0 {
			return fmt.Errorf("%s is below the minimum", path)
		}

		if maximum, ok := jsonNumber(node[keywordMaximum]); ok && number.compare(maximum) > 0 {
			return fmt.Errorf("%s is above the maximum", path)
		}
	}

	return nil
}

func validateObject(node map[string]any, value any, path string) error {
	object, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	properties, _ := node[keywordProperties].(map[string]any)
	if required, ok := node[keywordRequired].([]any); ok {
		for _, field := range required {
			name, ok := field.(string)
			if ok {
				if _, exists := object[name]; !exists {
					return fmt.Errorf("%s.%s is required", path, name)
				}
			}
		}
	}

	for name, child := range properties {
		if member, exists := object[name]; exists {
			if err := validateSchemaValue(child, member, path+"."+name); err != nil {
				return err
			}
		}
	}

	return validateAdditionalProperties(node, object, properties, path)
}

func validateAdditionalProperties(node, object, properties map[string]any, path string) error {
	for name, member := range object {
		if _, exists := properties[name]; exists {
			continue
		}

		switch additional := node[keywordAdditionalProperties].(type) {
		case bool:
			if !additional {
				return fmt.Errorf("%s.%s is not allowed", path, name)
			}

		case map[string]any:
			if err := validateSchemaValue(additional, member, path+"."+name); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateArray(node map[string]any, value any, path string) error {
	list, ok := value.([]any)
	if !ok {
		return nil
	}

	length := parseSchemaNumber(json.Number(strconv.Itoa(len(list))))
	if minimum, ok := jsonNumber(node[keywordMinItems]); ok && length.compare(minimum) < 0 {
		return fmt.Errorf("%s has fewer than minItems", path)
	}

	if maximum, ok := jsonNumber(node[keywordMaxItems]); ok && length.compare(maximum) > 0 {
		return fmt.Errorf("%s has more than maxItems", path)
	}

	if unique, _ := node[keywordUniqueItems].(bool); unique && !uniqueJSONValues(list) {
		return fmt.Errorf("%s must have uniqueItems", path)
	}

	if items, ok := node[keywordItems]; ok {
		for index, item := range list {
			if err := validateSchemaValue(items, item, fmt.Sprintf("%s[%d]", path, index)); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateOneOf(node map[string]any, value any, path string) error {
	alternatives, ok := node[keywordOneOf].([]any)
	if !ok {
		return nil
	}

	matches := 0
	for _, child := range alternatives {
		if validateSchemaValue(child, value, path) == nil {
			matches++
		}
	}

	if matches != 1 {
		return fmt.Errorf("%s must match exactly one oneOf branch", path)
	}

	return nil
}

func matchesSchemaTypes(kinds, value any) bool {
	switch kinds := kinds.(type) {
	case string:
		return matchesSchemaType(kinds, value)

	case []any:
		for _, kind := range kinds {
			if matchesSchemaType(kind.(string), value) {
				return true
			}
		}
	}

	return false
}

func knownSchemaType(kind string) bool {
	switch kind {
	case typeObject, typeArray, typeString, typeBoolean, typeInteger, typeNumber, typeNull:
		return true

	default:
		return false
	}
}

func matchesSchemaType(kind string, value any) bool {
	switch kind {
	case typeObject:
		_, ok := value.(map[string]any)

		return ok

	case typeArray:
		_, ok := value.([]any)

		return ok

	case typeString:
		_, ok := value.(string)

		return ok

	case typeBoolean:
		_, ok := value.(bool)

		return ok

	case typeInteger:
		number, ok := jsonNumber(value)

		return ok && number.exponent.Sign() >= 0

	case typeNumber:
		_, ok := jsonNumber(value)

		return ok

	case typeNull:
		return value == nil

	default:
		return false
	}
}

func jsonNumber(value any) (*schemaNumber, bool) {
	number, ok := value.(json.Number)
	if !ok {
		return nil, false
	}

	return parseSchemaNumber(number), true
}

func containsJSONValue(values []any, value any) bool {
	key := canonicalJSON(value)
	for _, candidate := range values {
		if key == canonicalJSON(candidate) {
			return true
		}
	}

	return false
}

func uniqueJSONValues(values []any) bool {
	seen := set.New[string]()
	for _, value := range values {
		key := canonicalJSON(value)
		if seen.Contains(key) {
			return false
		}

		seen.Add(key)
	}

	return true
}
