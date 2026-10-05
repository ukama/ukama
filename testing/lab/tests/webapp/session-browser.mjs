/* SPDX-License-Identifier: MPL-2.0 -- controlled fixtures, never live coverage. */
import test from 'node:test';
import assert from 'node:assert/strict';
import { readdir } from 'node:fs/promises';
import { run } from './browser-harness.mjs';
import { sessionFixture } from '../../adapters/webapp/test/session-fixture.mjs';
const files=(await readdir(new URL('../../scenarios/webapp/p0/auth/',import.meta.url))).filter(f=>f.endsWith('.yaml')).sort();
for(const file of files)test(`auth C/Chromium: ${file}`,async t=>{
 const r=await run(t,'auth/'+file,'',false,sessionFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');assert(r.report.checks.passed>0);
 assert.equal(r.journal.resources.length,0);assert(r.commands.filter(c=>c.direction==='request').every(c=>['init','web_session','web_session_equals','close'].includes(c.message.action)));
 if(file.includes('recover-'))assert(r.app.db.mints>0);
 if(file.includes('reject-token'))assert.equal(r.app.db.refresh,1);
 if(file.includes('logout'))assert.equal(r.app.db.logout,1);
 if(file.includes('welcome')){assert.equal(r.app.db.welcome,1);assert.equal(r.app.accounts.WELCOME.welcome,false)}
});
const short=text=>text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 3');
for(const [prefix,mode,label] of [['050','wrong-identity','Settings field'],['050','wrong-org','Organization'],['060','wrong-role','Settings field'],['051','recovery-broken','Surface'],['053','stale-access','Refresh observed'],['054','unauthorized-leak','Dashboard visible'],['056','logout-leak','Access blocked'],['057','stale-access','Access blocked'],['058','welcome-not-persisted','Surface'],['058','welcome-error','Surface'],['065','policy-leak','Control state'],['066','wrong-404','Surface'],['066','billing-link','Nav visible'],['064','auth-error','Surface']])test(`auth rejects ${mode} (${prefix})`,async t=>{
 const file=files.find(f=>f.startsWith('wb-'+prefix+'-'));const r=await run(t,'auth/'+file,mode,false,sessionFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.report.outcome,'FAIL');assert(r.report.results.some(x=>(mode==='auth-error'||x.label===label)&&x.state==='FAIL'),r.output);
});
