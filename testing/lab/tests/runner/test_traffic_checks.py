"""Traffic check regression coverage through the production C checker.

Requires libcurl and Jansson development files; accepts CPPFLAGS/LDFLAGS/LDLIBS.
SPDX-License-Identifier: MPL-2.0
"""
import importlib.util
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]


class TrafficChecks(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.build = tempfile.TemporaryDirectory(prefix='ulab-traffic-checks-')
        cls.addClassCleanup(cls.build.cleanup)
        build = Path(cls.build.name)
        (build / 'version.h').write_text('#define VERSION "traffic-checks-test"\n')
        cls.binary = build / 'probe'
        subprocess.run([
            *shlex.split(os.environ.get('CC', 'cc')),
            '-std=gnu11', '-D_POSIX_C_SOURCE=200809L', '-O0',
            '-Wall', '-Wextra', '-Werror', '-Wdeclaration-after-statement',
            '-ffunction-sections', '-fdata-sections',
            *shlex.split(os.environ.get('CPPFLAGS', '')),
            '-I' + str(ROOT / 'inc'), '-I' + str(build),
            str(ROOT / 'tests/runner/traffic_check_probe.c'),
            *(str(ROOT / 'src' / name) for name in (
                'check_runtime.c', 'bff.c', 'bff_queries.c', 'selector.c',
                'world.c', 'util.c', 'log.c')),
            '-Wl,--gc-sections', *shlex.split(os.environ.get('LDFLAGS', '')),
            *shlex.split(os.environ.get('LDLIBS', '-lcurl -ljansson -lm')),
            '-o', str(cls.binary),
        ], check=True)

    def check_traffic(self, check, traffic, pdns, passes):
        result = subprocess.run([str(self.binary), check, traffic, *pdns],
                                capture_output=True, text=True, timeout=10)
        self.assertEqual(result.returncode, 0 if passes else 1,
                         result.stdout + result.stderr)

    def test_service_off_accepts_confirmed_absence(self):
        self.check_traffic('unavailable', 'failed', ['absent'], True)

    def test_policy_block_still_requires_live_blocked_pdn(self):
        for state, passes in [('blocked', True), ('absent', False),
                              ('allowed', False), ('unknown', False)]:
            with self.subTest(state=state):
                self.check_traffic('blocked', 'failed', [state], passes)

    def test_unavailable_checks_every_ue(self):
        for states, passes in [(['absent', 'blocked'], True),
                               (['blocked', 'absent'], True),
                               (['absent', 'allowed'], False),
                               (['blocked', 'unknown'], False),
                               ([], False)]:
            with self.subTest(states=states):
                self.check_traffic('unavailable', 'failed', states, passes)

    def test_failed_traffic_alone_does_not_pass(self):
        for state in ['allowed', 'unknown']:
            with self.subTest(state=state):
                self.check_traffic('unavailable', 'failed', [state], False)

    def test_successful_traffic_cannot_pass_negative_checks(self):
        self.check_traffic('unavailable', 'ok', ['absent'], False)
        self.check_traffic('blocked', 'ok', ['blocked'], False)
        self.check_traffic('allowed', 'ok', ['allowed'], True)
        self.check_traffic('allowed', 'failed', ['allowed'], False)


class AbsenceObservation(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        spec = importlib.util.spec_from_file_location('ue_state',
                                                     ROOT / 'scripts/ue-state.py')
        cls.observer = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(cls.observer)

    def test_absence_requires_both_endpoints_and_no_permission_flows(self):
        state = {'IMSI': '001010000000001', 'UE_APN': 'internet'}
        agent = {'state': 'detaching'}
        cases = [((404, 404), set(), True),
                 ((404, 503), set(), False),
                 ((500, 404), set(), False),
                 ((404, 200), set(), False),
                 ((404, 404), {'src', 'src_normal'}, False)]
        for codes, flows, passes in cases:
            with self.subTest(codes=codes, flows=flows), \
                    patch.object(self.observer, 'agent_status', return_value=agent), \
                    patch.object(self.observer, 'api',
                                 side_effect=[(code, None) for code in codes]), \
                    patch.object(self.observer, 'ovs_permissions', return_value=flows):
                if passes:
                    self.assertEqual(self.observer.observe(state, 'absent'), agent)
                else:
                    with self.assertRaises(self.observer.Pending):
                        self.observer.observe(state, 'absent')


if __name__ == '__main__':
    unittest.main()
