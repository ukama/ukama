"""Compatibility checks for shared scripts used by ordinary local scenarios."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]


class LocalCompatibility(unittest.TestCase):
    def image_block(self, name, marker, image_exists, distributed=0):
        kind = 'ue' if name == 'start-ue.sh' else 'media'
        source = (ROOT / 'scripts' / name).read_text()
        start = source.index('if [ "${ULAB_DISTRIBUTED_RUNNER:-0}" = 1 ]')
        end = source.index('\nCSV=', start) if kind == 'ue' else source.index('\nfound=0', start)
        with tempfile.TemporaryDirectory() as t:
            root = Path(t); state = root / 'state'; state.mkdir()
            sentinel = state / f'.{kind}-image-built'
            if marker: sentinel.touch()
            env = {**os.environ, 'UE_STATE_DIR': str(state), 'MEDIA_STATE_DIR': str(state),
                   'UE_IMAGE': 'ukama/ue:dev', 'MEDIA_IMAGE': 'ukama/media:dev',
                   'UE_DIR': str(root),
                   'MOCK_EXISTS': str(int(image_exists)), 'TRACE': str(root / 'trace'),
                   'ULAB_DISTRIBUTED_RUNNER': str(distributed)}
            mock = '''set -eu
podman() {
    printf '%s\\n' "$*" >>"$TRACE"
    if [ "$1 $2" = 'image exists' ]; then [ "$MOCK_EXISTS" = 1 ]; else return 0; fi
}
'''
            result = subprocess.run(['sh', '-c', mock + source[start:end]], env=env, capture_output=True, text=True)
            calls = (root/'trace').read_text().splitlines() if (root/'trace').exists() else []
            return result.returncode, any(c.startswith('build ') for c in calls)

    def test_local_ue_preserves_per_run_rebuild(self):
        for marker in (False, True):
            for image in (False, True):
                with self.subTest(marker=marker, image=image):
                    rc, built = self.image_block('start-ue.sh', marker, image)
                    self.assertEqual(rc, 0)
                    self.assertEqual(built, not marker)

    def test_local_media_preserves_per_run_rebuild(self):
        for marker in (False, True):
            for image in (False, True):
                with self.subTest(marker=marker, image=image):
                    rc, built = self.image_block('start-media.sh', marker, image)
                    self.assertEqual(rc, 0)
                    self.assertEqual(built, not marker or not image)

    def test_worker_reuses_cached_ue_media_and_builds_only_missing(self):
        for name in ('start-ue.sh', 'start-media.sh'):
            for exists in (False, True):
                rc, built = self.image_block(name, False, exists, distributed=1)
                self.assertEqual(rc, 0)
                self.assertEqual(built, not exists)

    def test_stamping_network_restriction_is_opt_in(self):
        with tempfile.TemporaryDirectory() as t:
            root = Path(t); executable = root/'podman'
            executable.write_text('''#!/usr/bin/env python3
import json,os,sys
args=sys.argv[1:]
with open(os.environ['TRACE'],'a') as f: f.write(json.dumps(args)+'\\n')
if args[:2]==['image','exists']: sys.exit(0 if args[2]=='test-base' else 1)
''')
            executable.chmod(0o755)
            for distributed in (0, 1):
                trace = root/f'trace-{distributed}'
                env = {**os.environ, 'PATH': str(root)+':'+os.environ['PATH'], 'TRACE': str(trace),
                       'ULAB_DISTRIBUTED_RUNNER': str(distributed)}
                result = subprocess.run(['sh', str(ROOT/'scripts/stamp-node-image.sh'),
                                         'test-base', 'a-tnode-1', 'test-output'],
                                        env=env, capture_output=True, text=True)
                self.assertEqual(result.returncode, 0, result.stderr)
                args = next(json.loads(l) for l in trace.read_text().splitlines() if json.loads(l)[0]=='run')
                self.assertEqual('--network' in args, bool(distributed))
                self.assertEqual('--systemd=false' in args, bool(distributed))
