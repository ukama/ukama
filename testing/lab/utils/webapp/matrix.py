#!/usr/bin/env python3
"""Serial browser matrix and conservative requirement coverage. SPDX-License-Identifier: MPL-2.0"""
import argparse
from collections import Counter
import hashlib
import html
import json
import os
from pathlib import Path
import re
import signal
import subprocess
import time
import uuid
import yaml

ROOT = Path(__file__).resolve().parents[2]
CATALOG = 'docs/webapp/coverage.json'
BROWSERS = ('chromium', 'firefox', 'webkit')

class UniqueLoader(yaml.SafeLoader):
    pass

def mapping(loader, node, deep=False):
    result = {}
    for key, value in node.value:
        key = loader.construct_object(key, deep=deep)
        if key in result:
            raise ValueError(f'duplicate YAML key: {key}')
        result[key] = loader.construct_object(value, deep=deep)
    return result
UniqueLoader.add_constructor(yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, mapping)

def digest(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()

def fingerprint(root):
    """Source cohort; generated output, credentials and reports never enter this hash."""
    result = hashlib.sha256()
    files = [root / CATALOG, root / 'utils/webapp-worker.sh', root / 'utils/webapp/matrix.py']
    for part in ('src', 'inc', 'cmd', 'scenarios/webapp', 'adapters/webapp/src'):
        files += [p for p in (root / part).rglob('*') if p.is_file()]
    files += [root / 'adapters/webapp/package-lock.json']
    for p in sorted(files):
        result.update(str(p.relative_to(root)).encode() + b'\0' + p.read_bytes() + b'\0')
    return result.hexdigest()

def load_case(path):
    case = yaml.load(Path(path).read_text(), Loader=UniqueLoader)
    if not isinstance(case, dict) or case.get('version') != 2 or case.get('suite') != 'webapp':
        raise ValueError(f'not a v2 webapp scenario: {path}')
    return case

def controlled(case):
    """Injected response faults cannot be promoted to live product evidence."""
    return any((e.get('type') == 'web_action' and e.get('action') == 'status_fault') or (e.get('type') == 'web_onboard' and e.get('action') == 'arm_fault') or (e.get('type') == 'web_inventory' and e.get('action') in ('mask_home', 'stale_selection'))
               for phase in case.get('phases', []) for e in phase.get('events', []))

def planned(case):
    checks, events = [], []
    phases = case.get('phases', []) + [{'name': 'final', 'checks': case.get('final_checks', [])}]
    for p in phases:
        for c in p.get('checks', []):
            checks.append((p['name'], c['type'], c['requirement'], c['label']))
        for e in p.get('events', []):
            events.append((p['name'], e['type']))
    return Counter(checks), Counter(events)

def write_json(path, obj):
    tmp = path.with_suffix(path.suffix + '.tmp')
    tmp.write_text(json.dumps(obj, indent=2) + '\n')
    tmp.chmod(0o600)
    tmp.replace(path)

def generated_case(text, browser):
    generated, count = re.subn(r'^  browser:.*$', f'  browser: {browser}', text, flags=re.M)
    if count != 1:
        raise ValueError('scenario requires one explicit browser profile')
    return generated

def environment_hash(text):
    names = sorted(set(re.findall(r'\$\{([A-Za-z_][A-Za-z_0-9]*)\}', text)))
    return hashlib.sha256(json.dumps({name: os.environ.get(name) for name in names}, sort_keys=True).encode()).hexdigest()

def run_matrix(args):
    root = args.root.resolve()
    inventory = json.loads((root / CATALOG).read_text())
    paths = args.scenario or sorted({p for r in inventory['requirements'] for p in r['scenarios']})
    paths = sorted(set(paths))
    if not paths:
        raise ValueError('no scenarios selected')
    cases = {}
    for path in paths:
        full = (root / path).resolve()
        if not full.is_relative_to(root / 'scenarios/webapp'):
            raise ValueError('scenarios must be in scenarios/webapp')
        case = load_case(full)
        cases[str(full.relative_to(root))] = case
    if args.evidence == 'live' and (not args.app_build or not args.backend_build):
        raise ValueError('live evidence requires --app-build and --backend-build')
    binary = args.binary.resolve()
    if not binary.is_file():
        raise ValueError('build the lab binary before running the matrix')
    extra = args.lab_args[1:] if args.lab_args[:1] == ['--'] else args.lab_args
    allowed = {'--repo', '--scripts', '--warehouse-url', '--factory-url', '--asr-url', '--bff'}
    if len(extra) % 2 or any(extra[i] not in allowed or extra[i+1].startswith('--') for i in range(0, len(extra), 2)):
        raise ValueError('forward only endpoint/repo/scripts pairs after --; matrix owns run/output/worker flags')
    out = args.out.resolve()
    out.mkdir(mode=0o700, parents=True, exist_ok=False)
    manifest = {'schema_version': 1, 'matrix_id': str(uuid.uuid4()), 'started_ns': time.time_ns(),
                'evidence': args.evidence, 'app_build': args.app_build, 'backend_build': args.backend_build,
                'source_sha256': fingerprint(root), 'inventory_sha256': digest(root / CATALOG),
                'binary_sha256': digest(binary), 'attempts': []}
    for path, case in cases.items():
        for browser in args.browser:
            index = len(manifest['attempts']) + 1
            run_id = f'matrix-{index:04d}-{uuid.uuid4().hex[:8]}'
            attempt = {'scenario': path, 'browser': browser, 'run_id': run_id,
                       'scenario_sha256': digest(root / path), 'environment_sha256': environment_hash((root / path).read_text()), 'started_ns': time.time_ns(),
                       'state': 'pending', 'report': f'runs/{run_id}/report.json'}
            manifest['attempts'].append(attempt)
    write_json(out / 'matrix.json', manifest)
    failed = False
    for attempt in manifest['attempts']:
        case = cases[attempt['scenario']]
        if case.get('status') != 'active':
            attempt.update(state='skipped', reason='scenario is not active')
            failed = True
            write_json(out / 'matrix.json', manifest)
            continue
        original = (root / attempt['scenario']).read_text()
        original = generated_case(original, attempt['browser'])
        source = out / (attempt['run_id'] + '.yaml')
        source.write_text(original)
        source.chmod(0o600)
        attempt.update(state='running', started_ns=time.time_ns(), generated_case=source.name, generated_sha256=digest(source))
        write_json(out / 'matrix.json', manifest)
        command = [str(binary), 'run', str(source), '--out', str(out / 'runs'), '--run-id', attempt['run_id'],
                   '--webapp-worker', str(root / 'utils/webapp-worker.sh'), *extra]
        interrupted = False
        try:
            with (out / (attempt['run_id'] + '.log')).open('w') as log:
                child = subprocess.Popen(command, cwd=root, stdout=log, stderr=subprocess.STDOUT, start_new_session=True)
                try:
                    code = child.wait(timeout=int(case['webapp']['scenario_timeout_seconds']) + 120)
                except (subprocess.TimeoutExpired, KeyboardInterrupt) as exc:
                    # Give the lab its normal cleanup path before escalating this owned process group.
                    child.send_signal(signal.SIGTERM)
                    try:
                        code = child.wait(timeout=30)
                    except subprocess.TimeoutExpired:
                        os.killpg(child.pid, signal.SIGKILL)
                        code = child.wait()
                    interrupted = isinstance(exc, KeyboardInterrupt)
                    attempt['reason'] = 'interrupted' if interrupted else 'matrix deadline exceeded'
                    code = code or 124
            attempt.update(state='complete', exit_code=code, ended_ns=time.time_ns())
            report = out / attempt['report']
            if report.is_file():
                attempt['report_sha256'] = digest(report)
            failed |= code != 0
        except OSError as exc:
            attempt.update(state='launch_failed', reason=type(exc).__name__, ended_ns=time.time_ns())
            failed = True
        write_json(out / 'matrix.json', manifest)
        print(f"{attempt['browser']} {attempt['scenario']}: {attempt['state']} rc={attempt.get('exit_code', '-')} ", flush=True)
        if interrupted:
            return 130
    return int(failed)

def evaluate(root, manifest_path, manifest, attempt, source_hash, catalog_hash):
    """No credit from a summary flag alone: require every planned check and event."""
    try:
        if manifest.get('schema_version') != 1:
            return 'invalid_manifest'
        path = (root / attempt['scenario']).resolve()
        if not path.is_relative_to(root / 'scenarios/webapp'):
            return 'invalid_path'
        if manifest.get('source_sha256') != source_hash or manifest.get('inventory_sha256') != catalog_hash or attempt.get('scenario_sha256') != digest(path):
            return 'stale_source'
        if attempt.get('state') != 'complete' or attempt.get('exit_code') != 0:
            return attempt.get('state', 'incomplete') if attempt.get('state') != 'complete' else 'failed'
        if attempt.get('environment_sha256') != environment_hash(path.read_text()):
            return 'stale_environment'
        generated = (manifest_path.parent / attempt['generated_case']).resolve()
        if not generated.is_relative_to(manifest_path.parent.resolve()) or digest(generated) != attempt.get('generated_sha256') or generated.read_text() != generated_case(path.read_text(), attempt['browser']):
            return 'modified_case'
        report_path = (manifest_path.parent / attempt['report']).resolve()
        if not report_path.is_relative_to(manifest_path.parent.resolve()):
            return 'invalid_path'
        if digest(report_path) != attempt.get('report_sha256'):
            return 'report_hash_mismatch'
        report = json.loads(report_path.read_text())
        case = load_case(path)
        if case.get('status') != 'active' or report.get('status') != 'active':
            return 'skipped'
        if report.get('run_id') != attempt['run_id'] or report.get('scenario') != case['name'] or report.get('suite') != 'webapp':
            return 'wrong_report'
        if report.get('browser') != attempt['browser']:
            return 'wrong_browser'
        if report.get('outcome') != 'PASS' or report.get('passed') is not True or report.get('final_rc') != 0 or report.get('cleanup') != 'ok':
            return 'failed'
        want_checks, want_events = planned(case)
        checks = [r for r in report['results'] if r['kind'] == 'check']
        events = [r for r in report['results'] if r['kind'] == 'event']
        actual_checks = Counter((r['phase'], r['name'], r['requirement'], r['label']) for r in checks)
        actual_events = Counter((r['phase'], r['name']) for r in events)
        if not want_checks or want_checks != actual_checks or want_events != actual_events:
            return 'incomplete_assertions'
        if any(r['state'] != 'PASS' or r.get('actual') is None or r.get('expected') is None for r in checks) or any(r['state'] != 'PASS' for r in events):
            return 'failed'
        for r in checks:
            expected, actual = r['expected'], r['actual']
            if isinstance(expected, str) and isinstance(actual, str):
                expected, actual = ' '.join(expected.split()), ' '.join(actual.split())
            if r.get('match', 'equals') == 'contains':
                if not isinstance(expected, str) or not isinstance(actual, str) or not expected or expected not in actual:
                    return 'inconsistent_values'
            elif r.get('match', 'equals') != 'equals' or type(expected) is not type(actual) or expected != actual:
                return 'inconsistent_values'
        for name, rows in [('checks', checks), ('events', events)]:
            if report.get(name) != {'total': len(rows), 'passed': len(rows), 'failed': 0}:
                return 'inconsistent_totals'
        return 'passed'
    except (OSError, ValueError, KeyError, TypeError):
        return 'invalid_evidence'

def coverage(root, manifest_paths, browsers, app_build, backend_build):
    inventory = json.loads((root / CATALOG).read_text())
    requirements = inventory['requirements']
    ids = {r['id'] for r in requirements}
    if len(ids) != len(requirements) or not requirements:
        raise ValueError('empty or duplicate requirement inventory')
    source_hash, catalog_hash = fingerprint(root), digest(root / CATALOG)
    all_attempts = {}
    for path in sorted(set(p.resolve() for p in manifest_paths)):
        manifest = json.loads(path.read_text())
        if manifest.get('app_build') != app_build or manifest.get('backend_build') != backend_build:
            continue
        for a in manifest.get('attempts', []):
            if a.get('browser') not in browsers:
                continue
            key = (a.get('scenario'), a['browser'])
            state = evaluate(root, path, manifest, a, source_hash, catalog_hash)
            all_attempts.setdefault(key, []).append({'started_ns': a.get('started_ns', 0), 'state': state,
                'evidence': manifest.get('evidence', 'unverified'), 'run_id': a.get('run_id'), 'manifest': str(path)})
    rows = []
    for r in requirements:
        cells, verified = [], r['automation'] == 'implemented' and bool(r['scenarios'])
        for scenario in r['scenarios']:
            case = load_case(root / scenario)
            checks, _ = planned(case)
            # Mappings must reference real assertions, not merely a passing scenario.
            if r['id'] not in {k[2] for k in checks}:
                raise ValueError(f"catalog mapping has no assertion: {r['id']} {scenario}")
            if any(k[2] not in ids for k in checks):
                raise ValueError(f'unknown requirement in {scenario}')
            for browser in browsers:
                attempts = sorted(all_attempts.get((scenario, browser), []), key=lambda a: (a['started_ns'], a['manifest'], a['run_id'] or ''))
                latest = attempts[-1] if attempts else {'state': 'not_run', 'evidence': 'none'}
                live = [a['state'] for a in attempts if a['evidence'] == 'live' and a['state'] != 'stale_source']
                flaky = 'passed' in live and any(s != 'passed' for s in live)
                if controlled(case) and latest['evidence'] == 'live':
                    latest = {**latest, 'evidence': 'controlled_ui'}
                credit = latest['state'] == 'passed' and latest['evidence'] == 'live' and not flaky and bool(app_build and backend_build)
                verified &= credit
                cells.append({'scenario': scenario, 'browser': browser, **latest, 'flaky': flaky})
        rows.append({**r, 'verified': bool(verified), 'evidence': cells})
    def tally(selected):
        n = len(selected)
        return {'total': n, 'automated': sum(r['automation'] == 'implemented' for r in selected),
                'verified': sum(r['verified'] for r in selected),
                'verified_percent': round(100 * sum(r['verified'] for r in selected) / n, 2) if n else 0}
    totals, p0 = tally(rows), tally([r for r in rows if r['priority'] == 'p0'])
    return {'schema_version': 1, 'app_build': app_build, 'backend_build': backend_build, 'browsers': browsers,
            'source_sha256': source_hash, 'inventory_sha256': catalog_hash, 'totals': totals, 'p0': p0,
            'gate_passed': bool(p0['total'] and p0['verified'] == p0['total'] and totals['verified'] / totals['total'] >= .9),
            'requirements': rows}

def render_html(report):
    e = lambda x: html.escape(str(x), quote=True)
    rows = []
    for r in report['requirements']:
        states = sorted({a['state'] + '/' + a['evidence'] + (' (flaky)' if a['flaky'] else '') for a in r['evidence']})
        rows.append('<tr>' + ''.join('<td>' + e(v) + '</td>' for v in [r['id'], r['priority'], r['requirement'], r['automation'], 'yes' if r['verified'] else 'no', ', '.join(states) or 'not_run', r.get('implementation_note', '')]) + '</tr>')
    t = report['totals']
    return '<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>Ukama browser coverage</title><style>body{font:16px system-ui;margin:2rem}table{border-collapse:collapse;width:100%}th,td{text-align:left;vertical-align:top;padding:.6rem;border:1px solid #bbb}th{background:#eee}td{overflow-wrap:anywhere}</style><h1>Ukama browser coverage</h1><p>' + e(f"{t['automated']}/{t['total']} automated; {t['verified']}/{t['total']} live verified. Gate: {'PASS' if report['gate_passed'] else 'FAIL'}. Requires 100% P0 and 90% overall.") + '</p><p>' + e(f"App: {report['app_build'] or 'unspecified'}; backend: {report['backend_build'] or 'unspecified'}; browsers: {', '.join(report['browsers'])}") + '</p><p>Fixture, injected-response, skipped, incomplete, stale and flaky results earn no live credit. Build identifiers and evidence type are operator declarations; this report is not a signed attestation.</p><table><thead><tr><th>ID</th><th>Priority</th><th>Requirement</th><th>Automation</th><th>Verified</th><th>Evidence</th><th>Gaps / notes</th></tr></thead><tbody>' + ''.join(rows) + '</tbody></table></html>\n'

def main():
    os.umask(0o077)
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest='command', required=True)
    for name in ('run', 'report'):
        q = sub.add_parser(name)
        q.add_argument('--root', type=Path, default=ROOT)
        q.add_argument('--out', type=Path, required=True)
        q.add_argument('--browser', action='append', choices=BROWSERS)
        q.add_argument('--app-build', default='')
        q.add_argument('--backend-build', default='')
        if name == 'run':
            q.add_argument('--binary', type=Path, default=ROOT / 'bin/ukama-lab')
            q.add_argument('--scenario', action='append')
            q.add_argument('--evidence', choices=('unverified', 'fixture', 'live'), default='unverified')
            q.add_argument('lab_args', nargs=argparse.REMAINDER)
        else:
            q.add_argument('--manifest', type=Path, action='append', default=[])
            q.add_argument('--gate', action='store_true', help='exit 1 unless 100%% P0 / 90%% overall are live verified')
    args = p.parse_args()
    args.browser = list(dict.fromkeys(args.browser or BROWSERS))
    try:
        if args.command == 'run':
            return run_matrix(args)
        result = coverage(args.root.resolve(), args.manifest, args.browser, args.app_build, args.backend_build)
        args.out.mkdir(parents=True, mode=0o700, exist_ok=False)
        write_json(args.out / 'coverage.json', result)
        (args.out / 'coverage.html').write_text(render_html(result))
        print(json.dumps({k: result[k] for k in ('totals', 'p0', 'gate_passed')}))
        return int(args.gate and not result['gate_passed'])
    except (ValueError, OSError, KeyError, TypeError, yaml.YAMLError) as exc:
        p.error(str(exc))

if __name__ == '__main__':
    raise SystemExit(main())
