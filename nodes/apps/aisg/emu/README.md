# aisg-emu

Software-only AISG controller emulator.

The emulator implements the controller contract exposed by `aisg-ctrl`.
It does not emulate raw HDLC or RS-485. This keeps virtual-node tests
deterministic and keeps `aisgd` unchanged between lab, virtual, and
production flows.
