#!/usr/bin/env python3
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

"""Observe one fake UE and its PDN; never install or repair OVS flows."""

import argparse
import json
from pathlib import Path
import re
import shlex
import subprocess
import sys
import time
import uuid


class Pending(Exception):
    pass


def command(*args):
    try:
        result = subprocess.run(args, text=True, capture_output=True,
                                timeout=12, check=True)
    except (subprocess.SubprocessError, OSError) as error:
        raise Pending(f"command failed: {args[:3]}: {error}") from error
    return result.stdout


def load_state(key, run_dir):
    safe_key = re.sub(r"[^A-Za-z0-9_.-]", "-", key)
    path = Path(run_dir) / "runtime-ues" / (safe_key + ".env")
    state = {}
    for line in path.read_text().splitlines():
        if not line or line.startswith("#"):
            continue
        name, value = line.split("=", 1)
        words = shlex.split(value)
        state[name] = words[0] if len(words) == 1 else value
    for name in ("UE_CONTAINER", "TNODE_CONTAINER", "IMSI", "ICCID",
                 "UE_IP", "UE_ID"):
        if not state.get(name):
            raise ValueError(f"{name} missing in {path}")
    state.setdefault("UE_APN", "internet")
    return state


def api(state, port, path):
    output = command("podman", "exec", state["TNODE_CONTAINER"], "curl",
                     "-sS", "--max-time", "5", "-w", "\n%{http_code}",
                     f"http://127.0.0.1:{port}{path}")
    try:
        body, code = output.rsplit("\n", 1)
        code = int(code)
        return code, json.loads(body) if code == 200 else None
    except (ValueError, TypeError) as error:
        raise Pending("invalid HTTP/JSON observation") from error


def agent_status(state):
    output = command("podman", "exec", state["UE_CONTAINER"], "cat",
                     "/run/ukama-ue/status.json")
    try:
        agent = json.loads(output)
    except ValueError as error:
        raise Pending("invalid UE status") from error
    if (not isinstance(agent, dict) or agent.get("source") != "ue-agent"
            or any(agent.get(key) != state[name] for key, name in
                   (("imsi", "IMSI"), ("iccid", "ICCID"),
                    ("ip", "UE_IP"), ("requested_apn", "UE_APN")))):
        raise Pending("UE status does not match the requested UE/APN")
    return agent


def identity(body):
    try:
        values = tuple(body[key] for key in ("epc_instance_id", "pdn_id"))
        if any(str(uuid.UUID(value)) != value for value in values):
            raise ValueError("noncanonical identity")
        return values
    except (KeyError, TypeError, ValueError, AttributeError) as error:
        raise Pending("missing or invalid PDN identity") from error


def ovs_permissions(state):
    output = command("podman", "exec", state["TNODE_CONTAINER"],
                     "ovs-ofctl", "-O", "OpenFlow15", "dump-flows", "br0")
    found = set()
    ip = re.escape(state["UE_IP"])
    for line in output.splitlines():
        if not re.search(r"\bpriority=100(?:,|\s)", line):
            continue
        for direction in ("src", "dst"):
            if re.search(rf"\bnw_{direction}={ip}(?:/32)?(?:,|\s)",
                         line):
                found.add(direction)
                if re.search(r"\bactions=.*\bNORMAL\b", line):
                    found.add(direction + "_normal")
    return found


def observe(state, expected, apn="-", snapshot=None, same=False):
    agent = agent_status(state)
    epc_code, epc = api(state, 18028, "/v1/ue/" + state["IMSI"])
    pcrf_code, pcrf = api(state, 18030,
                          "/v1/pdn/imsi/" + state["IMSI"])
    requested = state["UE_APN"]
    if apn != "-" and requested.lower() != apn.lower():
        raise Pending("requested APN differs from the scenario expectation")

    if expected in ("detached", "rejected"):
        if agent.get("state") != expected:
            raise Pending(f"UE state={agent.get('state')}, need {expected}")
        if expected == "rejected" and (
                agent.get("http_status") != 403 or
                agent.get("reason") != "unsupported_apn"):
            raise Pending("no explicit unsupported_apn rejection")
        if epc_code != 404 or pcrf_code != 404:
            raise Pending("EPC/PCRF still has a PDN or cannot confirm absence")
        if ovs_permissions(state):
            raise Pending("OVS still has a UE permission flow")
        return agent

    if (epc_code != 200 or pcrf_code != 200 or
            not isinstance(epc, dict) or not isinstance(pcrf, dict)):
        raise Pending(f"PDN unavailable: EPC={epc_code} PCRF={pcrf_code}")
    if (agent.get("state") != "attached" or
            epc.get("state") != "attached" or pcrf.get("state") != "live"):
        raise Pending("PDN is not attached/live in all three components")
    for body, apn_key in ((agent, "apn"), (epc, "apn"),
                          (pcrf, "apn_name")):
        if (body.get("imsi") != state["IMSI"] or
                body.get("ip") != state["UE_IP"] or
                body.get(apn_key) != requested.lower()):
            raise Pending("PDN IMSI/IP/APN differs between components")
    if epc.get("iccid") != state["ICCID"]:
        raise Pending("EPC ICCID does not match this UE")
    binding = identity(agent)
    if identity(epc) != binding or identity(pcrf) != binding:
        raise Pending("PDN identity differs between UE, EPC and PCRF")
    if same:
        if not snapshot or not snapshot.exists():
            raise Pending("no initial PDN snapshot; wait for attach first")
        if identity(json.loads(snapshot.read_text())) != binding:
            raise Pending("PDN identity changed")
    if expected in ("allowed", "blocked"):
        if any(body.get("data_state") != expected
               for body in (agent, epc, pcrf)):
            raise Pending(f"data permission is not confirmed {expected}")
        flows = ovs_permissions(state)
        if expected == "blocked" and flows:
            raise Pending("blocked PDN retains an OVS permission flow")
        if expected == "allowed" and not {
                "src_normal", "dst_normal"}.issubset(flows):
            raise Pending("allowed PDN lacks both OVS permission flows")
    return agent


def wait(state, expected, seconds, apn, snapshot, same=False,
         remember=False):
    deadline = time.monotonic() + seconds
    while True:
        try:
            result = observe(state, expected, apn, snapshot, same)
            if remember and not snapshot.exists():
                temporary = snapshot.with_suffix(".tmp")
                temporary.write_text(json.dumps(result) + "\n")
                temporary.replace(snapshot)
            print(f"ue-pdn imsi={state['IMSI']} expected={expected} "
                  f"apn={state['UE_APN']} pdn={result.get('pdn_id', '')}")
            return
        except Pending as error:
            if time.monotonic() >= deadline:
                raise Pending(f"IMSI {state['IMSI']}: {error}") from error
            time.sleep(min(1, max(0, deadline - time.monotonic())))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("ue")
    parser.add_argument("run_dir")
    parser.add_argument("expected", choices=("attached", "allowed", "blocked",
                                           "rejected", "detached", "detach"))
    parser.add_argument("seconds", type=int, nargs="?", default=0)
    parser.add_argument("apn", nargs="?", default="-")
    parser.add_argument("same", type=int, choices=(0, 1), nargs="?", default=0)
    parser.add_argument("--remember", action="store_true")
    args = parser.parse_args()
    try:
        state = load_state(args.ue, args.run_dir)
        snapshot = (Path(args.run_dir) / "runtime-ues" /
                    (state["UE_ID"] + ".pdn.json"))
        if args.expected == "detach":
            # The agent stops forwarding and owns identity-pinned cleanup.
            # Do not kill it first or send a fresh IMSI-only DELETE here.
            status = agent_status(state)
            if status.get("state") != "detached":
                command("podman", "kill", "--signal", "USR1",
                        state["UE_CONTAINER"])
            args.expected = "detached"
        wait(state, args.expected, args.seconds, args.apn, snapshot,
             bool(args.same), args.remember)
    except (Pending, OSError, ValueError) as error:
        print(f"ue-state: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
