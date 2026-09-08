# Z80 Gap Closure Plan

Updated 2026-09-07. The bus interface, interrupt vectors, RETI notification, NMOS
LD A,I/R quirk, and ED mirrors are implemented and covered by tests. Interrupt
handling now uses one path for `Step` and `CheckInterrupts`. T-state callbacks
remain deferred, and arbitrary IM 0 instruction execution remains unsupported.

| Phase | Change | Status |
|-------|--------|--------|
| 1 | Unified Bus and 16-bit I/O addresses | Implemented; block output address ordering corrected |
| 2 | Device-supplied IM 0 restart vector | Implemented for all eight RST opcodes; other opcodes retain RST 38h fallback |
| 3 | IM 2 vector byte from the data bus | Implemented in both interrupt entry points |
| 4 | RETI notification | Implemented; RETN mirrors do not notify |
| 5 | LD A,I/R interrupt parity quirk | Implemented for NMOS behavior at instruction boundaries |
| 6 | T-state callbacks and contention | Deferred until a concrete system requires them |
| 7 | Undocumented ED mirrors | Already implemented; dedicated execution tests added |

## Bus and Compatibility

`NewWithBus(bus, options...)` uses the full `Bus` contract in
`arch/cpu/z80/memory.go`: memory access, `ReadPort(uint16)`,
`WritePort(uint16, uint8)`, `IRQData()`, and `OnRETI()`.

`New(memory, options...)` retains the legacy adapter. `WithIOHandler` receives only
the port's low eight bits; without a handler, reads return `FF` and writes are ignored.
The adapter supplies `FF` for interrupt data and ignores RETI notifications.
`WithIOHandler` applies to `New`; `NewWithBus` uses the bus's port methods.

| Operation | Full port address |
|-----------|-------------------|
| IN A,(n) / OUT (n),A | A before the transfer in the high byte; n in the low byte |
| IN r,(C) / OUT (C),r | B:C before the transfer |
| INI / IND / INIR / INDR | B:C before decrementing B |
| OUTI / OUTD / OTIR / OTDR | B:C after decrementing B |

Block output previously sent the old B value. The single-step runner now checks
full port addresses, values, direction, and ordering; its former eight-bit adapter
could not detect this error. IN F,(C) and OUT (C),0 also use the full address, with
OUT (C),0 following NMOS behavior.

Bus callbacks and execution hooks run under the CPU lock and must not call locking
CPU methods. Direct register access and interrupt configuration must be serialized
with execution. `State` is a debugging snapshot and does not contain every internal
latch needed for save/restore.

## Interrupt Acceptance

The duplicate interrupt implementation has been removed. `CheckInterrupts` uses the
same acceptance logic as `Step`, including bus-supplied vectors and the LD A,I/R quirk.
An accepted interrupt consumes its own step: the first handler instruction executes
on the following `Step` call. HALT checks pending interrupts before idling; accepted
interrupts release it. Idle HALT cycles and interrupt acknowledgments increment R's
low seven bits while preserving bit 7.

NMI has priority, preserves IFF2 across nested NMIs, and leaves a pending IRQ queued.
EI sets the enable flip-flops immediately but inhibits IRQ acceptance until the next
instruction completes. Consecutive EI instructions extend this delay, DI prevents
acceptance, and EI followed by HALT permits the pending IRQ to release HALT.
These behaviors follow the interrupt-response and HALT sections of the
[Zilog Z80 CPU User Manual](https://www.zilog.com/docs/z80/um0080.pdf).

| Event | Target | T-states |
|-------|--------|---------:|
| NMI | 0066 | 11 |
| IM 0 IRQ | RST target from IRQData; 0038 for the retained non-RST fallback | 13 |
| IM 1 IRQ | 0038 | 13 |
| IM 2 IRQ | Little-endian word at `(I << 8) | IRQData()` | 19 |

IM 0/IM 2 sample `IRQData` once before pushing the interrupted PC. IM 2 supports
odd vector addresses and wrapping the vector word from FFFF to 0000. IM 1 ignores
the device data byte. MEMPTR records the handler address.

The original Phase 2 title claimed “Full IM 0,” but both its proposal and implementation
only supported RST instructions. Arbitrary instructions and multi-byte device-supplied
sequences require further design; the existing RST 38h fallback is explicitly retained.

Only ED 4D emits `OnRETI`, after restoring the PC and IFF1 from IFF2. RETN and all
its mirrors restore state without notifying daisy-chain hardware. An accepted IRQ
immediately after LD A,I or LD A,R clears P/V; an intervening instruction ends this
quirk's window. This is instruction-boundary NMOS modeling, not sub-instruction sampling.

## ED Mirrors and Prefix Behavior

The opcode table already contained these variants. Dedicated unit tests now execute
every encoding and check results, flags or interrupt state, PC, and timing:

| Operation | ED opcode bytes |
|-----------|-----------------|
| NEG | 44, 4C, 54, 5C, 64, 6C, 74, 7C |
| IM 0 | 46, 4E, 66, 6E |
| IM 1 | 56, 76 |
| IM 2 | 5E, 7E |
| RETN | 45, 55, 5D, 65, 6D, 75, 7D |
| RETI | 4D |
| IN F,(C) / OUT (C),0 | 70 / 71 |

The former mirror list omitted NEG ED 74 and the IM 1/2 mirrors. These encodings
also occur in the external single-step corpus; the former claim that single-step
tests did not exercise them was incorrect.

The baseline single-step run failed four files: DD/FD-prefixed SCF and CCF. Ignored
index prefixes now clear Q before SCF/CCF derive the undocumented X/Y flags, matching
the corpus. A focused regression test reproduces the former flag loss.

## Validation

The [SingleStepTests Z80 corpus](https://github.com/SingleStepTests/z80) at checkout
`ebe1875d48f374bcfd4b505d8eb8ee751568b5f7` contains 1,604 files and 1,604,000 vectors.
The complete run passes with register, flag, memory, and full port-transaction checks.
The runner does not compare the final Q latch, memory bus cycles, or T-state timing,
and these vectors do not replace the dedicated interrupt tests.

```sh
# Download only when the corpora are absent.
make -C testdata z80

go test ./arch/cpu/z80 -run '^TestSingleStep$' -count=1 -timeout 180s
go test ./arch/cpu/z80 -run '^TestZex(doc|all)$' -count=1 -timeout 45m

go fmt ./...
make lint
make test
```

`go fmt ./...` completed, `make lint` reported zero issues from both linters, and
`make test` passed the repository's short tests with the race detector. ZEXDOC passed
all 67 groups in 397 seconds; ZEXALL passed all 67 groups, including undocumented
flags, in 318 seconds. The external corpus and exerciser runs used the commands above
without the race detector.
