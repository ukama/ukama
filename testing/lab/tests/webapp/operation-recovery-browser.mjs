/* SPDX-License-Identifier: MPL-2.0
 * C CLI + worker + Chromium. Controlled source-shaped DOM, not React/live coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { operationsFixture } from '../../adapters/webapp/test/operations-fixture.mjs';
import { run } from './browser-harness.mjs';
const files=['observed-site-restart','switch-port-isolation','node-status-read-failure','site-status-read-failure','optimistic-timeout','distinct-session-recovery','operation-prerequisites','busy-controller-ports'];
const execute=(t,index,mode='',transform=x=>x)=>run(t,`operations/wb-${110+index}-${files[index]}.yaml`,mode,false,operationsFixture,transform);
for(let i=0;i<files.length;i++)test(`C/browser operation recovery: ${files[i]}`,async t=>{
  const r=await execute(t,i);assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');
  assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
  const mutations=r.app.db.operations.filter(x=>x.ui);
  if(i===0)assert.deepEqual(mutations.map(x=>x.ui),['restart-site']);
  if(i===1)assert.deepEqual(mutations.map(x=>[x.port,x.value]),[[1,false],[1,true],[2,false],[2,true],[3,false],[3,true],[9,false],[9,true],[9,false],[9,true]]);
  if(i===2||i===3||i===6)assert.equal(mutations.length,0);
  if(i===4)assert.deepEqual(mutations.map(x=>x.ui),['restart-node','restart-node']);
  if(i===5){assert.equal(mutations.filter(x=>x.ui==='restart-node').length,0);assert.equal(mutations.filter(x=>x.ui==='update').length,1);}
});
for(const [index,mode,label] of [[0,'controller-outage','Restart observation'],[0,'no-restart-cycle','Restart observation'],[0,'service-resumed','Cellular state'],[1,'port-cross-talk','Port states'],[1,'port-not-persisted','Port states'],[2,'status-fail-open','Restart node'],[3,'status-fail-open','Restart site'],[4,'timeout-early','Optimistic timeout'],[4,'timeout-stuck','Optimistic timeout']])
  test(`Reject operation regression: ${mode}/${index}`,async t=>{
    const r=await execute(t,index,mode,text=>text.replaceAll('check_timeout_seconds: 15','check_timeout_seconds: 3').replaceAll('timeout_seconds: 15','timeout_seconds: 3').replace(/(label: "Optimistic timeout"[\s\S]*?timeout_seconds:) 30/,'$1 12'));
    assert.notEqual(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');
    assert(r.report.results.some(x=>x.label===label&&x.state==='FAIL'),JSON.stringify(r.report.results));
    if(index===4)assert.equal(r.app.db.operations.filter(x=>x.ui==='restart-node').length,1,'assertion failure cannot cause automatic retry');
  });
test('peer rejects the same saved login before a second actor can mutate',async t=>{
  const r=await execute(t,5,'same-peer');assert.notEqual(r.code,0);assert.match(r.output,/distinct unexpired ukama_session/);
  assert.equal(r.app.db.operations.filter(x=>x.ui).length,0);assert.equal(r.journal.cleanup,'complete');
});
test('missing controller cannot arm a vacuous restart observation',async t=>{
  const r=await execute(t,0,'missing-controller');assert.notEqual(r.code,0);assert.match(r.output,/Missing or ambiguous controller/);
  assert.equal(r.app.db.operations.filter(x=>x.ui).length,0);assert.equal(r.journal.cleanup,'complete');
});
test('a failure after peer creation retains traces for both isolated contexts',async t=>{
  const r=await execute(t,5,'',text=>text.replace('expected: "${ULAB_SOFTWARE_CURRENT_VERSION}"','expected: "wrong-installed-version"').replaceAll('timeout_seconds: 15','timeout_seconds: 3'));
  assert.notEqual(r.code,0);assert.equal(r.journal.cleanup,'complete');
  const failed=r.report.results.find(x=>x.state==='FAIL'&&x.label==='Current version');assert(failed,r.output);
  for(const name of ['trace.zip','trace-peer.zip','failure.png'])assert(failed.artifacts.some(p=>p.endsWith('/'+name)),name);
  assert.equal(r.app.db.operations.filter(x=>x.ui).length,0);
});
