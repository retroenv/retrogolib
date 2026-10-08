package cp1610

// OpcodeInfo gives the first word and base timing of an encoding.
// Size is in words. Some immediate instructions use one more word after SDBD.
type OpcodeInfo struct {
	Word   uint16
	Size   uint8
	Cycles uint8
}

// Instruction groups all encodings of one CP1610 mnemonic.
type Instruction struct {
	Name      string
	Encodings map[OperandKey]OpcodeInfo
}
