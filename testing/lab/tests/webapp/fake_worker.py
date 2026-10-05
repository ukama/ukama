#!/usr/bin/env python3
# SPDX-License-Identifier: MPL-2.0
# Controlled protocol faults for C transport tests; never product coverage.
import json
import os
import signal
import sys
import time

mode = os.environ.get('ULAB_FAKE_MODE', 'ok')
failed = False
signal.signal(signal.SIGTERM, lambda *_: sys.exit(1))
for line in sys.stdin:
    c = json.loads(line)
    action, inputs = c['action'], c['inputs']
    expected = inputs.get('expected', inputs.get('expected_count', inputs.get('available')))
    actual = expected
    if action == 'init':
        actual = {'authenticated': True}
    if action == 'web_tab':
        actual = {'executed': mode != 'bad_ack'}
    if action == 'web_field_equals' and inputs.get('match') == 'contains':
        actual = 'unrelated reason' if mode == 'lie' else expected + ' by fixture'
    if mode in ('hang', 'cancel') and action == 'web_open':
        while True:
            time.sleep(.1)
    if mode == 'crash' and action == 'web_open':
        sys.exit(9)
    if mode == 'malformed' and action == 'web_open':
        print('{broken', flush=True)
        continue
    if mode == 'lie' and action == 'web_action_available':
        actual = not expected
    failed = failed or (mode == 'fail' and action == 'web_action_available') or (action == 'close' and inputs.get('failed', False))
    r = dict(protocol=1, run_id=c['run_id'], command_id=c['command_id'], action=action,
             status='ok', run_status=('failed' if failed else 'passed') if action == 'close' else 'running',
             duration_ms=0, expected=expected, actual=actual, bindings=[], artifacts=[])
    if mode == 'wrong_id' and action == 'web_open':
        r['command_id'] += 10
    if failed and action != 'close':
        r.update(status='error', run_status='failed', error={'code': 'FIXTURE_FAILURE', 'message': 'controlled assertion failure'})
    print('fixture diagnostic on stderr', file=sys.stderr, flush=True)
    print(json.dumps(r), flush=True)
    if action == 'close':
        sys.exit(1 if failed else 0)
sys.exit(1)
