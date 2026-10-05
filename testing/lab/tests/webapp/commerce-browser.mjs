/* SPDX-License-Identifier: MPL-2.0
 * Full C lifecycle + Chromium controlled fixture; not live product coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { run } from './browser-harness.mjs';
import { commerceFixture } from '../../adapters/webapp/test/commerce-fixture.mjs';
const scenario='commerce/wb-020-plan-unit-roundtrip.yaml';
const short=text=>text.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 12').replaceAll('timeout_seconds: 30','timeout_seconds: 3').replaceAll('timeout_seconds: 15','timeout_seconds: 3');
test('C resolves plans, performs UI creation, checks all units and tears down exact owned IDs',async t=>{
 const r=await run(t,scenario,'',false,commerceFixture);assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.report.checks.passed,8);
 assert.equal(r.journal.resources.length,4);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));assert.equal(r.app.db.plans.length,0);assert.equal(r.app.db.networks.length,1);
 const creates=r.app.db.operations.filter(x=>x.op==='addPackage');assert.deepEqual(creates.map(x=>x.data.duration),[1440,10080,43200]);assert.deepEqual(creates.map(x=>[x.data.dataVolume,x.data.dataUnit]),[[512,'MB'],[1,'GB'],[2,'GB']]);
});
test('wrong duration returned to the app fails C acceptance and still cleans owned plans',async t=>{
 const r=await run(t,scenario,'wrong-minutes',false,commerceFixture,short);assert.notEqual(r.code,0);assert.equal(r.report.outcome,'FAIL');assert.equal(r.journal.cleanup,'complete');assert.equal(r.app.db.plans.length,0);
 assert(r.report.results.some(x => x.label === 'Plan terms' && /minutes/.test(x.actual)), 'failure must reach the incorrect visible validity');
});
test('late plan UI recovers its receipt and deletes the identified resource exactly once',async t=>{
 const r=await run(t,scenario,'late-plan',false,commerceFixture,short);assert.notEqual(r.code,0);assert.equal(r.journal.cleanup,'complete');assert.equal(r.app.db.operations.filter(x=>x.op==='addPackage').length,1);assert.equal(r.app.db.operations.filter(x=>x.op==='deletePackage').length,1);assert.equal(r.journal.uncertain_creation,false);
});
test('unidentified plan mutation fails cleanup without guessing resource IDs or retrying',async t=>{
 const r=await run(t,scenario,'opaque-plan',false,commerceFixture,short);assert.notEqual(r.code,0);assert.equal(r.journal.uncertain_creation,true);assert.equal(r.journal.cleanup,'failed');assert.equal(r.app.db.operations.filter(x=>x.op==='addPackage').length,1);assert.equal(r.app.db.operations.filter(x=>x.op==='deletePackage').length,0);
});
for(const [file,checks] of [['wb-021-customer-sim-allocation.yaml',10],['wb-022-cash-topup-receipt.yaml',34],['wb-027-customer-commerce-journey.yaml',36],['wb-023-sim-service-toggle.yaml',8],['wb-024-ue-usage.yaml',8],['wb-025-daily-allocation.yaml',10],['wb-026-monthly-allocation.yaml',10]]) {
 test(`C commerce lifecycle: ${file}`,async t=>{
  const r=await run(t,'commerce/'+file,'',false,commerceFixture);assert.equal(r.code,0,r.output);assert.equal(r.report.checks.passed,checks);assert.equal(r.journal.cleanup,'complete');assert.equal(r.app.db.sims.length,0);assert.equal(r.app.db.subscribers.length,0);assert.equal(r.app.db.plans.length,0);
  assert.equal(r.journal.sim_inventory.state,'visible');assert.equal(r.journal.sim_inventory.cleanup,'retained_factory_pool');
  if(file.includes('topup') || file.includes('journey')){assert.equal(r.app.db.payments.length,1);assert.equal(r.journal.resources.find(x=>x.kind==='payment').cleanup,'retained_ledger');}
  if(file.includes('usage') || file.includes('journey')){assert(r.app.db.operations.some(x=>x.runtime==='start-ue.sh'));assert(r.app.db.operations.some(x=>x.runtime==='traffic.sh'));assert(r.app.db.operations.some(x=>x.runtime==='cleanup-ue.sh'));}
 });
}
for(const [mode,file,kind] of [['late-sim','wb-021-customer-sim-allocation.yaml','sim'],['late-payment','wb-022-cash-topup-receipt.yaml','payment'],['opaque-payment','wb-022-cash-topup-receipt.yaml','payment']]) {
 test(`C recovers commerce mutation failure: ${mode}`,async t=>{
  const r=await run(t,'commerce/'+file,mode,false,commerceFixture,short);assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op===(kind==='sim'?'allocateSim':'addPayment')).length,1);
  if(mode==='opaque-payment'){assert.equal(r.journal.uncertain_creation,true);assert.equal(r.journal.cleanup,'failed');assert.equal(r.app.db.operations.filter(x=>x.op==='deleteSim').length,0);}
  else {assert.equal(r.journal.cleanup,'complete');assert.equal(r.journal.uncertain_creation,false);assert.equal(r.journal.resources.find(x=>x.kind===kind).cleanup,kind==='sim'?'deleted':'retained_ledger');assert.equal(r.app.db.sims.length,0);}
 });
}

for (const mode of ['wrong-revenue', 'wrong-plan-revenue']) {
 test(`C rejects incorrect Business totals after a valid payment: ${mode}`, async t => {
  const r = await run(t, 'commerce/wb-022-cash-topup-receipt.yaml', mode, false, commerceFixture,
    text => short(text).replaceAll('timeout_seconds: 15', 'timeout_seconds: 3'));
  assert.notEqual(r.code, 0, r.output);
  assert.equal(r.report.outcome, 'FAIL');
  assert.equal(r.app.db.payments.length, 1);
  assert.equal(r.journal.cleanup, 'complete');
  assert.equal(r.app.db.sims.length, 0);
  assert.equal(r.report.checks.failed, 1);
  assert(r.report.results.some(x => x.label === (mode === 'wrong-revenue' ? 'Revenue' : 'Performance revenue') && x.expected === '$25' && x.actual === '$26'));
 });
}
