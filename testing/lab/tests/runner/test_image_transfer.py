"""Exercise cache reuse, transfer, and actual node scripts with a fake Podman CLI."""
import gzip
import hashlib
import json
import os
from pathlib import Path
import platform
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'utils/runner'))
import images
import controller
import protocol

# A content-addressed image store, without starting containers or touching AWS.
PODMAN = r'''#!/usr/bin/env python3
import gzip,hashlib,json,os,sys
from pathlib import Path
args=sys.argv[1:]
p=Path(os.environ['PODMAN_STATE'])
db=json.loads(p.read_text()) if p.exists() else {}
def name(value): return value if value.startswith('sha256:') or value.startswith('localhost/') else 'localhost/'+value
def lookup(value):
    v=name(value)
    return db.get(v) or next((i for i in db.values() if i['Id']==value),None)
with open(os.environ['PODMAN_TRACE'],'a') as f: f.write(json.dumps(args)+'\n')
if args[:2]==['image','exists']: sys.exit(0 if lookup(args[2]) else 1)
if args[:2]==['image','inspect']:
    i=lookup(args[2])
    if not i: sys.exit(1)
    print(json.dumps([i]));sys.exit(0)
if args[0]=='save':
    Path(args[args.index('--output')+1]).write_text(json.dumps(lookup(args[-1])))
elif args[0]=='load':
    with gzip.open(args[args.index('--input')+1],'rt') as f: i=json.load(f)
    db[i['Id']]=i
elif args[0]=='tag': db[name(args[2])]=lookup(args[1])
elif args[0]=='commit':
    i={'Id':'sha256:'+hashlib.sha256(args[-1].encode()).hexdigest(),'Architecture':'amd64','Os':'linux'}
    db[name(args[-1])]=i
elif args[0] not in ('run','exec','rm','save'):
    raise SystemExit('unexpected mock Podman command: '+str(args))
p.write_text(json.dumps(db))
'''
BUILDER = r'''#!/usr/bin/env python3
import hashlib,json,os,sys
from pathlib import Path
args=sys.argv[1:]; node=args[args.index('--node-id')+1]
p=Path(os.environ['PODMAN_STATE']); db=json.loads(p.read_text())
i={'Id':'sha256:'+hashlib.sha256(('built:'+node).encode()).hexdigest(),'Architecture':'amd64','Os':'linux'}
db['localhost/testing/virtualnode:'+node]=i;p.write_text(json.dumps(db))
with open(os.environ['BUILD_TRACE'],'a') as f: f.write(node+'\n')
'''


class ImageTransfer(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(); self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name); self.lab = self.root/'lab'; self.repo = self.root/'ukama'
        for d in ('bin', 'scripts', 'utils/runner', 'scenarios/p0'):
            (self.lab/d).mkdir(parents=True)
        (self.repo/'testing/node').mkdir(parents=True)
        for name in ('build-node.sh', 'stamp-node-image.sh'):
            shutil.copy2(ROOT/'scripts'/name, self.lab/'scripts'/name)
        for name in ('lab-excludes.txt', 'ukama-excludes.txt', 'exclusive.txt', 'worker.sh'):
            shutil.copy2(ROOT/'utils/runner'/name, self.lab/'utils/runner'/name)
        binary=self.lab/'bin/ukama-lab'; binary.write_text('#!/bin/sh\nexit 0\n');binary.chmod(0o755)
        self.scenario=self.lab/'scenarios/p0/test.yaml'
        self.scenario.write_text('version: 1\nstatus: active\nworld:\n  networks: 1\n  sites_per_network: 1\n')
        (self.repo/'testing/node/source.c').write_text('original node source')
        self.commands = self.root/'commands'; self.commands.mkdir()
        podman=self.commands/'podman';podman.write_text(PODMAN);podman.chmod(0o755)
        builder=self.repo/'testing/node/mk_local_vnode.sh';builder.write_text(BUILDER);builder.chmod(0o755)
        self.state=self.root/'podman.json';self.trace=self.root/'trace';self.build_trace=self.root/'builds'
        self.db={}
        for kind in images.NODE_KINDS:
            self.db[self.base(kind)]={'Id':'sha256:'+hashlib.sha256(kind.encode()).hexdigest(),'Architecture':'amd64','Os':'linux'}
        self.write_db()
        env={'PATH':str(self.commands)+':'+os.environ['PATH'], 'UKAMA_REPO':str(self.repo),
             'PODMAN_STATE':str(self.state), 'PODMAN_TRACE':str(self.trace), 'BUILD_TRACE':str(self.build_trace),
             'P0_IMAGE_CACHE_DIR':str(self.root/'image-cache'), 'NODE_RUNTIME':'starter',
             'BASE_IMAGE_REPO':'localhost/testing/virtualnode-base', 'ULAB_REQUIRE_NODE_BASE':'0', 'ULAB_DISTRIBUTED_RUNNER':'0',
             'UE_IMAGE':'ukama/ue:dev', 'MEDIA_IMAGE':'ukama/media:dev',
             'ULAB_NET_PROBE_IMAGE':'docker.io/library/alpine:3.20'}
        self.addCleanup(patch.stopall)
        patch.dict(os.environ,env).start()
        patch.object(images,'ROOT',self.lab).start()
        patch.object(images.platform,'machine',return_value='x86_64').start()
        self.factory=patch.object(protocol,'factory_request',side_effect=AssertionError('unexpected Factory lookup')).start()
        patch.object(images.Store,'aws',side_effect=AssertionError('unexpected AWS/ECR call')).start()
        self.template={k:'a-'+k+'-1' for k in images.NODE_KINDS}

    def base(self,kind): return 'localhost/testing/virtualnode-base:'+kind+'-starter'
    def write_db(self): self.state.write_text(json.dumps(self.db))
    def calls(self): return [json.loads(s) for s in self.trace.read_text().splitlines()] if self.trace.exists() else []
    def prepare(self,name='input',**kw): return images.prepare(self.root/name,**kw)

    def test_existing_three_bases_exported_without_any_build_or_factory_lookup(self):
        manifest=protocol.read_json(self.prepare())
        self.assertEqual(set(manifest['images']),set(images.NODE_KINDS))
        self.assertEqual({k:v['id'] for k,v in manifest['images'].items()},
                         {k:self.db[self.base(k)]['Id'] for k in images.NODE_KINDS})
        self.assertEqual(len([c for c in self.calls() if c[0]=='save']),3)
        self.assertFalse(self.build_trace.exists());self.factory.assert_not_called()
        self.assertFalse(any(c[0] in ('build','pull','push','login') for c in self.calls()))

    def test_scenario_edits_change_bundle_but_reuse_existing_image_exports(self):
        a=protocol.read_json(self.prepare('first'))
        self.scenario.write_text(self.scenario.read_text()+'description: changed scenario only\n')
        b=protocol.read_json(self.prepare('second'))
        self.assertNotEqual(a['fingerprint'],b['fingerprint'])
        self.assertNotEqual(a['sources']['ukama-lab.tar.gz'],b['sources']['ukama-lab.tar.gz'])
        self.assertEqual(a['images'],b['images'])
        self.assertEqual(len([c for c in self.calls() if c[0]=='save']),3)
        self.assertFalse(self.build_trace.exists())

    def test_missing_one_base_uses_existing_builder_only_for_that_kind(self):
        del self.db[self.base('cnode')];self.write_db()
        manifest=protocol.read_json(self.prepare(template_ids=self.template))
        self.assertEqual(self.build_trace.read_text().splitlines(),['a-cnode-1'])
        self.assertEqual(manifest['images']['anode']['id'],self.db[self.base('anode')]['Id'])
        self.factory.assert_not_called()

    def test_missing_base_can_seed_existing_node_image_without_compiling(self):
        self.db['localhost/testing/virtualnode:a-cnode-1']=self.db.pop(self.base('cnode'));self.write_db()
        self.prepare(template_ids=self.template)
        self.assertFalse(self.build_trace.exists())
        self.assertIn(self.base('cnode'),json.loads(self.state.read_text()))

    def test_optional_runtime_images_are_copied_only_when_already_cached(self):
        self.db['localhost/ukama/ue:dev']=self.db[self.base('tnode')];self.write_db()
        manifest=protocol.read_json(self.prepare())
        self.assertIn('ue',manifest['images']);self.assertNotIn('media',manifest['images'])
        self.assertEqual(len([c for c in self.calls() if c[0]=='save']),3)
        self.assertEqual(images.env_for(manifest)['ULAB_REQUIRE_NODE_BASE'],'1')

    def transfer(self,manifest_path):
        worker=self.root/'worker';worker.mkdir()
        target=worker/'images.json';shutil.copy2(manifest_path,target)
        state=worker/'podman.json';state.write_text('{}')
        parent=manifest_path.parent
        class S3:
            def __init__(self): self.downloads=[]
            def download(self,key,dest):
                self.downloads.append(key);Path(dest).parent.mkdir(parents=True,exist_ok=True)
                shutil.copy2(parent/key.removeprefix('input/'),dest)
        return target,state,S3()

    def test_transfer_then_existing_scripts_stamp_two_scenarios_on_worker(self):
        path=self.prepare();target,state,s3=self.transfer(path)
        with patch.dict(os.environ,{'PODMAN_STATE':str(state)}):
            manifest=images.load_images(target,s3)
            env={**os.environ,**images.env_for(manifest)}
            for sequence in (1,2):
                for kind in images.NODE_KINDS:
                    node=f'scenario{sequence}-{kind}-node'
                    result=subprocess.run(['sh',str(ROOT/'scripts/build-node.sh'),str(self.repo),node,'starter'],
                                          env=env,capture_output=True,text=True)
                    self.assertEqual(result.returncode,0,result.stderr)
            images.load_images(target,s3)  # Already-cached content never downloads again.
        loaded=json.loads(state.read_text())
        for kind in images.NODE_KINDS:
            self.assertEqual(loaded[self.base(kind)]['Id'],self.db[self.base(kind)]['Id'])
            for sequence in (1,2):
                self.assertIn(f'localhost/testing/virtualnode:scenario{sequence}-{kind}-node',loaded)
        self.assertEqual(len(s3.downloads),3)
        self.assertEqual(len([c for c in self.calls() if c[0]=='commit']),6)
        self.assertFalse(self.build_trace.exists())

    def test_corrupt_archive_stops_before_load(self):
        path=self.prepare();manifest=protocol.read_json(path)
        (path.parent/next(iter(manifest['images'].values()))['archive']).write_bytes(b'corrupted')
        target,state,s3=self.transfer(path)
        with patch.dict(os.environ,{'PODMAN_STATE':str(state)}),self.assertRaisesRegex(ValueError,'checksum mismatch'):
            images.load_images(target,s3)
        self.assertFalse(any(c[0]=='load' for c in self.calls()))

    def test_wrong_architecture_stops_before_download(self):
        path=self.prepare();target,state,s3=self.transfer(path)
        manifest=protocol.read_json(target);manifest['architecture']='aarch64';protocol.atomic_json(target,manifest)
        with patch.dict(os.environ,{'PODMAN_STATE':str(state)}),self.assertRaisesRegex(ValueError,'architectures differ'):
            images.load_images(target,s3)
        self.assertEqual(s3.downloads,[])

    def test_missing_required_base_in_manifest_rejected(self):
        manifest=protocol.read_json(self.prepare());del manifest['images']['anode']
        with self.assertRaisesRegex(ValueError,'three node bases'):
            images.validate(manifest)

    def test_dry_run_ignores_stale_manifest_and_never_prepares_images(self):
        with patch.object(controller,'ROOT',self.lab),patch.dict(os.environ,{'P0_ARTIFACT_MANIFEST':'/obsolete/not-found.json',
                'LOCAL_RUNS_DIR':str(self.root/'runs')}),patch.object(controller,'prepare_images') as prep,\
                patch.object(controller,'snapshot_sources') as snapshot,\
                patch.object(sys,'argv',['controller.py','p0','--workers','10','--dry-run','--batch-id','dry-test']):
            controller.main()
        prep.assert_not_called();snapshot.assert_not_called();self.assertEqual(self.calls(),[])
        plan=protocol.read_json(self.root/'runs/dry-test/plan.json')
        self.assertEqual(len(plan['jobs']),1)

    def test_normal_run_automatically_uploads_bases_before_launch(self):
        events=[]
        class S3:
            def __init__(self,bucket,prefix): self.bucket=bucket;self.prefix=prefix
            @staticmethod
            def aws(*args,**kw):
                if args[:2]==('s3api','put-object'): return b'{"IfNoneMatch": ""}'
                if args[:2]==('secretsmanager','get-secret-value'):
                    return json.dumps({'VersionId':'pinned-test', 'SecretString':json.dumps({'UKAMA_IDENTIFIER':'test','UKAMA_PASSWORD':'synthetic','ULAB_VPN_CONFIG':'synthetic'})}).encode()
                raise AssertionError('unexpected cloud operation: '+str(args))
            def put(self,key,value,create=False): events.append(('put',key,json.loads(json.dumps(value))))
            def upload(self,source,key):
                if not Path(source).is_file(): raise AssertionError('missing input: '+str(source))
                events.append(('upload',key))
        def launch(c,worker,env):
            published=next(e[2] for e in events if e[:2]==('put','input/images.json'))
            for image in published['images'].values():
                self.assertIn(('upload','input/'+image['archive']),events)
            self.assertEqual(c.plan['artifact'],published['fingerprint'])
            events.append(('launch',worker))
        def finish(c,reason):
            self.assertEqual(reason,'')
            c.state['results']={j['id']:{'outcome':'PASS'} for j in c.plan['jobs']}
            c.fd.close()
        config={k:'synthetic' for k in ('AWS_REGION','S3_BUCKET','S3_PREFIX','AMI_ID','INSTANCE_TYPE',
                'INSTANCE_PROFILE_NAME','SUBNET_ID','SECURITY_GROUP_ID','SECRET_ID')}
        config.update(LOCAL_RUNS_DIR=str(self.root/'runs'),FACTORY_NODE_TARGET='0',
                      P0_ARTIFACT_MANIFEST='/obsolete/missing.json')
        payload={'nodes':[{'nodeId':nid} for nid in self.template.values()]}
        with patch.object(controller,'ROOT',self.lab),patch.dict(os.environ,config),\
                patch.object(controller,'Store',S3),patch.object(protocol,'factory_request',return_value=payload),\
                patch.object(controller.Controller,'acquire'),patch.object(controller.Controller,'launch',launch),\
                patch.object(controller.Controller,'watch',return_value=''),patch.object(controller.Controller,'finish',finish),\
                patch.object(sys,'argv',['controller.py','p0','--workers','10','--batch-id','automatic']):
            controller.main()
        self.assertIn(('launch','worker-01'),events)
        self.assertFalse(self.build_trace.exists())
        self.factory.assert_not_called()


if __name__=='__main__': unittest.main()
