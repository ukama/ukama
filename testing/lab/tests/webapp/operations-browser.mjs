/* SPDX-License-Identifier: MPL-2.0
 * Full C CLI + Playwright against a controlled operation fixture, not product coverage.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import { operationsFixture } from '../../adapters/webapp/test/operations-fixture.mjs';
import { run } from './browser-harness.mjs';
const execute=(t,file,mode='',transform)=>run(t,'operations/'+file,mode,false,operationsFixture,transform);
const files=['wb-010-node-restart.yaml','wb-011-site-operations.yaml','wb-012-software-update.yaml','wb-013-node-and-site-locks.yaml','wb-014-software-failure-retry.yaml','wb-015-stale-confirmation.yaml','wb-016-offline-recovery.yaml'];
for(const file of files)test(`C/browser operation journey: ${file}`,async t=>{
  const r=await execute(t,file);assert.equal(r.code,0,r.output);assert.equal(r.report.outcome,'PASS');
  assert.equal(r.journal.cleanup,'complete');assert(r.journal.resources.every(x=>x.cleanup==='deleted'));
  const mutations=r.app.db.operations.filter(x=>x.ui);
  if(file.includes('node-restart'))assert.equal(mutations.filter(x=>x.ui==='restart-node').length,1,'cancel must not mutate; confirm must submit once');
  if(file.includes('site-operations'))assert.deepEqual(mutations.map(x=>x.ui),['radio','radio','service','service','restart-site']);
  if(file.includes('software-failure'))assert.equal(mutations.filter(x=>x.ui==='update').length,2,'only the explicit retry may resubmit');
  if(file.includes('stale-confirmation')){assert.equal(mutations.filter(x=>x.ui==='restart-node').length,0);assert.equal(mutations.filter(x=>x.ui==='update').length,1);}
});
test('wrong visible target fails before any software mutation',async t=>{
  const r=await execute(t,files[2],'',text=>text.replace('tag: "${ULAB_SOFTWARE_TARGET_VERSION}"','tag: "unpublished-tag"'));
  assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.ui==='update').length,0);assert.equal(r.journal.cleanup,'complete');assert.match(r.output,/Visible software target differs/);
});
test('wrong terminal version fails even when the UI says Up to date',async t=>{
  const r=await execute(t,files[2],'wrong-terminal-version',text=>text.replace(/(label: "Current version"\n        expected: "\$\{ULAB_SOFTWARE_TARGET_VERSION\}"[\s\S]*?timeout_seconds:) 15/, '$1 3'));
  assert.notEqual(r.code,0);assert.equal(r.app.db.operations.filter(x=>x.ui==='update').length,1);assert.equal(r.journal.cleanup,'complete');
  assert(r.report.results.some(x=>x.label==='Current version'&&x.actual==='wrong-version'), 'failure must observe the wrong installed version');
});
test('operation-status error renders unavailable actions and a visible reason',async t=>{
  const r=await execute(t,files[6],'status-error',text=>text.slice(0,text.indexOf('phases:'))+`phases:
  - name: unknown_status
    events:
      - type: web_open
        view: network_node_detail
        networks: net-001
        nodes: tower-site-001-001
    checks:
      - type: web_action_available
        view: network_node_detail
        label: Restart node
        available: false
        requirement: WEB-OPS-013
      - type: web_field_equals
        view: network_node_detail
        label: Restart reason
        expected: Cannot verify operation status
        requirement: WEB-OPS-013
`);
  assert.equal(r.code,0,r.output);assert.equal(r.app.db.operations.filter(x=>x.ui).length,0);
});
