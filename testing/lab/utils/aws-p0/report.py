#!/usr/bin/env python3
"""Aggregate outcomes without treating missing workers or scenarios as passes."""
import collections
import csv
import json
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
shards = sorted((root / 'input/shards').glob('worker-*.txt'))
rows, infra = [], []
if not shards:
    infra.append(('controller', 'missing input/shards; cannot verify completeness'))
for shard in shards:
    worker = shard.stem
    directory = root / 'workers' / worker
    expected = [p[len('scenarios/p0/'):] for p in shard.read_text().splitlines() if p]
    found = {}
    reasons = []
    if (directory / (worker + '.failed')).exists():
        reasons.append('worker reported infrastructure failure')
    if not (directory / (worker + '.done')).exists():
        reasons.append('missing worker completion marker')
    reports = list(directory.glob('results/*/batch-report.json'))
    if len(reports) != 1:
        reasons.append('missing or ambiguous batch-report.json')
    else:
        try:
            data = json.loads(reports[0].read_text())
            for row in data['results']:
                scenario = row['scenario']
                if scenario not in expected or scenario in found:
                    reasons.append('duplicate or unassigned scenario: ' + scenario)
                    continue
                if row['outcome'] not in ('PASS', 'FAIL', 'SKIP'):
                    reasons.append('invalid outcome: ' + scenario)
                    continue
                found[scenario] = dict(row, worker=worker)
        except (OSError, ValueError, KeyError, TypeError) as exc:
            reasons.append('invalid batch report: ' + str(exc))
    for scenario in expected:
        row = found.get(scenario)
        if row is None:
            reasons.append('missing result: ' + scenario)
            row = dict(worker=worker, scenario=scenario, category=scenario.split('/')[0],
                       outcome='MISSING', duration_sec=0)
        # Remote report paths become usable paths in the collected output.
        for key in ('report', 'log'):
            value = row.get(key, '')
            if '/results/' in value:
                row[key] = str(directory / 'results' / value.split('/results/', 1)[1])
        rows.append(row)
    infra.extend((worker, reason) for reason in dict.fromkeys(reasons))
counts = collections.Counter(row['outcome'] for row in rows)
payload = dict(total=len(rows), completed=len(rows)-counts['MISSING'], passed=counts['PASS'],
               failed=counts['FAIL'], skipped=counts['SKIP'], missing=counts['MISSING'],
               infrastructure_failures=[dict(worker=w, reason=r) for w,r in infra],
               duration_sec=sum(row.get('duration_sec', 0) or 0 for row in rows), results=rows)
(root / 'batch-report.json').write_text(json.dumps(payload, indent=2) + '\n')
with (root / 'combined.tsv').open('w') as stream:
    writer = csv.writer(stream, delimiter='\t')
    keys = ('outcome', 'worker', 'category', 'scenario', 'duration_sec', 'report', 'log')
    writer.writerow(keys)
    writer.writerows([row.get(key, '') for key in keys] for row in rows)
with (root / 'infrastructure-failures.tsv').open('w') as stream:
    writer = csv.writer(stream, delimiter='\t')
    writer.writerow(('worker', 'reason'))
    writer.writerows(infra)
(root / 'failed.txt').write_text(''.join(f"{r['worker']}\t{r['scenario']}\t{r.get('log','')}\n"
                                       for r in rows if r['outcome'] in ('FAIL', 'MISSING')))
summary = (f"Ukama distributed P0 report\n\n"
           f"total={len(rows)} pass={counts['PASS']} fail={counts['FAIL']} "
           f"skip={counts['SKIP']} missing={counts['MISSING']} infra_issues={len(infra)}\n")
summary += ''.join(f"{r['outcome']:7} {r['worker']:10} {r['scenario']}\n" for r in rows)
if infra:
    summary += '\nInfrastructure failures:\n' + ''.join(f'{w}: {reason}\n' for w,reason in infra)
(root / 'summary.txt').write_text(summary)
print(summary)
sys.exit(1 if counts['FAIL'] or counts['MISSING'] or infra else 0)
