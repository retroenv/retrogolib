# Z80 Gap Closure Plan

Based on the cross-language comparison, this plan addresses the identified gaps while keeping
the existing instruction test suites (zexall, singlestep) working without requiring the new
interfaces.

## Design Principle

The CPU package currently has two host interfaces:
- `Memory` (required) -- for memory access
- `IOHandler` (optional, via `WithIOHandler`) -- for port I/O

The gaps require additional host interaction (interrupt data bus, T-state callbacks, RETI
notification). Rather than adding more optional handlers, we introduce a single `Bus` interface
that consolidates all host communication. The existing `Memory` + `IOHandler` approach continues
to work for simple use cases and tests.

---

## Phase 1: Unified Bus Interface

### Problem
Currently `Memory` and `IOHandler` are separate, and I/O uses 8-bit port addresses when the
real Z80 puts a full 16-bit address on the bus (high byte = register A or B depending on
instruction).

### Change

```go
// Bus provides the full hardware interface for a Z80 system.
// Implementations handle memory, I/O ports, and interrupt acknowledgment.
// For simple use cases (tests, basic emulation), use Memory + WithIOHandler instead.
type Bus interface {
    Memory

    // ReadPort reads from an I/O port. The full 16-bit address is provided
    // because the Z80 places register data on the upper address lines:
    // - IN A,(n):     address = A<<8 | n
    // - IN r,(C):     address = B<<8 | C
    // - INI/IND/etc:  address = B<<8 | C (after B decrement for some)
    ReadPort(address uint16) uint8

    // WritePort writes to an I/O port with full 16-bit address.
    WritePort(address uint16, value uint8)

    // IRQData returns the byte placed on the data bus during interrupt acknowledge.
    // For IM 0, this is an instruction opcode (typically RST 38h = 0xFF).
    // For IM 2, this is the low byte of the interrupt vector table address.
    // Only called when an IRQ is being serviced.
    IRQData() uint8

    // OnRETI is called when a RETI instruction executes.
    // Hardware (e.g., Z80 PIO/CTC daisy chain) monitors the bus for RETI
    // to manage interrupt priority.
    OnRETI()
}
```

### Constructor

```go
// NewWithBus creates a Z80 CPU with a full bus interface.
func NewWithBus(bus Bus, options ...Option) (*CPU, error)
```

The existing `New(memory Memory, ...)` constructor remains and wraps Memory + IOHandler into
an internal adapter that returns 0xFF for IRQData and no-ops for OnRETI -- exactly matching
current behavior.

### Internal Adapter

```go
// legacyBusAdapter wraps Memory + IOHandler into a Bus for backward compatibility.
type legacyBusAdapter struct {
    Memory
    ioHandler IOHandler
}

func (a *legacyBusAdapter) ReadPort(address uint16) uint8 {
    if a.ioHandler != nil {
        return a.ioHandler.ReadPort(uint8(address))
    }
    return 0xFF
}

func (a *legacyBusAdapter) WritePort(address uint16, value uint8) {
    if a.ioHandler != nil {
        a.ioHandler.WritePort(uint8(address), value)
    }
}

func (a *legacyBusAdapter) IRQData() uint8  { return 0xFF } // RST 38h
func (a *legacyBusAdapter) OnRETI()          {}
```

### Migration
- CPU internally stores `bus Bus` instead of separate `memory Memory`
- All `c.memory.Read/Write` calls unchanged (Bus embeds Memory)
- All `c.readPort(c.C)` / `c.writePort(c.C, val)` calls updated to pass
  full 16-bit address: `c.bus.ReadPort(c.bc())`, `c.bus.WritePort(c.bc(), val)`
- `c.opts.ioHandler` nil checks replaced by calling `c.bus.ReadPort/WritePort` directly
- Tests continue using `New(memory)` which internally creates the adapter

### Impact on Tests
None. `New(memory)` still works. zexall/singlestep tests don't need changes.

---

## Phase 2: Full IM 0 Interrupt Mode

### Problem
Current IM 0 is hardcoded to RST 38h. Real IM 0 reads an instruction from the data bus
and executes it. In practice this is almost always RST n, but the architecture should support
the general case.

### Change
In `handleInterrupts()`, replace:

```go
case 0:
    // Simplified: assumes RST 38H instruction on data bus
    c.PC = 0x0038
    c.cycles += 13
```

With:

```go
case 0:
    dataBusValue := c.bus.IRQData()
    // Most common case: RST instruction (0xC7, 0xCF, 0xD7, 0xDF, 0xE7, 0xEF, 0xF7, 0xFF)
    if dataBusValue&0xC7 == 0xC7 {
        // RST instruction: extract vector from bits 3-5
        vector := uint16(dataBusValue & 0x38)
        c.PC = vector
        c.MEMPTR = vector
        c.cycles += 13
    } else {
        // General case: execute the instruction from the data bus.
        // In practice this is extremely rare. For now, treat as RST 38h.
        c.PC = 0x0038
        c.cycles += 13
    }
```

This handles the real-world case (RST n from any device) while keeping a safe fallback for
the unlikely general case. The general case (arbitrary multi-byte instructions) is extremely
rare in real hardware and can be added later if needed.

### Impact on Tests
None. Default `IRQData()` returns 0xFF (RST 38h), matching current behavior.

---

## Phase 3: IM 2 Vector from Data Bus

### Problem
Current IM 2 reads the vector low byte from address 0xFFFF instead of from the data bus.

### Change
In `handleInterrupts()`, replace:

```go
case 2:
    // Simplified: reads vector low byte from 0xFFFF instead of data bus
    vector := uint16(c.I)<<8 | uint16(c.memory.Read(0xFFFF))
    c.PC = c.memory.ReadWord(vector)
    c.cycles += 19
```

With:

```go
case 2:
    vectorLow := c.bus.IRQData()
    vectorAddr := uint16(c.I)<<8 | uint16(vectorLow)
    c.PC = c.bus.ReadWord(vectorAddr)
    c.MEMPTR = c.PC
    c.cycles += 19
```

### Impact on Tests
None. Default `IRQData()` returns 0xFF, so vector becomes `I<<8 | 0xFF`. The zexall/singlestep
tests don't test IM 2 interrupt acknowledgment.

---

## Phase 4: RETI/RETN Notification

### Problem
Real Z80 hardware monitors the bus for the RETI instruction sequence (ED 4D) to manage
interrupt daisy-chain priority. Our emulator executes RETI but doesn't notify the host.

### Change
In the RETI handler, add:

```go
func edReti(c *CPU) error {
    c.PC = c.pop16()
    c.iff1 = c.iff2
    c.bus.OnRETI()
    return nil
}
```

### Impact on Tests
None. Default `OnRETI()` is a no-op.

---

## Phase 5: LD A,{I|R} Interrupt Bug

### Problem
On Zilog NMOS Z80s, if a maskable interrupt is accepted during LD A,I or LD A,R, the P/V
flag is reset to 0 instead of reflecting IFF2. This is a documented silicon bug.

### Change
Add a flag tracking whether the last instruction was LD A,I or LD A,R:

```go
// In CPU struct:
lastWasLdAIR bool
```

In `handleInterrupts()`, when servicing an IRQ:

```go
if c.triggerIrq && c.iff1 {
    // Zilog NMOS bug: LD A,I and LD A,R P/V flag reset on interrupt
    if c.lastWasLdAIR {
        c.Flags.P = false
    }
    // ... existing interrupt handling ...
}
```

In `Step()`, before executing each instruction:
```go
c.lastWasLdAIR = false
```

In `edLdAI` and `edLdAR` handlers:
```go
c.lastWasLdAIR = true
```

### Impact on Tests
None. The flag is only checked during interrupt handling, which zexall doesn't exercise
mid-instruction.

---

## Phase 6: Optional T-State Callbacks (Future)

This is the largest change and should be done only if a specific system (e.g., ZX Spectrum
48K with memory contention) requires it.

### Approach
Add an optional `CycleObserver` interface:

```go
// CycleObserver receives T-state notifications for cycle-exact emulation.
// This is optional -- most systems don't need sub-instruction timing.
type CycleObserver interface {
    // OnTick is called at each machine cycle boundary with the cycle type
    // and address on the bus. This enables memory contention modeling.
    OnTick(cycleType CycleType, address uint16)
}

type CycleType uint8

const (
    CycleOpcodeFetch CycleType = iota // M1 cycle: opcode fetch
    CycleMemRead                       // Memory read
    CycleMemWrite                      // Memory write
    CycleIORead                        // I/O port read
    CycleIOWrite                       // I/O port write
    CycleInternal                      // Internal operation (no bus activity)
)
```

This would be set via `WithCycleObserver(observer)` option. When nil (the default), all
callback sites are skipped with a simple nil check -- zero overhead for non-contention
systems.

### Impact
This requires instrumenting every memory/IO access point in the emulation code. It's
invasive but mechanical. The nil-check approach means no performance cost when not used.

This phase is **not recommended** until there's a concrete system that needs it.

---

## Summary

| Phase | Change | Breaking | Test Impact | Status |
|-------|--------|----------|-------------|--------|
| 1 | Bus interface + 16-bit I/O | No (adapter) | None | Done |
| 2 | Full IM 0 via IRQData() | No | None | Done |
| 3 | IM 2 vector from data bus | No | None | Done |
| 4 | RETI notification | No | None | Done |
| 5 | LD A,{I\|R} interrupt bug | No | None | Done |
| 6 | T-state callbacks | No (optional) | None | Planned |

All phases are backward-compatible. The existing `New(memory Memory, ...)` constructor
and all tests continue to work unchanged. New features are only active when the host
provides a `Bus` implementation.
