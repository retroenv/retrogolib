package analysis

import "testing"

func BenchmarkTransfer(b *testing.B) {
	for _, test := range []struct {
		name string
		code byte
	}{
		{
			name: "AND-immediate",
			code: 0x29,
		},
		{
			name: "ADC-immediate",
			code: 0x69,
		},
		{
			name: "SLO-indirect-X",
			code: 0x03,
		},
	} {
		b.Run(test.name, func(b *testing.B) {
			input := fixture(test.code, 2, 0x21)
			result := bufferedResult()
			b.ReportAllocs()
			for b.Loop() {
				Transfer(input, &result)
			}
		})
	}
}
