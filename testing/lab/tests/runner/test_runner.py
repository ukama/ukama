import copy
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import time
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'utils/runner'))
import controller
import protocol
spec = importlib.util.spec_from_file_location('worker_loop', ROOT / 'utils/runner/worker-loop.py')
worker = importlib.util.module_from_spec(spec)
spec.loader.exec_module(worker)


class MemoryStore:
    def __init__(self):
        self.data = {}; self.events = []; self.fail_sync = False
    def get(self, key):
        if key == 'controller/heartbeat.json':
            return {'updated_at': time.time()}
        return copy.deepcopy(self.data.get(key))
    def put(self, key, value, create=False):
        if create and key in self.data and self.data[key] != value:
            raise RuntimeError('conditional conflict')
        self.data[key] = copy.deepcopy(value)
        self.events.append(('put', key))
    def sync(self, source, key):
        self.events.append(('sync', key))
        if self.fail_sync:
            raise RuntimeError('upload failed')
    def keys(self, prefix):
        return [k for k in self.data if k.startswith(prefix)]
    def uri(self, key):
        return 's3://mock/batch/' + key


def job(ident='j-1', exclusive=False):
    return {'id': ident, 'run_id': 'r-' + ident, 'scenario': 'scenarios/p0/test.yaml', 'suite': 'p0',
            'exclusive': exclusive, 'selected': True, 'sites': 1,
            'bundles': {'site-001': {'tnode': 'a-tnode-1', 'cnode': 'a-cnode-1', 'anode': 'a-anode-1'}}}


class Scheduling(unittest.TestCase):
    def setUp(self):
        self.jobs = [job('a'), job('b'), job('exclusive', True), job('c')]
    def test_parallel_jobs_before_barrier(self):
        self.assertEqual([x['id'] for x in protocol.eligible(self.jobs, {}, set())], ['a', 'b'])
    def test_barrier_drains_fleet(self):
        self.assertEqual(protocol.eligible(self.jobs, {'a': {}}, {'b'}), [])
    def test_exclusive_alone(self):
        self.assertEqual([x['id'] for x in protocol.eligible(self.jobs, {'a': {}, 'b': {}}, set())], ['exclusive'])
    def test_no_jobs_during_exclusive(self):
        self.assertEqual(protocol.eligible(self.jobs, {'a': {}, 'b': {}}, {'exclusive'}), [])
    def test_resume_after_barrier(self):
        self.assertEqual([x['id'] for x in protocol.eligible(self.jobs, {'a': {}, 'b': {}, 'exclusive': {}}, set())], ['c'])
    def test_all_current_scenarios_plannable(self):
        paths = [*ROOT.glob('scenarios/p0/**/*.yaml'), *ROOT.glob('scenarios/resilience/**/*.yaml')]
        self.assertEqual(len(paths), 467)
        for p in paths:
            protocol.scenario_shape(p)
    def test_ids_are_unique_and_short_for_entire_suite(self):
        paths = [str(p.relative_to(ROOT)) for p in ROOT.glob('scenarios/p0/**/*.yaml')]
        jobs = controller.plan_jobs(paths, ROOT, 'p0-' + 'a' * 40)
        self.assertEqual(len({j['run_id'] for j in jobs}), len(jobs))
        self.assertLess(max(len(j['run_id']) for j in jobs), 40)
        self.assertEqual([j['scenario'] for j in jobs], paths)
    def test_known_global_scenarios_exclusive(self):
        paths = ['scenarios/p0/data-package/purchase/payment-service-failure-retry.yaml',
                 'scenarios/p0/console/analytics/rollup-retry-idempotency.yaml']
        self.assertTrue(all(j['exclusive'] for j in controller.plan_jobs(paths, ROOT, 'batch')))


class Reservations(unittest.TestCase):
    def setUp(self):
        self.payload = {'nodes': [{'nodeId': f'a-{kind}-{i}'} for i in range(3) for kind in ('tnode', 'cnode', 'anode')]}
    def test_complete_bundles_only(self):
        self.payload['nodes'].pop()
        self.assertEqual(len(protocol.bundles(self.payload)), 2)
    def test_provisioned_node_excluded(self):
        self.payload['nodes'][1]['isProvisioned'] = True
        self.assertEqual(len(protocol.bundles(self.payload)), 2)
    def test_disjoint_owned_bundles(self):
        jobs, journal, snapshots = [job('1'), job('2')], {}, []
        def api(base, path, method='GET'):
            return self.payload if method == 'GET' else {}
        with patch.object(protocol, 'factory_request', side_effect=api):
            protocol.reserve_jobs(jobs, 'https://mock', journal, lambda j: snapshots.append(copy.deepcopy(j)))
        ids = [n for j in jobs for b in j['bundles'].values() for n in b.values()]
        self.assertEqual(len(ids), len(set(ids)))
        self.assertEqual(journal['state'], 'RESERVED_IDS')
        self.assertEqual(len(snapshots), 1)
    def test_reservation_does_not_run_backend_setup(self):
        calls, journal = [], {}
        def api(base, path, method='GET'):
            calls.append(method)
            return self.payload
        with patch.object(protocol, 'factory_request', side_effect=api):
            protocol.reserve_jobs([job()], 'https://mock', journal, lambda j: None)
        self.assertEqual(calls, ['GET'])
        with self.assertRaises(RuntimeError):
            protocol.reserve_jobs([job()], 'https://mock', journal, lambda j: None)
    def test_shortage_before_any_patch(self):
        calls=[]
        def api(base, path, method='GET'):
            calls.append(method); return {'nodes': []}
        with patch.object(protocol, 'factory_request', side_effect=api), self.assertRaises(RuntimeError):
            protocol.reserve_jobs([job()], 'https://mock', {}, lambda j: None)
        self.assertEqual(calls, ['GET'])


class Reports(unittest.TestCase):
    def test_status_can_read_while_scheduler_owns_controller_lock(self):
        with tempfile.TemporaryDirectory() as tmp, patch.dict(os.environ, {'S3_BUCKET':'mock','S3_PREFIX':'test'}), patch.object(controller,'Store',return_value=MemoryStore()):
            owner=controller.Controller(Path(tmp),'batch')
            try:
                reader=controller.Controller(Path(tmp),'batch',readonly=True)
                reader.fd.close()
            finally:
                owner.fd.close()
    def test_valid_pass(self):
        self.assertEqual(protocol.classify({'ended_at': 1, 'cleanup': 'ok', 'passed': True}, 0)[0], 'PASS')
    def test_assertion_failure(self):
        self.assertEqual(protocol.classify({'ended_at': 1, 'cleanup': 'ok', 'passed': False}, 1)[0], 'FAIL')
    def test_incomplete_report_is_infra(self):
        self.assertEqual(protocol.classify({'passed': True}, 0)[0], 'INFRA')
    def test_timeout_cannot_be_pass(self):
        self.assertEqual(protocol.classify({'ended_at': 1, 'cleanup': 'ok', 'passed': True}, 0, 'scenario_deadline')[0], 'INFRA')
    def test_cleanup_failure_preserves_scenario_fail(self):
        self.assertEqual(protocol.classify({'ended_at': 1, 'cleanup': 'failed', 'passed': False}, 1)[0], 'FAIL')
    def test_skip(self):
        self.assertEqual(protocol.classify({'ended_at': 1, 'cleanup': 'ok', 'status': 'skip'}, 0)[0], 'SKIP')
    def test_missing_job_never_vanishes(self):
        p={'batch_id': 'test','artifact':'hash','jobs':[job('a'),job('b')]}
        r=controller.summarize(p, {'a': {'outcome': 'PASS'}})
        self.assertEqual(r['counts'], {'PASS': 1, 'PENDING': 1})


class WorkerExecution(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name); self.lab = self.root/'ukama-lab'
        (self.lab/'utils').mkdir(parents=True)
        script = self.lab/'utils/run-scenarios.sh'
        script.write_text('''#!/usr/bin/env python3
import json, os, pathlib, sys, time
if os.environ.get('MOCK_HANG'): time.sleep(60)
p=pathlib.Path(os.environ['P0_RUNS_DIR'])/'j-1'/'runs'/os.environ['ULAB_RUN_ID']
p.mkdir(parents=True)
r={'run_id':os.environ['ULAB_RUN_ID'],'passed':True,'ended_at':time.time(),'cleanup':'ok'}
if not os.environ.get('MOCK_NO_REPORT'): (p/'report.json').write_text(json.dumps(r))
''')
        script.chmod(0o755)
        self.env = patch.dict(os.environ, {'WORKER_ID':'worker-01', 'UKAMA_REPO':str(self.root/'ukama'),
                                          'P0_SCENARIO_TIMEOUT_MINUTES':'1','P0_CLEANUP_TIMEOUT_SECONDS':'2'})
        self.env.start()
        self.p1=patch.object(worker,'ROOT',self.root);self.p2=patch.object(worker,'LAB',self.lab)
        self.p1.start();self.p2.start(); worker.STOP.clear()
        self.store=MemoryStore()
        self.cmd={'job':job(),'sequence':1,'artifact':'hash'}
    def tearDown(self):
        worker.STOP.clear();self.p1.stop();self.p2.stop();self.env.stop();self.tmp.cleanup()
    def test_evidence_precedes_terminal_commit(self):
        row=worker.execute(self.cmd,self.store,'hash',1)
        self.assertEqual(row['outcome'],'PASS')
        self.assertEqual(self.store.events[-2][0],'sync')
        self.assertEqual(self.store.events[-1],('put','results/j-1.json'))
    def test_upload_failure_no_terminal_commit(self):
        self.store.fail_sync=True
        row = worker.execute(self.cmd,self.store,'hash',1)
        self.assertEqual(row['outcome'], 'PASS')
        self.assertNotIn('results/j-1.json', self.store.data)
        self.assertTrue((self.root/'outbox/j-1.json').is_file())
    def test_pending_upload_recovers_without_reexecution(self):
        self.store.fail_sync=True
        original=worker.execute(self.cmd,self.store,'hash',1)
        self.store.fail_sync=False
        worker.flush_pending(self.store)
        self.assertEqual(self.store.data['results/j-1.json'], original)
        self.assertFalse((self.root/'outbox/j-1.json').exists())
    def test_control_read_failure_does_not_kill_scenario(self):
        self.store.get=lambda key: (_ for _ in ()).throw(RuntimeError('temporary S3 outage'))
        row=worker.execute(self.cmd,self.store,'hash',1)
        self.assertEqual(row['outcome'],'PASS')
        self.assertIn('temporary S3 outage',(self.root/'results/diagnostics/agent.jsonl').read_text())
    def test_absent_controller_heartbeat_does_not_kill_scenario(self):
        self.store.get=lambda key: None
        self.assertEqual(worker.execute(self.cmd,self.store,'hash',1)['outcome'],'PASS')
    def test_worker_restart_does_not_reexecute(self):
        worker.execute(self.cmd,self.store,'hash',1)
        with self.assertRaises(RuntimeError): worker.execute(self.cmd,self.store,'hash',1)
    def test_missing_report_commits_infra(self):
        with patch.dict(os.environ,{'MOCK_NO_REPORT':'1'}):
            self.assertEqual(worker.execute(self.cmd,self.store,'hash',1)['outcome'],'INFRA')
    def test_network_failure_terminates_and_commits_infra(self):
        (self.root/'network.failed').write_text('dns')
        with patch.dict(os.environ,{'MOCK_HANG':'1'}):
            row=worker.execute(self.cmd,self.store,'hash',1)
        self.assertEqual(row['reason'],'network_unhealthy')
        self.assertEqual(row['outcome'],'INFRA')
    def test_controller_stop_is_observed(self):
        self.store.data['control/stop.json']={'reason':'stop'}
        with patch.dict(os.environ,{'MOCK_HANG':'1'}):
            self.assertEqual(worker.execute(self.cmd,self.store,'hash',1)['reason'],'controller_stopped_batch')
    def test_batch_deadline_is_enforced_without_controller(self):
        with patch.dict(os.environ, {'MOCK_HANG':'1','P0_BATCH_END_AT':'1'}):
            self.assertEqual(worker.execute(self.cmd,self.store,'hash',1)['reason'],'batch_deadline')
    def test_failed_scenario_does_not_retire_healthy_worker(self):
        (self.root/'images.json').write_text(json.dumps({'fingerprint':'hash'}))
        for seq in (1, 2):
            self.store.data[f'commands/worker-01/{seq:06d}.json'] = {'protocol':2,'sequence':seq,'artifact':'hash','job':job('j-'+str(seq))}
        attempts=[]
        def execute(cmd, store, artifact, seq):
            attempts.append(cmd['job']['id'])
            row={**cmd['job'],'outcome':'FAIL','cleanup':'failed'}
            protocol.atomic_json(self.root/'results'/row['id']/'result.json',row)
            if seq==2: worker.STOP.set()
            return row
        with patch.dict(os.environ, {'S3_BUCKET':'mock','S3_PREFIX':'test','BATCH_ID':'test'}), patch.object(worker,'Store',return_value=self.store), patch.object(worker,'validate_images'), patch.object(worker,'execute',side_effect=execute), patch.object(worker,'worker_health',return_value={'reusable':True,'problems':[]}):
            worker.main()
        self.assertEqual(attempts,['j-1','j-2'])


class ShellGuards(unittest.TestCase):
    def test_missing_worker_base_fails_without_full_build(self):
        with tempfile.TemporaryDirectory() as t:
            p=Path(t)/'podman';p.write_text('#!/bin/sh\nexit 1\n');p.chmod(0o755)
            env={**os.environ,'PATH':t+':'+os.environ['PATH'],'ULAB_REQUIRE_NODE_BASE':'1'}
            r=subprocess.run(['sh',str(ROOT/'scripts/build-node.sh'),'/missing-ukama','a-tnode-1'],env=env,capture_output=True,text=True)
            self.assertNotEqual(r.returncode,0)
            self.assertIn('runtime builds are disabled',r.stderr)
    def test_all_modified_shell_syntax(self):
        for p in [*ROOT.glob('utils/runner/*.sh'),ROOT/'utils/run-scenarios.sh',*ROOT.glob('scripts/*.sh')]:
            result=subprocess.run(['bash','-n',str(p)],capture_output=True,text=True)
            self.assertEqual(result.returncode,0,f'{p}: {result.stderr}')


if __name__=='__main__': unittest.main()

class ControllerLoop(unittest.TestCase):
    def simulate(self, outcomes=None, bad_worker=False):
        outcomes = outcomes or {}
        clock = [1000.0]
        running, commands, violations = {}, [], []
        store = MemoryStore()
        c = controller.Controller.__new__(controller.Controller)
        c.directory = Path('/unused')
        c.verified_commands = set()
        c.store = store
        c.state = {'workers': {f'w{i}': {'instance':f'i{i}','launched':1000} for i in (1,2)},
                   'assignments': {}, 'sequences': {}, 'results': {}, 'started':1000}
        jobs = [job('a'), job('b'), job('x',True), job('c')]
        for j in jobs: j['scenario'] += '-' + j['id']
        c.plan = {'batch_id':'test','artifact':'hash','jobs':jobs,'workers':2,'max_replacements':0,
                  'env':{'P0_BATCH_TIMEOUT_MINUTES':'20','P0_PREP_TIMEOUT_MINUTES':'45','P0_CLEANUP_TIMEOUT_SECONDS':'1'}}
        c.save=lambda: None
        instance_states={'i1':'running','i2':'running'}
        c.instances=lambda: [{'InstanceId':k,'State':{'Name':v}} for k,v in instance_states.items()]
        original_put=store.put
        for w in c.state['workers']:
            store.data['workers/'+w+'.json']={'state':'READY','sequence':0,'updated_at':1000}
        def put(key,value,create=False):
            if key.startswith('commands/') and key not in store.data:
                w=key.split('/')[1]; j=value['job']
                if w in running: violations.append('worker overlap')
                if j['exclusive'] and running: violations.append('exclusive overlap')
                if any(v[0]['job']['exclusive'] for v in running.values()): violations.append('ran beside exclusive')
                running[w]=(value,clock[0]+(10 if j['id']=='a' else 5))
                commands.append(j['id'])
                store.data['workers/'+w+'.json']={'state':'RUNNING','sequence':value['sequence']-1,'updated_at':clock[0]}
            return original_put(key,value,create)
        store.put=put
        def tick(seconds):
            clock[0]+=seconds
            for w,(cmd,end) in list(running.items()):
                if clock[0]>=end:
                    j=cmd['job'];outcome=outcomes.get(j['id'],'PASS')
                    store.data['results/'+j['id']+'.json']={**j,'artifact':'hash','outcome':outcome,'cleanup':'ok'}
                    store.data['workers/'+w+'.json']={'state':'FAILED' if bad_worker and outcome=='INFRA' else 'READY',
                                                      'sequence':cmd['sequence'],'updated_at':clock[0]}
                    del running[w]
        def aws(*args,**kw):
            if args[:2]==('ec2','terminate-instances'):
                instance_states[args[-1]]='terminated'
                return b'{}'
            raise AssertionError(args)
        with patch.object(controller.time,'time',side_effect=lambda:clock[0]),patch.object(controller.time,'sleep',side_effect=tick),patch.object(protocol.Store,'aws',side_effect=aws):
            reason=c.watch()
        return reason,commands,violations,c
    def test_dispatch_and_ack_with_fleet_barrier(self):
        reason,commands,violations,c=self.simulate()
        self.assertEqual(reason,'');self.assertEqual(commands,['a','b','x','c']);self.assertEqual(violations,[])
        self.assertEqual(len(c.state['results']),4)
    def test_scenario_failure_does_not_gate_queue(self):
        reason,commands,violations,c=self.simulate({'a':'FAIL'})
        self.assertEqual(reason,'');self.assertEqual(len(commands),4);self.assertEqual(violations,[])
        self.assertEqual(c.state['results']['a']['outcome'],'FAIL')
    def test_failed_worker_retires_others_finish(self):
        reason,commands,violations,c=self.simulate({'a':'INFRA'},True)
        self.assertEqual(reason,'');self.assertEqual(len(commands),4);self.assertEqual(violations,[])
        self.assertTrue(any(w.get('retired') for w in c.state['workers'].values()))
