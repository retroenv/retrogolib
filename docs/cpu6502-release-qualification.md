# 6502 release qualification

The emulator needs to produce both the right answer and the right instruction
cycle count before it can measure compiler performance. The ordinary tests are
fast and need no downloads. Release qualification additionally requires two
independent, pinned test collections; missing, altered, empty or wrong-revision
data fails the enabled gate instead of silently skipping it.

Provision the data in a new directory (about 1.3 GB for the NMOS vectors):

```sh
qualification_data=$(mktemp -d)
git clone --filter=blob:none --no-checkout https://github.com/SingleStepTests/65x02.git "$qualification_data/65x02"
git -C "$qualification_data/65x02" sparse-checkout set 6502
git -C "$qualification_data/65x02" checkout --detach 2f6980a2d95757486c7bee24355c360e40e2a224
git clone --filter=blob:none --no-checkout https://github.com/Klaus2m5/6502_65C02_functional_tests.git "$qualification_data/dormann"
git -C "$qualification_data/dormann" checkout --detach 7954e2dbb49c469ea286070bf46cdd71aeb29e4b
CPU6502_TESTDATA="$qualification_data/65x02" \
CPU6502_DORMANN_TESTDATA="$qualification_data/dormann" \
make test-6502-release
```

The gate verifies each checkout's commit and clean state, then requires:

- All 151 legal NMOS opcodes, 10,000 vectors each: **1,510,000 cases**.
  Each case checks the initial opcode, final registers and memory, and the
  instruction's cycle count against the external trace length. Reset cycles
  are excluded. KIL/JAM and other unofficial instructions are not legal NMOS
  instructions and are not included in this count.
- Dormann's NMOS functional binary reaching its success loop at `$3469`.
- The separate 65C02 extended-opcode regression reaching `$24F1`.

The target enables race detection and has a finite ten-minute deadline. It
must not use `-short`. Ordinary `make test` uses a configurable 60-second
per-package deadline because race-enabled CPU68000 tests need more headroom
than the previous ten-second limit; no assertions or cases are removed.

This qualifies observed instruction state and instruction cycle totals, not
every bus access or undocumented opcode. It is not qualification of a complete
NES: video, audio, DMA, interrupts in a running console, and cartridge hardware
need their own tests. The 65C02 result is a regression check, not a substitution
for the NMOS corpus.

Sources: [SingleStepTests/65x02](https://github.com/SingleStepTests/65x02),
[Klaus Dormann's functional tests](https://github.com/Klaus2m5/6502_65C02_functional_tests).
The immutable revisions above, rather than moving default branches, define the
release evidence.
