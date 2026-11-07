#!/usr/bin/env python3
# Copyright (c) 2026-present, Ukama Inc. MPL-2.0.
"""Try invalid addresses with the selected lab UE's IMSI and verify its binding."""

import ipaddress
import json
import re
import shlex
import subprocess
import sys
from pathlib import Path


def init_network_command(command):
    # starterd may execute ld-linux, with init-network as the loader's argument.
    # Encode each cmdline as one record, preserving NUL-separated argv entries.
    processes = command('sh', '-c', r"""for p in /proc/[0-9]*; do
    [ -r "$p/cmdline" ] || continue
    printf '%s ' "${p##*/}"
    od -An -v -tx1 "$p/cmdline" 2>/dev/null | tr -d '\n'
    printf '\n'
done""")
    matches = []
    app_names = {'init-network', 'init-network.d'}
    for record in processes.splitlines():
        pid, _, encoded = record.partition(' ')
        if not pid.isdecimal() or not encoded.strip():
            continue
        argv = bytes.fromhex(encoded).decode(errors='surrogateescape').split('\0')
        name = Path(argv[0]).name
        if name in app_names:
            matches.append([f'/proc/{pid}/exe'])
            continue
        if not (name == 'ld.so' or (name.startswith('ld-') and '.so' in name)):
            continue
        index = 1
        if len(argv) > 2 and argv[1] == '--library-path':
            index = 3
        if len(argv) <= index or Path(argv[index]).name not in app_names:
            continue
        app = argv[index]
        if not app.startswith('/'):
            app = f'/proc/{pid}/cwd/{app}'
        matches.append([f'/proc/{pid}/exe', *argv[1:index], app])
    if len(matches) != 1:
        raise ValueError('expected one running init-network app on the tower '
                         f'(direct or loader launch); found {len(matches)}')
    return matches[0]


def run(run_dir, ue_id):
    key = re.sub(r'[^A-Za-z0-9_.-]', '-', ue_id)
    path = Path(run_dir) / 'runtime-ues' / (key + '.env')
    state = {}
    for line in path.read_text().splitlines():
        name, sep, value = line.partition('=')
        if sep and name in ('TNODE_CONTAINER', 'IMSI', 'UE_IP'):
            state[name] = shlex.split(value)[0]

    def command(*args):
        return subprocess.check_output(['podman', 'exec', state['TNODE_CONTAINER'],
                                        *args], text=True, timeout=10)

    # Run the same parser used by the deployed daemon before it sets up OVS.
    validator = init_network_command(command)
    inputs = Path(__file__).resolve().parents[1] / 'scenarios/p0/ue/addressing/input'
    cases = json.loads((inputs / 'ue-subnets.json').read_text())
    for cidr, expected in cases.items():
        if expected not in ('valid', 'invalid'):
            raise ValueError(f'bad subnet expectation for {cidr}: {expected}')
        probe = subprocess.run(['podman', 'exec', state['TNODE_CONTAINER'],
                                *validator, '--check-ue-subnet', cidr],
                               text=True, capture_output=True, timeout=10)
        if probe.returncode != int(expected == 'invalid') or probe.stdout.strip() != expected:
            raise ValueError(f'init-network subnet {cidr!r}: expected {expected}, '
                             f'got exit {probe.returncode}, {probe.stdout!r} {probe.stderr!r}')
        print(f'PASS init-network subnet {cidr!r}: {expected}')

    ports = {}
    for line in command('cat', '/etc/services').splitlines():
        words = line.partition('#')[0].split()
        if len(words) >= 2 and words[1].endswith('/tcp'):
            for name in [words[0], *words[2:]]:
                ports[name] = int(words[1].split('/')[0])

    def request(service, endpoint, data=None):
        args = ['curl', '-sS', '--max-time', '5', '-w', '\\n%{http_code}']
        if data is not None:
            args.extend(['-H', 'Content-Type: application/json', '-d', json.dumps(data)])
        body, code = command(*args, f'http://127.0.0.1:{ports[service]}/v1/{endpoint}').rsplit('\n', 1)
        return int(code), json.loads(body)

    def get(service, endpoint):
        code, obj = request(service, endpoint)
        if code != 200:
            raise ValueError(f'{service}/{endpoint}: HTTP {code} {obj}')
        return obj

    imsi = state['IMSI']
    def binding():
        ue = get('epcemu', 'ue/' + imsi)
        pdn = get('pcrf', 'pdn/imsi/' + imsi)
        fields = ('imsi', 'ip', 'epc_instance_id', 'pdn_id')
        identity = tuple(ue[k] for k in fields)
        if (ue['state'] != 'attached' or pdn['state'] != 'live' or
                identity != tuple(pdn[k] for k in fields) or
                identity[:2] != (imsi, state['UE_IP'])):
            raise ValueError('expected the selected UE to have a live PDN')
        return identity

    original = binding()
    config = get('init-network', 'status')
    pool = ipaddress.IPv4Network(config['ue']['cidr'])
    invalid = {str(pool.network_address), str(pool.broadcast_address), '192.0.2.999'}
    if pool.prefixlen:
        invalid.add(str(ipaddress.IPv4Address(int(pool.network_address) ^ (1 << 31))))
    for cidr in (config['tun']['primaryCidr'], config['bridge']['cidr']):
        invalid.add(str(ipaddress.IPv4Interface(cidr).ip))
    if config.get('gateway', {}).get('enabled'):
        invalid.add(str(ipaddress.IPv4Interface(config['gateway']['address']).ip))
    for ip in sorted(invalid):
        code, reply = request('epcemu', 'ue/attach', {'imsi': imsi, 'ip': ip})
        if code != 400 or reply.get('error') != 'invalid_ue_address':
            raise ValueError(f'attach {ip}: expected HTTP 400 invalid_ue_address, got {code} {reply}')
    if binding() != original:
        raise ValueError('invalid address requests changed the original binding')
    print(f'PASS {imsi}: rejected {len(invalid)} invalid addresses, retained {original[-1]}')


if __name__ == '__main__':
    try:
        run(*sys.argv[1:])
    except (OSError, ValueError, KeyError, subprocess.SubprocessError) as exc:
        sys.exit(str(exc))
