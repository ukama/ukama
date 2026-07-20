# AISG Runbook

Ukama AISG supports one AISG v2 / TS 25.xxx single-antenna RET device.
AISG v2.0 is the umbrella standard; TS 25.461 / TS 25.462 / TS 25.463 win
whenever there is conflict.

Use this order:

```text
1. Protocol golden tests
2. Strict RET emulator ladder
3. Read-only hardware ladder
4. Optional movement hardware ladder
```

## 1. Build

Expected repo layout:

```text
nodes/
  ukamaOS/
  apps/
    aisg/
      aisgd/
      ctrl/
      emu/
```

From `nodes/apps/aisg`:

```bash
test -f ../../../ukamaOS/config.mk && echo "OK: ukamaOS config found"
make -C ctrl
make -C emu
make -C aisgd
```

## 2. Protocol-only golden tests

These tests use local stubs and do not need the full Ukama platform build.

```bash
make protocol-test
```

They validate HDLC, XID, RETAP framing, return codes, tilt endian, and config
segment size.

## 3. Strict RET emulator ladder

This validates `aisg-ctrl` against a standards-shaped single-antenna RET device.
It is the main software gate before touching hardware.

```bash
tests/scripts/run-ret-emu-ladder.sh
```

To keep logs:

```bash
AISG_KEEP_LOGS=1 tests/scripts/run-ret-emu-ladder.sh
```

The script starts:

```text
aisg-emu --mode ret --pty <tmp>/aisg-ret0
aisg-ctrl --backend raw-rs485 --device <tmp>/aisg-ret0
```

Then runs:

```text
get_status
scan/connect
get_info
get_alarm_status
get_tilt before calibration, expecting NotCalibrated
send_configuration_data
calibrate
get_tilt
set_tilt 0.5
get_tilt verify
```

Negative strict-emulator checks:

```bash
tests/scripts/run-ret-emu-negative.sh
```

This verifies the emulator does not respond to the old invalid scan and does
not accept RETAP before address assignment/SNRM.

## 4. Real hardware setup

Real setup:

```text
Linux machine
  -> USB-RS485 adapter exposed as /dev/ttyUSB0, /dev/ttyACM0, or /dev/serial/by-id/...
  -> AISG M16 harness
  -> real single-antenna RET device
```

Hardware checklist:

```text
Power: PSU +10–30 V to AISG Pin6 and PSU return to Pin7
RS485: A+ adapter terminal to AISG Pin3 (B/Vb); B- to Pin5 (A/Va)
GND/reference: adapter GND to AISG Pin4; do not bridge Pin4 to Pin7
Baud: start with 9600 8N1
Adapter: automatic TX/RX direction control preferred
Idle polarity: V(Pin3)-V(Pin5) must be positive, preferably above 200 mV
```

Stop ModemManager during testing:

```bash
sudo systemctl stop ModemManager 2>/dev/null || true
```

Give the user serial permissions:

```bash
sudo usermod -a -G dialout "$USER"
newgrp dialout
```

For one-time testing:

```bash
sudo chmod a+rw /dev/ttyUSB0
```

## 5. Read-only hardware ladder

Run only read-only operations first:

```bash
AISG_TTY=/dev/ttyUSB0 tests/scripts/run-real-hw-ladder.sh
```

The script runs:

```text
get_status
scan/connect
get_info
get_alarm_status
```

If this fails but the strict emulator passes, investigate:

```text
power/current
RS485 A/B polarity
RS485 reference/GND
USB-RS485 direction control
TX echo vs real RX
vendor config/pinout quirks
```

Logs are kept by default under `/tmp/aisg-real.*`.

## 6. Optional movement hardware ladder

Only run movement after read-only hardware tests pass.

```bash
AISG_TTY=/dev/ttyUSB0 \
AISG_MOVE=1 \
AISG_CONFIG_BLOB=/path/to/vendor-antenna.cfg \
tests/scripts/run-real-hw-ladder.sh
```

Movement sequence:

```text
send_configuration_data, if AISG_CONFIG_BLOB is set
calibrate
get_tilt
set_tilt 0.5
get_tilt verify
```

Do not repeatedly move the motor until the antenna tilt range and vendor config
file are confirmed.

## 7. Manual controller socket examples

Create a temporary `aisg-ctrl` config:

```bash
mkdir -p /tmp/aisg-real
cat >/tmp/aisg-real/aisg-ctrl.toml <<'EOF_CTRL'
[service]
socket = "/tmp/aisg-ctrl.sock"

[backend]
type = "raw-rs485"

[raw_rs485]
device = "/dev/ttyUSB0"
baud = 9600
EOF_CTRL
```

Start controller:

```bash
rm -f /tmp/aisg-ctrl.sock
./ctrl/aisg-ctrl -c /tmp/aisg-real/aisg-ctrl.toml -l TRACE \
  2>&1 | tee /tmp/aisg-real/aisg-ctrl.log
```

Send requests with `socat`.  `STDIO,ignoreeof` keeps the receive half open after
`printf` reaches EOF, which is required for operations such as `scan` that take
longer than an immediate status request:

```bash
printf '{"id":"status-1","type":"get_status","payload":{}}\n' \
  | socat -T 15 STDIO,ignoreeof UNIX-CONNECT:/tmp/aisg-ctrl.sock | jq

printf '{"id":"scan-1","type":"scan","payload":{}}\n' \
  | socat -T 15 STDIO,ignoreeof UNIX-CONNECT:/tmp/aisg-ctrl.sock | jq

printf '{"id":"info-1","type":"get_info","payload":{}}\n' \
  | socat -T 15 STDIO,ignoreeof UNIX-CONNECT:/tmp/aisg-ctrl.sock | jq

printf '{"id":"err-1","type":"get_alarm_status","payload":{}}\n' \
  | socat -T 15 STDIO,ignoreeof UNIX-CONNECT:/tmp/aisg-ctrl.sock | jq
```

## 8. Clean shutdown

Stop `aisgd`, `aisg-ctrl`, and `aisg-emu` with `Ctrl-C`, then remove temporary
sockets:

```bash
rm -f /tmp/aisg-ctrl.sock /tmp/aisg-ret0
```

Power down the external supply only after software is stopped.
