import importlib.util
import json
from pathlib import Path
import subprocess
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('health', ROOT/'utils/runner/worker-health.py')
health = importlib.util.module_from_spec(spec); spec.loader.exec_module(health)


class Health(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(); self.root = Path(self.tmp.name)
        (self.root/'runtime-systemd').mkdir()
        (self.root/'runtime-systemd/node.service').touch()
        self.containers=[]; self.networks=[]; self.active='inactive'; self.commands=[]
        def run(args, **kw):
            self.commands.append(args)
            if args[:2] == ('podman','ps'): value=json.dumps(self.containers)
            elif args[:3] == ('podman','network','ls'): value=json.dumps(self.networks)
            elif 'ActiveState' in ' '.join(args): value=self.active
            else: value='ok'
            return SimpleNamespace(stdout=value)
        self.mock=patch.object(health.subprocess,'run',side_effect=run); self.mock.start()
        self.disk=patch.object(health.os,'statvfs',return_value=SimpleNamespace(f_bavail=3*1024**3,f_frsize=1)); self.disk.start()
        self.env=patch.dict(health.os.environ, {}, clear=True); self.env.start()
    def tearDown(self): self.env.stop(); self.disk.stop(); self.mock.stop(); self.tmp.cleanup()
    def check(self): return health.check(self.root,'r-job',{'site-1':{'tnode':'a-tnode-1'}})
    def test_clean_worker_is_reusable_without_any_scenario_verdict(self):
        self.assertTrue(self.check()['reusable'])
        self.assertFalse(any(set(c)&{'rm','delete','del-br','stop','start'} for c in self.commands))
    def test_own_container_detected_but_unrelated_container_ignored(self):
        self.containers=[{'Names':['unrelated-container']}]
        self.assertTrue(self.check()['reusable'])
        self.containers.append({'Names':['ukama-vnode-a-tnode-1']})
        self.assertFalse(self.check()['reusable'])
    def test_owned_network_requires_retirement(self):
        self.networks=[{'name':'ukama-lab-r-job'}]
        self.assertFalse(self.check()['reusable'])
    def test_still_active_systemd_unit_requires_retirement(self):
        self.active='active'; self.assertFalse(self.check()['reusable'])
    def test_broken_systemd_or_podman_is_not_reusable(self):
        with patch.object(health.subprocess,'run',side_effect=subprocess.TimeoutExpired('systemctl',30)):
            self.assertFalse(self.check()['reusable'])
    def test_disk_pressure_prevents_next_scenario(self):
        with patch.object(health.os,'statvfs',return_value=SimpleNamespace(f_bavail=100,f_frsize=1)):
            self.assertFalse(self.check()['reusable'])
