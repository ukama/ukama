/* SPDX-License-Identifier: MPL-2.0
 * C runner + local Chromium + controlled runtime/BFF/DOM fixture.
 * This establishes integration behavior, not verified console-app coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { run } from './browser-harness.mjs';
const recovery='wb-001-sites-online-recovery.yaml', nodes='wb-002-node-list-detail.yaml';
test('UI provisions two sites, observes 2/2 -> 1/2 -> 2/2 without reload, and removes only owned resources',async t=>{
  const r=await run(t,recovery);assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');
  assert.deepEqual(r.report.results.filter(x=>x.kind==='check').map(x=>x.actual),['2/2','1/2','2/2']);
  assert.equal(r.app.db.documents,1);assert.equal(r.journal.resources.length,9);
  assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
  assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,2);
  assert.equal(r.commands.filter(x=>x.direction==='request'&&x.message.action==='web_reload').length,0);
});
test('all three node cards and detail views match world IDs, model, site and connectivity',async t=>{
  const r=await run(t,nodes);assert.equal(r.code,0,r.output);assert.equal(r.report.checks.passed,25);
  const checks=r.report.results.filter(x=>x.kind==='check'&&x.label==='Serial #');assert.equal(checks.length,6);
  assert(checks.every(x=>x.actual.startsWith('uk-')));assert.equal(r.journal.cleanup,'complete');
});
test('lost tower selection fails before site mutation and cleans the network and runtime',async t=>{
  const r=await run(t,recovery,'wrong-tower');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);
  assert.equal(r.journal.cleanup,'complete');assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.equal(r.journal.resources.length,4);assert(r.output.includes('lost the selected tower'));
});
test('successful mutation followed by UI timeout is journaled and cleaned without resubmission',async t=>{
  const r=await run(t,recovery,'late-network');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);
  assert.equal(r.journal.resources.length,1);assert.equal(r.journal.resources[0].cleanup,'deleted');assert.equal(r.journal.uncertain_creation,false);
});
test('partial runtime failure reloads early node identities for teardown',async t=>{
  const r=await run(t,recovery,'',true);assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);
  assert.equal(r.journal.resources.length,4);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
  assert(r.app.db.operations.some(x=>x.runtime==='cleanup-network.sh'));
});
test('BFF teardown success:false cannot produce a passing run',async t=>{
  const r=await run(t,recovery,'cleanup-false');assert.notEqual(r.code,0);assert.equal(r.report.outcome,'FAIL');assert.equal(r.journal.cleanup,'failed');
});

test('site UI timeout retains the successful creation receipt and cleans the linked trio',async t=>{
  const r=await run(t,recovery,'late-site');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,1);
  assert.equal(r.journal.resources.length,5);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.equal(r.journal.uncertain_creation,false);
});
test('unidentified submitted mutation is never retried or guessed away during cleanup',async t=>{
  const r=await run(t,recovery,'opaque-network');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);
  assert.equal(r.journal.uncertain_creation,true);assert.equal(r.journal.cleanup,'failed');assert.equal(r.app.db.operations.filter(x=>x.op==='deleteNetwork').length,0);
});
