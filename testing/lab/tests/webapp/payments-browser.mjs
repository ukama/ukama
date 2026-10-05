/* SPDX-License-Identifier: MPL-2.0 — fixtures, not live product coverage. */
import test from 'node:test';
import assert from 'node:assert/strict';
import {readFile,stat} from 'node:fs/promises';
import {run} from './browser-harness.mjs';
import {commerceFixture} from '../../adapters/webapp/test/commerce-fixture.mjs';
for(const file of ['wb-130-downloaded-receipt.yaml','wb-131-payment-rejection-retry.yaml','wb-132-rapid-topup.yaml','wb-133-no-sim-payment-guard.yaml'])test(`payment: ${file}`,async t=>{
 const r=await run(t,'payments/'+file,'',false,commerceFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');
 assert.equal(r.app.db.payments.length,file.includes('no-sim')?0:1);
 if(file.includes('rejection'))assert.deepEqual(r.journal.creation_intents.filter(i=>i.kind==='payment').map(i=>i.state),['not_submitted','identified']);
 if(file.includes('downloaded')) {
   const step=r.commands.find(c=>c.direction==='request'&&c.message.action==='web_commerce'&&c.message.inputs.action==='download_receipt');assert(step);
   const artifact=r.report.results.flatMap(row=>row.artifacts??[]).find(p=>/receipt-\d+\.pdf$/.test(p));assert(artifact,JSON.stringify(r.report));
   assert.equal((await stat(artifact)).mode&0o777,0o600);assert((await readFile(artifact)).subarray(0,5).equals(Buffer.from('%PDF-')));
 }
});
const short=s=>s.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 3');
for(const mode of ['wrong-pdf-total','wrong-pdf-id','html-download','old-receipt-date'])test(`receipt rejects ${mode}`,async t=>{
 const r=await run(t,'payments/wb-130-downloaded-receipt.yaml',mode,false,commerceFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');
 assert.equal(r.app.db.payments.length,1);
 assert.match(r.output,mode==='old-receipt-date'?/outside submission window/:mode==='html-download'?/not a PDF/:/receipt differs|different fields/);
});
test('rejected payment cannot hide its error and still claim success',async t=>{
 const r=await run(t,'payments/wb-131-payment-rejection-retry.yaml','hidden-payment-error',false,commerceFixture,short);
 assert.notEqual(r.code,0);assert.equal(r.app.db.payments.length,0);assert.equal(r.journal.cleanup,'complete');
});
test('rapid duplicate payment IDs fail and retain unresolved ownership for reconciliation',async t=>{
 const r=await run(t,'payments/wb-132-rapid-topup.yaml','double-payment',false,commerceFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.app.db.payments.length,2);
 assert.equal(r.journal.uncertain_creation,true);assert.equal(r.journal.cleanup,'failed');
 assert.match(r.output,/Multiple IDs observed/);
});
