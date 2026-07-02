# AISG Stack

This directory contains the Ukama AISG node stack.

- `aisgd`: northbound Ukama REST/API app
- `ctrl`: low-level AISG controller daemon
- `emu`: emulator

Supported protocol scope:

```text
AISG v2.0, following TS 25.461 / 25.462 / 25.463 on conflict
one RS485 bus
one single-antenna RET device
device type 0x01 only
```

Out of scope:

```text
multi-antenna RET device type 0x11
TMA device type 0x02
multiple RETs on one bus
daisy-chain discovery
firmware/software download mode
```

Emulator modes:

```text
aisg-emu --mode contract
  Existing controller-contract emulator for fast aisgd/API testing.

aisg-emu --mode ret
  Strict single-antenna RET secondary emulator over PTY/serial.
  Use this to validate real AISG/HDLC/RETAP behavior before hardware.
```

Quick checks:

```bash
make protocol-test
make -C ctrl
make -C emu
tests/scripts/run-ret-emu-ladder.sh
```

See `docs/test-plan.md` and `RUN.md` for full emulator and hardware ladders.
