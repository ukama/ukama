#!/usr/bin/env python3
"""
Standalone VE.Direct emulator launcher for supervisor/container use.

This wrapper reuses the existing implementation in test/vedirect_sim.py
and only adds stable runtime defaults suitable for controllerd running
inside the same container.

Examples:

  python3 utils/vedirect_emulator.py

  python3 utils/vedirect_emulator.py --mode realtime --city "Goma, DRC"

  python3 utils/vedirect_emulator.py --mode scenario --scenario sunny
"""

from __future__ import annotations

import os
import sys


def main() -> None:
    here = os.path.abspath(os.path.dirname(__file__))
    root = os.path.abspath(os.path.join(here, ".."))
    test_dir = os.path.join(root, "test")

    if test_dir not in sys.path:
        sys.path.insert(0, test_dir)

    import vedirect_sim  # noqa: WPS433

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
