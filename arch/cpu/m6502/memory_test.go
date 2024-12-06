package m6502

type testMemory struct {
	b [0x10000]byte
}

func (m *testMemory) Read(address uint16) uint8 {
	return m.b[address]
}

func (m *testMemory) Write(address uint16, value uint8) {
	m.b[address] = value
}
