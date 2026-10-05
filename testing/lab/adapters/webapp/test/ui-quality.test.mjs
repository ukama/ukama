/* SPDX-License-Identifier: MPL-2.0 — controlled browser diagnostics/DOM checks. */
import test from 'node:test';
import assert from 'node:assert/strict';
import {mkdtemp,readFile,rm} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {join} from 'node:path';
import {chromium} from 'playwright';
import {UIQuality,trackQualityReads} from '../dist/ui-quality.js';
import {Artifacts} from '../dist/artifacts.js';
import {Budget} from '../dist/contract.js';
import {teamSupportFixture} from './team-support-fixture.mjs';
async function setup(t,mode=''){
 const app=await teamSupportFixture(mode);t.after(()=>app.close());
 const browser=await chromium.launch({executablePath:process.env.ULAB_WEBAPP_EXECUTABLE_PATH,headless:true});t.after(()=>browser.close());
 const page=await browser.newPage({serviceWorkers:'block'});trackQualityReads(page);await page.goto(app.origin+'/business/manage/members');await page.locator('tbody tr').first().waitFor();
 return {page,app,ui:new UIQuality(page)};
}
test('member loading and errors stay distinct; UI retry recovers after releasing exact query fault',async t=>{
 const {ui}=await setup(t),run=(a,v='')=>ui.run('business_members',a,v,new Budget(4000));
 await run('member_fault','loading');assert.equal(await ui.observe('Members state',''),'loading');await run('member_release');
 await new Budget(2000).poll(()=>ui.observe('Members state',''),v=>v==='data','data');
 await run('member_fault','error');await new Budget(2000).poll(()=>ui.observe('Members state',''),v=>v==='error','error');
 await assert.rejects(run('member_retry'),/Release/);await run('member_release');await run('member_retry');
 await new Budget(2000).poll(()=>ui.observe('Members state',''),v=>v==='data','recovery');
});
test('member faults reject missing/ambiguous endpoint and do not intercept another operation',async t=>{
 const {page,ui}=await setup(t);await page.evaluate(()=>fetch('/other-graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'TeamList',query:'query TeamList { x }',variables:{}})}));
 await assert.rejects(ui.run('business_members','member_fault','error',new Budget(1000)),/one observed/);
});
test('focus traversal proves both directions and cannot reuse proof after dialog replacement',async t=>{
 const {page,ui}=await setup(t);await page.getByRole('button',{name:'Invite member',exact:true}).click();await ui.run('business_members','focus_cycle','',new Budget(3000));
 assert.equal(await ui.observe('Dialog keyboard cycle',''),'true');await page.keyboard.press('Escape');
 await page.getByRole('button',{name:'Invite member',exact:true}).click();assert.equal(await ui.observe('Dialog keyboard cycle',''),'false');
});
test('focus check catches missing traps without repairing focus',async t=>{
 const {page,ui}=await setup(t,'untrapped');await page.getByRole('button',{name:'Invite member',exact:true}).click();
 await assert.rejects(ui.run('business_members','focus_cycle','',new Budget(3000)),/escaped/);
});
test('reachable action rejects an overlay and dialog fit rejects horizontal overflow',async t=>{
 const {page,ui}=await setup(t,'wide-dialog');await page.setViewportSize({width:390,height:844});
 await page.getByRole('button',{name:'Invite member',exact:true}).click();assert.equal(await ui.observe('Dialog fits viewport',''),'false');await page.keyboard.press('Escape');
 await page.evaluate(()=>{const d=document.createElement('div');d.style.cssText='position:fixed;inset:0;z-index:999';document.body.append(d)});
 assert.equal(await ui.observe('Action reachable','Invite member'),'false');
});
test('step diagnostics retain request origin, omit secrets and preserve health after flushing',async t=>{
 const {page,app}=await setup(t);const dir=await mkdtemp(join(tmpdir(),'ulab-diag-'));t.after(()=>rm(dir,{recursive:true,force:true}));
 const evidence=new Artifacts(dir,'run');await evidence.create();evidence.attach(page);evidence.begin({command_id:20,action:'web_reload'});
 let release;const ready=new Promise(resolve=>release=resolve);await page.route('**/critical?**',async r=>{release(r)});
 await page.evaluate(()=>{void fetch('/critical?token=never-record-me').catch(()=>{})});const route=await ready;
 evidence.begin({command_id:21,action:'web_ui_equals'});const response=page.waitForResponse(r=>r.url().includes('/critical?'));await route.fulfill({status:503,body:'secret'});await response;
 const error=page.waitForEvent('pageerror');await page.evaluate(()=>{setTimeout(()=>{throw Error('never-record-me')},0)});await error;
 assert.equal(evidence.health(page),'errors');const result={command_id:21,artifacts:[]};await evidence.record(result);
 const diagnostic=await readFile(result.artifacts[0],'utf8');assert(!diagnostic.includes('never-record-me'));const events=JSON.parse(diagnostic).events;
 assert(events.some(e=>e.command.command_id===21&&e.initiated_by?.command_id===20&&e.detail.status===503));
 assert(events.some(e=>e.detail.type==='pageerror'));assert.equal(evidence.health(page),'errors');
});
test('ambiguous loading plus rows is not an empty/data success',async t=>{
 const {page,ui}=await setup(t);await page.locator('main').evaluate(e=>e.insertAdjacentHTML('beforeend','<div class="MuiSkeleton-root">loading</div>'));
 assert.equal(await ui.observe('Members state',''),'ambiguous');
});
