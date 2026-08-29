package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// parser handles parsing configuration data with structure tracking.
type parser struct {
	data           []byte
	line           int
	config         *Config
	currentSection string
	seenItems      set.Set[string] // Track seen sections and keys using composite keys
	itemLines      map[string]int  // Track line numbers for sections and keys
}

// parse parses the configuration data.
func (p *parser) parse() error {
	content := string(p.data)
	lines := strings.Split(content, "\n")

	if len(lines) > maxLines {
		return fmt.Errorf("%w: %d exceeds limit of %d", ErrTooManyLines, len(lines), maxLines)
	}

	for i, line := range lines {
		p.line = i + 1
		if err := p.parseLine(line); err != nil {
			return &ParseError{
				Line: p.line,
				Pos:  0,
				Msg:  err.Error(),
				Err:  err,
			}
		}
	}

	return nil
}

// parseLine parses a single line and tracks structure.
func (p *parser) parseLine(line string) error {
	original := line
	trimmed := strings.TrimSpace(line)

	// Track original structure element
	element := StructureElement{
		Line:    p.line,
		Content: original,
		Section: p.currentSection,
	}

	switch {
	case trimmed == "":
		element.Type = emptyLineElement

	case strings.HasPrefix(trimmed, "#"):
		element.Type = commentElement
		comment := Comment{
			Line:    p.line,
			Text:    strings.TrimSpace(trimmed[1:]),
			Section: p.currentSection,
		}
		p.config.comments = append(p.config.comments, comment)

	case strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"):
		element.Type = sectionElement
		if err := p.parseSection(trimmed, &element); err != nil {
			return err
		}

	case strings.Contains(trimmed, "="):
		element.Type = keyValueElement
		if err := p.parseKeyValue(trimmed, &element); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid line format: %s", trimmed)
	}

	p.config.structure = append(p.config.structure, element)
	return nil
}

// parseSection parses a section header and validates it.
func (p *parser) parseSection(trimmed string, element *StructureElement) error {
	sectionName := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	if sectionName == "" {
		return errors.New("empty section name")
	}

	if len(sectionName) > maxNameLength {
		return fmt.Errorf("%w: %d characters exceeds limit of %d", ErrSectionNameTooLong, len(sectionName), maxNameLength)
	}

	// Normalize section name to lowercase for case-insensitive comparison
	normalizedSection := strings.ToLower(sectionName)
	sectionKey := "section:" + normalizedSection
	if p.seenItems.Contains(sectionKey) {
		firstLine := p.itemLines[sectionKey]
		return fmt.Errorf("%w: section '%s' first defined at line %d", ErrDuplicateSection, sectionName, firstLine)
	}

	// Track this section using normalized name
	p.seenItems.Add(sectionKey)
	p.itemLines[sectionKey] = p.line

	p.currentSection = normalizedSection
	element.Section = normalizedSection
	if p.config.sections[normalizedSection] == nil {
		p.config.sections[normalizedSection] = make(Section)
	}

	return nil
}

// parseKeyValue parses a key-value pair.
func (p *parser) parseKeyValue(line string, element *StructureElement) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid key-value format: %s", line)
	}

	key := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	if key == "" {
		return errors.New("empty key name")
	}
	if len(key) > maxNameLength {
		return fmt.Errorf("%w: %d characters exceeds limit of %d", ErrKeyNameTooLong, len(key), maxNameLength)
	}

	// Allow root-level keys (when currentSection is empty)
	sectionName := p.currentSection
	if sectionName == "" {
		sectionName = "" // Use empty string for root-level keys
	}

	// Normalize key name to lowercase for case-insensitive comparison
	normalizedKey := strings.ToLower(key)
	keyItem := "key:" + sectionName + ":" + normalizedKey
	if p.seenItems.Contains(keyItem) {
		firstLine := p.itemLines[keyItem]
		if sectionName == "" {
			return fmt.Errorf("%w: key '%s' at root level first defined at line %d", ErrDuplicateKey, key, firstLine)
		}
		return fmt.Errorf("%w: key '%s' in section '%s' first defined at line %d", ErrDuplicateKey, key, sectionName, firstLine)
	}

	// Track this key using normalized name
	p.seenItems.Add(keyItem)
	p.itemLines[keyItem] = p.line

	element.Key = normalizedKey
	element.Section = sectionName

	value, err := p.parseValue(valueStr)
	if err != nil {
		return fmt.Errorf("parsing value for key %s: %w", key, err)
	}

	if p.config.sections[sectionName] == nil {
		p.config.sections[sectionName] = make(Section)
	}

	p.config.sections[sectionName][normalizedKey] = value
	return nil
}

// parseValue parses a configuration value and determines its type.
func (p *parser) parseValue(valueStr string) (Value, error) {
	if valueStr == "" {
		return Value{Raw: "", parsed: "", vtype: stringType}, nil
	}

	// Check for quoted string
	if strings.HasPrefix(valueStr, `"`) {
		unquoted, err := strconv.Unquote(valueStr)
		if err != nil {
			return Value{}, fmt.Errorf("invalid quoted string: %w", err)
		}
		return Value{Raw: unquoted, parsed: unquoted, vtype: stringType}, nil
	}

	// Check for boolean
	if valueStr == "true" || valueStr == "false" {
		parsed, _ := strconv.ParseBool(valueStr)
		return Value{Raw: valueStr, parsed: parsed, vtype: boolType}, nil
	}

	// Check for hexadecimal
	if strings.HasPrefix(valueStr, "0x") || strings.HasPrefix(valueStr, "0X") {
		parsed, err := strconv.ParseInt(valueStr, 0, 64)
		if err != nil {
			return Value{}, fmt.Errorf("invalid hex value: %w", err)
		}
		return Value{Raw: valueStr, parsed: int(parsed), vtype: hexType}, nil
	}

	// Check for integer
	if intVal, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return Value{Raw: valueStr, parsed: int(intVal), vtype: intType}, nil
	}

	// Check for float
	if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return Value{Raw: valueStr, parsed: floatVal, vtype: floatType}, nil
	}

	// Default to string
	return Value{Raw: valueStr, parsed: valueStr, vtype: stringType}, nil
}
