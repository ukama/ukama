/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { profile, command, canonical, safeURL, Budget } from '../dist/contract.js';
import { getView, VIEWS } from '../dist/console-app.js';
import { Worker } from '../dist/worker.js';

const config = { base_url: 'http://localhost:3000', auth_state: '.auth/owner.json' };
const request = { protocol: 1, run_id: 'test-1', command_id: 1, action: 'init', deadline_ms: Date.now() + 1000, inputs: {} };
test('profile defaults and validation reject mistaken or unsafe configuration', () => {
  assert.equal(profile(config).browser, 'chromium');
  assert.equal(profile(config).headless, true);
  for (const patch of [ { unknown: true }, { headless: 'false' }, { browser: 'magic' },
    { base_url: 'ftp://localhost' }, { base_url: 'http://name:secret@localhost' },
    { base_url: 'http://localhost/path' }, { base_url: 'http://localhost?token=secret' },
    { action_timeout_seconds: 0 }, { check_timeout_seconds: 901 },
    { scenario_timeout_seconds: 2 }, { auth_state: '' } ]) assert.throws(() => profile({ ...config, ...patch }));
});
test('protocol strictly validates version, IDs, action and envelope', () => {
  assert.equal(command(request).command_id, 1);
  for (const patch of [{ protocol: 2 }, { run_id: '../escape' }, { command_id: 0 },
    { command_id: 10001 }, { deadline_ms: 'tomorrow' }, { action: 'click_anything' }, { inputs: [] }, { extra: true }])
    assert.throws(() => command({ ...request, ...patch }));
});
test('fingerprints ignore JSON object key order and diagnostic URLs remove credentials/query/hash', () => {
  assert.equal(canonical({ b: 1, a: { d: 2, c: 3 } }), canonical({ a: { c: 3, d: 2 }, b: 1 }));
  assert.equal(safeURL('https://user:secret@example.com/page?token=secret#session'), 'https://example.com/page');
});
test('only mapped dashboard views have handlers', () => {
  assert.equal(Object.keys(VIEWS).length, 22);
  assert.equal(getView('network_node_detail').detail, 'node');
  assert.throws(() => getView('welcome'));
  assert.throws(() => getView('billing'));
});
test('poll observes changes without resubmitting an operation', async () => {
  let reads = 0;
  assert.equal(await new Budget(1000).poll(async () => ++reads, n => n === 3, 'missing'), 3);
  assert.equal(reads, 3);
});
test('worker rejects commands before init and refuses a success result on close', async () => {
  const worker = new Worker();
  const result = await worker.handle({ ...request, action: 'web_reload' });
  assert.equal(result.error.code, 'NOT_INITIALIZED'); assert.equal(result.run_status, 'failed');
  const close = await worker.handle({ ...request, command_id: 2, action: 'close' });
  assert.equal(close.run_status, 'failed'); assert.equal(close.status, 'error');
});
