# Motorola 68000 Gap Closure Plan

This plan addresses the identified accuracy gaps while keeping the existing test suite
working.

## Design Principle

The 68000 implementation covers instruction families, effective addressing, privilege
modes, and basic exception handling, but its SingleStepTests suite still fails. Full
instruction accuracy is not established. Known gaps include exception fidelity, alignment,
stack-pointer reporting, BCD/flag behavior, and timing. The phases below aim to preserve the
existing Memory/Bus interfaces while correcting behavior.

---

## Phase 1: Address Error Checking (High Priority)

### Problem
Word and long accesses at odd addresses should trigger an Address Error exception (vector 3).
The real 68000 generates a special stack frame with the faulting address, access type
(read/write), and function code. retrogolib silently allows misaligned access.

### Change
Add alignment validation to `readEA` and `writeEA` for word/long accesses:

```go
func (c *CPU) checkAlignment(address uint32, size OperandSize) error {
    if size != SizeByte && address&1 != 0 {
        return c.addressError(address, size, /* read/write info */)
    }
    return nil
}

func (c *CPU) addressError(address uint32, size OperandSize, info accessInfo) error {
    // Push Type 2 stack frame (see Phase 2)
    return c.processException(VectorAddressError)
}
```

### Impact on Tests
Existing tests should not use misaligned word/long accesses. Add new tests for:
- Word read at odd address triggers exception
- Long write at odd address triggers exception
- Exception vector is loaded from correct address
- Stack frame contains faulting address

---

## Phase 2: Exception Stack Frames (Medium Priority)

### Problem
The 68000 has two stack frame formats:
- **Type 0 (6 bytes):** SR + PC. Used for most exceptions.
- **Type 2 (14 bytes):** SR + PC + instruction word + fault address + access info.
  Used for bus error and address error.

retrogolib only implements Type 0 frames. Type 2 frames are needed for proper bus/address
error recovery.

### Change
Extend `processException` to support frame types:

```go
type exceptionFrame struct {
    frameType   uint8   // 0 or 2
    faultAddr   uint32  // Address that caused the fault
    instrWord   uint16  // First word of faulting instruction
    accessType  uint16  // Read/Write, instruction/data, function code
}

func (c *CPU) processException(vector uint8) error {
    // Existing: push PC and SR (Type 0)
    // ...
}

func (c *CPU) processExceptionType2(vector uint8, frame exceptionFrame) error {
    // 1. Save current SR
    // 2. Set supervisor mode
    // 3. Push additional info (14 bytes total):
    //    - Access type and function code (word)
    //    - Fault address (long)
    //    - Instruction word (word)
    //    - Status register (word)
    //    - Program counter (long)
    // 4. Load PC from vector table
}
```

### Impact on Tests
None for existing tests. Add new tests verifying Type 2 frame layout for address errors.

---

## Phase 3: Bus Error Support (Medium Priority)

### Problem
Bus errors (vector 2) occur when external hardware signals an invalid memory access. The Bus
interface has no mechanism for memory implementations to signal bus errors.

### Change
Extend the Memory or Bus interface to allow error signaling:

```go
// Option A: Error return from Memory methods
type MemoryWithErrors interface {
    Memory
    ReadByteChecked(address uint32) (uint8, error)
    WriteByteChecked(address uint32, value uint8) error
    // ... word and long variants
}

// Option B: Separate bus error trigger
type BusErrorHandler interface {
    // OnBusAccess is called for each bus access. Returns an error to trigger
    // a bus error exception.
    OnBusAccess(address uint32, size OperandSize, write bool) error
}
```

Option A is cleaner but changes the hot path. Option B keeps the existing Memory interface
unchanged and adds optional error checking. **Recommended: Option B** for backward
compatibility.

### Impact on Tests
None. Bus error checking is only active when the host provides a `BusErrorHandler`.

---

## Phase 4: Variable Cycle Timing (Low Priority)

### Problem
retrogolib assigns fixed cycle counts per instruction. The real 68000 has variable timing
based on addressing mode, operand size, and operand values.

### Change
Add EA cycle costs based on Motorola M68000UM Table 8-4:

```go
// eaCycles returns the additional cycles for an effective address calculation.
func eaCycles(mode AddressingMode, size OperandSize) int {
    // Data Register Direct: 0
    // Address Register Direct: 0
    // (An): 4/8 (word/long)
    // (An)+: 4/8
    // -(An): 6/10
    // d16(An): 8/12
    // d8(An,Xn): 10/14
    // (xxx).W: 8/12
    // (xxx).L: 12/16
    // d16(PC): 8/12
    // d8(PC,Xn): 10/14
    // #imm: 4/8
}
```

Also add data-dependent timing for:
- MULU: 38 + 2n cycles (n = number of 1-bits in source)
- MULS: 38 + 2n cycles (n = number of 01/10 bit transitions)
- DIVU: 140 cycles max (varies with quotient)
- DIVS: 158 cycles max (varies with quotient)

### Impact on Tests
Cycle count assertions in existing tests may need updating. Run tests first to identify
which are affected.

---

## Phase 5: Prefetch Queue Emulation (Low Priority, Optional)

### Problem
The real 68000 has a 2-word prefetch pipeline (IRC and IRD registers). Self-modifying code
and hardware register timing depend on prefetch behavior.

### Change
Model the 2-word prefetch pipeline:

```go
type CPU struct {
    // ... existing fields ...
    IRC uint16 // Instruction Register Capture (next word to decode)
    IRD uint16 // Instruction Register Decode (currently executing word)
}

// prefetch fetches the next instruction word into IRC.
func (c *CPU) prefetch() {
    c.IRC = c.bus.ReadWord(c.PC)
    c.PC += 2
}
```

### When to Implement
Only needed for:
- Amiga copper/blitter timing (critical for demos)
- Self-modifying code accuracy
- Certain copy protection schemes

**Not recommended** until a concrete Amiga or demo-accurate emulation use case requires it.

---

## Phase 6: SingleStepTests Integration

### Status: Runner Implemented, Conformance Failing

The runner exists in `arch/cpu/cpu68000/singlestep_test.go`, behind the `singlestep`
build tag. `make -C testdata cpu68000` downloads the vectors, and the root
`make test-integration` target includes this suite.

Run it directly with:

```sh
go test -tags singlestep ./arch/cpu/cpu68000 -run TestSingleStep -count=1 -timeout 120s
```

The 2026-09-07 review recorded 206,547 passing and 992 failing vectors among 207,539
executed. The runner stops each file after ten failures; this is not a full-corpus pass
rate. An overlay using the unchanged `dfc8a64` implementation produced the same counts
and failure patterns as the review fixes. See [the branch changelog](work-branch-changelog.md)
for the validation scope.

### Next Work

Use the existing failures to drive the fixes instead of adding another runner. Alongside
the exception phases, investigate active/saved stack-pointer reporting and BCD/flag
discrepancies. The runner compares architectural state and memory; it does not yet validate
the full bus transaction sequence or cycle timing.

---

## Summary

| Phase | Change | Breaking | Priority | Status |
|-------|--------|----------|----------|--------|
| 1 | Address error checking | No | High | Planned |
| 2 | Type 2 exception stack frames | No | Medium | Planned |
| 3 | Bus error support | No (optional interface) | Medium | Planned |
| 4 | Variable cycle timing | No (timing values change) | Low | Planned |
| 5 | Prefetch queue emulation | No | Low | Deferred |
| 6 | SingleStepTests integration | No | Medium | Runner implemented; conformance failures remain |

Phases 1-3 address exception fidelity. Phase 4 improves timing, and Phase 5 adds prefetch
behavior. The existing Phase 6 runner supplies regression evidence throughout this work;
passing it alone will not establish full bus-cycle accuracy.
