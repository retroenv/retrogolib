# Motorola 68000 Gap Closure Plan

Updated 2026-09-07. Phases 1–4 are implemented. The Phase 6 runner now executes the
entire corpus, with the remaining reference discrepancies recorded below. Phase 5
remains deferred until a concrete use case requires instruction prefetch emulation.

| Phase | Change | Status |
|-------|--------|--------|
| 1 | Alignment checks for CPU memory accesses | Implemented |
| 2 | Original 68000 bus/address-error frames | Implemented |
| 3 | Optional bus error callback | Implemented |
| 4 | Addressing-mode and operand-dependent timing | Implemented; complete bus timing remains unverified |
| 5 | Two-word instruction prefetch queue | Deferred |
| 6 | SingleStepTests integration | Full-corpus runner implemented; 3,739 reference mismatches remain |

## Phases 1–3: Memory Faults and Exceptions

CPU instruction fetches, extension words, operands, stack operations, and vector reads
use checked memory access. Odd word and long addresses raise vector 3 before reaching
the host memory implementation; byte accesses remain legal at odd addresses. Long
accesses use two word transfers on the 16-bit bus, preserving completed transfers if
the second word faults. Physical bus addresses wrap at 24 bits while address registers
and calculated effective addresses retain all 32 bits.

The original 68000 uses an ordinary six-byte SR/PC frame and a fourteen-byte bus/address
error frame. These frames have **no format word**; the former plan's “Type 0/Type 2”
terminology was incorrect for this processor. The fault frame starts at the new SSP:

| Offset | Size | Contents |
|--------|------|----------|
| 0 | Word | Special status word: read/write, instruction state, function code |
| 2 | Long | Fault address |
| 6 | Word | Instruction register |
| 8 | Word | Saved SR |
| 10 | Long | Saved PC |

The saved PC reflects the instruction's fault point and may differ from its starting
address. A guest handler must account for the extra eight bytes before using RTE;
this is not a 68010 restartable frame. Frame layout and exception behavior follow
[MC68000 User's Manual, sections 6.3.5 and 6.3.9–6.3.10](https://www.nxp.com/docs/en/reference-manual/MC68000UM.pdf).

Existing `Memory` and `Bus` interfaces are unchanged. A host bus can additionally
implement the optional capability in `arch/cpu/cpu68000/memory_access.go`:

```go
type BusErrorHandler interface {
    OnBusAccess(address uint32, size OperandSize, write bool) error
}
```

The callback runs before each CPU byte or word transfer, including each half of a long
access. Returning an error raises vector 2 and aborts the current instruction before
later side effects. It runs under the CPU lock, so it must not call locking CPU methods.
Direct host calls to `Memory` or the returned `Bus` bypass CPU checking.

A further fault during bus/address-error handling halts the CPU without recursively
building another frame. `Resume` cannot release this halt. The new `Reset` method
reloads SSP/PC from vectors 0/1, clears pending interrupts and execution halt/stop state,
and resets the cycle counter to 40. Reset-vector failures return an error wrapping
`ErrBusError` and the host cause. Construction through reset vectors uses the same
checks but retains the existing zero initial cycle count. The guest RESET instruction
continues to notify external devices through `Bus.OnReset`.

## Phase 4: Instruction Timing

`timing.go` combines instruction costs with effective-address costs for operand size
and addressing mode. It covers MOVE, arithmetic and logical operations, unary and bit
operations, branches and conditions, shifts, control transfers, and MOVEM/MOVEP.
MULU/MULS count source bits or transitions; DIVU/DIVS account for dividend, divisor,
and overflow. MOVEM adds the cost of the registers actually transferred. Operands and
register masks are read once, so timing calculation does not duplicate device reads.

Timing tables follow [MC68000 User's Manual, section 8](https://www.nxp.com/docs/en/reference-manual/MC68000UM.pdf).
The division calculations also use the algorithm documented in
[WinUAE's division timing routines](https://github.com/tonioni/WinUAE/blob/master/newcpu_common.cpp).
Address/bus errors add exception overhead to the tracked work preceding the fault.

These are instruction timing calculations, not full bus-cycle emulation. Prefetch,
wait states, and the precise internal/bus sequence at every fault point remain outside
this implementation's validated scope. Register bit-operation timings use the manual's
listed maximum values. Unit tests check representative addressing modes, branch outcomes,
shift counts, transfers, multiply/divide operands, exceptions, and fault timing.

## Phase 6: Full-Corpus Validation and Related Fixes

The existing runner in `arch/cpu/cpu68000/singlestep_test.go` remains behind the
`singlestep` build tag. It now limits diagnostics to ten failures per file while still
executing every vector. It compares registers, status, PC, stack pointers, and listed
memory contents; it does not compare the complete bus transaction sequence, prefetch
state, or cycle counts.

The corpus is [SingleStepTests/680x0](https://github.com/SingleStepTests/680x0), checkout
`e0d5ece9670205cc84a0101081837deb446f86a3`, directory `68000/v1`. Both runs below use the
same full-corpus runner; the baseline uses an overlay of the unchanged `754af72` CPU.

| Implementation | Passed | Failed | Total |
|----------------|-------:|-------:|------:|
| Before gap closure (`754af72`) | 660,384 | 339,676 | 1,000,060 |
| After gap closure | 996,321 | 3,739 | 1,000,060 |

Earlier counts in the branch changelog stopped each file after ten failures and are
not comparable to these complete runs. Corrections driven by the corpus include:

- Synchronizing the active stack pointer with exported SSP/USP after each step.
- Fixing overlapping MOVEP/bit-op and ABCD/SBCD/logical-op decoding, MOVEP direction,
  MOVEM predecrement register order, and preservation of complete effective addresses.
- Correcting extended carry/borrow, decimal adjustments, CHK flags, reserved SR bits,
  CMPM byte increments through A7, and LINK with A7 as its operand.
- Preserving partial register and memory updates at faults, including MOVE/MOVEM
  predecrement ordering, postincrement writes, and JSR target checking.
- Reading memory before CLR and MOVE-from-SR writes, checking privilege before immediate
  SR operands, and delaying tracing until the instruction after T becomes enabled.
- Saving the appropriate exception PC and distinguishing line A from line F traps.

### Remaining Reference Discrepancies

The suite still fails; no vectors or comparisons are skipped to obtain a passing result.

| File | Failures | Observed discrepancy |
|------|---------:|----------------------|
| ASR.b | 1,642 | Negative value, register count greater than 8: reference clears C/X after shifting sign-extension bits |
| ASR.w | 1,063 | Same discrepancy with count greater than 16 |
| ASR.l | 1,031 | Same discrepancy with count greater than 32 |
| ASL.b | 2 | Reference changes the destination register's upper 24 bits |
| DIVU | 1 | Zero-divisor reference saves the instruction's start PC instead of the following PC; N also differs |

All 3,736 ASR flag mismatches share the condition above. Motorola specifies sign
extension and the last shifted-out bit in C/X; see ASL/ASR in the
[M68000 Programmer's Reference Manual, pages 4-21–4-22](https://www.nxp.com/docs/en/reference-manual/M68000PRM.pdf).
For example, `ea23 [ASR.b D5, D3] 8` shifts negative byte `F3` by 12 and expects C/X
clear; the implementation leaves both set.

The two `e502 [ASL.b Q, D2]` cases, 1583 and 1761, change D2 from `CDFB7FBE` to
`2E5E4304` and from `417C7E7D` to `6461D390`, respectively. A byte operation must
preserve the upper 24 bits. The `80ef [DIVU (d16, A7), D0] 5745` case expects saved
PC `C00`, whereas the four-byte instruction ends at `C04`; the user manual's section
6.3.5 specifies the following instruction for this trap. Unit tests preserve the
documented behavior. These discrepancies need independent reference confirmation or
upstream corrections before full corpus conformance can be claimed.

## Validation Commands

The final implementation passed `go fmt ./...`, both linters in `make lint` with
zero issues, and all repository short tests under `make test` with the race detector.
Self-review confirmed that CPU memory transfers go through the checked access helpers;
`git diff --check` passed. The final full-corpus run reproduced the counts above.

```sh
go fmt ./...
make lint
make test
make -C testdata cpu68000
go test -tags singlestep ./arch/cpu/cpu68000 -run '^TestSingleStep$' -count=1 -timeout 120s
```

`make test` runs short tests with the race detector. The tagged integration command
currently fails with the reference discrepancies above. Full prefetch and bus-sequence
validation remain future work; passing architectural-state vectors alone would not
establish cycle accuracy.
