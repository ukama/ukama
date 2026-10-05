/* SPDX-License-Identifier: MPL-2.0 — controlled fixture, no live acceptance. */
import test from 'node:test';
import assert from 'node:assert/strict';
import {run} from './browser-harness.mjs';
import {teamSupportFixture} from '../../adapters/webapp/test/team-support-fixture.mjs';
const cases=['wb-160-member-loading-retry.yaml','wb-161-dialog-focus-cycle.yaml','wb-162-responsive-dialog.yaml','wb-163-persistent-plan-state.yaml','wb-164-browser-health.yaml'];
for(const file of cases)test(`UI quality: ${file}`,async t=>{
 const r=await run(t,'ui-quality/'+file,'',false,teamSupportFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');
});
const short=s=>s.replaceAll('timeout_seconds: 30','timeout_seconds: 2');
for(const [file,mode] of [[cases[0],'broken-retry'],[cases[1],'untrapped'],[cases[1],'lost-focus'],[cases[2],'wide-dialog'],[cases[4],'page-error']])test(`UI quality rejects ${mode}`,async t=>{
 const r=await run(t,'ui-quality/'+file,mode,false,teamSupportFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.report.outcome,'FAIL');assert(r.report.results.some(x=>x.state==='FAIL'));
});
