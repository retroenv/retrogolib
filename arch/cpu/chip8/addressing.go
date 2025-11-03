package chip8

// Mode specifies how a Chip-8 instruction accesses its operands.
// Multiple modes can be combined using bitwise OR for instructions that support variants.
//
// Available modes:
//   - ImpliedAddressing: No operands (CLS, RET)
//   - AbsoluteAddressing: 12-bit address (JP addr, CALL addr)
//   - V0AbsoluteAddressing: V0 + address (JP V0, addr)
//   - RegisterAddressing: Single register (Vx)
//   - RegisterValueAddressing: Register + byte value (LD Vx, byte, ADD Vx, byte)
//   - RegisterRegisterAddressing: Two registers (LD Vx, Vy, ADD Vx, Vy)
//   - RegisterRegisterNibbleAddressing: Two registers + nibble (DRW Vx, Vy, nibble)
//   - RegisterDTAddressing: Register from delay timer (LD Vx, DT)
//   - RegisterKAddressing: Register from key press (LD Vx, K)
//   - RegisterIndirectIAddressing: Register from [I] (LD Vx, [I])
//   - DTRegisterAddressing: Delay timer from register (LD DT, Vx)
//   - STRegisterAddressing: Sound timer from register (LD ST, Vx)
//   - FRegisterAddressing: Font location from register (LD F, Vx)
//   - BRegisterAddressing: BCD from register (LD B, Vx)
//   - IAbsoluteAddressing: I register + address (LD I, addr)
//   - IRegisterAddressing: I register + register (ADD I, Vx)
//   - IIndirectRegisterAddressing: [I] from register (LD [I], Vx)
type Mode int

const (
	NoAddressing                     Mode = 0
	ImpliedAddressing                Mode = 1 << iota // No operands (CLS, RET)
	AbsoluteAddressing                                // 12-bit address (JP addr)
	V0AbsoluteAddressing                              // V0 + address (JP V0, addr)
	RegisterAddressing                                // Single register (Vx)
	RegisterValueAddressing                           // Register + byte (LD Vx, byte)
	RegisterRegisterAddressing                        // Two registers (LD Vx, Vy)
	RegisterRegisterNibbleAddressing                  // Two registers + nibble (DRW Vx, Vy, nibble)
	RegisterDTAddressing                              // Register from delay timer (LD Vx, DT)
	RegisterKAddressing                               // Register from key (LD Vx, K)
	RegisterIndirectIAddressing                       // Register from [I] (LD Vx, [I])
	DTRegisterAddressing                              // Delay timer from register (LD DT, Vx)
	STRegisterAddressing                              // Sound timer from register (LD ST, Vx)
	FRegisterAddressing                               // Font location (LD F, Vx)
	BRegisterAddressing                               // BCD representation (LD B, Vx)
	IAbsoluteAddressing                               // I register + address (LD I, addr)
	IRegisterAddressing                               // I register + register (ADD I, Vx)
	IIndirectRegisterAddressing                       // [I] from register (LD [I], Vx)
)
