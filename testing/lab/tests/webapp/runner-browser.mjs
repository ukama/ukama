/* SPDX-License-Identifier: MPL-2.0
 * C CLI -> JSONL worker -> real Chromium -> source-derived DOM fixture.
 * These tests establish plumbing only, not console-app product coverage.
 */
import test, { before, after } from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, writeFile, readFile, rm } from 'node:fs/promises';
import { join, resolve } from 'node:path';
import { tmpdir } from 'node:os';
import { spawn } from 'node:child_process';
import { fixture } from '../../adapters/webapp/test/fixture.mjs';

const root = resolve(import.meta.dirname, '../..');
const binary = resolve(process.env.ULAB_TEST_BINARY || join(root, 'bin/ukama-lab'));
let app, directory, example, state;
before(async () => {
  app = await fixture(); directory = await mkdtemp(join(tmpdir(), 'ulab-c-browser-'));
  state = join(directory, 'auth.json'); await writeFile(state, JSON.stringify(app.state), { mode: 0o600 });
  example = (await readFile(join(root, 'scenarios/webapp/p0/session/wb-000-authenticated-members.yaml'), 'utf8'))
    .replace('http://localhost:3000', app.origin).replace('.auth/owner.json', state);
});
after(async () => { await app?.close(); if (directory) await rm(directory, { recursive: true, force: true }); });
async function run(name, text, cancel = false) {
  const scenario = join(directory, `${name}.yaml`); await writeFile(scenario, text);
  const runDir = join(directory, name);
  const child = spawn(binary, ['run', scenario, '--repo', '/test/ukama', '--out', directory, '--run-id', name,
    '--webapp-worker', join(root, 'utils/webapp-worker.sh')], { cwd: root, env: process.env });
  let output = ''; child.stdout.on('data', s => { output += s; }); child.stderr.on('data', s => { output += s; });
  const done = new Promise((resolve, reject) => { child.on('error', reject); child.on('close', code => resolve(code)); });
  const timer = setTimeout(() => child.kill('SIGKILL'), 30000);
  try {
    if (cancel) {
      const end = Date.now() + 15000; let started = false;
      while (Date.now() < end) {
        const journal = await readFile(join(runDir, 'webapp-commands.jsonl'), 'utf8').catch(() => '');
        if (journal.includes('web_action_available')) { started = true; break; }
        await new Promise(r => setTimeout(r, 20));
      }
      assert(started, output); child.kill('SIGTERM');
    }
    const code = await done;
    const report = JSON.parse(await readFile(join(runDir, 'report.json'), 'utf8'));
    const messages = (await readFile(join(runDir, 'webapp-commands.jsonl'), 'utf8')).trim().split('\n').map(JSON.parse);
    const summary = JSON.parse(await readFile(join(runDir, 'browser', name, 'worker-summary.json'), 'utf8'));
    return { code, report, messages, summary, output, runDir };
  } finally { clearTimeout(timer); if (child.exitCode === null && child.signalCode === null) child.kill('SIGKILL'); }
}
test('full CLI preserves phase order, explicit reload, visible checks and final result', async () => {
  const before = app.metrics.documents;
  const r = await run('browser-pass', example);
  assert.equal(r.code, 0, r.output); assert.equal(r.report.outcome, 'PASS');
  assert.equal(r.report.checks.passed, 2); assert.equal(r.summary.status, 'passed');
  const requests = r.messages.filter(x => x.direction === 'request').map(x => x.message);
  assert.deepEqual(requests.map(x => x.action), ['init', 'web_open', 'web_action_available', 'web_reload', 'web_action_available', 'close']);
  assert.equal(app.metrics.documents - before, 2, 'only initial navigation and explicit reload fetch the document');
  for (const c of r.report.results.filter(x => x.kind === 'check')) {
    assert.equal(c.expected, true); assert.equal(c.actual, true); assert.equal(c.requirement, 'WEB-TEAM-002');
  }
});
test('visible assertion failure retains trace and stops later phases', async () => {
  const r = await run('browser-fail', example.replaceAll('Invite member', 'Missing fixture control').replace('check_timeout_seconds: 30', 'check_timeout_seconds: 1'));
  assert.notEqual(r.code, 0, r.output); assert.equal(r.report.outcome, 'FAIL'); assert.equal(r.summary.status, 'failed');
  assert.equal(r.report.checks.failed, 1); assert.equal(r.report.events.total, 1);
  const check = r.report.results.find(x => x.kind === 'check');
  assert.equal(check.expected, true); assert(check.artifacts.some(p => p.endsWith('trace.zip')));
});
test('cancelling a live browser command retains complete evidence across repeated shutdown races', async () => {
  for (let attempt = 0; attempt < 5; attempt++) {
    const name = `browser-cancel-${attempt}`;
    const r = await run(name, example.replaceAll('Invite member', 'Missing fixture control'), true);
    assert.notEqual(r.code, 0, r.output); assert.equal(r.report.outcome, 'FAIL'); assert.equal(r.summary.status, 'failed');
    await readFile(join(r.runDir, 'browser', name, 'trace.zip'));
  }
});
