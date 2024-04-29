package assert

import "testing"

func TestInterfaceNilEqual(t *testing.T) {
	var values []int
	Equal(t, nil, values)
}
