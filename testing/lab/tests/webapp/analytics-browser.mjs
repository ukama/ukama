/* SPDX-License-Identifier: MPL-2.0 — controlled fixture integration only. */
import test from 'node:test';
import assert from 'node:assert/strict';
import {run} from './browser-harness.mjs';
import {analyticsFixture} from '../../adapters/webapp/test/analytics-fixture.mjs';
import {commerceFixture} from '../../adapters/webapp/test/commerce-fixture.mjs';
for(const file of ['wb-140-node-metric-observations.yaml','wb-141-site-metric-observations.yaml','wb-142-app-resource-identity.yaml','wb-143-revenue-trend-observation.yaml','wb-144-purchase-distribution.yaml','wb-145-zero-revenue.yaml'])test(`analytics: ${file}`,async t=>{
 const r=await run(t,'analytics/'+file,'',false,/14[45]/.test(file)?commerceFixture:analyticsFixture);
 assert.equal(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');
});
const short=s=>s.replace('scenario_timeout_seconds: 40','scenario_timeout_seconds: 15').replaceAll('timeout_seconds: 30','timeout_seconds: 2');
for(const mode of ['wrong-unit','zero-gap','wrong-time','wrong-tooltip','wrong-app'])test(`analytics rejects ${mode}`,async t=>{
 const r=await run(t,'analytics/'+(mode==='wrong-app'?'wb-142-app-resource-identity.yaml':'wb-140-node-metric-observations.yaml'),mode,false,analyticsFixture,short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.journal.cleanup,'complete');assert.match(r.output,/Visible UI state|Command deadline/);
});
