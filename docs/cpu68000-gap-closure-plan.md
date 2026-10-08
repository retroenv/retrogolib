# Motorola 68000 Gap Closure Plan

Implementation record: 2026-09-07. Validation instructions reviewed: 2026-10-02.
Review corrections applied: 2026-10-07 (see "Review Corrections").
Phases 1–4 are implemented. The Phase 6 runner now executes the
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

The 2026-09-07 implementation left the `Memory` and `Bus` interfaces unchanged.
Since 2026-10-07 `Memory` has only byte and word methods, because the CPU never
issued a long transfer; `BasicMemory` keeps `ReadLong` and `WriteLong` as host
conveniences. A host bus can additionally implement the optional capability in
`arch/cpu/cpu68000/memory_access.go`:

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
The table records the 2026-09-07 runs. It does not report a new test run.

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
clear; the implementation leaves both set. WinUAE and Musashi also set C/X to the
sign bit for a count that is not less than the operand width.

The two `e502 [ASL.b Q, D2]` cases, 1583 and 1761, change D2 from `CDFB7FBE` to
`2E5E4304` and from `417C7E7D` to `6461D390`, respectively. A byte operation must
preserve the upper 24 bits. The expected status words (`2713`, `271B`) also do not
match the result of the shift, so the two vectors are corrupt. The
`80ef [DIVU (d16, A7), D0] 5745` case expects saved PC `C00`, whereas the four-byte
instruction ends at `C04`; the user manual's section 6.3.5 specifies the following
instruction for this trap. It is the only zero-divisor vector in `DIVU.json.gz`, and
`DIVS.json.gz` has none, while the same generator saves the following instruction for
all 3,989 CHK traps and 4,095 TRAPV traps. The vector also expects N clear, whereas
[WinUAE's `divbyzero_special`](https://github.com/tonioni/WinUAE/blob/master/newcpu_common.cpp)
clears C, V, N, and Z and then sets N or Z from the high word of the dividend
(`A18E` here, so N is set) on the 68000. That routine is the independent reference
confirmation for the implementation; unit tests preserve the documented behavior.
Corrections to the corpus remain upstream work before full conformance can be claimed.

## Review Corrections (2026-10-07)

A review that compared the decoder with the official opcode map
(`testdata/cpu68000/680x0/map/68000.official.json`) and probed exception paths led
to these behavior changes. Each has a regression test.

- The decoder validates the addressing-mode class (data, memory, control, alterable)
  and size bits of every instruction. 7,935 of the 19,721 unassigned words, for
  example `JSR Dn`, `MOVE.B An,Dn`, `LEA Dn,An`, `ADDQ.B #n,An`, `AND.W An,Dn`, and
  memory shifts with bit 11 set, previously executed as instructions. They now decode
  to ILLEGAL. `TestDecoderMatchesOfficialOpcodeMap` checks all 65,536 words against
  the map: every `None` entry decodes to ILLEGAL and every other entry decodes to the
  mapped mnemonic family.
- Unassigned words raise the illegal instruction exception (vector 4) with the
  instruction address as saved PC instead of returning a host error.
  `ErrUnsupportedOpcode` and other unused error values are removed.
- Level 7 interrupts are edge-sensitive. The CPU accepts level 7 when the bus level
  rises to 7, when `TriggerIRQ` queues it, or when the mask drops below 7 while the
  bus still holds level 7. A bus that held level 7 previously re-entered the handler
  on every step. The acknowledged vector number is masked to eight bits.
- DIVS by zero clears N, V, and C and sets Z, as WinUAE's `divbyzero_special` does
  for the 68000. DIVU by zero keeps its high-word N/Z derivation.
- A trace exception follows TRAP, TRAPV, CHK, and zero divide, with the handler
  address as saved PC (WinUAE `exception_trace`). Illegal, privilege, and line A/F
  exceptions and interrupts drop the trace, as before.
- Scc reads its memory destination before the write, like CLR and MOVE from SR.
- The memory read, write, and read-modify-write instruction sets include the
  immediate, decimal, shift, stack, and read-before-write instructions that the
  handlers transfer through memory.
- `doc.go` documents that the hierarchical decoder has no encoder, so the package
  offers no instruction-to-opcode reverse mapping; the official-map test provides
  forward completeness.

The full corpus run after these changes (2026-10-07) reproduced 996,321 passed and
3,739 failed out of 1,000,060 with the same five discrepancy classes, so the decoder
validation introduced no new failures. The corpus does not exercise the trace bit,
level 7 interrupts, DIVS by zero, or unassigned opcode words.

## Validation Commands

The 2026-09-07 implementation passed `go fmt ./...`, both linters in `make lint` with
zero issues, and all repository short tests under `make test` with the race detector.
Self-review confirmed that CPU memory transfers go through the checked access helpers;
`git diff --check` passed. The final full-corpus run reproduced the counts above.

Run from the repository root. To use the recorded corpus revision without
changing an existing checkout:

```sh
set -eu
cpu68000_data=$(mktemp -d)
git clone --filter=blob:none --no-checkout https://github.com/SingleStepTests/680x0.git "$cpu68000_data/680x0"
git -C "$cpu68000_data/680x0" sparse-checkout set 68000/v1
git -C "$cpu68000_data/680x0" checkout --detach e0d5ece9670205cc84a0101081837deb446f86a3
CPU68000_TESTDATA="$cpu68000_data/680x0/68000/v1" \
go test -v -tags singlestep ./arch/cpu/cpu68000 -run '^TestSingleStep$' -count=1 -timeout 120s
```

`CPU68000_TESTDATA` must name the directory that contains the `.json.gz` files,
not the repository root. The runner skips when the directory or test files are
absent. It does not enforce the corpus revision or a clean checkout. Inspect
the verbose output and vector totals before reporting a result.

`make -C testdata cpu68000` downloads or updates the default corpus. It does
not pin the revision above. Its default runner path is
`testdata/cpu68000/680x0/68000/v1`, and the official opcode map used by
`TestDecoderMatchesOfficialOpcodeMap` is `testdata/cpu68000/680x0/map/68000.official.json`.
Both tests skip when their files are absent.

`make test` runs short tests with the race detector. The tagged integration command
currently fails with the reference discrepancies above. Full prefetch and bus-sequence
validation remain future work; passing architectural-state vectors alone would not
establish cycle accuracy.
