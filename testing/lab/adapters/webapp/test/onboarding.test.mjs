/* SPDX-License-Identifier: MPL-2.0 -- controlled fixture tests, not live coverage. */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp,rm,writeFile,readFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { Worker } from '../dist/worker.js';
import { profile } from '../dist/contract.js';
import { onboardingFixture } from './onboarding-fixture.mjs';
async function start(t,mode='onboarding'){
 const app=await onboardingFixture(),dir=await mkdtemp(join(tmpdir(),'ulab-onboard-')),worker=new Worker();let id=0;
 const auth=join(dir,'auth.json');await writeFile(auth,JSON.stringify(app.state));
 t.after(async()=>{await worker.shutdown(worker.failed);await app.close();await rm(dir,{recursive:true,force:true})});
 const send=(action,inputs={},ms=3000)=>worker.handle({protocol:1,run_id:'onboard-fixture',command_id:++id,action,deadline_ms:Date.now()+ms,inputs});
 const init=await send('init',{profile:{base_url:app.origin,auth_state:auth,session_mode:mode,headless:true},artifacts_dir:dir});
 const event=(action,rest={})=>send('web_onboard',{view:'configure',context:{},action,...rest});return{app,dir,worker,send,init,event};
}
const ok=r=>assert.equal(r.status,'ok',JSON.stringify(r));
test('onboarding init creates an isolated context and does not claim authentication',async t=>{
 const s=await start(t);ok(s.init);assert.equal(s.init.actual.authenticated,false);assert.equal(s.init.actual.session_mode,'onboarding');assert.equal(s.app.db.documents,0);assert.equal(s.app.db.operations.length,0);
});
test('onboarding rejects default-mode mutation commands',async t=>{
 const s=await start(t);ok(s.init);assert.equal((await s.send('web_create_network',{ref:'net',name:'unsafe-network'})).error.code,'SESSION_MODE');assert.equal(s.app.db.operations.length,0);
});
test('default authenticated mode rejects onboarding events',async t=>{
 const s=await start(t,'authenticated');ok(s.init);assert.equal((await s.event('open',{value:'sims'})).error.code,'SESSION_MODE');
});
test('onboarding requires captured auth and forbids a separate auth origin',()=>{
 assert.throws(()=>profile({base_url:'http://console.test',auth_state:'none',session_mode:'onboarding'}));
 assert.throws(()=>profile({base_url:'http://console.test',auth_state:'owner.json',session_mode:'onboarding',auth_origin:'http://auth.test'}));
});
test('validation probe refuses a valid name without sending a mutation',async t=>{
 const s=await start(t);ok(s.init);ok(await s.event('open',{value:'add_network'}));ok(await s.event('fill',{label:'Network name',value:'valid-name'}));
 assert.equal((await s.event('validate_name',{label:'Network name'})).error.code,'INVALID_INPUT');assert.equal(s.app.db.operations.length,0);
});
test('wrong planned name fails with a prepared receipt and no backend mutation',async t=>{
 const s=await start(t);ok(s.init);ok(await s.event('open',{value:'add_network'}));ok(await s.event('fill',{label:'Network name',value:'other-name'}));
 const r=await s.event('submit_network',{context:{network:{ref:'net-001',name:'planned-name',id:''}},creation:{kind:'network',ref:'net-001',name:'planned-name'}});
 assert.equal(r.error.code,'WRONG_NAME');const receipt=JSON.parse(await readFile(join(s.dir,'onboard-fixture','creation-4.json'),'utf8'));assert.equal(receipt.state,'prepared');assert.deepEqual(receipt.bindings,[]);assert.equal(s.app.db.operations.length,0);
});
test('arbitrary configure URLs are rejected',async t=>{
 const s=await start(t);ok(s.init);assert.equal((await s.event('open',{value:'https://other.test/'})).error.code,'INVALID_INPUT');assert.equal(s.app.db.documents,0);
});
