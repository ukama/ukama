/* SPDX-License-Identifier: MPL-2.0
 * C runner + local Chromium + controlled runtime/BFF/DOM fixture.
 * This establishes integration behavior, not verified console-app coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, writeFile, readFile, rm, mkdir, copyFile, chmod } from 'node:fs/promises';
import { join, resolve } from 'node:path';
import { tmpdir } from 'node:os';
import { spawn } from 'node:child_process';
import { provisioningFixture } from '../../adapters/webapp/test/provisioning-fixture.mjs';
const root=resolve(import.meta.dirname,'../..');
const binary=process.env.ULAB_TEST_BINARY || join(root,'bin/ukama-lab');
async function run(t, example, mode='', buildFail=false) {
  const app=await provisioningFixture(mode), dir=await mkdtemp(join(tmpdir(),'ulab-provision-'));
  t.after(async()=>{ await app.close(); if (!process.env.ULAB_KEEP_FIXTURE) await rm(dir,{recursive:true,force:true}); });
  const state=join(dir,'auth.json'); await writeFile(state,JSON.stringify(app.state));
  const scripts=join(dir,'scripts'); await mkdir(scripts);
  for(const name of ['ensure-network.sh','build-and-start-site.sh','wait-nodes-ready.sh','disconnect-node.sh','reconnect-node.sh','stop-node.sh','stop-media.sh','cleanup-network.sh']) {
    await writeFile(join(scripts,name), `#!/bin/sh
exec "${process.execPath}" "${join(root,'tests/webapp/runtime-fixture.mjs')}" "$0" "$@"
`); await chmod(join(scripts,name),0o700);
  }
  let text=(await readFile(join(root,'scenarios/webapp/p0/network',example),'utf8')).replace('http://localhost:3000',app.origin).replace('.auth/owner.json',state).replace('headless: false','headless: true');
  text=text.replace('scenario_timeout_seconds: 3600','scenario_timeout_seconds: 40').replaceAll('timeout_seconds: 900','timeout_seconds: 30');
  if(mode==='late-site'||mode==='late-network'||mode==='opaque-network') text=text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 10').replaceAll('timeout_seconds: 30','timeout_seconds: 5').replaceAll('timeout_seconds: 900','timeout_seconds: 5');
  const file=join(dir,'scenario.yaml'); await writeFile(file,text);
  const runDir=join(dir,'test-run');
  const child=spawn(binary,['run',file,'--repo','/test/ukama','--scripts',scripts,'--bff',app.origin+'/graphql','--out',dir,'--run-id','test-run','--webapp-worker',join(root,'utils/webapp-worker.sh')],{cwd:root,env:{...process.env,ULAB_FIXTURE_ORIGIN:app.origin,ULAB_FIXTURE_BUILD_FAIL:buildFail?'1':'0'}});
  let output=''; child.stdout.on('data',s=>output+=s); child.stderr.on('data',s=>output+=s);
  const timer=setTimeout(()=>child.kill('SIGTERM'),45000);
  const code=await new Promise((resolve,reject)=>{child.on('error',reject);child.on('close',resolve);});clearTimeout(timer);
  const report=JSON.parse(await readFile(join(runDir,'report.json'),'utf8').catch(()=>{throw new Error(output)}));
  const journal=JSON.parse(await readFile(join(runDir,'webapp-resources.json'),'utf8'));
  const commands=(await readFile(join(runDir,'webapp-commands.jsonl'),'utf8')).trim().split('\n').map(JSON.parse);
  return {app,code,report,journal,commands,output};
}
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
