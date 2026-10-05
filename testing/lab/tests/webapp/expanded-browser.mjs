/* SPDX-License-Identifier: MPL-2.0
 * Controlled C/Chromium integration; never live product coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { readdir } from 'node:fs/promises';
import http from 'node:http';
import { once } from 'node:events';
import { run } from './browser-harness.mjs';
import { commerceFixture } from '../../adapters/webapp/test/commerce-fixture.mjs';
import { provisioningFixture } from '../../adapters/webapp/test/provisioning-fixture.mjs';
async function sessionFixture(mode='') {
 const auth=http.createServer((req,res)=>res.end('<h1>Sign in</h1>'));auth.listen(0,'127.0.0.1');await once(auth,'listening');
 const authOrigin=`http://127.0.0.1:${auth.address().port}`;
 const app=await provisioningFixture(mode,{get(req,res){if(req.url==='/'&&!req.headers.cookie?.includes('session=owner')&&mode!=='session-leak'){res.writeHead(302,{Location:authOrigin});res.end();return true}return false}});
 const close=app.close;app.close=async()=>{await close();await new Promise(r=>auth.close(r))};app.authOrigin=authOrigin;
 app.state.cookies=[{name:'session',value:'owner',domain:'127.0.0.1',path:'/',expires:-1,httpOnly:true,secure:false,sameSite:'Lax'}];return app;
}
const files=(await readdir(new URL('../../scenarios/webapp/p0/expanded/',import.meta.url))).filter(f=>f.endsWith('.yaml')).sort();
for(const file of files)test(`expanded scenario: ${file}`,async t=>{
 const r=await run(t,'expanded/'+file,'',false,file.includes('session')?sessionFixture:commerceFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');assert(r.report.checks.passed>0);
 if(/duplicate|invalid-plan|keyboard/.test(file))assert.equal(r.app.db.operations.filter(o=>o.op==='addPackage').length,file.includes('duplicate')?1:0);
 if(/invalid-customer|missing-plan|empty-customers/.test(file))assert.equal(r.app.db.operations.filter(o=>o.op==='addSubscriber').length,0);
 if(file.includes('topup-requires'))assert.equal(r.app.db.operations.filter(o=>o.op==='addPayment').length,0);
});
const short=text=>text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 3').replaceAll('timeout_seconds: 15','timeout_seconds: 3');
for(const [file,mode,label] of [['wb-032-network-plan-isolation.yaml','scope-leak','Plan option present'],['wb-031-invalid-plan.yaml','validation-broken','Button enabled'],['wb-037-keyboard-dialog.yaml','lost-focus','Focus on button'],['wb-040-session-clearing.yaml','session-leak','Auth origin']])test(`expanded rejects ${mode}`,async t=>{
 const r=await run(t,'expanded/'+file,mode,false,file.includes('session')?sessionFixture:commerceFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.report.outcome,'FAIL');assert(r.report.results.some(x=>x.label===label&&x.state==='FAIL'),r.output);
});
