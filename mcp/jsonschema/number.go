package jsonschema

import (
	"encoding/json"
	"math/big"
	"strings"
)

// schemaNumber keeps powers separate to avoid expanding large exponents.
// Its input must be a valid JSON number.
type schemaNumber struct {
	negative bool
	digits   string
	exponent big.Int
}

func (number *schemaNumber) compare(other *schemaNumber) int {
	if number.negative != other.negative {
		if number.negative {
			return -1
		}

		return 1
	}

	result := number.compareMagnitude(other)
	if number.negative {
		result = -result
	}

	return result
}

func (number *schemaNumber) compareMagnitude(other *schemaNumber) int {
	switch {
	case number.digits == "0" && other.digits == "0":
		return 0

	case number.digits == "0":
		return -1

	case other.digits == "0":
		return 1
	}

	left := new(big.Int).Add(&number.exponent, big.NewInt(int64(len(number.digits))))
	right := new(big.Int).Add(&other.exponent, big.NewInt(int64(len(other.digits))))
	if order := left.Cmp(right); order != 0 {
		return order
	}

	for index := range max(len(number.digits), len(other.digits)) {
		a := byte('0')
		b := byte('0')
		if index < len(number.digits) {
			a = number.digits[index]
		}

		if index < len(other.digits) {
			b = other.digits[index]
		}

		switch {
		case a < b:
			return -1

		case a > b:
			return 1
		}
	}

	return 0
}

func parseSchemaNumber(raw json.Number) *schemaNumber {
	text := raw.String()
	number := &schemaNumber{}
	if strings.HasPrefix(text, "-") {
		number.negative = true
		text = text[1:]
	}

	if index := strings.IndexAny(text, "eE"); index >= 0 {
		number.exponent.SetString(text[index+1:], 10)
		text = text[:index]
	}

	integer, fraction, _ := strings.Cut(text, ".")
	digits := strings.TrimLeft(integer+fraction, "0")
	if digits == "" {
		number.digits = "0"
		number.negative = false
		number.exponent.SetInt64(0)

		return number
	}

	number.digits = strings.TrimRight(digits, "0")
	shift := len(digits) - len(number.digits) - len(fraction)
	number.exponent.Add(&number.exponent, big.NewInt(int64(shift)))

	return number
}

// canonicalJSON gives equal numbers, such as 1 and 1.0, the same key.
func canonicalJSON(value any) string {
	data, _ := json.Marshal(normalizeSchemaNumbers(value))

	return string(data)
}

func normalizeSchemaNumbers(value any) any {
	switch value := value.(type) {
	case json.Number:
		number := parseSchemaNumber(value)
		sign := ""
		if number.negative {
			sign = "-"
		}

		return json.Number(sign + number.digits + "e" + number.exponent.String())

	case []any:
		result := make([]any, len(value))
		for index, child := range value {
			result[index] = normalizeSchemaNumbers(child)
		}

		return result

	case map[string]any:
		result := make(map[string]any, len(value))
		for name, child := range value {
			result[name] = normalizeSchemaNumbers(child)
		}

		return result

	default:
		return value
	}
}
