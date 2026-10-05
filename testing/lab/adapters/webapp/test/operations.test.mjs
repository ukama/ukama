/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { Worker } from '../dist/worker.js';
import { operationsFixture } from './operations-fixture.mjs';
const ok=r=>assert.equal(r.status,'ok',JSON.stringify(r));
async function session(t,mode='') {
  const app=await operationsFixture(mode),dir=await mkdtemp(join(tmpdir(),'ulab-operation-'));
  app.db.sites.push({id:'site-1',name:'test-site',network_id:'existing'});
  app.db.nodes.push({id:'uk-tnode-1',type:'Tower node',site:'site-1'},{id:'uk-cnode-1',type:'Controller node',site:'site-1'});
  const worker=new Worker();let sequence=0;
  t.after(async()=>{await worker.shutdown(worker.failed);await app.close();await rm(dir,{force:true,recursive:true});});
  const make=(action,inputs,ms=5000)=>({protocol:1,run_id:'operations',command_id:++sequence,action,inputs,deadline_ms:Date.now()+ms});
  const send=(action,inputs,ms)=>worker.handle(make(action,inputs,ms));
  const state=join(dir,'auth.json');await writeFile(state,JSON.stringify(app.state));
  ok(await send('init',{profile:{base_url:app.origin,auth_state:state,scenario_timeout_seconds:30},artifacts_dir:dir}));
  const scope={view:'network_node_detail',network_name:'existing-network',entity:{ref:'tower-site-001-001',id:'uk-tnode-1',text:'uk-tnode-1'}};
  ok(await send('web_open',scope));return {app,worker,send,make,scope};
}
test('replaying a lost confirmation result cannot restart twice',async t=>{
  const s=await session(t);ok(await s.send('web_action',{...s.scope,action:'open_restart'}));
  const c=s.make('web_action',{...s.scope,action:'confirm_restart'}),result=await s.worker.handle(c);ok(result);
  assert.deepEqual(await s.worker.handle(c),result);
  assert.equal(s.app.db.operations.filter(x=>x.ui==='restart-node').length,1);
});
test('an action for another bound entity fails before clicking',async t=>{
  const s=await session(t);const r=await s.send('web_action',{...s.scope,entity:{...s.scope.entity,id:'uk-cnode-1'},action:'open_restart'});
  assert.equal(r.error.code,'WRONG_ENTITY');assert.equal(s.app.db.operations.filter(x=>x.ui).length,0);
});
test('a new tab cannot inherit another tab detail scope',async t=>{
  const s=await session(t);ok(await s.send('web_tab',{tab:'secondary'}));
  const r=await s.send('web_action',{...s.scope,action:'open_restart'});assert.equal(r.error.code,'WRONG_VIEW');
});
test('reason substring checks preserve the full visible reason',async t=>{
  const s=await session(t,'status-error');const r=await s.send('web_field_equals',{view:'network_node_detail',label:'Restart reason',expected:'Cannot verify',match:'contains',requirement:'WEB-OPS-013'});
  ok(r);assert.equal(r.expected,'Cannot verify');assert.equal(r.actual,'Cannot verify operation status');
});
test('reason substring mismatch fails with observed text and artifacts',async t=>{
  const s=await session(t,'status-error');const r=await s.send('web_field_equals',{view:'network_node_detail',label:'Restart reason',expected:'Controller busy',match:'contains',requirement:'WEB-OPS-013'},400);
  assert.equal(r.status,'error');assert.equal(r.actual,'Cannot verify operation status');assert(r.artifacts.some(p=>p.endsWith('failure.png')));
});

test('node open waits for rendered identity after a URL change before checking connectivity',async t=>{
  const s=await session(t,'slow-detail');
  ok(await s.send('web_field_equals',{view:'network_node_detail',label:'Connectivity',expected:'Online',requirement:'WEB-NODE-001'}));
});
