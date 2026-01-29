package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultPrefix = "default="
	// TagName is the struct tag name used for configuration field mapping.
	TagName = "config"
)

// Unmarshal unmarshalls configuration data into a struct.
func (c *Config) Unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%w: expected pointer to struct, got %T", ErrInvalidStruct, v)
	}

	return c.unmarshalStruct(rv.Elem(), "")
}

// Marshal marshals a struct into the configuration, preserving comments and structure.
func (c *Config) Marshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("%w: expected struct, got %T", ErrInvalidStruct, v)
	}

	return c.marshalStruct(rv, "")
}

// parseTag parses a struct tag and returns section, key, and default value information.
func (c *Config) parseTag(tag, parentSection string) tagInfo {
	// Split tag by comma to separate path from options
	parts := strings.Split(tag, ",")
	path := strings.TrimSpace(parts[0])

	info := tagInfo{}

	// Parse path to get section and key
	if strings.Contains(path, ".") {
		info.Section, info.Key = c.parseDottedPath(path)
	} else {
		info.Section, info.Key = c.parseSimplePath(path, parentSection)
	}

	// Parse options like default=value and required
	for i := 1; i < len(parts); i++ {
		option := strings.TrimSpace(parts[i])
		if strings.HasPrefix(option, defaultPrefix) {
			// Preserve whitespace in default value by working with original part
			originalOption := parts[i]
			trimmed := strings.TrimSpace(originalOption)
			if strings.HasPrefix(trimmed, defaultPrefix) {
				// Find where "default=" ends in the original string
				prefixIndex := strings.Index(originalOption, defaultPrefix)
				if prefixIndex != -1 {
					valueStart := prefixIndex + len(defaultPrefix)
					info.DefaultValue = originalOption[valueStart:]
					info.HasDefault = true
				}
			}
		} else if option == "required" {
			info.Required = true
		}
	}

	return info
}

// parseDottedPath parses a dotted path like "cpu.registers.accumulator" into section and key.
func (c *Config) parseDottedPath(path string) (section, key string) {
	// For deep paths like "cpu.registers.accumulator", we need to split properly
	// The last part is the key, everything before is the section
	lastDot := strings.LastIndex(path, ".")
	if lastDot != -1 {
		return strings.ToLower(path[:lastDot]), strings.ToLower(path[lastDot+1:])
	}

	// Fallback to original logic (should not happen since caller checked Contains)
	pathParts := strings.SplitN(path, ".", 2)
	return strings.ToLower(pathParts[0]), strings.ToLower(pathParts[1])
}

// parseSimplePath parses a simple path (no dots) into section and key.
func (c *Config) parseSimplePath(path, parentSection string) (section, key string) {
	if parentSection != "" {
		return strings.ToLower(parentSection), strings.ToLower(path)
	}
	return "", strings.ToLower(path) // Use empty string for root-level keys
}

// generateFieldTag creates an automatic config tag for fields without explicit tags.
func (c *Config) generateFieldTag(fieldName, parentSection string, isStruct bool) string {
	fieldNameLower := strings.ToLower(fieldName)

	if isStruct {
		// For nested structs, return just the field name (section will be combined with parent in unmarshal/marshal logic)
		return fieldNameLower
	}

	// For simple fields, create section.key format or root-level key
	if parentSection != "" {
		return parentSection + "." + fieldNameLower
	}
	return fieldNameLower // Return just the field name for root-level keys
}

// updateValue updates a configuration value with type conversion.
func (c *Config) updateValue(sectionName, key string, value any) error {
	if c.sections[sectionName] == nil {
		c.sections[sectionName] = make(section)
	}

	configValue, err := c.convertToValue(value)
	if err != nil {
		return fmt.Errorf("converting value for %s.%s: %w", sectionName, key, err)
	}

	c.sections[sectionName][key] = configValue
	return nil
}

// unmarshalNestedStruct handles unmarshalling of nested struct fields.
func (c *Config) unmarshalNestedStruct(field reflect.StructField, fieldValue reflect.Value, tag, parentSection string) error {
	var sectionName string

	// Handle tag parsing for nested structs
	if strings.Contains(tag, ".") {
		// For explicit nested tags like "system.cpu", use the full tag as section name
		sectionName = strings.ToLower(tag)
	} else {
		// For simple tags (automatic mapping), combine with parent section
		if parentSection != "" {
			sectionName = parentSection + "." + strings.ToLower(tag)
		} else {
			sectionName = strings.ToLower(tag)
		}
	}

	// Nested struct - recursively unmarshal with section as parent
	if err := c.unmarshalStruct(fieldValue, sectionName); err != nil {
		return &UnmarshalError{
			Field:   field.Name,
			Section: sectionName,
			Key:     "",
			Err:     err,
		}
	}
	return nil
}

// unmarshalSimpleField handles unmarshalling of simple (non-struct) fields.
func (c *Config) unmarshalSimpleField(field reflect.StructField, fieldValue reflect.Value, tag, parentSection string) error {
	tagInfo := c.parseTag(tag, parentSection)

	// Get value from configuration
	var value value
	var exists bool

	if c.sections[tagInfo.Section] != nil {
		value, exists = c.sections[tagInfo.Section][tagInfo.Key]
	}

	// If key doesn't exist and we have a default value, use it
	if !exists && tagInfo.HasDefault {
		parsedDefault, err := c.parseDefaultValue(tagInfo.DefaultValue, fieldValue.Type())
		if err != nil {
			return &UnmarshalError{
				Field:   field.Name,
				Section: tagInfo.Section,
				Key:     tagInfo.Key,
				Err:     fmt.Errorf("parsing default value: %w", err),
			}
		}
		value = parsedDefault
		exists = true
	}

	// Check if field is required but missing
	if !exists && tagInfo.Required {
		return &UnmarshalError{
			Field:   field.Name,
			Section: tagInfo.Section,
			Key:     tagInfo.Key,
			Err:     fmt.Errorf("%w: %s.%s", ErrRequiredField, tagInfo.Section, tagInfo.Key),
		}
	}

	if !exists {
		return nil // Key doesn't exist and no default, skip
	}

	if err := c.unmarshalField(fieldValue, value); err != nil {
		return &UnmarshalError{
			Field:   field.Name,
			Section: tagInfo.Section,
			Key:     tagInfo.Key,
			Err:     err,
		}
	}
	return nil
}

// unmarshalStruct processes a struct and populates it with configuration values.
func (c *Config) unmarshalStruct(rv reflect.Value, parentSection string) error {
	rt := rv.Type()

	for i := range rv.NumField() {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !fieldValue.CanSet() {
			continue // Skip unexported fields
		}

		tag := field.Tag.Get(TagName)
		if tag == "-" {
			continue // Skip fields explicitly marked to ignore
		}

		// Generate automatic tag if none provided
		if tag == "" {
			tag = c.generateFieldTag(field.Name, parentSection, fieldValue.Kind() == reflect.Struct)
		}

		if fieldValue.Kind() == reflect.Struct {
			if err := c.unmarshalNestedStruct(field, fieldValue, tag, parentSection); err != nil {
				return err
			}
		} else {
			if err := c.unmarshalSimpleField(field, fieldValue, tag, parentSection); err != nil {
				return err
			}
		}
	}

	return nil
}

// unmarshalField sets a struct field value from a configuration value.
func (c *Config) unmarshalField(fieldValue reflect.Value, value value) error {
	fieldType := fieldValue.Type()

	switch fieldType.Kind() {
	case reflect.String:
		if value.vtype != stringType {
			return fmt.Errorf("expected string, got %s: %w", value.vtype, ErrTypeMismatch)
		}
		fieldValue.SetString(value.parsed.(string))

	case reflect.Int, reflect.Int32, reflect.Int64:
		if value.vtype != intType && value.vtype != hexType {
			return fmt.Errorf("expected int, got %s: %w", value.vtype, ErrTypeMismatch)
		}
		fieldValue.SetInt(int64(value.parsed.(int)))

	case reflect.Bool:
		if value.vtype != boolType {
			return fmt.Errorf("expected bool, got %s: %w", value.vtype, ErrTypeMismatch)
		}
		fieldValue.SetBool(value.parsed.(bool))

	case reflect.Float32, reflect.Float64:
		if value.vtype != floatType && value.vtype != intType {
			return fmt.Errorf("expected float, got %s: %w", value.vtype, ErrTypeMismatch)
		}
		if value.vtype == intType {
			fieldValue.SetFloat(float64(value.parsed.(int)))
		} else {
			fieldValue.SetFloat(value.parsed.(float64))
		}

	default:
		return fmt.Errorf("unsupported field type %s: %w", fieldType, ErrUnsupportedType)
	}

	return nil
}

// marshalStruct processes a struct and updates configuration values.
func (c *Config) marshalStruct(rv reflect.Value, parentSection string) error {
	rt := rv.Type()

	for i := range rv.NumField() {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !fieldValue.CanInterface() {
			continue // Skip unexported fields
		}

		tag := field.Tag.Get(TagName)
		if tag == "-" {
			continue // Skip fields explicitly marked to ignore
		}

		// Generate automatic tag if none provided
		if tag == "" {
			tag = c.generateFieldTag(field.Name, parentSection, fieldValue.Kind() == reflect.Struct)
		}

		if fieldValue.Kind() == reflect.Struct {
			if err := c.marshalNestedStruct(field, fieldValue, tag, parentSection); err != nil {
				return err
			}
		} else {
			if err := c.marshalSimpleField(field, fieldValue, tag, parentSection); err != nil {
				return err
			}
		}
	}

	return nil
}

// marshalNestedStruct handles marshaling of nested struct fields.
func (c *Config) marshalNestedStruct(field reflect.StructField, fieldValue reflect.Value, tag, parentSection string) error {
	// Handle nested struct marshaling with deep nesting support
	var sectionName string
	if strings.Contains(tag, ".") {
		// For explicit nested tags like "system.cpu", use the full tag as section name
		sectionName = strings.ToLower(tag)
	} else {
		// For simple tags, combine with parent section
		if parentSection != "" {
			sectionName = parentSection + "." + strings.ToLower(tag)
		} else {
			sectionName = strings.ToLower(tag)
		}
	}

	// Nested struct - recursively marshal
	if err := c.marshalStruct(fieldValue, sectionName); err != nil {
		return &MarshalError{
			Field:   field.Name,
			Section: sectionName,
			Key:     "",
			Err:     err,
		}
	}
	return nil
}

// marshalSimpleField handles marshaling of simple (non-struct) fields.
func (c *Config) marshalSimpleField(field reflect.StructField, fieldValue reflect.Value, tag, parentSection string) error {
	// Handle simple field marshaling
	tagInfo := c.parseTag(tag, parentSection)

	// Update configuration value
	if err := c.updateValue(tagInfo.Section, tagInfo.Key, fieldValue.Interface()); err != nil {
		return &MarshalError{
			Field:   field.Name,
			Section: tagInfo.Section,
			Key:     tagInfo.Key,
			Err:     err,
		}
	}
	return nil
}

// parseDefaultValue parses a default value string based on the target field type.
func (c *Config) parseDefaultValue(defaultStr string, fieldType reflect.Type) (value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return value{Raw: defaultStr, parsed: defaultStr, vtype: stringType}, nil

	case reflect.Int, reflect.Int32, reflect.Int64:
		// Check for hex format first
		if strings.HasPrefix(defaultStr, "0x") || strings.HasPrefix(defaultStr, "0X") {
			parsed, err := strconv.ParseInt(defaultStr, 0, 64)
			if err != nil {
				return value{}, fmt.Errorf("invalid hex default value %q: %w", defaultStr, err)
			}
			return value{Raw: defaultStr, parsed: int(parsed), vtype: hexType}, nil
		}
		// Regular integer
		parsed, err := strconv.ParseInt(defaultStr, 10, 64)
		if err != nil {
			return value{}, fmt.Errorf("invalid int default value %q: %w", defaultStr, err)
		}
		return value{Raw: defaultStr, parsed: int(parsed), vtype: intType}, nil

	case reflect.Bool:
		parsed, err := strconv.ParseBool(defaultStr)
		if err != nil {
			return value{}, fmt.Errorf("invalid bool default value %q: %w", defaultStr, err)
		}
		return value{Raw: defaultStr, parsed: parsed, vtype: boolType}, nil

	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(defaultStr, 64)
		if err != nil {
			return value{}, fmt.Errorf("invalid float default value %q: %w", defaultStr, err)
		}
		return value{Raw: defaultStr, parsed: parsed, vtype: floatType}, nil

	default:
		return value{}, fmt.Errorf("unsupported field type for default value: %s", fieldType)
	}
}

// convertToValue converts a Go value to a Config value.
func (c *Config) convertToValue(val any) (value, error) {
	switch v := val.(type) {
	case string:
		return value{Raw: v, parsed: v, vtype: stringType}, nil
	case int:
		return value{Raw: strconv.Itoa(v), parsed: v, vtype: intType}, nil
	case int32:
		return value{Raw: strconv.Itoa(int(v)), parsed: int(v), vtype: intType}, nil
	case int64:
		return value{Raw: strconv.Itoa(int(v)), parsed: int(v), vtype: intType}, nil
	case bool:
		return value{Raw: strconv.FormatBool(v), parsed: v, vtype: boolType}, nil
	case float64:
		return value{Raw: strconv.FormatFloat(v, 'g', -1, 64), parsed: v, vtype: floatType}, nil
	case float32:
		f64 := float64(v)
		return value{Raw: strconv.FormatFloat(f64, 'g', -1, 32), parsed: f64, vtype: floatType}, nil
	default:
		return value{}, fmt.Errorf("%w: %T", ErrUnsupportedType, val)
	}
}
