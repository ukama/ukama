/* SPDX-License-Identifier: MPL-2.0 */
import test, {before,after} from 'node:test';
import assert from 'node:assert/strict';
import http from 'node:http';
import {once} from 'node:events';
import {chromium} from 'playwright';
import {Budget} from '../dist/contract.js';
import {statusFault, faultState, trackStatusRequests, watchOperation, observation} from '../dist/operation-observer.js';
let browser,server,origin;
before(async()=>{
  server=http.createServer((req,res)=>{
    res.setHeader('access-control-allow-origin',req.headers.origin??'*');res.setHeader('access-control-allow-credentials','true');res.setHeader('access-control-allow-headers','content-type');
    res.end(req.url==='/graphql'?JSON.stringify({data:{forwarded:true}}):'<main class="main"></main>');
  });server.listen(0,'0.0.0.0');await once(server,'listening');origin=`http://127.0.0.1:${server.address().port}`;
  browser=await chromium.launch({headless:true,...(process.env.ULAB_WEBAPP_EXECUTABLE_PATH?{executablePath:process.env.ULAB_WEBAPP_EXECUTABLE_PATH}:{})});
});
after(async()=>{await browser?.close();await new Promise(r=>server.close(r));});
const body=(id='owned')=>({operationName:'GetNodeOperationStatus',query:'query GetNodeOperationStatus($nodeId: String!) { getNodeOperationStatus(data: { nodeId: $nodeId }) { busy } }',variables:{nodeId:id}});
async function pageFor(t){const context=await browser.newContext();t.after(()=>context.close());const page=await context.newPage();trackStatusRequests(page);await page.goto(origin);return page;}
const request=(page,payload=body(),endpoint=origin)=>page.evaluate(async ({payload,endpoint})=>fetch(endpoint+'/graphql',{method:'POST',headers:{'content-type':'application/json'},credentials:'include',body:JSON.stringify(payload)}).then(r=>r.json()),{payload,endpoint});
const arm=(page,value)=>statusFault(page,false,{entity:{id:'owned'},value},new Budget(1000));
test('scoped read fault intercepts the observed cross-origin gateway, never other entities or mutations',async t=>{
  const page=await pageFor(t),gateway=origin.replace('127.0.0.1','localhost');
  await request(page,body(),gateway);await arm(page,'read_error');assert.equal(faultState(page),'read_error:armed');
  assert((await request(page,body(),gateway)).errors);assert.equal(faultState(page),'read_error:applied');
  for(const payload of [body('foreign'),[body()],{...body(),query:'mutation GetNodeOperationStatus($nodeId:String!) { restartNode { success } }'},{...body(),operationName:'RestartNode'}])
    assert.equal((await request(page,payload,gateway)).data.forwarded,true);
  assert.equal((await request(page,body(),origin)).data.forwarded,true,'same operation at another origin is not intercepted');
  await arm(page,'none');assert.equal(faultState(page),'none');assert.equal((await request(page,body(),gateway)).data.forwarded,true);
});
test('unseen or ambiguous gateway cannot arm a fault',async t=>{
  const page=await pageFor(t);
  await assert.rejects(statusFault(page,false,{entity:{id:'owned'},value:'idle'},new Budget(100)),/No operation-status read|deadline/i);
  await request(page);await request(page,body(),origin.replace('127.0.0.1','localhost'));
  await assert.rejects(arm(page,'read_error'),/multiple endpoints/);
  assert.equal(faultState(page),'none');
});
test('clearing and rearming faults resets application proof and restores actual reads',async t=>{
  const page=await pageFor(t);await request(page);await arm(page,'idle');assert.equal((await request(page)).data.getNodeOperationStatus.busy,false);
  await arm(page,'none');assert.equal((await request(page)).data.forwarded,true);
  await arm(page,'read_error');assert.equal(faultState(page),'read_error:armed');await request(page);assert.equal(faultState(page),'read_error:applied');
});
test('navigation invalidates old fault application proof and other detail paths pass through',async t=>{
  const page=await pageFor(t);await request(page);await arm(page,'read_error');await request(page);assert.equal(faultState(page),'read_error:applied');
  await page.goto(origin+'/other-detail');assert.equal(faultState(page),'read_error:wrong_scope');assert.equal((await request(page)).data.forwarded,true);
  await page.goto(origin);assert.equal(faultState(page),'read_error:armed');
  assert((await request(page)).errors);assert.equal(faultState(page),'read_error:applied');
});
const trio=['tower','amplifier','controller'].map(type=>({id:type,name:type,type}));
async function site(t){const page=await pageFor(t);await page.setContent(`<main class="main"><div class="pagehead"><button>Site actions</button></div>${trio.map(n=>`<div class="app-card" id="${n.id}"><div><span></span><span>${n.id}</span><span>is online and well</span></div></div>`).join('')}</main>`);return page;}
const setStatus=(page,type,status)=>page.locator(`#${type} > div > span:last-child`).evaluate((e,text)=>e.textContent=text,status);
test('a transient rendered controller outage remains failed after it returns online',async t=>{
  const page=await site(t);await watchOperation(page,'restart',{nodes:trio});
  await page.getByRole('button').evaluate(e=>e.textContent='Site actions • busy');
  await setStatus(page,'controller','is offline');await setStatus(page,'controller','is online and well');
  assert.match(await observation(page,'restart'),/failed: Controller/);
});
test('restart completion requires both offline transitions and recovery with an online controller',async t=>{
  const page=await site(t);await watchOperation(page,'restart',{nodes:trio});
  await page.getByRole('button').evaluate(e=>e.textContent='Site actions • busy');
  await setStatus(page,'tower','is offline');await setStatus(page,'amplifier','is offline');
  await page.getByRole('button').evaluate(e=>e.textContent='Site actions');
  assert.equal(await observation(page,'restart'),'observing');
  await setStatus(page,'tower','is online and well');assert.equal(await observation(page,'restart'),'observing');
  await setStatus(page,'amplifier','is online and well');assert.equal(await observation(page,'restart'),'complete');
});
test('controller observation continues after lock release while tower recovery is pending',async t=>{
  const page=await site(t);await watchOperation(page,'restart',{nodes:trio});
  await page.getByRole('button').evaluate(e=>e.textContent='Site actions • busy');
  await setStatus(page,'tower','is offline');await setStatus(page,'amplifier','is offline');
  await page.getByRole('button').evaluate(e=>e.textContent='Site actions');
  await setStatus(page,'controller','is offline');await setStatus(page,'controller','is online and well');
  await setStatus(page,'tower','is online and well');await setStatus(page,'amplifier','is online and well');
  assert.match(await observation(page,'restart'),/failed: Controller/);
});
test('an idle-only snapshot and navigation never pass a restart watch',async t=>{
  const page=await site(t);await watchOperation(page,'restart',{nodes:trio});assert.equal(await observation(page,'restart'),'armed');
  await page.reload();assert.equal(await observation(page,'restart'),'not armed');
});
test('missing and duplicate trio identities fail before arming',async t=>{
  const page=await site(t);await assert.rejects(watchOperation(page,'restart',{nodes:[trio[0],trio[0],trio[2]]}),/unique/);
  await page.locator('#controller').evaluate(e=>e.remove());
  await assert.rejects(watchOperation(page,'restart',{nodes:trio}),/Missing or ambiguous controller/);
});
test('early optimistic unlock is sticky even after the control becomes busy again',async t=>{
  const page=await pageFor(t);await page.setContent('<main class="main"><div class="pagehead"><button>Restart node</button></div></main>');
  await watchOperation(page,'timeout',{});
  await page.getByRole('button').evaluate(e=>{e.textContent='Restarting…';e.disabled=true;});
  await page.getByRole('button').evaluate(e=>{e.textContent='Restart node';e.disabled=false;});
  await page.getByRole('button').evaluate(e=>{e.textContent='Restarting…';e.disabled=true;});
  assert.match(await observation(page,'timeout'),/failed: Optimistic lock released/);
});
