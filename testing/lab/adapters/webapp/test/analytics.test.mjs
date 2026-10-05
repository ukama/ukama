/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {chromium} from 'playwright';
import {Analytics} from '../dist/analytics.js';
import {Budget} from '../dist/contract.js';
import {Interactions} from '../dist/interactions.js';
import {chartMarkup} from './analytics-fixture.mjs';
async function start(t,html=''){
 const browser=await chromium.launch({executablePath:process.env.ULAB_WEBAPP_EXECUTABLE_PATH,headless:true,args:['--no-sandbox']});t.after(()=>browser.close());
 const page=await browser.newPage();await page.setContent('<main class="main">'+html+'</main>');return {page,a:new Analytics(page)};
}
test('analytics refuses an absent tooltip without a hover observation',async t=>{
 const {a}=await start(t,chartMarkup('Temperature'));await assert.rejects(a.observe('Chart tooltip','Temperature'),/requires a hover/);
});
test('analytics reads rendered SVG labels and gap subpaths',async t=>{
 const {a}=await start(t,chartMarkup('Temperature'));assert.equal(await a.observe('Chart x axis','Temperature'),'12:00 PM|12:30 PM|1:00 PM');assert.equal(await a.observe('Chart y axis','Temperature'),'0°C|80°C');assert.equal(await a.observe('Chart segments','Temperature'),'2');
});
test('analytics distinguishes missing cards, empty series, errors and hidden plots',async t=>{
 const {a,page}=await start(t,chartMarkup('Memory')+chartMarkup('Backhaul downlink')+chartMarkup('Temperature'));
 assert.equal(await a.observe('Chart state','Missing'),null);assert.equal(await a.observe('Chart state','Memory'),'empty');assert.equal(await a.observe('Chart state','Backhaul downlink'),'error');
 await page.locator('svg').evaluate(e=>e.style.display='none');assert.equal(await a.observe('Chart state','Temperature'),null);
});
test('analytics rejects ambiguous cards instead of reading the first chart',async t=>{
 const {a}=await start(t,chartMarkup('Temperature')+chartMarkup('Temperature'));await assert.rejects(a.observe('Chart state','Temperature'),/multiple visible/);
});
test('analytics rejects unsupported hover/range and wrong-view access before interaction',async t=>{
 const {a,page}=await start(t,chartMarkup('Temperature'));
 for(const value of ['-1','101','NaN','1e2'])await assert.rejects(a.run('chart_hover','Temperature',value,new Budget(1000)),/percentage/);
 await assert.rejects(a.run('chart_range','Temperature','Year',new Budget(1000)),/range/);
 const app={assertView(){throw new Error('wrong current network')}};
 await assert.rejects(new Interactions(page,'http://fixture.test',app).check({view:'network_node_detail',label:'Chart state',subject:'Temperature',expected:'data',requirement:'WEB-UI-004'},new Budget(1000)),/wrong current network/);
});
