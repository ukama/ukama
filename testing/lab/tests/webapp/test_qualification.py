"""Read-only readiness, secret handling and frozen final gate. MPL-2.0."""
import argparse
import contextlib
import io
import json
import os
from pathlib import Path
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT=Path(__file__).resolve().parents[2]
sys.path.insert(0,str(ROOT/'utils/webapp'))
import qualify as q

class QualificationTest(unittest.TestCase):
    def test_full_inventory_is_frozen_and_contains_119_scenarios(self):
        catalog=q.inventory(ROOT);self.assertEqual(len(catalog['requirements']),132)
        rows=q.input_inventory(ROOT,catalog);self.assertEqual(len(rows),119)
        self.assertTrue(any(r['controlled'] for r in rows))

    def test_auth_structure_checks_do_not_claim_a_live_session(self):
        with tempfile.TemporaryDirectory() as folder:
            root=Path(folder);self.assertEqual(q.auth_state(root,'none'),'invalid')
            self.assertEqual(q.auth_state(root,'none',True),'not_required')
            self.assertEqual(q.auth_state(root,'absent'),'missing_or_invalid')
            (root/'state').write_text('{"cookies":[],"origins":[]}')
            self.assertEqual(q.auth_state(root,'state'),'present_unverified')
            (root/'state').write_text('{"cookies":"secret","origins":[]}')
            self.assertEqual(q.auth_state(root,'state'),'invalid')

    def test_endpoint_probe_is_tcp_only_and_rejects_credential_urls(self):
        for value in ('http://user:secret@localhost','http://localhost?token=secret','file:///tmp/state','http://localhost/#token'):
            self.assertEqual(q.endpoint(value,True),'invalid_url')
        with patch.object(q.socket,'create_connection') as connect:
            self.assertEqual(q.endpoint('http://localhost:3000'),'not_probed');connect.assert_not_called()
            self.assertEqual(q.endpoint('http://localhost:3000',True),'tcp_open')
            connect.assert_called_once_with(('localhost',3000),timeout=1)

    def test_browser_probe_missing_or_malformed_results_cannot_pass(self):
        for value in ('{}','{"browsers":[]}','not json'):
            with patch.object(q.subprocess,'run',return_value=argparse.Namespace(stdout=value)):
                rows=q.probe_browsers(ROOT,['chromium','firefox'])
                self.assertEqual([r['state'] for r in rows],['probe_failed','probe_failed'])

    def test_preflight_records_missing_inputs_and_never_executes_run_or_copies_values(self):
        with tempfile.TemporaryDirectory() as folder:
            root=Path(folder);source=root/'scenarios/webapp/p0/test.yaml';source.parent.mkdir(parents=True)
            source.write_text('''version: 2
suite: webapp
status: active
webapp:
  auth_state: ${ULAB_TEST_SECRET}
  base_url: http://localhost:3000
phases: []
''')
            catalog={'requirements':[{'scenarios':[str(source.relative_to(root))],'automation':'planned','priority':'p0'}]}
            binary=root/'lab';binary.write_text('#!/bin/sh\nexit 0\n');binary.chmod(0o700)
            args=argparse.Namespace(root=root,out=root/'out',binary=binary,app_build='',backend_build='',repo=root,scripts=None,
                                    browser=['chromium','firefox'],probe_endpoints=False,warehouse_url=None,factory_url=None,asr_url=None,bff=None)
            secret='DO-NOT-COPY-AUTH-PATH'
            with patch.object(q,'inventory',return_value=catalog),patch.object(q.m,'fingerprint',return_value='hash'),patch.object(q,'probe_browsers',return_value=[{'browser':'chromium','state':'ready'},{'browser':'firefox','state':'unavailable'}]),patch.dict(os.environ,{'ULAB_TEST_SECRET':secret}),patch.object(q.subprocess,'run',return_value=argparse.Namespace(returncode=0)) as process,contextlib.redirect_stdout(io.StringIO()):
                self.assertEqual(q.preflight(args),1)
                self.assertTrue(process.called)
                self.assertTrue(all(c.args[0][1]=='validate' for c in process.call_args_list))
            output=(args.out/'preflight.json').read_text()+(args.out/'inputs.env.example').read_text()
            self.assertNotIn(secret,output)
            data=json.loads((args.out/'preflight.json').read_text());self.assertFalse(data['execution_ready'])
            self.assertEqual(data['attempt_count'],2);self.assertEqual(data['live_verified'],0)

    def test_closeout_without_evidence_keeps_full_denominator_and_exits_failure(self):
        with tempfile.TemporaryDirectory() as folder,contextlib.redirect_stdout(io.StringIO()):
            out=Path(folder)/'closeout';args=argparse.Namespace(root=ROOT,out=out,manifest=[],browser=list(q.m.BROWSERS),app_build='',backend_build='')
            self.assertEqual(q.closeout(args),1)
            result=json.loads((out/'closeout.json').read_text())
            self.assertEqual(result['totals'],{'total':132,'automated':81,'verified':0,'verified_percent':0.0})
            self.assertEqual(result['p0']['total'],86);self.assertEqual(len(result['unimplemented_or_partial']),51)
            self.assertEqual(len(result['unverified_requirements']),132)
            with self.assertRaises(FileExistsError):q.closeout(args)

if __name__=='__main__':unittest.main()
