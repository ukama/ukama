/* SPDX-License-Identifier: MPL-2.0 — controlled fixture integration only. */
import test from 'node:test';
import assert from 'node:assert/strict';
import {run} from './browser-harness.mjs';
import {teamSupportFixture} from '../../adapters/webapp/test/team-support-fixture.mjs';
const cases=['wb-150-member-rows.yaml','wb-151-invitation-validation.yaml','wb-152-invitation-role-rejection.yaml','wb-153-map-preference-persistence.yaml','wb-154-customer-support-lookup.yaml','wb-155-network-support-lookup.yaml','wb-156-support-restart-policy.yaml','wb-157-support-no-match-read-failure.yaml','wb-158-existing-pending-invitation.yaml'];
for(const file of cases)test(`team/support: ${file}`,async t=>{
 const r=await run(t,'team-support/'+file,'',false,teamSupportFixture);assert.equal(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');assert.equal(r.app.db.interceptedWrites.length,0,'probes must never reach the server');
});
const short=s=>s.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 2');
for(const [n,mode] of [[150,'duplicate-member'],[150,'wrong-role'],[151,'invalid-invite-enabled'],[152,'wrong-invite-role'],[152,'hidden-invite-error'],[153,'lost-preference'],[154,'wrong-iccid'],[154,'stale-summary'],[155,'wrong-support-id'],[156,'unguarded-restart'],[157,'hidden-support-error']])test(`team/support rejects ${mode}`,async t=>{
 const r=await run(t,'team-support/'+cases.find(f=>f.includes(String(n))),mode,false,teamSupportFixture,short);assert.notEqual(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');assert.equal(r.app.db.interceptedWrites.length,0,'no invitation/restart may escape the guard');
});
test('support restart probe observes a disabled control without clicking it',async t=>{
 const r=await run(t,'team-support/'+cases[6],'locked-restart',false,teamSupportFixture,s=>s.replace('expected: "confirmation"','expected: "locked"'));assert.equal(r.code,0,r.output);assert.equal(r.app.db.interceptedWrites.length,0);
});
