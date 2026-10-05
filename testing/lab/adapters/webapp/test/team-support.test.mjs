/* SPDX-License-Identifier: MPL-2.0 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {chromium} from 'playwright';
import {TeamSupport,trackTeamReads} from '../dist/team-support.js';
import {Budget} from '../dist/contract.js';
import {teamSupportFixture} from './team-support-fixture.mjs';
async function start(t,mode='',path='/business/manage/members'){
 const app=await teamSupportFixture(mode);t.after(()=>app.close());app.db.sites.push({id:'site-1',name:'Site one',network_id:'existing'});app.db.nodes.push({id:'node-1',type:'Tower node',site:'site-1'});
 const browser=await chromium.launch({executablePath:process.env.ULAB_WEBAPP_EXECUTABLE_PATH,headless:true,args:['--no-sandbox']});t.after(()=>browser.close());
 const context=await browser.newContext({serviceWorkers:'block'}),page=await context.newPage();trackTeamReads(page);await page.goto(app.origin+path);await page.locator('main.main .card, main.main table').first().waitFor();return {page,app,team:new TeamSupport(page)};
}
test('invitation probes reject normal recipient domains and unknown roles before submission',async t=>{
 const {team,app}=await start(t);for(const [role,email] of [['Administrator','person@example.com'],['Superuser','person@example.test']])await assert.rejects(team.run('business_members','invite_probe',role,email,new Budget(1000)),/reserved/);assert.equal(app.db.interceptedWrites.length,0);
});
test('invitation write barrier persists after a completed probe',async t=>{
 const {page,team,app}=await start(t);await team.run('business_members','invite_probe','Administrator','probe@example.test',new Budget(3000));
 const blocked=await page.evaluate(()=>fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'CreateInvitation',query:'mutation CreateInvitation($data: Input!) { createInvitation(data: $data) { id } }',variables:{data:{email:'late@example.test',role:'ROLE_OWNER',name:'late'}}})}).then(()=>false,()=>true));assert.equal(blocked,true);assert.equal(app.db.interceptedWrites.length,0);
});
test('restart probe keeps delayed writes blocked after confirmation is cancelled',async t=>{
 const {page,team,app}=await start(t,'late-restart','/network/support');await team.run('network_support','support_node','','node-1',new Budget(3000));
 const failed=page.waitForEvent('requestfailed',{predicate:r=>r.method()==='POST'&&r.postDataJSON()?.operationName==='RestartNode'});
 await team.run('network_support','support_restart_probe','','Restart node',new Budget(3000));await failed;assert.equal(app.db.interceptedWrites.length,0);
});
test('support node selection uses the displayed identifier even for custom names',async t=>{
 const {team}=await start(t,'named-node','/network/support');await team.run('network_support','support_node','','node-1',new Budget(3000));assert.equal(await team.observe('network_support','Support identity',''),'Custom node-1');assert.equal(await team.observe('network_support','Support field','Node ID'),'node-1');
});
test('support fault requires an observed exact network and does not consume foreign reads',async t=>{
 const {page,team}=await start(t,'','/network/support');await assert.rejects(team.run('network_support','support_failure','','',new Budget(2000),'foreign'),/observed scoped/);
 await team.run('network_support','support_failure','','',new Budget(3000),'existing');
 const foreign=await page.evaluate(async()=>{const r=await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({operationName:'SitesList',query:'query SitesList($networkId: String!) { sitesView { networkId } }',variables:{networkId:'foreign'}})});return r.json()});assert.deepEqual(foreign,{data:{fixture:true}});
 assert.equal(await team.observe('network_support','Support fault',''),'applied');await team.run('network_support','support_clear_fault','','',new Budget(3000));assert.equal(await team.observe('network_support','Support fault',''),'not applied');
});
test('support summary cannot be asserted without copying the selected entity',async t=>{
 const {team}=await start(t,'','/network/support');await assert.rejects(team.observe('network_support','Support summary','ICCID'),/requires copying/);
});
test('team/support commands reject a declared network different from the opened one',async t=>{
 const {page,app}=await start(t,'','/network/support');
 const {ConsoleApp}=await import('../dist/console-app.js'),{Interactions}=await import('../dist/interactions.js');
 const consoleApp=new ConsoleApp(page,app.origin);await consoleApp.open({view:'network_support',network_name:'existing-network'},new Budget(3000));
 const ui=new Interactions(page,app.origin,consoleApp);
 await assert.rejects(ui.run({view:'network_support',action:'support_node',value:'node-1',network_name:'foreign'},new Budget(1000)),/differs from the opened network/);
 await assert.rejects(ui.check({view:'network_support',label:'Support field',subject:'Node ID',expected:'node-1',requirement:'WEB-SUPPORT-002',network_name:'foreign'},new Budget(1000)),/differs from the opened network/);
});
