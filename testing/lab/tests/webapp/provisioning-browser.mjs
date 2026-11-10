/* SPDX-License-Identifier: MPL-2.0
 * C runner + local Chromium + controlled runtime/BFF/DOM fixture.
 * This establishes integration behavior, not verified console-app coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { execFile } from 'node:child_process';
import { readFile, writeFile } from 'node:fs/promises';
import { dirname, join, resolve } from 'node:path';
import { promisify } from 'node:util';
import { run } from './browser-harness.mjs';
const recovery='wb-001-sites-online-recovery.yaml', nodes='wb-002-node-list-detail.yaml';
test('UI provisions two sites, observes 2/2 -> 1/2 -> 2/2 without reload, and removes only owned resources',async t=>{
  const r=await run(t,recovery,'multiple-components');assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');
  assert.deepEqual(r.report.results.filter(x=>x.kind==='check').map(x=>x.actual),['2/2','1/2','2/2']);
  assert.equal(r.app.db.documents,1);assert.equal(r.journal.resources.length,9);
  assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
  assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);const sites=r.app.db.operations.filter(x=>x.op==='addSite');assert.equal(sites.length,2);
  for(const {data} of sites) for(const kind of ['switch','backhaul','power']) assert.equal(data[kind+'_id'],kind+'-1');
  assert.equal(r.commands.filter(x=>x.direction==='request'&&x.message.action==='web_reload').length,0);
});
test('all three node cards and detail views match world IDs, model, site and connectivity',async t=>{
  const r=await run(t,nodes);assert.equal(r.code,0,r.output);assert.equal(r.report.checks.passed,25);
  const checks=r.report.results.filter(x=>x.kind==='check'&&x.label==='Serial #');assert.equal(checks.length,6);
  assert(checks.every(x=>x.actual.startsWith('uk-')));assert.equal(r.journal.cleanup,'complete');
});
test('WB-002 waits for the displayed location before naming the site',async t=>{
  const r=await run(t,nodes,'slow-location');
  assert.equal(r.code,0,r.output);assert.equal(r.report.checks.passed,25);
  const creates=r.app.db.operations.filter(x=>x.op==='addSite');
  assert.deepEqual(creates.map(x=>x.data.location),['51.05000, -114.07000']);
  assert.equal(r.journal.cleanup,'complete');
});
test('WB-002 stops before addSite if the wizard drops the displayed location',async t=>{
  const r=await run(t,nodes,'dropped-location');
  assert.notEqual(r.code,0);assert.match(r.report.error,/location/i);
  assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);
  assert.equal(r.journal.cleanup,'complete');
});
test('lost tower selection fails before site mutation and cleans the network and runtime',async t=>{
  const r=await run(t,recovery,'wrong-tower');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);
  assert.equal(r.journal.cleanup,'complete');assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.equal(r.journal.resources.length,7);assert.deepEqual(r.app.db.nodes,[]);assert(r.output.includes('lost the selected tower'));
});
test('successful mutation followed by UI timeout is journaled and cleaned without resubmission',async t=>{
  const r=await run(t,recovery,'late-network');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);
  assert.equal(r.journal.resources.length,7);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.deepEqual(r.app.db.nodes,[]);assert.equal(r.journal.uncertain_creation,false);
});
test('partial runtime failure reloads early node identities for teardown',async t=>{
  const r=await run(t,recovery,'',true);assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);
  assert.equal(r.journal.resources.length,3);assert.deepEqual(r.app.db.nodes,[]);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
  assert(r.app.db.operations.some(x=>x.runtime==='cleanup-network.sh'));
});
test('BFF teardown success:false cannot produce a passing run',async t=>{
  const r=await run(t,recovery,'cleanup-false');assert.notEqual(r.code,0);assert.equal(r.report.outcome,'FAIL');assert.equal(r.journal.cleanup,'failed');
});

test('site UI timeout retains the successful creation receipt and cleans the linked trio',async t=>{
  const r=await run(t,recovery,'late-site');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,1);
  assert.equal(r.journal.resources.length,8);assert.deepEqual(r.app.db.nodes,[]);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.equal(r.journal.uncertain_creation,false);
});
test('lost network identity is reconciled against the pre-run inventory without resubmission',async t=>{
  const r=await run(t,recovery,'opaque-network');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,1);
  assert.equal(r.journal.uncertain_creation,false);assert.equal(r.journal.cleanup,'complete');assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
});

test('rejected addSite fails WB-002 and removes its network and three unassigned nodes',async t=>{
  const r=await run(t,nodes,'rejected-site');
  assert.notEqual(r.code,0);assert.match(r.report.error,/Argument Validation Error/);
  assert.equal(r.journal.cleanup,'complete');assert.equal(r.journal.resources.length,4);
  assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
  assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
  assert.deepEqual(r.app.db.sites,[]);assert.deepEqual(r.app.db.nodes,[]);
  assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,1);
});

test('cleanup-webapp reclaims the network and nodes recorded by an older WB-002 run',async t=>{
  const r=await run(t,nodes,'rejected-site');
  const dir=r.report.artifacts.run_dir, path=join(dir,'webapp-resources.json');
  const journal=r.journal, net=journal.resources.find(x=>x.kind==='network');
  // Reproduce the old runner's skipped deletes and submitted-without-ID receipt.
  delete journal.global_names_absent;delete journal.bff_url;
  journal.uncertain_creation=true;journal.cleanup='failed';
  for(const intent of journal.creation_intents){delete intent.validation_rejected;if(intent.kind==='site')intent.state='pending';}
  for(const row of journal.resources){
    row.cleanup='failed';
    if(row.kind==='node')r.app.db.nodes.push({id:row.id,site:null});
  }
  r.app.db.networks.push({id:net.id,name:net.name});
  await writeFile(path,JSON.stringify(journal));
  const args=['cleanup-webapp',join(dirname(dir),'scenario.yaml'),'--out',dirname(dir),'--run-id',journal.run_id,'--bff',r.app.origin+'/graphql'];
  const binary=process.env.ULAB_TEST_BINARY||resolve('bin/ukama-lab');
  await promisify(execFile)(binary,args);
  const after=JSON.parse(await readFile(path,'utf8'));
  assert.equal(after.cleanup,'complete');assert(after.resources.every(x=>x.cleanup==='deleted'));
  assert.equal(after.run_result,'failed');assert.deepEqual(r.app.db.nodes,[]);
  assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
  const operations=r.app.db.operations.length;
  await promisify(execFile)(binary,args);
  assert.equal(r.app.db.operations.length,operations);
});

test('submitted site with a lost receipt is found in its run-owned network and deleted',async t=>{
  const shorten = text => text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15')
    .replaceAll('timeout_seconds: 30','timeout_seconds: 5');
  const r=await run(t,recovery,'opaque-site',false,undefined,shorten);
  assert.notEqual(r.code,0);assert.equal(r.journal.cleanup,'complete');
  const site=r.journal.resources.find(x=>x.kind==='site');
  assert.equal(site.observed_via,'cleanup_lookup');assert.equal(site.cleanup,'deleted');
  assert.deepEqual(r.app.db.nodes,[]);assert.deepEqual(r.app.db.sites,[]);
  assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
  assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,1);
});

test('a successful delete response with the network still visible fails the run',async t=>{
  const r=await run(t,nodes,'retained-network');
  assert.notEqual(r.code,0);assert.equal(r.report.cleanup,'failed');assert.equal(r.journal.cleanup,'failed');
  const network=r.journal.resources.find(x=>x.kind==='network');
  assert.equal(network.cleanup,'failed');assert.match(network.cleanup_error,/still in backend inventory/);
  assert.deepEqual(r.app.db.nodes,[]);assert.deepEqual(r.app.db.sites,[]);
});
