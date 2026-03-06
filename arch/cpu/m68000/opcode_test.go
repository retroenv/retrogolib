package m68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestDecodeLine0_ORI(t *testing.T) {
	// ORI.B #imm,<ea> = 0000 000 0 00 mmm rrr
	d, err := decodeOpcode(0x0000) // ORI.B #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insORI, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine0_ANDI(t *testing.T) {
	d, err := decodeOpcode(0x0240) // ANDI.W #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insANDI, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine0_SUBI(t *testing.T) {
	d, err := decodeOpcode(0x0480) // SUBI.L #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insSUBI, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine0_ADDI(t *testing.T) {
	d, err := decodeOpcode(0x0640) // ADDI.W #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insADDI, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine0_EORI(t *testing.T) {
	d, err := decodeOpcode(0x0A00) // EORI.B #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insEORI, d.Instruction)
}

func TestDecodeLine0_CMPI(t *testing.T) {
	d, err := decodeOpcode(0x0C00) // CMPI.B #imm,D0
	assert.NoError(t, err)
	assert.Equal(t, insCMPI, d.Instruction)
}

func TestDecodeLine0_BTSTReg(t *testing.T) {
	// BTST D0,D1 = 0000 000 1 00 000 001
	d, err := decodeOpcode(0x0101)
	assert.NoError(t, err)
	assert.Equal(t, insBTST, d.Instruction)
	assert.Equal(t, uint8(0), d.SrcReg)
	assert.Equal(t, uint8(1), d.DstReg)
}

func TestDecodeLine0_BSETReg(t *testing.T) {
	// BSET D0,D1 = 0000 000 1 11 000 001
	d, err := decodeOpcode(0x01C1)
	assert.NoError(t, err)
	assert.Equal(t, insBSET, d.Instruction)
}

func TestDecodeLine0_BTSTImm(t *testing.T) {
	// BTST #imm,D0 = 0000 100 0 00 000 000
	d, err := decodeOpcode(0x0800)
	assert.NoError(t, err)
	assert.Equal(t, insBTST, d.Instruction)
}

func TestDecodeLine1_MOVEB(t *testing.T) {
	// MOVE.B D0,D1 = 0001 001 000 000 000
	d, err := decodeOpcode(0x1200)
	assert.NoError(t, err)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine2_MOVEL(t *testing.T) {
	// MOVE.L D0,D1 = 0010 001 000 000 000
	d, err := decodeOpcode(0x2200)
	assert.NoError(t, err)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine2_MOVEAL(t *testing.T) {
	// MOVEA.L D0,A1 = 0010 001 001 000 000
	d, err := decodeOpcode(0x2240)
	assert.NoError(t, err)
	assert.Equal(t, insMOVEA, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine3_MOVEW(t *testing.T) {
	d, err := decodeOpcode(0x3200)
	assert.NoError(t, err)
	assert.Equal(t, insMOVE, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine4_NOP(t *testing.T) {
	d, err := decodeOpcode(0x4E71)
	assert.NoError(t, err)
	assert.Equal(t, insNOP, d.Instruction)
}

func TestDecodeLine4_RTS(t *testing.T) {
	d, err := decodeOpcode(0x4E75)
	assert.NoError(t, err)
	assert.Equal(t, insRTS, d.Instruction)
}

func TestDecodeLine4_RTE(t *testing.T) {
	d, err := decodeOpcode(0x4E73)
	assert.NoError(t, err)
	assert.Equal(t, insRTE, d.Instruction)
}

func TestDecodeLine4_TRAP(t *testing.T) {
	d, err := decodeOpcode(0x4E4F) // TRAP #15
	assert.NoError(t, err)
	assert.Equal(t, insTRAP, d.Instruction)
	assert.Equal(t, uint16(15), d.Extra)
}

func TestDecodeLine4_ILLEGAL(t *testing.T) {
	d, err := decodeOpcode(0x4AFC)
	assert.NoError(t, err)
	assert.Equal(t, insILLEGAL, d.Instruction)
}

func TestDecodeLine4_LEA(t *testing.T) {
	// LEA (A0),A1 = 0100 001 111 010 000
	d, err := decodeOpcode(0x43D0)
	assert.NoError(t, err)
	assert.Equal(t, insLEA, d.Instruction)
}

func TestDecodeLine4_CLR(t *testing.T) {
	// CLR.B D0 = 0100 001 0 00 000 000
	d, err := decodeOpcode(0x4200)
	assert.NoError(t, err)
	assert.Equal(t, insCLR, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine4_NEG(t *testing.T) {
	d, err := decodeOpcode(0x4400)
	assert.NoError(t, err)
	assert.Equal(t, insNEG, d.Instruction)
	assert.Equal(t, SizeByte, d.Size)
}

func TestDecodeLine4_SWAP(t *testing.T) {
	d, err := decodeOpcode(0x4840) // SWAP D0
	assert.NoError(t, err)
	assert.Equal(t, insSWAP, d.Instruction)
}

func TestDecodeLine4_EXT(t *testing.T) {
	d, err := decodeOpcode(0x4880) // EXT.W D0
	assert.NoError(t, err)
	assert.Equal(t, insEXT, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine5_ADDQ(t *testing.T) {
	// ADDQ.W #3,D0 = 0101 011 0 01 000 000
	d, err := decodeOpcode(0x5640)
	assert.NoError(t, err)
	assert.Equal(t, insADDQ, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
	assert.Equal(t, uint16(3), d.Extra)
}

func TestDecodeLine5_SUBQ(t *testing.T) {
	// SUBQ.L #1,D0 = 0101 001 1 10 000 000
	d, err := decodeOpcode(0x5380)
	assert.NoError(t, err)
	assert.Equal(t, insSUBQ, d.Instruction)
	assert.Equal(t, SizeLong, d.Size)
}

func TestDecodeLine5_DBcc(t *testing.T) {
	// DBcc D0 = 0101 cccc 11 001 000
	d, err := decodeOpcode(0x51C8) // DBRA D0 (false condition)
	assert.NoError(t, err)
	assert.Equal(t, insDBcc, d.Instruction)
}

func TestDecodeLine6_BRA(t *testing.T) {
	d, err := decodeOpcode(0x6000) // BRA with 16-bit displacement
	assert.NoError(t, err)
	assert.Equal(t, insBRA, d.Instruction)
}

func TestDecodeLine6_BSR(t *testing.T) {
	d, err := decodeOpcode(0x6100) // BSR with 16-bit displacement
	assert.NoError(t, err)
	assert.Equal(t, insBSR, d.Instruction)
}

func TestDecodeLine6_Bcc(t *testing.T) {
	d, err := decodeOpcode(0x6700) // BEQ with 16-bit displacement
	assert.NoError(t, err)
	assert.Equal(t, insBcc, d.Instruction)
	assert.Equal(t, uint16(7), d.Extra) // EQ condition
}

func TestDecodeLine7_MOVEQ(t *testing.T) {
	// MOVEQ #42,D3 = 0111 011 0 00101010
	d, err := decodeOpcode(0x762A)
	assert.NoError(t, err)
	assert.Equal(t, insMOVEQ, d.Instruction)
	assert.Equal(t, uint8(3), d.DstReg)
	assert.Equal(t, uint16(42), d.Extra)
}

func TestDecodeLine8_OR(t *testing.T) {
	// OR.W D0,D1 => 1000 001 001 000 000
	d, err := decodeOpcode(0x8240)
	assert.NoError(t, err)
	assert.Equal(t, insOR, d.Instruction)
	assert.Equal(t, SizeWord, d.Size)
}

func TestDecodeLine8_DIVU(t *testing.T) {
	// DIVU D0,D1 = 1000 001 011 000 000
	d, err := decodeOpcode(0x82C0)
	assert.NoError(t, err)
	assert.Equal(t, insDIVU, d.Instruction)
}

func TestDecodeLine8_DIVS(t *testing.T) {
	// DIVS D0,D1 = 1000 001 111 000 000
	d, err := decodeOpcode(0x83C0)
	assert.NoError(t, err)
	assert.Equal(t, insDIVS, d.Instruction)
}

func TestDecodeLine9_SUB(t *testing.T) {
	d, err := decodeOpcode(0x9040) // SUB.W D0,D1
	assert.NoError(t, err)
	assert.Equal(t, insSUB, d.Instruction)
}

func TestDecodeLine9_SUBA(t *testing.T) {
	d, err := decodeOpcode(0x90C0) // SUBA.W D0,A0
	assert.NoError(t, err)
	assert.Equal(t, insSUBA, d.Instruction)
}

func TestDecodeLineA(t *testing.T) {
	d, err := decodeOpcode(0xA000)
	assert.NoError(t, err)
	assert.Equal(t, insILLEGAL, d.Instruction)
}

func TestDecodeLineB_CMP(t *testing.T) {
	d, err := decodeOpcode(0xB040) // CMP.W D0,D1
	assert.NoError(t, err)
	assert.Equal(t, insCMP, d.Instruction)
}

func TestDecodeLineB_CMPA(t *testing.T) {
	d, err := decodeOpcode(0xB0C0) // CMPA.W D0,A0
	assert.NoError(t, err)
	assert.Equal(t, insCMPA, d.Instruction)
}

func TestDecodeLineB_EOR(t *testing.T) {
	// EOR.W D1,D0 = 1011 001 101 000 000
	d, err := decodeOpcode(0xB340)
	assert.NoError(t, err)
	assert.Equal(t, insEOR, d.Instruction)
}

func TestDecodeLineC_AND(t *testing.T) {
	d, err := decodeOpcode(0xC040) // AND.W D0,D1
	assert.NoError(t, err)
	assert.Equal(t, insAND, d.Instruction)
}

func TestDecodeLineC_MULU(t *testing.T) {
	d, err := decodeOpcode(0xC0C0) // MULU D0,D0
	assert.NoError(t, err)
	assert.Equal(t, insMULU, d.Instruction)
}

func TestDecodeLineC_MULS(t *testing.T) {
	d, err := decodeOpcode(0xC1C0) // MULS D0,D0
	assert.NoError(t, err)
	assert.Equal(t, insMULS, d.Instruction)
}

func TestDecodeLineD_ADD(t *testing.T) {
	d, err := decodeOpcode(0xD040) // ADD.W D0,D1
	assert.NoError(t, err)
	assert.Equal(t, insADD, d.Instruction)
}

func TestDecodeLineD_ADDA(t *testing.T) {
	d, err := decodeOpcode(0xD0C0) // ADDA.W D0,A0
	assert.NoError(t, err)
	assert.Equal(t, insADDA, d.Instruction)
}

func TestDecodeLineE_ASL(t *testing.T) {
	// ASL.W #1,D0 = 1110 001 1 01 0 00 000
	d, err := decodeOpcode(0xE340)
	assert.NoError(t, err)
	assert.Equal(t, insASL, d.Instruction)
}

func TestDecodeLineE_LSR(t *testing.T) {
	// LSR.B #1,D0 = 1110 001 0 00 0 01 000
	d, err := decodeOpcode(0xE208)
	assert.NoError(t, err)
	assert.Equal(t, insLSR, d.Instruction)
}

func TestDecodeLineE_ROL(t *testing.T) {
	// ROL.W #1,D0 = 1110 001 1 01 0 11 000
	d, err := decodeOpcode(0xE358)
	assert.NoError(t, err)
	assert.Equal(t, insROL, d.Instruction)
}

func TestDecodeLineF(t *testing.T) {
	d, err := decodeOpcode(0xF000)
	assert.NoError(t, err)
	assert.Equal(t, insILLEGAL, d.Instruction)
}
