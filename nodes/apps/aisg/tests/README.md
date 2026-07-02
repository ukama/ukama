# AISG Tests

## Protocol-only tests

```bash
make -C tests/protocol test
# or from aisg root
make protocol-test
```

These tests compile HDLC/XID/RETAP helpers with local stubs and do not require
the full Ukama platform tree.

## Strict RET emulator ladder

```bash
make -C ctrl
make -C emu
tests/scripts/run-ret-emu-ladder.sh
```

Keep logs:

```bash
AISG_KEEP_LOGS=1 tests/scripts/run-ret-emu-ladder.sh
```

## Negative strict-emulator checks

```bash
tests/scripts/run-ret-emu-negative.sh
```

## Real hardware ladder

Read-only first:

```bash
AISG_TTY=/dev/ttyUSB0 tests/scripts/run-real-hw-ladder.sh
```

Movement only after read-only passes:

```bash
AISG_TTY=/dev/ttyUSB0 AISG_MOVE=1 AISG_CONFIG_BLOB=/path/to/vendor.cfg \
  tests/scripts/run-real-hw-ladder.sh
```
