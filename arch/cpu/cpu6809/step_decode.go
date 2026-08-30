package cpu6809

type decodedOpcode struct {
	opcode      Opcode
	opcodeBytes [2]byte
	baseOffset  uint16
}
