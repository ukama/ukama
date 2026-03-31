#!/usr/bin/env python3
"""
Standalone VE.Direct emulator launcher for supervisor/container use.

This wrapper reuses the existing implementation in vedirect_sim.py and
adds stable runtime defaults suitable for controllerd running inside the
same container/image.

Examples:

  /usr/bin/python3 /sbin/vedirect_emulator.py

  /usr/bin/python3 /sbin/vedirect_emulator.py \
      --mode realtime --city "Goma, DRC"

  /usr/bin/python3 /sbin/vedirect_emulator.py \
      --mode scenario --scenario sunny
"""

from __future__ import annotations

import os
import sys


def main() -> None:
    here = os.path.abspath(os.path.dirname(__file__))

    if here not in sys.path:
        sys.path.insert(0, here)

    import vedirect_sim

    argv = sys.argv[1:]

    has_serial_link = any(
        arg == "--serial-link" or arg.startswith("--serial-link=")
        for arg in argv
    )
    has_ready_file = any(
        arg == "--ready-file" or arg.startswith("--ready-file=")
        for arg in argv
    )

    if not has_serial_link:
        argv.extend(["--serial-link", "/tmp/victron-tty"])

    if not has_ready_file:
        argv.extend(["--ready-file", "/tmp/victron-emu.ready"])

    sys.argv = [sys.argv[0]] + argv
    vedirect_sim.main()


if __name__ == "__main__":
    main()
