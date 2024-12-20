package chip8

// Mode defines an address mode.
type Mode int

// addressing modes.
const (
	NoAddressing Mode = 0
	// ImpliedAddressing has no parameters
	ImpliedAddressing Mode = 1 << iota
	// AbsoluteAddressing has an absolute address - addr
	AbsoluteAddressing
	// V0AbsoluteAddressing has V0 and an absolute address - V0, addr
	V0AbsoluteAddressing
	// RegisterAddressing has a register - Vx
	RegisterAddressing
	// RegisterValueAddressing has a register and a byte value - Vx, byte
	RegisterValueAddressing
	// RegisterRegisterAddressing has two registers - Vx, Vy
	RegisterRegisterAddressing
	// RegisterRegisterNibbleAddressing has two registers and a nibble value - Vx, Vy, nibble
	RegisterRegisterNibbleAddressing
	// RegisterDTAddressing has DT and a register - DT, Vx
	RegisterDTAddressing
	// RegisterKAddressing has K and a register - K, Vx
	RegisterKAddressing
	// RegisterIndirectIAddressing has indirect I and a register - Vx, [I]
	RegisterIndirectIAddressing
	// DTRegisterAddressing has DT and a register - DT, Vx
	DTRegisterAddressing
	// STRegisterAddressing has ST and a register - ST, Vx
	STRegisterAddressing
	// FRegisterAddressing has F and a register - F, Vx
	FRegisterAddressing
	// BRegisterAddressing has B and a register - B, Vx
	BRegisterAddressing
	// IAbsoluteAddressing has I and an absolute address - I, addr
	IAbsoluteAddressing
	// IRegisterAddressing has I and a register - I, Vx
	IRegisterAddressing
	// IIndirectRegisterAddressing has I and an indirect register - I, [Vx]
	IIndirectRegisterAddressing
)
