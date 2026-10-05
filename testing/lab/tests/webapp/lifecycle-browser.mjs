/* SPDX-License-Identifier: MPL-2.0
 * C/Chromium fixture qualification; never a live console result.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {readdir} from 'node:fs/promises';
import {run} from './browser-harness.mjs';
import {commerceFixture} from '../../adapters/webapp/test/commerce-fixture.mjs';
const files=(await readdir(new URL('../../scenarios/webapp/p0/lifecycle/',import.meta.url))).filter(f=>f.endsWith('.yaml')).sort();
for(const file of files)test(`lifecycle: ${file}`,async t=>{
 const r=await run(t,'lifecycle/'+file,'',false,commerceFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');
 assert.equal(r.app.db.plans.length,0);assert.equal(r.app.db.subscribers.length,0);assert.equal(r.app.db.sims.length,0);
 if(file.includes('rename'))assert.equal(r.app.db.operations.filter(o=>o.op==='updatePackage').length,1);
 if(file.includes('auto-assignment'))assert.equal(r.app.db.operations.filter(o=>o.op==='allocateSim').length,1);
});
const short=s=>s.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 3');
for(const [file,mode,label] of [
 ['wb-120-rename-plan.yaml','rename-terms-drift','Plan price'],
 ['wb-122-pool-allocation-choices.yaml','wrong-pool-count','Pool reconciliation'],
 ['wb-122-pool-allocation-choices.yaml','allocated-option-leak','SIM option present'],
 ['wb-123-inventory-read-failure.yaml','pool-fail-empty','Text visible'],
 ['wb-124-name-pending.yaml','name-fail-open','Button enabled'],
 ['wb-125-name-failure.yaml','name-fail-open','Button enabled'],
])test(`lifecycle rejects ${mode}: ${file}`,async t=>{
 const r=await run(t,'lifecycle/'+file,mode,false,commerceFixture,short);
 assert.notEqual(r.code,0,r.output);assert(r.report.results.some(c=>c.label===label&&c.state==='FAIL'),r.output);assert.equal(r.journal.cleanup,'complete');
});
test('cross-lens plan mismatch is rejected despite correct customer identity',async t=>{
 const r=await run(t,'expanded/wb-042-customer-cross-lens.yaml','wrong-customer-plan',false,commerceFixture,short);
 assert.notEqual(r.code,0,r.output);assert(r.report.results.some(c=>c.label==='Customer plan'&&c.actual==='Foreign plan'));assert.equal(r.journal.cleanup,'complete');
});
