/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import test, { before, after } from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, writeFile, readFile, stat, rm, readdir } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { spawn } from 'node:child_process';
import { once } from 'node:events';
import { randomUUID } from 'node:crypto';
import { Worker } from '../dist/worker.js';
import { VIEWS } from '../dist/console-app.js';
import { fixture } from './fixture.mjs';

let app, directory, state;
before(async () => {
  app = await fixture(); directory = await mkdtemp(join(tmpdir(), 'ulab-webapp-test-'));
  state = join(directory, 'state.json'); await writeFile(state, JSON.stringify(app.state), { mode: 0o600 });
});
after(async () => { await app?.close(); if (directory) await rm(directory, { recursive: true, force: true }); });
const settings = () => ({ base_url: app.origin, auth_state: state, scenario_timeout_seconds: 60 });
async function session(t, options = {}) {
  const worker = new Worker(); let id = 0; const run = `test-${randomUUID()}`;
  t.after(() => worker.shutdown(worker.failed));
  const make = (action, inputs = {}, ms = 5000) => ({ protocol: 1, run_id: run, command_id: ++id, action, deadline_ms: Date.now() + ms, inputs });
  const send = (action, inputs, ms) => worker.handle(make(action, inputs, ms));
  const result = await send('init', { profile: { ...settings(), ...options }, artifacts_dir: directory });
  return { worker, result, make, send, run };
}
const ok = result => assert.equal(result.status, 'ok', JSON.stringify(result));
const openNetwork = (s, network = 'lab-network', view = 'network_home', entity) => s.send('web_open', { view, network_name: network, ...(entity ? { entity } : {}) });
const check = (s, type = 'web_kpi_equals', inputs = {}, ms) => s.send(type, { view: 'network_home', label: 'Sites online', expected: '2/2', requirement: 'WEB-TEST-001', ...inputs }, ms);

test('real Chromium session uses UI navigation, visible KPI, node count and detail identity', async t => {
  const s = await session(t); ok(s.result); assert.equal(s.result.actual.authenticated, true);
  ok(await openNetwork(s));
  const documents = app.metrics.documents;
  const c = s.make('web_kpi_equals', { view: 'network_home', label: 'Sites online', expected: '2/2', requirement: 'WEB-TEST-001' });
  const result = await s.worker.handle(c); ok(result); assert.equal(result.actual, '2/2');
  assert.deepEqual(await s.worker.handle(c), result); // lost response replay
  assert.equal(app.metrics.documents, documents, 'assertion/replay must not reload');
  ok(await openNetwork(s, 'lab-network', 'network_nodes'));
  ok(await s.send('web_field_equals', { view: 'network_nodes', label: 'Nodes count', expected: '3', requirement: 'WEB-TEST-002' }));
  const detail = await openNetwork(s, 'lab-network', 'network_node_detail', { ref: 'tower-site-001-001', id: 'tower-001', text: 'tower-001' });
  ok(detail); assert.equal(detail.bindings[0].id, 'tower-001');
  ok(await s.send('web_field_equals', { view: 'network_node_detail', label: 'Serial #', expected: 'tower-001', requirement: 'WEB-TEST-003' }));
  ok(await s.send('web_field_equals', { view: 'network_node_detail', label: 'Empty field', expected: '', requirement: 'WEB-TEST-004' }));
  ok(await s.send('web_action_available', { view: 'network_node_detail', label: 'Restart node', available: true, requirement: 'WEB-TEST-005' }));
  ok(await s.send('web_reload', {})); assert.equal(app.metrics.documents, documents + 1);
  const close = await s.send('close', {}); ok(close); assert.equal(close.run_status, 'passed');
  assert(!close.artifacts.some(p => p.endsWith('trace.zip')), 'successful run discards raw trace');
});

test('all 22 mapped views navigate via the fixture UI, including leaving org-wide pages', async t => {
  const s = await session(t); ok(s.result);
  for (const [name, view] of Object.entries(VIEWS)) {
    const entity = view.detail ? { ref: `${view.detail}-ref`, id: `${view.detail === 'node' ? 'tower' : 'site'}-001`, text: `${view.detail === 'node' ? 'tower' : 'site'}-001` } : undefined;
    ok(await s.send('web_open', { view: name, network_name: 'lab-network', ...(entity ? { entity } : {}) }));
  }
});

test('delayed visible update converges without reload', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s, 'delayed-network'));
  const count = app.metrics.documents;
  ok(await check(s)); assert.equal(app.metrics.documents, count);
});

test('wrong value fails, keeps private screenshot/trace and does not leak diagnostic query tokens', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s, 'wrong-network'));
  const count = app.metrics.documents;
  const result = await check(s, undefined, {}, 350);
  assert.equal(result.status, 'error'); assert.equal(result.expected, '2/2');
  assert.equal(result.actual, '1/2');
  assert.equal(app.metrics.documents, count);
  for (const name of ['failure.png', 'trace.zip', 'diagnostics.json', 'results.jsonl']) {
    const path = result.artifacts.find(p => p.endsWith(name)); assert(path, name);
    assert((await stat(path)).size > 0); assert.equal((await stat(path)).mode & 0o777, 0o600);
  }
  const diagnostics = await readFile(result.artifacts.find(p => p.endsWith('diagnostics.json')), 'utf8');
  assert(!diagnostics.includes('do-not-record')); assert(!diagnostics.includes('fixture_session'));
  const close = await s.send('close', {}); ok(close); assert.equal(close.run_status, 'failed');
});

test('ambiguous visible KPI fails rather than taking the first match', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s, 'duplicate-network'));
  const result = await check(s); assert.equal(result.error.code, 'AMBIGUOUS_LOCATOR');
});

test('disabled action is valid unavailable state', async t => {
  const s = await session(t); ok(s.result);
  ok(await openNetwork(s, 'disabled-network', 'network_node_detail', { ref: 'node-ref', id: 'tower-001', text: 'tower-001' }));
  ok(await s.send('web_action_available', { view: 'network_node_detail', label: 'Restart node', available: false, requirement: 'WEB-TEST-005' }));
});

test('missing action is not accepted as disabled', async t => {
  const s = await session(t); ok(s.result);
  ok(await openNetwork(s, 'hidden-network', 'network_node_detail', { ref: 'node-ref', id: 'tower-001', text: 'tower-001' }));
  const result = await s.send('web_action_available', { view: 'network_node_detail', label: 'Restart node', available: false, requirement: 'WEB-TEST-005' }, 300);
  assert.equal(result.status, 'error'); assert.equal(result.expected, false);
});

test('explicit empty node list proves zero; loading does not', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s, 'empty-network', 'network_nodes'));
  const input = { view: 'network_nodes', label: 'Nodes count', expected: '0', requirement: 'WEB-TEST-006' };
  ok(await s.send('web_field_equals', input));
  ok(await openNetwork(s, 'loading-network', 'network_nodes'));
  assert.equal((await s.send('web_field_equals', input, 300)).status, 'error');
});

test('table count excludes headers and accepts a rendered zero-row table', async t => {
  const s = await session(t); ok(s.result);
  ok(await openNetwork(s, 'table-network', 'business_members'));
  const input = { view: 'business_members', label: 'Members', expected_count: 2, requirement: 'WEB-TEST-007' };
  ok(await s.send('web_table_count_equals', input));
  ok(await openNetwork(s, 'empty-network', 'business_members'));
  ok(await s.send('web_table_count_equals', { ...input, expected_count: 0 }));
  ok(await openNetwork(s, 'bad-table-network', 'business_members'));
  assert.equal((await s.send('web_table_count_equals', { ...input, expected_count: 0 }, 300)).status, 'error');
});

test('view and network scope are checked while observing', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s, 'switch-away'));
  const result = await check(s, undefined, { expected: 'never' });
  assert.equal(result.error.code, 'WRONG_NETWORK');
});

test('saved session is validated in a fresh context; an empty state cannot reuse another login', async t => {
  const s = await session(t); ok(s.result);
  const empty = join(directory, 'empty-state.json'); await writeFile(empty, '{"cookies":[],"origins":[]}');
  const other = await session(t, { auth_state: empty });
  assert.equal(other.result.error.code, 'AUTH_REQUIRED'); assert.equal(other.result.run_status, 'failed');
});

test('reusing a command ID for different content and changing run ID fail closed', async t => {
  const s = await session(t); ok(s.result); const c = s.make('web_open', { view: 'network_home', network_name: 'lab-network' });
  ok(await s.worker.handle(c));
  const result = await s.worker.handle({ ...c, inputs: { ...c.inputs, network_name: 'other-network' } });
  assert.equal(result.error.code, 'COMMAND_CONFLICT');
  const changed = await s.worker.handle({ ...s.make('close'), run_id: 'different-run' });
  assert.equal(changed.error.code, 'RUN_MISMATCH');
});

test('deadline failure stops execution and wrong-view assertion does not navigate', async t => {
  const s = await session(t); ok(s.result); ok(await openNetwork(s));
  const before = app.metrics.documents;
  const result = await check(s, undefined, { view: 'business_home' });
  assert.equal(result.error.code, 'WRONG_VIEW'); assert.equal(app.metrics.documents, before);
  const other = await session(t); ok(other.result);
  const expired = other.make('web_reload'); expired.deadline_ms = Date.now() - 1;
  assert.equal((await other.worker.handle(expired)).error.code, 'DEADLINE_EXCEEDED');
});

test('worker JSONL stdout stays machine-readable and close returns process success', async () => {
  const run = `pipe-${randomUUID()}`; const child = spawn(process.execPath, ['dist/cli.js', 'worker'], { env: process.env });
  let stdout = '', stderr = ''; child.stdout.on('data', b => { stdout += b; }); child.stderr.on('data', b => { stderr += b; });
  const completed = once(child, 'exit');
  const commands = [
    ['init', { profile: settings(), artifacts_dir: directory }],
    ['web_open', { view: 'network_home', network_name: 'lab-network' }],
    ['web_kpi_equals', { view: 'network_home', label: 'Sites online', expected: '2/2', requirement: 'WEB-TEST-001' }],
    ['close', {}],
  ].map(([action, inputs], i) => ({ protocol: 1, run_id: run, command_id: i + 1, action, inputs, deadline_ms: Date.now() + 30000 }));
  child.stdin.end(commands.map(c => JSON.stringify(c)).join('\n') + '\n');
  const [code] = await completed; assert.equal(code, 0, stderr + stdout);
  const results = stdout.trim().split('\n').map(JSON.parse); assert.equal(results.length, 4); results.forEach(ok);
  assert.equal(results.at(-1).run_status, 'passed');
});

test('EOF without close fails and preserves partial-run artifacts', async () => {
  const run = `eof-${randomUUID()}`; const child = spawn(process.execPath, ['dist/cli.js', 'worker'], { env: process.env });
  let stdout = ''; child.stdout.on('data', b => { stdout += b; }); child.stderr.resume(); const completed = once(child, 'exit');
  child.stdin.end(JSON.stringify({ protocol: 1, run_id: run, command_id: 1, action: 'init', inputs: { profile: settings(), artifacts_dir: directory }, deadline_ms: Date.now() + 30000 }) + '\n');
  const [code] = await completed; assert.equal(code, 1); ok(JSON.parse(stdout.trim()));
  assert((await readdir(join(directory, run))).includes('trace.zip'));
  const summary = JSON.parse(await readFile(join(directory, run, 'worker-summary.json'), 'utf8'));
  assert.equal(summary.status, 'failed'); assert.equal(summary.reason, 'INPUT_EOF');
});

test('SIGTERM closes an idle worker and records cancellation', async () => {
  const run = `signal-${randomUUID()}`;
  const child = spawn(process.execPath, ['dist/cli.js', 'worker'], { env: process.env });
  child.stderr.resume(); const completed = once(child, 'exit');
  const ready = new Promise(resolve => child.stdout.once('data', resolve));
  child.stdin.write(JSON.stringify({ protocol: 1, run_id: run, command_id: 1, action: 'init', inputs: { profile: settings(), artifacts_dir: directory }, deadline_ms: Date.now() + 30000 }) + '\n');
  ok(JSON.parse(String(await ready).trim())); child.kill('SIGTERM');
  const [code] = await completed; assert.equal(code, 143);
  const summary = JSON.parse(await readFile(join(directory, run, 'worker-summary.json'), 'utf8'));
  assert.equal(summary.reason, 'SIGTERM'); assert.equal(summary.status, 'failed');
});

test('scenario deadline terminates an idle protocol worker', async () => {
  const run = `deadline-${randomUUID()}`;
  const child = spawn(process.execPath, ['dist/cli.js', 'worker'], { env: process.env });
  let stdout = ''; child.stdout.on('data', b => { stdout += b; }); child.stderr.resume();
  const completed = once(child, 'exit');
  child.stdin.write(JSON.stringify({ protocol: 1, run_id: run, command_id: 1, action: 'init',
    inputs: { profile: { ...settings(), action_timeout_seconds: 1, check_timeout_seconds: 1, scenario_timeout_seconds: 1 }, artifacts_dir: directory }, deadline_ms: Date.now() + 30000 }) + '\n');
  const [code] = await completed; assert.equal(code, 1); ok(JSON.parse(stdout.trim()));
  const summary = JSON.parse(await readFile(join(directory, run, 'worker-summary.json'), 'utf8'));
  assert.equal(summary.reason, 'SCENARIO_DEADLINE');
});

test('standalone smoke command succeeds only on a matching visible value', async () => {
  const config = join(directory, 'profile.json'); await writeFile(config, JSON.stringify(settings()));
  for (const expected of ['2/2', '99/99']) {
    await writeFile(config, JSON.stringify({ ...settings(), check_timeout_seconds: 1 }));
    const child = spawn(process.execPath, ['dist/cli.js', 'smoke', '--profile', config, '--network', 'lab-network', '--view', 'network_home', '--label', 'Sites online', '--expected', expected, '--artifacts', directory], { env: process.env });
    let stdout = '', stderr = ''; child.stdout.on('data', b => { stdout += b; }); child.stderr.on('data', b => { stderr += b; });
    const [code] = await once(child, 'exit');
    assert.equal(code, expected === '2/2' ? 0 : 1, stderr);
    const results = stdout.trim().split('\n').map(JSON.parse);
    assert.equal(results.at(-1).run_status, expected === '2/2' ? 'passed' : 'failed');
    assert.match(stderr, /no scenario coverage credit/);
  }
});

test('C-owned failures retain browser evidence and cannot close as passed', async t => {
  const s = await session(t); ok(s.result);
  const result = await s.send('close', { failed: true, reason: 'C_RUNNER_FAILED' });
  ok(result); assert.equal(result.run_status, 'failed');
  assert(result.artifacts.some(p => p.endsWith('trace.zip')));
  const summary = JSON.parse(await readFile(result.artifacts.find(p => p.endsWith('worker-summary.json')), 'utf8'));
  assert.equal(summary.status, 'failed');
  assert.equal(summary.reason, 'C_RUNNER_FAILED');
});
