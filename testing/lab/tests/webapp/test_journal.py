"""Resource ownership, cleanup and runtime hooks using real C modules.
SPDX-License-Identifier: MPL-2.0
"""
import json
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]


class ResourceLifecycle(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.tmp = tempfile.TemporaryDirectory(prefix='ulab-journal-')
        cls.directory = Path(cls.tmp.name)
        (cls.directory / 'version.h').write_text('#define VERSION "journal-test"\n')
        cls.binary = cls.directory / 'probe'
        files = ['scenario', 'validate', 'scenario_webapp', 'world', 'report', 'util', 'log',
                 'webapp_client', 'webapp_journal', 'webapp_runner', 'workload_config']
        sources = [str(ROOT / f'src/{f}.c') for f in files]
        command = [os.environ.get('CC', 'cc'), '-std=gnu11', '-D_POSIX_C_SOURCE=200809L', '-Wall', '-Wextra', '-Werror',
                   '-ffunction-sections', '-fdata-sections', '-O1', '-I' + str(ROOT / 'inc'), '-I' + str(cls.directory),
                   *shlex.split(os.environ.get('ULAB_TEST_CFLAGS', '')), str(ROOT / 'tests/webapp/journal_probe.c'),
                   *sources, '-Wl,--gc-sections', '-o', str(cls.binary),
                   *shlex.split(os.environ.get('ULAB_TEST_LIBS', '-ljansson -lyaml -lm'))]
        subprocess.run(command, check=True)

    @classmethod
    def tearDownClass(cls):
        cls.tmp.cleanup()

    def probe(self, mode, *args):
        case = tempfile.TemporaryDirectory(dir=self.directory)
        self.addCleanup(case.cleanup)
        directory = Path(case.name)
        r = subprocess.run([str(self.binary), mode, str(directory), *map(str, args or ['unused'])],
                           cwd=ROOT, capture_output=True, text=True, timeout=15)
        self.assertEqual(r.returncode, 0, r.stdout + r.stderr)
        return directory

    def test_owned_resources_only_and_reverse_dependency_cleanup(self):
        d = self.probe('success')
        self.assertEqual((d / 'cleanup.log').read_text().splitlines(), ['site site-id', 'network network-id'])
        data = json.loads((d / 'webapp-resources.json').read_text())
        self.assertEqual([r['cleanup'] for r in data['resources']], ['deleted', 'deleted'])
        self.assertFalse(data['observations'][0]['owned'])

    def test_partial_cleanup_failure_retains_ids_and_continues(self):
        d = self.probe('partial')
        self.assertEqual((d / 'cleanup.log').read_text().splitlines(), ['site site-id', 'network network-id'])
        data = json.loads((d / 'webapp-resources.json').read_text())
        self.assertEqual([r['cleanup'] for r in data['resources']], ['deleted', 'failed'])
        self.assertEqual(data['resources'][1]['id'], 'site-id')

    def test_hung_cleanup_is_killed_and_later_resources_record_budget_exhaustion(self):
        d = self.probe('timeout')
        data = json.loads((d / 'webapp-resources.json').read_text())
        self.assertTrue(all(r['cleanup'] == 'failed' for r in data['resources']))

    def test_runtime_and_browser_events_are_serial(self):
        d = self.probe('runtime-ok', ROOT / 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml', ROOT / 'tests/webapp/fake_worker.py')
        data = json.loads((d / 'report.json').read_text())
        self.assertTrue(data['passed'])
        self.assertEqual([r['name'] for r in data['results']], ['web_open', 'web_action_available', 'disconnect_nodes', 'web_reload', 'web_action_available'])
        self.assertEqual((d / 'cleanup.log').read_text().splitlines(), ['runtime disconnect', 'runtime cleanup'])

    def test_runtime_failure_stops_browser_and_executes_final_cleanup(self):
        d = self.probe('runtime-fail', ROOT / 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml', ROOT / 'tests/webapp/fake_worker.py')
        data = json.loads((d / 'report.json').read_text())
        self.assertFalse(data['passed'])
        self.assertEqual(data['checks']['total'], 1)
        self.assertEqual((d / 'cleanup.log').read_text().splitlines(), ['runtime disconnect', 'runtime cleanup'])
        requests = [json.loads(x)['message'] for x in (d / 'webapp-commands.jsonl').read_text().splitlines() if json.loads(x)['direction'] == 'request']
        self.assertEqual(requests[-1]['action'], 'close')
        self.assertTrue(requests[-1]['inputs']['failed'])

    def test_cancellation_during_cleanup_finishes_cleanup_but_cannot_pass(self):
        d = self.probe('runtime-cancel-cleanup', ROOT / 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml', ROOT / 'tests/webapp/fake_worker.py')
        data = json.loads((d / 'report.json').read_text())
        self.assertEqual(data['outcome'], 'FAIL')
        self.assertEqual(data['cleanup'], 'ok')
        self.assertIn('cancelled during cleanup', data['error'])
        self.assertEqual(json.loads((d / 'webapp-resources.json').read_text())['run_result'], 'failed')

    def test_v2_dispatch_uses_kind_not_schema_version(self):
        self.probe('classify', ROOT / 'scenarios/workload/wl-001-console-single-site.yaml', '1')
        self.probe('classify', ROOT / 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml', '0')
        for name, text, expected in [
            ('quoted', 'version: 2\nkind: "workload"\n', '1'),
            ('json', '{"kind": "workload", "version": 2}', '1'),
            ('description', 'version: 2\ndescription: \'"kind": workload\'\n', '0'),
        ]:
            f = self.directory / (name + '.yaml'); f.write_text(text)
            self.probe('classify', f, expected)


if __name__ == '__main__':
    unittest.main()
