/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {mkdtemp,writeFile,rm} from 'node:fs/promises';
import {join} from 'node:path';
import {tmpdir} from 'node:os';
import {Worker} from '../dist/worker.js';
import {inventoryFixture} from './inventory-fixture.mjs';
async function start(t){
 const app=await inventoryFixture(),directory=await mkdtemp(join(tmpdir(),'inventory-worker-'));
 app.db.sites.push({id:'s1',name:'site-one',network_id:'existing'});
 app.db.nodes.push({id:'tnode-1',type:'Tower node',site:'s1'});
 const auth=join(directory,'auth.json');await writeFile(auth,JSON.stringify(app.state));
 const worker=new Worker();let id=0;
 t.after(async()=>{await worker.shutdown(worker.failed);await app.close();await rm(directory,{recursive:true,force:true})});
 const send=(action,inputs={},timeout=3000)=>worker.handle({protocol:1,run_id:'inventory-unit',command_id:++id,action,inputs,deadline_ms:Date.now()+timeout});
 const result=await send('init',{profile:{base_url:app.origin,auth_state:auth,scenario_timeout_seconds:30},artifacts_dir:directory});assert.equal(result.status,'ok',JSON.stringify(result));
 const context={network:{ref:'net-001',id:'existing',name:'existing-network'},sites:[{ref:'site-001',id:'s1',name:'site-one',network_ref:'net-001'}],nodes:[{ref:'tower-site-001-001',id:'tnode-1',network_ref:'net-001',site_ref:'site-001'}]};
 return {app,send,context};
}
const inputs=(s,extra={})=>({view:'network_home',context:s.context,...extra});
const check=(s,extra={})=>inputs(s,{label:'Map count',expected:'1',requirement:'WEB-TEST-001',...extra});
test('inventory missing watch cannot pass a no-leak assertion',async t=>{
 const s=await start(t);const r=await s.send('web_inventory_equals',check(s,{label:'Scope leaks',expected:'none'}));assert.equal(r.error.code,'MISSING_OBSERVER');
});
test('inventory direct navigation rejects arbitrary target views',async t=>{
 const s=await start(t);const r=await s.send('web_inventory',inputs(s,{view:'https://untrusted.example',action:'direct'}));assert.equal(r.error.code,'INVALID_INPUT');
});
test('inventory rejects a foreign entity reference before navigation',async t=>{
 const s=await start(t);s.context.sites[0].network_ref='net-002';const r=await s.send('web_inventory',inputs(s,{view:'network_site_detail',site_ref:'site-001',action:'direct'}));assert.equal(r.error.code,'WRONG_ENTITY');
});
test('inventory will not read a marker color without selecting that site',async t=>{
 const s=await start(t);const r=await s.send('web_inventory_equals',check(s,{site_ref:'site-001',label:'Map color',expected:'var(--uk-success-bright)'}));assert.equal(r.error.code,'WRONG_ENTITY');
});
test('inventory faults have no arbitrary GraphQL or mutation command',async t=>{
 const s=await start(t);const r=await s.send('web_inventory',inputs(s,{action:'mask_home',value:'deleteNetwork'}));assert.equal(r.error.code,'INVALID_INPUT');assert.equal(s.app.db.networks.length,1);
});
test('inventory membership ignores ordering but retains complete world identities',async t=>{
 const s=await start(t);assert.equal((await s.send('web_open',{view:'network_nodes',network_name:'existing-network'})).status,'ok');
 const r=await s.send('web_inventory_equals',check(s,{view:'network_nodes',label:'Node IDs',expected:['tnode-1']}));assert.equal(r.status,'ok',JSON.stringify(r));assert.deepEqual(r.actual,['tnode-1']);
});
test('inventory unknown assertion labels fail instead of counting missing data as empty',async t=>{
 const s=await start(t);const r=await s.send('web_inventory_equals',check(s,{label:'Arbitrary DOM',expected:''}));assert.equal(r.error.code,'INVALID_INPUT');
});
test('inventory scope watcher requires independent foreign identities',async t=>{
 const s=await start(t);const r=await s.send('web_inventory',inputs(s,{action:'switch_watch'}));assert.equal(r.error.code,'INVALID_INPUT');
});

test('inventory rejects a malformed extra card without a serial',async t=>{
 const s=await start(t);s.app.db.mode='malformed-node';assert.equal((await s.send('web_open',{view:'network_nodes',network_name:'existing-network'})).status,'ok');
 const r=await s.send('web_inventory_equals',check(s,{view:'network_nodes',label:'Node IDs',expected:['tnode-1']}));assert.equal(r.error.code,'MALFORMED_INVENTORY_CARD');
});
