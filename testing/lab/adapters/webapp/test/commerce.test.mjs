/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, writeFile, readFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { Worker } from '../dist/worker.js';
import { planValues } from '../dist/commerce.js';
import { commerceFixture } from './commerce-fixture.mjs';
const ok=r=>assert.equal(r.status,'ok',JSON.stringify(r));
const plan=(name='weekly',minutes=10080,mb=1024,unit='GB')=>({ref:name+'__net-001',id:'',name,data_mb:mb,duration_minutes:minutes,amount:10,currency:'USD',country:'USA',unit,organization:false});
async function session(t,mode='') {
 const app=await commerceFixture(mode),dir=await mkdtemp(join(tmpdir(),'ulab-commerce-')),worker=new Worker();let sequence=0;
 t.after(async()=>{await worker.shutdown(worker.failed);await app.close();await rm(dir,{force:true,recursive:true});});
 const make=(action,inputs,ms=5000)=>({protocol:1,run_id:'commerce',command_id:++sequence,action,inputs,deadline_ms:Date.now()+ms});
 const send=(action,inputs,ms)=>worker.handle(make(action,inputs,ms));
 const state=join(dir,'auth.json');await writeFile(state,JSON.stringify(app.state));
 ok(await send('init',{profile:{base_url:app.origin,auth_state:state,scenario_timeout_seconds:90},artifacts_dir:dir}));
 const scope={network_name:'existing-network',network_id:'existing'};
 const open=view=>send('web_open',{view,network_name:'existing-network'});
 const act=(action,extra={},ms)=>send('web_commerce',{view:action.includes('plan')?'business_data_plans':'customer_customers',...scope,action,...extra},ms);
 const check=(label,expected,extra={},ms)=>send('web_commerce_equals',{view:'customer_customers',label,expected,requirement:'WEB-COM-001',...extra},ms);
 return {app,dir,worker,make,send,open,act,check,scope};
}
async function createPlan(s,p) {ok(await s.open('business_data_plans'));const r=await s.act('create_plan',{plan:p,creation:{kind:'package',ref:p.ref,name:p.name}});ok(r);p.id=r.bindings[0].id;return p;}
async function allocated(s) {
 const p=await createPlan(s,plan());const c={ref:'sub-000001',id:'',name:'Lab User 1',email:'lab@example.test',iccid:'8901000000000000001',sim_id:''};
 s.app.db.pool.push({iccid:c.iccid});ok(await s.open('customer_customers'));
 let r=await s.act('create_customer',{customer:c,creation:{kind:'subscriber',ref:c.ref,name:c.name}});ok(r);c.id=r.bindings[0].id;
 ok(await s.act('open_customer',{customer:c}));r=await s.act('allocate_sim',{customer:c,plan:p,creation:{kind:'sim',ref:'ue-000001',name:c.iccid}});ok(r);c.sim_id=r.bindings[0].id;
 return {p,c,checkScope:{customer_name:c.name,iccid:c.iccid}};
}
test('day/minute and binary-unit adapter arithmetic rejects unsupported conversions',()=>{
 for(const [m,d] of [[1440,1],[10080,7],[43200,30]])assert.equal(planValues(plan('p',m)).days,d);
 assert.equal(planValues(plan()).volume,1);assert.equal(planValues(plan('p',1440,512,'MB')).volume,512);
 for(const p of [plan('p',60),plan('p',1441),plan('p',1440,1000,'GB'),{...plan(),amount:NaN}])assert.throws(()=>planValues(p));
});
test('all three validity choices survive reload and read-only edit; units remain exact',async t=>{
 const s=await session(t);
 for(const [minutes,days,amount,unit] of [[1440,'1 day',512,'MB'],[10080,'1 week',1024,'GB'],[43200,'1 month',2048,'GB']]){
  const p=await createPlan(s,plan('plan'+minutes,minutes,amount,unit));ok(await s.send('web_reload',{}));
  ok(await s.check('Plan terms',`${unit==='MB'?amount:amount/1024} ${unit} data · ${days} validity`,{view:'business_data_plans',plan_name:p.name}));
  ok(await s.act('edit_plan',{plan:p}));ok(await s.check('Validity',days,{view:'business_data_plans',plan_name:p.name}));
  ok(await s.act('close_dialog',{view:'business_data_plans'}));
 }
 assert.equal(s.app.db.operations.filter(x=>x.op==='addPackage').length,3);
});
test('wrong persisted minutes fail against the visible plan rather than mutation success',async t=>{
 const s=await session(t,'wrong-minutes'),p=await createPlan(s,plan());ok(await s.send('web_reload',{}));
 const r=await s.check('Plan terms','1 GB data · 1 week validity',{view:'business_data_plans',plan_name:p.name},500);assert.equal(r.status,'error');assert.match(r.actual,/7 minutes/);assert(r.artifacts.some(x=>x.endsWith('failure.png')));
});
test('successful plan submission followed by UI timeout preserves ownership receipt',async t=>{
 const s=await session(t,'late-plan'),p=plan();ok(await s.open('business_data_plans'));
 const r=await s.act('create_plan',{plan:p,creation:{kind:'package',ref:p.ref,name:p.name}},700);assert.equal(r.status,'error');
 const receipt=JSON.parse(await readFile(join(s.dir,'commerce','creation-3.json'),'utf8'));assert.equal(receipt.state,'identified');assert.equal(receipt.bindings[0].id,'plan-1');assert.equal(s.app.db.plans.length,1);
});
test('Factory CSV import uses file picker and reads exact ICCID status from pool UI',async t=>{
 const s=await session(t),csv=join(s.dir,'factory-sims.csv'),iccid='8901000000000000001';await writeFile(csv,`iccid,imsi\n${iccid},001010000000001\n`);
 ok(await s.send('web_import_sims',{csv_path:csv,iccids:[iccid]}));ok(await s.check('Pool status','Available',{view:'business_sim_pool',iccid}));
});
test('customer creation, exact SIM allocation, allowance and reload are verified on screen',async t=>{
 const s=await session(t),{p,c,checkScope}=await allocated(s);
 ok(await s.check('ICCID',c.iccid,checkScope));ok(await s.check('Cycle usage','0 B of 1 GB used this cycle',checkScope));ok(await s.check('Package status','Current',{...checkScope,plan_name:p.name}));
 ok(await s.send('web_reload',{}));ok(await s.act('open_customer',{customer:c}));ok(await s.check('ICCID',c.iccid,checkScope));
});
async function paid(s) {
 const top=await createPlan(s,{...plan('topup',43200,2048),amount:25}),a=await allocated(s);
 ok(await s.act('cancel_top_up',{customer:a.c,plan:top}));assert.equal(s.app.db.payments.length,0);
 const command=s.make('web_commerce',{...s.scope,view:'customer_customers',action:'top_up',customer:a.c,plan:top,creation:{kind:'payment',ref:'ue-000001',name:a.c.iccid}});
 const r=await s.worker.handle(command);ok(r);assert.deepEqual(await s.worker.handle(command),r);assert.equal(s.app.db.payments.length,1);
 return {...a,top,payment:r.bindings[0].id};
}
test('cancel and replay never duplicate cash top-up; entitlement and receipt reconcile',async t=>{
 const s=await session(t),a=await paid(s);ok(await s.check('Package count','1',{...a.checkScope,plan_name:a.top.name}));ok(await s.check('Package status','Upcoming',{...a.checkScope,plan_name:a.top.name}));
 ok(await s.act('open_receipt',{customer:a.c,plan:a.top}));
 for(const [label,value] of [['Receipt payment ID',a.payment],['Receipt total','$25.00'],['Receipt method','Cash'],['Receipt status','Completed'],['Receipt plan',a.top.name]])ok(await s.check(label,value,{...a.checkScope,plan_name:a.top.name}));
});
test('completed payment cannot pass a missing-entitlement assertion',async t=>{
 const s=await session(t,'no-entitlement'),a=await paid(s);const r=await s.check('Package count','1',{...a.checkScope,plan_name:a.top.name},500);assert.equal(r.status,'error');assert.equal(r.actual,'0');
});
test('receipt for another payment fails exact visible correlation',async t=>{
 const s=await session(t,'wrong-receipt'),a=await paid(s);ok(await s.act('open_receipt',{customer:a.c,plan:a.top}));const r=await s.check('Receipt payment ID',a.payment,a.checkScope,500);assert.equal(r.status,'error');assert.equal(r.actual,'foreign-payment');
});
test('SIM toggle and binary usage remain visible assertions',async t=>{
 const s=await session(t),a=await allocated(s);
 ok(await s.act('deactivate_sim',{customer:a.c}));ok(await s.check('SIM status','Inactive',a.checkScope));ok(await s.act('activate_sim',{customer:a.c}));ok(await s.check('SIM status','Active',a.checkScope));
 s.app.db.usage=64*1024*1024;ok(await s.send('web_reload',{}));ok(await s.act('open_customer',{customer:a.c}));ok(await s.check('Cycle usage','64 MB of 1 GB used this cycle',a.checkScope));ok(await s.check('Total usage','64 MB',a.checkScope));
});
test('wrong customer scope fails before a financial mutation',async t=>{
 const s=await session(t),a=await allocated(s);const r=await s.act('top_up',{customer:{...a.c,iccid:'8901000000000000002'},plan:a.p,creation:{kind:'payment',ref:'ue-000001',name:'8901000000000000002'}},500);assert.equal(r.status,'error');assert.equal(s.app.db.payments.length,0);
});
test('a new command ID cannot repeat the same commerce mutation reference',async t=>{
 const s=await session(t),a=await paid(s);
 const r=await s.act('top_up',{customer:a.c,plan:a.top,creation:{kind:'payment',ref:'ue-000001',name:a.c.iccid}});
 assert.equal(r.error.code,'DUPLICATE_MUTATION');assert.equal(s.app.db.payments.length,1);
});
test('auto-assignment refuses a foreign available option before any allocation',async t=>{
 const s=await session(t),p=await createPlan(s,plan());
 const c={ref:'sub-000001',id:'',name:'Lab User 1',email:'lab@example.test',iccid:'8901000000000000001',sim_id:''};
 s.app.db.pool.push({iccid:c.iccid},{iccid:'8901000000000000002'});ok(await s.open('customer_customers'));
 const made=await s.act('create_customer',{customer:c,creation:{kind:'subscriber',ref:c.ref,name:c.name}});ok(made);c.id=made.bindings[0].id;
 ok(await s.act('open_customer',{customer:c}));
 const r=await s.act('allocate_auto',{customer:c,plan:p,creation:{kind:'sim',ref:'ue-000001',name:c.iccid}});
 assert.equal(r.error.code,'UNSAFE_AUTO_ASSIGN');assert.equal(s.app.db.sims.length,0);
});
test('pool reconciliation checks all statuses, rejects incomplete filtered inventory',async t=>{
 const s=await session(t);s.app.db.pool.push({iccid:'8901000000000000001'},{iccid:'8901000000000000002',assigned:true},{iccid:'8901000000000000003',failed:true});
 ok(await s.open('business_sim_pool'));
 ok(await s.check('Pool reconciliation','matched',{view:'business_sim_pool'}));
 ok(await s.send('web_interact',{view:'business_sim_pool',action:'filter',label:'Status',value:'Available'}));
 const r=await s.check('Pool reconciliation','matched',{view:'business_sim_pool'},500);assert.equal(r.status,'error');assert.equal(r.actual,'incomplete inventory');
});
