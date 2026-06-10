# aisg-ctrl

Controller daemon for AISG/RET communication.

Responsibilities:

- UNIX socket controller contract server
- backend abstraction
- raw RS-485 backend for lab testing
- STM UART backend placeholder for production daughter-card integration
- HDLC framing
- AISG v2 framing/discovery
- 3GPP RETAP procedure encoding/parsing

Important files:

- `src/hdlc.c`: HDLC frame encode/decode and FCS.
- `src/aisg_v2.c`: AISG v2 bus/message handling.
- `src/retap.c`: RETAP envelope encode/decode.
- `src/retap_ops.c`: procedure-specific payloads.
- `src/backend_raw_rs485.c`: lab backend orchestration.
