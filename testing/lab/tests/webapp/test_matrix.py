"""Adversarial evidence/report checks. SPDX-License-Identifier: MPL-2.0"""
import copy
import importlib.util
import json
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('matrix', ROOT / 'utils/webapp/matrix.py')
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)

class CoverageTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.scenario = 'scenarios/webapp/p0/a.yaml'
        p = self.root / self.scenario
        p.parent.mkdir(parents=True)
        p.write_text('''version: 2
suite: webapp
name: example
status: active
webapp:
  browser: chromium
  scenario_timeout_seconds: 1
phases:
  - name: baseline
    events:
      - type: web_open
    checks:
      - type: web_ui_equals
        requirement: WEB-TEST-001
        label: Visible
        expected: 'true'
''')
        catalog = self.root / m.CATALOG
        catalog.parent.mkdir(parents=True)
        catalog.write_text(json.dumps({'requirements': [{'id':'WEB-TEST-001','automation':'implemented','priority':'p0','requirement':'Visible <script>alert(1)</script>','scenarios':[self.scenario]}]}))
        for path in ['utils/webapp-worker.sh','utils/webapp/matrix.py','adapters/webapp/package-lock.json']:
            f=self.root/path;f.parent.mkdir(parents=True,exist_ok=True);f.write_text('test')
        self.report = {'run_id':'run-1','scenario':'example','status':'active','suite':'webapp','browser':'chromium','outcome':'PASS','passed':True,'final_rc':0,'cleanup':'ok',
                       'checks':{'total':1,'passed':1,'failed':0},'events':{'total':1,'passed':1,'failed':0},
                       'results':[{'kind':'event','phase':'baseline','name':'web_open','state':'PASS'},
                                  {'kind':'check','phase':'baseline','name':'web_ui_equals','requirement':'WEB-TEST-001','label':'Visible','state':'PASS','expected':'true','actual':'true'}]}
        self.manifest = {'schema_version':1,'evidence':'live','app_build':'app1','backend_build':'backend1','source_sha256':m.fingerprint(self.root),'inventory_sha256':m.digest(catalog),
                         'attempts':[{'scenario':self.scenario,'browser':'chromium','state':'complete','exit_code':0,'run_id':'run-1','report':'report.json','scenario_sha256':m.digest(p),'started_ns':1,'environment_sha256':m.environment_hash(p.read_text())}]}
        self.write()

    def write(self, folder='one'):
        out=self.root/folder;out.mkdir(exist_ok=True)
        generated=out/'case.yaml';generated.write_text(m.generated_case((self.root/self.scenario).read_text(),self.manifest['attempts'][0]['browser']))
        self.manifest['attempts'][0].update(generated_case='case.yaml',generated_sha256=m.digest(generated))
        m.write_json(out/'report.json',self.report)
        self.manifest['attempts'][0]['report_sha256']=m.digest(out/'report.json')
        m.write_json(out/'matrix.json',self.manifest)
        return out/'matrix.json'

    def coverage(self, paths=None, browsers=None):
        return m.coverage(self.root,paths or [self.write()],browsers or ['chromium'],'app1','backend1')

    def test_complete_evidence_and_html_escaping(self):
        r=self.coverage();self.assertTrue(r['gate_passed']);self.assertEqual(r['totals']['verified'],1)
        page=m.render_html(r);self.assertNotIn('<script>',page);self.assertIn('&lt;script&gt;',page)

    def test_fixture_and_unverified_never_get_live_credit(self):
        for evidence in ['fixture','unverified','unknown']:
            self.manifest['evidence']=evidence
            self.assertEqual(self.coverage()['totals']['verified'],0)

    def test_partial_planned_requirement_never_gets_credit(self):
        p=self.root/m.CATALOG;catalog=json.loads(p.read_text());catalog['requirements'][0]['automation']='planned';p.write_text(json.dumps(catalog))
        self.manifest.update(inventory_sha256=m.digest(p),source_sha256=m.fingerprint(self.root))
        self.assertEqual(self.coverage()['totals']['verified'],0)

    def test_empty_partial_duplicate_null_and_failed_assertions_are_rejected(self):
        original=copy.deepcopy(self.report)
        variants=[[],original['results'][:1],original['results']+[original['results'][1]]]
        for field,value in [('state','FAIL'),('actual',None),('actual','false'),('requirement','WEB-UNKNOWN'),('label','Different')]:
            v=copy.deepcopy(original['results']);v[1][field]=value;variants.append(v)
        for rows in variants:
            with self.subTest(rows=rows):
                self.report['results']=rows;self.assertEqual(self.coverage()['totals']['verified'],0)

    def test_failed_cleanup_skips_and_wrong_report_are_rejected(self):
        for key,value in [('cleanup','failed'),('outcome','SKIP'),('status','skip'),('run_id','wrong'),('browser','webkit'),('final_rc',1),('passed',False),('checks',{'total':0,'passed':0,'failed':0})]:
            with self.subTest(key=key):
                old=self.report[key];self.report[key]=value;self.assertEqual(self.coverage()['totals']['verified'],0);self.report[key]=old

    def test_stale_or_tampered_evidence_is_rejected(self):
        path=self.write();(path.parent/'report.json').write_text('{}')
        self.assertEqual(self.coverage([path])['totals']['verified'],0)
        self.write();(self.root/self.scenario).write_text((self.root/self.scenario).read_text()+'\n# newer\n')
        self.assertEqual(self.coverage()['totals']['verified'],0)

    def test_latest_failure_and_flaky_later_pass_withhold_credit(self):
        first=self.write('first');self.manifest['attempts'][0].update(started_ns=2,exit_code=1)
        second=self.write('second');self.assertEqual(self.coverage([first,second])['totals']['verified'],0)
        self.manifest['attempts'][0].update(started_ns=3,exit_code=0)
        third=self.write('third');r=self.coverage([first,second,third]);self.assertEqual(r['totals']['verified'],0);self.assertTrue(r['requirements'][0]['evidence'][0]['flaky'])

    def test_all_browsers_and_same_build_are_required(self):
        self.assertEqual(self.coverage(browsers=['chromium','firefox'])['totals']['verified'],0)
        self.manifest['app_build']='another';self.assertEqual(self.coverage()['totals']['verified'],0)

    def test_missing_or_failed_attempt_cannot_reuse_earlier_pass(self):
        first=self.write('first')
        for state in ['pending','running','skipped','launch_failed']:
            self.manifest['attempts'][0].update(started_ns=2,state=state)
            second=self.write('second');self.assertEqual(self.coverage([first,second])['totals']['verified'],0)

    def test_catalog_mapping_and_duplicate_yaml_are_validated(self):
        p=self.root/self.scenario;p.write_text(p.read_text().replace('WEB-TEST-001','WEB-UNKNOWN'))
        with self.assertRaises(ValueError):self.coverage()
        p.write_text('version: 2\nversion: 2\nsuite: webapp\n')
        with self.assertRaises(ValueError):m.load_case(p)

    def test_report_command_enforces_gate_and_refuses_overwrite(self):
        out=self.root/'report-output'
        cmd=['python3',str(ROOT/'utils/webapp/matrix.py'),'report','--root',str(self.root),'--out',str(out),'--browser','chromium','--gate']
        p=subprocess.run(cmd,capture_output=True,text=True);self.assertEqual(p.returncode,1,p.stderr);self.assertTrue((out/'coverage.html').is_file())
        p=subprocess.run(cmd,capture_output=True,text=True);self.assertEqual(p.returncode,2)

    def test_matrix_records_launch_failure_and_private_artifacts(self):
        binary=self.root/'badbinary';binary.write_text('invalid executable');binary.chmod(0o700)
        out=self.root/'matrix-output'
        cmd=['python3',str(ROOT/'utils/webapp/matrix.py'),'run','--root',str(self.root),'--out',str(out),'--binary',str(binary),'--browser','chromium']
        p=subprocess.run(cmd,capture_output=True,text=True);self.assertEqual(p.returncode,1,p.stderr)
        manifest=json.loads((out/'matrix.json').read_text());self.assertEqual(manifest['attempts'][0]['state'],'launch_failed');self.assertEqual((out/'matrix.json').stat().st_mode&0o777,0o600)
        self.assertEqual(manifest['evidence'],'unverified')


class MatrixProcessTest(unittest.TestCase):
    setUp = CoverageTest.setUp
    write = CoverageTest.write
    def test_matrix_profiles_and_process_results_are_recorded(self):
        binary=self.root/'fake lab';binary.write_text('''#!/usr/bin/env python3
import json,sys,pathlib,yaml
args=sys.argv;source=pathlib.Path(args[2]);case=yaml.safe_load(source.read_text());run=args[args.index('--run-id')+1];out=pathlib.Path(args[args.index('--out')+1])/run;out.mkdir(parents=True)
(out/'report.json').write_text(json.dumps({'scenario':case['name'],'browser':case['webapp']['browser'],'run_id':run}))
''');binary.chmod(0o700)
        out=self.root/'profiles'
        p=subprocess.run(['python3',str(ROOT/'utils/webapp/matrix.py'),'run','--root',str(self.root),'--out',str(out),'--binary',str(binary),'--browser','chromium','--browser','firefox','--evidence','fixture','--','--repo','/path with spaces'],capture_output=True,text=True)
        self.assertEqual(p.returncode,0,p.stderr)
        manifest=json.loads((out/'matrix.json').read_text());self.assertEqual([a['browser'] for a in manifest['attempts']],['chromium','firefox']);self.assertEqual(len({a['run_id'] for a in manifest['attempts']}),2)
        for a in manifest['attempts']:
            self.assertEqual(a['exit_code'],0);self.assertEqual(a['report_sha256'],m.digest(out/a['report']));self.assertEqual(m.load_case(out/a['generated_case'])['webapp']['browser'],a['browser'])

if __name__=='__main__':unittest.main()
