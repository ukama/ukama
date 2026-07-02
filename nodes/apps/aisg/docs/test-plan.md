# AISG Test Plan

Ukama AISG supports one AISG v2 / TS 25.xxx single-antenna RET device.
Multi-antenna RET, TMA, multiple devices on one bus, daisy-chain discovery,
and firmware download mode are intentionally out of scope.

The validation order is:

```text
1. Protocol golden unit tests
2. Strict RET emulator tests
3. Read-only hardware ladder
4. Optional movement hardware ladder
```

## 1. Protocol golden tests

These compile only the protocol helpers with local stubs, so they do not need
the full Ukama platform/vendor build tree.

```bash
cd aisg
make protocol-test
```

Covered checks:

```text
HDLC:
  encode/decode roundtrip
  FCS failure rejected
  0x7E and 0x7D escaping
  I-frame ns/nr/poll helpers

XID:
  valid scan parsed
  address assignment parsed
  PI=3/bit-mask is rejected for address assignment matching
  single-RET device type matching

RETAP:
  GetInformation = 05 00 00
  GetTilt = 34 00 00
  SetTilt +3.2 = 33 02 00 20 00
  OK/FAIL response parsing
  FAIL=0x0B and NotCalibrated mapping
  config segment max = 70 bytes
```

## 2. Strict RET emulator ladder

Build `ctrl` and `emu` first:

```bash
cd aisg
make -C ctrl
make -C emu
```

Run the standards-shaped software ladder:

```bash
tests/scripts/run-ret-emu-ladder.sh
```

The script starts:

```text
aisg-emu --mode ret --pty <tmp>/aisg-ret0
aisg-ctrl --backend raw-rs485 --device <tmp>/aisg-ret0
```

Then it verifies:

```text
get_status
scan/connect
get_info
get_alarm_status
get_tilt before calibration fails with NotCalibrated
send_configuration_data
calibrate
get_tilt
set_tilt 0.5
get_tilt verify
```

Keep logs for debugging:

```bash
AISG_KEEP_LOGS=1 tests/scripts/run-ret-emu-ladder.sh
```

## 3. Strict RET emulator negative checks

```bash
tests/scripts/run-ret-emu-negative.sh
```

Covered checks:

```text
old fake scan FF BF 81 F0 00 produces no response
RETAP before address assignment/SNRM produces no response
```

These checks make sure the emulator is strict enough to catch the old
non-standard controller behavior.

## 4. Read-only hardware ladder

Connect the real RET through USB-RS485, then run:

```bash
AISG_TTY=/dev/ttyUSB0 tests/scripts/run-real-hw-ladder.sh
```

Read-only sequence:

```text
get_status
scan/connect
get_info
get_alarm_status
```

If scan fails:

```text
check power on AISG Pin6/Pin7
check enough current for motor movement
check RS485 GND/reference
try RS485 A/B swapped once
check that received bytes are not only TX echo
```

## 5. Optional movement hardware ladder

Only after read-only hardware tests pass:

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

Do not run movement repeatedly until the antenna tilt range and vendor config
file are confirmed.

## 6. Interpreting failures

```text
If protocol golden tests fail:
  HDLC/XID/RETAP helpers are wrong.

If ctrl fails aisg-emu --mode ret:
  controller protocol flow is wrong.

If ctrl passes aisg-emu --mode ret but hardware fails:
  investigate power, wiring, RS485 A/B polarity, adapter direction control,
  timing, vendor config file, or vendor-specific quirks.
```
