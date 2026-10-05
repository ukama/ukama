/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { run } from './browser-harness.mjs';
import { onboardingFixture } from '../../adapters/webapp/test/onboarding-fixture.mjs';
const cases=['070-network-create','071-network-validation','072-existing-selection','073-site-persistence','074-progress','075-missing-amplifier','076-offline-tower','077-unready-controller','078-missing-location','079-transport-once','080-rejection-retry','081-sim-guidance'];
const execute=(t,name,mode='',transform=text=>text)=>run(t,'onboarding/wb-'+name+'.yaml',mode,false,onboardingFixture,transform);
for(const name of cases)test('onboarding '+name+' executes staged C/browser assertions and cleans owned resources',async t=>{
 const r=await execute(t,name);assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');assert.equal(!!r.journal.uncertain_creation,false);assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
 assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
 if(name==='079-transport-once'||name==='080-rejection-retry')assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,1);
});
for(const[mode,name]of[['wrong-tower','073-site-persistence'],['wrong-coordinates','073-site-persistence'],['wrong-component','073-site-persistence'],['missing-confirm','074-progress'],['duplicate-site','079-transport-once'],['premature','075-missing-amplifier']])test('onboarding rejects '+mode+' instead of repairing or masking it',async t=>{
 const r=await execute(t,name,mode,text=>text.replaceAll('timeout_seconds: 15','timeout_seconds: 3'));assert.notEqual(r.code,0);assert.equal(r.report.outcome,'FAIL');
 if(mode==='wrong-tower'){assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);assert.equal(r.journal.cleanup,'complete')}
 if(mode!=='duplicate-site'){assert.equal(r.journal.cleanup,'complete');assert.equal(!!r.journal.uncertain_creation,false)}
});
test('onboarding never retries a site without the adapter controlled rejection',async t=>{
 const r=await execute(t,'074-progress','',text=>text.replace('action: submit_site','action: retry_site'));assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addSite').length,0);assert.equal(r.journal.cleanup,'complete');
});

test('invalid-name validator regression cannot leak an unowned network',async t=>{
 const r=await execute(t,'071-network-validation','broken-validation');assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.op==='addNetwork').length,0);assert.equal(r.journal.cleanup,'complete');
});
