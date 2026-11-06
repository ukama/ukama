#!/usr/bin/env python3
# Copyright (c) 2026-present, Ukama Inc. MPL-2.0.
"""Validate observed group state. A dispatched command is not completion."""

import json
import sys


def confirmed(expected, device, network, pcrf, epc):
    group = device.get("service_group", {})
    if (device.get("service") != expected or group.get("phase") != "idle" or
            group.get("current") != expected or group.get("desired") != expected):
        return False
    if (network.get("state") != expected or
            network.get("confirmed") is not True):
        return False
    epc_service = epc.get("service", {})
    for service in (pcrf, epc_service):
        if (service.get("state") != expected or
                service.get("cleanup_pending") is not False or
                service.get("confirmed") is not True):
            return False
        if expected == "off" and service.get("pdn_count") != 0:
            return False
        if service.get("admission") != (
                "enabled" if expected == "on" else "disabled"):
            return False
    if expected == "on":
        return (network.get("ready") is True and pcrf.get("ready") is True and
                epc.get("ready") is True and
                epc.get("pcrf", {}).get("ready") is True)
    return expected == "off"


def main(argv):
    if len(argv) != 6 or argv[1] not in ("on", "off"):
        return 2
    try:
        bodies = [json.loads(value) for value in argv[2:]]
        if not all(isinstance(body, dict) for body in bodies):
            return 1
        return 0 if confirmed(argv[1], *bodies) else 1
    except (TypeError, ValueError, AttributeError):
        return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv))
