package analysis

import "testing"

func BenchmarkTransfer(b *testing.B) {
	for _, test := range []struct {
		name string
		code byte
	}{
		{
			name: "LD-register",
			code: 0x41,
		},
		{
			name: "AND-immediate",
			code: 0xe6,
		},
		{
			name: "CP-immediate",
			code: 0xfe,
		},
	} {
		b.Run(test.name, func(b *testing.B) {
			input := fixture(0, test.code, 0x21)
			result := bufferedResult()
			b.ReportAllocs()
			for b.Loop() {
				Transfer(input, &result)
			}
		})
	}
}
