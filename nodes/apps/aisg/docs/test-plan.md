# Test Plan

1. Start `aisg-emu`.
2. Start `aisgd` against the emulator socket.
3. Verify `/v1/ping`, `/v1/version`, and `/v1/status`.
4. Run `/v1/reconcile`.
5. Exercise scan, info, alarms, self-test, config, and calibration.
6. Exercise get tilt, set tilt, and get device data.
7. Repeat with `aisg-ctrl` raw-rs485 backend and real RET hardware.
