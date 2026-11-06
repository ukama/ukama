#!/usr/bin/env python3
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

"""Set a virtual UE's PCRF test policy or observe its actual OVS meters."""

import ipaddress
import json
from pathlib import Path
import re
import shlex
import subprocess
import sys
import time
import uuid


def execute(tower, *args):
    proc = subprocess.run(['podman', 'exec', tower, *args], text=True,
                          capture_output=True, timeout=10)
    if proc.returncode:
        raise RuntimeError(proc.stderr.strip() or f'{args[0]} failed')
    return proc.stdout


def pcrf_url(tower):
    for line in execute(tower, 'cat', '/etc/services').splitlines():
        words = line.partition('#')[0].split()
        if len(words) >= 2 and 'pcrf' in [words[0], *words[2:]]:
            port, _, proto = words[1].partition('/')
            if proto == 'tcp' and port.isdecimal():
                return f'http://127.0.0.1:{port}/v1'
    raise ValueError('pcrf/tcp missing from node service table')


def retry(probe):
    deadline = time.monotonic() + 60
    while True:
        try:
            return probe()
        except (RuntimeError, subprocess.TimeoutExpired) as error:
            if time.monotonic() >= deadline:
                raise RuntimeError(str(error)) from error
            time.sleep(1)


def push(tower, imsi, ul, dl):
    base = pcrf_url(tower)
    curl = ('curl', '-fsS', '--max-time', '5')
    # Rate changes get a new policy ID. Send a complete eligible policy for
    # this run's disposable IMSI, just as a subscriber push would.
    now = int(time.time())
    policy = dict(uuid=str(uuid.uuid4()), ulbr=ul, dlbr=dl, burst=1500,
                  total_data=100 * 1024 * 1024, consumed_data=0,
                  start_time=now, end_time=now + 3600)
    execute(tower, *curl, '-X', 'PATCH', '-H', 'Content-Type: application/json',
            '-d', json.dumps({'imsi': imsi, 'policy': policy}), f'{base}/subscriber/imsi/{imsi}')


def installed_rates(flows, meters, ip):
    """Follow each exact UE flow's meter action to its drop band in kbps."""
    rates = []
    for field in ('nw_src', 'nw_dst'):
        matches = []
        for line in flows.splitlines():
            match, sep, actions = line.partition('actions=')
            if (sep and re.search(r'\btable=0\b', match) and
                    re.search(r'\bpriority=100\b', match) and
                    re.search(rf'\b{field}={re.escape(ip)}(?:/32)?(?=[,\s]|$)', match)):
                matches.append(actions.strip())
        if len(matches) != 1:
            raise RuntimeError(f'{ip} {field}: expected one UE flow, found {len(matches)}')
        action = re.fullmatch(r'meter:(\d+),NORMAL', matches[0], re.I)
        if not action:
            raise RuntimeError(f'{ip} {field}: unmetered or unexpected actions {matches[0]}')
        meter = action[1]
        blocks = re.split(r'(?=^\s*meter=\d+\b)', meters, flags=re.M)
        block = next((b for b in blocks if re.match(rf'^\s*meter={meter}\b', b)), '')
        bands = re.findall(r'\btype=drop[,\s]+rate=(\d+)\b', block)
        if not re.search(r'\bkbps\b', block) or len(bands) != 1:
            raise RuntimeError(f'{ip} {field}: expected one kbps drop band for meter {meter}')
        rates.append(int(bands[0]))
    return tuple(rates)


def main(run_dir, site, imsi, ip, operation, uplink, downlink):
    ip = str(ipaddress.IPv4Address(ip))
    if not re.fullmatch('[0-9]{5,15}', imsi):
        raise ValueError('invalid IMSI')
    safe = re.sub('[^A-Za-z0-9_.-]', '-', site)
    path = Path(run_dir) / 'runtime-sites' / (safe + '.env')
    tower = next((shlex.split(line.split('=', 1)[1])[0]
                  for line in path.read_text().splitlines()
                  if line.startswith('TNODE_CONTAINER=')), None)
    if not tower:
        raise ValueError(f'TNODE_CONTAINER missing in {path}')
    wanted = int(uplink), int(downlink)
    if operation == 'set':
        push(tower, imsi, *wanted)
    elif operation == 'wait':
        def observe():
            ovs = ('ovs-ofctl', '-O', 'OpenFlow15')
            got = installed_rates(execute(tower, *ovs, 'dump-flows', 'br0'),
                                  execute(tower, *ovs, 'dump-meters', 'br0'), ip)
            if got != wanted:
                raise RuntimeError(f'{ip}: expected UL/DL {wanted} kbps, got {got}')
        retry(observe)
    else:
        raise ValueError(f'unknown operation {operation}')
    print(f'{operation}: imsi={imsi} ip={ip} uplink={uplink} downlink={downlink} kbps')


if __name__ == '__main__':
    try:
        main(*sys.argv[1:])
    except (OSError, RuntimeError, ValueError, KeyError, subprocess.TimeoutExpired) as error:
        sys.exit(str(error))
