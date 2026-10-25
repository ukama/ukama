"""Node-state writer and shell readiness regression coverage.
SPDX-License-Identifier: MPL-2.0
"""
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
AFFECTED = {83, 92, 214, 215, 216, 218, 219, 225, 229, 230, 232, 233, 234, 235, 262}


class NodeState(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.build = tempfile.TemporaryDirectory(prefix='ulab-node-state-build-')
        cls.addClassCleanup(cls.build.cleanup)
        build = Path(cls.build.name)
        (build / 'version.h').write_text('#define VERSION "node-state-test"\n')
        cls.binary = build / 'probe'
        subprocess.run([
            *shlex.split(os.environ.get('CC', 'cc')),
            '-std=gnu11', '-D_POSIX_C_SOURCE=200809L', '-O0',
            '-Wall', '-Wextra', '-Werror', '-Wdeclaration-after-statement',
            '-ffunction-sections', '-fdata-sections',
            '-I' + str(ROOT / 'inc'), '-I' + str(build),
            str(ROOT / 'tests/runner/runtime_state_probe.c'),
            str(ROOT / 'src/runtime.c'), str(ROOT / 'src/util.c'),
            '-Wl,--gc-sections', '-o', str(cls.binary),
        ], check=True)

    def setUp(self):
        tmp = tempfile.TemporaryDirectory(prefix='ulab-node-state-')
        self.addCleanup(tmp.cleanup)
        self.directory = Path(tmp.name)
        self.run_dir = self.directory / 'run'
        self.prepare_run(self.run_dir)
        mock_bin = self.directory / 'bin'
        mock_bin.mkdir()
        podman = mock_bin / 'podman'
        podman.write_text('''#!/bin/sh
case "$1" in
    inspect) printf 'true\\n' ;;
    exec) printf '{"state":"READY"}\\n' ;;
    *) exit 1 ;;
esac
''')
        podman.chmod(0o755)
        self.env = {**os.environ, 'PATH': str(mock_bin) + ':' + os.environ['PATH'],
                    'ULAB_NODE_READY_TIMEOUT_SEC': '0'}

    def prepare_run(self, directory):
        (directory / 'runtime-sites').mkdir(parents=True)
        (directory / 'runtime-nodes').mkdir()
        (directory / 'runtime-sites/site-001.env').write_text(
            'TNODE_ID=factory-tower\nCNODE_ID=factory-controller\n'
            'ANODE_ID=factory-amplifier\nTNODE_CONTAINER=container-tower\n'
            'CNODE_CONTAINER=container-controller\nANODE_CONTAINER=container-amplifier\n')

    def write_state(self, *ids, run_dir=None):
        return subprocess.run([str(self.binary), str(run_dir or self.run_dir), *ids],
                              capture_output=True, text=True, timeout=10)

    def ready(self, node_id):
        result = subprocess.run(['sh', str(ROOT / 'scripts/wait-nodes-ready.sh'),
                                 node_id, str(self.run_dir)], env=self.env,
                                capture_output=True, text=True, timeout=10)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('node-ready logical=' + node_id, result.stdout)

    def test_reported_scenarios_reach_readiness_with_full_ids(self):
        scenarios = sorted(p for p in (ROOT / 'scenarios/resilience').rglob('r*.yaml')
                           if int(p.name[1:4]) in AFFECTED)
        self.assertEqual(len(scenarios), 15)
        for scenario in scenarios:
            slug = str(scenario.relative_to(ROOT / 'scenarios/resilience'))[:-5].replace('/', '-')
            run_id = 'p0-20261008t164402z-' + slug
            ids = [f'{run_id}-{kind}-site-001-001'
                   for kind in ('tower', 'amplifier', 'controller')]
            with self.subTest(scenario=scenario.name):
                result = self.write_state(*ids)
                self.assertEqual(result.returncode, 0, result.stderr)
                for node_id in ids:
                    self.assertTrue((self.run_dir / 'runtime-nodes' / (node_id + '.env')).is_file())
                    self.ready(node_id)

    def test_filename_boundaries(self):
        for length in (1, 126, 127, 128, 129, 251):
            node_id = 'n' * length
            with self.subTest(length=length):
                result = self.write_state(node_id)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.ready(node_id)

    def test_common_127_byte_prefix_stays_distinct(self):
        ids = ['n' * 127 + suffix for suffix in ('-tower', '-controller', '-amplifier')]
        result = self.write_state(*ids)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(len(list((self.run_dir / 'runtime-nodes').glob('*.env'))), 3)
        for node_id in ids:
            self.ready(node_id)

    def test_sanitization_matches_shell(self):
        node_id = 'node/caf\u00e9:ABC_09.-'
        result = self.write_state(node_id)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.ready(node_id)

    def test_sanitized_collision_preserves_original(self):
        result = self.write_state('node/name')
        self.assertEqual(result.returncode, 0, result.stderr)
        path = self.run_dir / 'runtime-nodes/node-name.env'
        original = path.read_bytes()
        result = self.write_state('node-name')
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('filename collision', result.stderr)
        self.assertIn('node/name', result.stderr)
        self.assertIn('node-name', result.stderr)
        self.assertEqual(path.read_bytes(), original)

    def test_same_identity_can_be_rewritten(self):
        node_id = 'n' * 128
        for _ in range(2):
            result = self.write_state(node_id)
            self.assertEqual(result.returncode, 0, result.stderr)
        self.ready(node_id)

    def test_unverifiable_existing_state_is_preserved(self):
        path = self.run_dir / 'runtime-nodes/node.env'
        for contents in ('', 'CONTAINER_NAME=other\n'):
            with self.subTest(contents=contents):
                path.write_text(contents)
                result = self.write_state('node')
                self.assertNotEqual(result.returncode, 0)
                self.assertIn('cannot verify existing runtime node state', result.stderr)
                self.assertEqual(path.read_text(), contents)

    def test_nonregular_state_is_rejected(self):
        path = self.run_dir / 'runtime-nodes/node.env'
        path.symlink_to(self.directory / 'missing')
        result = self.write_state('node')
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('cannot verify existing runtime node state', result.stderr)
        self.assertFalse((self.directory / 'missing').exists())

    def test_empty_and_overlong_filenames_fail_without_truncation(self):
        for length in (0, 252, 511):
            with self.subTest(length=length):
                result = self.write_state('n' * length)
                self.assertNotEqual(result.returncode, 0)
                self.assertIn('runtime node state filename', result.stderr)
        self.assertEqual(list((self.run_dir / 'runtime-nodes').iterdir()), [])

    def test_full_path_overflow_is_reported(self):
        run_dir = self.directory
        while len(str(run_dir)) < 900:
            run_dir /= 'd' * 100
        self.prepare_run(run_dir)
        result = self.write_state('n' * 200, run_dir=run_dir)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('runtime node state path too long', result.stderr)
        self.assertEqual(list((run_dir / 'runtime-nodes').iterdir()), [])


if __name__ == '__main__':
    unittest.main()
