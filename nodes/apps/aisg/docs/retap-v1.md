# RETAP V1

V1 implements the single-antenna RETAP procedures needed for lab bring-up:

- Get Information
- Get Alarm Status
- Clear Active Alarms
- Alarm Subscribe
- Self Test
- Send Configuration Data
- Calibrate
- Get Tilt
- Set Tilt
- Get Device Data
- Reset Software

The raw RS-485 backend includes module boundaries for HDLC, AISG v2, and
RETAP. Bench validation against the real RET is still required to tune exact
frame, address, control, and vendor configuration behavior.
