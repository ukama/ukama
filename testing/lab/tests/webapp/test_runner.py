"""Full CLI lifecycle tests using a fault-controlled worker (no product credit).

Build ukama-lab first, then set ULAB_TEST_BINARY to that executable.
SPDX-License-Identifier: MPL-2.0
"""
import json
import os
from pathlib import Path
import signal
import subprocess
import tempfile
import time
import unittest

ROOT = Path(__file__).resolve().parents[2]
BINARY = Path(os.environ.get('ULAB_TEST_BINARY', ROOT / 'bin/ukama-lab')).resolve()
EXAMPLE = ROOT / 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml'


class RunnerLifecycle(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        if not BINARY.is_file():
            raise RuntimeError('Build ukama-lab and set ULAB_TEST_BINARY before running lifecycle tests')

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='ulab runner "test" ')
        self.addCleanup(self.tmp.cleanup)
        self.directory = Path(self.tmp.name)
        self.scenario = self.directory / 'case.yaml'
        self.scenario.write_text(EXAMPLE.read_text())
        self.out = self.directory / 'out'
        self.run = self.out / 'test-run'

    def command(self, worker=None, run_id='test-run'):
        return [str(BINARY), 'run', str(self.scenario), '--repo', '/test/ukama',
                '--out', str(self.out), '--run-id', run_id,
                '--webapp-worker', str(worker or ROOT / 'tests/webapp/fake_worker.py')]

    def run_case(self, mode='ok', **kwargs):
        r = subprocess.run(self.command(**kwargs), cwd=ROOT, env={**os.environ, 'ULAB_FAKE_MODE': mode},
                           text=True, capture_output=True, timeout=35)
        report = json.loads((self.run / 'report.json').read_text()) if self.run.exists() and (self.run / 'report.json').exists() else None
        return r, report

    def messages(self):
        return [json.loads(line) for line in (self.run / 'webapp-commands.jsonl').read_text().splitlines()]

    def test_ordered_execution_reports_visible_expectations_and_private_artifacts(self):
        r, report = self.run_case()
        self.assertEqual(r.returncode, 0, r.stdout + r.stderr)
        self.assertTrue(report['passed'])
        self.assertEqual(report['outcome'], 'PASS')
        self.assertEqual(report['checks'], {'total': 2, 'passed': 2, 'failed': 0})
        requests = [x['message'] for x in self.messages() if x['direction'] == 'request']
        self.assertEqual([x['action'] for x in requests], ['init', 'web_open', 'web_action_available', 'web_reload', 'web_action_available', 'close'])
        self.assertEqual([x['command_id'] for x in requests], list(range(1, 7)))
        checks = [x for x in report['results'] if x['kind'] == 'check']
        self.assertTrue(all(x['expected'] is True and x['actual'] is True and x['requirement'] == 'WEB-TEAM-002' for x in checks))
        self.assertEqual(report['artifacts']['run_dir'], str(self.run))
        self.assertEqual(self.run.stat().st_mode & 0o777, 0o700)
        for filename in ('report.json', 'report.txt', 'webapp-resources.json', 'webapp-commands.jsonl', 'webapp-worker.log'):
            self.assertEqual((self.run / filename).stat().st_mode & 0o777, 0o600)
        self.assertFalse((self.run / 'created.json').exists(), 'browser run must not enter BFF setup')
        self.assertEqual(json.loads((self.run / 'webapp-resources.json').read_text())['resources'], [])

    def test_protocol_and_worker_failures_never_replay_or_pass(self):
        for mode in ('wrong_id', 'malformed', 'crash', 'lie', 'fail'):
            with self.subTest(mode=mode), tempfile.TemporaryDirectory(dir=self.directory) as out:
                self.out = Path(out); self.run = self.out / 'test-run'
                r, report = self.run_case(mode)
                self.assertNotEqual(r.returncode, 0, r.stdout + r.stderr)
                self.assertFalse(report['passed'])
                requests = [x['message'] for x in self.messages() if x['direction'] == 'request']
                self.assertEqual(sum(x['action'] == 'web_open' for x in requests), 1)
                self.assertFalse(any(x['action'] == 'web_reload' for x in requests))
                self.assertEqual(len({x['command_id'] for x in requests}), len(requests))

    def test_skip_is_not_pass_and_launches_no_worker(self):
        for status in ('skip', 'wip'):
            with self.subTest(status=status), tempfile.TemporaryDirectory(dir=self.directory) as out:
                self.out = Path(out); self.run = self.out / 'test-run'
                self.scenario.write_text(EXAMPLE.read_text().replace('status: active', f'status: {status}'))
                r, report = self.run_case(worker='/missing/worker')
                self.assertEqual(r.returncode, 0, r.stderr)
                self.assertEqual(report['outcome'], 'SKIP')
                self.assertFalse(report['passed'])
                self.assertEqual(report['checks']['total'], 0)
                self.assertFalse((self.run / 'webapp-commands.jsonl').exists())

    def test_existing_run_and_invalid_id_are_rejected_without_overwrite(self):
        self.assertEqual(self.run_case()[0].returncode, 0)
        before = (self.run / 'report.json').read_bytes()
        self.assertNotEqual(self.run_case()[0].returncode, 0)
        self.assertEqual((self.run / 'report.json').read_bytes(), before)
        for run_id in ('../escape', 'bad id', 'a' * 81):
            self.assertNotEqual(self.run_case(run_id=run_id)[0].returncode, 0)
        self.assertFalse((self.directory / 'escape').exists())

    def test_unimplemented_provisioning_and_xfail_fail_before_launch(self):
        text = (ROOT / 'scenarios/webapp/p0/network/wb-001-sites-online-recovery.yaml').read_text()
        for scenario in (text.replace('status: wip', 'status: active'), EXAMPLE.read_text().replace('status: active', 'status: xfail')):
            self.scenario.write_text(scenario)
            r, report = self.run_case()
            self.assertNotEqual(r.returncode, 0)
            self.assertIsNone(report)
            self.assertFalse(self.run.exists())

    def test_missing_worker_fails_and_retains_journal(self):
        r, report = self.run_case(worker='/missing/worker')
        self.assertNotEqual(r.returncode, 0)
        self.assertEqual(report['outcome'], 'FAIL')
        self.assertEqual(json.loads((self.run / 'webapp-resources.json').read_text())['run_result'], 'failed')

    def test_deadline_is_bounded(self):
        self.scenario.write_text(EXAMPLE.read_text().replace('action_timeout_seconds: 30', 'action_timeout_seconds: 1'))
        start = time.monotonic()
        r, report = self.run_case('hang')
        self.assertNotEqual(r.returncode, 0)
        self.assertFalse(report['passed'])
        self.assertLess(time.monotonic() - start, 20)

    def test_sigterm_cancels_active_worker_and_keeps_one_failed_result(self):
        with subprocess.Popen(self.command(), cwd=ROOT, env={**os.environ, 'ULAB_FAKE_MODE': 'cancel'},
                              stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True) as child:
            deadline = time.monotonic() + 8
            while time.monotonic() < deadline:
                journal = self.run / 'webapp-commands.jsonl'
                if journal.exists() and 'web_open' in journal.read_text():
                    break
                time.sleep(.02)
            else:
                child.kill(); self.fail('worker command did not start')
            child.send_signal(signal.SIGTERM)
            stdout, stderr = child.communicate(timeout=15)
            self.assertNotEqual(child.returncode, 0, stdout + stderr)
        report = json.loads((self.run / 'report.json').read_text())
        self.assertEqual(report['outcome'], 'FAIL')
        self.assertEqual(json.loads((self.run / 'webapp-resources.json').read_text())['cleanup'], 'complete')
        self.assertFalse(any(x['message']['action'] == 'web_reload' for x in self.messages()))


if __name__ == '__main__':
    unittest.main()
