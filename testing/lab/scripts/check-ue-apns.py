#!/usr/bin/env python3
# Copyright (c) 2026-present, Ukama Inc. MPL-2.0.
"""Validate APN admission on both services without changing a live UE binding."""

import ipaddress
import json
from pathlib import Path
import re
import shlex
import subprocess
import sys


def validate(request, state):
    imsi = state['IMSI']

    def binding():
        code, ue = request('epcemu', 'ue/' + imsi)
        pcode, pdn = request('pcrf', 'pdn/imsi/' + imsi)
        keys = ('imsi', 'ip', 'epc_instance_id', 'pdn_id')
        if code != 200 or pcode != 200:
            raise ValueError('expected an existing EPCEMU/PCRF binding')
        identity = tuple(ue[k] for k in keys)
        if (ue.get('state') != 'attached' or pdn.get('state') != 'live' or
                identity != tuple(pdn[k] for k in keys) or
                identity[:2] != (imsi, state['UE_IP']) or
                ue.get('apn') != pdn.get('apn_name')):
            raise ValueError('EPCEMU/PCRF binding differs from the selected UE')
        return ue, (*identity, ue['apn'])

    ue, original = binding()
    epc = {k: ue[k] for k in ('imsi', 'iccid', 'ip')}
    pcrf = {k: ue[k] for k in ('epc_instance_id', 'pdn_id')}
    pcrf.update(imsi=[int(c) for c in imsi],
                pdn_address=int(ipaddress.IPv4Address(ue['ip'])))
    longest = 'a' * 63 + '.' + 'b' * 35
    invalid = ['', ' internet', 'internet ', 'a..b', '-a', 'a-', '.a',
               'a.', 'a_b', 'a/b', 'a\tb', 'a\nb', 'a\0b', 'é',
               'a' * 64, longest + 'b', None, 7]
    for service, endpoint, base, key in (
            ('epcemu', 'ue/attach', epc, 'apn'),
            ('pcrf', 'session', pcrf, 'apn_name')):
        for value in invalid:
            code, reply = request(service, endpoint, {**base, key: value})
            if code != 400:
                raise ValueError(f'{service} accepted malformed APN {value!r}: {code} {reply}')
            # EPCEMU's JSON decoder rejects embedded NUL before APN parsing.
            # A wrong JSON type may be rejected by PCRF's HTTP binder first.
            errors = {'invalid_apn'}
            if service == 'epcemu' and isinstance(value, str) and '\0' in value:
                errors.add('invalid JSON body')
            if isinstance(value, str) and reply.get('error') not in errors:
                raise ValueError(f'{service} did not report APN rejection: {reply}')
        for value in ('unsupported.ukama', longest, 'a' * 63):
            if value == original[-1]:
                continue
            code, reply = request(service, endpoint, {**base, key: value})
            if code != 403 or reply.get('error') != 'unsupported_apn':
                raise ValueError(f'{service} unsupported APN {value!r}: {code} {reply}')
        code, reply = request(service, endpoint, {**base, key: original[-1].upper()})
        expected = 200 if service == 'epcemu' else 202
        if (code != expected or reply.get(key) != original[-1] or
                reply.get('pdn_id') != original[3]):
            raise ValueError(f'{service} normalized retry changed the binding: {code} {reply}')
        print(f'PASS {service}: malformed, unsupported, boundary and case-normalized APNs')

    code, reply = request('epcemu', 'ue/attach', epc)
    if code != 200 or reply.get('apn') != original[-1] or reply.get('pdn_id') != original[3]:
        raise ValueError(f'EPCEMU omitted APN did not select configured default: {code} {reply}')
    code, reply = request('pcrf', 'session', pcrf)
    if code != 400 or reply.get('error') != 'invalid_apn':
        raise ValueError(f'PCRF accepted missing selected APN: {code} {reply}')
    if binding()[1] != original:
        raise ValueError('APN checks changed the existing PDN')
    print(f'PASS default selection and unchanged PDN {original[3]}')


def run(run_dir, ue_id):
    key = re.sub(r'[^A-Za-z0-9_.-]', '-', ue_id)
    state = {}
    for line in (Path(run_dir) / 'runtime-ues' / (key + '.env')).read_text().splitlines():
        name, sep, value = line.partition('=')
        if sep and name in ('TNODE_CONTAINER', 'IMSI', 'UE_IP'):
            state[name] = shlex.split(value)[0]

    def command(*args):
        return subprocess.check_output(['podman', 'exec', state['TNODE_CONTAINER'],
                                        *args], text=True, timeout=10)

    ports = {}
    for line in command('cat', '/etc/services').splitlines():
        words = line.partition('#')[0].split()
        if len(words) >= 2 and words[1].endswith('/tcp'):
            for name in [words[0], *words[2:]]:
                ports[name] = int(words[1].split('/')[0])

    def request(service, endpoint, data=None):
        args = ['curl', '-sS', '--max-time', '5', '-w', '\\n%{http_code}']
        if data is not None:
            args += ['-H', 'Content-Type: application/json', '-d', json.dumps(data)]
        body, code = command(*args, f'http://127.0.0.1:{ports[service]}/v1/{endpoint}').rsplit('\n', 1)
        return int(code), json.loads(body)

    validate(request, state)


if __name__ == '__main__':
    try:
        run(*sys.argv[1:])
    except (OSError, ValueError, KeyError, TypeError, subprocess.SubprocessError) as exc:
        sys.exit(str(exc))
