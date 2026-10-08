package cpu68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOverlappingOpcodeFamilies(t *testing.T) {
	// Broad masks previously decoded these legal forms as BCD or bit operations.
	tests := []struct {
		opcode      uint16
		instruction *Instruction
	}{
		{opcode: 0xC92D, instruction: insAND},
		{opcode: 0x8D14, instruction: insOR},
		{opcode: 0x0108, instruction: insMOVEP},
		{opcode: 0x0188, instruction: insMOVEP},
		{opcode: 0xC101, instruction: insABCD},
		{opcode: 0x8101, instruction: insSBCD},
	}

	for _, tt := range tests {
		opcode := decodeOpcode(tt.opcode)
		assert.Equal(t, tt.instruction, opcode.Instruction)
	}
}

func TestDecodeRejectsForbiddenAddressingModes(t *testing.T) {
	// These encodings executed as instructions although the 68000 traps them.
	tests := []struct {
		name   string
		opcode uint16
	}{
		{name: "MOVE.B An,Dn", opcode: 0x1008},
		{name: "MOVE.B Dn,An", opcode: 0x1040},
		{name: "ORI with size bits 11", opcode: 0x00C0},
		{name: "ADDI #,#", opcode: 0x067C},
		{name: "ORI.L #,SR", opcode: 0x00BC},
		{name: "BTST #,#", opcode: 0x083C},
		{name: "BSET Dn,#", opcode: 0x01FC},
		{name: "LEA Dn,An", opcode: 0x41C0},
		{name: "LEA (An)+,An", opcode: 0x41D8},
		{name: "JSR Dn", opcode: 0x4E88},
		{name: "JMP -(An)", opcode: 0x4EE0},
		{name: "PEA An", opcode: 0x4848},
		{name: "MOVEM An,-(An)", opcode: 0x4888},
		{name: "MOVEM (An)+ to memory", opcode: 0x4898},
		{name: "MOVEM -(An) to registers", opcode: 0x4CA0},
		{name: "TST.W An", opcode: 0x4A48},
		{name: "TAS #", opcode: 0x4AFC},
		{name: "CLR An", opcode: 0x4248},
		{name: "MOVE SR,An", opcode: 0x40C8},
		{name: "MOVE An,CCR", opcode: 0x44C8},
		{name: "CHK An,Dn", opcode: 0x4188},
		{name: "ADDQ.B #8,An", opcode: 0x5008},
		{name: "Scc (d16,PC)", opcode: 0x50FA},
		{name: "OR.B An,Dn", opcode: 0x8008},
		{name: "OR.W An,Dn", opcode: 0x8048},
		{name: "AND.L An,Dn", opcode: 0xC088},
		{name: "OR Dn,Dn reversed form", opcode: 0x8100 | 0x0040},
		{name: "AND.L Dn,Dn reversed form", opcode: 0xC180},
		{name: "DIVU An,Dn", opcode: 0x80C8},
		{name: "MULS #,Dn with mode 7 reg 5", opcode: 0xC1FD},
		{name: "SUB.B An,Dn", opcode: 0x9008},
		{name: "ADDA mode 7 reg 6", opcode: 0xD0FE},
		{name: "CMP.B An,Dn", opcode: 0xB008},
		{name: "EOR Dn,(d16,PC)", opcode: 0xB13A},
		{name: "memory shift on Dn", opcode: 0xE0C0},
		{name: "memory shift with bit 11", opcode: 0xE8D0},
		{name: "MOVEQ with bit 8", opcode: 0x7100},
		{name: "MOVEC", opcode: 0x4E7A},
		{name: "MOVE CCR,<ea>", opcode: 0x42C0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := decodeOpcode(tt.opcode)
			assert.Equal(t, insILLEGAL, d.Instruction)
			assert.Equal(t, uint16(0), d.Extra)
		})
	}
}

func TestDecodeAcceptsAddressRegisterForms(t *testing.T) {
	// Word and long sources, and register-list forms, keep address registers.
	tests := []struct {
		opcode      uint16
		instruction *Instruction
	}{
		{opcode: 0x3008, instruction: insMOVE},          // MOVE.W A0,D0.
		{opcode: 0xD048, instruction: insADD},           // ADD.W A0,D0.
		{opcode: 0xB088, instruction: insCMP},           // CMP.L A0,D0.
		{opcode: 0x5048 | 0x0040, instruction: insADDQ}, // ADDQ.W #8,A0.
		{opcode: 0x48A0, instruction: insMOVEM},         // MOVEM.W list,-(A0).
		{opcode: 0x4C98, instruction: insMOVEM},         // MOVEM.W (A0)+,list.
		{opcode: 0x0138, instruction: insBTST},          // BTST D0,(xxx).W.
		{opcode: 0x013C, instruction: insBTST},          // BTST D0,#imm.
		{opcode: 0x4AC0, instruction: insTAS},           // TAS D0.
		{opcode: 0x4E90, instruction: insJSR},           // JSR (A0).
		{opcode: 0xE0D0, instruction: insASR},           // ASR.W (A0).
	}

	for _, tt := range tests {
		d := decodeOpcode(tt.opcode)
		assert.Equal(t, tt.instruction, d.Instruction)
	}
}

func TestDecodeLine0_ORI(t *testing.T) {
	// ORI.B #imm,<ea> = 0000 000 0 00 mmm rrr
	d := decodeOpcode(0x0000) // ORI.B #imm,D0
	assert.Equal(t, insORI, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine0_ANDI(t *testing.T) {
	d := decodeOpcode(0x0240) // ANDI.W #imm,D0
	assert.Equal(t, insANDI, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine0_SUBI(t *testing.T) {
	d := decodeOpcode(0x0480) // SUBI.L #imm,D0
	assert.Equal(t, insSUBI, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine0_ADDI(t *testing.T) {
	d := decodeOpcode(0x0640) // ADDI.W #imm,D0
	assert.Equal(t, insADDI, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine0_EORI(t *testing.T) {
	d := decodeOpcode(0x0A00) // EORI.B #imm,D0
	assert.Equal(t, insEORI, d.Instruction)
}

func TestDecodeLine0_CMPI(t *testing.T) {
	d := decodeOpcode(0x0C00) // CMPI.B #imm,D0
	assert.Equal(t, insCMPI, d.Instruction)
}

func TestDecodeLine0_BTSTReg(t *testing.T) {
	// BTST D0,D1 = 0000 000 1 00 000 001
	d := decodeOpcode(0x0101)
	assert.Equal(t, insBTST, d.Instruction)
	assert.Equal(t, uint8(0), d.SrcReg)
	assert.Equal(t, uint8(1), d.DstReg)
}

func TestDecodeLine0_BSETReg(t *testing.T) {
	// BSET D0,D1 = 0000 000 1 11 000 001
	d := decodeOpcode(0x01C1)
	assert.Equal(t, insBSET, d.Instruction)
}

func TestDecodeLine0_BTSTImm(t *testing.T) {
	// BTST #imm,D0 = 0000 100 0 00 000 000
	d := decodeOpcode(0x0800)
	assert.Equal(t, insBTST, d.Instruction)
}

func TestDecodeLine1_MOVEB(t *testing.T) {
	// MOVE.B D0,D1 = 0001 001 000 000 000
	d := decodeOpcode(0x1200)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine2_MOVEL(t *testing.T) {
	// MOVE.L D0,D1 = 0010 001 000 000 000
	d := decodeOpcode(0x2200)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine2_MOVEAL(t *testing.T) {
	// MOVEA.L D0,A1 = 0010 001 001 000 000
	d := decodeOpcode(0x2240)
	assert.Equal(t, insMOVEA, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine3_MOVEW(t *testing.T) {
	d := decodeOpcode(0x3200)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine4_NOP(t *testing.T) {
	d := decodeOpcode(0x4E71)
	assert.Equal(t, insNOP, d.Instruction)
}

func TestDecodeLine4_RTS(t *testing.T) {
	d := decodeOpcode(0x4E75)
	assert.Equal(t, insRTS, d.Instruction)
}

func TestDecodeLine4_RTE(t *testing.T) {
	d := decodeOpcode(0x4E73)
	assert.Equal(t, insRTE, d.Instruction)
}

func TestDecodeLine4_TRAP(t *testing.T) {
	d := decodeOpcode(0x4E4F) // TRAP #15
	assert.Equal(t, insTRAP, d.Instruction)
	assert.Equal(t, uint16(15), d.Extra)
}

func TestDecodeLine4_ILLEGAL(t *testing.T) {
	d := decodeOpcode(0x4AFC)
	assert.Equal(t, insILLEGAL, d.Instruction)
}

func TestDecodeLine4_LEA(t *testing.T) {
	// LEA (A0),A1 = 0100 001 111 010 000
	d := decodeOpcode(0x43D0)
	assert.Equal(t, insLEA, d.Instruction)
}

func TestDecodeLine4_CLR(t *testing.T) {
	// CLR.B D0 = 0100 001 0 00 000 000
	d := decodeOpcode(0x4200)
	assert.Equal(t, insCLR, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine4_NEG(t *testing.T) {
	d := decodeOpcode(0x4400)
	assert.Equal(t, insNEG, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine4_SWAP(t *testing.T) {
	d := decodeOpcode(0x4840) // SWAP D0
	assert.Equal(t, insSWAP, d.Instruction)
}

func TestDecodeLine4_EXT(t *testing.T) {
	d := decodeOpcode(0x4880) // EXT.W D0
	assert.Equal(t, insEXT, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine5_ADDQ(t *testing.T) {
	// ADDQ.W #3,D0 = 0101 011 0 01 000 000
	d := decodeOpcode(0x5640)
	assert.Equal(t, insADDQ, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
	assert.Equal(t, uint16(3), d.Extra)
}

func TestDecodeLine5_SUBQ(t *testing.T) {
	// SUBQ.L #1,D0 = 0101 001 1 10 000 000
	d := decodeOpcode(0x5380)
	assert.Equal(t, insSUBQ, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine5_DBcc(t *testing.T) {
	// DBcc D0 = 0101 cccc 11 001 000
	d := decodeOpcode(0x51C8) // DBRA D0 (false condition)
	assert.Equal(t, insDBcc, d.Instruction)
}

func TestDecodeLine6_BRA(t *testing.T) {
	d := decodeOpcode(0x6000) // BRA with 16-bit displacement
	assert.Equal(t, insBRA, d.Instruction)
}

func TestDecodeLine6_BSR(t *testing.T) {
	d := decodeOpcode(0x6100) // BSR with 16-bit displacement
	assert.Equal(t, insBSR, d.Instruction)
}

func TestDecodeLine6_Bcc(t *testing.T) {
	d := decodeOpcode(0x6700) // BEQ with 16-bit displacement
	assert.Equal(t, insBcc, d.Instruction)
	assert.Equal(t, uint16(7), d.Extra) // EQ condition
}

func TestDecodeLine7_MOVEQ(t *testing.T) {
	// MOVEQ #42,D3 = 0111 011 0 00101010
	d := decodeOpcode(0x762A)
	assert.Equal(t, insMOVEQ, d.Instruction)
	assert.Equal(t, uint8(3), d.DstReg)
	assert.Equal(t, uint16(42), d.Extra)
}

func TestDecodeLine8_OR(t *testing.T) {
	// OR.W D0,D1 => 1000 001 001 000 000
	d := decodeOpcode(0x8240)
	assert.Equal(t, insOR, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine8_DIVU(t *testing.T) {
	// DIVU D0,D1 = 1000 001 011 000 000
	d := decodeOpcode(0x82C0)
	assert.Equal(t, insDIVU, d.Instruction)
}

func TestDecodeLine8_DIVS(t *testing.T) {
	// DIVS D0,D1 = 1000 001 111 000 000
	d := decodeOpcode(0x83C0)
	assert.Equal(t, insDIVS, d.Instruction)
}

func TestDecodeLine9_SUB(t *testing.T) {
	d := decodeOpcode(0x9040) // SUB.W D0,D1
	assert.Equal(t, insSUB, d.Instruction)
}

func TestDecodeLine9_SUBA(t *testing.T) {
	d := decodeOpcode(0x90C0) // SUBA.W D0,A0
	assert.Equal(t, insSUBA, d.Instruction)
}

func TestDecodeLineA(t *testing.T) {
	d := decodeOpcode(0xA000)
	assert.Equal(t, insILLEGAL, d.Instruction)
}

func TestDecodeLineB_CMP(t *testing.T) {
	d := decodeOpcode(0xB040) // CMP.W D0,D1
	assert.Equal(t, insCMP, d.Instruction)
}

func TestDecodeLineB_CMPA(t *testing.T) {
	d := decodeOpcode(0xB0C0) // CMPA.W D0,A0
	assert.Equal(t, insCMPA, d.Instruction)
}

func TestDecodeLineB_EOR(t *testing.T) {
	// EOR.W D1,D0 = 1011 001 101 000 000
	d := decodeOpcode(0xB340)
	assert.Equal(t, insEOR, d.Instruction)
}

func TestDecodeLineC_AND(t *testing.T) {
	d := decodeOpcode(0xC040) // AND.W D0,D1
	assert.Equal(t, insAND, d.Instruction)
}

func TestDecodeLineC_MULU(t *testing.T) {
	d := decodeOpcode(0xC0C0) // MULU D0,D0
	assert.Equal(t, insMULU, d.Instruction)
}

func TestDecodeLineC_MULS(t *testing.T) {
	d := decodeOpcode(0xC1C0) // MULS D0,D0
	assert.Equal(t, insMULS, d.Instruction)
}

func TestDecodeLineD_ADD(t *testing.T) {
	d := decodeOpcode(0xD040) // ADD.W D0,D1
	assert.Equal(t, insADD, d.Instruction)
}

func TestDecodeLineD_ADDA(t *testing.T) {
	d := decodeOpcode(0xD0C0) // ADDA.W D0,A0
	assert.Equal(t, insADDA, d.Instruction)
}

func TestDecodeLineE_ASL(t *testing.T) {
	// ASL.W #1,D0 = 1110 001 1 01 0 00 000
	d := decodeOpcode(0xE340)
	assert.Equal(t, insASL, d.Instruction)
}

func TestDecodeLineE_LSR(t *testing.T) {
	// LSR.B #1,D0 = 1110 001 0 00 0 01 000
	d := decodeOpcode(0xE208)
	assert.Equal(t, insLSR, d.Instruction)
}

func TestDecodeLineE_ROL(t *testing.T) {
	// ROL.W #1,D0 = 1110 001 1 01 0 11 000
	d := decodeOpcode(0xE358)
	assert.Equal(t, insROL, d.Instruction)
}

func TestDecodeLineF(t *testing.T) {
	d := decodeOpcode(0xF000)
	assert.Equal(t, insILLEGAL, d.Instruction)
}
