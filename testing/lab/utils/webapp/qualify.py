#!/usr/bin/env python3
"""Read-only full-suite readiness and conservative final qualification. MPL-2.0."""
import argparse
from collections import Counter
import json
import os
from pathlib import Path
import re
import shutil
import socket
import subprocess
import sys
from urllib.parse import urlsplit

import matrix as m

VARIABLE = re.compile(r'\$\{([A-Za-z_][A-Za-z_0-9]*)\}')
SCRIPTS = ('start-media.sh','wait-media-ready.sh','start-ue.sh','wait-ues-attached.sh',
           'traffic.sh','cleanup-ue.sh','ensure-network.sh','build-and-start-site.sh',
           'wait-nodes-ready.sh','disconnect-node.sh','reconnect-node.sh','stop-node.sh',
           'stop-media.sh','cleanup-network.sh')

def expand(value):
    return VARIABLE.sub(lambda match: os.environ[match[1]], str(value))

def inventory(root):
    policy = json.loads((root/m.POLICY).read_text())
    if policy.get('schema_version')!=1 or policy.get('requirements')!=132 or policy.get('p0')!=86:
        raise ValueError('qualification requires the frozen 132 / 86 policy')
    catalog = json.loads((root/m.CATALOG).read_text())
    m.validate_inventory(root,catalog)
    # Validate every catalog mapping, including partial/unimplemented entries.
    m.coverage(root,[],list(m.BROWSERS),'','')
    return catalog

def input_inventory(root,catalog):
    rows=[]
    for relative in sorted({p for r in catalog['requirements'] for p in r['scenarios']}):
        path=(root/relative).resolve()
        if not path.is_relative_to(root/'scenarios/webapp'):
            raise ValueError('scenario outside scenarios/webapp')
        text=path.read_text();case=m.load_case(path)
        names=sorted(set(VARIABLE.findall(text)))
        missing=[n for n in names if not os.environ.get(n)]
        rows.append({'scenario':relative,'sha256':m.digest(path),'status':case.get('status'),
                     'controlled':m.controlled(case),'environment':names,'missing_environment':missing,
                     'case':case})
    return rows

def auth_state(root,path,allow_none=False):
    if path=='none':return 'not_required' if allow_none else 'invalid'
    try:
        p=(root/path).resolve()
        if not p.is_file() or p.stat().st_size>4*1024*1024:return 'missing_or_invalid'
        state=json.loads(p.read_text())
        if not isinstance(state,dict) or not isinstance(state.get('cookies'),list) or not isinstance(state.get('origins'),list):return 'invalid'
        return 'present_unverified'
    except (OSError,ValueError,TypeError):return 'missing_or_invalid'

def endpoint(value,probe=False):
    try:
        url=urlsplit(value)
        if url.scheme not in ('http','https') or not url.hostname or url.username or url.password or url.query or url.fragment:
            return 'invalid_url'
        port=url.port or (443 if url.scheme=='https' else 80)
        if not probe:return 'not_probed'
        with socket.create_connection((url.hostname,port),timeout=1):pass
        return 'tcp_open'
    except (ValueError,OSError):return 'unreachable'

def probe_browsers(root,browsers):
    try:
        p=subprocess.run(['node',str(root/'adapters/webapp/preflight.mjs'),*browsers],cwd=root,
                         capture_output=True,text=True,timeout=45)
        rows=json.loads(p.stdout)['browsers']
        if not isinstance(rows,list) or [r['browser'] for r in rows]!=browsers:raise ValueError()
        if any(r['state'] not in ('ready','unavailable','ambiguous_override') for r in rows):raise ValueError()
        return rows
    except (OSError,ValueError,KeyError,TypeError,subprocess.TimeoutExpired):
        return [{'browser':b,'state':'probe_failed'} for b in browsers]

def preflight(args):
    root=args.root.resolve();catalog=inventory(root);rows=input_inventory(root,catalog)
    checks=[]
    def check(name,ready,detail):checks.append({'check':name,'ready':bool(ready),'detail':detail})
    binary=args.binary.resolve()
    check('lab_binary',binary.is_file() and os.access(binary,os.X_OK),'Build the current C sources; existence is not build provenance.')
    sources=list((root/'adapters/webapp/src').glob('*.ts'))
    compiled=all((root/'adapters/webapp/dist'/p.with_suffix('.js').name).is_file() and
                 (root/'adapters/webapp/dist'/p.with_suffix('.js').name).stat().st_mtime_ns>=p.stat().st_mtime_ns for p in sources)
    check('worker_build',bool(sources) and compiled,'Build the worker after applying both patches; timestamps are only a freshness hint.')
    check('build_identifiers',bool(args.app_build and args.backend_build),'Specify the actually deployed console and backend build identifiers.')
    check('runtime_repo',args.repo is not None and args.repo.is_dir(),'Provide the local ukama runtime repository.')
    missing_scripts=[name for name in SCRIPTS if not args.scripts or not (args.scripts/name).is_file() or not os.access(args.scripts/name,os.X_OK)]
    check('runtime_scripts',not missing_scripts,{'missing':missing_scripts})
    for command in ('node','pdftotext','pdftoppm'):
        check(command,shutil.which(command) is not None,'Required by the worker or PDF receipt qualification.')
    browser_rows=probe_browsers(root,args.browser)
    check('browser_runtimes',all(r['state']=='ready' for r in browser_rows),browser_rows)
    missing=sorted({name for row in rows for name in row['missing_environment']})
    check('scenario_environment',not missing,{'missing_names':missing})
    needs_display=any(not row['case']['webapp'].get('headless',False) for row in rows)
    check('headed_display',not needs_display or sys.platform!='linux' or bool(os.environ.get('DISPLAY') or os.environ.get('WAYLAND_DISPLAY')),'Headed Linux profiles need a working display; the blank-page probe is headless.')
    endpoints={};auth=[]
    for row in rows:
        case=row.pop('case');profile=case['webapp']
        row['validation']='blocked'
        if not row['missing_environment']:
            state=auth_state(root,expand(profile['auth_state']),profile.get('session_mode')=='auth_test')
            auth.append({'scenario':row['scenario'],'scope':'primary','state':state})
            for phase in case.get('phases',[]):
                for event in phase.get('events',[]):
                    if 'auth_state' in event:
                        auth.append({'scenario':row['scenario'],'scope':'peer','state':auth_state(root,expand(event['auth_state']))})
            for key in ('base_url','auth_origin'):
                if key in profile:
                    value=expand(profile[key]);endpoints[value]=endpoint(value,args.probe_endpoints)
            # validate is the C parser's offline path; never run a scenario here.
            if binary.is_file() and os.access(binary,os.X_OK) and args.repo and args.repo.is_dir():
                try:
                    p=subprocess.run([str(binary),'validate',str(root/row['scenario']),'--repo',str(args.repo.resolve())],
                                     cwd=root,capture_output=True,timeout=15)
                    row['validation']='valid' if p.returncode==0 else 'invalid'
                except (OSError,subprocess.TimeoutExpired):row['validation']='unavailable'
        else:auth.append({'scenario':row['scenario'],'state':'blocked_by_environment'})
    check('auth_state_files',all(r['state'] in ('present_unverified','not_required') for r in auth),auth)
    check('offline_scenarios',all(r['validation']=='valid' and r['status']=='active' for r in rows),
          {'states':dict(Counter(r['validation'] for r in rows)),'inactive':[r['scenario'] for r in rows if r['status']!='active']})
    for name in ('warehouse_url','factory_url','asr_url','bff'):
        value=getattr(args,name)
        check(name,bool(value) and endpoint(value,args.probe_endpoints)=='tcp_open',
              endpoint(value,args.probe_endpoints) if value else 'not_configured')
    # Report scenario names and states, never endpoint credentials or env values.
    check('app_auth_endpoints',bool(endpoints) and all(v=='tcp_open' for v in endpoints.values()),
          {'endpoint_states':dict(Counter(endpoints.values())),'note':'TCP reachability does not validate HTTP, authentication or service readiness.'})
    matrix_out=args.out.resolve()/'live-matrix'
    command=[sys.executable,str(root/'utils/webapp/matrix.py'),'run','--root',str(root),'--binary',str(binary),
             '--out',str(matrix_out),'--evidence','live','--app-build',args.app_build,'--backend-build',args.backend_build]
    for browser in args.browser:command+=['--browser',browser]
    forwarded=[]
    for name in ('repo','scripts','warehouse_url','factory_url','asr_url','bff'):
        value=getattr(args,name)
        if value is not None:
            if name in ('repo','scripts'):value=value.resolve()
            # Never write supplied credential-bearing URLs into the plan.
            if name not in ('repo','scripts') and endpoint(value)=='invalid_url':continue
            forwarded+=['--'+name.replace('_','-'),str(value)]
    command+=['--',*forwarded]
    report={'schema_version':1,'execution_ready':all(c['ready'] for c in checks),
            'live_verified':0,'scope':'read-only preflight; no scenario executed; auth contents are not verified',
            'source_sha256':m.fingerprint(root),'browsers':args.browser,'requirements':132,'p0':86,
            'scenario_count':len(rows),'attempt_count':len(rows)*len(args.browser),'checks':checks,'scenarios':rows,
            'matrix_command':command,'qualification_gate_possible_with_current_automation':all(r['automation']=='implemented' for r in catalog['requirements'] if r['priority']=='p0') and sum(r['automation']=='implemented' for r in catalog['requirements'])/132>=.9,
            'role_policy':'Unagreed execution policy: WEB-SHELL-010 remains partial.'}
    args.out.mkdir(mode=0o700,parents=True,exist_ok=False)
    m.write_json(args.out/'preflight.json',report)
    names=sorted({name for row in rows for name in row['environment']})
    (args.out/'inputs.env.example').write_text('# Supply independent host values. No current values are copied.\n'+''.join(f'# {name}=\n' for name in names))
    print(json.dumps({'execution_ready':report['execution_ready'],'scenarios':len(rows),'attempts':report['attempt_count'],
                      'blocked_checks':[c['check'] for c in checks if not c['ready']]}))
    return int(not report['execution_ready'])

def closeout(args):
    root=args.root.resolve();inventory(root)
    result=m.coverage(root,args.manifest,args.browser,args.app_build,args.backend_build)
    args.out.mkdir(mode=0o700,parents=True,exist_ok=False)
    m.write_json(args.out/'coverage.json',result)
    (args.out/'coverage.html').write_text(m.render_html(result))
    gaps=[{'id':r['id'],'priority':r['priority'],'automation':r['automation'],
           'note':r.get('implementation_note','No complete automation.'),'scenarios':r['scenarios']}
          for r in result['requirements'] if r['automation']!='implemented']
    states=Counter(e['state'] for r in result['requirements'] for e in r['evidence'])
    summary={'schema_version':1,'gate_passed':result['gate_passed'],'totals':result['totals'],'p0':result['p0'],
             'browsers':args.browser,'manifest_count':len(set(p.resolve() for p in args.manifest)),
             'source_sha256':result['source_sha256'],'unimplemented_or_partial':gaps,
             'unverified_requirements':[r['id'] for r in result['requirements'] if not r['verified']],
             'evidence_states':dict(states),'note':'Fixture success is not live acceptance. Include every attempt manifest for the release cohort.'}
    m.write_json(args.out/'closeout.json',summary)
    print(json.dumps({k:summary[k] for k in ('gate_passed','totals','p0','browsers')}))
    return int(not result['gate_passed'])

def main():
    os.umask(0o077)
    parser=argparse.ArgumentParser(description=__doc__);sub=parser.add_subparsers(dest='command',required=True)
    for name in ('preflight','closeout'):
        p=sub.add_parser(name);p.add_argument('--root',type=Path,default=m.ROOT);p.add_argument('--out',type=Path,required=True)
        p.add_argument('--browser',choices=m.BROWSERS,action='append');p.add_argument('--app-build',default='');p.add_argument('--backend-build',default='')
        if name=='preflight':
            p.add_argument('--binary',type=Path,default=m.ROOT/'bin/ukama-lab')
            p.add_argument('--repo',type=Path);p.add_argument('--scripts',type=Path)
            for flag in ('warehouse-url','factory-url','asr-url','bff'):p.add_argument('--'+flag)
            p.add_argument('--probe-endpoints',action='store_true',help='read-only TCP checks, no HTTP requests')
        else:p.add_argument('--manifest',type=Path,action='append',default=[])
    args=parser.parse_args();args.browser=list(dict.fromkeys(args.browser or m.BROWSERS))
    try:return preflight(args) if args.command=='preflight' else closeout(args)
    except (OSError,ValueError,TypeError,KeyError,m.yaml.YAMLError) as error:
        parser.error('Qualification could not complete: '+str(error))

if __name__=='__main__':raise SystemExit(main())
