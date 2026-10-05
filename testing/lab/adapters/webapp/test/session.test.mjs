/* SPDX-License-Identifier: MPL-2.0 -- controlled lab acceptance tests only. */
import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp,rm,readFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { Worker } from '../dist/worker.js';
import { profile } from '../dist/contract.js';
import { sessionFixture } from './session-fixture.mjs';
async function session(t,options={},mode=''){
 const app=await sessionFixture(mode),dir=await mkdtemp(join(tmpdir(),'ulab-auth-'));const env=await app.prepare(dir),worker=new Worker();let id=0;
 t.after(async()=>{await worker.shutdown(worker.failed);await app.close();await rm(dir,{recursive:true,force:true})});
 const send=(action,inputs={},ms=4000)=>worker.handle({protocol:1,run_id:'auth-fixture',command_id:++id,action,deadline_ms:Date.now()+ms,inputs});
 const config={base_url:app.origin,auth_state:env.ULAB_AUTH_OWNER_STATE,session_mode:'auth_test',auth_origin:app.authOrigin,headless:true,...options};
 const result=await send('init',{profile:config,artifacts_dir:dir});return{app,dir,env,worker,send,result,config};
}
const ok=r=>assert.equal(r.status,'ok',JSON.stringify(r));
const check=(s,label,expected,subject,ms)=>s.send('web_session_equals',{view:'session',label,expected,requirement:'WEB-AUTH-002',...(subject?{subject}:{})},ms);
const event=(s,action,value)=>s.send('web_session',{view:'session',action,...(value?{value}:{})});
test('auth mode initializes an isolated unverified context without opening a protected page',async t=>{
 const s=await session(t);ok(s.result);assert.deepEqual(s.result.actual,{initialized:true,authenticated:false,session_mode:'auth_test',browser:'chromium',browser_version:s.result.actual.browser_version});assert.equal(s.app.db.documents,0);
 ok(await event(s,'navigate','/business/settings'));ok(await check(s,'Settings field','owner test','Full name'));ok(await s.send('close'));
});
test('empty context is permitted only in explicit auth_test mode',async t=>{
 const s=await session(t,{auth_state:'none'});ok(s.result);ok(await event(s,'navigate','/business/settings'));ok(await check(s,'Surface','auth'));ok(await check(s,'Dashboard visible','false'));
 assert.throws(()=>profile({base_url:s.app.origin,auth_state:'none'}));
});
test('auth origins and session mode are strict and do not change default authentication',()=>{
 for(const x of [{session_mode:'unknown'},{auth_origin:'http://auth.test'},{session_mode:'auth_test'},{session_mode:'auth_test',auth_origin:'http://console.test'},{session_mode:'auth_test',auth_origin:'https://name:password@auth.test'},{session_mode:'auth_test',auth_origin:'http://auth.test/path'}])assert.throws(()=>profile({base_url:'http://console.test',auth_state:'.auth/owner.json',...x}));
 assert.equal(profile({base_url:'http://console.test',auth_state:'.auth/owner.json'}).session_mode,'authenticated');
});
test('auth mode rejects ordinary navigation and resource mutations',async t=>{
 const s=await session(t);ok(s.result);const r=await s.send('web_open',{view:'business_members'});assert.equal(r.error.code,'SESSION_MODE');assert.equal(s.app.db.documents,0);
});
test('authenticated default mode rejects session commands',async t=>{
 const s=await session(t,{session_mode:'authenticated',auth_origin:undefined});ok(s.result);assert.equal((await event(s,'drop_token')).error.code,'SESSION_MODE');
});
for(const action of ['drop_token','invalidate_token','reject_token','expire_token'])test(`${action} changes only context state and never records raw tokens`,async t=>{
 const s=await session(t);ok(s.result);const before=await readFile(s.config.auth_state,'utf8');ok(await event(s,action));ok(await event(s,'navigate','/business/settings'));ok(await check(s,'Surface','dashboard'));if(action==='reject_token')ok(await event(s,'navigate','/business/settings'));ok(await check(s,'Settings field','owner test','Full name'));ok(await s.send('close'));assert.equal(await readFile(s.config.auth_state,'utf8'),before);assert(s.app.db.mints>0);
 const results=await readFile(join(s.dir,'auth-fixture','results.jsonl'),'utf8');for(const cookie of s.app.state.cookies)assert(!results.includes(cookie.value+'"'),'cookie value is not command evidence');
});
test('token fault with no real session fails its precondition',async t=>{
 const s=await session(t,{auth_state:'none'});ok(s.result);assert.equal((await event(s,'drop_token')).error.code,'AUTH_PRECONDITION');
});
test('off-origin and arbitrary routes cannot be injected',async t=>{
 const s=await session(t);ok(s.result);assert.equal((await event(s,'navigate','https://other.test/')).error.code,'INVALID_INPUT');assert.equal(s.app.db.documents,0);
});
test('welcome acknowledgement cannot execute outside the console',async t=>{
 const s=await session(t,{},'welcome-error');ok(s.result);
 // The success/failure welcome scenarios run the real C/Chromium path; verify
 // unknown commands here cannot turn an error page into a synthetic success.
 assert.equal((await event(s,'ack_welcome')).error.code,'WRONG_ORIGIN');
});
