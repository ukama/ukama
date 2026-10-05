/* SPDX-License-Identifier: MPL-2.0
 * C runner to Chromium qualification against controlled, source-shaped DOM.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { readdir } from 'node:fs/promises';
import { run } from './browser-harness.mjs';
import { inventoryFixture } from '../../adapters/webapp/test/inventory-fixture.mjs';
import { commerceFixture } from '../../adapters/webapp/test/commerce-fixture.mjs';
import { operationsFixture } from '../../adapters/webapp/test/operations-fixture.mjs';
const choose=file=>file.includes('106-')?commerceFixture:file.includes('107-')?operationsFixture:inventoryFixture;
const files=(await readdir(new URL('../../scenarios/webapp/p0/inventory/',import.meta.url))).filter(f=>f.endsWith('.yaml')).sort();
for(const file of files)test(`inventory scenario: ${file}`,async t=>{
 const r=await run(t,'inventory/'+file,'',false,choose(file));
 assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');assert.equal(r.journal.cleanup,'complete');assert(r.report.checks.passed>0);
 assert.deepEqual(r.app.db.networks,[{id:'existing',name:'existing-network'}]);
 if(file.includes('site-membership')){const checks=r.report.results.filter(x=>x.label==='Site names');assert(checks.every(c=>Array.isArray(c.actual)&&c.actual.length===2));}
 if(file.includes('transient'))assert(r.report.results.some(x=>x.label==='Scope leaks'&&x.actual==='none'));
});
const short=text=>text.replaceAll('timeout_seconds: 30','timeout_seconds: 3').replaceAll('timeout_seconds: 15','timeout_seconds: 3');
for(const [file,mode,label] of [
 ['wb-090-site-membership.yaml','extra-site','Site names'],['wb-090-site-membership.yaml','duplicate-site','Site names'],
 ['wb-092-site-identity.yaml','wrong-location','Site location'],['wb-092-site-identity.yaml','wrong-coordinates','Site coordinates'],['wb-092-site-identity.yaml','wrong-site-node','Site node IDs'],
 ['wb-093-site-detail-network-switch.yaml','no-detail-redirect','Path'],['wb-094-node-detail-network-switch.yaml','no-detail-redirect','Path'],
 ['wb-095-site-transient-scope.yaml','transient-leak','Scope leaks'],
 ['wb-099-map-selection.yaml','wrong-map-site','Map selection'],['wb-099-map-selection.yaml','wrong-map-link','Path'],['wb-099-map-selection.yaml','wrong-pin-color','Map color'],
 ['wb-101-node-network-membership.yaml','foreign-node','Node IDs'],['wb-102-missing-home-data.yaml','invented-zero','Sites online'],
 ['wb-106-network-customer-usage.yaml','foreign-network-usage','Data volume'],['wb-107-site-distinct-states.yaml','conflated-state','Site status']
])test(`inventory rejects ${mode}: ${file}`,async t=>{
 const r=await run(t,'inventory/'+file,mode,false,choose(file),short);
 assert.notEqual(r.code,0,r.output);assert.equal(r.report.outcome,'FAIL');assert(r.report.results.some(x=>x.label===label&&x.state==='FAIL'),r.output);assert.equal(r.journal.cleanup,'complete');
});
