# aisg-emu

`aisg-emu` has two modes.

## Contract mode

```sh
aisg-emu --mode contract
```

This is the original emulator. It speaks the `aisg-ctrl` JSON controller contract over the Unix socket and is useful for fast `aisgd` API tests. It does **not** emulate raw HDLC, RS485, XID, SNRM/UA, or RETAP packets.

## RET mode

```sh
aisg-emu --mode ret --pty /tmp/aisg-ret0 --log-level DEBUG
```

RET mode behaves as one strict AISG v2 / 3GPP TS 25.461/25.462/25.463 single-antenna RET secondary device. It exposes a pseudo-terminal symlink that `aisg-ctrl` can open like a serial device:

```sh
aisg-ctrl --backend raw-rs485 --device /tmp/aisg-ret0
```

RET mode scope:

```text
Supported:
  - one single-antenna RET device
  - device type 0x01
  - XID scan
  - XID address assignment
  - 3GPP release negotiation
  - AISG v2 negotiation
  - SNRM/UA link establishment
  - sequenced I-frame RETAP request/response
  - GetInformation
  - GetErrorStatus
  - ClearActiveAlarms
  - AlarmSubscribe
  - SendConfigurationData, max 70 bytes per request
  - Calibrate
  - GetTilt
  - SetTilt

Rejected/unsupported:
  - multi-antenna RET device type 0x11
  - TMA device type 0x02
  - multiple devices on one bus
  - RETAP before address assignment
  - RETAP before SNRM/UA
  - malformed RETAP without procedure/length/data envelope
```

Useful RET options:

```sh
--pty PATH
--vendor XX
--serial SERIAL
--requires-config true|false
--initial-tilt DEG
--min-tilt DEG
--max-tilt DEG
```

Example:

```sh
aisg-emu --mode ret \
  --pty /tmp/aisg-ret0 \
  --vendor UK \
  --serial UKAMA00000000001 \
  --requires-config true \
  --initial-tilt 3.0 \
  --min-tilt 0.0 \
  --max-tilt 10.0 \
  --log-level DEBUG
```
